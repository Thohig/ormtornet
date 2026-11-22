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

// Skriket: den som står bredvid hör det åt ett håll, de andra i vilken del
// av tornet, och mynten hörs när den döda bar guld.
func TestDeathScreamSaysWhere(t *testing.T) {
	g, cs := nameGame(t, 3)
	if g.phase != phasePlay {
		t.Fatalf("fas %s", g.phase)
	}
	a, b, c := cs[0].player, cs[1].player, cs[2].player
	x, y := g.center(), g.center()-1
	putInRoom(g, a, x, y)
	putInRoom(g, b, x+1, y)
	putInRoom(g, c, 0, g.size-1)
	a.Treasure = 30
	for _, cl := range cs {
		logs(cl)
	}
	a.HP = 1
	g.damage(a, 1, true)

	next := strings.Join(logs(cs[1]), "\n")
	far := strings.Join(logs(cs[2]), "\n")
	region := g.data.region(x, y, g.size)
	if !strings.Contains(next, "åt "+g.data.Riktningar[West]) || !strings.Contains(next, "mynt") {
		t.Fatalf("grannen hörde: %s", next)
	}
	if !strings.Contains(far, region) || !strings.Contains(far, "mynt") {
		t.Fatalf("den långt bort hörde: %s", far)
	}
	if len(b.deaths) != 1 || b.deaths[0].x != x || !b.deaths[0].coins {
		t.Fatalf("grannen minns %+v", b.deaths)
	}
	if len(c.deaths) != 1 || c.deaths[0].x != -1 || c.deaths[0].region != region {
		t.Fatalf("den långt bort minns %+v", c.deaths)
	}

	// Utan guld hörs inga mynt.
	b.HP = 1
	g.damage(b, 1, true)
	if far := strings.Join(logs(cs[2]), "\n"); strings.Contains(far, "mynt") {
		t.Fatalf("mynt utan guld: %s", far)
	}
}
