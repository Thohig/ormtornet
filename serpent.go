package main

// Skattkammaren och ormen. Varje gång någon tar skatt dras ett ormkort.
// Sovkort tas bort ur leken; vakna-korten blir kvar. Ju färre kort som är
// kvar, desto större är risken att ormen vaknar.

import (
	"math/rand/v2"
	"strconv"
)

type wakeMsg struct {
	T string `json:"t"`
}

// buildDeck blandar en ny lek. Varje gång ormen har vaknat blir det färre
// sovkort, så att den vaknar lättare nästa gång.
func (g *Game) buildDeck() {
	o := &g.data.Ormen
	g.deck = g.deck[:0]
	sleepTotal := 0
	for _, c := range o.Kort {
		if !c.Vaknar {
			sleepTotal += c.Antal
		}
	}
	// Ta bort sovkort ett i taget från den sort som har flest kvar.
	counts := make([]int, len(o.Kort))
	for i, c := range o.Kort {
		counts[i] = c.Antal
	}
	remove := min(g.wakes*o.FarreEfterVakning, max(0, sleepTotal-o.MinstSovkort))
	for range remove {
		most := -1
		for i, c := range o.Kort {
			if !c.Vaknar && counts[i] > 0 && (most < 0 || counts[i] > counts[most]) {
				most = i
			}
		}
		if most < 0 {
			break
		}
		counts[most]--
	}
	g.deckSleep = 0
	for i, c := range o.Kort {
		if !c.Vaknar {
			g.deckSleep += counts[i]
		}
		for range counts[i] {
			g.deck = append(g.deck, c.ID)
		}
	}
	rand.Shuffle(len(g.deck), func(i, j int) { g.deck[i], g.deck[j] = g.deck[j], g.deck[i] })
}

// sleepLevel säger hur djupt ormen sover: 0 djupt, 1 oroligt, 2 nära att vakna.
func (g *Game) sleepLevel() int {
	sleep := 0
	for _, id := range g.deck {
		if c := g.data.card(id); c != nil && !c.Vaknar {
			sleep++
		}
	}
	switch {
	case g.deckSleep == 0:
		return 2
	case sleep*3 > g.deckSleep*2:
		return 0
	case sleep*3 > g.deckSleep:
		return 1
	}
	return 2
}

// take låter spelaren ta en näve skatt och drar ett ormkort.
func (g *Game) take(p *Player) {
	if g.rooms[p.Y][p.X].Kind != "kammare" {
		return
	}
	if len(g.deck) == 0 {
		g.buildDeck()
	}
	i := rand.IntN(len(g.deck))
	c := g.data.card(g.deck[i])
	if c == nil {
		g.deck = append(g.deck[:i], g.deck[i+1:]...)
		return
	}
	if c.Vaknar {
		g.logTo(p, pick(c.Texter), "varn")
		p.Seen["ormen"] = true
		g.wakeSerpent()
		g.explainNow(p, "ormen")
		return
	}
	g.deck = append(g.deck[:i], g.deck[i+1:]...)

	o := &g.data.Ormen
	n := o.SkattMin + rand.IntN(max(1, o.SkattMax-o.SkattMin+1))
	p.Treasure += n
	g.soundTo(p, "skatt")
	g.logTo(p, fill(pick(c.Texter), "antal", strconv.Itoa(n)), "syn")
	g.logTo(p, g.data.Text("skatt_totalt", "skatt", strconv.Itoa(p.Treasure)), "dim")
	if other := pick(c.Andra); other != "" {
		g.logOthers(p, other, "info")
		g.soundAll("ormen_ror")
	}
}

// wakeSerpent: alla i kammaren förlorar skatt och liv och kastas ut.
// Väsandet hörs i hela tornet. Sedan somnar ormen om, lättare väckt än förut.
func (g *Game) wakeSerpent() {
	o := &g.data.Ormen
	g.soundAll("vasning")
	for _, p := range g.players {
		if p != nil {
			g.send(p.client, wakeMsg{T: "vaknar"})
		}
	}
	var inside []*Player
	for _, p := range g.players {
		if p == nil || p.Status != statusIn || p.Cellar > 0 {
			continue
		}
		if g.rooms[p.Y][p.X].Kind == "kammare" {
			inside = append(inside, p)
		} else {
			g.logTo(p, pick(o.VaknarAndra), "varn")
		}
	}
	for _, p := range inside {
		lost := int(float64(p.Treasure)*o.SkattForlust + 0.5)
		p.Treasure -= lost
		g.logTo(p, fill(pick(o.KastasUt), "forlust", strconv.Itoa(lost), "skada", strconv.Itoa(o.Skada)), "varn")
		g.damage(p, o.Skada, true)
		if p.Status != statusIn {
			continue
		}
		var doors []Dir
		for d := North; d <= West; d++ {
			if g.rooms[p.Y][p.X].Doors[d] {
				doors = append(doors, d)
			}
		}
		if len(doors) > 0 {
			p.Stuck = 0
			g.move(p, doors[rand.IntN(len(doors))], 2)
		}
	}
	g.wakes++
	g.buildDeck()
	g.logAll(pick(o.SomnarOm), "dim")
}
