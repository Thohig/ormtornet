package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
