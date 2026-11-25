package main

import "testing"

// Två föremål som hör ihop blir ett, och en plats i väskan blir ledig.
// Det går även när väskan är full.
func TestCombineFreesASlot(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	putInRoom(g, p, 4, 3)
	for _, id := range []string{"svard", "brynja", "mynt_hal"} {
		g.give(p, id)
	}
	g.data.Foremal.Lista = append(g.data.Foremal.Lista[:0:0], g.data.Foremal.Lista...)
	g.rooms[3][4].items = []string{"rep"}
	logs(a)
	g.pickUp(p)
	if len(p.items) != 3 || p.has("pendel") == nil || p.has("rep") != nil {
		t.Fatalf("väskan har %v", g.itemIDs(p))
	}
	if !containsPart(logs(a), pick(g.data.itemType("pendel").Fogas)) {
		t.Fatal("ingen rad om att de fogades ihop")
	}
	// Myntet verkar vidare i pendeln, med sin förbannelse.
	if p.has("mynt_hal") == nil {
		t.Fatal("myntet räknas inte i pendeln")
	}
	if g.itemGold(p) != 25 {
		t.Fatalf("pendeln är värd %d", g.itemGold(p))
	}
}

func TestCombineKeepsShownCurse(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	g.give(p, "rostkrona")
	g.curse(p, "rostkrona")
	g.give(p, "svart_parla")
	g.combine(p)
	it := p.has("krona_parla")
	if it == nil || len(p.items) != 1 {
		t.Fatalf("väskan har %v", g.itemIDs(p))
	}
	if !it.cursed["rostkrona"] || it.cursed["svart_parla"] {
		t.Fatalf("förbannelserna %v", it.cursed)
	}
	if p.has("svart_parla") != it || p.has("rostkrona") != it {
		t.Fatal("pärlan och kronan verkar inte vidare i kronan")
	}
	if g.itemGold(p) != 100 {
		t.Fatalf("kronan med pärlan är värd %d", g.itemGold(p))
	}
}

func TestPendulumPointsToGold(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	putInRoom(g, p, 4, 3)
	g.give(p, "rep")
	g.give(p, "mynt_hal")
	g.combine(p)
	if got := g.pendulum(p); got != g.data.Text("pendel_stilla") {
		t.Fatalf("pendeln utan guld: %q", got)
	}
	g.spillGold(4, 6, 30, "guld_lost")
	if got, want := g.pendulum(p), g.data.Text("pendel_pekar", "horn", g.data.Horn[5]); got != want {
		t.Fatalf("pendeln: %q, ville ha %q", got, want)
	}
	if st := g.stateFor(p); st.Pendulum == "" {
		t.Fatal("klienten ser inte pendeln")
	}
}

func TestLightMirror(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	openRoom(g, p, 4, 3)
	p.known = newGrid[bool](g.size)
	p.known[3][4] = true
	g.give(p, "veke")
	g.give(p, "spegel")
	g.combine(p)
	if p.has("ljusspegel") == nil || len(p.items) != 1 {
		t.Fatalf("väskan har %v", g.itemIDs(p))
	}
	g.mirror(p)
	n := 0
	for d := North; d <= West; d++ {
		if p.known[3+dirDY[d]][4+dirDX[d]] {
			n++
		}
	}
	if n != 2 || p.has("ljusspegel").uses != 2 {
		t.Fatalf("ljusspegeln visade %d rum och har %d kvar", n, p.has("ljusspegel").uses)
	}
	// Mörkret biter inte, och det kostar ingen användning.
	g.trigger(p, g.data.trap("morker"), 0)
	if p.Dark != 0 || p.has("ljusspegel").uses != 2 {
		t.Fatalf("mörker %d, %d kvar", p.Dark, p.has("ljusspegel").uses)
	}
}

func TestScaleMantleCalmsCreatures(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	m := &Monster{ID: 60, Type: g.data.Monster.Typer[0].ID}
	before := g.lynne(m, p)
	g.give(p, "skinn")
	g.give(p, "amulett")
	g.combine(p)
	if p.has("fjallmantel") == nil || p.has("skinn") == nil {
		t.Fatalf("väskan har %v", g.itemIDs(p))
	}
	if got := g.lynne(m, p); got != before/2 {
		t.Fatalf("lynnet %d med manteln, %d utan", got, before)
	}
}

// Pipan drar om två gånger och blir sedan benluren igen.
func TestPipeRedrawsTwiceThenHorn(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	g.give(p, "flojt")
	g.give(p, "benlur")
	g.combine(p)
	if p.has("pipa") == nil || p.has("benlur") == nil {
		t.Fatalf("väskan har %v", g.itemIDs(p))
	}
	wake := -1
	for i, id := range g.deck {
		if c := g.data.card(id); c != nil && c.Vaknar {
			wake = i
		}
	}
	g.flute(p, wake)
	if p.has("pipa") == nil || p.has("pipa").uses != 1 {
		t.Fatal("pipan förbrukades inte en gång")
	}
	g.flute(p, wake)
	if p.has("pipa") != nil || p.has("benlur") == nil || p.has("benlur").ID != "benlur" {
		t.Fatalf("efter två gånger har väskan %v", g.itemIDs(p))
	}
}
