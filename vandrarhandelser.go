package main

// Händelserna mellan vandrarna (docs/handelser.md, data/handelser.json):
// tornet drar tillbaka den som har gått ut, ormens blick, skottpengen,
// skatan, väggen som faller, kistan med två lås, spåren som glöder och
// golvet som vrider sig. Som de andra händelserna gynnar de bara den som
// drabbas eller hinner först; inget av dem belönar samarbete.

import "strconv"

// eventState skriver in det vandraren ska se av händelserna i tillståndet:
// blicken och skottpengen bland tecknen, och vem hen kan dra tillbaka.
func (g *Game) eventState(me *Player, st *stateMsg) {
	if me.Status != statusIn {
		return
	}
	if g.gazeOn == me.ID && g.gazeRounds > 0 {
		st.Blick = g.gazeRounds
	}
	if g.bountyOn == me.ID && g.bountyRounds > 0 {
		st.Skottpeng = g.bountyRounds
	}
	if g.canPull(me) {
		if t := g.playerByID(me.pullTarget); t != nil {
			st.Dra = t.Name
		}
	}
}

// betweenTick räknar ner blicken och skottpengen i slutet av rundan, och
// låter blicken gå vidare till den som står i samma rum som den som bär
// den.
func (g *Game) betweenTick() {
	if g.gazeOn != 0 {
		g.gazeTick()
	}
	if g.bountyOn != 0 {
		h := g.playerByID(g.bountyOn)
		switch {
		case h == nil || h.Status != statusIn:
			g.bountyOn, g.bountyRounds = 0, 0
		default:
			if g.bountyRounds--; g.bountyRounds <= 0 {
				g.bountyOn, g.bountyRounds = 0, 0
				g.logTo(h, g.data.Text("handelse_skottpeng_slut"), "info")
			}
		}
	}
}

// gazeTick: blicken räknas ner varje runda och går vidare till någon i
// samma rum, men inte tillbaka till den som just gav den så länge de står
// kvar tillsammans. När den släpper tar den guld från den som bär den.
func (g *Game) gazeTick() {
	h := g.playerByID(g.gazeOn)
	if h == nil || h.Status != statusIn {
		g.gazeOn, g.gazeRounds, g.gazeFrom = 0, 0, 0
		return
	}
	if g.gazeRounds--; g.gazeRounds <= 0 {
		g.gazeOn, g.gazeRounds, g.gazeFrom = 0, 0, 0
		g.gazeFine(h)
		return
	}
	if h.Cellar > 0 {
		return
	}
	var near []*Player
	from := false
	for _, o := range g.inside() {
		if o != h && o.Cellar == 0 && o.X == h.X && o.Y == h.Y {
			if o.ID == g.gazeFrom {
				from = true
				continue
			}
			near = append(near, o)
		}
	}
	if !from {
		g.gazeFrom = 0 // de har skilts åt; möts de igen kan blicken gå tillbaka
	}
	if len(near) == 0 {
		return
	}
	n := near[g.rng.IntN(len(near))]
	g.gazeOn, g.gazeFrom = n.ID, h.ID
	g.gazeCard(h, g.data.Text("handelse_blick_vidare", "namn", n.Name))
	g.gazeCard(n, g.data.Text("handelse_blick_fick", "namn", h.Name, "rundor", strconv.Itoa(g.gazeRounds)))
	g.chronicle("kronika_blick_vidare", "namn", h.Name, "offer", n.Name)
}

// gazeFine: när blicken släpper tar den guld från den som bar den.
func (g *Game) gazeFine(h *Player) {
	fine := 0
	if e := g.data.towerEvent("blick"); e != nil {
		fine = min(e.Guld, h.Treasure)
	}
	if fine > 0 {
		h.Treasure -= fine
		g.goldNote(h, "blick", -fine)
		g.gazeCard(h, g.data.Text("handelse_blick_slut_guld", "guld", strconv.Itoa(fine)))
		return
	}
	g.gazeCard(h, g.data.Text("handelse_blick_slut"))
}

