package main

// Läs om innehållet: spelledaren kan läsa in data/ igen utan att starta om
// tornet. Filerna kontrolleras direkt. I lobbyn gäller de genast och tornet
// byggs om; under och efter en vandring gäller de från nästa torn, så att
// ingen vandring byter innehåll mitt i.

// ReloadData tar emot nyss inlästa datafiler. Svarar true om de gäller
// redan nu, false om de väntar till nästa torn.
func (g *Game) ReloadData(d *Data) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.phase != phaseLobby {
		g.nextData = d
		return false
	}
	g.useData(d)
	for _, p := range g.players {
		if p != nil {
			p.HP = g.maxHP(p)
		}
	}
	g.sunMax = g.scaledSun()
	g.sun = g.sunMax
	g.newTower()
	g.broadcastState()
	return true
}

// useData byter till datafilerna d. AI-spelarna behåller sin personlighet
// om den finns kvar.
func (g *Game) useData(d *Data) {
	g.nextData = nil
	g.base = d
	g.cfg.Difficulty = d.difficultyID(g.cfg.Difficulty)
	g.data = d.withDifficulty(g.cfg.Difficulty)
	pers := g.data.AI.Personligheter
	for _, p := range g.players {
		if p == nil || p.bot == nil || len(pers) == 0 {
			continue
		}
		i := p.Slot % len(pers)
		if p.bot.pers != nil {
			for j := range pers {
				if pers[j].ID == p.bot.pers.ID {
					i = j
				}
			}
		}
		p.bot.pers = &pers[i]
	}
}
