package main

// Spelledarens ingrepp. De går bara i spelledarläget, och varje ingrepp
// skrivs in i spelledarens logg.

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// gmRequest är ett ingrepp från värdsidan.
type gmRequest struct {
	Typ     string `json:"typ"`
	Spelare int    `json:"spelare"` // vandrarens id, 0 = alla (för tornets röst)
	X       int    `json:"x"`
	Y       int    `json:"y"`
	ID      string `json:"id"`    // fällans eller varelsens id
	Text    string `json:"text"`  // tornets röst
	Antal   int    `json:"antal"` // hur mycket som ges (positivt) eller tas (negativt)
}

const voiceMaxLen = 200

var (
	errNotGM      = errors.New("Ingrepp går bara i spelledarläget.")
	errNotPlaying = errors.New("Det går bara medan en vandring pågår.")
	errNoPlayer   = errors.New("Den vandraren finns inte i tornet.")
	errNoRoom     = errors.New("Välj ett vanligt rum som någon redan har varit i.")
	errUnknown    = errors.New("Okänt ingrepp.")
)

// gmNote skriver ett ingrepp i spelledarens logg.
func (g *Game) gmNote(text string) {
	g.gmLog = append(g.gmLog, gmEntry{Round: g.round, Text: text, Class: "gm", To: []string{"Spelledaren"}})
	if len(g.gmLog) > gmLogMax {
		g.gmLog = append([]gmEntry(nil), g.gmLog[len(g.gmLog)-gmLogMax:]...)
	}
}

// Intervene utför ett ingrepp och berättar för alla vad de märker av det.
func (g *Game) Intervene(r gmRequest) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.cfg.GMMode {
		return errNotGM
	}
	if r.Typ != "rost" && g.phase != phasePlay {
		return errNotPlaying
	}
	var err error
	switch r.Typ {
	case "paus":
		err = g.gmPause()
	case "fortsatt":
		err = g.gmResume()
	case "mer_tid":
		err = g.gmMoreTime()
	case "avgor":
		if g.paused {
			return errors.New("Fortsätt först, rundan är pausad.")
		}
		g.gmNote("Avgjorde rundan direkt.")
		g.resolve()
		return nil
	case "sol":
		err = g.gmSun(r.Antal)
	case "rost":
		err = g.gmVoice(r.Spelare, r.Text)
	case "falla":
		err = g.gmTrap(r.X, r.Y, r.ID)
	case "varelse":
		err = g.gmMonster(r.X, r.Y, r.ID)
	case "vack":
		g.gmNote("Väckte ormen.")
		g.chronicle("kronika_vaknade_sjalv")
		g.wakeSerpent()
	case "lugna":
		err = g.gmCalm()
	case "niva_upp":
		err = g.gmLevelUp()
	case "niva_hall":
		err = g.gmHoldLevel(true)
	case "niva_slapp":
		err = g.gmHoldLevel(false)
	case "liv", "guld", "dryck":
		err = g.gmGive(r.Typ, r.Spelare, r.Antal)
	default:
		return errUnknown
	}
	if err != nil {
		return err
	}
	g.broadcastState()
	return nil
}

func (g *Game) gmPause() error {
	if g.paused {
		return nil
	}
	g.paused = true
	g.pauseLeft = max(time.Second, time.Until(g.deadline))
	if g.timer != nil {
		g.timer.Stop()
	}
	g.logAll(g.data.Text("gm_paus"), "rost")
	g.gmNote("Pausade rundan.")
	return nil
}

func (g *Game) gmResume() error {
	if !g.paused {
		return nil
	}
	g.paused = false
	g.armTimer(g.pauseLeft)
	g.logAll(g.data.Text("gm_fortsatt"), "rost")
	g.gmNote("Lät rundan fortsätta.")
	if g.allChosen() {
		g.resolve()
	}
	return nil
}

func (g *Game) gmMoreTime() error {
	const extra = 30 * time.Second
	if g.paused {
		g.pauseLeft += extra
	} else {
		g.armTimer(time.Until(g.deadline) + extra)
	}
	g.gmNote("Gav rundan 30 sekunder till.")
	return nil
}

