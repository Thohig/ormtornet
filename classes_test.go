package main

import (
	"testing"
)

// classGame ger ett torn med klasser i lobbyn och två anslutna vandrare.
func classGame(t *testing.T) (*Game, *Client, *Client) {
	t.Helper()
	g := testGame(t)
	g.cfg.Classes = true
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	return g, a, b
}

// withClass startar en vandring där Alva har klassen id.
func withClass(t *testing.T, id string) (*Game, *Player, *Player) {
	t.Helper()
	g, a, b := classGame(t)
	g.Ready(a)
	g.Ready(b)
	if g.phase != phaseChoose {
		t.Fatalf("klassvalet började inte, fasen är %q", g.phase)
	}
	g.ChooseClass(a, id)
	if a.player.classID != id {
		t.Fatalf("kunde inte välja %s", id)
	}
	g.startGame()
	if g.phase != phasePlay {
		t.Fatal("spelet startade inte")
	}
	return g, a.player, b.player
}

func TestClassChoosing(t *testing.T) {
	g, a, b := classGame(t)
	pa, pb := a.player, b.player
	if st := g.stateFor(pa); len(st.Offer) != 0 {
		t.Fatal("korten syns redan i lobbyn")
	}
	g.ChooseClass(a, "tjuv")
	if pa.classID != "" {
		t.Fatal("klass vald i lobbyn")
	}
	g.Ready(a)
	g.Ready(b)
	if g.phase != phaseChoose {
		t.Fatalf("fasen är %q, inte klassval", g.phase)
	}
	st := g.stateFor(pa)
	if len(st.Offer) != len(g.data.Klasser) || st.DeadlineMs <= 0 {
		t.Fatalf("%d kort, %d ms kvar", len(st.Offer), st.DeadlineMs)
	}
	id := g.data.Klasser[0].ID
	g.ChooseClass(a, id)
	if pa.classID != id || pa.HP != g.maxHP(pa) {
		t.Fatalf("valet gick inte igenom: %q, %d liv", pa.classID, pa.HP)
	}
	// Alva kan inte byta, och Bo kan inte ta samma kort.
	g.ChooseClass(a, g.data.Klasser[1].ID)
	g.ChooseClass(b, id)
	if pa.classID != id || pb.classID != "" || g.phase != phaseChoose {
		t.Fatalf("Alva %q, Bo %q, fasen %q", pa.classID, pb.classID, g.phase)
	}
	if v := g.stateFor(pb).Offer[0]; v.Tagen != "Alva" || v.Din {
		t.Fatalf("kortet visar inte vem som tog det: %+v", v)
	}
	// När alla har valt börjar vandringen.
	g.ChooseClass(b, g.data.Klasser[2].ID)
	if g.phase != phasePlay || pb.classID != g.data.Klasser[2].ID {
		t.Fatalf("fasen %q, Bo %q", g.phase, pb.classID)
	}
	if v := viewOf(g.stateFor(pa), pb); v.Class == "" {
		t.Fatal("klassen syns inte i listan")
	}
}

func TestClassChoosingTimesOut(t *testing.T) {
	g, a, b := classGame(t)
	g.Ready(a)
	g.Ready(b)
	g.ChooseClass(a, "tjuv")
	// Tiden är slut: Bo får ett av de lediga korten.
	g.startGame()
	if b.player.classID == "" || b.player.classID == "tjuv" {
		t.Fatalf("Bo fick %q", b.player.classID)
	}
}

func TestClassesOff(t *testing.T) {
	g := testGame(t)
	a := testClient()
	g.Join(a, "Alva")
	g.begin()
	if g.phase != phasePlay || g.classOf(a.player) != nil {
		t.Fatalf("fasen %q, klass %v fast klasserna är av", g.phase, g.classOf(a.player))
	}
}

// AI-spelaren väntar tills människan har valt och tar sedan den lediga
// klass den helst vill ha.
func TestBotPicksPreferredClass(t *testing.T) {
	g := testGame(t)
	g.cfg.Classes = true
	g.AddBots(1)
	a := testClient()
	g.Join(a, "Alva")
	g.Ready(a)
	var bot *Player
	for _, p := range g.players {
		if p != nil && p.bot != nil {
			bot = p
		}
	}
	if bot == nil {
		t.Fatal("ingen AI-spelare")
	}
	if bot.classID != "" {
		t.Fatalf("AI-spelaren valde %q före människan", bot.classID)
	}
	want := bot.bot.pers.Klasser
	if len(want) < 2 {
		t.Skip("personligheten har för få önskade klasser")
	}
	// Alva tar AI-spelarens favorit, så den får nöja sig med nästa.
	g.ChooseClass(a, want[0])
	if a.player.classID != want[0] {
		t.Fatalf("Alva fick %q, inte %q", a.player.classID, want[0])
	}
	if bot.classID != want[1] {
		t.Fatalf("AI-spelaren valde %q, inte %q", bot.classID, want[1])
	}
	if g.phase != phasePlay {
		t.Fatalf("fasen är %q efter att alla valt", g.phase)
	}
}

