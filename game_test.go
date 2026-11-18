package main

import (
	"testing"
	"time"
)

func testGame(t *testing.T) *Game {
	t.Helper()
	d, _, err := loadData("data")
	if err != nil {
		t.Fatal(err)
	}
	// Inga slumpade varelser eller fällor; testerna placerar dem själva.
	d.Monster.Chans = 0
	d.Regler.Handelsechans = 0
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
	g.deck = []string{"sover"} // ormen ska inte vakna i det här testet
	g.Act(a, actTake, 0)
	if g.round != 1 {
		t.Fatal("rundan utfördes innan alla valt")
	}
	g.Act(b, actWait, 0)
	if g.round != 2 || g.sun != 9 {
		t.Fatalf("runda %d sol %d", g.round, g.sun)
	}
	if pa.Treasure <= 0 {
		t.Fatal("ingen skatt togs")
	}

	// Ta skatt utanför kammaren går inte.
	g.Act(b, actTake, 0)
	if pb.Action != "" {
		t.Fatal("skatt kunde väljas utanför kammaren")
	}

	// Alva går ut genom sin port, Bo blir kvar tills solen går ner.
	pa.X, pa.Y = 0, 0
	g.Act(a, "W", 0)
	g.Act(b, actWait, 0)
	if pa.Status != statusOut {
		t.Fatalf("Alva är %s", pa.Status)
	}
	for g.phase == phasePlay {
		g.Act(b, actWait, 0)
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

// startTwo startar ett spel med två spelare och returnerar dem.
func startTwo(t *testing.T) (*Game, *Client, *Client) {
	t.Helper()
	g := testGame(t)
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	g.Ready(a)
	g.Ready(b)
	if g.phase != phasePlay {
		t.Fatal("spelet startade inte")
	}
	return g, a, b
}

// step låter båda välja drag så att rundan utförs.
func step(g *Game, a, b *Client, actA, actB string) {
	if a.player.Status == statusIn && a.player.Cellar == 0 {
		g.Act(a, actA, 0)
	}
	if b.player.Status == statusIn && b.player.Cellar == 0 {
		g.Act(b, actB, 0)
	}
}

// putInRoom ställer spelaren i ett vanligt rum mitt i tornet med dörrar åt alla håll.
func putInRoom(g *Game, p *Player, x, y int) *Room {
	r := g.reveal(x, y)
	r.Kind, r.Event, r.Searched = "rum", "", false
	p.X, p.Y = x, y
	return r
}

func TestTraps(t *testing.T) {
	cases := []struct {
		id    string
		check func(t *testing.T, g *Game, p *Player, sunBefore int)
	}{
		{"fallgaller", func(t *testing.T, g *Game, p *Player, _ int) {
			if p.Stuck != 2 {
				t.Fatalf("fast %d", p.Stuck)
			}
		}},
		{"pilar", func(t *testing.T, g *Game, p *Player, _ int) {
			if p.HP != g.data.Regler.Liv-2 {
				t.Fatalf("liv %d", p.HP)
			}
		}},
		{"dimma", func(t *testing.T, g *Game, p *Player, _ int) {
			if p.Poison != 2 || p.HP != g.data.Regler.Liv-1 {
				t.Fatalf("gift %d liv %d", p.Poison, p.HP)
			}
		}},
		{"ormgrop", func(t *testing.T, g *Game, p *Player, _ int) {
			if p.Stuck != 1 || p.HP != g.data.Regler.Liv-2 {
				t.Fatalf("fast %d liv %d", p.Stuck, p.HP)
			}
		}},
		{"skinn", func(t *testing.T, g *Game, p *Player, _ int) {
			if p.Stuck != 1 {
				t.Fatalf("fast %d", p.Stuck)
			}
		}},
		{"klocka", func(t *testing.T, g *Game, p *Player, sunBefore int) {
			if g.sun != sunBefore-2 {
				t.Fatalf("sol %d före %d", g.sun, sunBefore)
			}
		}},
		{"morker", func(t *testing.T, g *Game, p *Player, _ int) {
			if p.Dark != 4 {
				t.Fatalf("mörker %d", p.Dark)
			}
			st := g.stateFor(p)
			if len(st.Rooms) != 1 {
				t.Fatalf("ser %d rum i mörkret", len(st.Rooms))
			}
		}},
		{"rasande_golv", func(t *testing.T, g *Game, p *Player, _ int) {
			if p.Cellar != 2 {
				t.Fatalf("källare %d", p.Cellar)
			}
		}},
		{"roterande", func(t *testing.T, g *Game, p *Player, _ int) {
			if p.X == 4 && p.Y == 2 {
				t.Fatal("spelaren kastades inte ut")
			}
		}},
		{"eko", func(t *testing.T, g *Game, p *Player, _ int) {}},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			g, a, _ := startTwo(t)
			p := a.player
			r := putInRoom(g, p, 4, 2)
			r.Doors = [4]bool{true, true, true, true}
			for d := North; d <= West; d++ {
				g.reveal(4+dirDX[d], 2+dirDY[d]).Event = ""
			}
			before := g.sun
			g.trigger(p, g.data.trap(c.id), 0)
			c.check(t, g, p, before)
		})
	}
}

