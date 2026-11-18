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
	addHostRoutes(mux, nil, testGame(t), hostAddrs{Name: "http://test.local:8080"})

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
