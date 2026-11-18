package main

import "testing"

// fightReady ger två vandrare i samma vanliga rum, med strid påslagen.
func fightReady(t *testing.T) (*Game, *Client, *Client) {
	t.Helper()
	g, a, b := twoPlaying(t, 0)
	g.cfg.PvP = pvpOn
	putInRoom(g, a.player, 4, 1)
	putInRoom(g, b.player, 4, 1)
	g.deck = []string{"sover", "sover", "sover"}
	return g, a, b
}

func wakeCards(g *Game) int {
	n := 0
	for _, id := range g.deck {
		if c := g.data.card(id); c != nil && c.Vaknar {
			n++
		}
	}
	return n
}

func TestPvPExchangeStealsAThird(t *testing.T) {
	g, a, b := fightReady(t)
	alva, bo := a.player, b.player
	bo.Treasure = 30
	hp := bo.HP
	g.Act(a, actAttack, -bo.ID)
	g.Act(b, actDodge, -alva.ID) // anfall slår undanmanöver
	if bo.HP != hp-pvpSkada || bo.Treasure != 20 || alva.Treasure != 10 {
		t.Fatalf("Bo %d liv %d guld, Alva %d guld", bo.HP, bo.Treasure, alva.Treasure)
	}
	if wakeCards(g) != 1 {
		t.Fatalf("striden lade %d vakna-kort i leken, ville ha 1", wakeCards(g))
	}
}

func TestPvPAmbushAndKill(t *testing.T) {
	g, a, b := fightReady(t)
	alva, bo := a.player, b.player
	bo.Treasure = 50
	bo.HP = 1
	g.Act(a, actAttack, -bo.ID)
	g.Act(b, actWait, 0) // Bo försvarar sig inte
	if bo.Status != statusDead || bo.Treasure != 0 || alva.Treasure != 50 {
		t.Fatalf("Bo %s med %d guld, Alva %d guld", bo.Status, bo.Treasure, alva.Treasure)
	}
}

func TestPvPParryAgainstNoAttack(t *testing.T) {
	g, a, b := fightReady(t)
	bo := b.player
	hp := bo.HP
	g.Act(a, actParry, -bo.ID)
	g.Act(b, actWait, 0)
	if bo.HP != hp {
		t.Fatal("en parering skadade någon som inte anföll")
	}
}

func TestPvPRules(t *testing.T) {
	g, a, _ := fightReady(t)
	alva := a.player
	try := func() bool {
		alva.Action = ""
		g.actPlayer(alva, actAttack, 0)
		return alva.Action == actAttack
	}
	if !try() {
		t.Fatal("kunde inte anfalla med strid påslagen")
	}
	g.cfg.PvP = pvpOff
	if try() {
		t.Fatal("anföll fast strid var avstängd")
	}
	g.cfg.PvP = pvpHalf
	g.sun = g.cfg.SunMax
	if try() {
		t.Fatal("anföll innan solen sjunkit halvvägs")
	}
	g.sun = g.cfg.SunMax / 2
	if !try() {
		t.Fatal("kunde inte anfalla när solen sjunkit halvvägs")
	}
	g.cfg.PvP = pvpOn
	g.rooms[alva.Y][alva.X].Kind = "ingang"
	if try() {
		t.Fatal("anföll i ett fredat portrum")
	}
}

func TestBotsFightEachOtherOnlyWithoutHumans(t *testing.T) {
	g := testGame(t)
	g.cfg.PvP = pvpOn
	g.cfg.AI = 2
	g.AddBots(2)
	var bots []*Player
	for _, p := range g.players {
		if p != nil && p.bot != nil {
			bots = append(bots, p)
			p.Status = statusIn
			putInRoom(g, p, 4, 1)
		}
	}
	h := testClient()
	g.Join(h, "Alva")
	if g.botFoe(bots[0]) != nil {
		t.Fatal("en AI-spelare gav sig på en annan AI-spelare fast en människa var med")
	}
	g.Leave(h)
	if g.botFoe(bots[0]) != bots[1] {
		t.Fatal("utan människor ska AI-spelarna kunna slåss med varandra")
	}
}
