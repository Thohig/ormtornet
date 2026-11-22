package main

// Språk: tornet kan spelas på svenska eller engelska. De svenska filerna i
// data/ är grunden. En översättning ligger i en undermapp med språkets kod,
// till exempel data/en/, med samma filnamn. Den innehåller bara texter;
// siffror och val finns bara i de svenska filerna, så att en ändring av
// balansen görs på ett ställe. Listor som klasser, varelser och fällor
// kopplas ihop med sina id:n. Text som saknas i översättningen står kvar
// på svenska.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const (
	langSv = "sv"
	langEn = "en"
)

// languages är språken tornet kan visas på, med det första som förval.
var languages = []string{langSv, langEn}

// langOrDefault ger lang om tornet kan visas på det, annars svenska.
func langOrDefault(lang string) string {
	if slices.Contains(languages, lang) {
		return lang
	}
	return langSv
}

// readTranslation läser översättningen av datafilen name till lang, ur
// mappen dir/lang eller ur den inbyggda kopian. Svarar nil om filen saknas.
func readTranslation(dir, lang, name string) ([]byte, string, error) {
	if lang == langSv {
		return nil, "", nil
	}
	if dir != "" {
		path := filepath.Join(dir, lang, name)
		raw, err := os.ReadFile(path)
		if err == nil {
			return raw, path, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return nil, "", err
		}
	}
	raw, err := embeddedData.ReadFile("data/" + lang + "/" + name)
	if err != nil {
		return nil, "", nil // ingen översättning; texterna står kvar på svenska
	}
	return raw, lang + "/" + name + " (inbyggd)", nil
}

// translateJSON lägger texterna i tr över originalet orig och ger den
// sammanslagna filen. Översättningen får bara innehålla texter som finns i
// originalet.
func translateJSON(orig, tr []byte) ([]byte, error) {
	o, err := decodeTree(orig)
	if err != nil {
		return nil, err
	}
	t, err := decodeTree(tr)
	if err != nil {
		return nil, err
	}
	merged, err := overlay(o, t, "")
	if err != nil {
		return nil, err
	}
	return json.Marshal(merged)
}

// decodeTree läser JSON utan att röra talen.
func decodeTree(raw []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var tree any
	if err := dec.Decode(&tree); err != nil {
		return nil, explainJSON(raw, err)
	}
	return tree, nil
}

// overlay lägger t över o. Objekt slås ihop fält för fält, listor med
// objekt kopplas ihop med id, och listor med texter byts ut helt.
func overlay(o, t any, path string) (any, error) {
	switch tv := t.(type) {
	case map[string]any:
		ov, ok := o.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%s: ska inte vara ett objekt", pathName(path))
		}
		for k, e := range tv {
			if strings.HasPrefix(k, "_") {
				continue // kommentar
			}
			oe, found := ov[k]
			if !found {
				return nil, fmt.Errorf("%s: finns inte i den svenska filen", pathName(path+"."+k))
			}
			if k == "id" || isCode(path+"."+k) {
				if fmt.Sprint(e) != fmt.Sprint(oe) {
					return nil, fmt.Errorf("%s: är ingen text och översätts inte", pathName(path+"."+k))
				}
			}
			merged, err := overlay(oe, e, path+"."+k)
			if err != nil {
				return nil, err
			}
			ov[k] = merged
		}
		return ov, nil
	case []any:
		ol, ok := o.([]any)
		if !ok {
			return nil, fmt.Errorf("%s: ska inte vara en lista", pathName(path))
		}
		if hasIDs(ol) {
			for _, e := range tv {
				id := idOf(e)
				i := slices.IndexFunc(ol, func(x any) bool { return idOf(x) == id })
				if id == "" || i < 0 {
					return nil, fmt.Errorf("%s: id %q finns inte i den svenska filen", pathName(path), id)
				}
				merged, err := overlay(ol[i], e, fmt.Sprintf("%s[%s]", path, id))
				if err != nil {
					return nil, err
				}
				ol[i] = merged
			}
			return ol, nil
		}
		if !onlyTexts(tv) {
			return nil, fmt.Errorf("%s: en översättning får bara innehålla texter", pathName(path))
		}
		return tv, nil
	case string:
		if _, ok := o.(string); !ok {
			return nil, fmt.Errorf("%s: ska inte vara en text", pathName(path))
		}
		return tv, nil
	}
	return nil, fmt.Errorf("%s: en översättning får bara innehålla texter, inte tal eller sant/falskt", pathName(path))
}

// isCode svarar om fältet på path ser ut som text men är ett val som koden
// läser, till exempel förvalet bland svårigheterna.
func isCode(path string) bool {
	if path == ".svarighet.standard" || strings.HasPrefix(path, ".spelkort.") {
		return true
	}
	return strings.HasPrefix(path, ".ai.personligheter[") && strings.HasSuffix(path, "].klasser")
}

func pathName(path string) string {
	if path == "" {
		return "filen"
	}
	return strings.TrimPrefix(path, ".")
}

// hasIDs svarar om listan består av objekt med id.
func hasIDs(l []any) bool {
	return len(l) > 0 && idOf(l[0]) != ""
}

func idOf(v any) string {
	if m, ok := v.(map[string]any); ok {
		if id, ok := m["id"].(string); ok {
			return id
		}
	}
	return ""
}

