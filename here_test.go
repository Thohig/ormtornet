package main

import (
	"encoding/json"
	"testing"
)

// rooms ger rumsblocken som klienten har fått.
func rooms(c *Client) []roomMsg {
	var out []roomMsg
	for {
		select {
		case b := <-c.send:
			var m roomMsg
			if json.Unmarshal(b, &m) == nil && m.T == "rum" {
				out = append(out, m)
			}
		default:
			return out
		}
	}
}

func TestRoomBlockAtStart(t *testing.T) {
	g := testGame(t)
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	rooms(a)
	g.Ready(a)
	g.Ready(b)
	got := rooms(a)
	if len(got) != 1 || got[0].Text == "" || got[0].Har == "" || got[0].Kvar {
		t.Fatalf("vid starten väntades ett rumsblock med beskrivning, fick %+v", got)
	}
}

func TestRoomBlockTellsWhatIsHere(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	putInRoom(g, pa, 4, 4)
	putInRoom(g, pb, 4, 4)
	rooms(a)

	g.tellRoom(pa)
	got := rooms(a)
	want := g.data.Text("har_vandrare", "namn", "Bo")
	if len(got) != 1 || got[0].Kvar || got[0].Har != want {
		t.Fatalf("väntade ett block med %q, fick %+v", want, got)
	}

	g.tellRoom(pa)
	if got := rooms(a); len(got) != 0 {
		t.Fatalf("inget har ändrats, men fick %+v", got)
	}

	pb.X = 5
	g.tellRoom(pa)
	got = rooms(a)
	want = g.data.Text("har_vandrare_borta", "namn", "Bo")
	if len(got) != 1 || !got[0].Kvar || got[0].Text != "" || got[0].Har != want {
		t.Fatalf("väntade bara %q, fick %+v", want, got)
	}

	// Varelsen beskrivs med sin egen text, och samma text så länge den
	// står kvar.
	typ := g.data.Monster.Typer[0]
	g.monsters = append(g.monsters, &Monster{ID: 1, Type: typ.ID, X: 4, Y: 4, HP: 1})
	g.tellRoom(pa)
	got = rooms(a)
	want = typ.DykerUpp[1%len(typ.DykerUpp)]
	if len(got) != 1 || got[0].Har != want {
		t.Fatalf("väntade %q, fick %+v", want, got)
	}
	g.tellRoom(pa)
	if got := rooms(a); len(got) != 0 {
		t.Fatalf("varelsen står kvar, men fick %+v", got)
	}
}

// Ett tomt rum säger ingenting, och en fälla som har slagit till syns som
// ett spår utan sitt namn.
func TestRoomBlockShowsNotTells(t *testing.T) {
	g, a, _ := startTwo(t)
	pa := a.player
	r := putInRoom(g, pa, 3, 3)
	pa.told = told{}
	rooms(a)
	g.tellRoom(pa)
	got := rooms(a)
	if len(got) != 1 || got[0].Har != "" {
		t.Fatalf("ett tomt rum ska inte säga något, fick %+v", got)
	}
	trap := g.data.trap("roterande")
	r.Mark = trap.Tecken
	g.tellRoom(pa)
	got = rooms(a)
	if len(got) != 1 || got[0].Har != trap.SparEfter[0] {
		t.Fatalf("väntade spåret %q, fick %+v", trap.SparEfter[0], got)
	}
}

func TestRoomBlockWaitsInTheDark(t *testing.T) {
	g, a, _ := startTwo(t)
	pa := a.player
	putInRoom(g, pa, 4, 4)
	pa.Dark = 2
	rooms(a)
	g.tellRoom(pa)
	if got := rooms(a); len(got) != 0 {
		t.Fatalf("i mörkret ska inget block komma, fick %+v", got)
	}
}

func TestJourneyRoundHasRoom(t *testing.T) {
	g, a, b := startTwo(t)
	pa := a.player
	putInRoom(g, pa, 4, 4)
	putInRoom(g, b.player, 2, 2)
	step(g, a, b, actWait, actWait)
	if len(pa.journal) == 0 {
		t.Fatal("ingen runda i vandringen")
	}
	last := pa.journal[len(pa.journal)-1]
	if last.Room == "" {
		t.Fatalf("rundan i vandringen saknar rummet: %+v", last)
	}
}

// Det vandraren hör under rundan står sist i rumsblocket och inte som egna
// rader; står hen kvar kommer ett block bara med ljuden.
func TestRoomBlockCarriesWhatIsHeard(t *testing.T) {
	g, a, _ := startTwo(t)
	pa := a.player
	putInRoom(g, pa, 3, 3)
	pa.told = told{}
	g.tellRoom(pa)
	rooms(a)
	g.hearing = true
	g.hearTo(pa, "Stridsljud från norr.", "dim")
	g.hearTo(pa, "Stridsljud från norr.", "dim")
	g.hearing = false
	if l := logs(a); len(l) != 0 {
		t.Fatalf("ljudet kom som egna rader: %v", l)
	}
	g.tellRoom(pa)
	got := rooms(a)
	if len(got) != 1 || !got[0].Kvar || got[0].Hor != "Stridsljud från norr." {
		t.Fatalf("väntade ett block med ljudet en gång, fick %+v", got)
	}
}
