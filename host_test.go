package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHostViewListsPlayersAndBots(t *testing.T) {
	g := testGame(t)
	g.AddBots(1)
	a := testClient()
	g.Join(a, "Alva")

	v := g.HostView()
	if v.Phase != phaseLobby || len(v.Players) != 2 {
		t.Fatalf("fas %s, %d spelare", v.Phase, len(v.Players))
	}
	for _, p := range v.Players {
		if p.Bot && !p.Ready {
			t.Error("datorspelaren ska räknas som redo")
		}
		if p.Name == "Alva" && (p.Bot || p.Ready || !p.Online) {
			t.Errorf("Alva: %+v", p)
		}
	}
}

func TestHostPageOnlyLocal(t *testing.T) {
	mux := http.NewServeMux()
	addHostRoutes(mux, nil, testGame(t), hostAddrs{Name: "http://test.local:8080"}, hostHooks{Close: func() {}})

	for _, tc := range []struct {
		remote string
		want   int
	}{
		{"127.0.0.1:5000", http.StatusOK},
		{"[::1]:5000", http.StatusOK},
		{"203.0.113.7:5000", http.StatusForbidden},
	} {
		r := httptest.NewRequest("GET", "/spelledare/status", nil)
		r.RemoteAddr = tc.remote
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Errorf("%s: fick %d, ville ha %d", tc.remote, w.Code, tc.want)
		}
	}
}

func TestCloseNeedsHeader(t *testing.T) {
	mux := http.NewServeMux()
	closed := make(chan bool, 2)
	addHostRoutes(mux, nil, testGame(t), hostAddrs{}, hostHooks{Close: func() { closed <- true }})

	post := func(header string) int {
		r := httptest.NewRequest("POST", "/spelledare/stang", nil)
		r.RemoteAddr = "127.0.0.1:5000"
		if header != "" {
			r.Header.Set("X-Ormtornet", header)
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w.Code
	}
	if code := post(""); code != http.StatusBadRequest {
		t.Fatalf("utan huvud: %d", code)
	}
	if code := post("spelledare"); code != http.StatusNoContent {
		t.Fatalf("med huvud: %d", code)
	}
	if !<-closed {
		t.Fatal("tornet stängdes inte")
	}
}

func TestCloseTellsPlayers(t *testing.T) {
	g := testGame(t)
	a := testClient()
	g.Join(a, "Alva")
	for len(a.send) > 0 {
		<-a.send
	}
	g.Close()
	got := false
	for len(a.send) > 0 {
		if string(<-a.send) == `{"t":"stangt"}` {
			got = true
		}
	}
	if !got {
		t.Fatal("Alva fick inte veta att tornet stänger")
	}
}
