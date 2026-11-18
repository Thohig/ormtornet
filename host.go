package main

// Spelledarens sida: adressen spelarna ska öppna, vilka som har anslutit och
// vilka inställningar tornet har. Sidan och dess data visas bara på datorn
// som kör servern.

import (
	"encoding/json"
	"io/fs"
	"net"
	"net/http"
	"time"
)

type hostAddrs struct {
	Name string   `json:"name"` // t.ex. http://Stefans-MacBook.local:8080
	IPs  []string `json:"ips"`  // reservadresser med IP-nummer
}

type hostPlayer struct {
	Name   string `json:"name"`
	Slot   int    `json:"slot"`
	Bot    bool   `json:"bot"`
	Online bool   `json:"online"`
	Ready  bool   `json:"ready"`
	Status string `json:"status"`
	AI     bool   `json:"ai"`          // en AI-spelare styr medan människan är borta
	CanAI  bool   `json:"canTakeOver"` // spelledaren kan låta AI:n ta över
}

type hostView struct {
	Addrs      hostAddrs    `json:"addrs"`
	Phase      string       `json:"phase"`
	Round      int          `json:"round"`
	Sun        int          `json:"sun"`
	SunMax     int          `json:"sunMax"`
	TurnSecs   int          `json:"turnSecs"`
	MinPlayers int          `json:"minPlayers"`
	MaxPlayers int          `json:"maxPlayers"`
	Reveal     bool         `json:"reveal"`
	Players    []hostPlayer `json:"players"`
	Settings   Settings     `json:"settings"`
	Bots       int          `json:"bots"` // datorspelare som är med just nu
}

// HostView sammanfattar spelet för spelledaren. Den visar inte skatter,
// positioner eller drag, så spelledaren som också spelar får inget övertag.
func (g *Game) HostView() hostView {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.lastActive = time.Now()
	v := hostView{
		Phase: g.phase, Round: g.round, Sun: g.sun, SunMax: g.cfg.SunMax,
		TurnSecs: int(g.cfg.TurnTime / time.Second), MinPlayers: g.cfg.MinPlayers,
		MaxPlayers: maxPlayers, Reveal: g.cfg.Reveal, Players: []hostPlayer{},
		Settings: settingsFrom(g.cfg), Bots: g.botCount(),
	}
	for _, p := range g.players {
		if p == nil {
			continue
		}
		v.Players = append(v.Players, hostPlayer{
			Name: p.Name, Slot: p.Slot, Bot: p.bot != nil, Online: p.present(),
			Ready: p.Ready || p.bot != nil, Status: p.Status,
			AI: p.takenOver, CanAI: g.canTakeOver(p),
		})
	}
	return v
}

// introInfo är det startsidan behöver innan spelaren har anslutit.
type introInfo struct {
	Avslojande bool     `json:"avslojande"`
	Rykten     []string `json:"rykten"`
}

// Intro väljer rykten till startsidan. Med avslöjandet avstängt, för spelare
// som redan känner till tornets hemlighet, blir ryktena ofarliga.
func (g *Game) Intro() introInfo {
	g.mu.Lock()
	defer g.mu.Unlock()
	key := "rykten"
	if !g.cfg.Reveal {
		key = "rykten_invigda"
	}
	return introInfo{Avslojande: g.cfg.Reveal, Rykten: g.data.Texter[key]}
}

// IdleFor säger hur länge tornet har stått tomt: ingen människa ansluten och
// ingen spelledarsida öppen. Används för att stänga ett bortglömt torn.
func (g *Game) IdleFor() time.Duration {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.humansPresent() > 0 {
		g.lastActive = time.Now()
	}
	return time.Since(g.lastActive)
}

