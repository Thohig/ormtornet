package main

import (
	"testing"
	"time"
)

func testGame(t *testing.T) *Game {
	t.Helper()
	d, _, err := loadData("data/texter.json")
	if err != nil {
		t.Fatal(err)
	}
	return NewGame(d, Config{SunMax: 10, TurnTime: time.Hour, MinPlayers: 2})
}

func testClient() *Client {
	return &Client{send: make(chan []byte, 10000)}
}

func TestRoundTreasureAndWinner(t *testing.T) {
	g := testGame(t)
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	g.Ready(a)
	if g.phase != phaseLobby {
		t.Fatal("spelet startade innan alla var redo")
	}
	g.Ready(b)
	if g.phase != phasePlay || g.round != 1 {
		t.Fatalf("fas %s runda %d", g.phase, g.round)
	}

	// Alva står i kammaren och tar skatt, Bo väntar.
	pa, pb := a.player, b.player
	g.reveal(center, center)
	pa.X, pa.Y = center, center
	g.Act(a, actTake)
	if g.round != 1 {
		t.Fatal("rundan utfördes innan alla valt")
	}
	g.Act(b, actWait)
	if g.round != 2 || g.sun != 9 {
		t.Fatalf("runda %d sol %d", g.round, g.sun)
	}
	if pa.Treasure <= 0 {
		t.Fatal("ingen skatt togs")
	}

	// Ta skatt utanför kammaren går inte.
	g.Act(b, actTake)
	if pb.Action != "" {
		t.Fatal("skatt kunde väljas utanför kammaren")
	}

	// Alva går ut genom sin port, Bo blir kvar tills solen går ner.
	pa.X, pa.Y = 0, 0
	g.Act(a, "W")
	g.Act(b, actWait)
	if pa.Status != statusOut {
		t.Fatalf("Alva är %s", pa.Status)
	}
	for g.phase == phasePlay {
		g.Act(b, actWait)
	}
	if pb.Status != statusLost {
		t.Fatalf("Bo är %s", pb.Status)
	}
	res := g.results()
	if res[0].Name != "Alva" || !res[0].Winner || res[1].Winner {
		t.Fatalf("fel resultat: %+v", res)
	}

	// Båda trycker Enter: ett nytt torn byggs och spelet börjar om.
	g.Ready(a)
	g.Ready(b)
	if g.phase != phasePlay || pa.Treasure != 0 || pa.Status != statusIn || g.sun != 10 {
		t.Fatalf("omstart misslyckades: fas %s skatt %d", g.phase, pa.Treasure)
	}
}

func TestRoomsStayConsistent(t *testing.T) {
	g := testGame(t)
	for y := range gridSize {
		for x := range gridSize {
			g.reveal(x, y)
		}
	}
	for y := range gridSize {
		for x := range gridSize {
			r := &g.rooms[y][x]
			for d := North; d <= West; d++ {
				nx, ny := x+dirDX[d], y+dirDY[d]
				if !inBounds(nx, ny) {
					if r.Doors[d] {
						t.Fatalf("dörr ut ur tornet vid %d,%d", x, y)
					}
					continue
				}
				if r.Doors[d] != g.rooms[ny][nx].Doors[d.Opp()] {
					t.Fatalf("dörrarna stämmer inte mellan %d,%d och %d,%d", x, y, nx, ny)
				}
			}
		}
	}
}
