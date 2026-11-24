package main

import (
	"strconv"
	"strings"
	"testing"
)

func TestChronicle(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player

	// Runda 1: Alva tar skatt i kammaren, Bo väntar.
	g.reveal(g.center(), g.center())
	pa.X, pa.Y = g.center(), g.center()
	g.deck = []string{"sover"}
	step(g, a, b, actTake, actWait)

	// Runda 2: Bo går i pilfällan med bara ett liv kvar.
	putInRoom(g, pb, 3, 3)
	pb.HP = 1
	g.trigger(pb, g.data.trap("pilar"), 0)
	if pb.Status != statusDead {
		t.Fatalf("Bo är %s", pb.Status)
	}

	// Alva går ut genom sin port, och vandringen är slut.
	pa.X, pa.Y = 0, 0
	g.Act(a, "W", 0)
	if g.phase != phaseEnd {
		t.Fatalf("fas %s", g.phase)
	}

	var rows []string
	for _, l := range g.chronicleLines() {
		rows = append(rows, l.Text)
	}
	all := strings.Join(rows, "\n")
	gold := pa.Treasure
	for _, want := range []string{
		"Krönikan över vandringen",
		"Runda 2. Bo dog av pilar ur små hål i väggen.",
		"Runda 2. Alva gick ut med " + strconv.Itoa(gold) + " guld.",
		"Runda 1. Alva nådde skattkammaren först.",
		"Alva: kammaren en gång, " + strconv.Itoa(gold) + " guld",
		"Bo: fällor 1",
	} {
		if !strings.Contains(all, want) {
			t.Errorf("krönikan saknar %q:\n%s", want, all)
		}
	}
	// Nollor står inte med.
	if strings.Contains(all, " 0") {
		t.Errorf("krönikan har nollor:\n%s", all)
	}
	// Krönikan står på skärmen, inte i loggen, och epilogen följer med.
	if strings.Contains(strings.Join(logs(a), "\n"), "Krönikan över vandringen") {
		t.Error("krönikan står i loggen")
	}
	if st := g.stateFor(pa); st.Journey == nil || len(st.Epilogue) == 0 {
		t.Error("slutet saknar den egna vandringen eller epilogen")
	}

	// En ny vandring börjar med en tom krönika.
	g.Ready(a)
	g.Ready(b)
	if g.phase != phasePlay || len(g.chron) != 0 || pa.stats.takes != 0 {
		t.Fatalf("fas %s, krönika %d, nävar %d", g.phase, len(g.chron), pa.stats.takes)
	}
}

// Den som dör ser sin egen vandring, men inget om de andra utom hur många
// som är kvar. Den som lever ser den inte förrän vandringen är slut.
func TestJourneyForTheDead(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	putInRoom(g, pa, 3, 3)
	putInRoom(g, pb, 5, 5)
	step(g, a, b, actSearch, actWait)
	if st := g.stateFor(pa); st.Journey != nil {
		t.Fatal("den levande ser sin vandring mitt i vandringen")
	}
	pb.HP = 1
	g.trigger(pb, g.data.trap("pilar"), 0)
	st := g.stateFor(pb)
	if st.Journey == nil {
		t.Fatal("den döda ser ingen vandring")
	}
	if !strings.Contains(st.Journey.Death, "Bo dog av pilar ur små hål i väggen") {
		t.Errorf("döden: %q", st.Journey.Death)
	}
	if st.Journey.In != 1 || st.Journey.Out != 0 {
		t.Errorf("kvar %d, ute %d", st.Journey.In, st.Journey.Out)
	}
	if !strings.Contains(st.Journey.Stats, "Bo:") {
		t.Errorf("siffrorna: %q", st.Journey.Stats)
	}
	if len(pa.journal) == 0 {
		t.Error("Alvas runda med sökandet skrevs inte ner")
	}
}

// Den som dör med guld får guldet som gick förlorat på samma rad.
func TestChronicleDeathLosesGold(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.Treasure = 45
	p.hurt("falla", "pilar ur små hål i väggen")
	g.kill(p)
	last := g.chron[len(g.chron)-1].text
	if want := g.data.Text("kronika_dod_guld", "guld", "45"); !strings.HasSuffix(last, want) {
		t.Fatalf("väntade %q sist i %q", want, last)
	}
}
