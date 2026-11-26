package main

import (
	"io"
	"log"
	"os"
	"testing"
	"time"
)

// Tornet stänger aldrig in någon: från varje port når man kammaren genom
// dörrarna, hur rummen än byggs. Rummen byggs här i slumpad ordning från
// en port, som när vandrare går omkring, och ibland ett rum på avstånd,
// som när en spegelskärva eller Kartritaren ser in i ett grannrum.
func TestTowerNeverSealsAPortIn(t *testing.T) {
	d, _, err := loadData("data")
	if err != nil {
		t.Fatal(err)
	}
	type pos struct{ x, y int }
	for i := 0; i < 3000; i++ {
		g := NewGame(d, Config{TurnTime: time.Hour, MinPlayers: 1})
		g.buildTower([]int{9, 11, 13}[i%3])
		cn := g.corners()[i%4]
		seen := map[pos]bool{{cn.x, cn.y}: true}
		frontier := []pos{{cn.x, cn.y}}
		chamber := false
		// rescan letar efter dörrar som tornet har öppnat i rum man redan
		// har varit i, som en vandrare ser när hen går tillbaka.
		rescan := func() {
			for c := range seen {
				r := &g.rooms[c.y][c.x]
				for dd := North; dd <= West; dd++ {
					n := pos{c.x + dirDX[dd], c.y + dirDY[dd]}
					if r.Doors[dd] && g.inBounds(n.x, n.y) && !seen[n] {
						seen[n] = true
						frontier = append(frontier, n)
					}
				}
			}
		}
		for len(frontier) > 0 || func() bool { rescan(); return len(frontier) > 0 }() {
			k := g.rng.IntN(len(frontier))
			c := frontier[k]
			frontier = append(frontier[:k], frontier[k+1:]...)
			r := g.reveal(c.x, c.y)
			if r.Kind == "kammare" {
				chamber = true
				break
			}
			for dd := North; dd <= West; dd++ {
				nx, ny := c.x+dirDX[dd], c.y+dirDY[dd]
				if !g.inBounds(nx, ny) {
					continue
				}
				if !r.Doors[dd] {
					// Ibland ser någon in i rummet bakom en vägg.
					if g.rng.IntN(10) == 0 {
						g.reveal(nx, ny)
					}
					continue
				}
				if n := (pos{nx, ny}); !seen[n] {
					seen[n] = true
					frontier = append(frontier, n)
				}
			}
		}
		if !chamber {
			t.Fatalf("torn %d (%d rum i bredd): instängd i %d rum från porten %d", i, g.size, len(seen), i%4)
		}
	}
}

// Tornet stänger inte in någon när många utforskar samtidigt: alla portar
// bygger rum om vartannat, rum byggs på avstånd (som när någon kommer upp ur
// källarvalven eller en varelse vandrar), och ibland flyttar kammaren. När
// allt som går att nå är byggt når varje port kammaren.
func TestTowerNeverSealsAnyoneIn(t *testing.T) {
	d, _, err := loadData("data")
	if err != nil {
		t.Fatal(err)
	}
	type pos struct{ x, y int }
	// reaches följer dörrarna från (x, y) och säger om man kommer till
	// kammaren, och vilka okända rum man kan se in i på vägen.
	reaches := func(g *Game, x, y int) (bool, []pos) {
		cx, cy := g.chamber()
		seen := map[pos]bool{{x, y}: true}
		q := []pos{{x, y}}
		var unknown []pos
		found := false
		for k := 0; k < len(q); k++ {
			c := q[k]
			if c.x == cx && c.y == cy {
				found = true
			}
			r := &g.rooms[c.y][c.x]
			if !r.Revealed {
				unknown = append(unknown, c)
				continue
			}
			for dd := North; dd <= West; dd++ {
				n := pos{c.x + dirDX[dd], c.y + dirDY[dd]}
				if r.Doors[dd] && g.inBounds(n.x, n.y) && !seen[n] {
					seen[n] = true
					q = append(q, n)
				}
			}
		}
		return found, unknown
	}
	for i := 0; i < 1500; i++ {
		g := NewGame(d, Config{TurnTime: time.Hour, MinPlayers: 1, Seed: uint64(i + 1)})
		g.buildTower([]int{9, 11, 13, 15}[i%4])
		moveAt := -1
		if i%3 == 0 {
			moveAt = 5 + g.rng.IntN(40)
		}
		for step := 0; ; step++ {
			var front []pos
			for _, cn := range g.corners() {
				_, u := reaches(g, cn.x, cn.y)
				front = append(front, u...)
			}
			if len(front) == 0 {
				break
			}
			c := front[g.rng.IntN(len(front))]
			g.reveal(c.x, c.y)
			if g.rng.IntN(6) == 0 {
				x, y := g.rng.IntN(g.size), g.rng.IntN(g.size)
				g.build(x, y, false)
				if ok, u := reaches(g, x, y); !ok && len(u) == 0 {
					t.Fatalf("torn %d: rummet (%d, %d) som byggdes på avstånd är instängt", i, x, y)
				}
			}
			if step == moveAt {
				g.moveChamber()
			}
		}
		for ci, cn := range g.corners() {
			if ok, _ := reaches(g, cn.x, cn.y); !ok {
				t.Fatalf("torn %d (%d rum i bredd): porten %d når inte kammaren", i, g.size, ci)
			}
		}
	}
}

// Ingen AI-vandrare blir instängd under riktiga vandringar: från där hen
// står finns alltid en väg till kammaren eller till ett rum som inte är
// byggt än.
func TestNoWalkerSealedInDuringPlay(t *testing.T) {
	d, _, err := loadData("data")
	if err != nil {
		t.Fatal(err)
	}
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)
	type pos struct{ x, y int }
	for seed := uint64(1); seed <= 3000; seed++ {
		n := 2 + int(seed%7)
		g := NewGame(d, Config{Seed: seed, AI: n, MinPlayers: 1, TurnTime: time.Hour})
		g.mu.Lock()
		g.AddBots(n)
		simDealClasses(g, seed)
		g.begin()
		if g.phase == phaseChoose {
			g.startGame()
		}
		for guard := 0; g.phase == phasePlay && guard < 20*g.sunMax; guard++ {
			round := g.round
			for _, p := range g.players {
				if p != nil && p.bot != nil {
					g.botAct(p, round)
				}
			}
			if g.phase == phasePlay && g.round == round {
				g.resolve()
			}
			if g.phase != phasePlay || len(g.shutDoors) > 0 {
				continue // dörrar som tornet har slutit en stund öppnas igen
			}
			cx, cy := g.chamber()
			for _, p := range g.players {
				if p == nil || p.Status != statusIn || p.Cellar != 0 {
					continue
				}
				seen := map[pos]bool{{p.X, p.Y}: true}
				q := []pos{{p.X, p.Y}}
				ok := false
				for k := 0; k < len(q) && !ok; k++ {
					c := q[k]
					r := &g.rooms[c.y][c.x]
					if (c.x == cx && c.y == cy) || !r.Revealed {
						ok = true
						break
					}
					for dd := North; dd <= West; dd++ {
						nn := pos{c.x + dirDX[dd], c.y + dirDY[dd]}
						if r.Doors[dd] && g.inBounds(nn.x, nn.y) && !seen[nn] {
							seen[nn] = true
							q = append(q, nn)
						}
					}
				}
				if !ok {
					g.mu.Unlock()
					t.Fatalf("frö %d (%d AI), runda %d: %s är instängd i %d rum", seed, n, g.round, p.Name, len(seen))
				}
			}
		}
		g.mu.Unlock()
	}
}
