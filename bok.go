package main

import (
	"slices"
)

// Vandrarens bok (docs/vandrarens-bok.md): det spelarsidan behöver för
// bokens flikar, på tornets språk. Grundernas texter står bland sidans egna
// texter (bok_* i data/spelarsidan.json); här står det som kommer ur
// innehållet i data/. Inga namn på rum, fällor eller varelser och inget om
// tornets hemlighet.
//
// Boken fylls i med det vandraren har mött. Tornet säger det i
// tillståndet (bok, som "varelse:krypare"), och spelarsidan sparar det i
// webbläsaren, så att boken växer från vandring till vandring.

type bookView struct {
	Sprak    string       `json:"sprak"`
	Klasser  []*classView `json:"klasser"`
	Foremal  []bookCard   `json:"foremal"`
	Varelser []bookCard   `json:"varelser"`
	Fallor   []bookCard   `json:"fallor"`
}

// bookCard är ett kort i boken. Spelarsidan visar det bara när vandraren
// har mött det; forbannelse och vander bara när hen har varit med om det.
type bookCard struct {
	ID          string `json:"id"`
	Bild        string `json:"bild"`
	Rubrik      string `json:"rubrik"`
	Text        string `json:"text,omitempty"`
	Om          string `json:"om,omitempty"`   // vad den gör eller hur den är, i ord
	Spar        string `json:"spar,omitempty"` // spåret efter en fälla
	Forbannelse string `json:"forbannelse,omitempty"`
	Vander      string `json:"vander,omitempty"`
}

func first(list []string) string {
	if len(list) == 0 {
		return ""
	}
	return list[0]
}

// Book ger bokens innehåll.
func (g *Game) Book() bookView {
	g.mu.Lock()
	defer g.mu.Unlock()
	v := bookView{Sprak: g.base.Lang, Klasser: []*classView{}, Foremal: []bookCard{}, Varelser: []bookCard{}, Fallor: []bookCard{}}
	for i := range g.data.Klasser {
		v.Klasser = append(v.Klasser, viewOfClass(&g.data.Klasser[i]))
	}
	for _, t := range g.data.Foremal.Lista {
		if len(t.Av) > 0 {
			continue // de hopfogade kommer i ett senare steg
		}
		c := bookCard{ID: t.ID, Bild: cardImage(t.Bild, t.ID), Rubrik: upperFirst(t.Namn), Text: t.Beskrivning, Vander: first(t.Vander)}
		if t.Forbannad {
			c.Forbannelse = first(t.Forbannelse)
		}
		v.Foremal = append(v.Foremal, c)
	}
	for _, t := range g.data.Monster.Typer {
		v.Varelser = append(v.Varelser, bookCard{ID: t.ID, Bild: cardImage(t.Bild, t.ID), Rubrik: upperFirst(t.Kort), Text: first(t.DykerUpp), Om: t.Bok})
	}
	for _, t := range g.data.Fallor {
		v.Fallor = append(v.Fallor, bookCard{ID: t.ID, Bild: cardImage(t.Bild, t.ID), Rubrik: upperFirst(t.Kort), Om: t.Bok, Spar: first(t.SparEfter)})
	}
	return v
}

// learn skriver in i vandrarens bok att hen har mött något: slag är
// foremal, varelse, falla, forbannelse eller vander.
func (g *Game) learn(p *Player, slag, id string) {
	if p == nil || id == "" {
		return
	}
	if p.bok == nil {
		p.bok = map[string]bool{}
	}
	p.bok[slag+":"+id] = true
}

// learned är det vandraren har mött, i ordning.
func (p *Player) learned() []string {
	var out []string
	for k := range p.bok {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}
