package main

// Klasser: varje vandrare kan vara något, med en egen styrka och en svaghet.
// Namn, texter och siffror står i data/klasser.json. Här står vad varje
// klass gör, styrt av dess id, och hur man väljer klass efter lobbyn.

import (
	"slices"
	"strconv"
	"time"
)

type Class struct {
	ID          string             `json:"id"`
	Namn        string             `json:"namn"`
	Kort        string             `json:"kort"`   // med liten bokstav, för listan och texterna
	Tecken      string             `json:"tecken"` // klassens tecken på kartan
	Bild        string             `json:"bild"`   // klasskortets motiv (web/kort.js); utan bild gäller id:t
	Beskrivning string             `json:"beskrivning"`
	Styrka      string             `json:"styrka"`
	Svaghet     string             `json:"svaghet"`
	Liv         int                `json:"liv"` // läggs till svårighetens antal liv
	Drycker     int                `json:"drycker"`
	Kraft       int                `json:"kraft"`   // hur lätt vandraren träffar i strid
	Forsvar     int                `json:"forsvar"` // hur svår vandraren är att träffa
	Varden      map[string]float64 `json:"varden"`
	Modell      *ModelLook         `json:"modell,omitempty"` // färgen på det klassen bär och storleken, på kartan snett ovanifrån
}

// classView är det klienten får veta om en klass.
type classView struct {
	ID          string `json:"id"`
	Namn        string `json:"namn"`
	Tecken      string `json:"tecken,omitempty"`
	Bild        string `json:"bild"`
	Beskrivning string `json:"beskrivning"`
	Styrka      string `json:"styrka"`
	Svaghet     string `json:"svaghet"`
	Liv         int    `json:"liv,omitempty"`
	Drycker     int    `json:"drycker,omitempty"`
	Tagen       string `json:"tagen,omitempty"` // vem som har valt klassen, under klassvalet
	Din         bool   `json:"din,omitempty"`   // klassen är den egna
}

func viewOfClass(c *Class) *classView {
	if c == nil {
		return nil
	}
	return &classView{ID: c.ID, Namn: c.Namn, Tecken: c.Tecken, Bild: cardImage(c.Bild, c.ID), Beskrivning: c.Beskrivning, Styrka: c.Styrka, Svaghet: c.Svaghet, Liv: c.Liv, Drycker: c.Drycker}
}

func (d *Data) class(id string) *Class {
	for i := range d.Klasser {
		if d.Klasser[i].ID == id {
			return &d.Klasser[i]
		}
	}
	return nil
}

// classOf ger vandrarens klass, eller nil när klasserna är avstängda.
func (g *Game) classOf(p *Player) *Class {
	if p == nil || !g.cfg.Classes || p.classID == "" {
		return nil
	}
	return g.data.class(p.classID)
}

// is är sant när vandraren har klassen id.
func (g *Game) is(p *Player, id string) bool {
	c := g.classOf(p)
	return c != nil && c.ID == id
}

// classVal ger ett värde ur klassens "varden", eller def.
func (g *Game) classVal(p *Player, key string, def float64) float64 {
	if c := g.classOf(p); c != nil {
		if v, ok := c.Varden[key]; ok {
			return v
		}
	}
	return def
}

// maxHP är vandrarens liv: svårighetens liv och klassens tillägg.
func (g *Game) maxHP(p *Player) int {
	n := g.data.Regler.Liv
	if c := g.classOf(p); c != nil {
		n += c.Liv
	}
	return max(1, n)
}

func (g *Game) startPotions(p *Player) int {
	if c := g.classOf(p); c != nil {
		return c.Drycker
	}
	return 0
}

// takenClasses är klasserna som andra redan har valt, med vem som valde.
func (g *Game) takenClasses(p *Player) map[string]*Player {
	taken := map[string]*Player{}
	for _, o := range g.players {
		if o != nil && o != p && o.classID != "" {
			taken[o.classID] = o
		}
	}
	return taken
}

// freeClasses är klasserna som ingen annan än p har valt.
func (g *Game) freeClasses(p *Player) []string {
	taken := g.takenClasses(p)
	var free []string
	for _, c := range g.data.Klasser {
		if taken[c.ID] == nil {
			free = append(free, c.ID)
		}
	}
	return free
}

