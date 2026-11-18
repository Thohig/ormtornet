package main

import (
	"strings"
	"testing"
)

func TestWhisperMatching(t *testing.T) {
	answers := []string{"tornet är ormen", "inuti ormen"}
	for _, c := range []struct {
		w  string
		ok bool
	}{
		{"Tornet ÄR ormen!", true},
		{"jag tror att tornet är ormen", true},
		{"vi är inuti ormen", true},
		{"ormen sover i tornet", false},
		{"tornet är ormens", false},
		{"", false},
	} {
		if got := whisperMatches(c.w, answers); got != c.ok {
			t.Errorf("%q: %v", c.w, got)
		}
	}
}

func TestWhisperInitiates(t *testing.T) {
	g, a, b := startTwo(t)
	pa := a.player
	g.Whisper(a, "tornet är ormen")
	if pa.initiated {
		t.Fatal("invigd på nivå 1")
	}
	g.level = 2
	logs(b)
	g.Whisper(a, "tornet är ett torn")
	if pa.initiated || pa.guessLevel != 2 {
		t.Fatal("fel viskning gick igenom")
	}
	g.Whisper(a, "tornet är ormen")
	if pa.initiated {
		t.Fatal("två viskningar på samma nivå")
	}
	g.level = 3
	g.Whisper(a, "Tornet är ormen.")
	if !pa.initiated {
		t.Fatal("rätt viskning gick inte igenom")
	}
	for _, line := range logs(b) {
		if strings.Contains(line, "ormen") || strings.Contains(line, pa.Name) {
			t.Fatalf("Bo fick veta något: %q", line)
		}
	}
	if st := g.stateFor(pa); !st.Initiated || st.Heart == "" || st.CanWhisper {
		t.Fatalf("klienten: %+v", st.Initiated)
	}
	// Utan avslöjande går det inte att viska.
	g2, c, _ := startTwo(t)
	g2.cfg.Reveal = false
	g2.level = 3
	g2.Whisper(c, "tornet är ormen")
	if c.player.initiated {
		t.Fatal("invigd utan avslöjande")
	}
}

func TestHeartbeatPointsToChamber(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.initiated = true
	putInRoom(g, p, 1, 1)
	if want := g.data.Text("hjartslag", "horn", "sydöst"); g.heartbeat(p) != want {
		t.Fatalf("%q", g.heartbeat(p))
	}
	p.X, p.Y = g.center(), g.center()
	if g.heartbeat(p) != g.data.Text("hjartslag_har") {
		t.Fatal("inget dunk under kammaren")
	}
}

func TestPortWaitsForInitiated(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	pa.initiated = true
	putInRoom(g, pa, 4, 3)
	putInRoom(g, pb, 2, 3)
	g.sun = 1
	g.resolve()
	if g.phase != phasePlay || pa.Status != statusIn || pb.Status != statusLost {
		t.Fatalf("fas %s, Alva %s, Bo %s", g.phase, pa.Status, pb.Status)
	}
	g.resolve()
	if g.phase != phaseEnd || pa.Status != statusLost {
		t.Fatal("porten väntade mer än en runda")
	}
}
