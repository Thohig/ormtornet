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
			g := NewGame(d, Config{Reveal: true, SunMax: 60, TurnTime: time.Hour, MinPlayers: 1, GMMode: true, AI: 4, Sight: sight, PvP: pvpOn, Classes: true, Size: size})
			g.AddBots(4)
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
