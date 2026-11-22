package main

import "testing"

func TestCarryThreeAndLeaveInRoom(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	putInRoom(g, pa, 4, 3)
	for range 3 {
		if !g.give(pa, "veke") {
			t.Fatal("väskan tog inte tre föremål")
		}
	}
	if g.give(pa, "rep") {
		t.Fatal("väskan tog ett fjärde")
	}
	// Ett föremål som inte får plats blir kvar i rummet.
	g.data.Foremal.Lista = []ItemType{*g.data.itemType("rep")}
	g.data.Foremal.Lista[0].Vikt = 1
	g.findItem(pa)
	if r := &g.rooms[3][4]; len(r.items) != 1 || r.items[0] != "rep" {
		t.Fatalf("rummet har %v", r.items)
	}
	// Bo tar det när han kommer in.
	putInRoom(g, pb, 4, 3)
	g.pickUp(pb)
	if pb.has("rep") == nil || len(g.rooms[3][4].items) != 0 {
		t.Fatal("Bo tog inte repet")
	}
	if st := g.stateFor(pb); len(st.Items) != 1 {
		t.Fatal("klienten ser inte föremålet")
	}
}

func TestRopeAmuletWick(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	putInRoom(g, p, 4, 3)

	g.give(p, "rep")
	g.trigger(p, g.data.trap("rasande_golv"), 0)
	if p.Cellar != 0 || p.has("rep") != nil {
		t.Fatalf("källaren %d, repet kvar %v", p.Cellar, p.has("rep") != nil)
	}

	g.give(p, "amulett")
	hp := p.HP
	g.trigger(p, g.data.trap("pilar"), 0)
	if p.HP != hp || p.has("amulett") != nil {
		t.Fatal("amuletten tog inte fällan")
	}
	if g.rooms[3][4].Mark == "" {
		t.Fatal("fällans tecken syns inte")
	}

	g.give(p, "veke")
	g.trigger(p, g.data.trap("morker"), 0)
	if p.Dark != 0 {
		t.Fatal("veken räckte inte")
	}
}

func TestShieldBreaksAfterThree(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.HP = 50
	g.give(p, "skold")
	for i := range 3 {
		if got := g.shieldDamage(p, 2); got != 1 {
			t.Fatalf("träff %d: %d skada", i+1, got)
		}
	}
	if p.has("skold") != nil {
		t.Fatal("skölden höll för mer än tre träffar")
	}
	if got := g.shieldDamage(p, 2); got != 2 {
		t.Fatal("skada utan sköld")
	}
}

func TestDoubleBottomBag(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	g.give(pb, "vaska")
	pb.Treasure = 90
	pa.Action, pb.Action = actAttack, actWait
	putInRoom(g, pa, 4, 3)
	putInRoom(g, pb, 4, 3)
	g.hit(pa, pb, 30, false)
	if pb.Treasure != 75 || pa.Treasure != 15 {
		t.Fatalf("Bo har %d, Alva %d", pb.Treasure, pa.Treasure)
	}
	if pb.has("vaska") == nil {
		t.Fatal("väskan försvann")
	}
}

func TestCompass(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	putInRoom(g, p, 2, 3)
	if g.compass(p) != "" {
		t.Fatal("kompass utan kompass")
	}
	g.give(p, "kompass")
	// Närmaste port från 2,3 i ett torn på 9 är nordväst (0,0).
	if want := g.data.Text("kompass_pekar", "horn", "nordväst"); g.compass(p) != want {
		t.Fatalf("kompassen: %q", g.compass(p))
	}
	if st := g.stateFor(p); st.Compass == "" {
		t.Fatal("klienten ser inte kompassen")
	}
}

func TestItemsResetWithNewTower(t *testing.T) {
	g, a, _ := startTwo(t)
	g.give(a.player, "rep")
	g.newTower()
	if len(a.player.items) != 0 {
		t.Fatal("föremålen följde med till nästa torn")
	}
}

// openRoom ställer p i ett rum med dörrar åt alla håll och tomma grannrum.
func openRoom(g *Game, p *Player, x, y int) {
	r := putInRoom(g, p, x, y)
	r.Doors = [4]bool{true, true, true, true}
	for d := North; d <= West; d++ {
		n := g.reveal(x+dirDX[d], y+dirDY[d])
		n.Kind, n.Event, n.Walked = "rum", "", true
		n.Doors[d.Opp()] = true
	}
}

func TestSmokeFleesWithoutBlow(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	openRoom(g, p, 4, 3)
	typ := &g.data.Monster.Typer[0]
	typ.Vanor = map[string]int{actAttack: 1}
	g.monsters = append(g.monsters, &Monster{ID: 60, Type: typ.ID, X: 4, Y: 3, HP: 9})
	g.give(p, "rokpulver")
	hp := p.HP
	g.actPlayer(p, "E", 0)
	g.playTurn(p)
	if p.HP != hp || p.X != 5 || p.has("rokpulver") != nil {
		t.Fatalf("liv %d av %d, står på %d", p.HP, hp, p.X)
	}
}

func TestItemActionsAreGone(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	g.give(p, "krita")
	g.actPlayer(p, "f1E", 0)
	if p.Action != "" {
		t.Fatalf("draget med föremål togs emot: %q", p.Action)
	}
}

