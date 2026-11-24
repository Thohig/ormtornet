package main

// Kontroll av datafilerna. Ett felstavat fält eller en hänvisning till
// något som inte finns ska synas när tornet startar, inte mitt i en vandring.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
)

// decodeStrict läser JSON och nekar fält som inte finns i strukturen.
// Fält som börjar med _ är kommentarer och hoppas över.
func decodeStrict(raw []byte, v any) error {
	var tree any
	if err := json.Unmarshal(raw, &tree); err != nil {
		return explainJSON(raw, err)
	}
	clean, err := json.Marshal(dropComments(tree))
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(clean))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return explainJSON(nil, err)
	}
	return nil
}

// dropComments tar bort alla fält som börjar med _.
func dropComments(v any) any {
	switch t := v.(type) {
	case map[string]any:
		for k, e := range t {
			if strings.HasPrefix(k, "_") {
				delete(t, k)
			} else {
				t[k] = dropComments(e)
			}
		}
	case []any:
		for i, e := range t {
			t[i] = dropComments(e)
		}
	}
	return v
}

// explainJSON gör JSON-felen begripliga, med rad där det går.
func explainJSON(raw []byte, err error) error {
	msg := err.Error()
	if field, ok := strings.CutPrefix(msg, "json: unknown field "); ok {
		return fmt.Errorf("okänt fält %s (felstavat, eller på fel ställe?)", field)
	}
	var syn *json.SyntaxError
	if errors.As(err, &syn) {
		return fmt.Errorf("rad %d: %s", lineAt(raw, syn.Offset), msg)
	}
	var typ *json.UnmarshalTypeError
	if errors.As(err, &typ) {
		return fmt.Errorf("fältet %s har fel sorts värde (%s, inte %s)", typ.Field, typ.Value, typ.Type)
	}
	return err
}

func lineAt(raw []byte, offset int64) int {
	if offset > int64(len(raw)) {
		offset = int64(len(raw))
	}
	return bytes.Count(raw[:offset], []byte("\n")) + 1
}

