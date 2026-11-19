package main

// Kontroll av datafilerna. Ett felstavat fält eller en hänvisning till
// något som inte finns ska synas när tornet startar, inte mitt i en vandring.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	var ids []string
	for _, r := range d.Rumstyper {
		ids = append(ids, r.ID)
	}
	unique("texter.json", "rumstypen", ids)

	ids = nil
	for _, t := range d.Fallor {
		ids = append(ids, t.ID)
	}
	unique("fallor.json", "fällan", ids)
	chance("fallor.json", "handelsechans", d.Regler.Handelsechans)

	ids = nil
	for _, c := range d.Klasser {
		ids = append(ids, c.ID)
	}
	classes := unique("klasser.json", "klassen", ids)

	ids = nil
	for _, it := range d.Foremal.Lista {
		ids = append(ids, it.ID)
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
	ids = nil
	for _, t := range m.Typer {
		ids = append(ids, t.ID)
		if t.Vandrar != nil {
			chance("monster.json", t.ID+": vandrar", *t.Vandrar)
		}
		for k := range t.Vanor {
			if _, ok := m.Val[k]; !ok {
				bad("monster.json", "%s: vanan %q är inget val (%s)", t.ID, k, keys(m.Val))
			}
		}
		if t.Liv <= 0 {
			bad("monster.json", "%s: liv ska vara minst 1", t.ID)
		}
	}
	unique("monster.json", "varelsen", ids)

	for _, p := range d.AI.Personligheter {
		chance("ai.json", p.ID+": sok_chans", p.SokChans)
		chance("ai.json", p.ID+": stridslust", p.Stridslust)
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
	if len(a.Ordning) != 3 {
		bad("avslojande.json", "ordning: ska ha tre ord")
	}
	for _, lvl := range []string{"2", "3", "4"} {
		if len(a.Delar[lvl]) == 0 {
			bad("avslojande.json", "delar: nivå %s saknar texter", lvl)
		}
		for _, s := range a.Delar[lvl] {
			if !strings.Contains(s, "{del}") {
				bad("avslojande.json", "delar %s: %q saknar {del}", lvl, s)
			}
		}
	}
	for _, s := range a.SvarMedNamn {
		if !strings.Contains(s, "{namn}") {
			bad("avslojande.json", "svar_med_namn: %q saknar {namn}", s)
		}
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