func TestStuckPoisonCellarAndDeath(t *testing.T) {
	g, a, b := startTwo(t)
	p := a.player
	putInRoom(g, p, 4, 2)
	g.rooms[2][4].Doors[North] = true
	g.reveal(4, 1)

	// Fast: draget norrut misslyckas så länge gallret är nere.
	g.trigger(p, g.data.trap("fallgaller"), 0)
	step(g, a, b, "N", actWait)
	step(g, a, b, "N", actWait)
	if p.Y != 2 {
		t.Fatal("kom loss för tidigt")
	}
	step(g, a, b, "N", actWait)
	if p.Y != 1 {
		t.Fatalf("kom inte loss, fast %d", p.Stuck)
	}

	// Gift tickar en gång per runda; en dryck tar bort det.
	hp := p.HP
	g.trigger(p, g.data.trap("dimma"), 0)
	step(g, a, b, actWait, actWait)
	if p.HP != hp-2 {
		t.Fatalf("liv %d, väntade %d", p.HP, hp-2)
	}
	p.Potions = 1
	step(g, a, b, actDrink, actWait)
	if p.Poison != 0 || p.Potions != 0 {
		t.Fatal("drycken hjälpte inte")
	}

	// Källaren: dold för andra, sedan upp i ett annat rum.
	g.trigger(p, g.data.trap("rasande_golv"), 0)
	if !g.stateFor(b.player).Players[0].Hidden {
		t.Fatal("spelaren i källaren syns för andra")
	}
	step(g, a, b, actWait, actWait)
	step(g, a, b, actWait, actWait)
	if p.Cellar != 0 || !g.rooms[p.Y][p.X].Revealed {
		t.Fatal("kom inte upp ur källaren")
	}

	// Döden.
	g.damage(p, 100, true)
	if p.Status != statusDead {
		t.Fatal("dog inte")
	}
}

func TestSearch(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	seen := map[string]bool{}
	for i := 0; i < 300; i++ {
		r := putInRoom(g, p, 4, 2)
		r.Doors = [4]bool{}
		p.Status, p.HP, p.Cellar, p.X, p.Y = statusIn, 99, 0, 4, 2
		gold, pot := p.Treasure, p.Potions
		g.search(p)
		switch {
		case p.Treasure > gold:
			seen["skatt"] = true
		case p.Potions > pot:
			seen["dryck"] = true
		case r.Doors != [4]bool{}:
			seen["gang"] = true
		}
		if !r.Searched {
			t.Fatal("rummet markerades inte som genomsökt")
		}
	}
	for _, k := range []string{"skatt", "dryck", "gang"} {
		if !seen[k] {
			t.Errorf("hittade aldrig %s", k)
		}
	}
}

func TestSerpentDeck(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	g.reveal(center, center)
	for d := North; d <= West; d++ {
		g.reveal(center+dirDX[d], center+dirDY[d]).Event = ""
	}
	full := len(g.deck)
	sleepFull := g.deckSleep
	if g.sleepLevel() != 0 {
		t.Fatal("ormen ska sova djupt från början")
	}

	// Båda står i kammaren och tar skatt tills ormen vaknar.
	taken := 0
	for g.wakes == 0 {
		pa.X, pa.Y, pb.X, pb.Y = center, center, center, center
		pa.HP, pb.HP = 99, 99
		g.take(pa)
		taken++
		if taken > full {
			t.Fatal("ormen vaknade aldrig")
		}
		if g.wakes == 0 && len(g.deck) != full-taken {
			t.Fatalf("sovkortet togs inte bort: %d kort kvar", len(g.deck))
		}
	}
	// Båda kastades ut ur kammaren, förlorade liv och hälften av skatten.
	for _, p := range []*Player{pa, pb} {
		if g.rooms[p.Y][p.X].Kind == "kammare" {
			t.Fatalf("%s är kvar i kammaren", p.Name)
		}
		if p.HP != 99-g.data.Ormen.Skada {
			t.Fatalf("%s har %d liv", p.Name, p.HP)
		}
	}
	// Ny lek med färre sovkort.
	if g.deckSleep != sleepFull-g.data.Ormen.FarreEfterVakning {
		t.Fatalf("sovkort %d, före %d", g.deckSleep, sleepFull)
	}
}

func TestSerpentLosesHalfTreasure(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	g.reveal(center, center)
	p.X, p.Y, p.Treasure = center, center, 100
	g.deck = []string{"vaknar"}
	g.take(p)
	if p.Treasure != 50 {
		t.Fatalf("skatt %d, väntade 50", p.Treasure)
	}
}

func TestSleepLevelRises(t *testing.T) {
	g, _, _ := startTwo(t)
	levels := []int{g.sleepLevel()}
	for len(g.deck) > 0 {
		if c := g.data.card(g.deck[0]); c.Vaknar {
			g.deck = g.deck[1:]
			g.deck = append(g.deck, c.ID)
			if allWake(g) {
				break
			}
			continue
		}
		g.deck = g.deck[1:]
		levels = append(levels, g.sleepLevel())
	}
	if levels[0] != 0 || levels[len(levels)-1] != 2 {
		t.Fatalf("nivåer %v", levels)
	}
}

