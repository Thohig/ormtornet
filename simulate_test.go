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
		cfg := Config{Reveal: true, SunMax: 40, Classes: true}
		if err := runSimulation(d, cfg, 20, 4, 5, &b); err != nil {
			t.Fatal(err)
		}
		// Raden med tiden skiljer sig mellan körningarna.
		lines := strings.Split(b.String(), "\n")
		return strings.Join(append(lines[:1], lines[2:]...), "\n")
	}
	a, b := run(), run()
	for _, want := range []string{"20 vandringar med 4 AI-spelare", "Kom ut med skatt", "Klasserna"} {
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
