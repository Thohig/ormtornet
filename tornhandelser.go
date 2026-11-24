package main

// Tornets slumpade händelser (docs/handelser.md, data/handelser.json).
// Från en viss runda kan en händelse hända i slutet av varje runda. Alla
// inne i tornet får ett spelkort med händelsen, och åt vilket håll den
// hände; den som står där får en egen text. Händelserna gynnar bara den
// som drabbas eller hinner dit först.

// eventSounds är ljudet för varje händelse, bland ljuden i web/audio.js.
var eventSounds = map[string]string{
	"andas": "ormen_ror", "guld": "hitta", "ljus": "rost", "skinn": "skinn",
	"ras": "rasande_golv", "doden": "monster", "klockslag": "klocka",
}

// towerEventTick slår för en händelse i slutet av rundan och räknar ner
// det som pågår (ljuset och raset).
func (g *Game) towerEventTick() {
	if g.lightRounds > 0 {
		g.lightRounds--
	}
	for y := range g.size {
		for x := range g.size {
			if r := &g.rooms[y][x]; r.blocked > 0 {
				r.blocked--
			}
		}
	}
	h := &g.data.Handelser
	if len(h.Lista) == 0 || g.round < h.FranRunda || g.rng.Float64() >= h.Chans {
		return
	}
	w := map[string]int{}
	for _, e := range h.Lista {
		if e.Vikt > 0 && (e.Max <= 0 || g.eventCount[e.ID] < e.Max) && g.eventPossible(e.ID) {
			w[e.ID] = e.Vikt
		}
	}
	if len(w) == 0 {
		return
	}
	g.runTowerEvent(weighted(g.rng, w))
}

// eventPossible säger om händelsen kan hända just nu.
func (g *Game) eventPossible(id string) bool {
	switch id {
	case "andas":
		_, _, _, ok := g.newDoorSpot()
		return ok
	case "doden":
		_, _, ok := g.deathSpot()
		return ok
	case "ras":
		_, _, ok := g.caveSpot()
		return ok
	}
	return true
}

// runTowerEvent låter händelsen id hända. Den returnerar false om den inte
// kunde hända.
func (g *Game) runTowerEvent(id string) bool {
	e := g.data.towerEvent(id)
	if e == nil {
		return false
	}
	x, y := -1, -1 // där det händer; -1 för hela tornet
	switch id {
	case "andas":
		ax, ay, d, ok := g.newDoorSpot()
		if !ok {
			return false
		}
		g.rooms[ay][ax].Doors[d] = true
		g.rooms[ay+dirDY[d]][ax+dirDX[d]].Doors[d.Opp()] = true
		x, y = ax, ay
	case "guld":
		x, y = g.quietSpot()
		g.spillGold(x, y, e.GuldMin+g.rng.IntN(max(1, e.GuldMax-e.GuldMin+1)), "handelse_guld_tar")
	case "ljus":
		g.lightRounds = max(1, e.Rundor)
		for _, p := range g.players {
			if p != nil && p.Status == statusIn {
				p.Dark = 0
			}
		}
	case "skinn":
		x, y = g.quietSpot()
		g.rooms[y][x].items = append(g.rooms[y][x].items, "skinn")
	case "ras":
		cx, cy, ok := g.caveSpot()
		if !ok {
			return false
		}
		x, y = cx, cy
		g.rooms[y][x].blocked = max(1, e.Rundor)
	case "doden":
		dx, dy, ok := g.deathSpot()
		if !ok {
			return false
		}
		x, y = dx, dy
		if t := g.data.monsterType(e.Varelse); t != nil {
			g.addMonster(t, x, y)
		}
	case "klockslag":
		g.sun = max(0, g.sun-max(1, e.Sol))
		x, y = g.chamber()
	default:
		return false
	}
	g.eventCount[id]++
	g.chronicle("kronika_handelse_" + id)
	g.soundAll(eventSounds[id])
	for _, p := range g.players {
		if p == nil || p.Status != statusIn {
			continue
		}
		here := x >= 0 && p.X == x && p.Y == y && p.Cellar == 0
		text := pick(e.Text)
		if here && len(e.Har) > 0 {
			text = pick(e.Har)
		} else if x >= 0 {
			text = fill(text, "horn", g.hornToward(p.X, p.Y, x, y))
		}
		if here {
			g.eventHere(p, id)
		}
		g.logTo(p, text, "varn")
		g.addCard(p, gameCard{Slag: "handelse", Bild: e.Bild, Rubrik: e.Rubrik, Text: text})
	}
	return true
}

