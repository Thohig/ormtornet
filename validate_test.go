package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// dataWith kopierar datamappen och ändrar en fil med change.
func dataWith(t *testing.T, file string, change func(string) string) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range dataFiles {
		raw, err := os.ReadFile(filepath.Join("data", name))
		if err != nil {
			t.Fatal(err)
		}
		s := string(raw)
		if name == file {
			s = change(s)
		}
		if err := os.WriteFile(filepath.Join(dir, name), []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestDataFilesAreValid(t *testing.T) {
	if _, _, err := loadData("data"); err != nil {
		t.Fatal(err)
	}
}

func TestDataCheckFindsMistakes(t *testing.T) {
	cases := []struct {
		file, old, new, want string
	}{
		{"monster.json", `"vandra_chans"`, `"vandar_chans"`, `okänt fält "vandar_chans"`},
		{"monster.json", `"vandra_chans": 0.25`, `"vandra_chans": 2.5`, "vandra_chans ska vara mellan 0 och 1"},
		{"monster.json", `"vanor": { "anfall": 1`, `"vanor": { "anfal": 1`, `vanan "anfal"`},
		{"monster.json", `"id": "yngel"`, `"id": "krypare"`, `varelsen "krypare" finns två gånger`},
		{"ai.json", `"legoknekt"`, `"legoknekten"`, `klassen "legoknekten" finns inte`},
		{"monster.json", `"liv": 2,`, `"liv": "två",`, "har fel sorts värde"},
		{"monster.json", `"chans": 0.2,`, `"chans": 0.2`, "rad 4"},
		{"monster.json", `"farg": "#d8d0bf"`, `"farg": "blek"`, `modell.farg ska skrivas som #rrggbb`},
		{"klasser.json", `"farg": "#7a5a32", "storlek": 1`, `"farg": "#7a5a32", "storlek": 5`, `modell.storlek ska vara mellan 0,5 och 2`},
		{"foremal.json", `"tecken": "➰",`, ``, `föremålet "rep" saknar tecken`},
	}
	for _, c := range cases {
		dir := dataWith(t, c.file, func(s string) string {
			if !strings.Contains(s, c.old) {
				t.Fatalf("%s: hittar inte %s", c.file, c.old)
			}
			return strings.Replace(s, c.old, c.new, 1)
		})
		_, _, err := loadData(dir)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s med %s: fick %v, ville ha %q", c.file, c.new, err, c.want)
		}
	}
}

// Varje textnyckel som koden frågar efter ska finnas i data/. En nyckel som
// saknas syns annars som [nyckel] mitt i spelet.
func TestTextKeysExist(t *testing.T) {
	d, _, err := loadData("data")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	files, _ := filepath.Glob("*.go")
	checked := 0
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, f, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			lit, ok := call.Args[0].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true // sammansatta nycklar kontrolleras inte här
			}
			key, _ := strconv.Unquote(lit.Value)
			var have bool
			switch sel.Sel.Name {
			case "Text":
				have = len(d.Texter[key]) > 0
			case "MonsterText":
				have = len(d.Monster.Texter[key]) > 0
			default:
				return true
			}
			checked++
			if !have {
				t.Errorf("%s: texten %q saknas i data/", fset.Position(lit.Pos()), key)
			}
			return true
		})
	}
	if checked < 50 {
		t.Fatalf("hittade bara %d textnycklar, letar testet rätt?", checked)
	}
}
