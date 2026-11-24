package main

import (
	"errors"
	"log"
	"math/rand/v2"
)

// Fröet: samma frö och samma drag ger samma rum, fällor, varelser och
// utfall. Tornet vet alltid vilket frö nästa vandring får, så att
// spelledaren kan se det, skriva in ett eget eller slumpa ett nytt på
// värdsidan (avsnittet Fröet). -fro ger det första.

// maxSeed är det största fröet som går att skriva in.
const maxSeed = 999999999

// errBadSeed: fröet saknas eller är 0.
var errBadSeed = errors.New("inget frö")

// randomSeed ger ett slumpat frö, 1–maxSeed.
func randomSeed() uint64 {
	return 1 + rand.Uint64()%maxSeed
}

// reseed ger tornet en ny slump med fröet seed.
func (g *Game) reseed(seed uint64) {
	g.seed = seed
	g.rng = rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
}

// seedWandering tar nästa frö när en vandring börjar och slumpar ett nytt
// till vandringen efter.
func (g *Game) seedWandering() {
	seed := g.nextSeed
	if seed == 0 {
		seed = randomSeed()
	}
	g.nextSeed = randomSeed()
	g.reseed(seed)
	log.Printf("Vandringen börjar. Frö: %d (starta med -fro %d för samma torn)", seed, seed)
}

// SetNextSeed sätter fröet nästa vandring får, eller slumpar ett nytt när
// seed är 0. I lobbyn byggs tornet om med fröet, så att också portarna
// följer det. Ett frö utanför 1–maxSeed ger false.
func (g *Game) SetNextSeed(seed uint64) bool {
	if seed > maxSeed {
		return false
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if seed == 0 {
		seed = randomSeed()
	}
	changed := seed != g.nextSeed
	g.nextSeed = seed
	if g.phase == phaseEnd && changed {
		g.towerChanged() // bygger tornet med det nya fröet
	} else if g.phase == phaseLobby {
		g.reseed(seed)
		g.newTower()
	}
	g.broadcastState()
	return true
}

// towerForNext bygger ett nytt torn efter en vandring med nästa frö, så att
// också portarna blir desamma som när vandringen spelades förut.
func (g *Game) towerForNext() {
	g.reseed(g.nextSeed)
	g.newTower()
}

// towerChanged stänger slutskärmen när spelledaren ändrar förutsättningarna
// (fröet, AI-spelarna eller svårigheten) efter en vandring. Alla kommer till
// lobbyn med ett nytt torn, så att placeringen aldrig visar en vandring som
// inte längre stämmer.
func (g *Game) towerChanged() {
	if g.phase != phaseEnd {
		return
	}
	g.toLobby("torn_andrat")
}