// isLocal är sant när förfrågan kommer från den här datorn, oavsett om den
// gick via localhost, datorns namn eller dess IP-adress.
func isLocal(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	if ip.IsLoopback() {
		return true
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, a := range addrs {
		if n, ok := a.(*net.IPNet); ok && n.IP.Equal(ip) {
			return true
		}
	}
	return false
}

func localOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isLocal(r) {
			http.Error(w, "Spelledarens sida visas bara på datorn som kör tornet.", http.StatusForbidden)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		h(w, r)
	}
}

// Close säger till alla spelare att tornet stänger. Klienterna slutar då
// försöka ansluta på nytt.
func (g *Game) Close() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.timer != nil {
		g.timer.Stop()
	}
	g.logAll(g.data.Text("stanger"), "syn")
	for _, p := range g.players {
		if p != nil {
			g.send(p.client, map[string]string{"t": "stangt"})
		}
	}
}

// hostHooks är det spelledarens sida behöver från resten av programmet.
type hostHooks struct {
	Close func()               // stänger tornet
	Save  func(Settings) error // sparar inställningarna till nästa gång
}

// hostAction kräver förutom att förfrågan kommer från den här datorn ett
// eget huvud, så att en annan webbsida i samma webbläsare inte kan göra
// något här (den skulle stoppas av CORS).
func hostAction(h http.HandlerFunc) http.HandlerFunc {
	return localOnly(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Ormtornet") != "spelledare" {
			http.Error(w, "fel förfrågan", http.StatusBadRequest)
			return
		}
		h(w, r)
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

// addHostRoutes lägger till spelledarens sida och det den kan göra.
func addHostRoutes(mux *http.ServeMux, web fs.FS, game *Game, addrs hostAddrs, hooks hostHooks) {
	mux.HandleFunc("POST /spelledare/stang", hostAction(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		go hooks.Close()
	}))
	mux.HandleFunc("POST /spelledare/starta", hostAction(func(w http.ResponseWriter, r *http.Request) {
		if !game.ForceStart() {
			writeJSON(w, http.StatusConflict, map[string]string{"fel": "Det går inte att starta nu. Minst en människa måste vara i tornet, och ingen vandring får pågå."})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	mux.HandleFunc("POST /spelledare/ta-over", hostAction(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Slot int `json:"slot"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256)).Decode(&req); err != nil || !game.TakeOverSlot(req.Slot) {
			writeJSON(w, http.StatusConflict, map[string]string{"fel": "AI:n kan bara ta över en vandrare vars spelare är borta, mitt i en vandring."})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	mux.HandleFunc("POST /spelledare/installningar", hostAction(func(w http.ResponseWriter, r *http.Request) {
		var s Settings
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&s); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"fel": "Kunde inte läsa inställningarna."})
			return
		}
		got, err := game.ApplySettings(s)
		if err != nil {
			writeJSON(w, http.StatusConflict, map[string]any{"fel": "Inställningarna kan ändras när omgången är slut.", "settings": got})
			return
		}
		if hooks.Save != nil {
			if err := hooks.Save(got); err != nil {
				writeJSON(w, http.StatusOK, map[string]any{"settings": got, "varning": "Gäller nu, men kunde inte sparas till nästa gång: " + err.Error()})
				return
			}
		}
		writeJSON(w, http.StatusOK, map[string]any{"settings": got})
	}))
	mux.HandleFunc("GET /intro", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, http.StatusOK, game.Intro())
	})
	mux.HandleFunc("GET /spelledare", localOnly(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, web, "spelledare.html")
	}))
	mux.HandleFunc("GET /spelledare/oversikt", localOnly(func(w http.ResponseWriter, r *http.Request) {
		if !game.Settings().Spelledare {
			writeJSON(w, http.StatusForbidden, map[string]string{"fel": "Översikten visas bara i spelledarläget."})
			return
		}
		writeJSON(w, http.StatusOK, game.GMView())
	}))
	mux.HandleFunc("GET /spelledare/status", localOnly(func(w http.ResponseWriter, r *http.Request) {
		v := game.HostView()
		v.Addrs = addrs
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(v)
	}))
}
