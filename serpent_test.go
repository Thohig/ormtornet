package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// messagesOf plockar ut meddelandena av sorten typ ur det klienten har fått.
func messagesOf[M any](c *Client, typ string) []M {
	var out []M
	for {
		select {
		case b := <-c.send:
			var head struct {
				T string `json:"t"`
			}
			var m M
			if json.Unmarshal(b, &head) == nil && head.T == typ && json.Unmarshal(b, &m) == nil {
				out = append(out, m)
			}
		default:
			return out
		}
	}
}

// Den som tar en näve får ormkortet som spelkort i rundans sammanfattning,
// med namn och text. De andra får inget kort.
func TestTakeShowsCard(t *testing.T) {
	g, a, b := startTwo(t)
	pa := a.player
	pa.X, pa.Y = g.center(), g.center()
	g.reveal(pa.X, pa.Y)
	g.deck = []string{"sover"}
	pa.parts = [3]bool{true, true, true} // ormens del ger annars ett kort till
	summaries(a)
	summaries(b)
	snap := g.snapshot()
	g.take(pa)
	g.summarize(snap, g.wakes)
	got := summaries(a)
	if len(got) != 1 || len(got[0].Kort) != 1 {
		t.Fatalf("rundan: %+v", got)
	}
	if c := got[0].Kort[0]; c.Slag != "ormen" || c.Bild != "sover" || c.Rubrik != "Ormen sover" || c.Text == "" {
		t.Fatalf("kortet: %+v", c)
	}
	for _, s := range summaries(b) {
		if len(s.Kort) != 0 {
			t.Fatalf("Bo fick kort: %+v", s.Kort)
		}
	}
}

// När ormen vaknar får alla kortet: den i kammaren med vad den förlorade,
// den utanför med det den hörde.
func TestWakeShowsCardToAll(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	pa.X, pa.Y = g.center(), g.center()
	g.reveal(pa.X, pa.Y)
	putInRoom(g, pb, 0, 0)
	pa.Treasure, pa.HP = 40, 9
	messagesOf[wakeMsg](a, "vaknar")
	messagesOf[wakeMsg](b, "vaknar")
	g.wakeSerpent()
	ga, gb := messagesOf[wakeMsg](a, "vaknar"), messagesOf[wakeMsg](b, "vaknar")
	if len(ga) != 1 || !ga[0].Inne || ga[0].Forlust != 40-pa.Treasure || ga[0].Skada == 0 || ga[0].Namn != "Ormen vaknar" || ga[0].Text == "" {
		t.Fatalf("i kammaren: %+v", ga)
	}
	if len(gb) != 1 || gb[0].Inne || gb[0].Forlust != 0 || gb[0].Namn != "Ormen vaknar" || gb[0].Text == "" {
		t.Fatalf("utanför: %+v", gb)
	}
}

// Kortet vandraren drog är spelkortet i rundans sammanfattning, och krönikan har
// alla kort vandraren drog på raden med siffror.
func TestDrawnCardsInSummaryAndChronicle(t *testing.T) {
	g, a, _ := startTwo(t)
	pa := a.player
	pa.X, pa.Y = g.center(), g.center()
	g.reveal(pa.X, pa.Y)
	g.deck = []string{"ror_sig", "sover"}
	summaries(a)
	snap := g.snapshot()
	g.take(pa)
	g.summarize(snap, g.wakes)
	found := false
	for _, s := range summaries(a) {
		for _, c := range s.Kort {
			if c.Slag == "ormen" {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("rundans kort saknar ormkortet")
	}
	g.take(pa)
	if tags := g.cardTags(pa); len(tags) != 2 || tags[0].Namn == "" {
		t.Fatalf("krönikans kort: %+v", tags)
	}
	if v := g.journeyOf(pa); len(v.Cards) != 2 {
		t.Fatalf("Din vandring: %+v", v.Cards)
	}
}

// När ormen vaknar kryper den till en ny kammare nära mitten, och den
// gamla blir ett vanligt rum med guld kvar på golvet.
func TestSerpentMovesWhenWaking(t *testing.T) {
	g, a, _ := startTwo(t)
	g.data.Ormen.Spridning = 1
	ox, oy := g.chamber()
	g.reveal(ox, oy)
	g.wakeSerpent()
	nx, ny := g.chamber()
	if nx == ox && ny == oy {
		t.Fatal("kammaren flyttade inte")
	}
	if abs(nx-g.center()) > 1 || abs(ny-g.center()) > 1 {
		t.Fatalf("den nya kammaren på %d,%d", nx, ny)
	}
	if g.rooms[ny][nx].Kind != "kammare" || g.rooms[oy][ox].Kind != "rum" || g.rooms[oy][ox].loose <= 0 {
		t.Fatalf("ny %s, gammal %s med %d guld", g.rooms[ny][nx].Kind, g.rooms[oy][ox].Kind, g.rooms[oy][ox].loose)
	}
	// Guldet tas av den som kommer dit.
	p := a.player
	p.X, p.Y, p.Treasure = ox, oy, 0
	g.pickUpGold(p)
	if p.Treasure <= 0 || g.rooms[oy][ox].loose != 0 {
		t.Fatalf("guldet togs inte: %d", p.Treasure)
	}
}

// När någon tar en näve hör de andra åt vilket håll; tjuvens näve hörs
// inte.
func TestHandfulIsHeard(t *testing.T) {
	g, a, b := startTwo(t)
	p, o := a.player, b.player
	cx, cy := g.chamber()
	g.reveal(cx, cy)
	p.X, p.Y = cx, cy
	putInRoom(g, o, cx, 0)
	g.deck = []string{"sover", "sover"}
	logs(b) // töm det Bo redan har fått
	g.take(p)
	want := g.hornToward(o.X, o.Y, cx, cy)
	got := logs(b)
	heard := false
	for _, l := range got {
		heard = heard || strings.Contains(l, want)
	}
	if !heard {
		t.Fatalf("Bo hörde inte näven åt %s: %q", want, got)
	}
}
