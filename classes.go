package main

// Klasser: varje vandrare kan vara något, med en egen styrka och en svaghet.
// Namn, texter och siffror står i data/klasser.json. Här står vad varje
// klass gör, styrt av dess id, och hur man väljer klass i lobbyn.

import (
	"math/rand/v2"
	"slices"
	"strconv"
)

type Class struct {
	ID          string             `json:"id"`
	Namn        string             `json:"namn"`
	Kort        string             `json:"kort"` // med liten bokstav, för listan och texterna
	Beskrivning string             `json:"beskrivning"`
	Styrka      string             `json:"styrka"`
	Svaghet     string             `json:"svaghet"`
	Liv         int                `json:"liv"` // läggs till svårighetens antal liv
	Drycker     int                `json:"drycker"`
	Varden      map[string]float64 `json:"varden"`
}

// classView är det klienten får veta om en klass.
type classView struct {
	ID          string `json:"id"`
	Namn        string `json:"namn"`
	Beskrivning string `json:"beskrivning"`
	Styrka      string `json:"styrka"`
	Svaghet     string `json:"svaghet"`
}

func viewOfClass(c *Class) *classView {
	if c == nil {
		return nil
	}
	return &classView{ID: c.ID, Namn: c.Namn, Beskrivning: c.Beskrivning, Styrka: c.Styrka, Svaghet: c.Svaghet}
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

// takenClasses är klasserna som andra redan har valt.
func (g *Game) takenClasses(except *Player) map[string]bool {
	taken := map[string]bool{}
	for _, o := range g.players {
		if o != nil && o != except && o.classID != "" {
			taken[o.classID] = true
		}
	}
	return taken
}

// offerClasses drar två klasser som ingen annan har valt eller fått att
// välja mellan, om det går.
func (g *Game) offerClasses(p *Player) {
	if !g.cfg.Classes || p.classID != "" || len(g.data.Klasser) == 0 {
		return
	}
	taken := g.takenClasses(p)
	offered := map[string]bool{}
	for _, o := range g.players {
		if o != nil && o != p {
			for _, id := range o.offer {
				offered[id] = true
			}
		}
	}
	var fresh, free []string
	for _, c := range g.data.Klasser {
		if taken[c.ID] {
			continue
		}
		free = append(free, c.ID)
		if !offered[c.ID] {
			fresh = append(fresh, c.ID)
		}
	}
	rand.Shuffle(len(fresh), func(i, j int) { fresh[i], fresh[j] = fresh[j], fresh[i] })
	rand.Shuffle(len(free), func(i, j int) { free[i], free[j] = free[j], free[i] })
	p.offer = nil
	for _, id := range append(fresh, free...) {
		if len(p.offer) == 2 {
			break
		}
		if !slices.Contains(p.offer, id) {
			p.offer = append(p.offer, id)
		}
	}
	if p.bot != nil {
		g.botChooseClass(p)
		return
	}
	if len(p.offer) > 0 {
		g.logTo(p, g.data.Text("klass_val"), "info")
	}
}

// ChooseClass tar emot spelarens val av klass i lobbyn.
func (g *Game) ChooseClass(c *Client, id string) {
	if c.player == nil || g.phase == phasePlay {
		return
	}
	g.chooseClass(c.player, id)
	g.broadcastState()
}

func (g *Game) chooseClass(p *Player, id string) bool {
	if !g.cfg.Classes || !slices.Contains(p.offer, id) {
		return false
	}
	cl := g.data.class(id)
	if cl == nil {
		return false
	}
	if g.takenClasses(p)[id] {
		g.logTo(p, g.data.Text("klass_upptagen", "klass", cl.Kort), "varn")
		g.offerClasses(p)
		return false
	}
	p.classID = id
	p.offer = nil
	p.HP = g.maxHP(p)
	g.logTo(p, g.data.Text("klass_vald", "klass", cl.Kort, "styrka", cl.Styrka), "syn")
	g.logOthers(p, g.data.Text("klass_vald_andra", "namn", p.Name, "klass", cl.Kort), "dim")
	// Den som fick samma klass att välja på får dra nya.
	for _, o := range g.players {
		if o != nil && o != p && o.classID == "" && slices.Contains(o.offer, id) {
			g.offerClasses(o)
		}
	}
	return true
}

// botChooseClass: AI-spelaren väljer den klass dess personlighet helst vill ha.
func (g *Game) botChooseClass(p *Player) {
	if len(p.offer) == 0 {
		return
	}
	pick := p.offer[0]
	if p.bot != nil && p.bot.pers != nil {
		for _, want := range p.bot.pers.Klasser {
			if slices.Contains(p.offer, want) {
				pick = want
				break
			}
		}
	}
	g.chooseClass(p, pick)
}

// ensureClasses ger alla som inte har valt en klass när vandringen börjar.
func (g *Game) ensureClasses() {
	if !g.cfg.Classes {
		return
	}
	for _, p := range g.players {
		if p == nil || p.classID != "" {
			continue
		}
		if len(p.offer) == 0 {
			g.offerClasses(p)
		}
		if p.classID == "" && len(p.offer) > 0 {
			g.chooseClass(p, p.offer[rand.IntN(len(p.offer))])
		}
	}
}

// resetClasses tar bort allas klasser inför ett nytt torn.
func (g *Game) resetClasses() {
	for _, p := range g.players {
		if p != nil {
			p.classID, p.offer = "", nil
		}
	}
}

// offerAll låter alla som saknar klass välja, till exempel när spelledaren
// slår på klasserna i lobbyn.
func (g *Game) offerAll() {
	for _, p := range g.players {
		if p != nil && p.classID == "" && len(p.offer) == 0 {
			g.offerClasses(p)
		}
	}
}

// classLabel är klassen i listan över vandrare, till exempel "kartritaren".
func (g *Game) classLabel(p *Player) string {
	if c := g.classOf(p); c != nil {
		return c.Kort
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
	if g.cfg.Sight != sightOwn {
		return
	}
	for _, o := range foes {
		if !g.is(o, "kartritare") || g.is(p, "kartritare") {
			continue
		}
		for y := range g.size {
			for x := range g.size {
				if g.rooms[y][x].Revealed {
					p.known[y][x] = true
				}
			}
		}
		g.logTo(p, g.data.Text("klass_kartritare_karta", "namn", o.Name), "info")
		g.logTo(o, g.data.Text("klass_kartritare_ger", "namn", p.Name), "dim")
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
