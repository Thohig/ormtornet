package main

// Skattkammaren och ormen. Varje gång någon tar skatt dras ett ormkort.
// Sovkort tas bort ur leken; vakna-korten blir kvar. Ju färre kort som är
// kvar, desto större är risken att ormen vaknar.

import (
	"strconv"
)

// wakeMsg säger till alla att ormen har vaknat, med kortet som vänds upp
// för alla: namnet, det var och en märker, och för den som stod i kammaren
// vad den förlorade.
type wakeMsg struct {
	T       string `json:"t"`
	Namn    string `json:"namn"`
	Text    string `json:"text"`
	Inne    bool   `json:"inne,omitempty"` // stod i kammaren och kastades ut
	Forlust int    `json:"forlust,omitempty"`
	Skada   int    `json:"skada,omitempty"`
}

// cardMsg är ormkortet vandraren drog, som vänds upp på kartan.
type cardMsg struct {
	T      string `json:"t"`
	ID     string `json:"kort"`
	Namn   string `json:"namn"`
	Text   string `json:"text"`
	Guld   int    `json:"guld"`
	Vaknar bool   `json:"vaknar,omitempty"`
}

// buildDeck blandar en ny lek. Varje gång ormen har vaknat blir det färre
// sovkort, så att den vaknar lättare nästa gång.
func (g *Game) buildDeck() {
	o := &g.data.Ormen
	g.deck = g.deck[:0]
	sleepTotal := 0
	for _, c := range o.Kort {
		if !c.Vaknar {
			sleepTotal += c.Antal
		}
	}
	// Ta bort sovkort ett i taget från den sort som har flest kvar.
	counts := make([]int, len(o.Kort))
	for i, c := range o.Kort {
		counts[i] = c.Antal
	}
	remove := min(g.wakes*o.FarreEfterVakning, max(0, sleepTotal-o.MinstSovkort))
	for range remove {
		most := -1
		for i, c := range o.Kort {
			if !c.Vaknar && counts[i] > 0 && (most < 0 || counts[i] > counts[most]) {
				most = i
			}
		}
		if most < 0 {
			break
		}
		counts[most]--
	}
	// Tornets nåd håller i hela vandringen: färre vakna-kort också i en ny lek.
	for left, i := g.mercy, 0; left > 0 && i < len(o.Kort); i++ {
		if o.Kort[i].Vaknar {
			n := min(left, counts[i])
			counts[i] -= n
			left -= n
		}
	}
	// Fler vandrare tömmer leken fortare: lägg till sovkort för var och en
	// över fyra, så att risken växer i samma takt som med fyra.
	if extra := o.SovkortPerVandrare * max(0, g.vandrare()-4); extra > 0 {
		for i, c := range o.Kort {
			if !c.Vaknar {
				counts[i] += extra
				break
			}
		}
	}
	g.deckSleep = 0
	for i, c := range o.Kort {
		if !c.Vaknar {
			g.deckSleep += counts[i]
		}
		for range counts[i] {
			g.deck = append(g.deck, c.ID)
		}
	}
	g.rng.Shuffle(len(g.deck), func(i, j int) { g.deck[i], g.deck[j] = g.deck[j], g.deck[i] })
}

// sleepLevel säger hur djupt ormen sover: 0 djupt, 1 oroligt, 2 nära att vakna.
func (g *Game) sleepLevel() int {
	sleep := 0
	for _, id := range g.deck {
		if c := g.data.card(id); c != nil && !c.Vaknar {
			sleep++
		}
	}
	switch {
	case g.deckSleep == 0:
		return 2
	case sleep*3 > g.deckSleep*2:
		return 0
	case sleep*3 > g.deckSleep:
		return 1
	}
	return 2
}

// wakeRisk är chansen i procent att ormen vaknar vid nästa näve. En tom
// lek blandas om vid nästa näve, så då gäller en ny lek.
func (g *Game) wakeRisk() int {
	deck := g.deck
	if len(deck) == 0 {
		saved, savedSleep := append([]string(nil), g.deck...), g.deckSleep
		g.buildDeck()
		deck, g.deck, g.deckSleep = g.deck, saved, savedSleep
	}
	wake := 0
	for _, id := range deck {
		if c := g.data.card(id); c != nil && c.Vaknar {
			wake++
		}
	}
	if len(deck) == 0 {
		return 0
	}
	return (wake*100 + len(deck)/2) / len(deck)
}

