package main

import "testing"

// noWander stänger av varelsernas egen vandring.
func noWander(d *Data) {
	d.Monster.VandraChans = 0
	for i := range d.Monster.Typer {
		d.Monster.Typer[i].Vandrar = nil
	}
}

// wanderGame ger ett torn där varelser alltid vandrar och där ett rum
// (2, 1) bara har en dörr, österut mot ett rum som inte finns än.
func wanderGame(t *testing.T) (*Game, *Monster) {
	t.Helper()
	g := testGame(t)
	g.data.Monster.VandraChans = 1
	g.data.Monster.Chans = 1 // ett rum som en varelse bygger får ändå ingen egen
	r := g.reveal(2, 1)
	r.Kind, r.Event, r.Doors = "rum", "", [4]bool{}
	r.Doors[East] = true
	if g.rooms[1][3].Revealed {
		t.Fatal("rummet österut finns redan")
	}
	g.monsters = nil
	m := &Monster{ID: 7, Type: "krypare", X: 2, Y: 1, HP: 2}
	g.monsters = append(g.monsters, m)
	return g, m
}

func TestMonsterWandersIntoUnknownRoom(t *testing.T) {
	g, m := wanderGame(t)
	g.wanderMonsters()
	if m.X != 3 || m.Y != 1 {
		t.Fatalf("varelsen står på %d,%d", m.X, m.Y)
	}
	n := &g.rooms[1][3]
	if !n.Revealed || n.Walked || n.Kind != "rum" {
		t.Fatalf("rummet: byggt %v, hittat %v, %q", n.Revealed, n.Walked, n.Kind)
	}
	if len(g.monsters) != 1 {
		t.Fatalf("%d varelser, rummet fick en egen", len(g.monsters))
	}
	if !n.Doors[West] {
		t.Fatal("rummet har ingen dörr tillbaka")
	}
}

func TestMonsterStaysWithWanderer(t *testing.T) {
	g, m := wanderGame(t)
	c := testClient()
	g.Join(c, "Alva")
	putInRoom(g, c.player, 2, 1)
	g.rooms[1][2].Doors = [4]bool{}
	g.rooms[1][2].Doors[East] = true
	c.player.Status = statusIn
	g.wanderMonsters()
	if m.X != 2 {
		t.Fatal("varelsen lämnade striden")
	}
}

func TestMonsterDoesNotWanderToWanderer(t *testing.T) {
	g, m := wanderGame(t)
	c := testClient()
	g.Join(c, "Alva")
	putInRoom(g, c.player, 3, 1)
	c.player.Status = statusIn
	g.rooms[1][3].Doors[West] = true
	for i := 0; i < 20; i++ {
		g.wanderMonsters()
	}
	if m.X != 2 {
		t.Fatal("varelsen vandrade in till Alva")
	}
}

func TestMonsterNeverEntersPortOrChamber(t *testing.T) {
	g, m := wanderGame(t)
	// Varelsen står bredvid ett portrum och kammaren och får välja fritt.
	for i := 0; i < 200; i++ {
		g.wanderMonsters()
		if r := &g.rooms[m.Y][m.X]; r.Kind != "rum" {
			t.Fatalf("varelsen gick in i %q på %d,%d", r.Kind, m.X, m.Y)
		}
	}
}
