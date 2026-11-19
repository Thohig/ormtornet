package main

// Avslöjandet. Tornet visar gradvis vad det är, i fyra nivåer som följer
// solen. Alla texter ligger i data/avslojande.json och i rumsbeskrivningarna
// i data/texter.json. Spelarna får olika ledtrådar, så att de måste prata med
// varandra för att förstå hela bilden.

import (
	"strconv"
	"strings"
)

// levelFor räknar ut nivån för ett visst solvärde.
func (g *Game) levelFor(sun int) int {
	if !g.cfg.Reveal {
		return 1
	}
	frac := float64(sun) / float64(g.sunMax)
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
	if next > g.level && !g.levelHeld {
		for g.level < next {
			g.level++
		}
		g.onLevelUp()
		return
	}
	n := g.data.Avslojande.LedtradVarNRunda
	switch {
	case g.round%n == 0:
		g.dealClues(g.data.Avslojande.Ledtradar, "ledtrad", "info")
	case g.round%n == n/2:
		// Pilgrimen får ledtrådar dubbelt så ofta.
		g.dealCluesTo(g.data.Avslojande.Ledtradar, "ledtrad", "info", func(p *Player) bool { return g.is(p, "pilgrim") })
	}
}

func (g *Game) onLevelUp() {
	g.soundAll("niva")
	g.dealClues(g.data.Avslojande.Nivabyte, "nivabyte", "syn")
	g.dealClues(g.data.Avslojande.Ledtradar, "ledtrad", "info")
	g.dealNameParts()
	for _, p := range g.players {
		if p != nil && p.Status == statusIn && g.is(p, "pilgrim") {
			p.Potions++
			g.logTo(p, g.data.Text("klass_pilgrim_dryck"), "syn")
		}
	}
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
	g.dealCluesTo(pool, kind, class, func(*Player) bool { return true })
}

// dealCluesTo delar ut ledtrådar bara till dem som want godkänner.
func (g *Game) dealCluesTo(pool map[string][]string, kind, class string, want func(*Player) bool) {
	list := pool[strconv.Itoa(g.level)]
	if len(list) == 0 {
		return
	}
	var receivers []*Player
	for _, p := range g.players {
		// Ledtrådarna går bara till människor; datorspelare behöver dem inte.
		if p != nil && p.Status == statusIn && p.bot == nil && want(p) {
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
	g.rng.Shuffle(len(free), func(i, j int) { free[i], free[j] = free[j], free[i] })
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
	g.rng.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
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
	if g.usesName() && len(a.EpilogNamn) > 0 {
		g.logAll(strings.ReplaceAll(pick(a.EpilogNamn), "{namn}", g.serpentNameText()), "syn")
	}
	for _, p := range g.players {
		if p == nil {
			continue
		}
		if p.initiated {
			g.logTo(p, pick(a.EpilogInvigd), "rost")
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