// gazeCard ger vandraren en rad och ett spelkort om ormens blick.
func (g *Game) gazeCard(p *Player, text string) {
	e := g.data.towerEvent("blick")
	if e == nil {
		return
	}
	g.logTo(p, text, "varn")
	g.soundTo(p, "ormen_ror")
	g.addCard(p, gameCard{Slag: "handelse", Bild: e.Bild, Rubrik: e.Rubrik, Text: text})
}

// giveGaze lägger ormens blick på p i rounds rundor. Bar någon annan den
// släpper den, utan att ta guld, och hen får veta det.
func (g *Game) giveGaze(p *Player, rounds int) {
	if o := g.playerByID(g.gazeOn); o != nil && o != p && o.Status == statusIn && g.gazeRounds > 0 {
		g.gazeCard(o, g.data.Text("handelse_blick_annan"))
	}
	g.gazeOn, g.gazeRounds, g.gazeFrom = p.ID, max(1, rounds), 0
}

// gazeWakes slår för om näven väcker ormen för den som bär dess blick.
func (g *Game) gazeWakes(p *Player) bool {
	if g.gazeOn != p.ID || g.gazeRounds <= 0 {
		return false
	}
	e := g.data.towerEvent("blick")
	return e != nil && g.rng.Float64() < e.Vacker
}

// gazeCandidate är den ormen öppnar ögat mot: någon som är inne och inte
// i valven.
func (g *Game) gazeCandidate() *Player {
	var c []*Player
	for _, p := range g.inside() {
		if p.Cellar == 0 {
			c = append(c, p)
		}
	}
	if len(c) < 2 {
		return nil
	}
	return c[g.rng.IntN(len(c))]
}

// ---- Tornet drar tillbaka ----

// pullable säger om tornet kan dra tillbaka t: ute med guld, inte dragen
// förut, porten hen gick ut genom är öppen och solen räcker för att gå ut
// igen.
func (g *Game) pullable(t *Player) bool {
	if t == nil || t.Status != statusOut || t.pulled || t.Treasure <= 0 || !g.inBounds(t.X, t.Y) {
		return false
	}
	r := &g.rooms[t.Y][t.X]
	e := g.data.towerEvent("tillbaka")
	return r.HasExit && !r.Shut && e != nil && g.sun >= e.SolMinst
}

