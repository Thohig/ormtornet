package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// I lobbyn gäller det nya innehållet direkt.
func TestReloadInLobby(t *testing.T) {
	g := testGame(t)
	g.AddBots(1)
	dir := dataWith(t, "texter.json", func(s string) string {
		return strings.Replace(s, `"ljud_varelse": [`, `"ljud_varelse": ["Nytt ljud åt {riktning}.",`, 1)
	})
	d, _, err := loadData(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !g.ReloadData(d) {
		t.Fatal("innehållet väntar fast tornet är i lobbyn")
	}
	if g.base != d || g.data.Texter["ljud_varelse"][0] != "Nytt ljud åt {riktning}." {
		t.Fatal("tornet använder inte det nya innehållet")
	}
	for _, p := range g.players {
		if p == nil || p.bot == nil {
			continue
		}
		found := false
		for i := range g.data.AI.Personligheter {
			found = found || p.bot.pers == &g.data.AI.Personligheter[i]
		}
		if !found {
			t.Fatal("AI-spelaren pekar på den gamla personligheten")
		}
	}
}

// Under en vandring väntar det nya innehållet till nästa torn.
func TestReloadWaitsForNextTower(t *testing.T) {
	g, _, _ := startTwo(t)
	old := g.base
	d, _, err := loadData("data")
	if err != nil {
		t.Fatal(err)
	}
	if g.ReloadData(d) {
		t.Fatal("innehållet bytte mitt i vandringen")
	}
	if g.base != old || g.nextData != d {
		t.Fatal("det nya innehållet väntar inte")
	}
	g.toLobby()
	if g.base != d || g.nextData != nil {
		t.Fatal("det nya innehållet gäller inte i nästa torn")
	}
}

// Har datafilerna fel säger värdsidan vad, och tornet behåller det gamla.
func TestReloadRoute(t *testing.T) {
	g := testGame(t)
	old := g.base
	dir := dataWith(t, "monster.json", func(s string) string {
		return strings.Replace(s, `"vandra_chans"`, `"vandar_chans"`, 1)
	})
	mux := http.NewServeMux()
	addHostRoutes(mux, nil, g, hostAddrs{}, hostHooks{Close: func() {}, Load: func(string) (*Data, error) {
		d, _, err := loadData(dir)
		return d, err
	}})
	r := httptest.NewRequest("POST", "/spelledare/las-om", nil)
	r.RemoteAddr = "127.0.0.1:5000"
	r.Header.Set("X-Ormtornet", "spelledare")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if w.Code != http.StatusUnprocessableEntity || !strings.Contains(body["fel"], "vandar_chans") {
		t.Fatalf("svar %d: %v", w.Code, body)
	}
	if g.base != old {
		t.Fatal("tornet bytte till trasigt innehåll")
	}
}
