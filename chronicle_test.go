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
	g.Act(b, actWait, 0) // Bo styr varelserna och låter dem vila
	if g.phase != phaseEnd {
		t.Fatalf("fas %s", g.phase)
	}

	all := strings.Join(logs(a), "\n")
	gold := pa.Treasure
	for _, want := range []string{
		"Krönikan över vandringen",
		"Runda 2. Bo dog av pilfällan.",
		"Runda 2. Alva gick ut med " + strconv.Itoa(gold) + " guld.",
		"Alva: kammaren 1 ggr, " + strconv.Itoa(gold) + " guld · fällor 0",
		"Bo: kammaren 0 ggr, 0 guld · fällor 1",
	} {
		if !strings.Contains(all, want) {
			t.Errorf("krönikan saknar %q:\n%s", want, all)
		}
	}
	// Krönikan kommer före epilogen.
	if strings.Index(all, "Krönikan") > strings.Index(all, g.data.Text("spela_igen")) {
		t.Error("krönikan kom efter slutet")
	}

	// En ny vandring börjar med en tom krönika.
	g.Ready(a)
	g.Ready(b)
	if g.phase != phasePlay || len(g.chron) != 0 || pa.stats.takes != 0 {
		t.Fatalf("fas %s, krönika %d, nävar %d", g.phase, len(g.chron), pa.stats.takes)
	}
}
