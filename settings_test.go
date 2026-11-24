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
	if _, err := g.ApplySettings(Settings{AI: 3, Franvaro: 4}); err != nil {
		t.Fatal(err)
	}
	if n := countBots(g); n != 3 {
		t.Fatalf("%d datorspelare, ville ha 3", n)
	}
	if g.cfg.AbsentRounds != 4 {
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

// En plats finns alltid kvar åt en människa: högst sju AI-spelare.
func TestAtMostSevenBots(t *testing.T) {
	g := testGame(t)
	got, _ := g.ApplySettings(Settings{AI: 8})
	if got.AI != 7 || countBots(g) != 7 {
		t.Fatalf("%d inställt, %d vid portarna, ville ha 7", got.AI, countBots(g))
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
	g.ApplySettings(Settings{AI: 7})
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
	want := Settings{AI: 2, Franvaro: 1, Sprak: langEn}
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
	if s, err = loadSettings(path); err != nil || s.AI != 1 {
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
	if w.Code != http.StatusOK || saved.AI != 2 {
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
	g.cfg.MinPlayers = 1
	g.AddBots(6)
	if g.size != 11 || g.sunMax != 50 || g.sun != 50 {
		t.Fatalf("automatiskt: torn %d, sol %d av %d", g.size, g.sun, g.sunMax)
	}
	if v := g.HostView(false); v.SunMax != 50 || v.TurnSecs != 30 {
		t.Fatalf("värdsidan: sol %d, turtid %d", v.SunMax, v.TurnSecs)
	}
}

// endedGame ger ett torn där en vandring med Alva och två AI-spelare är
// slut, så att slutskärmen visas.
func endedGame(t *testing.T) (*Game, *Client) {
	t.Helper()
	g := testGame(t)
	a := testClient()
	g.Join(a, "Alva")
	if _, err := g.ApplySettings(Settings{AI: 2, Franvaro: 2}); err != nil {
		t.Fatal(err)
	}
	g.ForceStart()
	if g.phase != phasePlay {
		t.Fatalf("fas %s, ville ha vandringen", g.phase)
	}
	g.endGame()
	if g.phase != phaseEnd {
		t.Fatalf("fas %s, ville ha slut", g.phase)
	}
	return g, a
}

// Ändras antalet AI-spelare i slutet stängs slutskärmen och ett nytt torn
// i rätt storlek väntar i lobbyn; placeringen tappar aldrig någon.
func TestChangedAIClosesEndScreen(t *testing.T) {
	g, a := endedGame(t)
	if _, err := g.ApplySettings(Settings{AI: 5, Franvaro: 2}); err != nil {
		t.Fatal(err)
	}
	if g.phase != phaseLobby {
		t.Fatalf("fas %s, ville ha lobbyn", g.phase)
	}
	if n := countBots(g); n != 5 {
		t.Fatalf("%d AI-spelare, ville ha 5", n)
	}
	if g.size != g.towerSize() {
		t.Errorf("tornet är %d stort, ville ha %d", g.size, g.towerSize())
	}
	if !strings.Contains(strings.Join(drainAll(a), "\n"), g.data.Texter["torn_andrat"][0]) {
		t.Error("ingen rad om att tornet ändrades")
	}
	if a.player.Ready {
		t.Error("Alva är redo i det nya tornet utan att ha sagt det")
	}
}

// Frånvarande vandrare ändrar inte tornet, så slutskärmen står kvar.
func TestAbsentSettingKeepsEndScreen(t *testing.T) {
	g, _ := endedGame(t)
	g.ApplySettings(Settings{AI: 2, Franvaro: 5})
	if g.phase != phaseEnd || countBots(g) != 2 {
		t.Fatalf("fas %s med %d AI-spelare", g.phase, countBots(g))
	}
}

// Ett nytt frö i slutet ger ett nytt torn med det fröet i lobbyn, samma
// torn som en lobby med samma frö ger.
func TestNewSeedClosesEndScreen(t *testing.T) {
	g, _ := endedGame(t)
	if !g.SetNextSeed(4242) {
		t.Fatal("fröet togs inte emot")
	}
	if g.phase != phaseLobby || g.nextSeed != 4242 || countBots(g) != 2 {
		t.Fatalf("fas %s frö %d med %d AI-spelare", g.phase, g.nextSeed, countBots(g))
	}
	// Samma frö en gång till ändrar inget.
	ports := g.corners()
	g.SetNextSeed(4242)
	if g.phase != phaseLobby || len(g.corners()) != len(ports) {
		t.Fatal("samma frö igen byggde om tornet")
	}
}
