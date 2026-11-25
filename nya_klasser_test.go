package main

import "testing"

// Gravplundraren hör var någon dör och hur långt bort, och får dubbelt så
// mycket av den dödes guld.
func TestGraveRobber(t *testing.T) {
	g, pa, pb := withClass(t, "gravplundrare")
	putInRoom(g, pa, 1, 1)
	putInRoom(g, pb, 4, 1)
	pb.Treasure = 20
	logs(pa.client)
	g.kill(pb)
	want := g.data.Text("klass_gravplundrare_dod", "namn", pb.Name, "antal", "3", "horn", g.hornToward(1, 1, 4, 1))
	if !containsPart(logs(pa.client), want) {
		t.Fatalf("gravplundraren hörde inte döden: ville ha %q", want)
	}
	putInRoom(g, pa, 4, 1)
	before := pa.Treasure
	g.pickUpGold(pa)
	if pa.Treasure != before+40 {
		t.Fatalf("gravplundraren fick %d av den dödes 20 guld", pa.Treasure-before)
	}
	m := &Monster{ID: 80, Type: g.data.Monster.Typer[0].ID}
	if got, other := g.lynne(m, pa), g.lynne(m, pb); got != min(100, other+int(g.classVal(pa, "extra_lynne", 0))) || got == other {
		t.Fatalf("lynnet %d mot gravplundraren, %d mot andra", got, other)
	}
}

// Lyssnaren hör två rum bort, tre med benluren, ser sömnmätaren överallt
// och blir utan ljus dubbelt så länge.
func TestListener(t *testing.T) {
	g, pa, pb := withClass(t, "lyssnare")
	putInRoom(g, pa, 1, 1)
	putInRoom(g, pb, 1, 1)
	if _, ok := g.hearsFrom(pa, 3, 1); !ok {
		t.Fatal("lyssnaren hör inte två rum bort")
	}
	if _, ok := g.hearsFrom(pb, 3, 1); ok {
		t.Fatal("den andra hör två rum bort")
	}
	if _, ok := g.hearsFrom(pa, 4, 1); ok {
		t.Fatal("lyssnaren hör tre rum bort utan benluren")
	}
	g.give(pa, "benlur")
	if _, ok := g.hearsFrom(pa, 4, 1); !ok {
		t.Fatal("lyssnaren hör inte tre rum bort med benluren")
	}
	if st := g.stateFor(pa); st.Sleep < 0 || st.SleepText == "" {
		t.Fatal("lyssnaren ser inte sömnmätaren utanför kammaren")
	}
	if st := g.stateFor(pb); st.Sleep >= 0 {
		t.Fatal("den andra ser sömnmätaren utanför kammaren")
	}
	trap := g.data.trap("morker")
	g.trigger(pa, trap, 0)
	if pa.Dark != 2*trap.Rundor {
		t.Fatalf("mörker %d för lyssnaren, fällan ger %d", pa.Dark, trap.Rundor)
	}
}

// Ormens präst: det första vakna-kortet blir ett sovkort, en näve per
// besök i kammaren, och ingen skada när ormen vaknar.
func TestSerpentsPriest(t *testing.T) {
	g, pa, _ := withClass(t, "prast")
	cx, cy := g.chamber()
	pa.X, pa.Y = cx, cy
	g.reveal(cx, cy)
	var wake, sleep string
	for _, c := range g.data.Ormen.Kort {
		if c.Vaknar && wake == "" {
			wake = c.ID
		}
		if !c.Vaknar && sleep == "" {
			sleep = c.ID
		}
	}
	// Leken har ett vakna-kort och ett sovkort. Prästen drar sovkortet
	// vilket hen än får först; det första vakna-kortet blir ett sovkort.
	for range 50 {
		g.deck = []string{wake, sleep}
		pa.tookVisit = false
		g.take(pa)
		if len(g.deck) != 1 || g.deck[0] != wake {
			t.Fatalf("prästen väckte ormen: leken %v", g.deck)
		}
		if pa.priestCalm {
			break
		}
	}
	if !pa.priestCalm {
		t.Fatal("prästen lugnade aldrig ormen")
	}
	// Bara första gången: sedan väcker vakna-kortet ormen som för andra.
	g.deck = []string{wake}
	pa.tookVisit = false
	g.take(pa)
	if len(g.deck) == 1 {
		t.Fatal("prästen lugnade ormen en andra gång")
	}
	pa.X, pa.Y = cx, cy
	g.deck = []string{sleep, sleep}
	pa.tookVisit = false
	g.take(pa)
	took := pa.Treasure
	g.take(pa)
	if pa.Treasure != took || len(g.deck) != 1 {
		t.Fatal("prästen tog två nävar på samma besök")
	}
	pa.X, pa.Y = cx, cy
	// Ut ur kammaren och tillbaka: en ny näve.
	pa.tookVisit = false
	g.take(pa)
	if len(g.deck) != 0 {
		t.Fatal("prästen fick ingen ny näve efter att ha gått ut")
	}

	pa.X, pa.Y = cx, cy
	hp := pa.HP
	g.wakeSerpent()
	if pa.HP != hp {
		t.Fatalf("ormen skadade prästen: liv %d -> %d", hp, pa.HP)
	}
}

// AI-prästen går ut ur kammaren när hen har tagit sin näve.
func TestPriestBotLeavesChamber(t *testing.T) {
	g, pa, _ := withClass(t, "prast")
	cx, cy := g.chamber()
	pa.X, pa.Y = cx, cy
	g.reveal(cx, cy)
	pa.tookVisit = true
	if _, ok := g.anyDoor(pa); !ok {
		t.Fatal("kammaren har ingen dörr")
	}
}
