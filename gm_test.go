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
	v := g.GMView()
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
