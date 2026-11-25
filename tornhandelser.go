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
	"stannar": "klocka", "draget": "morker", "klirrar": "skatt",
	"packning": "hitta", "brev": "hitta", "sluter": "ormen_ror",
	"yngel": "monster", "drommer": "ormen_ror", "kryp": "monster",
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
	g.reopenDoors()
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
	case "stannar":
		return g.sun > 0 && g.sun < g.sunMax
	case "draget":
		return g.lightRounds == 0
	case "klirrar":
		return g.richest() != nil
	case "packning":
		return g.plainItem() != ""
	case "brev":
		return len(g.letterReaders()) > 0
	case "sluter":
		_, ok := g.shutSpot()
		return ok
	case "yngel":
		return g.loneWanderer() != nil
	case "drommer":
		return len(g.deck) > 0
	case "kryp":
		return len(g.crawlWitnesses()) > 0
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
	// hit är de det händer, när det inte är de som står där. other är sant
	// när det gick åt andra hållet (HarAnnars).
	var hit map[*Player]bool
	other := false
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
	case "stannar":
		g.sun = min(g.sunMax, g.sun+max(1, e.Sol))
	case "draget":
		for _, p := range g.inside() {
			if g.is(p, "lyktbarare") || g.keepsLight(p) || g.useItem(p, "veke") {
				continue
			}
			p.Dark = max(p.Dark, g.darkRounds(p, max(1, e.Rundor)))
		}
	case "klirrar":
		p := g.richest()
		if p == nil {
			return false
		}
		x, y = p.X, p.Y
		hit = map[*Player]bool{p: true}
	case "packning":
		id := g.plainItem()
		if id == "" {
			return false
		}
		x, y = g.quietSpot()
		g.rooms[y][x].items = append(g.rooms[y][x].items, id)
	case "brev":
		r := g.letterReaders()
		if len(r) == 0 {
			return false
		}
		p := r[g.rng.IntN(len(r))]
		x, y = p.X, p.Y
		hit = map[*Player]bool{p: true}
	case "sluter":
		sd, ok := g.shutSpot()
		if !ok {
			return false
		}
		sd.rounds = max(1, e.Rundor)
		nx, ny := sd.x+dirDX[sd.d], sd.y+dirDY[sd.d]
		g.rooms[sd.y][sd.x].Doors[sd.d] = false
		g.rooms[ny][nx].Doors[sd.d.Opp()] = false
		g.shutDoors = append(g.shutDoors, sd)
		x, y = sd.x, sd.y
		hit = map[*Player]bool{}
		for _, p := range g.inside() {
			if p.Cellar == 0 && (p.X == sd.x && p.Y == sd.y || p.X == nx && p.Y == ny) {
				hit[p] = true
			}
		}
	case "yngel":
		p := g.loneWanderer()
		t := g.data.monsterType(e.Varelse)
		if p == nil || t == nil {
			return false
		}
		g.addMonster(t, p.X, p.Y)
		x, y = p.X, p.Y
		hit = map[*Player]bool{p: true}
	case "drommer":
		other = g.rng.IntN(2) == 1
		for range max(1, e.Ormkort) {
			if other {
				g.stirSerpent()
			} else {
				g.calmSerpent()
			}
		}
		x, y = g.chamber()
	case "kryp":
		w := g.crawlWitnesses()
		if len(w) == 0 {
			return false
		}
		x, y = g.chamber()
		hit = map[*Player]bool{}
		for _, p := range w {
			hit[p] = true
		}
	default:
		return false
	}
	g.eventCount[id]++
	g.chronicle("kronika_handelse_"+id, "bild", e.Bild)
	for _, p := range g.inside() {
		at := x >= 0 && p.X == x && p.Y == y && p.Cellar == 0
		here := at
		if hit != nil {
			here = hit[p]
		}
		var text string
		switch {
		case here && other && len(e.HarAnnars) > 0:
			text = pick(e.HarAnnars)
		case here && len(e.Har) > 0:
			text = pick(e.Har)
		case at && len(e.Bredvid) > 0:
			text = pick(e.Bredvid)
		case len(e.Text) > 0:
			text = pick(e.Text)
		}
		if text == "" {
			continue // bara den det händer märker något
		}
		if x >= 0 && (p.X != x || p.Y != y) {
			text = fill(text, "horn", g.hornToward(p.X, p.Y, x, y))
		}
		if here {
			g.eventHere(p, id)
		}
		g.soundTo(p, eventSounds[id])
		g.logTo(p, text, "varn")
		g.addCard(p, gameCard{Slag: "handelse", Bild: e.Bild, Rubrik: e.Rubrik, Text: text})
	}
	return true
}

// inside är vandrarna som är kvar i tornet.
func (g *Game) inside() []*Player {
	var in []*Player
	for _, p := range g.players {
		if p != nil && p.Status == statusIn {
			in = append(in, p)
		}
	}
	return in
}

// richest är den i tornet som bär mest guld, utanför valven, när de är
// minst två där inne och den bär något. Lika mycket avgörs av slumpen.
func (g *Game) richest() *Player {
	in := g.inside()
	if len(in) < 2 {
		return nil
	}
	var best []*Player
	for _, p := range in {
		if p.Cellar > 0 || p.Treasure <= 0 {
			continue
		}
		switch {
		case len(best) == 0 || p.Treasure > best[0].Treasure:
			best = []*Player{p}
		case p.Treasure == best[0].Treasure:
			best = append(best, p)
		}
	}
	if len(best) == 0 {
		return nil
	}
	return best[g.rng.IntN(len(best))]
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
	case "skinn", "packning":
		g.pickUp(p)
	case "brev":
		g.readLetter(p)
	case "yngel":
		g.announceMonster(p)
	}
}

