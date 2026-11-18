package main

import "testing"

func countCards(d *Data, wakes bool) int {
	n := 0
	for _, c := range d.Ormen.Kort {
		if c.Vaknar == wakes {
			n += c.Antal
		}
	}
	return n
}

func TestDifficultyPresets(t *testing.T) {
	d, _, err := loadData("data")
	if err != nil {
		t.Fatal(err)
	}
	baseSleep, baseLiv := countCards(d, false), d.Regler.Liv

	// Normal är datafilerna som de är.
	if d.withDifficulty("normal") != d {
		t.Error("normal ska inte ändra något")
	}
	if d.difficultyID("finns_inte") != "normal" || d.difficultyID("") != "normal" {
		t.Error("okänt förval ska ge standard")
	}

	for _, s := range d.Svarighet.Nivaer {
		e := d.withDifficulty(s.ID)
		if s.Liv != nil && e.Regler.Liv != *s.Liv {
			t.Errorf("%s: liv %d", s.ID, e.Regler.Liv)
		}
		if s.Sovkort != nil && countCards(e, false) != *s.Sovkort {
			t.Errorf("%s: sovkort %d", s.ID, countCards(e, false))
		}
		if s.Vaknakort != nil && countCards(e, true) != *s.Vaknakort {
			t.Errorf("%s: vaknakort %d", s.ID, countCards(e, true))
		}
	}
	// Originalet är orört.
	if countCards(d, false) != baseSleep || d.Regler.Liv != baseLiv {
		t.Error("förvalen ändrade datafilerna")
	}

	// Svårare ska vara svårare.
	easy, hard := d.withDifficulty("latt"), d.withDifficulty("mardrom")
	if !(easy.Regler.Liv > baseLiv && hard.Regler.Liv < baseLiv) {
		t.Errorf("liv: lätt %d, mardröm %d", easy.Regler.Liv, hard.Regler.Liv)
	}
	if !(countCards(easy, false) > countCards(hard, false)) {
		t.Error("mardröm ska ha färre sovkort än lätt")
	}
}

func TestDifficultyAppliesToNextGame(t *testing.T) {
	g := testGame(t)
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	s := settingsFrom(g.cfg)
	s.Svarighet = "mardrom"
	got, err := g.ApplySettings(s)
	if err != nil || got.Svarighet != "mardrom" {
		t.Fatalf("%v %+v", err, got)
	}
	if a.player.HP != 4 {
		t.Fatalf("liv i lobbyn %d", a.player.HP)
	}
	g.Ready(a)
	g.Ready(b)
	if g.phase != phasePlay || a.player.HP != 4 || g.deckSleep != 9 {
		t.Fatalf("fas %s liv %d sovkort %d", g.phase, a.player.HP, g.deckSleep)
	}
	if !contains(logs(a), "Svårighet: Mardröm. Du har 4 liv.") {
		t.Error("vandrarna fick inte veta svårigheten")
	}
	// Under vandringen går det inte att byta.
	s.Svarighet = "latt"
	if _, err := g.ApplySettings(s); err == nil || g.cfg.Difficulty != "mardrom" {
		t.Error("svårigheten byttes mitt i vandringen")
	}

	// Ett okänt förval blir standard.
	g2 := testGame(t)
	s.Svarighet = "finns_inte"
	got, _ = g2.ApplySettings(s)
	if got.Svarighet != "normal" {
		t.Errorf("okänt förval gav %q", got.Svarighet)
	}
}
