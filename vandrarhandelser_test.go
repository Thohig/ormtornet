package main

import (
	"strings"
	"testing"
)

// Händelserna mellan vandrarna gör det de ska (vandrarhandelser.go).
func TestBetweenEvents(t *testing.T) {
	t.Run("tillbaka", func(t *testing.T) {
		g, a, b := startTwo(t)
		pa, pb := a.player, b.player
		// Bo går ut med 40 guld; Alva, som inte bär något, får välja.
		pb.Treasure = 40
		cx, cy := g.corners()[pb.Corner].x, g.corners()[pb.Corner].y
		pb.X, pb.Y = cx, cy
		pb.Status = statusOut
		g.sun = g.sunMax
		if !g.runTowerEvent("tillbaka") {
			t.Fatal("händelsen hände inte")
		}
		if !g.canPull(pa) || g.stateFor(pa).Dra != "Bo" {
			t.Fatal("Alva kan inte dra tillbaka Bo")
		}
		g.Act(a, actPull, 0)
		if pb.Status != statusIn || pb.Treasure != 20 || pb.banked != 20 || !pb.pulled {
			t.Fatalf("Bo: %s, bär %d, säkert %d", pb.Status, pb.Treasure, pb.banked)
		}
		if g.gazeOn != pa.ID {
			t.Fatal("Alva fick inte ormens blick")
		}
		if g.pullable(pb) {
			t.Fatal("Bo kan dras tillbaka två gånger")
		}
		// Ute igen får Bo tillbaka det som var säkert.
		pb.Treasure = 25
		g.move(pb, g.rooms[cy][cx].Exit, 0)
		if pb.Status != statusOut || pb.Treasure != 45 {
			t.Fatalf("ute igen med %d", pb.Treasure)
		}
	})
	t.Run("tillbaka och dör", func(t *testing.T) {
		g, a, b := startTwo(t)
		pb := b.player
		pb.Status, pb.Treasure, pb.banked, pb.pulled = statusDead, 5, 30, true
		a.player.Status, a.player.Treasure = statusOut, 10
		var bo resultView
		for _, r := range g.results() {
			if r.Name == "Bo" {
				bo = r
			}
		}
		if bo.Place != 1 || bo.Treasure != 30 {
			t.Fatalf("Bo plats %d med %d", bo.Place, bo.Treasure)
		}
	})
	t.Run("blick", func(t *testing.T) {
		g, a, b := startTwo(t)
		pa, pb := a.player, b.player
		putInRoom(g, pa, 3, 3)
		putInRoom(g, pb, 5, 5)
		g.monsters = nil
		g.giveGaze(pa, 3)
		pa.Treasure = 30
		if g.stateFor(pa).Blick != 3 || g.stateFor(pb).Blick != 0 {
			t.Fatal("blicken syns inte bland tecknen")
		}
		// I samma rum går blicken vidare.
		putInRoom(g, pb, 3, 3)
		g.betweenTick()
		if g.gazeOn != pb.ID {
			t.Fatal("blicken gick inte vidare")
		}
		// När den släpper tar den guld.
		putInRoom(g, pb, 5, 5)
		pb.Treasure = 30
		for range 5 {
			g.betweenTick()
		}
		e := g.data.towerEvent("blick")
		if g.gazeOn != 0 || pb.Treasure != 30-e.Guld {
			t.Fatalf("blicken på %d, Bo bär %d", g.gazeOn, pb.Treasure)
		}
	})
	t.Run("skottpeng", func(t *testing.T) {
		g, a, b := startTwo(t)
		pa, pb := a.player, b.player
		pa.Treasure, pb.Treasure = 60, 0
		before := g.loot(pb, pa)
		if !g.runTowerEvent("skottpeng") || g.bountyOn != pa.ID {
			t.Fatal("inget pris på den rike")
		}
		if got := g.loot(pb, pa); got != 2*before {
			t.Fatalf("bytet %d, utan pris %d", got, before)
		}
		if g.stateFor(pa).Skottpeng == 0 {
			t.Fatal("priset syns inte bland tecknen")
		}
	})
	t.Run("skata", func(t *testing.T) {
		g, a, b := startTwo(t)
		pa, pb := a.player, b.player
		putInRoom(g, pa, 3, 3)
		putInRoom(g, pb, 5, 5)
		pb.items = nil
		pa.items = nil
		g.give(pa, "rep")
		if !g.runTowerEvent("skata") {
			t.Fatal("händelsen hände inte")
		}
		if pa.has("rep") != nil || pb.has("rep") == nil {
			t.Fatal("repet bytte inte väska")
		}
		if !hasEventCard(pa, "skata") || !hasEventCard(pb, "skata") {
			t.Fatal("inget spelkort")
		}
	})
	t.Run("vagg", func(t *testing.T) {
		g, a, b := startTwo(t)
		pa, pb := a.player, b.player
		ra := putInRoom(g, pa, 3, 3)
		rb := putInRoom(g, pb, 4, 3)
		ra.Doors[East], rb.Doors[West] = false, false
		if !g.runTowerEvent("vagg") || !ra.Doors[East] || !rb.Doors[West] {
			t.Fatal("väggen föll inte")
		}
		if !pa.known[3][4] || !pb.known[3][3] {
			t.Fatal("rummet bredvid syns inte")
		}
	})
	t.Run("kista", func(t *testing.T) {
		g, a, b := startTwo(t)
		pa, pb := a.player, b.player
		r := putInRoom(g, pa, 3, 3)
		g.monsters = nil
		r.chest = 50
		r.Searched = true
		// Ensam går kistan inte upp.
		putInRoom(g, pb, 5, 5)
		step(g, a, b, actSearch, actWait)
		if r.chest != 50 {
			t.Fatal("kistan gick upp för en")
		}
		putInRoom(g, pa, 3, 3)
		putInRoom(g, pb, 3, 3)
		g.monsters = nil
		before := pa.Treasure
		step(g, a, b, actSearch, actWait)
		if r.chest != 0 || pa.Treasure != before+50 {
			t.Fatalf("kistan %d, Alva bär %d", r.chest, pa.Treasure)
		}
	})
	t.Run("glod", func(t *testing.T) {
		g, a, b := startTwo(t)
		pa, pb := a.player, b.player
		putInRoom(g, pb, 5, 5)
		putInRoom(g, pa, 3, 3)
		g.monsters = nil
		if !g.runTowerEvent("glod") {
			t.Fatal("händelsen hände inte")
		}
		g.round++
		nr := g.reveal(5, 4)
		nr.Kind = "rum"
		g.rooms[5][5].Doors[North], nr.Doors[South] = true, true
		g.move(pb, North, 0)
		if !g.knowsRoom(pa, 5, 4) || !g.knowsRoom(pa, 5, 5) {
			t.Fatal("Alva ser inte de glödande spåren")
		}
		g.round = g.glowUntil + 1
		if g.knowsRoom(pa, 5, 4) {
			t.Fatal("spåren glöder fortfarande")
		}
	})
	t.Run("vrider", func(t *testing.T) {
		g, a, b := startTwo(t)
		pa, pb := a.player, b.player
		putInRoom(g, pa, 3, 3)
		putInRoom(g, pb, 5, 5)
		g.monsters = nil
		if !g.runTowerEvent("vrider") {
			t.Fatal("händelsen hände inte")
		}
		if pa.X != 5 || pa.Y != 5 || pb.X != 3 || pb.Y != 3 {
			t.Fatal("de bytte inte plats")
		}
		if !pa.known[5][5] {
			t.Fatal("det nya rummet syns inte på kartan")
		}
	})
}

// Alla nya händelser har bok, rubrik och bild, och en text för den som
// drabbas eller för alla.
func TestBetweenEventsData(t *testing.T) {
	g := testGame(t)
	for _, id := range []string{"tillbaka", "blick", "skottpeng", "skata", "vagg", "kista", "glod", "vrider"} {
		e := g.data.towerEvent(id)
		if e == nil {
			t.Fatalf("%s saknas", id)
		}
		if e.Bok == "" || e.Rubrik == "" || len(e.Text)+len(e.Har) == 0 {
			t.Fatalf("%s saknar texter", id)
		}
		if strings.Contains(e.Bok, "{") {
			t.Fatalf("%s: boken har en mall", id)
		}
	}
}