func TestMirrorShowsANewRoom(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	openRoom(g, p, 4, 3)
	p.known = newGrid[bool](g.size)
	for d := North; d <= West; d++ {
		p.known[3+dirDY[d]][4+dirDX[d]] = d != East
	}
	g.rooms[3][5].Event = "pilar"
	g.give(p, "spegel")
	logs(a)
	g.mirror(p)
	if !p.known[3][5] || !containsPart(logs(a), g.data.Text("spegel_falla", "falla", g.data.trap("pilar").Kort)) {
		t.Fatal("spegeln visade inte rummet österut med fällan")
	}
	if p.has("spegel").uses != 2 {
		t.Fatal("spegeln förbrukades inte")
	}
	g.mirror(p)
	if p.has("spegel").uses != 2 {
		t.Fatal("spegeln verkade fast alla grannrum var kända")
	}
}

func TestKnifeChalkHourglass(t *testing.T) {
	g, a, b := startTwo(t)
	p, o := a.player, b.player

	openRoom(g, p, 4, 3)
	g.monsters = append(g.monsters, &Monster{ID: 61, Type: g.data.Monster.Typer[0].ID, X: 4, Y: 3, HP: 1})
	g.give(p, "kastkniv")
	gold := p.Treasure
	g.announceMonster(p)
	if g.monsterAt(4, 3) != nil || p.Treasure <= gold || p.has("kastkniv") != nil {
		t.Fatal("kniven fällde inte varelsen")
	}

	openRoom(g, p, 4, 3)
	g.give(p, "krita")
	g.move(p, North, 0)
	if r := g.rooms[3][4]; !r.hasChalk || r.chalk != North || p.has("krita").uses != 4 {
		t.Fatal("ingen pil när vandraren gick norrut")
	}
	openRoom(g, o, 4, 4)
	logs(b)
	g.move(o, North, 0)
	if !containsPart(logs(b), "krita") {
		t.Fatal("Bo såg inte pilen")
	}
	p.items = nil

	g.give(p, "timglas")
	g.sun = 5
	g.hourglass()
	if g.sun != 5 || p.has("timglas") == nil {
		t.Fatal("timglaset vändes för tidigt")
	}
	g.sun = 3
	g.hourglass()
	if g.sun != 4 || p.has("timglas") != nil {
		t.Fatal("timglaset stoppade inte solen")
	}
}

func TestFluteRedraws(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	g.give(p, "flojt")
	wake := -1
	for i, id := range g.deck {
		if c := g.data.card(id); c != nil && c.Vaknar {
			wake = i
		}
	}
	g.flute(p, wake)
	if p.has("flojt") != nil {
		t.Fatal("flöjten användes inte")
	}
}

func TestCursedItems(t *testing.T) {
	g, a, b := startTwo(t)
	p, o := a.player, b.player

	// Pärlan: varje näve lägger ett vakna-kort i leken. Värd 40 guld ute.
	g.give(p, "svart_parla")
	p.X, p.Y = g.center(), g.center()
	g.reveal(p.X, p.Y)
	for i, id := range g.deck { // dra ett sovkort
		if c := g.data.card(id); c != nil && !c.Vaknar {
			g.deck[0], g.deck[i] = g.deck[i], g.deck[0]
			break
		}
	}
	g.deck = g.deck[:1]
	g.take(p)
	if g.wakeCardsLeft() != 1 {
		t.Fatalf("%d vakna-kort efter en näve med pärlan", g.wakeCardsLeft())
	}
	if g.itemGold(p) != 40 {
		t.Fatalf("pärlan är värd %d", g.itemGold(p))
	}

	// Kronan: syns från grannrummet.
	p.items = nil
	g.give(p, "rostkrona")
	putInRoom(g, p, 4, 3)
	putInRoom(g, o, 5, 3)
	if v := viewOf(g.stateFor(o), p); v.Hidden {
		t.Fatal("kronan syns inte från grannrummet")
	}
}

func TestPilgrimLeavesCursed(t *testing.T) {
	g, p, _ := withClass(t, "pilgrim")
	putInRoom(g, p, 4, 3)
	g.rooms[3][4].items = []string{"mynt_hal"}
	g.pickUp(p)
	if p.has("mynt_hal") != nil || len(g.rooms[3][4].items) != 1 {
		t.Fatal("pilgrimen tog det förbannade myntet")
	}
}

func TestLetterAndDeathDrop(t *testing.T) {
	g, a, b := startTwo(t)
	p, o := a.player, b.player
	putInRoom(g, p, 4, 3)
	logs(a)
	g.give(p, "brev")
	if len(p.items) != 0 || len(logs(a)) < 2 {
		t.Fatal("brevet lästes inte, eller lades i väskan")
	}
	// Den som dör tappar sina föremål i rummet, och nästa tar dem.
	g.give(p, "rep")
	g.give(p, "veke")
	g.damage(p, 99, false)
	if len(g.rooms[3][4].items) != 2 {
		t.Fatalf("rummet har %v", g.rooms[3][4].items)
	}
	putInRoom(g, o, 4, 3)
	g.pickUp(o)
	if o.has("rep") == nil || o.has("veke") == nil {
		t.Fatal("Bo tog inte det Alva tappade")
	}
}

func TestMapMakerChalk(t *testing.T) {
	g, p, _ := withClass(t, "kartritare")
	g.give(p, "krita")
	if p.has("krita").uses != 8 {
		t.Fatalf("kritan har %d streck", p.has("krita").uses)
	}
}
