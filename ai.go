package main

// Datorstyrda vandrare. En datorspelare är en vanlig spelare utan webbläsare.
// Den fattar beslut utifrån samma vy som en människa får (stateFor), så den
// ser inte okända rum, andras skatt eller ormkorten.

import (
	"math/rand/v2"
	"slices"
	"strconv"
	"time"
)

// botBrain är det en datorspelare minns mellan rundorna.
type botBrain struct {
	pers   *Personality
	habits map[string]map[string]int // varelsetyp -> val -> antal gånger sett
}

// saw minns vad en varelse valde i en strid, precis som en människa kan läsa
// det i loggen.
func (b *botBrain) saw(monsterType, choice string) {
	if b.habits[monsterType] == nil {
		b.habits[monsterType] = map[string]int{}
	}
	b.habits[monsterType][choice]++
}

// counter väljer det drag som slår varelsens troligaste val. Utgångsläget är
// det en människa kan gissa ur varelsens beskrivning (vanorna); sedan väger
// det den själv har sett i strid tyngre.
func (b *botBrain) counter(rng *rand.Rand, t *MonsterType) string {
	if t == nil {
		return []string{actAttack, actParry, actDodge}[rng.IntN(3)]
	}
	weights := map[string]int{}
	for k, n := range t.Vanor {
		weights[k] += n
	}
	for k, n := range b.habits[t.ID] {
		weights[k] += n * 2
	}
	// Välj det val varelsen oftast gör, med lite slump vid lika.
	likely, best := "", -1
	for _, k := range []string{actAttack, actParry, actDodge} {
		if w := weights[k]; w > best || (w == best && rng.IntN(2) == 0) {
			likely, best = k, w
		}
	}
	switch likely {
	case actAttack:
		return actParry
	case actParry:
		return actDodge
	}
	return actAttack
}

// AddBots fyller upp till n lediga platser med datorspelare, från sista
// hörnet och bakåt så att den första människan får det första hörnet.
func (g *Game) AddBots(n int) int {
	ai := &g.data.AI
	added := 0
	for added < n {
		slot := -1
		for i := maxPlayers - 1; i >= 0; i-- {
			if g.players[i] == nil {
				slot = i
				break
			}
		}
		if slot < 0 || len(ai.Personligheter) == 0 {
			break
		}
		name := ""
		for _, cand := range ai.Namn {
			if g.playerByName(cand) == nil {
				name = cand
				break
			}
		}
		if name == "" {
			break
		}
		g.nextID++
		p := &Player{
			ID: g.nextID, Name: name, Slot: slot,
			bot: &botBrain{pers: &ai.Personligheter[g.botCount()%len(ai.Personligheter)], habits: map[string]map[string]int{}},
		}
		g.players[slot] = p
		g.fitTower(p)
		if g.phase == phaseChoose {
			g.botsChoose()
		}
		added++
	}
	return added
}

func (g *Game) playerByName(name string) *Player {
	for _, p := range g.players {
		if p != nil && p.Name == name {
			return p
		}
	}
	return nil
}

// scheduleBots låter datorspelarna tänka en kort stund och sedan välja.
func (g *Game) scheduleBots() {
	lo, hi := 700, 2200
	if t := g.data.AI.BetanketidMs; len(t) == 2 && t[1] >= t[0] {
		lo, hi = t[0], t[1]
	}
	for _, p := range g.players {
		g.scheduleBotIn(p, lo, hi)
	}
}

// scheduleBot låter en datorspelare välja sitt drag i den pågående rundan,
// till exempel när den just har tagit över en frånvarande vandrare.
func (g *Game) scheduleBot(p *Player) {
	lo, hi := 700, 2200
	if t := g.data.AI.BetanketidMs; len(t) == 2 && t[1] >= t[0] {
		lo, hi = t[0], t[1]
	}
	g.scheduleBotIn(p, lo, hi)
}

func (g *Game) scheduleBotIn(p *Player, lo, hi int) {
	if p == nil || p.bot == nil {
		return
	}
	if !g.isController(p) && (p.Status != statusIn || p.Cellar > 0) {
		return
	}
	round := g.round
	delay := time.Duration(lo+rand.IntN(hi-lo+1)) * time.Millisecond
	time.AfterFunc(delay, func() {
		g.mu.Lock()
		defer g.mu.Unlock()
		g.botAct(p, round)
	})
}

