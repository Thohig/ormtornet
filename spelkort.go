package main

// Spelkorten: allt som händer en vandrare under en runda blir ett kort ur
// samma lek som ormkorten (docs/spelkort.md). Händelserna lägger sina kort
// här medan rundan avgörs, och rundans sammanfattning skickar dem med
// rundans rader. Klienten vänder upp det första och visar resten som
// minikort. Bilderna ritas av klienten (web/kort.js) efter motivet i bild.

import (
	"slices"
	"strings"
)

// CardData styr spelkorten (spelkort i data/texter.json).
type CardData struct {
	// Ordning är korten efter hur viktiga de är; det första vänds upp.
	Ordning []string `json:"ordning"`
}

// gameCard är ett spelkort. Slag är sorten, för ordningen (falla, varelse,
// foremal …), Bild motivet och Rubrik det som står överst. Fällor och
// varelser har sin korta beskrivning som rubrik, aldrig namnet.
type gameCard struct {
	Slag   string `json:"slag"`
	Bild   string `json:"bild"`
	Rubrik string `json:"rubrik"`
	Text   string `json:"text,omitempty"`
	// Svag: kortet visar något som inte hände, till exempel en fälla man
	// klev förbi. Bilden tonas ned.
	Svag bool `json:"svag,omitempty"`
	// Stor: kortet vänds upp stort mitt på kartan, som uppvaknandet.
	Stor bool `json:"stor,omitempty"`
}

// addCard lägger ett kort till vandrarens runda.
func (g *Game) addCard(p *Player, c gameCard) {
	if p == nil || p.client == nil {
		return // ingen skärm att visa kortet på
	}
	c.Rubrik = upperFirst(c.Rubrik)
	p.roundCards = append(p.roundCards, c)
}

// roundCardsOf ger rundans kort i ordningen ur data, det viktigaste först.
// Kort av en sort som ordningen inte känner kommer sist.
func (g *Game) roundCardsOf(p *Player) []gameCard {
	cards := slices.Clone(p.roundCards)
	rank := func(slag string) int {
		if i := slices.Index(g.data.Spelkort.Ordning, slag); i >= 0 {
			return i
		}
		return len(g.data.Spelkort.Ordning)
	}
	slices.SortStableFunc(cards, func(a, b gameCard) int { return rank(a.Slag) - rank(b.Slag) })
	return cards
}

// cardImages är motiven web/kort.js kan rita (BILDER i kort.js). Ett
// motiv i data/ som inte finns här ritas som facklan; kontrollen av
// datafilerna säger till (validate.go).
var cardImages = []string{
	"fackla", "sover", "ror_sig", "vaknar",
	"fallgaller", "rasande_golv", "roterande", "pilar", "dimma", "ormgrop", "skinn", "klocka", "morker", "eko",
}

// cardImage är motivet för innehåll i data/: fältet bild, annars id:t.
func cardImage(bild, id string) string {
	if strings.TrimSpace(bild) != "" {
		return bild
	}
	return id
}
