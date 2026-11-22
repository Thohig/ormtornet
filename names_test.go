package main

import (
	"testing"
	"unicode/utf8"
)

func TestNamesFitInBothLanguages(t *testing.T) {
	for _, lang := range []string{langSv, "en"} {
		d, _, err := loadDataLang("data", lang)
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range d.Vandrarnamn.Fornamn {
			for _, s := range d.Vandrarnamn.Tillnamn {
				if n := utf8.RuneCountInString(f + " " + s); n > maxNameLen {
					t.Errorf("%s: %q är %d tecken, högst %d", lang, f+" "+s, n, maxNameLen)
				}
			}
		}
	}
}

func TestJoinGivesUniqueNames(t *testing.T) {
	g := testGame(t)
	seen := map[string]bool{}
	for i := 0; i < maxPlayers; i++ {
		c := testClient()
		g.Join(c, "")
		if c.player == nil {
			t.Fatalf("vandrare %d kom inte in", i+1)
		}
		f := firstName(c.player.Name)
		if seen[f] {
			t.Fatalf("förnamnet %q kom två gånger", f)
		}
		seen[f] = true
		if !contains(logs(c), g.data.Text("namn_givet", "namn", c.player.Name)) {
			t.Fatal("loggen sa inte vad tornet kallar vandraren")
		}
	}
}

func TestJoinWithTakenNameGetsANewOne(t *testing.T) {
	g := testGame(t)
	a, b := testClient(), testClient()
	g.Join(a, "")
	g.Join(b, a.player.Name)
	if b.player == nil || b.player == a.player || firstName(b.player.Name) == firstName(a.player.Name) {
		t.Fatal("den andra fliken fick inte ett eget namn")
	}
	// I lobbyn lämnar den som tappar kontakten sin plats, men får samma
	// namn igen när sidan kommer tillbaka.
	name := a.player.Name
	g.Leave(a)
	c := testClient()
	g.Join(c, name)
	if c.player == nil || c.player.Name != name {
		t.Fatal("namnet följde inte med tillbaka")
	}
}
