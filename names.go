package main

import (
	"strings"
	"unicode/utf8"
)

// Namnen: vandrarna skriver inte sina namn, tornet ger dem ett förnamn och
// ett tillnamn ur namn.json. Två vandrare får aldrig samma förnamn, så att
// det går att skilja dem åt i loggen.

// NameData är listorna i namn.json.
type NameData struct {
	Fornamn  []string `json:"fornamn"`
	Tillnamn []string `json:"tillnamn"`
}

// firstName är det första ordet i ett namn, i små bokstäver.
func firstName(name string) string {
	f, _, _ := strings.Cut(name, " ")
	return strings.ToLower(f)
}

// randomName lottar ett namn vars förnamn ingen annan vandrare har.
func (g *Game) randomName() string {
	used := map[string]bool{}
	for _, p := range g.players {
		if p != nil {
			used[firstName(p.Name)] = true
		}
	}
	n := g.data.Vandrarnamn
	for _, i := range g.rng.Perm(len(n.Fornamn)) {
		f := n.Fornamn[i]
		if used[strings.ToLower(f)] {
			continue
		}
		for _, j := range g.rng.Perm(len(n.Tillnamn)) {
			if name := f + " " + n.Tillnamn[j]; utf8.RuneCountInString(name) <= maxNameLen {
				return name
			}
		}
		return f
	}
	return n.Fornamn[g.rng.IntN(len(n.Fornamn))] // fler vandrare än förnamn: händer inte med 8
}