func allWake(g *Game) bool {
	for _, id := range g.deck {
		if !g.data.card(id).Vaknar {
			return false
		}
	}
	return true
}

func TestBeats(t *testing.T) {
	if !beats(actParry, actAttack) || !beats(actAttack, actDodge) || !beats(actDodge, actParry) {
		t.Fatal("fel ordning i sten-sax-påse")
	}
	for _, a := range []string{actAttack, actParry, actDodge} {
		if beats(a, a) {
			t.Fatalf("%s slår sig själv", a)
		}
	}
}

// fightSetup ställer Alva i ett rum med en varelse som alltid väljer mc.
func fightSetup(t *testing.T, mc string) (*Game, *Client, *Client, *Monster) {
	g, a, b := startTwo(t)
	putInRoom(g, a.player, 4, 2)
	typ := g.data.monsterType("fjallvakt")
	typ.Vanor = map[string]int{mc: 1}
	m := &Monster{ID: 99, Type: "fjallvakt", X: 4, Y: 2, HP: typ.Liv}
	g.monsters = append(g.monsters, m)
	return g, a, b, m
}

func TestDuel(t *testing.T) {
	cases := []struct {
		player, monster string
		playerLoses     bool
		monsterLoses    bool
	}{
		{actParry, actAttack, false, true},
		{actAttack, actParry, true, false},
		{actAttack, actDodge, false, true},
		{actDodge, actParry, false, true},
		{actAttack, actAttack, true, true},
		{actParry, actParry, false, false},
		{actDodge, actDodge, false, false},
	}
	for _, c := range cases {
		t.Run(c.player+"_mot_"+c.monster, func(t *testing.T) {
			g, a, b, m := fightSetup(t, c.monster)
			p := a.player
			hp, mhp := p.HP, m.HP
			step(g, a, b, c.player, actWait)
			if (p.HP < hp) != c.playerLoses {
				t.Fatalf("spelarens liv %d -> %d", hp, p.HP)
			}
			if (m.HP < mhp) != c.monsterLoses {
				t.Fatalf("varelsens liv %d -> %d", mhp, m.HP)
			}
		})
	}
}

func TestDefeatMonsterGivesGold(t *testing.T) {
	g, a, b, m := fightSetup(t, actDodge)
	p := a.player
	for i := 0; i < 10 && g.monsterByID(m.ID) != nil; i++ {
		step(g, a, b, actAttack, actWait)
	}
	if g.monsterByID(m.ID) != nil {
		t.Fatal("varelsen dog aldrig")
	}
	if p.Treasure != g.data.monsterType("fjallvakt").Guld {
		t.Fatalf("skatt %d", p.Treasure)
	}
}

func TestFleeAndUnprepared(t *testing.T) {
	g, a, b, _ := fightSetup(t, actAttack)
	p := a.player
	g.rooms[2][4].Doors[North] = true
	g.reveal(4, 1)
	hp := p.HP
	// Den som söker igenom rummet medan varelsen står där blir träffad.
	step(g, a, b, actSearch, actWait)
	if p.HP >= hp || g.rooms[2][4].Searched {
		t.Fatal("varelsen fick inget fritt anfall")
	}
	// Flykt: ett fritt anfall till, sedan är spelaren ute ur rummet.
	hp = p.HP
	step(g, a, b, "N", actWait)
	if p.HP >= hp || p.Y != 1 {
		t.Fatalf("flykten gick fel: liv %d -> %d, y %d", hp, p.HP, p.Y)
	}
}

func TestDeadPlayerControlsMonsters(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	putInRoom(g, pb, 4, 2)
	r := putInRoom(g, pa, 4, 1)
	r.Doors[South] = true
	g.rooms[2][4].Doors[North] = true
	m := &Monster{ID: 7, Type: "krypare", X: 4, Y: 2, HP: 2}
	g.monsters = append(g.monsters, m)
	pb.X, pb.Y = 4, 3
	g.reveal(4, 3)

	g.damage(pa, 100, false)
	if !g.isController(pa) {
		t.Fatal("den döda spelaren styr inte varelserna")
	}
	// Rundan väntar på den döda spelaren också.
	g.Act(b, actWait, 0)
	if g.round != 1 {
		t.Fatal("rundan utfördes utan den döda spelarens val")
	}
	// Varelsen kan inte gå till ett okänt rum eller in i en vägg.
	if g.canMonsterGo(m, East) && !g.rooms[2][5].Revealed {
		t.Fatal("varelsen kan gå till ett okänt rum")
	}
	g.Act(a, "N", m.ID)
	if g.round != 2 || m.Y != 1 {
		t.Fatalf("varelsen flyttades inte: runda %d, y %d", g.round, m.Y)
	}
	// Bara den första som dör blir styrande.
	g.damage(pb, 100, false)
	if g.isController(pb) {
		t.Fatal("två spelare styr varelserna")
	}
}
