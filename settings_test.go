package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func countBots(g *Game) int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.botCount()
}

func TestSettingsAddAndRemoveBots(t *testing.T) {
	g := testGame(t)
	if _, err := g.ApplySettings(Settings{AI: 3, Avslojande: false}); err != nil {
		t.Fatal(err)
	}
	if n := countBots(g); n != 3 {
		t.Fatalf("%d datorspelare, ville ha 3", n)
	}
	if g.cfg.Reveal {
		t.Errorf("inställningarna slog inte igenom: %+v", g.cfg)
	}
	// Solen, turtiden och minsta antalet är inga inställningar: det som
	// flaggorna har gett är kvar.
	if g.cfg.SunMax != 10 || g.cfg.TurnTime != time.Hour || g.cfg.MinPlayers != 2 {
		t.Errorf("solen, turtiden eller minsta antalet ändrades: %+v", g.cfg)
	}
	g.ApplySettings(Settings{AI: 1})
	if n := countBots(g); n != 1 {
		t.Fatalf("%d datorspelare, ville ha 1", n)
	}
}

func TestAllBotsOnlyInGMMode(t *testing.T) {
	g := testGame(t)
	got, _ := g.ApplySettings(Settings{AI: 8})
	if got.AI != 7 {
		t.Fatalf("utan spelledarläge: %d AI-spelare, ville ha högst 7", got.AI)
	}
	got, _ = g.ApplySettings(Settings{AI: 8, Spelledare: true})
	if got.AI != 8 || countBots(g) != 8 {
		t.Fatalf("i spelledarläget: %d inställt, %d vid portarna, ville ha 8", got.AI, countBots(g))
	}
	got, _ = g.ApplySettings(Settings{AI: 8})
	if got.AI != 7 || countBots(g) != 7 {
		t.Fatalf("när spelledarläget slogs av: %d inställt, %d vid portarna, ville ha 7", got.AI, countBots(g))
	}
}

func TestSettingsClamped(t *testing.T) {
	g := testGame(t)
	got, _ := g.ApplySettings(Settings{AI: 9, Franvaro: 20})
	if got.AI != 7 || got.Franvaro != 10 {
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
	if _, err := g.ApplySettings(Settings{AI: 2}); err != errPlaying {
		t.Fatalf("fel: %v", err)
	}
	if countBots(g) != 0 || g.cfg.SunMax != 10 {
		t.Fatal("inställningarna ändrades mitt i spelet")
	}
}

func TestBotStepsAsideForHuman(t *testing.T) {
	g := testGame(t)
	g.ApplySettings(Settings{AI: 7, Avslojande: true})
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	if a.player == nil || countBots(g) != 7 {
		t.Fatalf("Alva kom inte in bredvid sju datorspelare (%d)", countBots(g))
	}
	g.Join(b, "Bo")
	if b.player == nil {
		t.Fatal("Bo kom inte in fast en datorspelare kunde gå")
	}
	if n := countBots(g); n != 6 {
		t.Fatalf("%d datorspelare, ville ha 6", n)
	}
	g.Leave(b)
	if n := countBots(g); n != 7 {
		t.Fatalf("%d datorspelare efter att Bo gick, ville ha 7", n)
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
	want := Settings{AI: 2, Avslojande: false, Franvaro: 1, Sprak: langEn}
	if err := saveSettings(path, want); err != nil {
		t.Fatal(err)
	}
	if s, err = loadSettings(path); err != nil || s != want {
		t.Fatalf("fick %+v %v", s, err)
	}
	// En fil från förr, med inställningar som inte finns längre, går att läsa.
	old := `{"ai": 1, "min": 3, "sol": 20, "turtid": 15, "strid": "halva", "klasser": false, "storlek": 13, "invigning": "svar", "forandras": false, "karta": "platt", "avslojande": true}`
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	if s, err = loadSettings(path); err != nil || s.AI != 1 || !s.Avslojande {
		t.Fatalf("gammal fil: %+v %v", s, err)
	}
}

func TestSettingsEndpointSaves(t *testing.T) {
	g := testGame(t)
	var saved Settings
	mux := http.NewServeMux()
	addHostRoutes(mux, nil, g, hostAddrs{}, hostHooks{Close: func() {}, Save: func(s Settings) error { saved = s; return nil }})

	r := httptest.NewRequest("POST", "/spelledare/installningar", strings.NewReader(`{"ai":2,"avslojande":true}`))
	r.RemoteAddr = "127.0.0.1:5000"
	r.Header.Set("X-Ormtornet", "spelledare")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK || saved.AI != 2 || !saved.Avslojande {
		t.Fatalf("kod %d, sparat %+v", w.Code, saved)
	}
	if n := countBots(g); n != 2 {
		t.Fatalf("%d datorspelare", n)
	}
}

func TestAutoSunAndTurn(t *testing.T) {
	for size, want := range map[int]int{9: 40, 11: 50, 13: 60} {
		if got := autoSun(size); got != want {
			t.Errorf("autoSun(%d) = %d, ville ha %d", size, got, want)
		}
	}
	for humans, want := range map[int]int{0: 30, 1: 30, 4: 30, 5: 35, 6: 35, 7: 40, 8: 40} {
		if got := int(autoTurn(humans).Seconds()); got != want {
			t.Errorf("autoTurn(%d) = %d s, ville ha %d", humans, got, want)
		}
	}

	// Solen följer tornets storlek.
	g := testGame(t)
	g.cfg.SunMax = 0
	g.cfg.TurnTime = 0
	g.cfg.GMMode = true
	g.cfg.MinPlayers = 1
	g.AddBots(6)
	if g.size != 11 || g.sunMax != 50 || g.sun != 50 {
		t.Fatalf("automatiskt: torn %d, sol %d av %d", g.size, g.sun, g.sunMax)
	}
	if v := g.HostView(false); v.SunMax != 50 || v.TurnSecs != 30 {
		t.Fatalf("värdsidan: sol %d, turtid %d", v.SunMax, v.TurnSecs)
	}
}
