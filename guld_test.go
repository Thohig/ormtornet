package main

import (
	"strings"
	"testing"
)

// Guldet den döda bar ligger kvar, och den som står kvar i rummet tar det.
func TestDeadGoldIsLeftAndTaken(t *testing.T) {
	g := testGame(t)
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	g.Ready(a)
	g.Ready(b)
	pa, pb := a.player, b.player
	x, y := g.center()-1, g.center()
	g.reveal(x, y)
	pa.X, pa.Y, pb.X, pb.Y = x, y, x, y
	pa.Treasure, pb.Treasure = 40, 5

	pa.HP = 1
	g.damage(pa, 1, true)
	if pa.Status != statusDead || pa.Treasure != 0 {
		t.Fatalf("Alva är %s med %d guld", pa.Status, pa.Treasure)
	}
	if r := g.rooms[y][x]; len(r.gold) != 1 || r.gold[0].gold != 40 {
		t.Fatalf("i rummet ligger %+v", r.gold)
	}

	g.Act(b, actWait, 0)
	if pb.Treasure != 45 {
		t.Fatalf("Bo bär %d guld", pb.Treasure)
	}
	if len(g.rooms[y][x].gold) != 0 {
		t.Fatal("guldet ligger kvar efter att Bo tog det")
	}
	found := false
	for _, e := range g.chron {
		if e.key == "kronika_guld_kvar" && strings.Contains(e.text, "Bo") && strings.Contains(e.text, "Alva") {
			found = true
		}
	}
	if !found {
		t.Fatal("krönikan säger inte att Bo tog Alvas guld")
	}

	pb.Status = statusOut
	g.goldEpilogue()
	if len(pa.epilog) != 1 || !strings.Contains(pa.epilog[0].Text, "Bo") {
		t.Fatalf("Alvas epilog: %+v", pa.epilog)
	}
}

// Den som går in i rummet tar guldet.
func TestDeadGoldTakenOnEntering(t *testing.T) {
	g := testGame(t)
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	g.Ready(a)
	g.Ready(b)
	pa, pb := a.player, b.player
	g.rooms[3][4].gold = []leftGold{{by: pa.ID, name: "Alva", gold: 12}}
	pb.X, pb.Y = 4, 3
	g.pickUp(pb)
	if pb.Treasure != 12 || len(g.rooms[3][4].gold) != 0 {
		t.Fatalf("Bo bär %d guld", pb.Treasure)
	}
	if len(pa.goldTakenBy) != 1 || pa.goldTakenBy[0] != pb.ID {
		t.Fatalf("Alvas guld togs av %v", pa.goldTakenBy)
	}
}
