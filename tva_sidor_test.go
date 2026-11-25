package main

import "testing"

// Ormens öga pekar mot kammaren och vänder på nivå 3: sedan lägger varje
// näve ett vakna-kort i leken.
func TestSerpentsEye(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	putInRoom(g, p, 0, 0)
	g.give(p, "ormens_oga")
	cx, cy := g.chamber()
	if got, want := g.eye(p), g.data.Text("oga_pekar", "horn", g.hornToward(0, 0, cx, cy)); got != want {
		t.Fatalf("ögat: %q, ville ha %q", got, want)
	}
	if g.stateFor(p).Eye == "" {
		t.Fatal("klienten ser inte ögat")
	}
	g.level = 2
	g.ageItems(p)
	if p.has("ormens_oga").turned {
		t.Fatal("ögat vände på nivå 2")
	}
	g.level = 3
	g.ageItems(p)
	if !p.has("ormens_oga").turned || !g.eyeWakes(p) {
		t.Fatal("ögat vände inte på nivå 3")
	}

	// I kammaren lägger näven ett vakna-kort i leken.
	p.X, p.Y = cx, cy
	g.reveal(cx, cy)
	for i, id := range g.deck {
		if c := g.data.card(id); c != nil && !c.Vaknar {
			g.deck[0], g.deck[i] = g.deck[i], g.deck[0]
			break
		}
	}
	g.deck = g.deck[:1]
	g.take(p)
	if g.wakeCardsLeft() != 1 {
		t.Fatalf("%d vakna-kort efter en näve med ögat", g.wakeCardsLeft())
	}
}

// Ägget är värt guld, men kläcks efter tolv rundor.
func TestEggHatches(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	putInRoom(g, p, 4, 3)
	g.give(p, "agg")
	if g.itemGold(p) != 50 {
		t.Fatalf("ägget är värt %d", g.itemGold(p))
	}
	for range 11 {
		g.ageItems(p)
	}
	if p.has("agg") == nil || g.monsterAt(4, 3) != nil {
		t.Fatal("ägget kläcktes för tidigt")
	}
	g.ageItems(p)
	m := g.monsterAt(4, 3)
	if p.has("agg") != nil || m == nil || m.Type != "yngel" {
		t.Fatal("ägget kläcktes inte till ett yngel i rummet")
	}
}

// Kniven som minns ger Kraft för varje dödande, högst tre, och gör
// varelserna lynnigare.
func TestKnifeRemembers(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	putInRoom(g, p, 4, 3)
	g.give(p, "kniv_minns")
	typ := g.data.monsterType("yngel")
	m := &Monster{ID: 70, Type: typ.ID}
	k0, _ := g.itemBonus(p)
	l0 := g.lynne(m, p)
	for i := range 5 {
		g.addMonster(typ, 4, 3)
		g.defeat(p, g.monsterAt(4, 3))
		if k, _ := g.itemBonus(p); k != k0+min(i+1, knifeMax) {
			t.Fatalf("efter %d dödanden Kraft +%d", i+1, k-k0)
		}
	}
	if !p.has("kniv_minns").turned {
		t.Fatal("kniven vände inte")
	}
	if got := g.lynne(m, p); got != min(100, l0+30) {
		t.Fatalf("lynnet %d med kniven, %d utan", got, l0)
	}
}

// Glödmasken håller mörkret borta och visar varelser i grannrummen, och
// efter tio rundor syns den som bär den från grannrummet.
func TestGlowWorm(t *testing.T) {
	g, a, b := startTwo(t)
	p, o := a.player, b.player
	openRoom(g, p, 4, 3)
	g.give(p, "glodmask")
	g.trigger(p, g.data.trap("morker"), 0)
	if p.Dark != 0 {
		t.Fatal("mörkret bet på glödmasken")
	}
	g.addMonster(g.data.monsterType("yngel"), 5, 3)
	if !g.seesMonsterAt(p, 5, 3) {
		t.Fatal("glödmasken visar inte varelsen bredvid")
	}
	putInRoom(g, o, 3, 3)
	if g.nearLantern(o, p) {
		t.Fatal("glödmasken lyste genom väskan för tidigt")
	}
	for range 10 {
		g.ageItems(p)
	}
	if !g.nearLantern(o, p) {
		t.Fatal("glödmasken lyser inte genom väskan efter tio rundor")
	}
}