// När tiden tar slut får människan sitt kort före AI-spelarna.
func TestBotWaitsWhenTimeRunsOut(t *testing.T) {
	g := testGame(t)
	g.cfg.Classes = true
	g.AddBots(2)
	a := testClient()
	g.Join(a, "Alva")
	g.Ready(a)
	g.startGame()
	for _, p := range g.players {
		if p != nil && p.classID == "" {
			t.Fatalf("%s fick ingen klass", p.Name)
		}
	}
}

func TestLegoknektShield(t *testing.T) {
	g, pa, _ := withClass(t, "legoknekt")
	if pa.HP != g.data.Regler.Liv+1 {
		t.Fatalf("legoknekten har %d liv", pa.HP)
	}
	putInRoom(g, pa, 4, 3)
	typ := &g.data.Monster.Typer[0]
	typ.Vanor = map[string]int{actAttack: 1}
	m := &Monster{ID: 50, Type: typ.ID, X: 4, Y: 3, HP: 10}
	g.monsters = append(g.monsters, m)
	hp := pa.HP
	g.duel(pa, m, actAttack)
	if pa.HP != hp || m.HP != 10-g.data.Monster.Spelarskada {
		t.Fatalf("liv %d, varelsen %d", pa.HP, m.HP)
	}
}

func TestOrtkvinna(t *testing.T) {
	g, pa, _ := withClass(t, "ortkvinna")
	if pa.Potions != 2 {
		t.Fatalf("%d drycker", pa.Potions)
	}
	putInRoom(g, pa, 4, 3)
	g.trigger(pa, g.data.trap("dimma"), 0)
	if pa.Poison != 0 {
		t.Fatal("giftet bet")
	}
	hp := pa.HP
	g.drink(pa)
	if pa.HP != min(g.maxHP(pa), hp+3) {
		t.Fatalf("drycken gav %d liv", pa.HP-hp)
	}
}

func TestThiefLootAndPilgrimPeace(t *testing.T) {
	g, pa, pb := withClass(t, "tjuv")
	pb.Treasure = 100
	if got := g.loot(pa, pb); got != 50 {
		t.Fatalf("tjuven tog %d", got)
	}
	if got := g.loot(pb, pa); got != 0 {
		t.Fatalf("utan skatt: %d", got)
	}

	g2, pilgrim, other := withClass(t, "pilgrim")
	other.Treasure = 90
	putInRoom(g2, pilgrim, 4, 3)
	putInRoom(g2, other, 4, 3)
	pilgrim.Action = actAttack
	hp := other.HP
	g2.ambush(pilgrim, other)
	if other.HP != hp || other.Treasure != 90 {
		t.Fatal("pilgrimen slog")
	}
}

func TestKartritareSeesAll(t *testing.T) {
	g, pa, pb := withClass(t, "kartritare")
	g.cfg.Sight = sightOwn
	putInRoom(g, pb, 6, 6)
	if !hasRoom(g.stateFor(pa), 6, 6) {
		t.Fatal("kartritaren ser inte rummet")
	}
	// Bo får kartan när de möts.
	putInRoom(g, pa, 6, 6)
	pb.met = ""
	g.reveal(2, 2).Walked = true
	g.meetIfNew(pb)
	if !pb.known[2][2] {
		t.Fatal("Bo fick inte kartan")
	}
}

// Kartritaren ser genom dörrarna in i grannrummen, men de räknas inte som
// hittade: ingen annan ser dem, och fällan där väntar på den som går in.
func TestKartritarePeeks(t *testing.T) {
	for _, sight := range []string{sightOwn, sightAll} {
		g, pa, pb := withClass(t, "kartritare")
		g.cfg.Sight = sight
		r := putInRoom(g, pa, 4, 4)
		r.Doors = [4]bool{false, true, false, false}
		g.noteRoom(pa)
		if !hasRoom(g.stateFor(pa), 5, 4) {
			t.Fatalf("%s: kartritaren ser inte in i rummet åt öster", sight)
		}
		if hasRoom(g.stateFor(pa), 3, 4) {
			t.Fatalf("%s: kartritaren ser genom väggen", sight)
		}
		if hasRoom(g.stateFor(pb), 5, 4) {
			t.Fatalf("%s: Bo ser rummet kartritaren bara tittade in i", sight)
		}
		n := &g.rooms[4][5]
		n.Kind, n.Event = "rum", g.data.Fallor[0].ID
		n.Doors[West] = true
		putInRoom(g, pb, 4, 4)
		pb.HP = 50
		g.move(pb, East, 0)
		if n.Event != "" {
			t.Fatalf("%s: fällan gick inte av för Bo", sight)
		}
		if !hasRoom(g.stateFor(pa), 5, 4) {
			t.Fatalf("%s: rummet försvann från kartritarens karta", sight)
		}
	}
}