// botAct fattar och lämnar datorspelarens beslut för rundan.
func (g *Game) botAct(p *Player, round int) {
	// En människa kan ha kommit tillbaka och tagit över efter att draget
	// planerades, och då finns ingen AI kvar som kan välja.
	if p.bot == nil || g.phase != phasePlay || g.round != round || p.Action != "" {
		return
	}
	if !g.isController(p) && (p.Status != statusIn || p.Cellar > 0) {
		return
	}
	a, m := g.botDecide(p)
	g.actPlayer(p, a, m)
	if g.phase == phasePlay && g.round == round && p.Action == "" {
		g.actPlayer(p, actWait, 0)
	}
}

type pos struct{ x, y int }

// botMap är det datorspelaren vet om tornet: kända rum, varelser och
// avstånd dit den kan gå.
type botMap struct {
	rooms    map[pos]roomView
	monsters map[pos]monsterView
	dist     map[pos]int
	first    map[pos]Dir // första steget på vägen dit
	rng      *rand.Rand
}

// places ger rummen i m i en fast ordning, så att samma frö ger samma val.
func places[V any](m map[pos]V) []pos {
	out := make([]pos, 0, len(m))
	for p := range m {
		out = append(out, p)
	}
	slices.SortFunc(out, func(a, b pos) int {
		if a.y != b.y {
			return a.y - b.y
		}
		return a.x - b.x
	})
	return out
}

// explore räknar ut avstånd från start genom kända rum. Okända rum och rum
// med varelser får ett avstånd men vägen fortsätter inte genom dem.
func explore(rng *rand.Rand, st *stateMsg, start pos) *botMap {
	bm := &botMap{
		rooms: map[pos]roomView{}, monsters: map[pos]monsterView{},
		dist: map[pos]int{start: 0}, first: map[pos]Dir{}, rng: rng,
	}
	for _, r := range st.Rooms {
		bm.rooms[pos{r.X, r.Y}] = r
	}
	for _, m := range st.Monsters {
		bm.monsters[pos{m.X, m.Y}] = m
	}
	queue := []pos{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		r, known := bm.rooms[cur]
		if !known {
			continue
		}
		if _, foe := bm.monsters[cur]; foe && cur != start {
			continue
		}
		for d := North; d <= West; d++ {
			if !r.D[d] {
				continue
			}
			n := pos{cur.x + dirDX[d], cur.y + dirDY[d]}
			if n.x < 0 || n.y < 0 || n.x >= st.Size || n.y >= st.Size {
				continue
			}
			if _, done := bm.dist[n]; done {
				continue
			}
			bm.dist[n] = bm.dist[cur] + 1
			if cur == start {
				bm.first[n] = d
			} else {
				bm.first[n] = bm.first[cur]
			}
			queue = append(queue, n)
		}
	}
	return bm
}

// toExit ger första steget och antalet rundor ut ur tornet, eller -1.
func (bm *botMap) toExit(start pos) (Dir, int) {
	best, step := -1, Dir(0)
	for _, p := range places(bm.rooms) {
		r := bm.rooms[p]
		if r.Exit < 0 {
			continue
		}
		dist, ok := bm.dist[p]
		if !ok {
			continue
		}
		if _, foe := bm.monsters[p]; foe && p != start {
			continue
		}
		total := dist + 1
		if best < 0 || total < best {
			best = total
			if p == start {
				step = Dir(r.Exit)
			} else {
				step = bm.first[p]
			}
		}
	}
	return step, best
}

// toward ger första steget mot det billigaste rum som uppfyller want.
// Kostnaden är vägens längd plus extra(rummet).
func (bm *botMap) toward(want func(pos) bool, extra func(pos) int) (Dir, bool) {
	best, step, found := 0, Dir(0), false
	for _, p := range places(bm.dist) {
		d := bm.dist[p]
		if d == 0 || !want(p) {
			continue
		}
		cost := d + extra(p)
		if !found || cost < best || (cost == best && bm.rng.IntN(2) == 0) {
			best, step, found = cost, bm.first[p], true
		}
	}
	return step, found
}