// shutDoor är en dörr som tornet har slutit: från rummet (x, y) åt d, i
// rounds rundor till.
type shutDoor struct {
	x, y   int
	d      Dir
	rounds int
}

// shutSpot är en dörr mellan två vanliga rum som någon har sett, där det
// ändå finns en annan väg mellan rummen. Tornet stänger aldrig den enda
// vägen någonstans.
func (g *Game) shutSpot() (shutDoor, bool) {
	var spots []shutDoor
	for y := range g.size {
		for x := range g.size {
			r := &g.rooms[y][x]
			if !r.Revealed || r.Kind != "rum" {
				continue
			}
			for _, d := range []Dir{East, South} {
				nx, ny := x+dirDX[d], y+dirDY[d]
				if !g.inBounds(nx, ny) || !r.Doors[d] {
					continue
				}
				if n := &g.rooms[ny][nx]; n.Revealed && n.Kind == "rum" && g.otherWay(x, y, d) {
					spots = append(spots, shutDoor{x: x, y: y, d: d})
				}
			}
		}
	}
	if len(spots) == 0 {
		return shutDoor{}, false
	}
	return spots[g.rng.IntN(len(spots))], true
}

// otherWay säger om man kan gå från rummet (x, y) till grannen åt d genom
// rum man har sett, utan dörren mellan dem.
func (g *Game) otherWay(x, y int, d Dir) bool {
	tx, ty := x+dirDX[d], y+dirDY[d]
	type pos struct{ x, y int }
	seen := map[pos]bool{{x, y}: true}
	queue := []pos{{x, y}}
	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]
		r := &g.rooms[c.y][c.x]
		for dd := North; dd <= West; dd++ {
			if c.x == x && c.y == y && dd == d || !r.Doors[dd] {
				continue
			}
			n := pos{c.x + dirDX[dd], c.y + dirDY[dd]}
			if !g.inBounds(n.x, n.y) || seen[n] || !g.rooms[n.y][n.x].Revealed || g.rooms[n.y][n.x].blocked > 0 {
				continue
			}
			if n.x == tx && n.y == ty {
				return true
			}
			seen[n] = true
			queue = append(queue, n)
		}
	}
	return false
}

// reopenDoors räknar ner de slutna dörrarna och öppnar dem igen.
func (g *Game) reopenDoors() {
	left := g.shutDoors[:0]
	for _, sd := range g.shutDoors {
		if sd.rounds--; sd.rounds > 0 {
			left = append(left, sd)
			continue
		}
		g.rooms[sd.y][sd.x].Doors[sd.d] = true
		g.rooms[sd.y+dirDY[sd.d]][sd.x+dirDX[sd.d]].Doors[sd.d.Opp()] = true
	}
	g.shutDoors = left
}

// loneWanderer är den som står ensam längst in i tornet, långt från
// portarna, där inget annat finns: där ynglet kryper ut ur väggen.
func (g *Game) loneWanderer() *Player {
	var best *Player
	far := -1
	for _, p := range g.inside() {
		if p.Cellar > 0 || g.isChamber(p.X, p.Y) || g.cornerIndex(p.X, p.Y) >= 0 || g.monsterAt(p.X, p.Y) != nil {
			continue
		}
		alone := true
		for _, o := range g.inside() {
			if o != p && o.Cellar == 0 && o.X == p.X && o.Y == p.Y {
				alone = false
			}
		}
		if !alone {
			continue
		}
		near := 1 << 30
		for _, c := range g.corners() {
			near = min(near, abs(c.x-p.X)+abs(c.y-p.Y))
		}
		if near > far {
			best, far = p, near
		}
	}
	return best
}

// plainItem drar ett föremål som inte är förbannat och inte är ett brev,
// för packningen som någon glömde.
func (g *Game) plainItem() string {
	for range 20 {
		id := g.data.randomItem(g.rng)
		if t := g.data.itemType(id); t != nil && !t.Forbannad && !t.Brev {
			return id
		}
	}
	return ""
}

// letterReaders är de som kan få ett brev under dörren: inne, inte i
// valven och inte i mörkret.
func (g *Game) letterReaders() []*Player {
	var r []*Player
	for _, p := range g.inside() {
		if p.Cellar == 0 && p.Dark == 0 && g.cornerIndex(p.X, p.Y) < 0 {
			r = append(r, p)
		}
	}
	return r
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

// crawlWitnesses är de som ser krypen fly från kammaren: inne, inte i
// valven och inte i kammaren, och högst avstand rum från den.
func (g *Game) crawlWitnesses() []*Player {
	e := g.data.towerEvent("kryp")
	if e == nil {
		return nil
	}
	cx, cy := g.chamber()
	var w []*Player
	for _, p := range g.inside() {
		if p.Cellar > 0 || g.isChamber(p.X, p.Y) {
			continue
		}
		if abs(p.X-cx)+abs(p.Y-cy) <= max(1, e.Avstand) {
			w = append(w, p)
		}
	}
	return w
}