func TestLanternSeenFromNextRoom(t *testing.T) {
	g, pa, pb := withClass(t, "lyktbarare")
	g.cfg.Sight = sightOwn
	putInRoom(g, pa, 4, 3)
	putInRoom(g, pb, 5, 3)
	if v := viewOf(g.stateFor(pb), pa); v.Hidden {
		t.Fatal("lyktbäraren syns inte från grannrummet")
	}
	if v := viewOf(g.stateFor(pa), pb); !v.Hidden {
		t.Fatal("lyktbäraren ser Bo i grannrummet")
	}
	putInRoom(g, pb, 7, 3)
	if v := viewOf(g.stateFor(pb), pa); !v.Hidden {
		t.Fatal("lyktbäraren syns långt bort")
	}
}

func TestOrmtjusareLessDamage(t *testing.T) {
	g, pa, _ := withClass(t, "ormtjusare")
	pa.X, pa.Y = g.center(), g.center()
	g.reveal(g.center(), g.center())
	hp := pa.HP
	g.wakeSerpent()
	if want := hp - max(0, g.data.Ormen.Skada-1); pa.HP != want {
		t.Fatalf("liv %d, ville ha %d", pa.HP, want)
	}
}

func TestSparfinnareNamesTrail(t *testing.T) {
	g, pa, pb := withClass(t, "sparfinnare")
	g.cfg.Sight = sightOwn
	putInRoom(g, pa, 4, 4)
	g.rooms[4][4].trail = trail{by: pb.ID, round: g.round, dir: West}
	b := pa.client
	logs(b)
	g.noticeTrail(pa)
	if !containsPart(logs(b), pb.Name) {
		t.Fatal("spårfinnaren ser inte vems spåren är")
	}
}

// Nytt torn leder tillbaka till lobbyn och sedan till ett nytt klassval.
func TestNewTowerOffersNewClasses(t *testing.T) {
	g, a, b := classGame(t)
	g.Ready(a)
	g.Ready(b)
	g.startGame()
	if g.phase != phasePlay || a.player.classID == "" {
		t.Fatal("vandringen startade inte med klasser")
	}
	g.level = 3
	g.endGame()
	g.Ready(a)
	g.Ready(b)
	if g.level != 1 {
		t.Fatalf("lobbyn har kvar nivå %d från förra vandringen", g.level)
	}
	if g.phase != phaseLobby {
		t.Fatalf("nytt torn ledde till %q, inte lobbyn", g.phase)
	}
	for _, c := range []*Client{a, b} {
		p := c.player
		if p.Ready || p.Status != statusIn || p.Treasure != 0 {
			t.Fatalf("%s står inte redo att börja om vid porten", p.Name)
		}
	}
	g.Ready(a)
	g.Ready(b)
	if g.phase != phaseChoose || a.player.classID != "" || b.player.classID != "" {
		t.Fatalf("fasen %q, klasser %q %q", g.phase, a.player.classID, b.player.classID)
	}
	g.ChooseClass(a, "tjuv")
	g.ChooseClass(b, "kartritare")
	if g.phase != phasePlay {
		t.Fatal("nästa vandring startade inte efter klassvalet")
	}
}

// Den som laddar om sidan under klassvalet behåller sin plats, och
// vandringen väntar på hen fast de andra redan har valt.
func TestReloadDuringClassChoosing(t *testing.T) {
	g, a, b := classGame(t)
	g.Ready(a)
	g.Ready(b)
	g.ChooseClass(a, "tjuv")
	pb := b.player
	g.Leave(b)
	if g.phase != phaseChoose {
		t.Fatalf("fasen är %q: vandringen började utan Bo", g.phase)
	}
	b2 := testClient()
	g.Join(b2, "Bo")
	if b2.player != pb {
		t.Fatal("Bo fick inte tillbaka sin plats")
	}
	g.ChooseClass(b2, "pilgrim")
	if g.phase != phasePlay || pb.classID != "pilgrim" {
		t.Fatalf("fasen är %q, Bo har %q", g.phase, pb.classID)
	}
}

// Är hen fortfarande borta när tiden är slut följer hen med som frånvarande
// och kan komma in igen under vandringen.
func TestAbsentWhenClassChoosingEnds(t *testing.T) {
	g, a, b := classGame(t)
	g.Ready(a)
	g.Ready(b)
	g.ChooseClass(a, "tjuv")
	pb := b.player
	g.Leave(b)
	g.startGame() // tiden är slut
	if g.players[pb.Slot] != pb || pb.classID == "" || pb.awayRound != 1 {
		t.Fatalf("Bo är inte med som frånvarande: klass %q, borta sedan runda %d", pb.classID, pb.awayRound)
	}
	b2 := testClient()
	g.Join(b2, "Bo")
	if b2.player != pb || pb.awayRound != 0 {
		t.Fatal("Bo kom inte tillbaka in i vandringen")
	}
}
