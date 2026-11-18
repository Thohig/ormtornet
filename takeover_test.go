package main

import "testing"

// twoPlaying startar en vandring med Alva och Bo.
func twoPlaying(t *testing.T, absentRounds int) (*Game, *Client, *Client) {
	t.Helper()
	g := testGame(t)
	g.cfg.AbsentRounds = absentRounds
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	g.Ready(a)
	g.Ready(b)
	if g.phase != phasePlay {
		t.Fatal("vandringen startade inte")
	}
	return g, a, b
}

func TestAITakesOverAfterRounds(t *testing.T) {
	g, a, b := twoPlaying(t, 2)
	bo := b.player
	g.Leave(b) // Bo går i runda 1
	g.Act(a, actWait, 0)
	if g.round != 2 || bo.takenOver {
		t.Fatalf("runda %d, övertagen %v: för tidigt", g.round, bo.takenOver)
	}
	g.Act(a, actWait, 0)
	if g.round != 3 || !bo.takenOver || bo.bot == nil {
		t.Fatalf("runda %d: AI:n tog inte över Bo efter två rundor", g.round)
	}

	// Bo kommer tillbaka och får sina steg.
	b2 := testClient()
	g.Join(b2, "Bo")
	if b2.player != bo || bo.bot != nil || bo.takenOver || bo.client != b2 {
		t.Fatal("Bo fick inte tillbaka sin vandrare")
	}
}

func TestNoTakeOverWhenOff(t *testing.T) {
	g, a, b := twoPlaying(t, 0)
	bo := b.player
	g.Leave(b)
	for range 4 {
		g.Act(a, actWait, 0)
	}
	if bo.takenOver {
		t.Fatal("AI:n tog över fast inställningen var av")
	}
}

func TestHostTakeOverNow(t *testing.T) {
	g, a, b := twoPlaying(t, 0)
	if g.TakeOverSlot(a.player.Slot) {
		t.Fatal("AI:n tog över en spelare som är här")
	}
	bo := b.player
	g.Leave(b)
	if !g.TakeOverSlot(bo.Slot) || !bo.takenOver {
		t.Fatal("spelledaren kunde inte låta AI:n ta över")
	}
	if g.TakeOverSlot(bo.Slot) {
		t.Fatal("AI:n tog över två gånger")
	}
}

func TestTakenOverAtEndAndNextTower(t *testing.T) {
	g, _, b := twoPlaying(t, 0)
	bo := b.player
	g.Leave(b)
	g.TakeOverSlot(bo.Slot)
	g.endGame()
	if bo.takenOver || bo.bot != nil {
		t.Fatal("AI:n släppte inte Bo när vandringen var slut")
	}
	found := false
	for _, r := range g.results() {
		if r.Name == "Bo" && r.AI {
			found = true
		}
	}
	if !found {
		t.Fatal("resultatet visar inte att AI:n spelade åt Bo")
	}
	if !g.ForceStart() {
		t.Fatal("nästa vandring startade inte")
	}
	for _, p := range g.players {
		if p != nil && p.Name == "Bo" {
			t.Fatal("Bo, som aldrig kom tillbaka, är med i nästa vandring")
		}
	}
}

// En AI-spelare som har planerat ett drag får inte krascha tornet när
// människan kommer tillbaka innan draget hinner väljas.
func TestHandBackBeforeBotActs(t *testing.T) {
	g, _, b := twoPlaying(t, 0)
	bo := b.player
	g.Leave(b)
	if !g.TakeOverSlot(bo.Slot) {
		t.Fatal("AI:n tog inte över Bo")
	}
	round := g.round
	g.Join(testClient(), "Bo")
	if bo.bot != nil {
		t.Fatal("Bo fick inte tillbaka sin vandrare")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.botAct(bo, round) // timern från övertagandet löser ut
	if bo.Action != "" {
		t.Fatalf("en AI som inte finns valde %q åt Bo", bo.Action)
	}
}
