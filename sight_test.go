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
	g.rooms[3][4].Doors[East] = true // det nya rummet får en dörr tillbaka
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
	cn := g.corners()[pb.Slot]
	if !hasRoom(g.stateFor(pb), cn.x, cn.y) {
		t.Fatal("Bo ser inte sin egen port")
	}
	// Men inte Alvas port, där han aldrig har varit.
	if ca := g.corners()[pa.Slot]; hasRoom(g.stateFor(pb), ca.x, ca.y) {
		t.Fatal("Bo ser en annans port")
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

func TestTrail(t *testing.T) {
	g, a, b := startTwo(t)
	g.cfg.Sight = sightOwn
	g.data.Regler.SparRundor = 3
	pa, pb := a.player, b.player
	putInRoom(g, pa, 4, 4)
	putInRoom(g, pb, 2, 2)
	putInRoom(g, pa, 3, 4) // rummet åt väster finns
	putInRoom(g, pa, 4, 4)
	g.rooms[4][4].Doors[West] = true
	g.rooms[4][3].Doors[East] = true

	// Alva går västerut. Bo kommer in efter henne och ser spåren.
	g.move(pa, West, 0)
	pb.X, pb.Y = 5, 4
	g.rooms[4][5] = g.rooms[4][4]
	g.rooms[4][5].Doors[West] = true
	g.rooms[4][4].Doors[East] = true
	logs(b)
	g.move(pb, West, 0)
	got := logs(b)
	if !containsPart(got, "väster") {
		t.Fatalf("Bo såg inga spår: %q", got)
	}

	// Gamla spår syns inte.
	g.round += 10
	pb.X, pb.Y = 5, 4
	logs(b)
	g.move(pb, West, 0)
	if containsPart(logs(b), "spår") {
		t.Fatal("Bo såg gamla spår")
	}

	// Med synligheten alla finns inga spår.
	g.cfg.Sight = sightAll
	g.round -= 10
	pb.X, pb.Y = 5, 4
	logs(b)
	g.move(pb, West, 0)
	if containsPart(logs(b), "spår") {
		t.Fatal("spår fast alla syns")
	}
}

func TestNoiseThroughWall(t *testing.T) {
	g, a, b := startTwo(t)
	g.cfg.Sight = sightOwn
	pa, pb := a.player, b.player
	putInRoom(g, pa, 4, 4)
	putInRoom(g, pb, 4, 3) // norr om Alva
	pa.Potions = 1
	logs(b)
	g.drink(pa)
	got := logs(b)
	if !containsPart(got, "söder") {
		t.Fatalf("Bo hörde inte att Alva drack: %q", got)
	}

	// Längre bort hörs det inte.
	putInRoom(g, pb, 1, 1)
	pa.Potions = 1
	logs(b)
	g.drink(pa)
	if got := logs(b); len(got) > 0 {
		t.Fatalf("Bo hörde något långt bort: %q", got)
	}

	// Strid i grannrummet hörs med riktning i stället för del av tornet.
	putInRoom(g, pb, 5, 4) // öster om Alva
	logs(b)
	g.hear(pa.X, pa.Y, "ljud_strid", "allmän text", "dim", pa)
	got = logs(b)
	if containsPart(got, "allmän text") || !containsPart(got, "väster") {
		t.Fatalf("fel stridsljud: %q", got)
	}
}

func containsPart(list []string, part string) bool {
	for _, s := range list {
		if strings.Contains(s, part) {
			return true
		}
	}
	return false
}

// hiddenGame startar en vandring med Alva och Bo där vandrarna inte ser varandra.
func hiddenGame(t *testing.T) (*Game, *Client, *Client, []string) {
	t.Helper()
	g := testGame(t)
	g.cfg.Sight = sightOwn
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	joinLog := logs(a)
	g.Ready(a)
	g.Ready(b)
	if g.phase != phasePlay {
		t.Fatal("spelet startade inte")
	}
	return g, a, b, joinLog
}

func TestLobbyHidesOthers(t *testing.T) {
	g := testGame(t)
	g.cfg.Sight = sightOwn
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	if v := viewOf(g.stateFor(a.player), b.player); !v.Hidden {
		t.Fatal("Bo syns i lobbyn")
	}
	// Alva ser bara porten hon står vid, inte de andra.
	st := g.stateFor(a.player)
	for _, cn := range g.corners() {
		mine := cn.x == a.player.X && cn.y == a.player.Y
		if hasRoom(st, cn.x, cn.y) != mine {
			t.Fatalf("porten i %d,%d: syns %v, borde %v", cn.x, cn.y, !mine, mine)
		}
	}
	// Med synligheten alla syns Bo som förut.
	g.cfg.Sight = sightAll
	if v := viewOf(g.stateFor(a.player), b.player); v.Hidden {
		t.Fatal("Bo syns inte med synligheten alla")
	}
}

func TestJoinTextWithoutCorner(t *testing.T) {
	_, _, _, joinLog := hiddenGame(t)
	for _, line := range joinLog {
		if strings.Contains(line, "Bo") {
			for _, h := range []string{"nordväst", "nordöst", "sydväst", "sydöst"} {
				if strings.Contains(line, h) {
					t.Fatalf("Alva fick veta Bos port: %q", line)
				}
			}
		}
	}
}

func TestDrawnCorners(t *testing.T) {
	seen := map[int]bool{}
	for range 40 {
		g, a, b, _ := hiddenGame(t)
		pa, pb := a.player, b.player
		if pa.Corner == pb.Corner {
			t.Fatal("två vandrare vid samma port")
		}
		for _, p := range []*Player{pa, pb} {
			cn := g.corners()[p.Corner]
			if p.X != cn.x || p.Y != cn.y || !p.known[cn.y][cn.x] {
				t.Fatalf("%s står inte vid sin port: %d,%d", p.Name, p.X, p.Y)
			}
		}
		if !containsPart(logs(a), g.data.Horn[pa.Corner]) {
			t.Fatal("Alva fick inte veta sin port")
		}
		seen[pa.Corner] = true
	}
	if len(seen) < 3 {
		t.Fatalf("Alva började bara vid portarna %v", seen)
	}
}

func TestPortSound(t *testing.T) {
	g, a, b, _ := hiddenGame(t)
	pa, pb := a.player, b.player
	cn := g.corners()[0]
	pa.X, pa.Y = cn.x, cn.y
	putInRoom(g, pb, cn.x+1, cn.y) // öster om porten
	logs(b)
	g.move(pa, g.corners()[0].exit, 0)
	if pa.Status != statusOut {
		t.Fatal("Alva kom inte ut")
	}
	// Porten ligger väster om Bo. Raden om att Alva har gått ut säger inte riktningen.
	if !containsPart(logs(b), "väster") {
		t.Fatal("Bo hörde inte porten")
	}
}

// Spelledarens översikt får inte krascha med egen karta.
func TestGMViewWithOwnMap(t *testing.T) {
	g, _, _ := startTwo(t)
	g.cfg.Sight = sightOwn
	g.cfg.GMMode = true
	g.monsters = append(g.monsters, &Monster{ID: 7, Type: g.data.Monster.Typer[0].ID, X: 4, Y: 4, HP: 1})
	g.reveal(4, 4)
	v := g.GMView(false)
	if len(v.Map.Monsters) != 1 {
		t.Fatalf("översikten ser %d varelser", len(v.Map.Monsters))
	}
}
