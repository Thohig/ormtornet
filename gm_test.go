package main

import (
	"strings"
	"testing"
)

func gmSettings(g *Game, on bool) Settings {
	s := settingsFrom(g.cfg)
	s.Spelledare = on
	return s
}

func TestGMCannotJoinFromHostComputer(t *testing.T) {
	g := testGame(t)
	g.ApplySettings(gmSettings(g, true))
	host := testClient()
	host.local = true
	g.Join(host, "Spelledaren")
	if host.player != nil {
		t.Fatal("spelledaren kom in som vandrare från sin egen dator")
	}
	other := testClient()
	g.Join(other, "Alva")
	if other.player == nil {
		t.Fatal("en spelare från en annan dator kom inte in")
	}
}

func TestGMModeKicksLocalPlayer(t *testing.T) {
	g := testGame(t)
	host := testClient()
	host.local = true
	g.Join(host, "Stefan")
	for len(host.send) > 0 {
		<-host.send
	}
	g.ApplySettings(gmSettings(g, true))
	if host.player != nil || g.playerByName("Stefan") != nil {
		t.Fatal("vandraren från spelledarens dator är kvar")
	}
	told := false
	for len(host.send) > 0 {
		if strings.Contains(string(<-host.send), `"t":"spelledare"`) {
			told = true
		}
	}
	if !told {
		t.Fatal("vandraren fick inte veta varför")
	}
}

func TestGMStartsWithOnlyBots(t *testing.T) {
	g := testGame(t)
	s := settingsFrom(g.cfg)
	s.AI = 2
	g.ApplySettings(s)
	if g.ForceStart() {
		t.Fatal("startade med bara AI-spelare utan spelledarläge")
	}
	s.Spelledare = true
	g.ApplySettings(s)
	if g.phase != phaseLobby {
		t.Fatal("vandringen startade av sig själv utan att spelledaren bad om det")
	}
	if !g.ForceStart() || g.phase != phasePlay {
		t.Fatal("spelledaren kunde inte starta med bara AI-spelare")
	}
}

func TestGMViewShowsEverything(t *testing.T) {
	g, a, b := twoPlaying(t, 0)
	a.player.Treasure = 42
	r := putInRoom(g, b.player, 4, 4)
	r.Event = "pilar"
	g.Act(a, actWait, 0)
	g.logAll("Ett test hörs i hela tornet.", "info")
	v := g.GMView(false)
	found := false
	for _, p := range v.Players {
		if p.Name == "Alva" && p.Treasure == 42 {
			found = true
		}
	}
	if !found {
		t.Fatal("spelledaren ser inte Alvas skatt")
	}
	if v.Traps == 0 {
		t.Fatal("spelledaren ser inte fällan som väntar")
	}
	last := v.Log[len(v.Log)-1]
	if last.Text != "Ett test hörs i hela tornet." || len(last.To) != 2 {
		t.Fatalf("loggen samlar inte samma rad till alla: %+v", last)
	}
	if !strings.HasPrefix(g.actionText(a.player), "väntar") {
		t.Fatalf("draget beskrivs som %q", g.actionText(a.player))
	}
}

// När vandringen är slut har spelledarbordet resultatet och krönikan, och
// introt på tornets dator leder dit i spelledarläget.
func TestGMViewEndAndIntro(t *testing.T) {
	g := testGame(t)
	if g.Intro().Leder {
		t.Fatal("introt leder till spelledarbordet utan spelledarläge")
	}
	s := gmSettings(g, true)
	s.AI = 2
	g.ApplySettings(s)
	if !g.Intro().Leder {
		t.Fatal("introt leder inte till spelledarbordet i spelledarläget")
	}
	if v := g.GMView(false); v.Results != nil || v.Chronicle != nil {
		t.Fatal("resultat och krönika före vandringen")
	}
	if !g.ForceStart() {
		t.Fatal("vandringen startade inte")
	}
	g.mu.Lock()
	g.endGame()
	g.mu.Unlock()
	v := g.GMView(false)
	if len(v.Results) != 2 || len(v.Chronicle) < 3 || v.Chronicle[0].Class != "rubrik" {
		t.Fatalf("efter vandringen: %d resultat, krönikan %+v", len(v.Results), v.Chronicle)
	}
}

// Spelledarbordet visar ormens namn, vandrarnas delar, varslen och portarna.
func TestGMViewNameOmensPorts(t *testing.T) {
	g, cs := nameGame(t, 2)
	g.cfg.GMMode = true
	g.level = 2
	g.dealNameParts()
	g.data.Regler.Handelsechans = 1
	g.dealOmens()
	v := g.GMView(false)
	if v.SerpentName != g.serpentNameText() || v.SerpentName == "" {
		t.Fatalf("namnet %q", v.SerpentName)
	}
	for _, p := range v.Players {
		if len(p.Parts) != 3 {
			t.Fatalf("%s: delarna %v", p.Name, p.Parts)
		}
	}
	if len(v.Omens) != len(cs) || v.Omens[0].What == "" {
		t.Fatalf("varslen %+v", v.Omens)
	}
	if v.Ports != 4 || v.OpenPorts != 4 || v.Mercy {
		t.Fatalf("portar %d/%d, nåd %v", v.OpenPorts, v.Ports, v.Mercy)
	}
}
