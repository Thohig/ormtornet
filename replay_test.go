package main

import "testing"

// Reprisen: var alla stod när vandringen började, efter varje runda och
// när den var slut, och tornet som det blev. Den kommer bara i slutet.
func TestReplay(t *testing.T) {
	g, a, b := startTwo(t)
	if len(g.replayFrames) != 1 || len(g.replayWho) != 2 {
		t.Fatalf("i början: %d rundor, %d vandrare", len(g.replayFrames), len(g.replayWho))
	}
	if st := g.stateFor(a.player); st.Replay != nil {
		t.Fatal("reprisen syns under vandringen")
	}
	for range 3 {
		step(g, a, b, actWait, actWait)
	}
	if len(g.replayFrames) != 4 {
		t.Fatalf("efter tre rundor: %d", len(g.replayFrames))
	}
	b.player.Treasure = 17
	g.endGame()
	v := g.stateFor(a.player).Replay
	if v == nil || len(v.Frames) != 5 || len(v.People) != 2 || len(v.Rooms) == 0 {
		t.Fatalf("reprisen: %+v", v)
	}
	last := v.Frames[len(v.Frames)-1]
	if last[1].Gold != 17 || last[0].Status != statusLost {
		t.Fatalf("sista rundan: %+v", last)
	}
	if v.People[0].Name != a.player.Name || v.Size != g.size {
		t.Fatalf("vandrarna: %+v", v.People)
	}
}
