package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHostViewListsPlayersAndBots(t *testing.T) {
	g := testGame(t)
	g.AddBots(1)
	a := testClient()
	g.Join(a, "Alva")

	v := g.HostView(false)
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

func TestAutostop(t *testing.T) {
	const empty, abandoned = 10 * time.Minute, 12 * time.Hour
	g := testGame(t)
	hourAgo := time.Now().Add(-time.Hour)

	g.lastActive, g.lastReal = hourAgo, hourAgo
	if !g.ShouldClose(empty, abandoned) {
		t.Fatal("ett tomt torn stängde inte")
	}

	// En gömd flik med värdsidan håller inte tornet öppet.
	g.lastActive, g.lastReal = hourAgo, hourAgo
	g.HostView(false)
	if !g.ShouldClose(empty, abandoned) {
		t.Fatal("en gömd flik höll tornet öppet")
	}
	// En synlig gör det.
	g.HostView(true)
	if g.ShouldClose(empty, abandoned) {
		t.Fatal("en synlig spelledarsida räknades inte")
	}

	// En ansluten spelare håller tornet öppet, men inte i all evighet om
	// ingen gör något.
	a := testClient()
	g.Join(a, "Alva")
	g.lastActive = hourAgo
	if g.ShouldClose(empty, abandoned) {
		t.Fatal("tornet stängde fast en spelare är ansluten")
	}
	g.lastReal = time.Now().Add(-13 * time.Hour)
	if !g.ShouldClose(empty, abandoned) {
		t.Fatal("tornet stängde inte fast ingen har gjort något på 13 timmar")
	}
}

func TestAutostopWarnsAMinuteBefore(t *testing.T) {
	g := testGame(t)
	a := testClient()
	g.Join(a, "Alva")
	for len(a.send) > 0 {
		<-a.send
	}
	g.lastReal = time.Now().Add(-12*time.Hour + 30*time.Second)
	if g.ShouldClose(10*time.Minute, 12*time.Hour) {
		t.Fatal("stängde för tidigt")
	}
	warned := false
	for len(a.send) > 0 {
		if strings.Contains(string(<-a.send), "en minut") {
			warned = true
		}
	}
	if !warned || g.closeAt.IsZero() {
		t.Fatal("ingen varning en minut innan")
	}
	g.Touch()
	if !g.closeAt.IsZero() {
		t.Fatal("varningen släpptes inte när någon kom tillbaka")
	}
}

func TestIntroRumors(t *testing.T) {
	g := testGame(t)
	if v := g.Intro(); len(v.Rykten) == 0 {
		t.Fatalf("inga rykten: %+v", v)
	}
}
