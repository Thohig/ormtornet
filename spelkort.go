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
	// Berattelse styr berättelsens höjdpunkter i slutet (berattelse.go).
	Berattelse StoryData `json:"berattelse"`
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
	// Forbannad: ett föremål som visar sin baksida. Kortet vänds en gång
	// till och blir mörkt.
	Forbannad bool `json:"forbannad,omitempty"`
	// Fallen: en varelse som vandraren fällde. Bilden är överkorsad.
	Fallen bool `json:"fallen,omitempty"`
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

// cardMsg är ett spelkort som kommer mellan rundorna, till exempel ett
// varsel eller viskningen. Klienten visar det som ett eget spelkort, efter
// det som står uppe.
type cardMsg struct {
	T     string   `json:"t"`
	Round int      `json:"round"`
	Kort  gameCard `json:"kort"`
}

// sendCard visar ett spelkort genast, utanför rundans sammanfattning.
func (g *Game) sendCard(p *Player, c gameCard) {
	if p == nil || p.client == nil {
		return
	}
	c.Rubrik = upperFirst(c.Rubrik)
	g.send(p.client, cardMsg{T: "kort", Round: g.round, Kort: c})
}

// bigMsg är ett stort kort mitt på kartan för alla som är kvar inne, som
// uppvaknandet: en port som stängs eller solen som går ner.
type bigMsg struct {
	T    string   `json:"t"`
	Kort gameCard `json:"kort"`
}

// bigCardInside visar kortet stort för alla som är kvar i tornet.
func (g *Game) bigCardInside(c gameCard) {
	c.Rubrik = upperFirst(c.Rubrik)
	c.Stor = true
	for _, p := range g.players {
		if p != nil && p.Status == statusIn {
			g.send(p.client, bigMsg{T: "storkort", Kort: c})
		}
	}
}

// endCard är vandringens sista kort: ute, död eller kvar när solen gick
// ner. Det står överst i Din vandring, och den som kommer ut ser det stort
// på kartan.
func (g *Game) endCard(p *Player, c gameCard) {
	c.Rubrik = upperFirst(c.Rubrik)
	p.lastCard = &c
	if p.Status == statusOut || p.Status == statusDead {
		g.addCard(p, c)
	}
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
	"krypare", "yngel", "fjallvakt", "gallgast", "benvaktare", "vandrare",
	"rep", "amulett", "veke", "skold", "vaska", "kompass", "rokpulver", "spegel", "kastkniv", "benlur",
	"flojt", "timglas", "svart_parla", "mynt_hal", "rostkrona", "brev", "dryck", "guld",
	"svard", "brynja", "tarst_svard",
	"andas", "spricka", "gluggar", "ormskinn", "ras", "doden", "klockslag",
	"kallare", "varsel", "namndel", "invigd",
	"ut", "dod", "sol", "port",
	"legoknekt", "tjuv", "ortkvinna", "kartritare", "ormtjusare", "sparfinnare", "lyktbarare", "pilgrim",
}

// cardImage är motivet för innehåll i data/: fältet bild, annars id:t.
func cardImage(bild, id string) string {
	if strings.TrimSpace(bild) != "" {
		return bild
	}
	return id
}
