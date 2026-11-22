package main

// Tornet förändras med avslöjandets nivåer, när avslöjandet är på: från
// nivå 3 blir fällorna som hör till ormens kropp vanligare i nya rum, och på
// nivå 4 sluter sig portarna en i taget. Den sista porten stängs aldrig.

import (
	"strconv"
)

// changes är sant när tornet förändras i den här vandringen.
func (g *Game) changes() bool {
	return g.cfg.Reveal
}

// randomTrap väljer en fälla för ett nytt rum. När tornet förändras gäller
// vikterna för den högsta nivån som inte är över avslöjandets.
func (g *Game) randomTrap() *Trap {
	if !g.changes() || g.level < 3 {
		return g.data.randomTrap(g.rng)
	}
	var over map[string]int
	for lvl := g.level; lvl >= 3 && over == nil; lvl-- {
		over = g.data.Avslojande.Forandring.Fallvikter[strconv.Itoa(lvl)]
	}
	w := map[string]int{}
	for _, t := range g.data.Fallor {
		w[t.ID] = t.Vikt
		if v, ok := over[t.ID]; ok {
			w[t.ID] = v
		}
	}
	return g.data.trap(weighted(g.rng, w))
}

// openPorts är index för portarna som fortfarande är öppna.
func (g *Game) openPorts() []int {
	var out []int
	for i, c := range g.corners() {
		if !g.rooms[c.y][c.x].Shut {
			out = append(out, i)
		}
	}
	return out
}

// closePorts körs efter varje runda. På nivå 4 varnar tornet för en port
// och stänger den när port_stangs_var rundor har gått, sedan nästa. Först
// stängs den port som ligger längst från vandrarna som är kvar inne.
func (g *Game) closePorts() {
	if !g.changes() || g.level < 4 {
		return
	}
	every := max(1, g.data.Avslojande.Forandring.PortStangsVar)
	g.portClock++
	corners := g.corners()
	if g.portNext >= 0 {
		if g.portClock < every {
			return
		}
		c := corners[g.portNext]
		g.reveal(c.x, c.y).Shut = true
		horn := g.data.Horn[g.portNext]
		text := g.data.Text("port_stangs", "horn", horn)
		g.tellInside(text, "varn")
		g.bigCardInside(gameCard{Slag: "port", Bild: "port", Rubrik: g.data.Text("kort_port"), Text: text})
		g.soundAll("niva")
		g.chronicle("kronika_port", "horn", horn)
		g.portNext, g.portClock = -1, 0
		return
	}
	if g.portClock < every-1 {
		return
	}
	open := g.openPorts()
	if len(open) <= 1 {
		return
	}
	best, bestD := -1, -1
	for _, i := range open {
		c := corners[i]
		near := -1
		for _, p := range g.players {
			if p != nil && p.Status == statusIn {
				if d := abs(p.X-c.x) + abs(p.Y-c.y); near < 0 || d < near {
					near = d
				}
			}
		}
		if near < 0 {
			return // ingen är kvar inne
		}
		if near > bestD || (near == bestD && g.rng.IntN(2) == 0) {
			best, bestD = i, near
		}
	}
	g.portNext = best
	g.tellInside(g.data.Text("port_darrar", "horn", g.data.Horn[best]), "varn")
}

// tellInside skriver en rad till alla som är kvar i tornet.
func (g *Game) tellInside(text, class string) {
	for _, p := range g.players {
		if p != nil && p.Status == statusIn {
			g.logTo(p, text, class)
		}
	}
}
