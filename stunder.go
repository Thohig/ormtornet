package main

// Vandringens största stunder: när vandringen är slut läggs några av
// krönikans händelser fram som stora spelkort, till exempel det första
// uppvaknandet, den första som dog och portarna som stängdes. Vilka
// händelser som räknas, i vilken ordning och hur många står under
// spelkort.stunder i data/texter.json.

import (
	"slices"
	"sort"
)

// MomentData styr vandringens största stunder.
type MomentData struct {
	Antal   int          `json:"antal"`    // så många kort som läggs fram
	PerSlag int          `json:"per_slag"` // så många av samma slag
	Slag    []MomentKind `json:"slag"`     // det viktigaste först
}

// MomentKind är ett slags stund: krönikans nycklar som hör dit, motivet och
// rubriken (en nyckel i texterna).
type MomentKind struct {
	ID          string   `json:"id"`
	Nycklar     []string `json:"nycklar"`
	Bild        string   `json:"bild"`
	Rubrik      string   `json:"rubrik"`
	BaraVinnare bool     `json:"bara_vinnare,omitempty"` // bara när den det handlar om vann
}

// momentView är en stund i krönikan: rundan och kortet.
type momentView struct {
	Round int      `json:"runda"`
	Kort  gameCard `json:"kort"`
}

// moments väljer vandringens största stunder: de viktigaste slagen först,
// de tidigaste först inom ett slag, högst PerSlag av varje och Antal i
// allt. De läggs fram i den ordning de hände.
func (g *Game) moments() []momentView {
	md := g.data.Spelkort.Stunder
	winners := map[string]bool{}
	for _, r := range g.results() {
		if r.Winner {
			winners[r.Name] = true
		}
	}
	type pick struct {
		rank, order int
		m           momentView
	}
	var all []pick
	for i, e := range g.chron {
		rank := slices.IndexFunc(md.Slag, func(k MomentKind) bool {
			return slices.Contains(k.Nycklar, e.key) && (!k.BaraVinnare || winners[e.who])
		})
		if rank < 0 {
			continue
		}
		k := md.Slag[rank]
		bild := k.Bild
		if e.bild != "" {
			bild = e.bild
		}
		c := gameCard{Slag: k.ID, Bild: bild, Rubrik: upperFirst(g.data.Text(k.Rubrik)), Text: e.text, Fallen: k.ID == "varelse"}
		all = append(all, pick{rank, i, momentView{Round: e.round, Kort: c}})
	}
	sort.SliceStable(all, func(a, b int) bool { return all[a].rank < all[b].rank })
	var chosen []pick
	per := map[int]int{}
	for _, p := range all {
		if len(chosen) >= md.Antal {
			break
		}
		if md.PerSlag > 0 && per[p.rank] >= md.PerSlag {
			continue
		}
		per[p.rank]++
		chosen = append(chosen, p)
	}
	sort.SliceStable(chosen, func(a, b int) bool { return chosen[a].order < chosen[b].order })
	out := make([]momentView, 0, len(chosen))
	for _, p := range chosen {
		out = append(out, p.m)
	}
	return out
}
