package main

// Den dödes guld: guldet en vandrare bar när hen dog ligger kvar i rummet,
// och den som kommer dit, eller står kvar där, tar det. Föremålen ligger
// redan kvar (dropAll i items.go).

import (
	"slices"
	"strconv"
)

// leftGold är guld som en död vandrare bar.
type leftGold struct {
	by   int // den döda vandrarens id
	name string
	gold int
}

// heardDeath är det en vandrare hörde när någon dog: var, så noga som hen
// kunde höra det, och om mynt rasslade. AI-spelarna går efter det (ai.go);
// en människa läser samma sak i loggen.
type heardDeath struct {
	x, y   int    // rummet, om vandraren stod där eller bredvid; annars -1
	region string // tornets del
	coins  bool
	name   string // den döda, som rumsblocket säger "Här föll" om
	round  int

	// För AI-spelaren: om den går dit (1), inte går (-1) eller inte har
	// bestämt sig (0), och rummen den har stått i sedan skriket.
	chase   int
	checked map[pos]bool
}

// announceDeath låter de andra höra att p dog. Den som står i rummet ser
// det, den som står bredvid hör det genom väggen åt ett håll, och alla
// andra hör i vilken del av tornet. Bar den döda guld hörs mynten.
func (g *Game) announceDeath(p *Player) {
	region := g.data.region(p.X, p.Y, g.size)
	coins := p.Treasure > 0
	for _, o := range g.players {
		if o == nil || o == p {
			continue
		}
		heard := heardDeath{x: -1, y: -1, region: region, coins: coins, name: p.Name, round: g.round}
		var text string
		dir, hears := g.hearsFrom(o, p.X, p.Y)
		switch {
		case o.Status == statusIn && o.Cellar == 0 && o.X == p.X && o.Y == p.Y:
			text = g.data.Text("dod_andra", "namn", p.Name)
			heard.x, heard.y = p.X, p.Y
		case hears:
			text = g.data.Text("ljud_dod", "namn", p.Name, "riktning", dir)
			heard.x, heard.y = p.X, p.Y
		default:
			text = g.data.Text("dod_andra_region", "namn", p.Name, "region", region)
		}
		if coins {
			text += " " + g.data.Text("dod_mynt")
		}
		g.logTo(o, text, "varn")
		if o.Status == statusIn {
			o.deaths = append(o.deaths, heard)
		}
	}
}

// dropGold lägger guldet den döda bar i rummet hen dog i.
func (g *Game) dropGold(p *Player) {
	if !g.inBounds(p.X, p.Y) {
		return
	}
	g.rooms[p.Y][p.X].fell = append(g.rooms[p.Y][p.X].fell, p.Name)
	if p.Treasure <= 0 {
		return
	}
	r := &g.rooms[p.Y][p.X]
	r.gold = append(r.gold, leftGold{by: p.ID, name: p.Name, gold: p.Treasure})
	g.logTo(p, g.data.Text("guld_kvar_du", "guld", strconv.Itoa(p.Treasure)), "dim")
	lost := p.Treasure
	p.Treasure = 0
	g.goldNote(p, "tappade", -lost)
}

// pickUpGold låter en levande vandrare ta guldet som ligger i rummet.
func (g *Game) pickUpGold(p *Player) {
	if p.Status != statusIn || p.Cellar > 0 || !g.inBounds(p.X, p.Y) {
		return
	}
	r := &g.rooms[p.Y][p.X]
	for _, l := range r.gold {
		p.Treasure += l.gold
		g.goldNote(p, "dod_guld", l.gold, "namn", l.name)
		n := strconv.Itoa(l.gold)
		g.soundTo(p, "skatt")
		text := g.data.Text("guld_kvar", "namn", l.name, "antal", n)
		g.logTo(p, text, "syn")
		g.addCard(p, gameCard{Slag: "guld", Bild: "guld", Rubrik: g.data.Text("kort_guld_kvar"), Text: text})
		g.chronicle("kronika_guld_kvar", "namn", p.Name, "dod", l.name, "guld", n)
		if dead := g.playerByID(l.by); dead != nil {
			dead.goldTakenBy = append(dead.goldTakenBy, p.ID)
		}
	}
	if len(r.gold) > 0 {
		r.goldBy = p.ID
	}
	r.gold = nil
	if r.loose > 0 {
		n := strconv.Itoa(r.loose)
		p.Treasure += r.loose
		kind := "lost"
		if r.looseKey == "ormen_guld_kvar" {
			kind = "lost_orm"
		}
		g.goldNote(p, kind, r.loose)
		g.soundTo(p, "skatt")
		text := g.data.Text(r.looseKey, "antal", n)
		g.logTo(p, text, "syn")
		g.addCard(p, gameCard{Slag: "guld", Bild: "guld", Rubrik: g.data.Text("kort_guld_lost"), Text: text})
		g.chronicle("kronika_guld_lost", "namn", p.Name, "guld", n)
		r.loose, r.looseKey = 0, ""
	}
}

