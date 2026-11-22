package main

// Fällor, sökande och drycker. Texter och siffror kommer från data/fallor.json;
// här står bara vad varje fälla gör.

import (
	"strconv"
)

type pingMsg struct {
	T    string `json:"t"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Slot int    `json:"slot"`
}

// damage drar liv. Når livet noll dör spelaren.
func (g *Game) damage(p *Player, n int, announce bool) {
	if n <= 0 || p.Status != statusIn {
		return
	}
	p.HP = max(0, p.HP-n)
	g.soundTo(p, "skada")
	if announce && p.HP > 0 {
		g.logTo(p, g.data.Text("skada", "skada", strconv.Itoa(n), "liv", strconv.Itoa(p.HP)), "varn")
	}
	if p.HP == 0 {
		g.kill(p)
	}
}

func (g *Game) kill(p *Player) {
	g.dropAll(p)
	p.Status = statusDead
	g.chronicleDeath(p)
	p.Stuck, p.Cellar, p.Poison, p.Dark = 0, 0, 0, 0
	p.Action = ""
	g.logTo(p, g.data.Text("dod"), "varn")
	g.logOthers(p, g.data.Text("dod_andra", "namn", p.Name), "varn")
	g.soundAll("dod")
}

// trigger låter en fälla slå till mot spelaren.
func (g *Game) trigger(p *Player, t *Trap, depth int) {
	if t == nil || p.Status != statusIn {
		return
	}
	room := &g.rooms[p.Y][p.X]
	region := g.data.region(p.X, p.Y, g.size)
	roomName := g.roomShort(room)
	room.Mark = t.Tecken
	// Fjällamuletten tar nästa fälla. Tecknet syns ändå.
	if g.useItem(p, "amulett", "falla", t.Kort) {
		return
	}
	p.Seen[t.ID] = true
	p.stats.traps++
	p.hurt("falla", t.Kort)

	// Fällan har inget namn för spelarna: texten säger vad som händer.
	g.logTo(p, pick(t.Texter), "varn")
	others := t.Andra
	if len(t.AndraDolt) > 0 {
		others = t.AndraDolt
	}
	general := ""
	if other := pick(others); other != "" {
		general = fill(other, "namn", p.Name, "region", region, "rum", roomName)
	}
	if t.ID == "eko" || t.ID == "klocka" {
		// Ekot och klockan hörs lika i hela tornet.
		if general != "" {
			g.hearOthers(p, general, "info")
		}
	} else {
		g.hear(p.X, p.Y, "ljud_falla", general, "info", p)
	}

	// Ljud som hörs i hela tornet, eller bara för den som drabbas.
	switch t.ID {
	case "klocka", "eko", "fallgaller", "rasande_golv":
		g.soundAll(t.ID)
	default:
		g.soundTo(p, t.ID)
	}

	switch t.ID {
	case "fallgaller", "skinn":
		p.Stuck = max(p.Stuck, t.Rundor)
	case "ormgrop":
		p.Stuck = max(p.Stuck, t.Rundor)
		g.damage(p, t.Skada, true)
	case "pilar":
		g.damage(p, t.Skada, true)
	case "dimma":
		if g.is(p, "ortkvinna") {
			g.logTo(p, g.data.Text("klass_ortkvinna_gift"), "info")
		} else {
			p.Poison = max(p.Poison, t.Rundor)
		}
		g.damage(p, t.Skada, true)
	case "rasande_golv":
		g.damage(p, t.Skada, true)
		if p.Status == statusIn && !g.useItem(p, "rep") {
			p.Cellar = max(1, t.Rundor)
			p.Stuck = 0
		}
	case "roterande":
		var doors []Dir
		for d := North; d <= West; d++ {
			if room.Doors[d] {
				doors = append(doors, d)
			}
		}
		if len(doors) > 0 && depth < 2 {
			g.move(p, doors[g.rng.IntN(len(doors))], depth+1)
		}
	case "klocka":
		g.sun -= max(1, t.Sol)
	case "morker":
		if g.is(p, "lyktbarare") {
			g.logTo(p, g.data.Text("klass_lyktbarare_morker"), "info")
		} else if !g.useItem(p, "veke") {
			p.Dark = max(p.Dark, t.Rundor)
		}
	case "eko":
		for _, o := range g.players {
			if o != nil {
				g.send(o.client, pingMsg{T: "ping", X: p.X, Y: p.Y, Slot: p.Slot})
			}
		}
	}
	g.explainNow(p, t.ID)
}

// leaveCellar låter spelaren klättra upp i ett slumpat rum.
func (g *Game) leaveCellar(p *Player) {
	var spots [][2]int
	for y := range g.size {
		for x := range g.size {
			if g.cornerIndex(x, y) < 0 && !(x == g.center() && y == g.center()) && (x != p.X || y != p.Y) {
				spots = append(spots, [2]int{x, y})
			}
		}
	}
	s := spots[g.rng.IntN(len(spots))]
	p.X, p.Y = s[0], s[1]
	room := g.reveal(p.X, p.Y)
	room.Event = "" // ingen ny fälla direkt efter källaren
	room.Walked = true
	g.noteRoom(p)
	g.newRoomSound(p)
	g.logTo(p, g.data.Text("kallare_ut", "rum", g.roomShort(room)), "info")
	g.describeLater(p, room, p.X, p.Y)
	g.noticeTrail(p)
	g.pickUp(p)
	g.logOthers(p, g.hiddenText("kallare_ut_andra", "namn", p.Name, "region", g.data.region(p.X, p.Y, g.size)), "info")
	g.announceMonster(p)
}

// search letar igenom rummet. Varje rum kan bara genomsökas en gång.
func (g *Game) search(p *Player) {
	room := &g.rooms[p.Y][p.X]
	if room.Kind != "rum" || room.Searched {
		return
	}
	room.Searched = true
	p.stats.searched++
	outcome := weighted(g.rng, g.data.Sok.Vikter)
	if outcome == "inget" && g.rng.Float64() < g.classVal(p, "foremal_i_stallet", 0) {
		outcome = "foremal" // tjuven hittar oftare något
	}
	switch outcome {
	case "skatt":
		lo, hi := g.data.Regler.SokSkattMin, g.data.Regler.SokSkattMax
		n := lo + g.rng.IntN(max(1, hi-lo+1))
		n = max(1, int(float64(n)*g.classVal(p, "sok_guld", 1)))
		p.Treasure += n
		g.soundTo(p, "skatt")
		g.logTo(p, g.data.SearchText("skatt", "antal", strconv.Itoa(n), "skatt", strconv.Itoa(p.Treasure)), "syn")
	case "dryck":
		if g.classVal(p, "inga_drycker", 0) > 0 {
			g.soundTo(p, "inget")
			g.logTo(p, g.data.SearchText("inget"), "dim")
			return
		}
		p.Potions++
		g.soundTo(p, "hitta")
		g.logTo(p, g.data.SearchText("dryck"), "syn")
	case "gang":
		if d, ok := g.secretPassage(p.X, p.Y); ok {
			g.soundAll("gang")
			g.logTo(p, g.data.SearchText("gang", "riktning", g.data.Riktningar[d]), "syn")
		} else {
			g.soundTo(p, "inget")
			g.logTo(p, g.data.SearchText("gang_ingen"), "dim")
		}
	case "foremal":
		g.findItem(p)
	case "falla":
		g.logTo(p, g.data.SearchText("falla"), "varn")
		g.trigger(p, g.data.randomTrap(g.rng), 0)
	default:
		// Kartritaren hittar alltid en hemlig gång om det går att göra en.
		if g.is(p, "kartritare") {
			if d, ok := g.secretPassage(p.X, p.Y); ok {
				g.soundAll("gang")
				g.logTo(p, g.data.SearchText("gang", "riktning", g.data.Riktningar[d]), "syn")
				return
			}
		}
		g.soundTo(p, "inget")
		g.logTo(p, g.data.SearchText("inget"), "dim")
	}
}

// secretPassage öppnar en vägg mot ett grannrum. Hörn och kammare har fasta
// dörrar och deras grannar har redan öppningar mot dem.
func (g *Game) secretPassage(x, y int) (Dir, bool) {
	room := &g.rooms[y][x]
	var cand []Dir
	for d := North; d <= West; d++ {
		nx, ny := x+dirDX[d], y+dirDY[d]
		if room.Doors[d] || !g.inBounds(nx, ny) {
			continue
		}
		if _, fixed := g.fixedDoors(nx, ny); fixed {
			continue
		}
		cand = append(cand, d)
	}
	if len(cand) == 0 {
		return 0, false
	}
	d := cand[g.rng.IntN(len(cand))]
	room.Doors[d] = true
	if n := &g.rooms[y+dirDY[d]][x+dirDX[d]]; n.Revealed {
		n.Doors[d.Opp()] = true
	}
	return d, true
}

func (g *Game) drink(p *Player) {
	if p.Potions <= 0 {
		return
	}
	p.Potions--
	p.HP = min(g.maxHP(p), p.HP+int(g.classVal(p, "dryck_helar", float64(g.data.Regler.DryckHelar))))
	g.soundTo(p, "drick")
	g.logTo(p, g.data.Text("drick", "liv", strconv.Itoa(p.HP)), "info")
	g.hear(p.X, p.Y, "ljud_drick", "", "dim", p)
	if p.Poison > 0 {
		p.Poison = 0
		g.logTo(p, g.data.Text("drick_gift"), "info")
	}
}
