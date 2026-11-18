package main

// Föremål: saker vandrarna hittar när de söker igenom ett rum och bär med
// sig, högst några stycken. Namn, texter och siffror står i
// data/foremal.json. Här står vad varje föremål gör, styrt av dess id.

import "strconv"

type ItemType struct {
	ID           string   `json:"id"`
	Namn         string   `json:"namn"`
	Kort         string   `json:"kort"` // med liten bokstav, för texterna
	Beskrivning  string   `json:"beskrivning"`
	Vikt         int      `json:"vikt"`
	Anvandningar int      `json:"anvandningar"` // 0 = verkar hela vandringen
	Verkar       []string `json:"verkar"`
	Slut         []string `json:"slut"`
}

type ItemData struct {
	MaxBara int        `json:"max_bara"`
	Lista   []ItemType `json:"lista"`
}

// Item är ett föremål i en väska. uses är hur många gånger det verkar till,
// 0 för föremål som verkar hela vandringen.
type Item struct {
	ID   string
	uses int
}

type itemView struct {
	Namn        string `json:"namn"`
	Beskrivning string `json:"beskrivning"`
	Kvar        int    `json:"kvar,omitempty"`
}

func (d *Data) itemType(id string) *ItemType {
	for i := range d.Foremal.Lista {
		if d.Foremal.Lista[i].ID == id {
			return &d.Foremal.Lista[i]
		}
	}
	return nil
}

func (d *Data) maxItems() int {
	if d.Foremal.MaxBara > 0 {
		return d.Foremal.MaxBara
	}
	return 3
}

// randomItem drar ett föremål efter vikterna, eller "" om det inte finns några.
func (d *Data) randomItem() string {
	w := map[string]int{}
	for _, t := range d.Foremal.Lista {
		if t.Vikt > 0 {
			w[t.ID] = t.Vikt
		}
	}
	if len(w) == 0 {
		return ""
	}
	return weighted(w)
}

// has ger föremålet id i vandrarens väska, eller nil.
func (p *Player) has(id string) *Item {
	for _, it := range p.items {
		if it.ID == id {
			return it
		}
	}
	return nil
}

// give lägger föremålet i väskan. Den är falsk när väskan redan är full.
func (g *Game) give(p *Player, id string) bool {
	if len(p.items) >= g.data.maxItems() {
		return false
	}
	t := g.data.itemType(id)
	if t == nil {
		return false
	}
	p.items = append(p.items, &Item{ID: id, uses: t.Anvandningar})
	p.stats.items++
	g.chronicle("kronika_foremal", "namn", p.Name, "foremal", lowerFirst(t.Namn))
	return true
}

// findItem: vandraren hittar ett föremål när den söker. Är väskan full blir
// föremålet kvar i rummet, åt den som kommer efter.
func (g *Game) findItem(p *Player) {
	id := g.data.randomItem()
	t := g.data.itemType(id)
	if t == nil {
		g.soundTo(p, "inget")
		g.logTo(p, g.data.SearchText("inget"), "dim")
		return
	}
	g.soundTo(p, "hitta")
	if g.give(p, id) {
		g.logTo(p, g.data.Text("foremal_hittar", "foremal", lowerFirst(t.Namn), "beskrivning", t.Beskrivning), "syn")
		return
	}
	room := &g.rooms[p.Y][p.X]
	room.items = append(room.items, id)
	g.logTo(p, g.data.Text("foremal_full", "foremal", lowerFirst(t.Namn), "max", strconv.Itoa(g.data.maxItems())), "info")
}

// pickUp: den som kommer in i ett rum där något ligger tar det, om det får
// plats i väskan.
func (g *Game) pickUp(p *Player) {
	if p.Status != statusIn || p.Cellar > 0 {
		return
	}
	room := &g.rooms[p.Y][p.X]
	var left []string
	for _, id := range room.items {
		t := g.data.itemType(id)
		if t == nil {
			continue
		}
		if g.give(p, id) {
			g.logTo(p, g.data.Text("foremal_ligger", "foremal", lowerFirst(t.Namn), "beskrivning", t.Beskrivning), "syn")
		} else {
			left = append(left, id)
			g.logTo(p, g.data.Text("foremal_ligger_full", "foremal", lowerFirst(t.Namn)), "dim")
		}
	}
	room.items = left
}

// useItem förbrukar en användning av föremålet id, om vandraren har det, och
// berättar vad som händer. Föremål som verkar hela vandringen förbrukas inte.
func (g *Game) useItem(p *Player, id string, vars ...string) bool {
	it := p.has(id)
	if it == nil {
		return false
	}
	t := g.data.itemType(id)
	if t != nil {
		if text := pick(t.Verkar); text != "" {
			g.logTo(p, fill(text, vars...), "info")
		}
	}
	if it.uses > 0 {
		it.uses--
		if it.uses == 0 {
			if t != nil {
				if text := pick(t.Slut); text != "" {
					g.logTo(p, text, "dim")
				}
			}
			g.dropItem(p, it)
		}
	}
	return true
}

func (g *Game) dropItem(p *Player, it *Item) {
	for i, o := range p.items {
		if o == it {
			p.items = append(p.items[:i], p.items[i+1:]...)
			return
		}
	}
}

// hidden halverar en förlust för den som bär väskan med dubbel botten.
func (g *Game) hiddenShare(p *Player, loss int) int {
	if loss <= 0 || p.has("vaska") == nil {
		return loss
	}
	kept := loss / 2
	g.useItem(p, "vaska", "antal", strconv.Itoa(kept))
	return loss - kept
}

// shieldDamage: skölden tar en skada av varje träff från en varelse.
func (g *Game) shieldDamage(p *Player, dmg int) int {
	if dmg <= 0 || p.has("skold") == nil {
		return dmg
	}
	g.useItem(p, "skold")
	return dmg - 1
}

// compass är vad kompassen visar: riktningen mot närmaste port, som ett av
// orden under "horn" i texter.json.
func (g *Game) compass(p *Player) string {
	if p.has("kompass") == nil || p.Status != statusIn || p.Cellar > 0 {
		return ""
	}
	best, bx, by := -1, 0, 0
	for _, cn := range g.corners() {
		d := abs(cn.x-p.X) + abs(cn.y-p.Y)
		if best < 0 || d < best {
			best, bx, by = d, cn.x, cn.y
		}
	}
	if best == 0 {
		return g.data.Text("kompass_har")
	}
	sx, sy := sign(bx-p.X), sign(by-p.Y)
	// Samma ordning som "horn": nordväst, nordöst, sydväst, sydöst, norr, söder, väster, öster.
	idx := map[[2]int]int{{-1, -1}: 0, {1, -1}: 1, {-1, 1}: 2, {1, 1}: 3, {0, -1}: 4, {0, 1}: 5, {-1, 0}: 6, {1, 0}: 7}[[2]int{sx, sy}]
	return g.data.Text("kompass_pekar", "horn", g.data.Horn[idx])
}

func sign(v int) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}

func (g *Game) itemViews(p *Player) []itemView {
	var out []itemView
	for _, it := range p.items {
		if t := g.data.itemType(it.ID); t != nil {
			out = append(out, itemView{Namn: t.Namn, Beskrivning: t.Beskrivning, Kvar: it.uses})
		}
	}
	return out
}

// itemNames är vandrarens föremål för spelledarens översikt.
func (g *Game) itemNames(p *Player) []string {
	var out []string
	for _, it := range p.items {
		if t := g.data.itemType(it.ID); t != nil {
			out = append(out, t.Kort)
		}
	}
	return out
}
