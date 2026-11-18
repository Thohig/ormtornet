package main

import (
	"testing"
	"time"
)

func TestTowerSizes(t *testing.T) {
	for _, size := range []int{9, 11, 13} {
		g := testGame(t)
		g.cfg.Size = size
		g.newTower()
		if g.size != size || len(g.rooms) != size || len(g.rooms[0]) != size {
			t.Fatalf("tornet är %d, ville ha %d", g.size, size)
		}
		c := g.center()
		if c != size/2 {
			t.Fatalf("mitten %d", c)
		}
		if d, ok := g.fixedDoors(c, c); !ok || d != [4]bool{true, true, true, true} {
			t.Fatal("kammaren har inte fyra dörrar")
		}
		for i, cn := range g.corners() {
			r := g.reveal(cn.x, cn.y)
			if r.Kind != "ingang" || !r.HasExit || g.cornerIndex(cn.x, cn.y) != i {
				t.Fatalf("porten i %d,%d", cn.x, cn.y)
			}
		}
		if g.data.region(size-1, size-1, size) != g.data.Regioner[2][2] {
			t.Fatal("fel del av tornet i hörnet")
		}
	}
}

func TestSizeChangeInLobby(t *testing.T) {
	g := testGame(t)
	a := testClient()
	g.Join(a, "Alva")
	s := settingsFrom(g.cfg)
	s.Storlek = 13
	if _, err := g.ApplySettings(s); err != nil {
		t.Fatal(err)
	}
	if g.size != 13 {
		t.Fatalf("tornet är %d", g.size)
	}
	cn := g.corners()[a.player.Corner]
	if a.player.X != cn.x || a.player.Y != cn.y || len(a.player.known) != 13 {
		t.Fatal("Alva står inte vid sin port i det nya tornet")
	}
	if st := g.stateFor(a.player); st.Size != 13 {
		t.Fatalf("klienten får storleken %d", st.Size)
	}
}

// En hel vandring med fyra AI-spelare i varje storlek, för att se att inget
// går utanför tornet.
func TestBotsPlayEverySize(t *testing.T) {
	for _, size := range []int{9, 11, 13} {
		for _, sight := range []string{sightAll, sightOwn} {
			d, _, err := loadData("data")
			if err != nil {
				t.Fatal(err)
			}
			d.AI.BetanketidMs = []int{3600000, 3600000}
			g := NewGame(d, Config{Reveal: true, SunMax: 40, TurnTime: time.Hour, MinPlayers: 1, GMMode: true, AI: 8, Sight: sight, PvP: pvpOn, Classes: true, Size: size})
			g.AddBots(8)
			g.startGame()
			for g.phase == phasePlay && g.round < 200 {
				r := g.round
				for _, p := range g.players {
					if p != nil && g.round == r {
						g.botAct(p, r)
					}
				}
				if g.phase == phasePlay && g.round == r {
					g.resolve()
				}
			}
			if g.phase != phaseEnd {
				t.Fatalf("storlek %d: vandringen tog inte slut", size)
			}
		}
	}
}

func TestAutoSizeAndPorts(t *testing.T) {
	g := testGame(t)
	g.cfg.GMMode = true
	g.cfg.MinPlayers = 1
	for _, c := range []struct{ bots, size, ports int }{{4, 9, 4}, {5, 11, 8}, {6, 11, 8}, {7, 13, 8}, {8, 13, 8}} {
		for g.botCount() < c.bots {
			g.AddBots(1)
		}
		if g.size != c.size || len(g.corners()) != c.ports {
			t.Fatalf("%d vandrare: tornet %d med %d portar, ville ha %d med %d", c.bots, g.size, len(g.corners()), c.size, c.ports)
		}
		seen := map[int]bool{}
		for _, p := range g.players {
			if p == nil {
				continue
			}
			if seen[p.Corner] {
				t.Fatalf("%d vandrare: två vid samma port", c.bots)
			}
			seen[p.Corner] = true
			cn := g.corners()[p.Corner]
			if p.X != cn.x || p.Y != cn.y {
				t.Fatalf("%s står inte vid sin port", p.Name)
			}
		}
	}
	// En mittport har tre dörrar inåt och en utgång utåt.
	cn := g.corners()[4]
	r := g.reveal(cn.x, cn.y)
	open := 0
	for _, d := range r.Doors {
		if d {
			open++
		}
	}
	if r.Kind != "ingang" || open != 3 || r.Exit != North {
		t.Fatalf("mittporten: %s, %d dörrar, utgång %d", r.Kind, open, r.Exit)
	}
	// Solen och ormleken växer med tornet.
	g.cfg.SunMax = 40
	g.startGame()
	if g.sunMax != 60 {
		t.Fatalf("solen %d i ett torn på 13, ville ha 60", g.sunMax)
	}
	if g.deckSleep != 13+2*4 {
		t.Fatalf("%d sovkort med åtta vandrare", g.deckSleep)
	}
}

func TestManualSizeNeverTooSmall(t *testing.T) {
	g := testGame(t)
	g.cfg.GMMode = true
	g.cfg.Size = 9
	g.AddBots(6)
	if g.size != 11 {
		t.Fatalf("tornet är %d för sex vandrare", g.size)
	}
	g.cfg.Size = 13
	g.newTower()
	if g.size != 13 {
		t.Fatalf("spelledaren valde 13, fick %d", g.size)
	}
}