// Klassvalet: när alla är redo i lobbyn lägger tornet fram ett kort för
// varje klass. Var och en väljer ett kort som ingen annan har tagit.
// Människorna väljer först; AI-spelarna tar sina kort när alla människor har
// valt. Den som inte har valt när tiden är slut får ett av de lediga korten,
// och vandringen börjar när alla har en klass.
func (g *Game) startChoosing() {
	g.phase = phaseChoose
	g.resetClasses()
	g.chooseGen++
	gen := g.chooseGen
	wait := time.Duration(g.data.KlassvalSekunder) * time.Second
	if wait <= 0 {
		wait = 30 * time.Second
	}
	g.deadline = time.Now().Add(wait)
	if g.timer != nil {
		g.timer.Stop()
	}
	g.timer = time.AfterFunc(wait, func() {
		g.mu.Lock()
		defer g.mu.Unlock()
		if g.phase == phaseChoose && g.chooseGen == gen {
			g.startGame()
		}
	})
	for _, p := range g.players {
		if p == nil {
			continue
		}
		if p.bot == nil && p.client != nil {
			g.logTo(p, g.data.Text("klass_val"), "info")
		}
	}
	g.soundAll("val")
	if !g.checkChosen() {
		g.broadcastState()
	}
}

// checkChosen startar vandringen när alla som är med har valt klass, också
// den som har tappat anslutningen under klassvalet: hen väntar tornet på
// tills tiden är slut. Svarar true om den startade.
func (g *Game) checkChosen() bool {
	if g.phase != phaseChoose {
		return false
	}
	g.botsChoose()
	for _, p := range g.players {
		if p != nil && p.classID == "" && g.cfg.Classes {
			return false
		}
	}
	g.startGame()
	return true
}

// ChooseClass tar emot spelarens val av klass under klassvalet.
func (g *Game) ChooseClass(c *Client, id string) {
	p := c.player
	if p == nil || g.phase != phaseChoose || p.classID != "" {
		return
	}
	if g.chooseClass(p, id) {
		g.soundTo(p, "val")
	}
	if !g.checkChosen() {
		g.broadcastState()
	}
}

func (g *Game) chooseClass(p *Player, id string) bool {
	if !g.cfg.Classes {
		return false
	}
	cl := g.data.class(id)
	if cl == nil {
		return false
	}
	if g.takenClasses(p)[id] != nil {
		g.logTo(p, g.data.Text("klass_upptagen", "klass", cl.Kort), "varn")
		return false
	}
	p.classID = id
	p.HP = g.maxHP(p)
	g.logTo(p, g.data.Text("klass_vald", "klass", cl.Kort, "styrka", cl.Styrka), "syn")
	g.logOthers(p, g.data.Text("klass_vald_andra", "namn", p.Name, "klass", cl.Kort), "dim")
	return true
}

// botsChoose låter AI-spelarna välja klass när alla människor som är med
// har valt, så att människorna får välja först.
func (g *Game) botsChoose() {
	for _, p := range g.players {
		if p != nil && p.bot == nil && p.client != nil && p.classID == "" {
			return
		}
	}
	for _, p := range g.players {
		if p != nil && p.bot != nil {
			g.botChooseClass(p)
		}
	}
}

// botChooseClass: AI-spelaren väljer den lediga klass dess personlighet
// helst vill ha, annars någon av de lediga.
func (g *Game) botChooseClass(p *Player) {
	if !g.cfg.Classes || p.classID != "" {
		return
	}
	free := g.freeClasses(p)
	if len(free) == 0 {
		return
	}
	pick := free[g.rng.IntN(len(free))]
	if p.bot != nil && p.bot.pers != nil {
		for _, want := range p.bot.pers.Klasser {
			if slices.Contains(free, want) {
				pick = want
				break
			}
		}
	}
	g.chooseClass(p, pick)
}

