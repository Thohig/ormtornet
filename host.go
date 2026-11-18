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
}

// HostView sammanfattar spelet för spelledaren. Den visar inte skatter,
// positioner eller drag, så spelledaren som också spelar får inget övertag.
func (g *Game) HostView() hostView {
	g.mu.Lock()
	defer g.mu.Unlock()
	v := hostView{
		Phase: g.phase, Round: g.round, Sun: g.sun, SunMax: g.cfg.SunMax,
		TurnSecs: int(g.cfg.TurnTime / time.Second), MinPlayers: g.cfg.MinPlayers,
		MaxPlayers: maxPlayers, Reveal: g.cfg.Reveal, Players: []hostPlayer{},
	}
	for _, p := range g.players {
		if p == nil {
			continue
		}
		v.Players = append(v.Players, hostPlayer{
			Name: p.Name, Slot: p.Slot, Bot: p.bot != nil, Online: p.present(),
			Ready: p.Ready || p.bot != nil, Status: p.Status,
		})
	}
	return v
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

func addHostRoutes(mux *http.ServeMux, web fs.FS, game *Game, addrs hostAddrs) {
	mux.HandleFunc("GET /spelledare", localOnly(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, web, "spelledare.html")
	}))
	mux.HandleFunc("GET /spelledare/status", localOnly(func(w http.ResponseWriter, r *http.Request) {
		v := game.HostView()
		v.Addrs = addrs
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(v)
	}))
}
