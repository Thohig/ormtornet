package main

import (
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
