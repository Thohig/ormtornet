package main

// Varslen. Från nivå 2 är varannan ledtråd ett varsel: tornet förråder ett
// rum som ingen har varit i, två eller tre steg från vandraren. Vad som finns
// i rummet bestäms när varslet ges, med samma chanser som annars, och blir
// så när någon går in. Varslet säger aldrig var kammaren ligger.

import (
	"strings"
)

// omenNote är ett varsel som har getts, för spelledaren.
type omenNote struct {
	X, Y  int
	What  string // "falla:id", "varelse:id" eller "tomt"
	To    string // vandraren som fick det
	Round int
}

// randomTrap väljer en fälla för ett nytt rum.
func (g *Game) randomTrap() *Trap {
	return g.data.randomTrap(g.rng)
}

// rollOmen bestämmer vad som väntar i ett rum, med samma chanser som när
// rummet byggs: en varelse, annars kanske en fälla, annars ingenting.
func (g *Game) rollOmen() string {
	if g.rng.Float64() < g.data.Monster.Chans {
		w := map[string]int{}
		for _, t := range g.data.Monster.Typer {
			w[t.ID] = t.Vikt
		}
		if id := weighted(g.rng, w); id != "" {
			return "varelse:" + id
		}
	}
	if g.rng.Float64() < g.data.Regler.Handelsechans {
		if t := g.randomTrap(); t != nil {
			return "falla:" + t.ID
		}
	}
	return "tomt"
}

// fulfilOmen gör det varslet lovade när rummet byggs.
func (g *Game) fulfilOmen(r *Room, x, y int) {
	kind, id, _ := strings.Cut(r.omen, ":")
	switch kind {
	case "falla":
		r.Event = id
	case "varelse":
		if t := g.data.monsterType(id); t != nil && g.monsterAt(x, y) == nil {
			g.addMonster(t, x, y)
		}
	}
}

// omenRoom väljer ett rum som ingen har varit i, två eller tre steg från p,
// som inte är en port eller kammaren och inte redan har ett varsel.
func (g *Game) omenRoom(p *Player) (int, int, bool) {
	var cands [][2]int
	for y := range g.size {
		for x := range g.size {
			d := abs(x-p.X) + abs(y-p.Y)
			if d < 2 || d > 3 {
				continue
			}
			r := &g.rooms[y][x]
			if r.Revealed || r.omen != "" {
				continue
			}
			if _, fixed := g.fixedDoors(x, y); fixed {
				continue
			}
			cands = append(cands, [2]int{x, y})
		}
	}
	if len(cands) == 0 {
		return 0, 0, false
	}
	c := cands[g.rng.IntN(len(cands))]
	return c[0], c[1], true
}

// dealOmens ger varje människa i tornet ett varsel om sin omgivning. Den som
// inte har något rum att få varsel om får en vanlig ledtråd.
func (g *Game) dealOmens() {
	v := &g.data.Avslojande.Varsel
	for _, p := range g.players {
		if p == nil || p.bot != nil || p.Status != statusIn {
			continue
		}
		x, y, ok := g.omenRoom(p)
		if !ok || p.Cellar > 0 {
			g.dealCluesTo(g.data.Avslojande.Ledtradar, "ledtrad", "info", func(o *Player) bool { return o == p })
			continue
		}
		r := &g.rooms[y][x]
		r.omen = g.rollOmen()
		kind, id, _ := strings.Cut(r.omen, ":")
		list := v.Tomt
		switch kind {
		case "falla":
			list = v.Fallor[id]
		case "varelse":
			list = v.Varelser[id]
		}
		dist := abs(x-p.X) + abs(y-p.Y)
		avstand := ""
		if i := dist - 2; i >= 0 && i < len(v.Avstand) {
			avstand = v.Avstand[i]
		}
		text := strings.NewReplacer("{avstand}", avstand, "{riktning}", g.data.Horn[hornIndex(x-p.X, y-p.Y)]).Replace(pick(list))
		g.logTo(p, text, "info")
		g.omens = append(g.omens, omenNote{X: x, Y: y, What: r.omen, To: p.Name, Round: g.round})
	}
}

// hornIndex är väderstrecket åt (dx, dy), som index i horn.
func hornIndex(dx, dy int) int {
	return map[[2]int]int{{-1, -1}: 0, {1, -1}: 1, {-1, 1}: 2, {1, 1}: 3, {0, -1}: 4, {0, 1}: 5, {-1, 0}: 6, {1, 0}: 7}[[2]int{sign(dx), sign(dy)}]
}
