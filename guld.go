package main

// Den dödes guld: guldet en vandrare bar när hen dog ligger kvar i rummet,
// och den som kommer dit, eller står kvar där, tar det. Föremålen ligger
// redan kvar (dropAll i items.go).

import "strconv"

// leftGold är guld som en död vandrare bar.
type leftGold struct {
	by   int // den döda vandrarens id
	name string
	gold int
}

// dropGold lägger guldet den döda bar i rummet hen dog i.
func (g *Game) dropGold(p *Player) {
	if p.Treasure <= 0 || !g.inBounds(p.X, p.Y) {
		return
	}
	r := &g.rooms[p.Y][p.X]
	r.gold = append(r.gold, leftGold{by: p.ID, name: p.Name, gold: p.Treasure})
	g.logTo(p, g.data.Text("guld_kvar_du", "guld", strconv.Itoa(p.Treasure)), "dim")
	p.Treasure = 0
}

// pickUpGold låter en levande vandrare ta guldet som ligger i rummet.
func (g *Game) pickUpGold(p *Player) {
	if p.Status != statusIn || p.Cellar > 0 || !g.inBounds(p.X, p.Y) {
		return
	}
	r := &g.rooms[p.Y][p.X]
	for _, l := range r.gold {
		p.Treasure += l.gold
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
	r.gold = nil
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
