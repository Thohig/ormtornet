package main

import (
	"fmt"
	"strings"
	"testing"
)

// fingerprint sammanfattar en vandring: var alla är, vad de bär och hur
// tornet ser ut.
func fingerprint(g *Game) string {
	var b strings.Builder
	fmt.Fprintf(&b, "runda %d, sol %d\n", g.round, g.sun)
	for _, p := range g.players {
		if p != nil {
			fmt.Fprintf(&b, "%s %s %s %d,%d liv %d skatt %d\n", p.Name, p.classID, p.Status, p.X, p.Y, p.HP, p.Treasure)
		}
	}
	for y := range g.rooms {
		for x := range g.rooms[y] {
			if r := &g.rooms[y][x]; r.Revealed {
				fmt.Fprintf(&b, "%d,%d %s %v %s\n", x, y, r.TypeID, r.Doors, r.Event)
			}
		}
	}
	for _, m := range g.monsters {
		fmt.Fprintf(&b, "%s %d,%d %d\n", m.Type, m.X, m.Y, m.HP)
	}
	return b.String()
}

// Samma frö och samma drag ger samma vandring, rum för rum.
func TestSameSeedSameWandering(t *testing.T) {
	for _, seed := range []uint64{1, 42, 123456789} {
		a, _ := simulate(t, 40, seed)
		b, _ := simulate(t, 40, seed)
		if a.seed != seed || b.seed != seed {
			t.Fatalf("fröet blev %d och %d, inte %d", a.seed, b.seed, seed)
		}
		fa, fb := fingerprint(a), fingerprint(b)
		if fa != fb {
			t.Fatalf("frö %d gav olika vandringar:\n%s\n---\n%s", seed, fa, fb)
		}
	}
	x, _ := simulate(t, 40, 7)
	y, _ := simulate(t, 40, 8)
	if fingerprint(x) == fingerprint(y) {
		t.Fatal("två olika frön gav samma vandring")
	}
}

// Fröet från -fro gäller bara den första vandringen.
func TestSeedOnlyForFirstWandering(t *testing.T) {
	g := testGame(t)
	g.cfg.Seed = 99
	g.seedWandering()
	if g.seed != 99 {
		t.Fatalf("fröet är %d", g.seed)
	}
	g.seedWandering()
	if g.seed == 99 {
		t.Fatal("nästa vandring fick samma frö")
	}
}
