package main

import (
	"strings"
	"testing"
)

func viewOf(st stateMsg, p *Player) playerView {
	for _, v := range st.Players {
		if v.ID == p.ID {
			return v
		}
	}
	return playerView{}
}

func TestSightMeetHidesOthers(t *testing.T) {
	g, a, b := startTwo(t)
	g.cfg.Sight = sightMeet
	pa, pb := a.player, b.player
	putInRoom(g, pa, 2, 3)
	putInRoom(g, pb, 5, 3)

	st := g.stateFor(pa)
	if v := viewOf(st, pb); !v.Hidden || v.X != -1 {
		t.Fatalf("Bo syns i ett annat rum: %+v", v)
	}
	if v := viewOf(st, pa); v.Hidden || v.X != 2 {
		t.Fatalf("Alva ser inte sig själv: %+v", v)
	}

	// I samma rum syns de.
	putInRoom(g, pb, 2, 3)
	if v := viewOf(g.stateFor(pa), pb); v.Hidden || v.X != 2 || v.Y != 3 {
		t.Fatalf("Bo syns inte i samma rum: %+v", v)
	}

	// Med "alla" syns alla som förut.
	putInRoom(g, pb, 5, 3)
	g.cfg.Sight = sightAll
	if v := viewOf(g.stateFor(pa), pb); v.Hidden || v.X != 5 {
		t.Fatalf("Bo syns inte med synligheten alla: %+v", v)
	}
}

func TestSightMeetControllerAndEnd(t *testing.T) {
	g, a, b := startTwo(t)
	g.cfg.Sight = sightMeet
	pa, pb := a.player, b.player
	putInRoom(g, pa, 2, 3)
	putInRoom(g, pb, 5, 3)

	// Den döda spelaren som styr varelserna ser alla.
	pa.Status = statusDead
	g.controller = pa.ID
	if v := viewOf(g.stateFor(pa), pb); v.Hidden {
		t.Fatal("den som styr varelserna ser inte Bo")
	}

	// När vandringen är slut syns alla.
	g.controller = 0
	g.phase = phaseEnd
	if v := viewOf(g.stateFor(pa), pb); v.Hidden {
		t.Fatal("Bo syns inte i slutet")
	}
}

func TestSightMeetVagueTexts(t *testing.T) {
	g, a, b := startTwo(t)
	g.cfg.Sight = sightMeet
	putInRoom(g, a.player, 4, 3)
	putInRoom(g, b.player, 1, 1)
	g.rooms[3][5] = Room{}
	logs(b)
	g.move(a.player, East, 0)
	for _, line := range logs(b) {
		if strings.Contains(line, "Alva") {
			t.Fatalf("namnet syns i texten om ett nytt rum: %q", line)
		}
	}
}

func TestSightSetting(t *testing.T) {
	if s := (Settings{}).normalize(); s.Synlighet != sightOwn {
		t.Fatalf("standard %q", s.Synlighet)
	}
	if c := (Settings{Synlighet: sightAll}).config(); c.Sight != sightAll {
		t.Fatalf("config %q", c.Sight)
	}
	if s := settingsFrom(Config{Sight: sightAll}); s.Synlighet != sightAll {
		t.Fatalf("settingsFrom %q", s.Synlighet)
	}
}

func hasRoom(st stateMsg, x, y int) bool {
	for _, r := range st.Rooms {
		if r.X == x && r.Y == y {
			return true
		}
	}
	return false
}

func TestSightOwnMap(t *testing.T) {
	g, a, b := startTwo(t)
	g.cfg.Sight = sightOwn
	pa, pb := a.player, b.player
	putInRoom(g, pa, 4, 3)
	putInRoom(g, pb, 1, 1)
	g.rooms[3][5] = Room{}
	g.noteRoom(pa)
	g.noteRoom(pb)
	logs(b)

	// Alva går in i ett nytt rum. Bo får inte veta något och ser inte rummet.
	g.move(pa, East, 0)
	if !hasRoom(g.stateFor(pa), 5, 3) {
		t.Fatal("Alva ser inte rummet hon gick in i")
	}
	if hasRoom(g.stateFor(pb), 5, 3) || hasRoom(g.stateFor(pb), 4, 3) {
		t.Fatal("Bo ser rum han aldrig har varit i")
	}
	for _, line := range logs(b) {
		if strings.Contains(line, "nytt rum") || strings.Contains(line, "ny väg") {
			t.Fatalf("Bo fick höra om det nya rummet: %q", line)
		}
	}
	// Bo ser sin egen startport.
	cn := corners[pb.Slot]
	if !hasRoom(g.stateFor(pb), cn.x, cn.y) {
		t.Fatal("Bo ser inte sin egen port")
	}

	// En varelse i ett rum Bo inte känner till syns inte för honom.
	g.monsters = append(g.monsters, &Monster{ID: 99, Type: g.data.Monster.Typer[0].ID, X: 5, Y: 3, HP: 1})
	for _, m := range g.stateFor(pb).Monsters {
		if m.ID == 99 {
			t.Fatal("Bo ser en varelse i ett okänt rum")
		}
	}

	// När Bo går dit blir rummet hans.
	putInRoom(g, pb, 5, 3)
	if !hasRoom(g.stateFor(pb), 5, 3) {
		t.Fatal("Bo ser inte rummet han står i")
	}

	// I slutet syns hela tornet.
	g.phase = phaseEnd
	if !hasRoom(g.stateFor(pb), 4, 3) {
		t.Fatal("hela tornet syns inte i slutet")
	}
}

func TestSightOwnController(t *testing.T) {
	g, a, b := startTwo(t)
	g.cfg.Sight = sightOwn
	pa, pb := a.player, b.player
	putInRoom(g, pb, 6, 6)
	pa.Status = statusDead
	g.controller = pa.ID
	if !hasRoom(g.stateFor(pa), 6, 6) {
		t.Fatal("den som styr varelserna ser inte hela tornet")
	}
}
