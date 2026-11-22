package main

import "time"

// Frånvarande vandrare: när en människa har lämnat mitt i en vandring kan en
// AI-spelare ta över vandraren, med dess plats, skatt, liv och tillstånd.
// Kommer människan tillbaka med samma namn får den genast tillbaka sina steg.
// Detsamma gäller den som är kvar men står still, utan att välja något drag,
// i lika många rundor; den får tillbaka sina steg med nästa drag.

// canTakeOver är sant för en vandrare vars människa är borta och som
// fortfarande har något att göra i tornet.
func (g *Game) canTakeOver(p *Player) bool {
	return g.phase == phasePlay && p != nil && p.client == nil && p.bot == nil &&
		p.Status == statusIn
}

func (g *Game) takeOver(p *Player) {
	pers := g.data.AI.Personligheter
	if len(pers) == 0 {
		return
	}
	p.bot = &botBrain{pers: &pers[p.Slot%len(pers)], habits: map[string]map[string]int{}}
	p.takenOver = true
	p.aiPlayed = true
	p.still = 0
	if len(g.data.AI.TarOver) > 0 {
		g.logOthers(p, fill(pick(g.data.AI.TarOver), "namn", p.Name), "info")
	}
	if p.client != nil {
		g.logTo(p, g.data.Text("stilla_tar_over"), "syn")
	}
	g.scheduleBot(p)
}

// autoTakeOver körs när en runda börjar och låter AI:n ta över dem som har
// varit borta, eller stått still, i så många hela rundor som spelledaren har valt.
func (g *Game) autoTakeOver() {
	if g.cfg.AbsentRounds <= 0 {
		return
	}
	for _, p := range g.players {
		if g.canTakeOver(p) && p.awayRound > 0 && g.round-p.awayRound >= g.cfg.AbsentRounds {
			g.takeOver(p)
		} else if g.standsStill(p) {
			g.takeOver(p)
		}
	}
}

// countStill körs när rundan avgörs och räknar hur många rundor i rad varje
// ansluten människa har låtit tiden gå ut utan att välja något drag. Den som
// ligger i källaren kan inte välja och räknas inte.
func (g *Game) countStill(active []*Player) {
	for _, p := range active {
		switch {
		case p.client == nil || p.bot != nil || p.Cellar > 0:
		case p.Action == "":
			p.still++
		default:
			p.still = 0
		}
	}
}

// standsStill är sant för en ansluten människa som har stått still så länge
// att AI:n ska ta över.
func (g *Game) standsStill(p *Player) bool {
	return g.phase == phasePlay && p != nil && p.client != nil && p.bot == nil &&
		p.Status == statusIn && p.still >= g.cfg.AbsentRounds
}

// reclaim ger tillbaka vandraren till människan som stod still, när hen
// själv väljer ett drag.
func (g *Game) reclaim(p *Player) {
	p.bot = nil
	p.takenOver = false
	p.still = 0
	g.logTo(p, g.data.Text("stilla_tillbaka"), "syn")
	if len(g.data.AI.Tillbaka) > 0 {
		g.logOthers(p, fill(pick(g.data.AI.Tillbaka), "namn", p.Name), "info")
	}
	g.broadcastState()
}

// TakeOverSlot låter spelledaren be AI:n ta över en frånvarande vandrare direkt.
func (g *Game) TakeOverSlot(slot int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if slot < 0 || slot >= maxPlayers || !g.canTakeOver(g.players[slot]) {
		return false
	}
	g.takeOver(g.players[slot])
	g.broadcastState()
	return true
}

// handBack ger tillbaka vandraren till människan som kom tillbaka.
func (g *Game) handBack(p *Player, c *Client) {
	p.bot = nil
	p.takenOver = false
	p.awayRound, p.still = 0, 0
	p.client, c.player = c, p
	g.checkAwait()
	g.broadcastState()
	g.logTo(p, g.data.Text("tillbaka"), "syn")
	if len(g.data.AI.Tillbaka) > 0 {
		g.logOthers(p, fill(pick(g.data.AI.Tillbaka), "namn", p.Name), "info")
	}
}

// releaseTakenOver körs när vandringen är slut: AI:n släpper de vandrare den
// har tagit över, så att platserna blir lediga om människorna inte kommer tillbaka.
func (g *Game) releaseTakenOver() {
	for _, p := range g.players {
		if p != nil && p.takenOver {
			p.bot = nil
			p.takenOver = false
		}
	}
}

// dropAbsent tar bort människor som lämnade under förra vandringen och inte
// har kommit tillbaka, innan en ny vandring börjar. Den säger hur många.
func (g *Game) dropAbsent() int {
	n := 0
	for i, p := range g.players {
		if p != nil && p.bot == nil && p.client == nil {
			g.players[i] = nil
			n++
		}
	}
	return n
}

// shouldAwait är sant när ingen människa är ansluten men någon människas
// vandrare fortfarande är kvar i tornet. Utan någon som ser på skulle
// rundorna annars gå så fort AI:n hinner välja. En vandring med bara
// AI-spelare, som spelledaren kan leda, går som vanligt.
func (g *Game) shouldAwait() bool {
	if g.phase != phasePlay {
		return false
	}
	waiting := false
	for _, p := range g.players {
		if p == nil {
			continue
		}
		if p.client != nil {
			return false
		}
		if (p.bot == nil || p.takenOver) && p.Status == statusIn {
			waiting = true
		}
	}
	return waiting
}

// checkAwait pausar rundan när den sista människan har gått, och låter den
// fortsätta när någon kommer tillbaka.
func (g *Game) checkAwait() {
	switch should := g.shouldAwait(); {
	case should && !g.paused:
		g.paused = true
		g.pauseLeft = max(time.Second, time.Until(g.deadline))
		if g.timer != nil {
			g.timer.Stop()
		}
	case !should && g.paused:
		g.paused = false
		g.armTimer(g.pauseLeft)
	}
}
