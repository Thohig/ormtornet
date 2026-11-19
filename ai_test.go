package main

import "testing"

// En vilsen AI-spelare med skatt letar sig ut mot portarna, inte in mot
// kammaren.
func TestLostBotHeadsForPorts(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.bot = &botBrain{pers: &g.data.AI.Personligheter[0], habits: map[string]map[string]int{}}
	for y := range g.size {
		for x := range g.size {
			g.rooms[y][x].Revealed, g.rooms[y][x].Walked = false, false
		}
	}
	r := putInRoom(g, p, 2, g.center())
	r.Doors = [4]bool{false, true, false, true} // öster mot mitten, väster mot kanten
	r.Searched = true
	p.Treasure = 10
	if a, _ := g.botDecide(p); a != "W" {
		t.Fatalf("gick %q, ville ha W mot portarna", a)
	}
	p.Treasure = 0
	g.sun = 4 * g.size
	if a, _ := g.botDecide(p); a != "E" {
		t.Fatalf("utan skatt och med gott om tid gick den %q, ville ha E mot kammaren", a)
	}
}
