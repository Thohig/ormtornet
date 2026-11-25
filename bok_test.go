package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestBookHasAllClasses(t *testing.T) {
	g := testGame(t)
	mux := http.NewServeMux()
	addHostRoutes(mux, nil, g, hostAddrs{}, hostHooks{Close: func() {}})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/bok", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	var v bookView
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}
	if len(v.Klasser) != len(g.data.Klasser) || len(v.Klasser) == 0 {
		t.Fatalf("%d klasser i boken, %d i data", len(v.Klasser), len(g.data.Klasser))
	}
	for _, c := range v.Klasser {
		if c.Namn == "" || c.Beskrivning == "" || c.Bild == "" {
			t.Errorf("klassen %s saknar något: %+v", c.ID, c)
		}
	}
}

// Grundernas avsnitt står i bok_grunder, och vart och ett har en rubrik och
// en text på båda språken.
func TestBookBasicsTexts(t *testing.T) {
	for _, lang := range []string{"sv", "en"} {
		d, _, err := loadDataLang("", lang)
		if err != nil {
			t.Fatal(err)
		}
		list := d.Spelarsidan["bok_grunder"]
		if list == "" {
			t.Fatalf("%s: bok_grunder saknas", lang)
		}
		for _, id := range strings.Fields(list) {
			for _, k := range []string{"bok_" + id + "_rubrik", "bok_" + id} {
				if d.Spelarsidan[k] == "" {
					t.Errorf("%s: %s saknas", lang, k)
				}
			}
		}
	}
}

// Tornet skriver in i boken det vandraren möter, och skickar det i
// tillståndet: föremål, förbannelser, varelser och fällor.
func TestBookLearnsWhatYouMeet(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	putInRoom(g, p, 3, 3)
	g.give(p, "svart_parla")
	g.curse(p, "svart_parla")
	g.trigger(p, g.data.trap("pilar"), 0)
	g.addMonster(g.data.monsterType("krypare"), 3, 3)
	g.announceMonster(p)
	got := g.stateFor(p).Book
	for _, want := range []string{"foremal:svart_parla", "forbannelse:svart_parla", "falla:pilar", "varelse:krypare"} {
		if !slices.Contains(got, want) {
			t.Errorf("%s saknas i boken: %v", want, got)
		}
	}
	if slices.Contains(got, "varelse:yngel") || slices.Contains(got, "foremal:rep") {
		t.Errorf("boken har sådant vandraren inte har mött: %v", got)
	}
}

// Boken visar korta beskrivningar och aldrig namnen på fällor och
// varelser, och inga hopfogade föremål än.
func TestBookHasNoNames(t *testing.T) {
	g := testGame(t)
	v := g.Book()
	if len(v.Foremal) == 0 || len(v.Varelser) != len(g.data.Monster.Typer) || len(v.Fallor) != len(g.data.Fallor) {
		t.Fatalf("%d föremål, %d varelser, %d fällor", len(v.Foremal), len(v.Varelser), len(v.Fallor))
	}
	for _, c := range append(v.Varelser, v.Fallor...) {
		if c.Om == "" || c.Rubrik == "" {
			t.Errorf("%s saknar rubrik eller text: %+v", c.ID, c)
		}
		for _, m := range g.data.Monster.Typer {
			if strings.Contains(c.Rubrik+c.Text+c.Om, m.Namn) {
				t.Errorf("%s nämner varelsen %s", c.ID, m.Namn)
			}
		}
		for _, f := range g.data.Fallor {
			if strings.Contains(c.Rubrik+c.Text+c.Om, f.Namn) {
				t.Errorf("%s nämner fällan %s", c.ID, f.Namn)
			}
		}
	}
	for _, c := range v.Foremal {
		if it := g.data.itemType(c.ID); it == nil || len(it.Av) > 0 {
			t.Errorf("%s ska inte stå bland föremålen", c.ID)
		}
	}
}
