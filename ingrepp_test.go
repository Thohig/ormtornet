package main

import (
	"strings"
	"testing"
)

// gmPlaying ger en pågående vandring med Alva och Bo i spelledarläge.
func gmPlaying(t *testing.T) (*Game, *Client, *Client) {
	t.Helper()
	g, a, b := twoPlaying(t, 0)
	g.cfg.GMMode = true
	return g, a, b
}

func TestInterveneNeedsGMMode(t *testing.T) {
	g, _, _ := twoPlaying(t, 0)
	if err := g.Intervene(gmRequest{Typ: "paus"}); err != errNotGM {
		t.Fatalf("fel: %v", err)
	}
}

func TestPauseHoldsTheRound(t *testing.T) {
	g, a, b := gmPlaying(t)
	if err := g.Intervene(gmRequest{Typ: "paus"}); err != nil {
		t.Fatal(err)
	}
	g.Act(a, actWait, 0)
	g.Act(b, actWait, 0)
	if g.round != 1 {
		t.Fatal("rundan avgjordes fast den var pausad")
	}
	if st := g.stateFor(a.player); !st.Paused {
		t.Fatal("spelarna ser inte att rundan är pausad")
	}
	g.Intervene(gmRequest{Typ: "fortsatt"})
	if g.round != 2 {
		t.Fatal("rundan avgjordes inte när den fick fortsätta och alla hade valt")
	}
}

func TestMoveTheSun(t *testing.T) {
	g, _, _ := gmPlaying(t)
	g.Intervene(gmRequest{Typ: "sol", Antal: -3})
	if g.sun != g.cfg.SunMax-3 {
		t.Fatalf("solen %d", g.sun)
	}
	g.Intervene(gmRequest{Typ: "sol", Antal: 100})
	if g.sun != g.cfg.SunMax {
		t.Fatalf("solen steg över sin startpunkt: %d", g.sun)
	}
	g.Intervene(gmRequest{Typ: "sol", Antal: -100})
	if g.sun != 1 || g.phase != phasePlay {
		t.Fatalf("solen %d, fas %s: spelledaren ska inte kunna låta den gå ner", g.sun, g.phase)
	}
}

func TestTowerVoice(t *testing.T) {
	g, a, b := gmPlaying(t)
	for len(a.send) > 0 {
		<-a.send
	}
	for len(b.send) > 0 {
		<-b.send
	}
	if err := g.Intervene(gmRequest{Typ: "rost", Spelare: a.player.ID, Text: "Någon viskar ditt namn."}); err != nil {
		t.Fatal(err)
	}
	got := func(c *Client) bool {
		ok := false
		for len(c.send) > 0 {
			if strings.Contains(string(<-c.send), "Någon viskar ditt namn.") {
				ok = true
			}
		}
		return ok
	}
	if !got(a) || got(b) {
		t.Fatal("tornets röst nådde fel vandrare")
	}
	if err := g.Intervene(gmRequest{Typ: "rost", Text: strings.Repeat("a", voiceMaxLen+1)}); err == nil {
		t.Fatal("för lång text gick igenom")
	}
}

func TestPlaceTrapAndMonster(t *testing.T) {
	g, a, _ := gmPlaying(t)
	if err := g.Intervene(gmRequest{Typ: "falla", X: 4, Y: 4, ID: "pilar"}); err == nil {
		t.Fatal("lade en fälla i ett rum som inte finns än")
	}
	putInRoom(g, a.player, 4, 1)
	if err := g.Intervene(gmRequest{Typ: "falla", X: 4, Y: 1, ID: "pilar"}); err != nil || g.rooms[1][4].Event != "pilar" {
		t.Fatalf("fällan lades inte: %v", err)
	}
	typ := g.data.Monster.Typer[0].ID
	if err := g.Intervene(gmRequest{Typ: "varelse", X: 4, Y: 1, ID: typ}); err != nil || g.monsterAt(4, 1) == nil {
		t.Fatalf("varelsen släpptes inte in: %v", err)
	}
	if err := g.Intervene(gmRequest{Typ: "varelse", X: 4, Y: 1, ID: typ}); err == nil {
		t.Fatal("två varelser i samma rum")
	}
}

func TestCalmSerpent(t *testing.T) {
	g, _, _ := gmPlaying(t)
	g.deck = []string{"sover"}
	for _, c := range g.data.Ormen.Kort {
		if c.Vaknar {
			g.deck = append(g.deck, c.ID)
			break
		}
	}
	g.Intervene(gmRequest{Typ: "lugna"})
	if wakeCards(g) != 0 {
		t.Fatal("ormen lugnades inte")
	}
}

func TestGiveAndTake(t *testing.T) {
	g, a, _ := gmPlaying(t)
	alva := a.player
	alva.HP = 3
	alva.Treasure = 5
	g.Intervene(gmRequest{Typ: "liv", Spelare: alva.ID, Antal: 100})
	if alva.HP != g.data.Regler.Liv {
		t.Fatalf("liv %d, ska stanna vid %d", alva.HP, g.data.Regler.Liv)
	}
	g.Intervene(gmRequest{Typ: "guld", Spelare: alva.ID, Antal: -10})
	if alva.Treasure != 0 {
		t.Fatalf("guld %d", alva.Treasure)
	}
	g.Intervene(gmRequest{Typ: "dryck", Spelare: alva.ID, Antal: 1})
	if alva.Potions != 1 {
		t.Fatalf("drycker %d", alva.Potions)
	}
	g.Intervene(gmRequest{Typ: "liv", Spelare: alva.ID, Antal: -100})
	if alva.Status != statusDead {
		t.Fatal("Alva dog inte när allt liv togs")
	}
	if !strings.Contains(g.gmLog[len(g.gmLog)-1].Text, "Tog") {
		t.Fatalf("ingreppet loggades inte: %+v", g.gmLog[len(g.gmLog)-1])
	}
}
