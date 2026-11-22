package main

import (
	"strings"
	"testing"
)

// Med benluren hör man en varelse i grannrummet när man kommer in.
func TestHornHearsMonsterNextDoor(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	x, y := g.center()-1, g.center()
	putInRoom(g, p, x, y)
	g.monsters = append(g.monsters, &Monster{ID: 71, Type: g.data.Monster.Typer[0].ID, X: x - 1, Y: y, HP: 2})
	logs(a)
	g.hornListen(p)
	if len(logs(a)) != 0 {
		t.Fatal("hörde något utan benluren")
	}
	g.give(p, "benlur")
	g.hornListen(p)
	got := strings.Join(logs(a), "\n")
	if !strings.Contains(got, g.data.Riktningar[West]) {
		t.Fatalf("hörde inte varelsen åt väster: %q", got)
	}
}

// Med benluren hörs ljud två rum bort med väderstreck; utan den bara
// tornets del.
func TestHornHearsTwoRoomsAway(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	x, y := g.center()-1, g.center()-1
	putInRoom(g, pa, x, y)
	putInRoom(g, pb, x+1, y+1) // två steg bort, snett nedåt höger
	logs(b)
	g.hear(x, y, "ljud_falla", "ALLMÄNT", "info", pa)
	if got := logs(b); !contains(got, "ALLMÄNT") {
		t.Fatalf("utan benluren: %q", got)
	}
	g.give(pb, "benlur")
	g.hear(x, y, "ljud_falla", "ALLMÄNT", "info", pa)
	got := strings.Join(logs(b), "\n")
	if strings.Contains(got, "ALLMÄNT") || !strings.Contains(got, g.data.Horn[0]) {
		t.Fatalf("med benluren: %q", got)
	}
}

// Spårfinnaren hör genom väggarna med benluren.
func TestHornLetsTrackerHear(t *testing.T) {
	g, p, _ := withClass(t, "sparfinnare")
	x, y := g.center()-1, g.center()
	putInRoom(g, p, x, y)
	if _, ok := g.hearsFrom(p, x+1, y); ok {
		t.Fatal("spårfinnaren hörde genom väggen utan benluren")
	}
	g.give(p, "benlur")
	if dir, ok := g.hearsFrom(p, x+1, y); !ok || dir != g.data.Riktningar[East] {
		t.Fatalf("med benluren: %q %v", dir, ok)
	}
}