// take låter spelaren ta en näve skatt och drar ett ormkort.
func (g *Game) take(p *Player) {
	if g.rooms[p.Y][p.X].Kind != "kammare" {
		return
	}
	if len(g.deck) == 0 {
		g.buildDeck()
	}
	i := g.rng.IntN(len(g.deck))
	i = g.flute(p, i)
	// Ormen känner den som vet: ibland drar den invigde ett vakna-kort.
	if p.initiated && g.rng.Float64() < g.data.Avslojande.Invigd.Vacker {
		for j, id := range g.deck {
			if c := g.data.card(id); c != nil && c.Vaknar {
				i = j
				break
			}
		}
	}
	c := g.data.card(g.deck[i])
	if c == nil {
		g.deck = append(g.deck[:i], g.deck[i+1:]...)
		return
	}
	p.stats.draws = append(p.stats.draws, c.ID)
	if c.Vaknar {
		g.chronicle("kronika_vaknade", "namn", p.Name)
		g.logTo(p, pick(c.Texter), "varn")
		p.Seen["ormen"] = true
		g.wakeSerpent()
		g.explainNow(p, "ormen")
		return
	}
	g.deck = append(g.deck[:i], g.deck[i+1:]...)
	// Den svarta pärlan: varje näve lägger ett vakna-kort i leken.
	if g.curse(p, "svart_parla") {
		for _, c := range g.data.Ormen.Kort {
			if c.Vaknar {
				g.deck = append(g.deck, c.ID)
				break
			}
		}
	}

	o := &g.data.Ormen
	n := o.SkattMin + g.rng.IntN(max(1, o.SkattMax-o.SkattMin+1))
	n = int(float64(n) * g.classVal(p, "kammar_guld", 1))
	p.Treasure += n
	p.stats.takes++
	p.stats.takeGold += n
	g.soundTo(p, "skatt")
	text := fill(pick(c.Texter), "antal", strconv.Itoa(n))
	g.logTo(p, text, "syn")
	g.send(p.client, cardMsg{T: "ormkort", ID: c.ID, Namn: c.Namn, Text: text, Guld: n})
	g.logTo(p, g.data.Text("skatt_totalt", "skatt", strconv.Itoa(p.Treasure)), "dim")
	if other := pick(c.Andra); other != "" {
		if g.is(p, "tjuv") {
			g.logTo(p, g.data.Text("klass_tjuv_tyst"), "dim")
		} else {
			g.hearOthers(p, other, "info")
			g.soundAll("ormen_ror")
		}
	}
}

// wakeSerpent: alla i kammaren förlorar skatt och liv och kastas ut.
// Väsandet hörs i hela tornet. Sedan somnar ormen om, lättare väckt än förut.
func (g *Game) wakeSerpent() {
	o := &g.data.Ormen
	g.soundAll("vasning")
	name := ""
	for _, c := range o.Kort {
		if c.Vaknar {
			name = c.Namn
			break
		}
	}
	// Var och en får kortet när det är klart vad som hände: alla i tornet
	// det de märker, de i kammaren vad de förlorade.
	msgs := map[*Player]wakeMsg{}
	general := pick(o.VaknarAndra)
	var inside []*Player
	for _, p := range g.players {
		if p == nil || p.Status != statusIn || p.Cellar > 0 {
			continue
		}
		if g.rooms[p.Y][p.X].Kind == "kammare" {
			inside = append(inside, p)
		} else {
			text := pick(o.VaknarAndra)
			g.logTo(p, text, "varn")
			msgs[p] = wakeMsg{T: "vaknar", Namn: name, Text: text}
		}
	}
	for _, p := range inside {
		share := min(1, o.SkattForlust+g.classVal(p, "extra_forlust", 0))
		lost := g.hiddenShare(p, int(float64(p.Treasure)*share+0.5))
		p.Treasure -= lost
		dmg := max(0, o.Skada-int(g.classVal(p, "ormskada", 0)))
		text := fill(pick(o.KastasUt), "forlust", strconv.Itoa(lost), "skada", strconv.Itoa(dmg))
		g.logTo(p, text, "varn")
		msgs[p] = wakeMsg{T: "vaknar", Namn: name, Text: text, Inne: true, Forlust: lost, Skada: dmg}
		p.hurt("ormen", "")
		g.damage(p, dmg, true)
		if p.Status != statusIn {
			continue
		}
		var doors []Dir
		for d := North; d <= West; d++ {
			if g.rooms[p.Y][p.X].Doors[d] {
				doors = append(doors, d)
			}
		}
		if len(doors) > 0 {
			p.Stuck = 0
			g.move(p, doors[g.rng.IntN(len(doors))], 2)
		}
	}
	g.wakes++
	g.buildDeck()
	g.logAll(pick(o.SomnarOm), "dim")
	for _, p := range g.players {
		if p == nil {
			continue
		}
		m, ok := msgs[p]
		if !ok {
			m = wakeMsg{T: "vaknar", Namn: name, Text: general} // ute, död eller i källarvalven
		}
		g.send(p.client, m)
	}
}