// validate letar efter hänvisningar som inte går ihop och värden utanför
// det tillåtna. Den ger en rad för varje problem.
// hexColor är en färg som #rrggbb.
var hexColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func (d *Data) validate() []string {
	var out []string
	bad := func(file, format string, args ...any) {
		out = append(out, file+": "+fmt.Sprintf(format, args...))
	}
	chance := func(file, name string, v float64) {
		if v < 0 || v > 1 {
			bad(file, "%s ska vara mellan 0 och 1, inte %v", name, v)
		}
	}
	unique := func(file, what string, ids []string) map[string]bool {
		seen := map[string]bool{}
		for _, id := range ids {
			if id == "" {
				bad(file, "%s utan id", what)
			} else if seen[id] {
				bad(file, "%s %q finns två gånger", what, id)
			}
			seen[id] = true
		}
		return seen
	}

	// cardImageOK: spelkortets motiv måste vara ett som web/kort.js kan rita.
	cardImageOK := func(file, what, bild string) {
		if !slices.Contains(cardImages, bild) {
			bad(file, "%s: bilden %q finns inte; spelkortet kan rita %s", what, bild, strings.Join(cardImages, ", "))
		}
	}

	var ids []string
	for _, r := range d.Rumstyper {
		ids = append(ids, r.ID)
		if r.Modell == nil {
			continue
		}
		for _, f := range r.Modell.Inredning {
			if !slices.Contains(furnitureNames, f) {
				bad("texter.json", "rumstypen %q: inredningen %q finns inte; kartan kan rita %s", r.ID, f, strings.Join(furnitureNames, ", "))
			}
		}
	}
	unique("texter.json", "rumstypen", ids)

	// Titlarna på raden om var och en: måttet måste finnas och texten också.
	for _, t := range d.Spelkort.Berattelse.Titlar {
		if titleMeasures[t.Matt] == nil {
			bad("texter.json", "berättelsens titel %q: måttet %q finns inte", t.Text, t.Matt)
		}
		if len(d.Texter[t.Text]) == 0 {
			bad("texter.json", "berättelsens titel: texten %q finns inte", t.Text)
		}
	}
	// Berättelsens höjdpunkter: motiven, rubrikerna, krönikans nycklar och
	// berättelsens text för varje nyckel måste finnas.
	for _, k := range d.Spelkort.Berattelse.Slag {
		cardImageOK("texter.json", "berättelsens slag "+k.ID, k.Bild)
		if len(d.Texter[k.Rubrik]) == 0 {
			bad("texter.json", "berättelsens slag %q: rubriken %q finns inte bland texterna", k.ID, k.Rubrik)
		}
		if k.Flera != "" && len(d.Texter[k.Flera]) == 0 {
			bad("texter.json", "berättelsens slag %q: texten %q finns inte", k.ID, k.Flera)
		}
		for _, n := range k.Nycklar {
			if len(d.Texter[n]) == 0 {
				bad("texter.json", "berättelsens slag %q: krönikans nyckel %q finns inte bland texterna", k.ID, n)
			}
			if len(d.Texter[storyKey(n)]) == 0 {
				bad("texter.json", "berättelsens slag %q: berättelsens text %q finns inte", k.ID, storyKey(n))
			}
		}
	}

	ids = nil
	for _, t := range d.Fallor {
		ids = append(ids, t.ID)
		cardImageOK("fallor.json", "fällan "+t.ID, cardImage(t.Bild, t.ID))
	}
	unique("fallor.json", "fällan", ids)
	chance("fallor.json", "handelsechans", d.Regler.Handelsechans)

	// Modellerna på kartan snett ovanifrån: en färg som #rrggbb och en
	// storlek som går att se.
	model := func(file, what, id string, m *ModelLook) {
		if m == nil {
			return
		}
		if m.Farg != "" && !hexColor.MatchString(m.Farg) {
			bad(file, "%s %q: modell.farg ska skrivas som #rrggbb, inte %q", what, id, m.Farg)
		}
		if m.Storlek != 0 && (m.Storlek < 0.5 || m.Storlek > 2) {
			bad(file, "%s %q: modell.storlek ska vara mellan 0,5 och 2, inte %v", what, id, m.Storlek)
		}
	}

	ids = nil
	for _, c := range d.Klasser {
		ids = append(ids, c.ID)
		model("klasser.json", "klassen", c.ID, c.Modell)
		cardImageOK("klasser.json", "klassen "+c.ID, cardImage(c.Bild, c.ID))
	}
	classes := unique("klasser.json", "klassen", ids)

	ids = nil
	for _, it := range d.Foremal.Lista {
		cardImageOK("foremal.json", "föremålet "+it.ID, cardImage(it.Bild, it.ID))
		ids = append(ids, it.ID)
		if it.Tecken == "" && !it.Brev {
			bad("foremal.json", "föremålet %q saknar tecken (det står framför namnet i väskan)", it.ID)
		}
	}
	unique("foremal.json", "föremålet", ids)

	ids = nil
	for _, c := range d.Ormen.Kort {
		ids = append(ids, c.ID)
	}
	unique("ormen.json", "ormkortet", ids)

	m := &d.Monster
	chance("monster.json", "chans", m.Chans)
	chance("monster.json", "vandra_chans", m.VandraChans)
	if r := m.Traff; r.Min < 0 || r.Max > 100 || r.Min > r.Max {
		bad("monster.json", "traff: min och max ska vara 0–100, min först")
	}
	ids = nil
	for _, t := range m.Typer {
		model("monster.json", "varelsen", t.ID, t.Modell)
		cardImageOK("monster.json", "varelsen "+t.ID, cardImage(t.Bild, t.ID))
		ids = append(ids, t.ID)
		if t.Vandrar != nil {
			chance("monster.json", t.ID+": vandrar", *t.Vandrar)
		}
		if t.Lynne < 0 || t.Lynne > 100 {
			bad("monster.json", "%s: lynne ska vara 0–100", t.ID)
		}
		if t.Liv <= 0 {
			bad("monster.json", "%s: liv ska vara minst 1", t.ID)
		}
	}
	unique("monster.json", "varelsen", ids)

	ids = nil
	for _, e := range d.Handelser.Lista {
		ids = append(ids, e.ID)
		cardImageOK("handelser.json", "händelsen "+e.ID, e.Bild)
		if e.Varelse != "" && d.monsterType(e.Varelse) == nil {
			bad("handelser.json", "%s: varelsen %q finns inte i monster.json", e.ID, e.Varelse)
		}
		if len(e.Text) == 0 && len(e.Har) == 0 {
			bad("handelser.json", "%s: text och har saknas", e.ID)
		}
	}
	unique("handelser.json", "händelsen", ids)
	chance("handelser.json", "chans", d.Handelser.Chans)

	for _, p := range d.AI.Personligheter {
		chance("ai.json", p.ID+": sok_chans", p.SokChans)
		chance("ai.json", p.ID+": stridslust", p.Stridslust)
		chance("ai.json", p.ID+": guldlust", p.Guldlust)
		for _, c := range p.Klasser {
			if len(d.Klasser) > 0 && !classes[c] {
				bad("ai.json", "%s: klassen %q finns inte i klasser.json", p.ID, c)
			}
		}
	}
	if len(d.AI.Personligheter) == 0 {
		bad("ai.json", "inga personligheter")
	}

	ids = nil
	for _, s := range d.Svarighet.Nivaer {
		ids = append(ids, s.ID)
		for name, v := range map[string]*float64{"fallchans": s.Fallchans, "varelsechans": s.Varelsechans, "skatt_forlust": s.SkattForlust} {
			if v != nil {
				chance("svarighet.json", s.ID+": "+name, *v)
			}
		}
	}
	a := &d.Avslojande
	if len(a.Namn) == 0 {
		bad("avslojande.json", "namn: inga namn")
	}
	for i, n := range a.Namn {
		if len(n) != 3 {
			bad("avslojande.json", "namn %d: ska ha tre delar, har %d", i+1, len(n))
		}
		for _, part := range n {
			if part == "" || strings.Trim(strings.ToLower(part), "abcdefghijklmnopqrstuvwxyz") != "" {
				bad("avslojande.json", "namn %d: delen %q ska ha bara bokstäver a–z", i+1, part)
			}
		}
	}
	if len(a.Platser) != 3 {
		bad("avslojande.json", "platser: ska ha tre listor")
	}
	for i, l := range a.Platser {
		if len(l) == 0 {
			bad("avslojande.json", "platser %d: inga texter", i+1)
		}
	}
	for _, src := range partSources {
		if len(a.Delar[src]) == 0 {
			bad("avslojande.json", "delar: %s saknar texter", src)
		}
		for _, s := range a.Delar[src] {
			if !strings.Contains(s, "{del}") || !strings.Contains(s, "{plats}") {
				bad("avslojande.json", "delar %s: %q saknar {del} eller {plats}", src, s)
			}
		}
	}
	chance("avslojande.json", "ristning_chans", a.RistningChans)
	for _, s := range a.SvarMedNamn {
		if !strings.Contains(s, "{namn}") {
			bad("avslojande.json", "svar_med_namn: %q saknar {namn}", s)
		}
	}
	chance("avslojande.json", "invigd.forst_guld", a.Invigd.ForstGuld)
	for _, t := range d.Fallor {
		if len(a.Varsel.Fallor[t.ID]) == 0 {
			bad("avslojande.json", "varsel: fällan %q saknar text", t.ID)
		}
	}
	for _, t := range d.Monster.Typer {
		if len(a.Varsel.Varelser[t.ID]) == 0 {
			bad("avslojande.json", "varsel: varelsen %q saknar text", t.ID)
		}
	}
	for lvl, w := range a.Forandring.Fallvikter {
		for id := range w {
			if d.trap(id) == nil {
				bad("avslojande.json", "forandring: fallvikter %s: fällan %q finns inte", lvl, id)
			}
		}
	}
	if a.Forandring.PortStangsVar < 1 {
		bad("avslojande.json", "forandring: port_stangs_var ska vara minst 1")
	}
	if len(a.Varsel.Tomt) == 0 || len(a.Varsel.Avstand) != 2 {
		bad("avslojande.json", "varsel: tomt ska ha texter och avstand två ord")
	}
	if len(a.EpilogNamn) == 0 {
		bad("avslojande.json", "epilog_namn: inga texter")
	}
	levels := unique("svarighet.json", "förvalet", ids)
	if d.Svarighet.Standard != "" && !levels[d.Svarighet.Standard] {
		bad("svarighet.json", "standard %q är inget av förvalen", d.Svarighet.Standard)
	}
	return out
}

// keys ger nycklarna i en tabell, för felmeddelanden.
func keys[V any](m map[string]V) string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return strings.Join(out, ", ")
}