// onlyTexts svarar om listan bara innehåller texter, eller listor med texter.
func onlyTexts(l []any) bool {
	for _, e := range l {
		switch ev := e.(type) {
		case string:
		case []any:
			if !onlyTexts(ev) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

var placeholderRe = regexp.MustCompile(`\{[^{}\s]+\}`)

// placeholders ger platshållarna i texterna under v, till exempel {namn}.
func placeholders(v any, into map[string]bool) {
	switch tv := v.(type) {
	case string:
		for _, p := range placeholderRe.FindAllString(tv, -1) {
			into[p] = true
		}
	case []any:
		for _, e := range tv {
			placeholders(e, into)
		}
	}
}

// translationReport jämför en översättning med originalet. Den räknar hur
// många av originalets texter som är översatta och listar platshållare i
// översättningen som originalet inte har, eftersom de skulle synas som
// de står i spelet.
type translationReport struct {
	texts, done int
	problems    []string
}

func (r *translationReport) compare(o, t any, path string) {
	switch ov := o.(type) {
	case map[string]any:
		tm, _ := t.(map[string]any)
		for k, e := range ov {
			if strings.HasPrefix(k, "_") || k == "id" || k == "tecken" || k == "modell" || k == "bild" || isCode(path+"."+k) || path+"."+k == ".ai.namn" || path+"."+k == ".vandrarnamn.fornamn" || path+"."+k == ".avslojande.namn" {
				continue // inget att översätta: id, tecken och modeller på kartan, val i koden, AI-namnen och ormens namn
			}
			var te any
			if tm != nil {
				te = tm[k]
			}
			r.compare(e, te, path+"."+k)
		}
	case []any:
		if hasIDs(ov) {
			tl, _ := t.([]any)
			for _, e := range ov {
				id := idOf(e)
				var te any
				if i := slices.IndexFunc(tl, func(x any) bool { return idOf(x) == id }); i >= 0 {
					te = tl[i]
				}
				r.compare(e, te, fmt.Sprintf("%s[%s]", path, id))
			}
			return
		}
		if !onlyTexts(ov) {
			return
		}
		n := countTexts(ov)
		r.texts += n
		if t == nil {
			return
		}
		r.done += n
		r.checkPlaceholders(o, t, path)
	case string:
		r.texts++
		if t == nil {
			return
		}
		r.done++
		r.checkPlaceholders(o, t, path)
	}
}

func (r *translationReport) checkPlaceholders(o, t any, path string) {
	have := map[string]bool{}
	placeholders(o, have)
	got := map[string]bool{}
	placeholders(t, got)
	var extra []string
	for p := range got {
		if !have[p] {
			extra = append(extra, p)
		}
	}
	slices.Sort(extra)
	if len(extra) > 0 {
		r.problems = append(r.problems, fmt.Sprintf("%s: %s finns inte i den svenska texten", pathName(path), strings.Join(extra, ", ")))
	}
}

func countTexts(l []any) int {
	n := 0
	for _, e := range l {
		if sub, ok := e.([]any); ok {
			n += countTexts(sub)
		} else {
			n++
		}
	}
	return n
}

// checkTranslation jämför översättningen tr av en datafil med originalet.
func checkTranslation(orig, tr []byte) (translationReport, error) {
	var r translationReport
	o, err := decodeTree(orig)
	if err != nil {
		return r, err
	}
	var t any
	if tr != nil {
		if t, err = decodeTree(tr); err != nil {
			return r, err
		}
	}
	r.compare(o, t, "")
	return r, nil
}

// LangStatus är hur långt översättningen till ett språk har kommit, för
// -kontrollera.
type LangStatus struct {
	File        string
	Texts, Done int
}

// translationStatus räknar de översatta texterna i varje datafil.
func translationStatus(dir, lang string) ([]LangStatus, error) {
	var out []LangStatus
	for _, name := range dataFiles {
		orig, _, err := readDataFile(dir, name)
		if err != nil {
			return nil, err
		}
		tr, _, err := readTranslation(dir, lang, name)
		if err != nil {
			return nil, err
		}
		r, err := checkTranslation(orig, tr)
		if err != nil {
			return nil, fmt.Errorf("%s/%s: %w", lang, name, err)
		}
		out = append(out, LangStatus{File: name, Texts: r.texts, Done: r.done})
	}
	return out, nil
}

// langNames är språkens namn i utskrifterna i terminalen.
var langNames = map[string]string{langSv: "Svenska", langEn: "Engelska"}

// printTranslationStatus kontrollerar översättningen till lang och skriver
// hur många texter i varje fil som är översatta, för -kontrollera.
func printTranslationStatus(dir, lang string) error {
	if _, _, err := loadDataLang(dir, lang); err != nil {
		return err
	}
	status, err := translationStatus(dir, lang)
	if err != nil {
		return err
	}
	texts, done := 0, 0
	for _, s := range status {
		texts += s.Texts
		done += s.Done
	}
	fmt.Printf("%s: %d av %d texter är översatta. Resten visas på svenska.\n", langNames[lang], done, texts)
	for _, s := range status {
		fmt.Printf("  %-16s %5d av %5d\n", s.File, s.Done, s.Texts)
	}
	return nil
}
