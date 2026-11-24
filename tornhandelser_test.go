package main

import (
	"strings"
	"testing"
)

// hasEventCard är sant när vandraren fick händelsens spelkort.
func hasEventCard(p *Player, bild string) bool {
	for _, c := range p.roundCards {
		if c.Slag == "handelse" && c.Bild == bild {
			return true
		}
	}
	return false
}

// Varje händelse gör det den ska, och alla inne får ett spelkort.
func TestTowerEvents(t *testing.T) {
	t.Run("andas", func(t *testing.T) {
		g, a, _ := startTwo(t)
		r1, r2 := g.reveal(3, 3), g.reveal(4, 3)
		r1.Kind, r2.Kind = "rum", "rum"
		r1.Doors[East], r2.Doors[West] = false, false
		for y := range g.size {
			for x := range g.size {
				if r := &g.rooms[y][x]; r.Revealed && r.Kind == "rum" && (x != 3 || y != 3) && (x != 4 || y != 3) {
					r.Kind = "ingang" // bara ett par att öppna
				}
			}
		}
		if !g.runTowerEvent("andas") || !g.rooms[3][3].Doors[East] || !g.rooms[3][4].Doors[West] {
			t.Fatal("ingen ny öppning")
		}
		if !hasEventCard(a.player, "andas") {
			t.Fatal("inget spelkort")
		}
	})
	t.Run("guld", func(t *testing.T) {
		g, a, _ := startTwo(t)
		g.runTowerEvent("guld")
		n := 0
		for y := range g.size {
			for x := range g.size {
				n += g.rooms[y][x].loose
			}
		}
		e := g.data.towerEvent("guld")
		if n < e.GuldMin || n > e.GuldMax || !hasEventCard(a.player, "spricka") {
			t.Fatalf("%d guld i tornet", n)
		}
	})
	t.Run("ljus", func(t *testing.T) {
		g, a, _ := startTwo(t)
		p := a.player
		p.Dark = 3
		g.runTowerEvent("ljus")
		if p.Dark != 0 || g.lightRounds == 0 {
			t.Fatalf("mörker %d, ljus %d", p.Dark, g.lightRounds)
		}
		putInRoom(g, p, 4, 3)
		g.trigger(p, g.data.trap("morker"), 0)
		if p.Dark != 0 {
			t.Fatal("mörkret bet trots ljuset")
		}
	})
	t.Run("skinn", func(t *testing.T) {
		g, _, _ := startTwo(t)
		g.runTowerEvent("skinn")
		found := false
		for y := range g.size {
			for x := range g.size {
				for _, id := range g.rooms[y][x].items {
					found = found || id == "skinn"
				}
			}
		}
		if !found {
			t.Fatal("inget skinn i tornet")
		}
	})
	t.Run("ras", func(t *testing.T) {
		g, a, _ := startTwo(t)
		p := a.player
		putInRoom(g, p, 4, 3) // det enda rum någon har gått i
		g.rooms[3][4].Walked = true
		if !g.runTowerEvent("ras") || g.rooms[3][4].blocked == 0 || p.Cellar == 0 {
			t.Fatalf("raset: stängt %d, källaren %d", g.rooms[3][4].blocked, p.Cellar)
		}
	})
	t.Run("doden", func(t *testing.T) {
		g, _, _ := startTwo(t)
		if g.eventPossible("doden") {
			t.Fatal("den döde vandrar fast ingen har dött")
		}
		g.reveal(2, 4).fell = []string{"Bo"}
		if !g.runTowerEvent("doden") || g.monsterAt(2, 4) == nil {
			t.Fatal("ingen varelse där Bo dog")
		}
	})
	t.Run("klockslag", func(t *testing.T) {
		g, a, _ := startTwo(t)
		p := a.player
		cx, cy := g.chamber()
		g.reveal(cx, cy)
		p.X, p.Y = cx, cy
		sun := g.sun
		g.runTowerEvent("klockslag")
		if g.sun != sun-1 || !p.safeDraw {
			t.Fatalf("solen %d -> %d, säker näve %v", sun, g.sun, p.safeDraw)
		}
	})
}

// Klockan stannar, draget och det som klirrar.
func TestTowerEventsMore(t *testing.T) {
	t.Run("stannar", func(t *testing.T) {
		g, a, _ := startTwo(t)
		if g.eventPossible("stannar") {
			t.Fatal("klockan stannar fast solen är full")
		}
		g.sun = g.sunMax - 3
		if !g.runTowerEvent("stannar") || g.sun != g.sunMax-2 || !hasEventCard(a.player, "klocka_stilla") {
			t.Fatalf("solen %d av %d", g.sun, g.sunMax)
		}
	})
	t.Run("draget", func(t *testing.T) {
		g, a, b := withClass(t, "lyktbarare")
		if !g.runTowerEvent("draget") || a.Dark != 0 || b.Dark == 0 {
			t.Fatalf("mörker %d och %d", a.Dark, b.Dark)
		}
		g.lightRounds = 2
		if g.eventPossible("draget") {
			t.Fatal("draget medan ljuset faller in")
		}
	})
	t.Run("klirrar", func(t *testing.T) {
		g, a, b := startTwo(t)
		if g.eventPossible("klirrar") {
			t.Fatal("det klirrar fast ingen bär något")
		}
		putInRoom(g, a.player, 2, 2)
		putInRoom(g, b.player, 4, 4)
		a.player.Treasure, b.player.Treasure = 40, 10
		if !g.runTowerEvent("klirrar") {
			t.Fatal("inget klirr")
		}
		e := g.data.towerEvent("klirrar")
		ca, cb := a.player.roundCards[len(a.player.roundCards)-1], b.player.roundCards[len(b.player.roundCards)-1]
		if ca.Text != e.Har[0] || !strings.Contains(cb.Text, g.hornToward(4, 4, 2, 2)) {
			t.Fatalf("texterna: %q / %q", ca.Text, cb.Text)
		}
	})
}

