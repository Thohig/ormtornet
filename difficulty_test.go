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

// Skadan läggs till fällor och varelser som skadar, och solen gångras när
// den är automatisk.
func TestDifficultyDamageAndSun(t *testing.T) {
	d, _, err := loadData("data")
	if err != nil {
		t.Fatal(err)
	}
	one, minus := 1, -5
	d.Svarighet.Nivaer = append(d.Svarighet.Nivaer,
		DifficultyStep{ID: "hard", Skada: &one},
		DifficultyStep{ID: "mild", Skada: &minus})
	hard, mild := d.withDifficulty("hard"), d.withDifficulty("mild")
	for i, tr := range d.Fallor {
		want, wantMild := tr.Skada, tr.Skada
		if tr.Skada > 0 {
			want, wantMild = tr.Skada+1, 1
		}
		if hard.Fallor[i].Skada != want || mild.Fallor[i].Skada != wantMild {
			t.Errorf("fällan %s: %d och %d, ville ha %d och %d", tr.ID, hard.Fallor[i].Skada, mild.Fallor[i].Skada, want, wantMild)
		}
	}
	for i, m := range d.Monster.Typer {
		if m.Skada > 0 && hard.Monster.Typer[i].Skada != m.Skada+1 {
			t.Errorf("varelsen %s: %d", m.ID, hard.Monster.Typer[i].Skada)
		}
	}
	half := 1.5
	d.Svarighet.Nivaer = append(d.Svarighet.Nivaer, DifficultyStep{ID: "rik", Kammarskatt: &half})
	rich := d.withDifficulty("rik")
	if rich.Ormen.SkattMin != (d.Ormen.SkattMin*3+1)/2 || rich.Ormen.SkattMax != (d.Ormen.SkattMax*3+1)/2 {
		t.Errorf("kammarskatt 1,5: %d–%d av %d–%d", rich.Ormen.SkattMin, rich.Ormen.SkattMax, d.Ormen.SkattMin, d.Ormen.SkattMax)
	}

	orig, _, _ := loadData("data")
	for i := range d.Fallor {
		if d.Fallor[i].Skada != orig.Fallor[i].Skada {
			t.Errorf("fällan %s ändrades i originalet", d.Fallor[i].ID)
		}
	}

	g := testGame(t)
	g.cfg.SunMax = 0 // automatisk sol
	sun := g.scaledSun()
	g.data.Regler.Sol = 1.25
	if got := g.scaledSun(); got != autoSun(g.size)*5/4 {
		t.Errorf("solen med 1,25 gav %d, förut %d", got, sun)
	}
	g.cfg.SunMax = 20
	if got := g.scaledSun(); got != 20 {
		t.Errorf("-sol 20 gav %d", got)
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
	if a.player.HP != 6 {
		t.Fatalf("liv i lobbyn %d", a.player.HP)
	}
	g.Ready(a)
	g.Ready(b)
	if g.phase != phasePlay || a.player.HP != 6 || g.deckSleep != 9 {
		t.Fatalf("fas %s liv %d sovkort %d", g.phase, a.player.HP, g.deckSleep)
	}
	if !contains(logs(a), "Svårighet: Mardröm. Du har 6 liv.") {
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
