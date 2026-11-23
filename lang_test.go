package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// withTranslation ger en datamapp med de svenska filerna och en engelsk fil.
func withTranslation(t *testing.T, name, content string) string {
	t.Helper()
	dir := dataWith(t, "", func(s string) string { return s })
	if err := os.MkdirAll(filepath.Join(dir, langEn), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, langEn, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestDataLoadsInEveryLanguage(t *testing.T) {
	for _, lang := range languages {
		if _, _, err := loadDataLang("data", lang); err != nil {
			t.Fatalf("%s: %v", lang, err)
		}
	}
}

// En översättning byter texterna men inte siffrorna, och det som inte är
// översatt står kvar på svenska.
func TestTranslationKeepsNumbers(t *testing.T) {
	dir := withTranslation(t, "klasser.json", `{"klasser": [
		{"id": "kartritare", "namn": "The Mapmaker", "styrka": "Sees {inget} through doors."}
	]}`)
	if _, _, err := loadDataLang(dir, langEn); err == nil || !strings.Contains(err.Error(), "{inget}") {
		t.Fatalf("en okänd platshållare ska ge fel, fick %v", err)
	}
	dir = withTranslation(t, "klasser.json", `{"klasser": [
		{"id": "kartritare", "namn": "The Mapmaker", "beskrivning": "Has drawn maps."}
	]}`)
	sv, _, err := loadDataLang(dir, langSv)
	if err != nil {
		t.Fatal(err)
	}
	en, src, err := loadDataLang(dir, langEn)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(src, filepath.Join(langEn, "klasser.json")) {
		t.Fatalf("källorna nämner inte översättningen: %s", src)
	}
	a, b := sv.class("kartritare"), en.class("kartritare")
	if b.Namn != "The Mapmaker" || b.Beskrivning != "Has drawn maps." {
		t.Fatalf("inte översatt: %+v", b)
	}
	if b.Styrka != a.Styrka || b.Liv != a.Liv || b.Drycker != a.Drycker || len(b.Varden) != len(a.Varden) {
		t.Fatalf("det som inte är översatt ska vara som på svenska: %+v", b)
	}
	if en.class("legoknekt").Namn != sv.class("legoknekt").Namn {
		t.Fatal("en klass som inte är översatt ska stå kvar på svenska")
	}
}

func TestTranslationOnlyTexts(t *testing.T) {
	for _, c := range []struct{ name, content, want string }{
		{"klasser.json", `{"klasser": [{"id": "kartritare", "liv": 2}]}`, "bara innehålla texter"},
		{"klasser.json", `{"klasser": [{"id": "trollkarl", "namn": "Wizard"}]}`, "trollkarl"},
		{"klasser.json", `{"klasser": [{"id": "kartritare", "nam": "Mapmaker"}]}`, "nam"},
		{"svarighet.json", `{"svarighet": {"standard": "svar"}}`, "översätts inte"},
		{"ai.json", `{"ai": {"personligheter": [{"id": "djarv", "klasser": ["thief"]}]}}`, "översätts inte"},
	} {
		dir := withTranslation(t, c.name, c.content)
		if _, _, err := loadDataLang(dir, langEn); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: ville ha ett fel om %q, fick %v", c.content, c.want, err)
		}
	}
}

// Svenska läser aldrig översättningarna, även om de är trasiga.
func TestSwedishIgnoresTranslations(t *testing.T) {
	dir := withTranslation(t, "klasser.json", `{"klasser": [{"id": "kartritare", "liv": 2}]}`)
	if _, _, err := loadDataLang(dir, langSv); err != nil {
		t.Fatal(err)
	}
}

func TestTranslationStatus(t *testing.T) {
	status, err := translationStatus("data", langEn)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range status {
		if s.Done > s.Texts {
			t.Errorf("%s: %d av %d", s.File, s.Done, s.Texts)
		}
		if s.Texts == 0 || s.Done != s.Texts {
			t.Errorf("%s: %d av %d texter är översatta till engelska", s.File, s.Done, s.Texts)
		}
	}
}

// Byter spelledaren språk på värdsidan läses texterna in på nytt.
func TestSettingsChangeLanguage(t *testing.T) {
	g := testGame(t)
	var asked string
	var saved Settings
	mux := http.NewServeMux()
	addHostRoutes(mux, nil, g, hostAddrs{}, hostHooks{
		Close: func() {},
		Save:  func(s Settings) error { saved = s; return nil },
		Load: func(lang string) (*Data, error) {
			asked = lang
			d, _, err := loadDataLang("data", lang)
			return d, err
		},
	})
	post := func(body string) map[string]any {
		r := httptest.NewRequest("POST", "/spelledare/installningar", strings.NewReader(body))
		r.RemoteAddr = "127.0.0.1:5000"
		r.Header.Set("X-Ormtornet", "spelledare")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		var out map[string]any
		json.NewDecoder(w.Body).Decode(&out)
		if w.Code != http.StatusOK {
			t.Fatalf("kod %d: %v", w.Code, out)
		}
		return out
	}
	out := post(`{"ai":0,"min":2,"sol":40,"turtid":30,"sprak":"en","svarighet":"latt"}`)
	if asked != langEn || saved.Sprak != langEn || g.Settings().Sprak != langEn {
		t.Fatalf("frågade efter %q, sparade %q", asked, saved.Sprak)
	}
	if out["text"] == nil {
		t.Fatalf("svaret säger inte att språket gäller: %v", out)
	}
	if s := g.base.difficulty("latt"); s == nil || s.Namn != "Easy" {
		t.Fatalf("svårigheten heter %+v", s)
	}
	asked = ""
	post(`{"ai":0,"min":2,"sol":40,"turtid":30,"sprak":"en"}`)
	if asked != "" {
		t.Fatal("texterna lästes om fast språket var detsamma")
	}
	post(`{"ai":0,"min":2,"sol":40,"turtid":30,"sprak":"xx"}`)
	if g.Settings().Sprak != langSv || g.base.difficulty("latt").Namn != "Lätt" {
		t.Fatal("ett okänt språk ska bli svenska")
	}
}

// Varje textnyckel som sidorna och koden använder finns i datafilerna, på
// svenska och, när översättningen är klar, på engelska. Nycklar som byggs
// ihop i koden fångas inte här.
func TestPageTextKeysExist(t *testing.T) {
	sv, _, err := loadDataLang("data", langSv)
	if err != nil {
		t.Fatal(err)
	}
	re := regexp.MustCompile(`(?:\bt\(|data-t(?:-html|-title|-placeholder|-aria)?=)"([a-z0-9_]+)"`)
	for file, texts := range map[string]map[string]string{
		"web/index.html": sv.Spelarsidan, "web/app.js": sv.Spelarsidan,
		"web/spelledare.html": sv.Vardsidan,
		"web/torn.html":       sv.Vardsidan, "web/torn.js": sv.Vardsidan,
	} {
		raw, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range re.FindAllStringSubmatch(string(raw), -1) {
			if m[1] == "nyckel" || strings.HasSuffix(m[1], "_") {
				continue // ett exempel i en kommentar, eller början på en nyckel som byggs ihop
			}
			if _, ok := texts[m[1]]; !ok {
				t.Errorf("%s använder %q, som inte finns i datafilerna", file, m[1])
			}
		}
	}
	goRe := regexp.MustCompile(`HostText\("([a-z0-9_]+)"|"(server_[a-z0-9_]+)"`)
	files, _ := filepath.Glob("*.go")
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		raw, _ := os.ReadFile(file)
		for _, m := range goRe.FindAllStringSubmatch(string(raw), -1) {
			key := m[1] + m[2]
			if _, ok := sv.Vardsidan[key]; !ok {
				t.Errorf("%s använder %q, som inte finns i vardsidan.json", file, key)
			}
		}
	}
	mapRaw, _ := os.ReadFile("web/map.js")
	for _, m := range regexp.MustCompile(`name: "(zoom_[a-z]+)"`).FindAllStringSubmatch(string(mapRaw), -1) {
		if _, ok := sv.Spelarsidan[m[1]]; !ok {
			t.Errorf("web/map.js använder %q, som inte finns i spelarsidan.json", m[1])
		}
	}
}