// gmSun flyttar solen. Den går inte under 1, så att det är rundan som får
// solen att gå ner, och inte över sin startpunkt.
func (g *Game) gmSun(delta int) error {
	if delta == 0 {
		return nil
	}
	before := g.sun
	g.sun = max(1, min(g.cfg.SunMax, g.sun+delta))
	if g.sun == before {
		return nil
	}
	if g.sun < before {
		g.logAll(g.data.Text("gm_sol_ner"), "rost")
		g.sunWarning(before)
		g.updateLevel()
		g.gmNote("Lät solen sjunka till " + strconv.Itoa(g.sun) + ".")
	} else {
		g.logAll(g.data.Text("gm_sol_upp"), "rost")
		g.gmNote("Lyfte solen till " + strconv.Itoa(g.sun) + ".")
	}
	return nil
}

// gmVoice: tornets röst, till en vandrare eller alla.
func (g *Game) gmVoice(to int, text string) error {
	text = strings.TrimSpace(strings.ReplaceAll(text, "\n", " "))
	if text == "" {
		return errors.New("Skriv något för tornet att säga.")
	}
	if utf8.RuneCountInString(text) > voiceMaxLen {
		return errors.New("Tornets röst får vara högst " + strconv.Itoa(voiceMaxLen) + " tecken.")
	}
	if to == 0 {
		g.logAll(text, "rost")
		g.soundAll("rost")
		return nil
	}
	p := g.playerByID(to)
	if p == nil {
		return errNoPlayer
	}
	g.logTo(p, text, "rost")
	g.soundTo(p, "rost")
	return nil
}

// gmRoom ger ett vanligt rum som redan finns i tornet.
func (g *Game) gmRoom(x, y int) (*Room, error) {
	if !g.inBounds(x, y) {
		return nil, errNoRoom
	}
	r := &g.rooms[y][x]
	if !r.Revealed || r.Kind != "rum" {
		return nil, errNoRoom
	}
	return r, nil
}

func (g *Game) gmTrap(x, y int, id string) error {
	r, err := g.gmRoom(x, y)
	if err != nil {
		return err
	}
	t := g.data.trap(id)
	if t == nil {
		return errors.New("Okänd fälla.")
	}
	if r.Event != "" {
		return errors.New("Det väntar redan en fälla i det rummet.")
	}
	r.Event = id
	r.Mark = ""
	g.gmNote("Lade " + lowerFirst(t.Namn) + " i " + g.roomName(r) + ".")
	return nil
}

func (g *Game) gmMonster(x, y int, id string) error {
	if _, err := g.gmRoom(x, y); err != nil {
		return err
	}
	t := g.data.monsterType(id)
	if t == nil {
		return errors.New("Okänd varelse.")
	}
	if g.monsterAt(x, y) != nil {
		return errors.New("Det finns redan en varelse i det rummet.")
	}
	g.nextMonsterID++
	m := &Monster{ID: g.nextMonsterID, Type: t.ID, X: x, Y: y, HP: t.Liv}
	g.monsters = append(g.monsters, m)
	for _, p := range g.players {
		if p != nil && p.Status == statusIn && p.Cellar == 0 && p.X == x && p.Y == y {
			g.soundTo(p, "monster")
			g.logTo(p, g.data.MonsterText("monster_kommer", g.monsterVars(m)...), "varn")
		}
	}
	g.gmNote("Släppte in " + lowerFirst(t.Namn) + " i " + g.roomName(&g.rooms[y][x]) + ".")
	return nil
}

// gmCalm tar bort ett vakna-kort ur ormens lek, eller lägger till ett
// sovkort om det inte finns något att ta bort.
func (g *Game) gmCalm() error {
	for i, id := range g.deck {
		if c := g.data.card(id); c != nil && c.Vaknar {
			g.deck = append(g.deck[:i], g.deck[i+1:]...)
			g.gmNote("Lugnade ormen: ett vakna-kort mindre.")
			return nil
		}
	}
	for _, c := range g.data.Ormen.Kort {
		if !c.Vaknar {
			g.deck = append(g.deck, c.ID)
			g.gmNote("Lugnade ormen: ett sovkort till.")
			return nil
		}
	}
	return nil
}

