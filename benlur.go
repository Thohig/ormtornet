package main

// Benluren: ett krökt horn av ben. Den som bär den hör om något andas i
// rummen bredvid när hen kommer in i ett rum, och hör ljuden från rummen
// två steg bort, med riktning, som andra bara hör från grannrummen. Med
// benluren hör också spårfinnaren genom väggarna.

// hearsFrom säger om o hör vad som händer i rummet (x, y) med riktning, och
// åt vilket håll: från ett grannrum, eller två rum bort med benluren.
func (g *Game) hearsFrom(o *Player, x, y int) (string, bool) {
	horn := o.has("benlur") != nil
	if d, ok := g.nextTo(o, x, y); ok && (horn || !g.is(o, "sparfinnare")) {
		return g.data.Riktningar[d], true
	}
	if !horn || o.Status != statusIn || o.Cellar > 0 || abs(x-o.X)+abs(y-o.Y) != 2 {
		return "", false
	}
	return g.hornToward(o.X, o.Y, x, y), true
}

// hornToward är väderstrecket från (fx, fy) mot (tx, ty), som ett av orden
// under "horn" i texter.json.
func (g *Game) hornToward(fx, fy, tx, ty int) string {
	// Samma ordning som "horn": nordväst, nordöst, sydväst, sydöst, norr, söder, väster, öster.
	idx := map[[2]int]int{{-1, -1}: 0, {1, -1}: 1, {-1, 1}: 2, {1, 1}: 3, {0, -1}: 4, {0, 1}: 5, {-1, 0}: 6, {1, 0}: 7}[[2]int{sign(tx - fx), sign(ty - fy)}]
	return g.data.Horn[idx]
}

// hornListen: den som bär benluren och kommer in i ett rum hör om en
// varelse står i något av rummen bredvid, men inte vad det är.
func (g *Game) hornListen(p *Player) {
	if p.has("benlur") == nil || p.Status != statusIn || p.Cellar > 0 {
		return
	}
	for d := North; d <= West; d++ {
		nx, ny := p.X+dirDX[d], p.Y+dirDY[d]
		if g.inBounds(nx, ny) && g.monsterAt(nx, ny) != nil {
			g.hearTo(p, g.data.Text("benlur_varelse", "riktning", g.data.Riktningar[d]), "info")
		}
	}
}
