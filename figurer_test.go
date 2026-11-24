package main

import "testing"

// Kartan får veta vad vandraren och varelsen gjorde under rundan, och den
// varelse som föll står i fallna rundan den föll, men inte rundan efter.
func TestMovesAndFallenForTheMap(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	typ := g.data.Monster.Typer[0]
	m := &Monster{ID: 77, Type: typ.ID, X: p.X, Y: p.Y, HP: 1}
	g.monsters = append(g.monsters, m)
	for i := 0; i < 40 && m.HP > 0; i++ {
		// Varelsen står kvar i rummet och tål ett slag.
		m.X, m.Y, m.HP = p.X, p.Y, 1
		p.HP, p.Action = 50, actAttack
		g.resolve()
		if p.Status != statusIn {
			t.Fatal("vandraren är inte kvar i tornet")
		}
	}
	if m.HP > 0 {
		t.Fatal("varelsen föll aldrig")
	}
	st := g.stateFor(p)
	if me := viewOf(st, p); me.Did != actAttack {
		t.Fatalf("gjorde = %q", me.Did)
	}
	if len(st.Fallen) != 1 || st.Fallen[0].ID != 77 || !st.Fallen[0].Hit {
		t.Fatalf("fallna = %+v", st.Fallen)
	}
	for _, v := range st.Monsters {
		if v.ID == 77 {
			t.Fatal("den fallna står bland varelserna")
		}
	}
	p.Action = actWait
	g.resolve()
	st = g.stateFor(p)
	if len(st.Fallen) != 0 || viewOf(st, p).Did != "" {
		t.Fatalf("rundan efter: fallna %+v, gjorde %q", st.Fallen, viewOf(st, p).Did)
	}
}