// spillGold lägger löst guld i ett rum, som den som kommer först tar.
// key är texten när det tas, med {antal}.
func (g *Game) spillGold(x, y, amount int, key string) {
	if amount <= 0 || !g.inBounds(x, y) {
		return
	}
	r := &g.rooms[y][x]
	r.loose += amount
	r.looseKey = key
}

// goldEpilogue säger till den döda vart guldet tog vägen, i slutet.
func (g *Game) goldEpilogue() {
	for _, p := range g.players {
		if p == nil {
			continue
		}
		for _, id := range p.goldTakenBy {
			o := g.playerByID(id)
			if o == nil {
				continue
			}
			key := "guld_kvar_ute"
			if o.Status != statusOut {
				key = "guld_kvar_inte_ute"
			}
			text := g.data.Text(key, "namn", o.Name)
			g.logTo(p, text, "dim")
			p.epilog = append(p.epilog, chronLine{Text: text, Class: "dim"})
		}
	}
}

// botChaseRounds är hur många rundor en AI-spelare letar efter guldet
// innan den ger upp.
const botChaseRounds = 10

// botTowardGold ger AI-spelarens första steg mot guldet där någon dog, om
// den vill gå dit. Den vet bara det en människa vet: skriket (rummet, om
// den stod bredvid, annars tornets del), om mynt rasslade, och rummen den
// har sett. Den går bara med gott om liv och tid, och personligheten
// avgör hur gärna (guldlust i data/ai.json).
//
// maxDist > 0 är en omväg på vägen ut: då går den bara till ett känt rum
// högst så många steg bort.
func (g *Game) botTowardGold(p *Player, bm *botMap, st *stateMsg, pers *Personality, plenty bool, maxDist int) (Dir, bool) {
	here := pos{p.X, p.Y}
	room := &g.rooms[p.Y][p.X]
	var target *heardDeath
	for i := range p.deaths {
		h := &p.deaths[i]
		if h.chase < 0 {
			continue
		}
		if h.checked == nil {
			h.checked = map[pos]bool{}
		}
		h.checked[here] = true
		// Rumsblocket säger "Här föll" i rummet: där var det, och guldet är
		// taget nu (pickUp), av den här vandraren eller någon annan.
		if !h.coins || slices.Contains(room.fell, h.name) || g.round-h.round > botChaseRounds {
			h.chase = -1
			continue
		}
		if h.chase == 0 {
			h.chase = -1
			if g.rng.Float64() < pers.Guldlust {
				h.chase = 1
			}
		}
		if h.chase > 0 && target == nil {
			target = h
		}
	}
	if target == nil || !plenty || st.HP*2 < st.MaxHP {
		return 0, false
	}
	want := func(q pos) bool {
		if maxDist > 0 {
			return q.x == target.x && q.y == target.y && bm.dist[q] <= maxDist
		}
		if target.x >= 0 {
			return q.x == target.x && q.y == target.y
		}
		return !target.checked[q] && g.data.region(q.x, q.y, g.size) == target.region
	}
	d, ok := bm.toward(want, func(pos) int { return 0 })
	if !ok && maxDist == 0 {
		target.chase = -1 // ingen väg dit som den känner till
	}
	return d, ok
}
