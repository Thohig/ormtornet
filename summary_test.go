package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// summaries plockar ut rundans sammanfattningar ur det klienten har fått.
func summaries(c *Client) []summaryMsg {
	var out []summaryMsg
	for {
		select {
		case b := <-c.send:
			var m summaryMsg
			if json.Unmarshal(b, &m) == nil && m.T == "runda" {
				out = append(out, m)
			}
		default:
			return out
		}
	}
}

func TestRoundSummary(t *testing.T) {
	g, a, b := startTwo(t)
	pa := a.player
	r := putInRoom(g, pa, 4, 4)
	r.Doors[East] = true
	n := g.reveal(5, 4)
	n.Kind, n.Event, n.Doors[West] = "rum", "", true
	g.monsters = nil
	summaries(a)
	pa.Action = "E"
	snap := g.snapshot()
	g.move(pa, East, 0)
	pa.HP -= 2
	pa.Treasure += 15
	g.summarize(snap, g.wakes)
	got := summaries(a)
	if len(got) != 1 {
		t.Fatalf("%d sammanfattningar", len(got))
	}
	lines := got[0].Lines
	if lines[0].C != "varn" || !strings.Contains(lines[0].Text, "−2 liv") {
		t.Fatalf("det farliga står inte först: %+v", lines)
	}
	var text []string
	for _, l := range lines {
		text = append(text, l.Text)
		if strings.HasPrefix(l.Text, "[") {
			t.Fatalf("text saknas i data: %s", l.Text)
		}
	}
	all := strings.Join(text, " ")
	if !strings.Contains(all, "österut till "+g.roomName(n)) || !strings.Contains(all, "+15 guld") {
		t.Fatalf("sammanfattningen säger inte vad som hände: %q", all)
	}
	bo := summaries(b)
	if len(bo) != 1 || strings.Contains(bo[0].Lines[0].Text, "guld") || bo[0].Lines[0].Text != g.data.Text("sam_inget") {
		t.Fatalf("Bos sammanfattning: %+v", bo)
	}
}