// ensureClasses ger alla som inte har valt en klass när vandringen börjar
// ett av de lediga korten. Människorna får sina först, sedan väljer
// AI-spelarna bland det som är kvar.
func (g *Game) ensureClasses() {
	if !g.cfg.Classes {
		return
	}
	for _, p := range g.players {
		if p == nil || p.classID != "" || p.bot != nil {
			continue
		}
		free := g.freeClasses(p)
		if len(free) == 0 {
			continue
		}
		if p.bot == nil {
			g.logTo(p, g.data.Text("klass_slump"), "info")
		}
		g.chooseClass(p, free[g.rng.IntN(len(free))])
	}
	g.botsChoose()
}

// resetClasses tar bort allas klasser inför ett nytt klassval.
func (g *Game) resetClasses() {
	for _, p := range g.players {
		if p != nil {
			p.classID = ""
		}
	}
}

// classCards är korten under klassvalet: alla klasser, och vem som har tagit
// vilken.
func (g *Game) classCards(me *Player) []*classView {
	var cards []*classView
	for i := range g.data.Klasser {
		c := &g.data.Klasser[i]
		v := viewOfClass(c)
		for _, o := range g.players {
			if o != nil && o.classID == c.ID {
				v.Tagen = o.Name
				v.Din = o == me
			}
		}
		cards = append(cards, v)
	}
	return cards
}

// classLabel är klassen i listan över vandrare, till exempel "kartritaren".
func (g *Game) classLabel(p *Player) string {
	if c := g.classOf(p); c != nil {
		return c.Kort
	}
	return ""
}

// classModel ger färgen och storleken på vandrarens figur, eller nil.
func (g *Game) classModel(p *Player) *ModelLook {
	if c := g.classOf(p); c != nil {
		return c.Modell
	}
	return nil
}

// classID ger id för vandrarens klass, eller tomt.
func (g *Game) classID(p *Player) string {
	if c := g.classOf(p); c != nil {
		return c.ID
	}
	return ""
}

// wakeCardsLeft räknar vakna-korten i ormleken, för ormtjusaren.
func (g *Game) wakeCardsLeft() int {
	n := 0
	for _, id := range g.deck {
		if c := g.data.card(id); c != nil && c.Vaknar {
			n++
		}
	}
	return n
}

// shareMap: kartritaren ger sin karta till den hen möter. Med egen karta
// får den andra alla rum som någon har hittat.
func (g *Game) shareMap(p *Player, foes []*Player) {
	for _, o := range foes {
		if !g.is(o, "kartritare") || g.is(p, "kartritare") {
			continue
		}
		fresh := 0
		for y := range g.size {
			for x := range g.size {
				if g.knowsRoom(o, x, y) && !g.knowsRoom(p, x, y) {
					p.known[y][x] = true
					fresh++
				}
			}
		}
		// Bara när kartan ger något nytt, så att det inte upprepas varje
		// runda som de står tillsammans.
		if fresh == 0 {
			continue
		}
		g.logTo(p, g.data.Text("klass_kartritare_karta", "namn", o.Name), "info")
		g.logTo(o, g.data.Text("klass_kartritare_ger", "namn", p.Name), "dim")
	}
}

// peekNeighbors: kartritaren ser genom dörrarna in i rummen bredvid det
// hen står i, och de kommer med på hens karta. Rummen räknas inte som
// hittade, så ingen annan ser dem, och en fälla där väntar fortfarande.
func (g *Game) peekNeighbors(p *Player) {
	if g.phase != phasePlay || !g.is(p, "kartritare") || p.Dark > 0 {
		return
	}
	r := &g.rooms[p.Y][p.X]
	for d := North; d <= West; d++ {
		nx, ny := p.X+dirDX[d], p.Y+dirDY[d]
		if r.Doors[d] && g.inBounds(nx, ny) {
			g.reveal(nx, ny)
			p.known[ny][nx] = true
		}
	}
}

// nearLantern är sant när o är lyktbäraren och står i ett grannrum till me.
func (g *Game) nearLantern(me, o *Player) bool {
	if !(g.is(o, "lyktbarare") || o.has("rostkrona") != nil) || o.Status != statusIn || o.Cellar > 0 {
		return false
	}
	_, ok := g.nextTo(me, o.X, o.Y)
	return ok
}

func itoa(n int) string { return strconv.Itoa(n) }
