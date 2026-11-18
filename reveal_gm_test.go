package main

import "testing"

func gmGame(t *testing.T) (*Game, *Client, *Client) {
	t.Helper()
	g, a, b := startTwo(t)
	g.cfg.GMMode = true
	return g, a, b
}

func TestGMLevelUpAndHold(t *testing.T) {
	g, a, _ := gmGame(t)
	if g.level != 1 {
		t.Fatalf("nivå %d", g.level)
	}
	logs(a)
	if err := g.Intervene(gmRequest{Typ: "niva_upp"}); err != nil {
		t.Fatal(err)
	}
	if g.level != 2 || len(logs(a)) == 0 {
		t.Fatalf("nivå %d, eller inga nya texter", g.level)
	}

	// Håll kvar: solen sjunker långt, men nivån står still.
	if err := g.Intervene(gmRequest{Typ: "niva_hall"}); err != nil {
		t.Fatal(err)
	}
	g.sun = 1
	g.updateLevel()
	if g.level != 2 {
		t.Fatalf("nivån steg till %d fast den hölls kvar", g.level)
	}
	if v := g.GMView(false); !v.Held || v.SunLvl != 4 {
		t.Fatalf("översikten: hålls %v, solen ger %d", v.Held, v.SunLvl)
	}
	// Nästa nivå går ändå att ta för hand.
	g.Intervene(gmRequest{Typ: "niva_upp"})
	if g.level != 3 {
		t.Fatalf("nivå %d", g.level)
	}
	// Släpp: nivån hinner ikapp solen efter rundan.
	g.Intervene(gmRequest{Typ: "niva_slapp"})
	g.updateLevel()
	if g.level != 4 {
		t.Fatalf("nivå %d efter släpp", g.level)
	}
	if err := g.Intervene(gmRequest{Typ: "niva_upp"}); err == nil {
		t.Error("kunde gå förbi nivå 4")
	}

	// Utan avslöjande går det inte.
	g2, _, _ := gmGame(t)
	g2.cfg.Reveal = false
	if err := g2.Intervene(gmRequest{Typ: "niva_upp"}); err != errNoReveal {
		t.Errorf("fel %v", err)
	}
}

func TestWakeRisk(t *testing.T) {
	g, _, _ := gmGame(t)
	g.deck = []string{"sover", "sover", "sover", "vaknar"}
	if r := g.wakeRisk(); r != 25 {
		t.Fatalf("risk %d", r)
	}
	if v := g.GMView(false); v.Serpent.Risk != 25 {
		t.Fatalf("översikten visar %d", v.Serpent.Risk)
	}
	// En tom lek blandas om, och leken i spelet är orörd.
	g.deck = nil
	sleep := g.deckSleep
	if r := g.wakeRisk(); r != 13 || len(g.deck) != 0 || g.deckSleep != sleep {
		t.Fatalf("risk %d, lek %d, sovkort %d", r, len(g.deck), g.deckSleep)
	}
}
