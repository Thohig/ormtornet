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
	if !strings.Contains(all, "+15 guld") {
		t.Fatalf("sammanfattningen säger inte vad som hände: %q", all)
	}
	// Vart vandraren gick och solen står i rumsblocket och översta raden.
	if strings.Contains(all, "österut") || strings.Contains(all, g.roomShort(n)) {
		t.Fatalf("kortet säger vart vandraren gick: %q", all)
	}
	// Bo valde inget. Det säger loggen, så något kort kommer inte, men
	// raden står i Din vandring.
	if bo := summaries(b); len(bo) != 0 {
		t.Fatalf("Bo fick ett kort fast inget hände: %+v", bo)
	}
	if j := b.player.journal; len(j) == 0 || len(j[len(j)-1].Lines) != 1 || j[len(j)-1].Lines[0].Text != g.data.Text("sam_inget") {
		t.Fatalf("Bos vandring: %+v", j)
	}
}

func TestSummaryNamesTheFoe(t *testing.T) {
	g, a, b := fightReady(t)
	alva, bo := a.player, b.player
	summaries(a)
	summaries(b)
	g.Act(a, actAttack, -bo.ID)
	g.Act(b, actAttack, -alva.ID)
	has := func(c *Client, want string) bool {
		for _, m := range summaries(c) {
			for _, l := range m.Lines {
				if l.Text == want {
					return true
				}
			}
		}
		return false
	}
	if want := g.data.Text("sam_anfall_vandrare", "namn", bo.Name); !has(a, want) {
		t.Fatalf("Alvas kort saknar %q", want)
	}
	if want := g.data.Text("sam_anfall_vandrare", "namn", alva.Name); !has(b, want) {
		t.Fatalf("Bos kort saknar %q", want)
	}
}

// En runda med rummet men utan rader ger en tom lista, inte null, så att
// Din vandring på spelarsidan kan läsa den.
func TestJourneyRoundHasLines(t *testing.T) {
	g, a, b := startTwo(t)
	g.Act(a, actWait, 0)
	g.Act(b, actWait, 0)
	for _, jr := range a.player.journal {
		if jr.Lines == nil {
			t.Fatalf("runda %d har inga rader alls (null)", jr.Round)
		}
	}
	if len(a.player.journal) == 0 {
		t.Fatal("ingen runda i vandringen")
	}
}

// Tysta rundor slås ihop i Din vandring, och skadan säger vad den kom av.
func TestJourneyQuietRoundsAndCause(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.journal = nil
	stood := []summaryLine{{Text: g.data.Text("sam_inget"), journey: true}}
	for r := 10; r <= 12; r++ {
		g.round = r
		g.addJourney(p, stood, nil, false, false)
	}
	if len(p.journal) != 1 || p.journal[0].Round != 10 || p.journal[0].To != 12 {
		t.Fatalf("tre tysta rundor väntades bli en, fick %+v", p.journal)
	}
	// Tornet styrde nästa runda: en egen rad, med märket.
	g.round = 13
	g.addJourney(p, nil, nil, true, false)
	if len(p.journal) != 2 || !p.journal[1].Tower {
		t.Fatalf("rundan då tornet styrde väntades stå för sig: %+v", p.journal)
	}

	trap := g.data.trap("pilar")
	if trap == nil {
		t.Skip("ingen fälla med pilar")
	}
	g.round = 14
	before := g.snapshot()[p]
	p.HP -= 2
	p.hurt("falla", trap.Kort)
	lines := g.summaryFor(p, before, false)
	g.addJourney(p, lines, nil, false, false)
	last := p.journal[len(p.journal)-1]
	if len(last.Lines) == 0 || !strings.HasPrefix(last.Lines[0].Text, upperFirst(trap.Kort)) {
		t.Fatalf("skadan sa inte vad den kom av: %+v", last)
	}
	if strings.HasPrefix(lines[0].Text, upperFirst(trap.Kort)) {
		t.Fatalf("kortet fick orsaken också: %+v", lines)
	}
}

// En tyst runda då vandraren satt fast (en grop, ett galler) säger det, och
// slås inte ihop med rundor då hen stod kvar av sig själv.
func TestJourneyStuckRound(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.journal = nil
	stood := []summaryLine{{Text: g.data.Text("sam_inget"), journey: true}}
	g.round = 3
	g.addJourney(p, stood, nil, false, true)
	g.round = 4
	g.addJourney(p, stood, nil, false, false)
	if len(p.journal) != 2 {
		t.Fatalf("väntade två rundor, fick %+v", p.journal)
	}
	if got := p.journal[0].Lines[0].Text; got != g.data.Text("sam_tyst_fast") {
		t.Fatalf("runda 3: %q", got)
	}
	if got := p.journal[1].Lines[0].Text; got != g.data.Text("sam_tyst") {
		t.Fatalf("runda 4: %q", got)
	}
}