// Den glömda packningen och brevet under dörren.
func TestTowerEventsFound(t *testing.T) {
	t.Run("packning", func(t *testing.T) {
		g, a, _ := startTwo(t)
		if !g.runTowerEvent("packning") || !hasEventCard(a.player, "packning") {
			t.Fatal("ingen packning")
		}
		n := 0
		for y := range g.size {
			for x := range g.size {
				for _, id := range g.rooms[y][x].items {
					if it := g.data.itemType(id); it != nil && !it.Forbannad && !it.Brev {
						n++
					}
				}
			}
		}
		if n != 1 {
			t.Fatalf("%d föremål i tornet", n)
		}
	})
	t.Run("brev", func(t *testing.T) {
		g, a, b := startTwo(t)
		putInRoom(g, a.player, 2, 2)
		putInRoom(g, b.player, 4, 4)
		if !g.runTowerEvent("brev") {
			t.Fatal("inget brev")
		}
		got := 0
		for _, p := range []*Player{a.player, b.player} {
			for _, c := range p.roundCards {
				if c.Slag == "brev" {
					got++
					if !hasEventCard(p, "brev_dorr") {
						t.Fatal("brevet kom utan händelsens kort")
					}
				}
			}
		}
		if got != 1 {
			t.Fatalf("%d fick ett brev", got)
		}
	})
}

// Tornet sluter sig och ynglet ur väggen.
func TestTowerEventsShut(t *testing.T) {
	t.Run("sluter", func(t *testing.T) {
		g, a, _ := startTwo(t)
		// Fyra rum i en ring: varje dörr har en annan väg runt.
		for y := range g.size {
			for x := range g.size {
				if r := &g.rooms[y][x]; r.Revealed && r.Kind == "rum" {
					r.Kind = "ingang"
				}
			}
		}
		ring := [][2]int{{3, 3}, {4, 3}, {4, 4}, {3, 4}}
		for _, c := range ring {
			r := g.reveal(c[0], c[1])
			r.Kind, r.Doors = "rum", [4]bool{}
		}
		g.rooms[3][3].Doors[East], g.rooms[3][4].Doors[West] = true, true
		g.rooms[3][4].Doors[South], g.rooms[4][4].Doors[North] = true, true
		g.rooms[4][4].Doors[West], g.rooms[4][3].Doors[East] = true, true
		g.rooms[4][3].Doors[North], g.rooms[3][3].Doors[South] = true, true
		putInRoom(g, a.player, 3, 3)
		g.rooms[3][3].Kind = "rum"
		if !g.runTowerEvent("sluter") || len(g.shutDoors) != 1 {
			t.Fatal("ingen dörr slöts")
		}
		sd := g.shutDoors[0]
		if g.rooms[sd.y][sd.x].Doors[sd.d] {
			t.Fatal("dörren står öppen")
		}
		// Den slutna dörren var inte den enda vägen: det går fortfarande runt.
		if !g.otherWay(sd.x, sd.y, sd.d) {
			t.Fatal("tornet stängde den enda vägen")
		}
		for range sd.rounds {
			g.reopenDoors()
		}
		if !g.rooms[sd.y][sd.x].Doors[sd.d] || len(g.shutDoors) != 0 {
			t.Fatal("dörren öppnades inte igen")
		}
	})
	t.Run("yngel", func(t *testing.T) {
		g, a, b := startTwo(t)
		x, y := 4, 4
		if g.isChamber(x, y) {
			x = 3
		}
		putInRoom(g, a.player, x, y)
		putInRoom(g, b.player, x, y)
		if g.eventPossible("yngel") {
			t.Fatal("ynglet kryper fram fast ingen är ensam")
		}
		putInRoom(g, b.player, 1, 0)
		if !g.runTowerEvent("yngel") || g.monsterAt(x, y) == nil || !hasEventCard(a.player, "yngel_vagg") {
			t.Fatal("inget yngel hos den som står längst in")
		}
		if hasEventCard(b.player, "yngel_vagg") {
			t.Fatal("den andra märkte ynglet")
		}
	})
}

// Ett rasat rum går inte att gå in i, och varelser går inte dit heller.
func TestCaveInBlocks(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	putInRoom(g, p, 4, 3)
	g.rooms[3][4].Doors[East] = true
	g.reveal(5, 3).blocked = 2
	g.move(p, East, 0)
	if p.X != 4 {
		t.Fatal("vandraren gick in i det rasade rummet")
	}
}
