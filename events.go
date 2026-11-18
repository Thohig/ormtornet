package main

// Fällor, sökande och drycker. Texter och siffror kommer från data/fallor.json;
// här står bara vad varje fälla gör.

import (
	"math/rand/v2"
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
	p.Status = statusDead
	p.Stuck, p.Cellar, p.Poison, p.Dark = 0, 0, 0, 0
	p.Action = ""
	g.logTo(p, g.data.Text("dod"), "varn")
	g.logOthers(p, g.data.Text("dod_andra", "namn", p.Name), "varn")
	g.soundAll("dod")
	g.makeController(p)
}

// trigger låter en fälla slå till mot spelaren.
func (g *Game) trigger(p *Player, t *Trap, depth int) {
	if t == nil || p.Status != statusIn {
		return
	}
	room := &g.rooms[p.Y][p.X]
	region := g.data.region(p.X, p.Y)
	roomName := g.roomName(room)
	room.Mark = t.Tecken

	g.logTo(p, t.Namn+"!", "varn")
	g.logTo(p, pick(t.Texter), "")
	if other := pick(t.Andra); other != "" {
		g.logOthers(p, fill(other, "namn", p.Name, "region", region, "rum", roomName), "info")
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
		p.Poison = max(p.Poison, t.Rundor)
		g.damage(p, t.Skada, true)
	case "rasande_golv":
		g.damage(p, t.Skada, true)
		if p.Status == statusIn {
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
			g.move(p, doors[rand.IntN(len(doors))], depth+1)
		}
	case "klocka":
		g.sun -= max(1, t.Sol)
	case "morker":
		p.Dark = max(p.Dark, t.Rundor)
	case "eko":
		for _, o := range g.players {
			if o != nil {
				g.send(o.client, pingMsg{T: "ping", X: p.X, Y: p.Y, Slot: p.Slot})
			}
		}
	}
}

// leaveCellar låter spelaren klättra upp i ett slumpat rum.
func (g *Game) leaveCellar(p *Player) {
	var spots [][2]int
	for y := range gridSize {
		for x := range gridSize {
			if cornerIndex(x, y) < 0 && !(x == center && y == center) && (x != p.X || y != p.Y) {
				spots = append(spots, [2]int{x, y})
			}
		}
	}
	s := spots[rand.IntN(len(spots))]
	p.X, p.Y = s[0], s[1]
	room := g.reveal(p.X, p.Y)
	room.Event = "" // ingen ny fälla direkt efter källaren
	g.soundAll("nyttrum")
	g.logTo(p, g.data.Text("kallare_ut", "rum", g.roomName(room)), "info")
	g.logTo(p, g.describe(room), "")
	g.logOthers(p, g.data.Text("kallare_ut_andra", "namn", p.Name, "region", g.data.region(p.X, p.Y)), "info")
	g.announceMonster(p)
}

// search letar igenom rummet. Varje rum kan bara genomsökas en gång.
func (g *Game) search(p *Player) {
	room := &g.rooms[p.Y][p.X]
	if room.Kind != "rum" || room.Searched {
		return
	}
	room.Searched = true
	switch weighted(g.data.Sok.Vikter) {
	case "skatt":
		lo, hi := g.data.Regler.SokSkattMin, g.data.Regler.SokSkattMax
		n := lo + rand.IntN(max(1, hi-lo+1))
		p.Treasure += n
		g.soundTo(p, "skatt")
		g.logTo(p, g.data.SearchText("skatt", "antal", strconv.Itoa(n), "skatt", strconv.Itoa(p.Treasure)), "syn")
	case "dryck":
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
	case "falla":
		g.logTo(p, g.data.SearchText("falla"), "varn")
		g.trigger(p, g.data.randomTrap(), 0)
	default:
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
		if room.Doors[d] || !inBounds(nx, ny) {
			continue
		}
		if _, fixed := fixedDoors(nx, ny); fixed {
			continue
		}
		cand = append(cand, d)
	}
	if len(cand) == 0 {
		return 0, false
	}
	d := cand[rand.IntN(len(cand))]
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
	p.HP = min(g.data.Regler.Liv, p.HP+g.data.Regler.DryckHelar)
	g.soundTo(p, "drick")
	g.logTo(p, g.data.Text("drick", "liv", strconv.Itoa(p.HP)), "info")
	if p.Poison > 0 {
		p.Poison = 0
		g.logTo(p, g.data.Text("drick_gift"), "info")
	}
}