// eventHere är det som händer den som står där händelsen sker.
func (g *Game) eventHere(p *Player, id string) {
	switch id {
	case "ras":
		p.hurt("tornet", "")
		g.damage(p, 1, false)
		if p.Status == statusIn && !g.useItem(p, "rep") {
			p.Cellar = 2
			p.Stuck = 0
		}
	case "klockslag":
		if g.rooms[p.Y][p.X].Kind == "kammare" {
			p.safeDraw = true
		}
	case "guld":
		g.pickUpGold(p)
	case "skinn":
		g.pickUp(p)
	}
}

// quietSpot är ett vanligt rum där ingen står och ingen har varit: där
// guldet eller skinnet hamnar, så att man måste gå dit.
func (g *Game) quietSpot() (int, int) {
	var spots [][2]int
	for y := range g.size {
		for x := range g.size {
			if g.cornerIndex(x, y) >= 0 || g.isChamber(x, y) || g.rooms[y][x].Walked || g.playerIn(x, y) {
				continue
			}
			spots = append(spots, [2]int{x, y})
		}
	}
	if len(spots) == 0 {
		cx, cy := g.chamber()
		return cx, cy
	}
	s := spots[g.rng.IntN(len(spots))]
	return s[0], s[1]
}

// newDoorSpot är två grannrum som båda finns men saknar dörr mellan sig,
// där tornet kan öppna en ny.
func (g *Game) newDoorSpot() (int, int, Dir, bool) {
	type spot struct {
		x, y int
		d    Dir
	}
	var spots []spot
	for y := range g.size {
		for x := range g.size {
			r := &g.rooms[y][x]
			if !r.Revealed || r.Kind != "rum" {
				continue
			}
			for _, d := range []Dir{East, South} {
				nx, ny := x+dirDX[d], y+dirDY[d]
				if !g.inBounds(nx, ny) || r.Doors[d] {
					continue
				}
				if n := &g.rooms[ny][nx]; n.Revealed && n.Kind == "rum" {
					spots = append(spots, spot{x, y, d})
				}
			}
		}
	}
	if len(spots) == 0 {
		return 0, 0, 0, false
	}
	s := spots[g.rng.IntN(len(spots))]
	return s.x, s.y, s.d, true
}

// caveSpot är ett vanligt rum som någon har gått igenom, där taket kan
// rasa. Inte där en varelse står.
func (g *Game) caveSpot() (int, int, bool) {
	var spots [][2]int
	for y := range g.size {
		for x := range g.size {
			r := &g.rooms[y][x]
			if r.Kind == "rum" && r.Walked && r.blocked == 0 && g.monsterAt(x, y) == nil {
				spots = append(spots, [2]int{x, y})
			}
		}
	}
	if len(spots) == 0 {
		return 0, 0, false
	}
	s := spots[g.rng.IntN(len(spots))]
	return s[0], s[1], true
}

// deathSpot är ett rum där någon har dött och ingen varelse står.
func (g *Game) deathSpot() (int, int, bool) {
	var spots [][2]int
	for y := range g.size {
		for x := range g.size {
			if len(g.rooms[y][x].fell) > 0 && g.monsterAt(x, y) == nil && !g.isChamber(x, y) && g.cornerIndex(x, y) < 0 {
				spots = append(spots, [2]int{x, y})
			}
		}
	}
	if len(spots) == 0 {
		return 0, 0, false
	}
	s := spots[g.rng.IntN(len(spots))]
	return s[0], s[1], true
}
