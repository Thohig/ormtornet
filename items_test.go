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
	g.hit(pa, pb, 30)
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
