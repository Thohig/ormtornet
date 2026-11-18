package main

// Avslöjandet. Tornet visar gradvis vad det är, i fyra nivåer som följer
// solen. Alla texter ligger i data/avslojande.json och i rumsbeskrivningarna
// i data/texter.json. Spelarna får olika ledtrådar, så att de måste prata med
// varandra för att förstå hela bilden.

import (
	"math/rand/v2"
	"strconv"
)

// levelFor räknar ut nivån för ett visst solvärde.
func (g *Game) levelFor(sun int) int {
	if !g.cfg.Reveal {
		return 1
	}
	frac := float64(sun) / float64(g.cfg.SunMax)
	level := 1
	for _, limit := range g.data.Avslojande.Granser {
		if frac <= limit {
			level++
		}
	}
	return min(level, 4)
}

// updateLevel körs efter varje runda: höjer nivån när solen har sjunkit
// tillräckligt och delar ut nya ledtrådar med jämna mellanrum.
func (g *Game) updateLevel() {
	if !g.cfg.Reveal {
		return
	}
	next := g.levelFor(g.sun)
	if next > g.level {
		for g.level < next {
			g.level++
		}
		g.onLevelUp()
		return
	}
	if g.round%g.data.Avslojande.LedtradVarNRunda == 0 {
		g.dealClues(g.data.Avslojande.Ledtradar, "ledtrad", "info")
	}
}

func (g *Game) onLevelUp() {
	g.soundAll("niva")
	g.dealClues(g.data.Avslojande.Nivabyte, "nivabyte", "syn")
	g.dealClues(g.data.Avslojande.Ledtradar, "ledtrad", "info")
	if g.level == 4 {
		for _, p := range g.players {
			if p != nil && p.Status != statusOut {
				g.explainAll(p)
			}
		}
	}
}

// dealClues ger varje spelare i tornet en egen text för nivån, så att ingen
// får samma som någon annan. Texterna delas ut som kort ur en lek; när leken
// är slut blandas den om.
func (g *Game) dealClues(pool map[string][]string, kind, class string) {
	list := pool[strconv.Itoa(g.level)]
	if len(list) == 0 {
		return
	}
	var receivers []*Player
	for _, p := range g.players {
		// Ledtrådarna går bara till människor; datorspelare behöver dem inte.
		if p != nil && p.Status == statusIn && p.bot == nil {
			receivers = append(receivers, p)
		}
	}
	if len(receivers) == 0 {
		return
	}
	key := func(i int) string { return kind + ":" + strconv.Itoa(g.level) + ":" + strconv.Itoa(i) }
	var free []int
	for i := range list {
		if !g.usedClues[key(i)] {
			free = append(free, i)
		}
	}
	if len(free) < len(receivers) {
		for i := range list {
			delete(g.usedClues, key(i))
		}
		free = free[:0]
		for i := range list {
			free = append(free, i)
		}
	}
	rand.Shuffle(len(free), func(i, j int) { free[i], free[j] = free[j], free[i] })
	for n, p := range receivers {
		i := free[n%len(free)]
		g.usedClues[key(i)] = true
		g.logTo(p, list[i], class)
	}
}

// explainNow ger den nya förklaringen direkt när något händer på nivå 4.
func (g *Game) explainNow(p *Player, id string) {
	if !g.cfg.Reveal || g.level < 4 {
		return
	}
	if text := g.data.Avslojande.Forklaringar[id]; text != "" {
		g.logTo(p, text, "syn")
	}
}

// explainAll förklarar allt spelaren har varit med om, när nivå 4 nås.
func (g *Game) explainAll(p *Player) {
	var lines []string
	for id, seen := range p.Seen {
		if text := g.data.Avslojande.Forklaringar[id]; seen && text != "" {
			lines = append(lines, text)
		}
	}
	if len(lines) == 0 {
		return
	}
	rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
	g.logTo(p, g.data.Avslojande.ForklaringRubrik, "dim")
	for _, l := range lines {
		g.logTo(p, l, "syn")
	}
}

// epilogue säger hemligheten rakt ut när spelet är slut.
func (g *Game) epilogue() {
	if !g.cfg.Reveal {
		return
	}
	g.level = 4
	a := &g.data.Avslojande
	g.logAll(pick(a.Epilog), "syn")
	for _, p := range g.players {
		if p == nil {
			continue
		}
		switch p.Status {
		case statusOut:
			g.logTo(p, pick(a.EpilogUte), "info")
		case statusLost:
			g.logTo(p, pick(a.EpilogForlorad), "varn")
		case statusDead:
			g.logTo(p, pick(a.EpilogDod), "varn")
		}
	}
}
