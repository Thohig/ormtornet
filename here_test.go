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
	if len(got) != 1 || !got[0].Kvar || got[0].Text != "" || got[0].Har != g.data.Text("har_tomt") {
		t.Fatalf("väntade ett kort block om att rummet är tomt, fick %+v", got)
	}

	g.monsters = append(g.monsters, &Monster{ID: 1, Type: g.data.Monster.Typer[0].ID, X: 4, Y: 4, HP: 1})
	g.tellRoom(pa)
	got = rooms(a)
	want = g.data.Text("har_varelse", "varelse", g.data.Monster.Typer[0].Namn)
	if len(got) != 1 || got[0].Har != want {
		t.Fatalf("väntade %q, fick %+v", want, got)
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
