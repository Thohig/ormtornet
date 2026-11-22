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
