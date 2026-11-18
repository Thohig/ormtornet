package main

// Spår och ljud genom väggarna. När vandrarna inte ser varandra kan de
// ändå ana var de andra är: färska fotspår i ett rum, och ljud från
// grannrummen med riktning men utan namn.

import "slices"

// hear låter de andra vandrarna höra något som händer i rummet (x, y). Med
// dold synlighet hör den som står i ett grannrum texten key med riktningen
// dit. Alla andra får general, om den inte är tom.
func (g *Game) hear(x, y int, key, general, class string, skip ...*Player) {
	for _, o := range g.players {
		if o == nil || slices.Contains(skip, o) {
			continue
		}
		if d, ok := g.nextTo(o, x, y); ok && g.hidesOthers() {
			g.logTo(o, g.data.Text(key, "riktning", g.data.Riktningar[d]), class)
		} else if general != "" {
			g.logTo(o, general, class)
		}
	}
}

// nextTo ger riktningen från o till rummet (x, y) när o står i ett grannrum.
func (g *Game) nextTo(o *Player, x, y int) (Dir, bool) {
	if o.Status != statusIn || o.Cellar > 0 {
		return 0, false
	}
	for d := North; d <= West; d++ {
		if o.X+dirDX[d] == x && o.Y+dirDY[d] == y {
			return d, true
		}
	}
	return 0, false
}

// noticeTrail berättar för p om färska spår i rummet den just kom in i.
func (g *Game) noticeTrail(p *Player) {
	if !g.hidesOthers() || p.Status != statusIn || p.Cellar > 0 || p.Dark > 0 {
		return
	}
	room := &g.rooms[p.Y][p.X]
	t := room.trail
	rounds := g.data.Regler.SparRundor
	if t.by == 0 || t.by == p.ID || rounds <= 0 || g.round-t.round > rounds {
		return
	}
	for _, o := range g.players {
		if o != nil && o.ID == t.by && o.Status == statusIn && o.X == p.X && o.Y == p.Y {
			return // den som lämnade spåren står här, och då ser man den
		}
	}
	if room.HasExit && room.Exit == t.dir {
		g.logTo(p, g.data.Text("spar_ut"), "dim")
		return
	}
	g.logTo(p, g.data.Text("spar", "riktning", g.data.Riktningar[t.dir]), "dim")
}
