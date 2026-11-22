package main

import (
	"encoding/json"
	"strings"
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
	noWander(d)
	return NewGame(d, Config{Reveal: true, SunMax: 10, TurnTime: time.Hour, MinPlayers: 2})
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
	g.reveal(g.center(), g.center())
	pa.X, pa.Y = g.center(), g.center()
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
	for y := range g.size {
		for x := range g.size {
			g.reveal(x, y)
		}
	}
	for y := range g.size {
		for x := range g.size {
			r := &g.rooms[y][x]
			for d := North; d <= West; d++ {
				nx, ny := x+dirDX[d], y+dirDY[d]
				if !g.inBounds(nx, ny) {
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
	r.Kind, r.Event, r.Searched, r.Walked = "rum", "", false, true
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
	g.reveal(g.center(), g.center())
	for d := North; d <= West; d++ {
		g.reveal(g.center()+dirDX[d], g.center()+dirDY[d]).Event = ""
	}
	full := len(g.deck)
	sleepFull := g.deckSleep
	if g.sleepLevel() != 0 {
		t.Fatal("ormen ska sova djupt från början")
	}

	// Båda står i kammaren och tar skatt tills ormen vaknar.
	taken := 0
	for g.wakes == 0 {
		pa.X, pa.Y, pb.X, pb.Y = g.center(), g.center(), g.center(), g.center()
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
	g.reveal(g.center(), g.center())
	p.X, p.Y, p.Treasure = g.center(), g.center(), 100
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

// Den som dör väntar ingen på, och varelserna sköts av tornet.
func TestDeadPlayerIsNotWaitedFor(t *testing.T) {
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
	g.Act(a, "N", m.ID) // en död spelare kan inte välja något
	if pa.Action != "" {
		t.Fatalf("den döda spelaren valde %q", pa.Action)
	}
	g.Act(b, actWait, 0)
	if g.round != 2 {
		t.Fatal("rundan väntade på den döda spelaren")
	}
	if st := g.stateFor(pa); st.Phase != phasePlay {
		t.Fatalf("fas %s", st.Phase)
	}
	// Varelsen kan inte gå genom en vägg, men väl till ett rum som inte finns än.
	g.rooms[2][4].Doors[East], g.rooms[2][4].Doors[West] = false, true
	if g.canMonsterGo(m, East) {
		t.Fatal("varelsen kan gå genom väggen")
	}
	if !g.rooms[2][3].Revealed && !g.canMonsterGo(m, West) {
		t.Fatal("varelsen kan inte gå till ett rum som inte finns än")
	}
}

// logs plockar ut loggtexterna en klient har fått.
func logs(c *Client) []string {
	var out []string
	for {
		select {
		case b := <-c.send:
			var m logMsg
			if json.Unmarshal(b, &m) == nil && m.T == "log" {
				out = append(out, m.Text)
			}
		default:
			return out
		}
	}
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func TestRevealLevels(t *testing.T) {
	g := testGame(t)
	g.cfg.SunMax = 100
	g.sunMax = 100
	want := map[int]int{100: 1, 71: 1, 70: 2, 41: 2, 40: 3, 16: 3, 15: 4, 0: 4}
	for sun, lvl := range want {
		if got := g.levelFor(sun); got != lvl {
			t.Errorf("sol %d: nivå %d, väntade %d", sun, got, lvl)
		}
	}
	g.cfg.Reveal = false
	if g.levelFor(0) != 1 {
		t.Error("avstängt avslöjande ska stanna på nivå 1")
	}
}

func TestRevealDifferentCluesAndExplanations(t *testing.T) {
	g, a, b := startTwo(t)
	a.player.Seen["fallgaller"] = true
	logs(a)
	logs(b)
	// Solen sjunker under alla gränser: nivå 4 på en gång.
	g.sun = 1
	g.updateLevel()
	if g.level != 4 {
		t.Fatalf("nivå %d", g.level)
	}
	la, lb := logs(a), logs(b)
	lvl4 := g.data.Avslojande.Ledtradar["4"]
	var ca, cb string
	for _, c := range lvl4 {
		if contains(la, c) {
			ca = c
		}
		if contains(lb, c) {
			cb = c
		}
	}
	if ca == "" || cb == "" || ca == cb {
		t.Fatalf("ledtrådarna ska finnas och skilja sig: %q och %q", ca, cb)
	}
	expl := g.data.Avslojande.Forklaringar["fallgaller"]
	if !contains(la, expl) || contains(lb, expl) {
		t.Fatal("bara den som mött fallgallret ska få dess förklaring")
	}
	// Beskrivningarna följer nivån.
	room := g.reveal(g.center(), g.center())
	if d := g.describe(a.player, room, g.center(), g.center()); !contains(g.data.Skattkammare.Beskrivningar["4"], d) {
		t.Fatalf("fel beskrivning på nivå 4: %q", d)
	}
}

func TestRevealDisabled(t *testing.T) {
	g, a, _ := startTwo(t)
	g.cfg.Reveal = false
	g.sun = 1
	g.updateLevel()
	g.endGame()
	for _, l := range logs(a) {
		if contains(g.data.Avslojande.Epilog, l) {
			t.Fatal("epilogen visades fast avslöjandet är avstängt")
		}
	}
	if g.level != 1 {
		t.Fatalf("nivå %d", g.level)
	}
}

func TestEpilogue(t *testing.T) {
	g, a, _ := startTwo(t)
	logs(a)
	g.endGame()
	l := logs(a)
	if !contains(l, g.data.Avslojande.Epilog[0]) {
		t.Fatal("epilogen saknas")
	}
	if !contains(l, g.data.Avslojande.EpilogForlorad[0]) {
		t.Fatal("den som är kvar inne ska få veta var den är")
	}
}

// simulate spelar en hel omgång med en människa som bara väntar och tre
// datorspelare. Timrar används inte; testet driver rundorna själv.
func simulate(t *testing.T, sun int, seed ...uint64) (*Game, *Player) {
	t.Helper()
	d, _, err := loadData("data")
	if err != nil {
		t.Fatal(err)
	}
	d.AI.BetanketidMs = []int{3600000, 3600000} // timrarna ska aldrig hinna gå
	cfg := Config{Reveal: true, SunMax: sun, TurnTime: time.Hour, MinPlayers: 2}
	if len(seed) > 0 {
		cfg.Seed = seed[0]
	}
	g := NewGame(d, cfg)
	if n := g.AddBots(3); n != 3 {
		t.Fatalf("lade till %d datorspelare", n)
	}
	c := testClient()
	g.Join(c, "Människa")
	g.Ready(c)
	if g.phase != phasePlay {
		t.Fatal("spelet startade inte med datorspelare")
	}
	for guard := 0; g.phase == phasePlay; guard++ {
		if guard > 10*sun {
			t.Fatal("omgången tog aldrig slut")
		}
		round := g.round
		for _, p := range g.players {
			if p != nil && p.bot != nil {
				g.botAct(p, round)
			}
		}
		if g.phase == phasePlay && g.round == round && c.player.Status == statusIn && c.player.Cellar == 0 {
			g.Act(c, actWait, 0)
		}
		if g.phase == phasePlay && g.round == round {
			g.resolve() // som när turtiden går ut
		}
		for len(c.send) > 0 {
			<-c.send
		}
	}
	return g, c.player
}

func TestBotsPlayWholeGames(t *testing.T) {
	escaped, rich, dead, total := 0, 0, 0, 0
	for i := 0; i < 300; i++ {
		g, _ := simulate(t, 40)
		for _, p := range g.players {
			if p == nil || p.bot == nil {
				continue
			}
			total++
			switch p.Status {
			case statusOut:
				escaped++
				if p.Treasure > 0 {
					rich++
				}
			case statusDead:
				dead++
			}
		}
	}
	t.Logf("datorspelare: %d av %d kom ut, %d med skatt, %d dog", escaped, total, rich, dead)
	if rich*4 < total {
		t.Errorf("för få datorspelare kom ut med skatt: %d av %d", rich, total)
	}
}

func TestBotsNeverTakeHumanNamesOrNeedReady(t *testing.T) {
	g := testGame(t)
	g.AddBots(2)
	c := testClient()
	g.Join(c, g.data.AI.Namn[0])
	if c.player == nil || c.player.bot != nil || c.player.Name == g.data.AI.Namn[0] {
		t.Fatal("en människa tog över en datorspelare via namnet")
	}
	if g.phase != phaseLobby {
		t.Fatal("spelet startade innan människan var redo")
	}
	g.Ready(c)
	if g.phase != phasePlay {
		t.Fatal("datorspelarna hindrade starten")
	}
}

// Den enda som inte är redo laddar om sidan i lobbyn. Vandringen ska inte
// starta utan hen medan hen är på väg tillbaka, men väl om hen inte kommer.
func TestReloadInLobbyWaits(t *testing.T) {
	g := testGame(t)
	g.leaveWait = 50 * time.Millisecond
	// Servern håller låset när den tar emot något; timern tar det själv.
	locked := func(f func()) { g.mu.Lock(); defer g.mu.Unlock(); f() }
	phase := func() (ph string) { locked(func() { ph = g.phase }); return }
	a, b, c, b2 := testClient(), testClient(), testClient(), testClient()
	locked(func() {
		g.Join(a, "Alva")
		g.Join(b, "Bo")
		g.Join(c, "Cia")
		g.Ready(a)
		g.Ready(c)
		g.Leave(b)
	})
	if ph := phase(); ph != phaseLobby {
		t.Fatalf("fasen är %q: vandringen började medan Bo laddade om", ph)
	}
	locked(func() { g.Join(b2, "Bo") })
	time.Sleep(100 * time.Millisecond)
	if ph := phase(); ph != phaseLobby || b2.player == nil {
		t.Fatalf("fasen är %q: vandringen började fast Bo kom tillbaka", ph)
	}
	locked(func() { g.Leave(b2) }) // Bo går och kommer inte tillbaka
	time.Sleep(100 * time.Millisecond)
	if phase() == phaseLobby {
		t.Fatal("vandringen startade inte när Bo inte kom tillbaka")
	}
}

// Ett tryck mot porten utan guld är en fråga; det andra går ut. Med guld
// går man ut direkt.
func TestEmptyExitAsksFirst(t *testing.T) {
	g := testGame(t)
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	g.Ready(a)
	g.Ready(b)
	pa, pb := a.player, b.player
	exit := func(p *Player) string {
		r := g.rooms[p.Y][p.X]
		if !r.HasExit {
			t.Fatalf("%s står inte vid en port", p.Name)
		}
		for code, d := range dirCodes {
			if d == r.Exit {
				return code
			}
		}
		return ""
	}
	g.Act(a, exit(pa), 0)
	if pa.Action != "" {
		t.Fatalf("Alva valde %q utan att porten frågade", pa.Action)
	}
	g.round++ // rundan hann ta slut innan Alva svarade
	g.Act(a, exit(pa), 0)
	if pa.Action != exit(pa) {
		t.Fatal("det andra trycket, i nästa runda, gick inte ut")
	}
	pb.Treasure = 5
	g.Act(b, exit(pb), 0)
	if g.phase == phasePlay && pb.Status != statusOut {
		t.Fatal("Bo med guld kom inte ut på ett tryck")
	}
	if pa.Status != statusOut {
		t.Fatal("Alva gick inte ut efter två tryck")
	}
}

// Den som kommer när vandringen redan är slut väntar på nästa torn och står
// inte i resultatet, krönikans siffror eller på krönikans skärm.
func TestJoinAfterEnd(t *testing.T) {
	g, _, _ := startTwo(t)
	g.endGame()
	c := testClient()
	g.Join(c, "Sen")
	p := c.player
	if p == nil || p.Status != statusNew {
		t.Fatalf("status %v", p)
	}
	for _, r := range g.results() {
		if r.Name == "Sen" {
			t.Fatal("den sena står i resultatet")
		}
	}
	for _, l := range g.chronicleLines() {
		if strings.Contains(l.Text, "Sen") {
			t.Fatalf("den sena står i krönikan: %q", l.Text)
		}
	}
	if st := g.stateFor(p); st.Journey != nil {
		t.Fatal("den sena ska inte ha någon egen vandring")
	}
	if l := logs(c); !contains(l, g.data.Text("sen_ankomst", "namn", "Sen")) {
		t.Fatalf("saknar texten om sen ankomst: %v", l)
	}
	// Med i nästa torn.
	for _, q := range g.players {
		if q != nil && q.client != nil {
			g.Ready(q.client)
		}
	}
	if g.phase == phaseEnd || p.Status != statusIn {
		t.Fatalf("fas %s, status %s", g.phase, p.Status)
	}
}

// Ett steg mot en vägg avvisas direkt och kostar inte rundan, utom i mörker.
func TestWallRejectedAtOnce(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	room := &g.rooms[p.Y][p.X]
	wall := Dir(-1)
	for d := North; d <= West; d++ {
		if !room.Doors[d] && !(room.HasExit && room.Exit == d) {
			wall = d
			break
		}
	}
	if wall < 0 {
		t.Skip("rummet har ingen vägg")
	}
	code := "NESW"[wall : wall+1]
	round := g.round
	g.Act(a, code, 0)
	if p.Action != "" || g.round != round {
		t.Fatalf("draget mot väggen godtogs: %q, runda %d", p.Action, g.round)
	}
	p.Dark = 2
	g.Act(a, code, 0)
	if p.Action != code {
		t.Fatalf("i mörkret ska man kunna pröva väggen, drag %q", p.Action)
	}
}

// Grannrum har aldrig samma rumstyp, typerna används jämnt, och rummen
// bredvid en port har ingen fälla.
func TestRoomTypesAndSafePorts(t *testing.T) {
	g, _, _ := startTwo(t)
	g.data.Regler.Handelsechans = 1
	for y := range g.size {
		for x := range g.size {
			g.reveal(x, y)
		}
	}
	count := map[string]int{}
	for y := range g.size {
		for x := range g.size {
			r := &g.rooms[y][x]
			if r.Kind != "rum" {
				continue
			}
			count[r.TypeID]++
			if x+1 < g.size && g.rooms[y][x+1].TypeID == r.TypeID {
				t.Fatalf("två %s bredvid varandra vid %d,%d", r.TypeID, x, y)
			}
			if y+1 < g.size && g.rooms[y+1][x].TypeID == r.TypeID {
				t.Fatalf("två %s ovanför varandra vid %d,%d", r.TypeID, x, y)
			}
			if g.nextToPort(x, y) && r.Event != "" {
				t.Fatalf("fälla %s bredvid porten vid %d,%d", r.Event, x, y)
			}
		}
	}
	lo, hi := 1000, 0
	for _, n := range count {
		lo, hi = min(lo, n), max(hi, n)
	}
	if len(count) != len(g.data.Rumstyper) || hi-lo > 2 {
		t.Fatalf("ojämnt fördelade rumstyper: %v", count)
	}
}

// Sök i kammaren tar en näve, så att spelarna bara behöver en tangent.
func TestSearchTakesInTheChamber(t *testing.T) {
	g, a, b := startTwo(t)
	pa := a.player
	g.reveal(g.center(), g.center())
	pa.X, pa.Y = g.center(), g.center()
	g.deck = []string{"sover"}
	g.Act(a, actSearch, 0)
	if pa.Action != actTake {
		t.Fatalf("sök i kammaren blev %q, inte en näve", pa.Action)
	}
	g.Act(b, actWait, 0)
	if pa.Treasure <= 0 {
		t.Fatal("ingen skatt togs")
	}
}