// pullCandidate är den som kom ut med mest guld och kan dras tillbaka.
func (g *Game) pullCandidate() *Player {
	var best []*Player
	for _, p := range g.players {
		if p == nil || !g.pullable(p) {
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

// poorest är den inne, utanför valven, som bär minst guld: den som får
// välja att dra tillbaka någon.
func (g *Game) poorest() *Player {
	var best []*Player
	for _, p := range g.inside() {
		if p.Cellar > 0 {
			continue
		}
		switch {
		case len(best) == 0 || p.Treasure < best[0].Treasure:
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

// canPull säger om p kan dra tillbaka någon den här rundan.
func (g *Game) canPull(p *Player) bool {
	return p != nil && p.pullOffer >= g.round && p.Status == statusIn && p.Cellar == 0 &&
		g.pullable(g.playerByID(p.pullTarget))
}

// pullBack: p drar tillbaka den som gick ut med mest guld. Hälften av
// guldet är säkert ute, och hälften bär hen in i tornet igen, några rum in
// längs vägen hen gick ut (pullSpot). p får en del av det som följer med
// in, och ormens blick för det.
func (g *Game) pullBack(p *Player) {
	t := g.playerByID(p.pullTarget)
	p.pullOffer = 0
	if !g.pullable(t) {
		g.logTo(p, g.data.Text("handelse_tillbaka_inte"), "dim")
		return
	}
	e := g.data.towerEvent("tillbaka")
	gx, gy := t.X, t.Y
	carried := t.Treasure / 2
	t.banked = t.Treasure - carried
	share := int(float64(carried) * e.Andel)
	t.Treasure = carried
	t.Status = statusIn
	t.pulled = true
	t.HP = g.maxHP(t)
	t.Stuck, t.Cellar, t.Poison, t.Dark = 0, 0, 0, 0
	t.items = nil // föremålen blev guld vid porten
	t.exitAsked, t.still = 0, 0
	t.Action = ""
	t.Ready = t.bot != nil
	t.lastCard = nil
	inside := false
	if x, y, ok := g.pullSpot(t); ok {
		t.X, t.Y = x, y
		inside = true
	}
	t.forgetRoom()
	g.noteRoom(t)
	g.goldNote(t, "tillbaka", -t.banked)
	if share > 0 {
		t.Treasure -= share
		g.goldNote(t, "tillbaka_andel", -share, "namn", p.Name)
		p.Treasure += share
		g.goldNote(p, "drog", share, "namn", t.Name)
	}
	g.chronicle("kronika_tillbaka", "namn", p.Name, "offer", t.Name)

	key := "handelse_tillbaka_du_port"
	if inside {
		key = "handelse_tillbaka_du"
	}
	text := g.data.Text(key, "antal", strconv.Itoa(t.banked), "bar", strconv.Itoa(t.Treasure))
	if share > 0 {
		text += " " + g.data.Text("handelse_tillbaka_andel_du", "namn", p.Name, "guld", strconv.Itoa(share))
	}
	g.logTo(t, text, "varn")
	g.soundTo(t, "ormen_ror")
	g.learn(t, "handelse", "tillbaka")
	g.addCard(t, gameCard{Slag: "handelse", Bild: e.Bild, Rubrik: e.Rubrik, Text: text, Stor: true})

	g.giveGaze(p, g.gazeLength())
	drog := g.data.Text("handelse_tillbaka_drog", "namn", t.Name, "rundor", strconv.Itoa(g.gazeRounds))
	if share > 0 {
		drog += " " + g.data.Text("handelse_tillbaka_andel_drog", "guld", strconv.Itoa(share))
		g.soundTo(p, "skatt")
	}
	g.gazeCard(p, drog)
	g.learn(p, "handelse", "blick")

	for _, o := range g.inside() {
		if o == p || o == t || o.Cellar > 0 {
			continue
		}
		heard := g.data.Text("handelse_tillbaka_hors", "horn", g.hornToward(o.X, o.Y, gx, gy))
		g.logTo(o, heard, "varn")
		g.addCard(o, gameCard{Slag: "handelse", Bild: e.Bild, Rubrik: e.Rubrik, Text: heard})
	}
}

// pullSpot är rummet där den som dras tillbaka hamnar: så långt in som
// rum_in säger, men aldrig längre, längs rum hen själv har gått igenom och
// dörrar som är öppna. Kammaren, portrummen, rasade rum och rum med en
// varelse eller en annan vandrare hoppas över. Finns inget sådant rum
// står hen kvar i porten.
func (g *Game) pullSpot(t *Player) (int, int, bool) {
	e := g.data.towerEvent("tillbaka")
	if e == nil || e.RumIn <= 0 || t.known == nil {
		return 0, 0, false
	}
	type at struct{ x, y int }
	dist := map[at]int{{t.X, t.Y}: 0}
	queue := []at{{t.X, t.Y}}
	best, bestD := []at(nil), 0
	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]
		d := dist[c]
		if d > 0 && g.pullRoomOK(c.x, c.y) {
			switch {
			case d > bestD:
				best, bestD = []at{c}, d
			case d == bestD:
				best = append(best, c)
			}
		}
		if d >= e.RumIn {
			continue
		}
		for dir := North; dir <= West; dir++ {
			if !g.rooms[c.y][c.x].Doors[dir] {
				continue
			}
			n := at{c.x + dirDX[dir], c.y + dirDY[dir]}
			if !g.inBounds(n.x, n.y) || !t.known[n.y][n.x] || !g.rooms[n.y][n.x].Doors[dir.Opp()] {
				continue
			}
			if _, seen := dist[n]; seen || g.rooms[n.y][n.x].blocked > 0 {
				continue
			}
			dist[n] = d + 1
			queue = append(queue, n)
		}
	}
	if len(best) == 0 {
		return 0, 0, false
	}
	c := best[g.rng.IntN(len(best))]
	return c.x, c.y, true
}

// pullRoomOK säger om den som dras tillbaka kan hamna i rummet (x, y).
func (g *Game) pullRoomOK(x, y int) bool {
	r := &g.rooms[y][x]
	return r.Kind == "rum" && r.blocked == 0 && !g.isChamber(x, y) && g.cornerIndex(x, y) < 0 &&
		g.monsterAt(x, y) == nil && !g.playerIn(x, y)
}

// gazeLength är hur länge ormens blick ligger kvar.
func (g *Game) gazeLength() int {
	if e := g.data.towerEvent("blick"); e != nil {
		return max(1, e.Rundor)
	}
	return 6
}

// ---- Kistan med två lås ----

// chestSpot är ett vanligt rum som någon har gått i, där ingen står och
// inget annat finns: där kistan dyker upp.
func (g *Game) chestSpot() (int, int, bool) {
	var spots [][2]int
	for y := range g.size {
		for x := range g.size {
			r := &g.rooms[y][x]
			if r.Kind == "rum" && r.Walked && r.blocked == 0 && r.chest == 0 && !g.playerIn(x, y) && g.monsterAt(x, y) == nil {
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

// openChests: kistan går bara upp när minst två står i rummet, och då tar
// den av dem som bände (S) allt. Bände flera tar slumpen en av dem.
func (g *Game) openChests(active []*Player) {
	type pos struct{ x, y int }
	takers := map[pos][]*Player{}
	var order []pos
	for _, p := range active {
		if p.Status != statusIn || p.Cellar > 0 || p.did != actTake || g.rooms[p.Y][p.X].chest == 0 {
			continue
		}
		at := pos{p.X, p.Y}
		if len(takers[at]) == 0 {
			order = append(order, at)
		}
		takers[at] = append(takers[at], p)
	}
	for _, at := range order {
		r := &g.rooms[at.y][at.x]
		var here []*Player
		for _, o := range g.inside() {
			if o.Cellar == 0 && o.X == at.x && o.Y == at.y {
				here = append(here, o)
			}
		}
		if len(here) < 2 {
			for _, p := range takers[at] {
				g.logTo(p, g.data.Text("handelse_kista_ensam"), "dim")
			}
			continue
		}
		w := takers[at][g.rng.IntN(len(takers[at]))]
		gold := r.chest
		r.chest = 0
		w.Treasure += gold
		g.goldNote(w, "kista", gold)
		e := g.data.towerEvent("kista")
		text := g.data.Text("handelse_kista_du", "guld", strconv.Itoa(gold))
		g.logTo(w, text, "syn")
		g.soundTo(w, "skatt")
		g.addCard(w, gameCard{Slag: "guld", Bild: e.Bild, Rubrik: e.Rubrik, Text: text})
		g.chronicle("kronika_kista", "namn", w.Name, "guld", strconv.Itoa(gold))
		for _, o := range here {
			if o != w {
				g.logTo(o, g.data.Text("handelse_kista_annan", "namn", w.Name), "varn")
			}
		}
	}
}

// ---- Spåren glöder ----

// markGlow: medan spåren glöder lyser rummen där någon går.
func (g *Game) markGlow(fx, fy, tx, ty int) {
	if g.phase != phasePlay || g.round > g.glowUntil {
		return
	}
	g.rooms[fy][fx].glowed = true
	g.rooms[ty][tx].glowed = true
}

// glowRoom säger om rummet (x, y) syns för alla för att spåren glöder där.
func (g *Game) glowRoom(x, y int) bool {
	return g.phase == phasePlay && g.round <= g.glowUntil && g.rooms[y][x].glowed
}

// ---- Skatan, väggen och golvet ----

// magpie väljer den skatan tar ett föremål från, föremålets plats i
// väskan och den som får det: två vandrare i olika rum.
func (g *Game) magpie() (*Player, int, *Player, bool) {
	var victims []*Player
	for _, p := range g.inside() {
		if p.Cellar == 0 && g.magpieItem(p) >= 0 {
			victims = append(victims, p)
		}
	}
	if len(victims) == 0 {
		return nil, 0, nil, false
	}
	v := victims[g.rng.IntN(len(victims))]
	var takers []*Player
	for _, o := range g.inside() {
		if o != v && o.Cellar == 0 && (o.X != v.X || o.Y != v.Y) {
			takers = append(takers, o)
		}
	}
	if len(takers) == 0 {
		return nil, 0, nil, false
	}
	return v, g.magpieItem(v), takers[g.rng.IntN(len(takers))], true
}

// magpieItem är ett föremål skatan kan ta från p: inte förbannat, för det
// släpper inte. -1 om inget finns.
func (g *Game) magpieItem(p *Player) int {
	var idx []int
	for i, it := range p.items {
		if t := g.data.itemType(it.ID); t != nil && !t.Forbannad {
			idx = append(idx, i)
		}
	}
	if len(idx) == 0 {
		return -1
	}
	return idx[g.rng.IntN(len(idx))]
}

// magpieAny säger om skatan har någon att ta från och någon att ge till,
// utan att slå något.
func (g *Game) magpieAny() bool {
	for _, p := range g.inside() {
		if p.Cellar > 0 {
			continue
		}
		has := false
		for _, it := range p.items {
			if t := g.data.itemType(it.ID); t != nil && !t.Forbannad {
				has = true
			}
		}
		if !has {
			continue
		}
		for _, o := range g.inside() {
			if o != p && o.Cellar == 0 && (o.X != p.X || o.Y != p.Y) {
				return true
			}
		}
	}
	return false
}

// wallPair är två vandrare i grannrum utan dörr emellan, där väggen kan
// falla, och riktningen från den förste till den andre.
func (g *Game) wallPair() (*Player, *Player, Dir, bool) {
	type pair struct {
		a, b *Player
		d    Dir
	}
	var pairs []pair
	in := g.inside()
	for i, a := range in {
		for _, b := range in[i+1:] {
			if a.Cellar > 0 || b.Cellar > 0 {
				continue
			}
			for d := North; d <= West; d++ {
				if a.X+dirDX[d] != b.X || a.Y+dirDY[d] != b.Y {
					continue
				}
				ra, rb := &g.rooms[a.Y][a.X], &g.rooms[b.Y][b.X]
				if !ra.Doors[d] && ra.Kind == "rum" && rb.Kind == "rum" {
					pairs = append(pairs, pair{a, b, d})
				}
			}
		}
	}
	if len(pairs) == 0 {
		return nil, nil, 0, false
	}
	p := pairs[g.rng.IntN(len(pairs))]
	return p.a, p.b, p.d, true
}

// swapPair är två vandrare i olika vanliga rum, utan varelser, som golvet
// kan vrida sig under: de byter plats.
func (g *Game) swapPair() (*Player, *Player, bool) {
	var c []*Player
	for _, p := range g.inside() {
		if p.Cellar == 0 && g.cornerIndex(p.X, p.Y) < 0 && !g.isChamber(p.X, p.Y) && g.monsterAt(p.X, p.Y) == nil && g.rooms[p.Y][p.X].Kind == "rum" {
			c = append(c, p)
		}
	}
	var pairs [][2]*Player
	for i, a := range c {
		for _, b := range c[i+1:] {
			if a.X != b.X || a.Y != b.Y {
				pairs = append(pairs, [2]*Player{a, b})
			}
		}
	}
	if len(pairs) == 0 {
		return nil, nil, false
	}
	p := pairs[g.rng.IntN(len(pairs))]
	return p[0], p[1], true
}
