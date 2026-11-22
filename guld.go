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

// heardDeath är det en vandrare hörde när någon dog: var, så noga som hen
// kunde höra det, och om mynt rasslade. AI-spelarna går efter det (ai.go);
// en människa läser samma sak i loggen.
type heardDeath struct {
	x, y   int    // rummet, om vandraren stod där eller bredvid; annars -1
	region string // tornets del
	coins  bool
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
		heard := heardDeath{x: -1, y: -1, region: region, coins: coins}
		var text string
		d, next := g.nextTo(o, p.X, p.Y)
		switch {
		case o.Status == statusIn && o.Cellar == 0 && o.X == p.X && o.Y == p.Y:
			text = g.data.Text("dod_andra", "namn", p.Name)
			heard.x, heard.y = p.X, p.Y
		case next && !g.is(o, "sparfinnare"):
			text = g.data.Text("ljud_dod", "namn", p.Name, "riktning", g.data.Riktningar[d])
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
