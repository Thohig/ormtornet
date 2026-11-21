package main

import "testing"

// På nivå 4 sluter sig portarna en i taget, med en varning innan. Den sista
// stängs aldrig, och genom en stängd port kommer ingen ut.
func TestPortsCloseOneByOne(t *testing.T) {
	g, a, _ := startTwo(t)
	g.data.Avslojande.Forandring.PortStangsVar = 2
	g.level = 3
	g.closePorts()
	if len(g.openPorts()) != 4 {
		t.Fatal("en port stängdes före nivå 4")
	}
	g.level = 4
	logs(a)
	g.closePorts() // varning
	if g.portNext < 0 || len(logs(a)) == 0 {
		t.Fatal("ingen varning innan porten stängs")
	}
	warned := g.portNext
	g.closePorts() // stängs
	if open := g.openPorts(); len(open) != 3 || g.portNext != -1 {
		t.Fatalf("efter första stängningen: öppna %v", open)
	}
	if c := g.corners()[warned]; !g.rooms[c.y][c.x].Shut {
		t.Fatal("en annan port stängdes än den som darrade")
	}
	for range 20 {
		g.closePorts()
	}
	if len(g.openPorts()) != 1 {
		t.Fatalf("%d portar öppna, ville ha en", len(g.openPorts()))
	}
	// Den som står i en stängd port kommer inte ut.
	c := g.corners()[warned]
	p := a.player
	putInRoom(g, p, c.x, c.y)
	g.move(p, c.exit, 0)
	if p.Status != statusIn {
		t.Fatal("ut genom en stängd port")
	}
	if st := g.stateFor(p); !roomShut(st, c.x, c.y) {
		t.Fatal("klienten ser inte att porten är stängd")
	}
}

func roomShut(st stateMsg, x, y int) bool {
	for _, r := range st.Rooms {
		if r.X == x && r.Y == y {
			return r.Shut
		}
	}
	return false
}

// När tornet förändras gäller nivåns vikter för fällorna i nya rum.
func TestTrapWeightsFollowLevel(t *testing.T) {
	g, _, _ := startTwo(t)
	g.data.Avslojande.Forandring.Fallvikter = map[string]map[string]int{"3": {"ormgrop": 1000000000}}
	g.level = 2
	other := 0
	for range 200 {
		if g.randomTrap().ID != "ormgrop" {
			other++
		}
	}
	if other < 50 {
		t.Fatalf("vikterna gällde före nivå 3 (%d andra fällor)", other)
	}
	g.level = 4 // nivå 4 har inga egna vikter här: nivå 3 gäller
	for range 200 {
		if id := g.randomTrap().ID; id != "ormgrop" {
			t.Fatalf("fällan %s trots vikten", id)
		}
	}
}
