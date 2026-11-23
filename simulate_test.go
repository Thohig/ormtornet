package main

import (
	"strings"
	"testing"
)

func TestSimulation(t *testing.T) {
	run := func() string {
		d, _, err := loadData("data")
		if err != nil {
			t.Fatal(err)
		}
		var b strings.Builder
		cfg := Config{SunMax: 40, Classes: true}
		if err := runSimulation(d, cfg, 20, 4, 5, &b); err != nil {
			t.Fatal(err)
		}
		// Raden med tiden skiljer sig mellan körningarna.
		lines := strings.Split(b.String(), "\n")
		return strings.Join(append(lines[:1], lines[2:]...), "\n")
	}
	a, b := run(), run()
	for _, want := range []string{"20 vandringar med 4 AI-spelare", "Kom ut med skatt", "Personligheterna", "Klasserna"} {
		if !strings.Contains(a, want) {
			t.Fatalf("saknar %q i:\n%s", want, a)
		}
	}
	if a != b {
		t.Fatalf("samma frön gav olika resultat:\n%s\n---\n%s", a, b)
	}
}

func TestSimCause(t *testing.T) {
	for kind, want := range map[string]string{"falla": "fälla", "vandrare": "en annan vandrare", "": "okänt", "gift": "gift"} {
		if got := simCause(kind); got != want {
			t.Errorf("simCause(%q) = %q, vill ha %q", kind, got, want)
		}
	}
}

// Simuleringen ger varje klass lika ofta, och lika ofta till varje
// personlighet.
func TestSimDealsClassesEvenly(t *testing.T) {
	d, _, err := loadData("data")
	if err != nil {
		t.Fatal(err)
	}
	d.AI.BetanketidMs = []int{3600000, 3600000}
	k := len(d.Klasser)
	// Lika många AI-spelare för varje personlighet, så att de turas om jämnt.
	bots := len(d.AI.Personligheter)
	for bots*2 <= maxPlayers {
		bots *= 2
	}
	count := map[[2]string]int{}
	for seed := uint64(1); seed <= uint64(2*k); seed++ {
		r, err := simulateOne(d, Config{SunMax: 40, Classes: true}, bots, seed)
		if err != nil {
			t.Fatal(err)
		}
		for _, w := range r.walkers {
			count[[2]string{w.class, w.pers}]++
		}
	}
	want := 2 * k * bots / (k * len(d.AI.Personligheter))
	for _, c := range d.Klasser {
		for _, p := range d.AI.Personligheter {
			if n := count[[2]string{c.Kort, p.Namn}]; n != want {
				t.Errorf("%s spelades %d gånger av den %s, ville ha %d", c.Kort, n, p.Namn, want)
			}
		}
	}
}
