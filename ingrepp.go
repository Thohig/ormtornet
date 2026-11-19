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
	errNotGM      = hostErr("server_bara_gm")
	errNotPlaying = hostErr("server_bara_under_vandring")
	errNoPlayer   = hostErr("server_ingen_vandrare")
	errNoRoom     = hostErr("server_valj_rum")
	errUnknown    = hostErr("server_okant_ingrepp")
)

// gmNote skriver ett ingrepp i spelledarens logg.
// Texten står i data/vardsidan.json under key.
func (g *Game) gmNote(key string, vars ...string) {
	text := g.data.HostText(key, vars...)
	g.gmLog = append(g.gmLog, gmEntry{Round: g.round, Text: text, Class: "gm", To: []string{g.data.HostText("server_spelledaren")}})
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
			return hostErr("server_pausad")
		}
		g.gmNote("server_not_avgjorde")
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
		g.gmNote("server_not_vackte")
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
	case "invig":
		err = g.gmInitiate(r.Spelare)
	case "foremal":
		err = g.gmGiveItem(r.Spelare, r.ID)
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
	g.gmNote("server_not_pausade")
	return nil
}

func (g *Game) gmResume() error {
	if !g.paused {
		return nil
	}
	g.paused, g.awaiting = false, false
	g.armTimer(g.pauseLeft)
	g.logAll(g.data.Text("gm_fortsatt"), "rost")
	g.gmNote("server_not_fortsatte")
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
	g.gmNote("server_not_mer_tid", "sekunder", strconv.Itoa(int(extra/time.Second)))
	return nil
}

// gmSun flyttar solen. Den går inte under 1, så att det är rundan som får
// solen att gå ner, och inte över sin startpunkt.
func (g *Game) gmSun(delta int) error {
	if delta == 0 {
		return nil
	}
	before := g.sun
	g.sun = max(1, min(g.sunMax, g.sun+delta))
	if g.sun == before {
		return nil
	}
	if g.sun < before {
		g.logAll(g.data.Text("gm_sol_ner"), "rost")
		g.sunWarning(before)
		g.updateLevel()
		g.gmNote("server_not_sol_ner", "sol", strconv.Itoa(g.sun))
	} else {
		g.logAll(g.data.Text("gm_sol_upp"), "rost")
		g.gmNote("server_not_sol_upp", "sol", strconv.Itoa(g.sun))
	}
	return nil
}

// gmVoice: tornets röst, till en vandrare eller alla.
func (g *Game) gmVoice(to int, text string) error {
	text = strings.TrimSpace(strings.ReplaceAll(text, "\n", " "))
	if text == "" {
		return hostErr("server_rost_tom")
	}
	if utf8.RuneCountInString(text) > voiceMaxLen {
		return hostErr("server_rost_lang", "max", strconv.Itoa(voiceMaxLen))
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
		return hostErr("server_okand_falla")
	}
	if r.Event != "" {
		return hostErr("server_falla_finns")
	}
	r.Event = id
	r.Mark = ""
	g.gmNote("server_not_falla", "falla", lowerFirst(t.Namn), "rum", g.roomName(r))
	return nil
}

func (g *Game) gmMonster(x, y int, id string) error {
	if _, err := g.gmRoom(x, y); err != nil {
		return err
	}
	t := g.data.monsterType(id)
	if t == nil {
		return hostErr("server_okand_varelse")
	}
	if g.monsterAt(x, y) != nil {
		return hostErr("server_varelse_finns")
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
	g.gmNote("server_not_varelse", "varelse", lowerFirst(t.Namn), "rum", g.roomName(&g.rooms[y][x]))
	return nil
}

// gmCalm tar bort ett vakna-kort ur ormens lek, eller lägger till ett
// sovkort om det inte finns något att ta bort.
func (g *Game) gmCalm() error {
	switch g.calmSerpent() {
	case "vakna":
		g.gmNote("server_not_lugna_vakna")
	case "sov":
		g.gmNote("server_not_lugna_sov")
	}
	return nil
}

var errNoReveal = hostErr("server_avslojande_av")

// gmLevelUp flyttar fram avslöjandet en nivå, med allt som hör till: ljud,
// nya texter och ledtrådar, och förklaringar på nivå 4.
func (g *Game) gmLevelUp() error {
	if !g.cfg.Reveal {
		return errNoReveal
	}
	if g.level >= 4 {
		return hostErr("server_avslojande_sista")
	}
	g.level++
	g.onLevelUp()
	g.gmNote("server_not_niva_upp", "niva", strconv.Itoa(g.level))
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
		g.gmNote("server_not_niva_hall", "niva", strconv.Itoa(g.level))
	} else {
		g.gmNote("server_not_niva_slapp")
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
		g.gmNote(verb(n, "server_not_gav_liv", "server_not_tog_liv"), "antal", strconv.Itoa(abs(n)), "namn", p.Name)
	case "guld":
		before := p.Treasure
		p.Treasure = max(0, p.Treasure+n)
		g.logTo(p, g.data.Text("gm_guld_"+sign), "rost")
		g.gmNote(verb(n, "server_not_gav_guld", "server_not_tog_guld"), "antal", strconv.Itoa(abs(p.Treasure-before)), "namn", p.Name)
	case "dryck":
		before := p.Potions
		p.Potions = max(0, p.Potions+n)
		if p.Potions == before {
			return hostErr("server_inga_drycker", "namn", p.Name)
		}
		g.logTo(p, g.data.Text("gm_dryck_"+sign), "rost")
		g.gmNote(verb(n, "server_not_gav_dryck", "server_not_tog_dryck"), "namn", p.Name)
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

// gmInitiate invigar en vandrare, som om den hade viskat hemligheten, eller
// tar tillbaka det.
func (g *Game) gmInitiate(id int) error {
	p := g.playerByID(id)
	if p == nil {
		return errNoPlayer
	}
	if !g.cfg.Reveal {
		return errNoReveal
	}
	p.initiated = !p.initiated
	if p.initiated {
		g.logTo(p, g.data.Text("invigd"), "rost")
		g.gmNote("server_not_invigde", "namn", p.Name)
	} else {
		g.gmNote("server_not_avinvigde", "namn", p.Name)
	}
	return nil
}

// gmGiveItem lägger ett föremål i en vandrares väska.
func (g *Game) gmGiveItem(id int, item string) error {
	p := g.playerByID(id)
	if p == nil {
		return errNoPlayer
	}
	if p.Status != statusIn {
		return errors.New(p.Name + " är inte kvar i tornet.")
	}
	t := g.data.itemType(item)
	if t == nil || t.Brev {
		return hostErr("server_foremal_finns_inte")
	}
	if !g.give(p, item) {
		return hostErr("server_bar_inte_mer", "namn", p.Name)
	}
	g.logTo(p, g.data.Text("foremal_fran_tornet", "foremal", lowerFirst(t.Namn), "beskrivning", t.Beskrivning), "rost")
	g.gmNote("server_not_gav_foremal", "foremal", lowerFirst(t.Namn), "namn", p.Name)
	return nil
}
