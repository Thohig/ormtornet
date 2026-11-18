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

func TestInitiatedBenefits(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.initiated = true
	g.data.Avslojande.Invigd.Falla = 1
	// Fällan i ett nytt rum: den invigde kliver förbi.
	putInRoom(g, p, 4, 3)
	g.rooms[3][4].Doors[East] = true
	g.rooms[3][5] = Room{}
	g.reveal(5, 3).Event = "pilar"
	g.rooms[3][5].Revealed = false
	g.rooms[3][5].Event = "pilar"
	hp := p.HP
	g.move(p, East, 0)
	if p.HP != hp || g.rooms[3][5].Mark == "" {
		t.Fatalf("liv %d av %d, tecken %q", p.HP, hp, g.rooms[3][5].Mark)
	}

	// Lugna ormen i kammaren, en gång.
	p.X, p.Y = g.center(), g.center()
	g.reveal(p.X, p.Y)
	before := g.wakeCardsLeft()
	if !g.canWhisper(p) {
		t.Fatal("kan inte viska i kammaren")
	}
	g.Whisper(a, "sov")
	if g.wakeCardsLeft() != before-1 || !p.calmed || g.canWhisper(p) {
		t.Fatalf("vakna-kort %d, före %d", g.wakeCardsLeft(), before)
	}

	// Guld när den invigde kommer ut.
	cn := g.corners()[0]
	p.X, p.Y, p.Treasure = cn.x, cn.y, 100
	g.reveal(cn.x, cn.y)
	g.move(p, cn.exit, 0)
	if p.Status != statusOut || p.Treasure != 125 {
		t.Fatalf("%s med %d", p.Status, p.Treasure)
	}
}

func TestInitiatedWakesSerpentMoreEasily(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.initiated = true
	g.data.Avslojande.Invigd.Vacker = 1
	p.X, p.Y = g.center(), g.center()
	g.reveal(p.X, p.Y)
	wakes := g.wakes
	g.take(p)
	if g.wakes != wakes+1 {
		t.Fatal("ormen vaknade inte")
	}
}

func TestGMInitiateAndGiveItem(t *testing.T) {
	g, a, _ := gmGame(t)
	p := a.player
	if err := g.Intervene(gmRequest{Typ: "invig", Spelare: p.ID}); err != nil || !p.initiated {
		t.Fatalf("inviga: %v", err)
	}
	if err := g.Intervene(gmRequest{Typ: "invig", Spelare: p.ID}); err != nil || p.initiated {
		t.Fatalf("ta tillbaka: %v", err)
	}
	if err := g.Intervene(gmRequest{Typ: "foremal", Spelare: p.ID, ID: "rep"}); err != nil || p.has("rep") == nil {
		t.Fatalf("ge föremål: %v", err)
	}
	if err := g.Intervene(gmRequest{Typ: "foremal", Spelare: p.ID, ID: "finns_inte"}); err == nil {
		t.Fatal("gav ett föremål som inte finns")
	}
	if v := g.GMView(false); len(v.Items) == 0 || !v.Map.Reveal {
		t.Fatal("översikten saknar föremålen")
	}
}