var dirLetter = [4]string{"N", "E", "S", "W"}

// botDecide väljer datorspelarens drag. Den returnerar draget och, för den
// som styr varelserna, vilken varelse som ska flyttas.
func (g *Game) botDecide(p *Player) (string, int) {
	if g.isController(p) {
		return g.botControl()
	}
	st := g.stateFor(p)
	pers := p.bot.pers
	here := pos{p.X, p.Y}
	bm := explore(g.rng, &st, here)
	room, known := bm.rooms[here]
	exitDir, exitDist := bm.toExit(here)

	// Strid: slå tillbaka, eller fly om nästa slag kan döda.
	if foe, ok := bm.monsters[here]; ok {
		typ := g.data.monsterType(g.monsterTypeID(foe.ID))
		if typ != nil && st.HP <= typ.Skada+pers.FlyMarginal && foe.HP > 1 {
			if exitDist > 0 {
				return g.botFlee(p, exitDir), 0
			}
			if known {
				var doors []Dir
				for d := North; d <= West; d++ {
					if room.D[d] {
						doors = append(doors, d)
					}
				}
				if len(doors) > 0 {
					return g.botFlee(p, doors[g.rng.IntN(len(doors))]), 0
				}
			}
		}
		c := p.bot.counter(g.rng, typ)
		if c == actDodge && g.is(p, "kartritare") {
			c = actParry // kartritaren vet att undanmanövern aldrig lyckas
		}
		return c, 0
	}

	// Strid mot en annan vandrare: ibland, beroende på hur stridslysten den är.
	if foe := g.botFoe(p); foe != nil && g.rng.Float64() < pers.Stridslust {
		a := actAttack
		if r := g.rng.Float64(); r > 0.8 {
			a = actDodge
		} else if r > 0.6 {
			a = actParry
		}
		return a, -foe.ID
	}

	if st.HP <= pers.DrickVidLiv && st.Potions > 0 {
		return actDrink, 0
	}

	// Kastkniven mot en varelse i ett grannrum, innan den går in.
	if slot := botItemSlot(p, "kastkniv"); slot >= 0 && known {
		for d := North; d <= West; d++ {
			if _, foe := bm.monsters[pos{p.X + dirDX[d], p.Y + dirDY[d]}]; foe && room.D[d] {
				return "f" + strconv.Itoa(slot+1) + dirLetter[d], 0
			}
		}
	}
	// Timglaset när solen hinner ner innan den är ute.
	if slot := botItemSlot(p, "timglas"); slot >= 0 && exitDist >= 0 && st.Sun <= exitDist {
		return "f" + strconv.Itoa(slot+1), 0
	}

	// Gå ut i tid, när väskan är full nog eller när livet tryter.
	mustLeave := exitDist >= 0 && st.Sun <= exitDist+pers.Marginal
	hurt := st.HP <= pers.GaHemVidLiv
	if exitDist >= 0 && (mustLeave || p.Treasure >= pers.SkattMal || hurt) {
		return dirLetter[exitDir], 0
	}

	// Vilse: ingen känd väg ut, till exempel efter ett fall genom golvet.
	// Med skatt i väskan, eller när solen börjar sjunka, letar den sig ut
	// mot portarna i stället för in mot kammaren.
	if exitDist < 0 && (p.Treasure > 0 || st.Sun <= 2*g.size) {
		if d, ok := g.botTowardPorts(bm); ok {
			return dirLetter[d], 0
		}
	}

	// I kammaren: ta skatt så länge ormen sover tillräckligt djupt.
	if known && room.K == "kammare" {
		chance := 0.0
		if st.Sleep >= 0 && len(pers.TaVidSomn) > 0 {
			chance = pers.TaVidSomn[min(st.Sleep, len(pers.TaVidSomn)-1)]
		}
		if st.HP <= g.data.Ormen.Skada {
			chance = 0
		}
		if g.rng.Float64() < chance || (p.Treasure == 0 && st.HP > g.data.Ormen.Skada) {
			return actTake, 0
		}
		if exitDist >= 0 {
			return dirLetter[exitDir], 0
		}
	}

	// Leta ibland igenom rummet om det finns tid.
	plenty := exitDist < 0 || st.Sun > exitDist+pers.Marginal+6
	if known && room.K == "rum" && !room.S && plenty && g.rng.Float64() < pers.SokChans {
		return actSearch, 0
	}

	// Mot kammaren om vägen är känd, annars utforska närmaste okända rum.
	isCenter := func(q pos) bool { return q.x == g.center() && q.y == g.center() }
	none := func(pos) int { return 0 }
	if d, ok := bm.toward(isCenter, none); ok {
		return dirLetter[d], 0
	}
	// Utforska helst åt mitten till, där kammaren är.
	unknown := func(q pos) bool { _, k := bm.rooms[q]; return !k }
	toCenter := func(q pos) int { return abs(q.x-g.center()) + abs(q.y-g.center()) }
	if d, ok := bm.toward(unknown, toCenter); ok {
		return dirLetter[d], 0
	}

	// Ingen väg framåt: leta efter hemliga gångar, eller vänta.
	if known && room.K == "rum" && !room.S {
		return actSearch, 0
	}
	if known {
		var doors []Dir
		for d := North; d <= West; d++ {
			if room.D[d] {
				doors = append(doors, d)
			}
		}
		if len(doors) > 0 {
			return dirLetter[doors[g.rng.IntN(len(doors))]], 0
		}
	}
	return actWait, 0
}

