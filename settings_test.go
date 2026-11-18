package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func countBots(g *Game) int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.botCount()
}

func TestSettingsAddAndRemoveBots(t *testing.T) {
	g := testGame(t)
	if _, err := g.ApplySettings(Settings{AI: 3, Min: 2, Sol: 30, Turtid: 45, Avslojande: false}); err != nil {
		t.Fatal(err)
	}
	if n := countBots(g); n != 3 {
		t.Fatalf("%d datorspelare, ville ha 3", n)
	}
	if g.cfg.SunMax != 30 || g.sun != 30 || g.cfg.TurnTime.Seconds() != 45 || g.cfg.Reveal {
		t.Errorf("inställningarna slog inte igenom: %+v sol %d", g.cfg, g.sun)
	}
	g.ApplySettings(Settings{AI: 1, Min: 2, Sol: 30, Turtid: 45})
	if n := countBots(g); n != 1 {
		t.Fatalf("%d datorspelare, ville ha 1", n)
	}
}

func TestSettingsClamped(t *testing.T) {
	g := testGame(t)
	got, _ := g.ApplySettings(Settings{AI: 9, Min: 0, Sol: 1, Turtid: 0})
	if got.AI != 3 || got.Min != 1 || got.Sol != 4 || got.Turtid != 3 {
		t.Fatalf("%+v", got)
	}
}

func TestSettingsLockedDuringPlay(t *testing.T) {
	g := testGame(t)
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	g.Ready(a)
	g.Ready(b)
	if g.phase != phasePlay {
		t.Fatal("spelet startade inte")
	}
	if _, err := g.ApplySettings(Settings{AI: 2, Min: 2, Sol: 20, Turtid: 30}); err != errPlaying {
		t.Fatalf("fel: %v", err)
	}
	if countBots(g) != 0 || g.cfg.SunMax != 10 {
		t.Fatal("inställningarna ändrades mitt i spelet")
	}
}

func TestBotStepsAsideForHuman(t *testing.T) {
	g := testGame(t)
	g.ApplySettings(Settings{AI: 3, Min: 2, Sol: 10, Turtid: 30, Avslojande: true})
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	if a.player == nil || countBots(g) != 3 {
		t.Fatalf("Alva kom inte in bredvid tre datorspelare (%d)", countBots(g))
	}
	g.Join(b, "Bo")
	if b.player == nil {
		t.Fatal("Bo kom inte in fast en datorspelare kunde gå")
	}
	if n := countBots(g); n != 2 {
		t.Fatalf("%d datorspelare, ville ha 2", n)
	}
	g.Leave(b)
	if n := countBots(g); n != 3 {
		t.Fatalf("%d datorspelare efter att Bo gick, ville ha 3", n)
	}
}

func TestForceStart(t *testing.T) {
	g := testGame(t)
	if g.ForceStart() {
		t.Fatal("startade utan människor")
	}
	a := testClient()
	g.Join(a, "Alva") // inte redo, och min är 2
	if !g.ForceStart() || g.phase != phasePlay {
		t.Fatal("Starta nu startade inte")
	}
	if g.ForceStart() {
		t.Fatal("startade en gång till mitt i spelet")
	}
}

func TestSettingsSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Ormtornet", "installningar.json")
	s, err := loadSettings(path)
	if err != nil || s != defaultSettings {
		t.Fatalf("utan fil: %+v %v", s, err)
	}
	want := Settings{AI: 2, Min: 3, Sol: 60, Turtid: 15, Avslojande: false, Franvaro: 1, Strid: pvpHalf}
	if err := saveSettings(path, want); err != nil {
		t.Fatal(err)
	}
	if s, err = loadSettings(path); err != nil || s != want {
		t.Fatalf("fick %+v %v", s, err)
	}
}

func TestSettingsEndpointSaves(t *testing.T) {
	g := testGame(t)
	var saved Settings
	mux := http.NewServeMux()
	addHostRoutes(mux, nil, g, hostAddrs{}, hostHooks{Close: func() {}, Save: func(s Settings) error { saved = s; return nil }})

	r := httptest.NewRequest("POST", "/spelledare/installningar", strings.NewReader(`{"ai":2,"min":2,"sol":20,"turtid":30,"avslojande":true}`))
	r.RemoteAddr = "127.0.0.1:5000"
	r.Header.Set("X-Ormtornet", "spelledare")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK || saved.AI != 2 || saved.Sol != 20 {
		t.Fatalf("kod %d, sparat %+v", w.Code, saved)
	}
	if n := countBots(g); n != 2 {
		t.Fatalf("%d datorspelare", n)
	}
}
