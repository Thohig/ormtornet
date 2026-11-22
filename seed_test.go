package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
	g.nextSeed = 99
	g.seedWandering()
	if g.seed != 99 {
		t.Fatalf("fröet är %d", g.seed)
	}
	g.seedWandering()
	if g.seed == 99 {
		t.Fatal("nästa vandring fick samma frö")
	}
}

// Tornet vet alltid nästa frö. Vandringen tar det och slumpar ett nytt
// till vandringen efter.
func TestNextSeedIsKnownAndUsed(t *testing.T) {
	g := testGame(t)
	if g.nextSeed == 0 || g.nextSeed > maxSeed {
		t.Fatalf("nästa frö är %d", g.nextSeed)
	}
	if !g.SetNextSeed(4242) {
		t.Fatal("4242 gick inte att sätta")
	}
	if g.seed != 4242 {
		t.Fatalf("lobbyn byggdes inte om med fröet: %d", g.seed)
	}
	g.seedWandering()
	if g.seed != 4242 {
		t.Fatalf("vandringen fick %d", g.seed)
	}
	if g.nextSeed == 0 || g.nextSeed == 4242 {
		t.Fatalf("nästa frö efter vandringen är %d", g.nextSeed)
	}
}

func TestSetNextSeedLimits(t *testing.T) {
	g := testGame(t)
	if g.SetNextSeed(maxSeed + 1) {
		t.Fatal("ett för stort frö gick att sätta")
	}
	before := g.nextSeed
	if !g.SetNextSeed(0) || g.nextSeed == 0 {
		t.Fatal("slumpa gav inget frö")
	}
	_ = before
}

func TestSeedRoute(t *testing.T) {
	g := testGame(t)
	mux := http.NewServeMux()
	addHostRoutes(mux, nil, g, hostAddrs{}, hostHooks{})
	post := func(body string) int {
		r := httptest.NewRequest("POST", "/spelledare/fro", strings.NewReader(body))
		r.RemoteAddr = "127.0.0.1:5000"
		r.Header.Set("X-Ormtornet", "spelledare")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w.Code
	}
	if code := post(`{"fro":777}`); code != http.StatusNoContent || g.nextSeed != 777 {
		t.Fatalf("fro 777: %d, nästa %d", code, g.nextSeed)
	}
	if code := post(`{"slumpa":true}`); code != http.StatusNoContent || g.nextSeed == 777 {
		t.Fatalf("slumpa: %d, nästa %d", code, g.nextSeed)
	}
	for _, bad := range []string{`{"fro":0}`, `{"fro":-3}`, `{"fro":1000000000}`, `{"fro":"abc"}`, `nej`} {
		if code := post(bad); code != http.StatusBadRequest {
			t.Errorf("%s gav %d", bad, code)
		}
	}
}

// Samma frö efter en vandring ger samma portar i lobbyn som förut.
func TestSameSeedSameLobby(t *testing.T) {
	a := testGame(t)
	a.SetNextSeed(31337)
	b := testGame(t)
	b.SetNextSeed(5)
	b.SetNextSeed(31337)
	if fmt.Sprint(a.corners()) != fmt.Sprint(b.corners()) {
		t.Fatal("samma frö gav olika portar")
	}
}