// botTowardPorts ger första steget mot det okända rum som ligger närmast
// en port. Var portarna sitter vet alla vandrare, så det kräver inga rum
// som AI-spelaren inte har sett.
func (g *Game) botTowardPorts(bm *botMap) (Dir, bool) {
	unknown := func(q pos) bool { _, k := bm.rooms[q]; return !k }
	toPort := func(q pos) int {
		best := -1
		for _, c := range g.corners() {
			if d := abs(q.x-c.x) + abs(q.y-c.y); best < 0 || d < best {
				best = d
			}
		}
		return best
	}
	return bm.toward(unknown, toPort)
}

func (g *Game) monsterTypeID(id int) string {
	if m := g.monsterByID(id); m != nil {
		return m.Type
	}
	return ""
}

// botControl: en död datorspelare flyttar den varelse som kommer närmast en
// levande vandrare.
func (g *Game) botControl() (string, int) {
	var targets []pos
	for _, p := range g.players {
		if p != nil && p.Status == statusIn && p.Cellar == 0 {
			targets = append(targets, pos{p.X, p.Y})
		}
	}
	if len(targets) == 0 {
		return actWait, 0
	}
	near := func(x, y int) int {
		best := -1
		for _, t := range targets {
			d := abs(t.x-x) + abs(t.y-y)
			if best < 0 || d < best {
				best = d
			}
		}
		return best
	}
	bestGain, bestID, bestDir := 0, 0, Dir(0)
	for _, m := range g.monsters {
		if !g.rooms[m.Y][m.X].Revealed {
			continue
		}
		now := near(m.X, m.Y)
		if now == 0 {
			continue // står redan hos någon
		}
		for d := North; d <= West; d++ {
			if !g.canMonsterGo(m, d) {
				continue
			}
			gain := now - near(m.X+dirDX[d], m.Y+dirDY[d])
			if gain > bestGain || (gain == bestGain && gain > 0 && g.rng.IntN(2) == 0) {
				bestGain, bestID, bestDir = gain, m.ID, d
			}
		}
	}
	if bestID == 0 {
		return actWait, 0
	}
	return dirLetter[bestDir], bestID
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// botItemSlot ger platsen i väskan för föremålet id, eller -1.
func botItemSlot(p *Player, id string) int {
	for i, it := range p.items {
		if it.ID == id {
			return i
		}
	}
	return -1
}

// botFlee flyr åt hållet d, med rökpulver om AI-spelaren har det.
func (g *Game) botFlee(p *Player, d Dir) string {
	if slot := botItemSlot(p, "rokpulver"); slot >= 0 {
		return "f" + strconv.Itoa(slot+1) + dirLetter[d]
	}
	return dirLetter[d]
}