var errNoReveal = errors.New("Avslöjandet är avstängt.")

// gmLevelUp flyttar fram avslöjandet en nivå, med allt som hör till: ljud,
// nya texter och ledtrådar, och förklaringar på nivå 4.
func (g *Game) gmLevelUp() error {
	if !g.cfg.Reveal {
		return errNoReveal
	}
	if g.level >= 4 {
		return errors.New("Avslöjandet är redan på sista nivån.")
	}
	g.level++
	g.onLevelUp()
	g.gmNote("Flyttade fram avslöjandet till nivå " + strconv.Itoa(g.level) + ".")
	return nil
}

// gmHoldLevel håller kvar avslöjandets nivå, eller släpper den så att den
// följer solen igen. Släpps den när solen redan har sjunkit förbi en gräns
// hinner nivån ikapp efter rundan.
func (g *Game) gmHoldLevel(hold bool) error {
	if !g.cfg.Reveal {
		return errNoReveal
	}
	if g.levelHeld == hold {
		return nil
	}
	g.levelHeld = hold
	if hold {
		g.gmNote("Höll kvar avslöjandet på nivå " + strconv.Itoa(g.level) + ".")
	} else {
		g.gmNote("Släppte avslöjandet, så att det följer solen igen.")
	}
	return nil
}

// gmGive ger eller tar liv, guld eller drycker. Vandraren märker att något
// händer, men inte varifrån.
func (g *Game) gmGive(what string, id, n int) error {
	p := g.playerByID(id)
	if p == nil {
		return errNoPlayer
	}
	if n == 0 {
		return nil
	}
	if p.Status != statusIn {
		return errors.New(p.Name + " är inte kvar i tornet.")
	}
	sign := "mer"
	if n < 0 {
		sign = "mindre"
	}
	switch what {
	case "liv":
		if n > 0 {
			p.HP = min(g.maxHP(p), p.HP+n)
		}
		g.logTo(p, g.data.Text("gm_liv_"+sign, "liv", strconv.Itoa(max(0, p.HP+min(0, n)))), "rost")
		if n < 0 {
			p.hurt("tornet", "")
			g.damage(p, -n, false)
		}
		g.gmNote(verb(n, "Gav", "Tog") + " " + strconv.Itoa(abs(n)) + " liv " + verb(n, "till", "från") + " " + p.Name + ".")
	case "guld":
		before := p.Treasure
		p.Treasure = max(0, p.Treasure+n)
		g.logTo(p, g.data.Text("gm_guld_"+sign), "rost")
		g.gmNote(verb(n, "Gav", "Tog") + " " + strconv.Itoa(abs(p.Treasure-before)) + " guld " + verb(n, "till", "från") + " " + p.Name + ".")
	case "dryck":
		before := p.Potions
		p.Potions = max(0, p.Potions+n)
		if p.Potions == before {
			return errors.New(p.Name + " har inga drycker att ta.")
		}
		g.logTo(p, g.data.Text("gm_dryck_"+sign), "rost")
		g.gmNote(verb(n, "Gav", "Tog") + " en dryck " + verb(n, "till", "från") + " " + p.Name + ".")
	}
	return nil
}

func verb(n int, plus, minus string) string {
	if n < 0 {
		return minus
	}
	return plus
}

// gmCatalog är vad spelledaren kan lägga i ett rum.
type gmCatalog struct {
	ID   string `json:"id"`
	Namn string `json:"namn"`
}

func (g *Game) gmCatalogs() (traps, monsters []gmCatalog) {
	for _, t := range g.data.Fallor {
		traps = append(traps, gmCatalog{t.ID, t.Namn})
	}
	for _, t := range g.data.Monster.Typer {
		monsters = append(monsters, gmCatalog{t.ID, t.Namn})
	}
	return traps, monsters
}
