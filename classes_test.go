package main

import (
	"slices"
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
	a.player.offer = []string{id}
	if !g.chooseClass(a.player, id) {
		t.Fatalf("kunde inte välja %s", id)
	}
	g.Ready(a)
	g.Ready(b)
	if g.phase != phasePlay {
		t.Fatal("spelet startade inte")
	}
	return g, a.player, b.player
}

func TestClassOfferAndChoose(t *testing.T) {
	g, a, b := classGame(t)
	pa, pb := a.player, b.player
	if len(pa.offer) != 2 || pa.offer[0] == pa.offer[1] {
		t.Fatalf("Alva fick inte två olika klasser: %v", pa.offer)
	}
	if st := g.stateFor(pa); len(st.Offer) != 2 {
		t.Fatal("klienten ser inte valet")
	}
	id := pa.offer[0]
	g.ChooseClass(a, id)
	if pa.classID != id || pa.HP != g.maxHP(pa) {
		t.Fatalf("valet gick inte igenom: %q, %d liv", pa.classID, pa.HP)
	}
	// Bo kan inte välja samma klass.
	pb.offer = []string{id}
	if g.chooseClass(pb, id) || pb.classID == id {
		t.Fatal("två vandrare med samma klass")
	}
	if slices.Contains(pb.offer, id) {
		t.Fatal("Bo fick samma klass att välja på igen")
	}
	// Den som inte väljer får en klass när vandringen börjar.
	g.Ready(a)
	g.Ready(b)
	if pb.classID == "" {
		t.Fatal("Bo fick ingen klass")
	}
	if v := viewOf(g.stateFor(pa), pb); v.Class == "" {
		t.Fatal("klassen syns inte i listan")
	}
}

func TestClassesOff(t *testing.T) {
	g := testGame(t)
	a := testClient()
	g.Join(a, "Alva")
	if len(a.player.offer) != 0 || g.classOf(a.player) != nil {
		t.Fatal("klass fast klasserna är av")
	}
}

func TestBotPicksPreferredClass(t *testing.T) {
	g := testGame(t)
	g.cfg.Classes = true
	g.AddBots(1)
	var bot *Player
	for _, p := range g.players {
		if p != nil {
			bot = p
		}
	}
	if bot == nil || bot.classID == "" {
		t.Fatal("AI-spelaren valde ingen klass")
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
	g.reveal(2, 2)
	g.meetIfNew(pb)
	if !pb.known[2][2] {
		t.Fatal("Bo fick inte kartan")
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
