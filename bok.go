package main

// Vandrarens bok (docs/vandrarens-bok.md): det spelarsidan behöver för
// bokens flikar, på tornets språk. Grundernas texter står bland sidans egna
// texter (bok_* i data/spelarsidan.json); här står det som kommer ur
// innehållet i data/. Inga namn på rum, fällor eller varelser och inget om
// tornets hemlighet.

type bookView struct {
	Sprak   string       `json:"sprak"`
	Klasser []*classView `json:"klasser"`
}

// Book ger bokens innehåll.
func (g *Game) Book() bookView {
	g.mu.Lock()
	defer g.mu.Unlock()
	v := bookView{Sprak: g.base.Lang, Klasser: []*classView{}}
	for i := range g.data.Klasser {
		v.Klasser = append(v.Klasser, viewOfClass(&g.data.Klasser[i]))
	}
	return v
}
