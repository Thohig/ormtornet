package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// Testerna här kommer från provspelet 2026-11-22 (docs/provspel-2026-11-22.md).

// Den som gick ut tomhänt står före de döda och förlorade i resultatet.
func TestResultsEmptyHandedBeforeDead(t *testing.T) {
	g := testGame(t)
	for _, n := range []string{"Död", "Förlorad", "Tomhänt", "Rik"} {
		g.Join(testClient(), n)
	}
	status := map[string]string{"Död": statusDead, "Förlorad": statusLost, "Tomhänt": statusOut, "Rik": statusOut}
	for _, p := range g.players {
		if p != nil {
			p.Status = status[p.Name]
			if p.Name == "Rik" || p.Name == "Död" {
				p.Treasure = 30
			}
		}
	}
	var order []string
	for _, r := range g.results() {
		order = append(order, r.Name)
	}
	if got := strings.Join(order, " "); got != "Rik Tomhänt Död Förlorad" {
		t.Fatalf("ordningen blev %s", got)
	}
}

// Varje föremål har en obestämd form för texterna, på båda språken, så att
// det inte står "Sigge hittade extra veke".
func TestItemsHaveIndefiniteForm(t *testing.T) {
	for _, lang := range []string{langSv, langEn} {
		d, _, err := loadDataLang("data", lang)
		if err != nil {
			t.Fatal(err)
		}
		for _, it := range d.Foremal.Lista {
			if it.Hittad == "" {
				t.Errorf("%s: %s saknar hittad", lang, it.ID)
			}
		}
	}
}

// Svaret på ett drag mot en vägg är märkt, så att spelarsidan lägger det i
// rundan som står i bild.
func TestWallReplyIsMarked(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	room := &g.rooms[p.Y][p.X]
	wall := Dir(-1)
	for d := North; d <= West; d++ {
		if !room.Doors[d] && !(room.HasExit && room.Exit == d) {
			wall = d
			break
		}
	}
	if wall < 0 {
		t.Skip("rummet har ingen vägg")
	}
	for len(a.send) > 0 {
		<-a.send
	}
	g.Act(a, "NESW"[wall:wall+1], 0)
	found := false
	for len(a.send) > 0 {
		var m logMsg
		if json.Unmarshal(<-a.send, &m) == nil && m.T == "log" {
			found = found || m.Svar
		}
	}
	if !found {
		t.Fatal("svaret om väggen var inte märkt")
	}
}

// En runda fast säger "en runda", inte "1 rundor".
func TestStuckOneRound(t *testing.T) {
	g := testGame(t)
	if txt := g.data.Text("sam_fast_en"); txt == "" || strings.Contains(txt, "rundor") {
		t.Fatalf("sam_fast_en: %q", txt)
	}
}

// Den som laddar om sidan i slutet står kvar i resultatet och kommer
// tillbaka till sin egen vandring, inte som en nykomling.
func TestReloadOnEndScreenKeepsResult(t *testing.T) {
	g, a, b := startTwo(t)
	g.leaveWait = 0
	pa := a.player
	pa.X, pa.Y = 0, 0
	g.endGame()
	g.Leave(a)
	named := func() bool {
		for _, r := range g.results() {
			if r.Name == "Alva" {
				return true
			}
		}
		return false
	}
	if !named() {
		t.Fatal("Alva försvann ur resultatet när sidan laddades om")
	}
	// Bo trycker Enter medan Alva är borta: det nya tornet väntar inte på
	// den som har gått, men Alva står kvar tills tornet byggs.
	a2 := testClient()
	g.Join(a2, "Alva")
	if a2.player != pa || pa.Status == statusNew {
		t.Fatalf("Alva kom tillbaka som en ny vandrare: %+v", a2.player)
	}
	if st := g.stateFor(pa); st.Journey == nil {
		t.Fatal("Alva fick inte tillbaka sin vandring")
	}
	g.Ready(a2)
	g.Ready(b)
	if g.phase == phaseEnd {
		t.Fatal("nästa torn startade inte")
	}
}

// Den som lämnar i slutet och inte kommer tillbaka är inte med i nästa torn.
func TestLeftOnEndScreenDroppedFromNextTower(t *testing.T) {
	g, a, b := startTwo(t)
	g.leaveWait = 0
	g.cfg.MinPlayers = 1
	g.endGame()
	g.Leave(a)
	g.Ready(b)
	for _, p := range g.players {
		if p != nil && p.Name == "Alva" {
			t.Fatalf("Alva är kvar i fas %s", g.phase)
		}
	}
}

// Den som laddar om sidan mitt i vandringen får rummet hen står i igen.
func TestReloadMidGameRetellsRoom(t *testing.T) {
	g, a, _ := startTwo(t)
	pa := a.player
	g.Leave(a)
	a2 := testClient()
	g.Join(a2, "Alva")
	if a2.player != pa {
		t.Fatal("Alva fick inte tillbaka sin vandrare")
	}
	found := false
	for len(a2.send) > 0 {
		var m roomMsg
		if json.Unmarshal(<-a2.send, &m) == nil && m.T == "rum" && m.Text != "" {
			found = true
		}
	}
	if !found {
		t.Fatal("rummet skickades inte igen")
	}
}
