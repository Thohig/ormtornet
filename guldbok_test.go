package main

import (
	"strings"
	"testing"
)

// Guldboken: varje ändring blir en rad med summan efteråt, flera nävar i
// rad blir en mening, en förlust säger vad som var kvar, och porten och
// placeringen står sist. Ingen text är trasig, och den som aldrig fick
// något guld får en bok som säger det.
func TestGoldBook(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	note := func(p *Player, round int, kind string, n int, vars ...string) {
		g.round = round
		p.Treasure += n
		g.goldNote(p, kind, n, vars...)
	}
	note(pa, 2, "sok", 8, "rum", "ett kallt rum")
	note(pa, 5, "kammaren", 20, "kort", "Ormen sover")
	note(pa, 6, "kammaren", 15, "kort", "Ormen sover")
	note(pa, 7, "kammaren", 30, "kort", "Ormen rör sig")
	note(pa, 8, "ormen", -10)
	note(pb, 3, "sok", 5, "rum", "en trång gång")
	note(pa, 9, "slag", 5, "namn", pb.Name)
	note(pb, 9, "slagen", -5, "namn", pa.Name)
	note(pa, 12, "foremal", 4)
	pa.Status = statusOut
	pb.Status = statusDead
	g.goldNote(pa, "nothing", 0) // noll skrivs inte ner
	g.buildStory()

	if len(g.books) != 2 || g.books[0].Name != pa.Name {
		t.Fatalf("böckerna: %+v", g.books)
	}
	ba := g.books[0]
	if len(ba.Rows) != 7 {
		t.Fatalf("raderna: %+v", ba.Rows)
	}
	if r := ba.Rows[4]; r.Amount != -10 || r.Total != 63 || r.Round != 8 {
		t.Fatalf("förlusten: %+v", r)
	}
	if last := ba.Rows[6]; last.Total != pa.Treasure || pa.Treasure != 72 {
		t.Fatalf("sista raden: %+v, bär %d", last, pa.Treasure)
	}
	all := strings.Join(ba.Story, "\n")
	for _, want := range []string{pa.Name, "72", "ett kallt rum", g.data.Text("tal_3"), "30", "63", pb.Name, "1 av 2"} {
		if !strings.Contains(all, want) {
			t.Fatalf("berättelsen saknar %q:\n%s", want, all)
		}
	}
	for _, bk := range g.books {
		for _, s := range append(bk.Story, func() (out []string) {
			for _, r := range bk.Rows {
				out = append(out, r.Text)
			}
			return
		}()...) {
			if strings.Contains(s, "{") || strings.Contains(s, "[") || s == "" {
				t.Fatalf("trasig text: %q", s)
			}
		}
	}
	bb := g.books[1]
	if len(bb.Rows) != 2 || bb.Rows[1].Amount != -5 || bb.Rows[1].Total != 0 {
		t.Fatalf("den slagnas rader: %+v", bb.Rows)
	}
	if !strings.Contains(strings.Join(bb.Story, " "), pa.Name) {
		t.Fatalf("den slagnas berättelse: %q", bb.Story)
	}
}
