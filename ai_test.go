package main

import "testing"

// En vilsen AI-spelare med skatt letar sig ut mot portarna, inte in mot
// kammaren.
func TestLostBotHeadsForPorts(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.bot = &botBrain{pers: &g.data.AI.Personligheter[0]}
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

// Livsgränserna följer vandrarens maxliv, så att en försiktig AI-spelare med
// få liv inte vänder redan i porten.
func TestBotLifeLimitFollowsMaxHP(t *testing.T) {
	g := testGame(t)
	// Ett jämnt antal liv, så att hälften går jämnt upp.
	g.base.Regler.Liv = 6
	base := g.base.Regler.Liv
	for _, c := range []struct{ limit, maxHP, want int }{
		{3, base, 3},
		{3, base * 2, 6},
		{3, base / 2, 2},
		{2, base / 2, 1},
		{3, 0, 3},
	} {
		if got := g.botLifeLimit(c.limit, c.maxHP); got != c.want {
			t.Errorf("gräns %d med %d liv gav %d, ville ha %d", c.limit, c.maxHP, got, c.want)
		}
	}
}

// En datorspelare som hör ormen letar åt det hållet, och glömmer det när
// ormen flyttar.
func TestBotFollowsSerpentClues(t *testing.T) {
	g := testGame(t)
	g.AddBots(1)
	var bot *Player
	for _, p := range g.players {
		if p != nil && p.bot != nil {
			bot = p
		}
	}
	bot.heardSerpent(2, 6)
	if !bot.bot.hasHint || bot.bot.hint != [2]int{2, 6} {
		t.Fatalf("ledtråden: %+v", bot.bot)
	}
	g.data.Ormen.Spridning = 1 // annars finns ingen annan plats att flytta till
	g.reveal(g.chamber())
	g.moveChamber()
	if bot.bot.hasHint {
		t.Fatal("datorspelaren minns den gamla platsen")
	}
}

// En AI-spelare med ägget går ut innan det kläcks, och en med ormens öga
// som har vänt tar inga fler nävar när den har guld.
func TestBotKnowsEggAndEye(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.bot = &botBrain{pers: &g.data.AI.Personligheter[0]}
	g.sun = 4 * g.size
	room := &g.rooms[p.Y][p.X]
	if !room.HasExit {
		t.Fatal("vandraren står inte vid sin port")
	}
	exit := dirLetter[room.Exit]
	g.give(p, "agg")
	p.has("agg").rounds = 11
	if act, _ := g.botDecide(p); act != exit {
		t.Fatalf("med ägget nära att kläckas valde den %q, ville ha %q ut", act, exit)
	}

	p.items = nil
	cx, cy := g.chamber()
	p.X, p.Y = cx, cy
	g.reveal(cx, cy)
	p.known[cy][cx] = true
	p.Treasure = 10
	g.give(p, "ormens_oga")
	p.has("ormens_oga").turned = true
	for range 30 {
		if act, _ := g.botDecide(p); act == actTake {
			t.Fatal("tog en näve fast ormens öga har vänt")
		}
	}
}
