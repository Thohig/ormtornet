package main

import (
	"encoding/json"
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

// Den som tar en näve får kortet den drog, med namn, text och guld. De
// andra får inget kort.
func TestTakeShowsCard(t *testing.T) {
	g, a, b := startTwo(t)
	pa := a.player
	pa.X, pa.Y = g.center(), g.center()
	g.reveal(pa.X, pa.Y)
	g.deck = []string{"sover"}
	messagesOf[cardMsg](a, "ormkort")
	messagesOf[cardMsg](b, "ormkort")
	before := pa.Treasure
	g.take(pa)
	got := messagesOf[cardMsg](a, "ormkort")
	if len(got) != 1 || got[0].ID != "sover" || got[0].Namn != "Ormen sover" || got[0].Guld != pa.Treasure-before || got[0].Text == "" {
		t.Fatalf("kortet: %+v", got)
	}
	if n := len(messagesOf[cardMsg](b, "ormkort")); n != 0 {
		t.Fatalf("Bo fick %d kort", n)
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

// Kortet vandraren drog står i rundans sammanfattning, och krönikan har
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
		for _, l := range s.Lines {
			if l.Ormkort != "" {
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
