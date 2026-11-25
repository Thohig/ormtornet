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

// stirSerpent gör ormen oroligare: ett vakna-kort till, på en slumpad plats
// i leken.
func (g *Game) stirSerpent() {
	for _, c := range g.data.Ormen.Kort {
		if c.Vaknar {
			g.deck = append(g.deck, c.ID)
			i, j := len(g.deck)-1, g.rng.IntN(len(g.deck))
			g.deck[i], g.deck[j] = g.deck[j], g.deck[i]
			return
		}
	}
}

// calmSerpent låter ormen sova djupare: ett sovkort till, på en slumpad
// plats i leken.
func (g *Game) calmSerpent() {
	for _, c := range g.data.Ormen.Kort {
		if !c.Vaknar {
			g.deck = append(g.deck, c.ID)
			i, j := len(g.deck)-1, g.rng.IntN(len(g.deck))
			g.deck[i], g.deck[j] = g.deck[j], g.deck[i]
			return
		}
	}
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
	// Ormens präst tar bara en näve per besök i kammaren.
	if g.classVal(p, "en_nave", 0) > 0 && p.tookVisit {
		g.logTo(p, g.data.Text("klass_prast_en_nave"), "dim")
		return
	}
	p.tookVisit = true
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
	// Ormens skinn: ormen känner lukten, och ibland väcker näven den.
	if p.has("skinn") != nil && g.curse(p, "skinn") && g.rng.Float64() < 0.25 {
		for j, id := range g.deck {
			if c := g.data.card(id); c != nil && c.Vaknar {
				i = j
				break
			}
		}
	}
	// Den invigde som har viskat till ormen: nästa näve väcker den inte,
	// om leken har något sovkort kvar. Leken är densamma för de andra.
	if p.safeDraw {
		p.safeDraw = false
		if c := g.data.card(g.deck[i]); c != nil && c.Vaknar {
			for j, id := range g.deck {
				if c := g.data.card(id); c != nil && !c.Vaknar {
					i = j
					break
				}
			}
		}
	}
	// Ormens präst: det första vakna-kortet hen drar blir ett sovkort.
	if g.is(p, "prast") && !p.priestCalm {
		if c := g.data.card(g.deck[i]); c != nil && c.Vaknar {
			for j, id := range g.deck {
				if c := g.data.card(id); c != nil && !c.Vaknar {
					i, p.priestCalm = j, true
					g.logTo(p, g.data.Text("klass_prast_lugnar"), "rost")
					break
				}
			}
		}
	}
	c := g.data.card(g.deck[i])
	if c == nil {
		g.deck = append(g.deck[:i], g.deck[i+1:]...)
		return
	}
	// Den förste som sträcker sig efter guldet nådde kammaren först, också
	// när näven väcker ormen.
	first := true
	for _, o := range g.players {
		if o != nil && len(o.stats.draws) > 0 {
			first = false
		}
	}
	if first {
		g.chronicle("kronika_kammaren", "namn", p.Name)
	}
	p.stats.draws = append(p.stats.draws, c.ID)
	g.learnPart(p, srcChamber) // ormen väser ett ord i sömnen
	if c.Vaknar {
		g.chronicle("kronika_vaknade", "namn", p.Name)
		g.logTo(p, pick(c.Texter), "varn")
		p.Seen["ormen"] = true
		g.wakeSerpent()
		g.explainNow(p, "ormen")
		return
	}
	g.deck = append(g.deck[:i], g.deck[i+1:]...)
	// Den svarta pärlan och ormens öga när det har vänt: varje näve lägger
	// ett vakna-kort i leken.
	if g.curse(p, "svart_parla") || g.eyeWakes(p) {
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
	g.goldNote(p, "kammaren", n, "kort", c.Namn)
	p.stats.takes++
	p.stats.takeGold += n
	g.soundTo(p, "skatt")
	text := fill(pick(c.Texter), "antal", strconv.Itoa(n))
	g.logTo(p, text, "syn")
	// Guldet och vad vandraren bär står längst ner på spelkortet.
	g.addCard(p, gameCard{Slag: "ormen", Bild: c.ID, Rubrik: c.Namn, Text: text})
	other := pick(c.Andra)
	if g.is(p, "tjuv") {
		if other != "" {
			g.logTo(p, g.data.Text("klass_tjuv_tyst"), "dim")
		}
		return
	}
	g.heardTaking(p, other)
	if other != "" {
		g.soundAll("ormen_ror")
	}
}

// heardTaking: när någon tar en näve hör alla andra i tornet åt vilket håll
// det kom ifrån, och de i kammaren ser vem det var (docs/ormens-plats.md).
// Rör sig ormen står det först. Den som är först i kammaren blir känd.
func (g *Game) heardTaking(p *Player, other string) {
	for _, o := range g.players {
		if o == nil || o == p || o.Status != statusIn {
			continue
		}
		var text string
		switch {
		case o.Cellar == 0 && o.X == p.X && o.Y == p.Y:
			text = g.data.Text("nave_ser", "namn", p.Name)
		default:
			text = g.data.Text("nave_hors", "horn", g.hornToward(o.X, o.Y, p.X, p.Y))
		}
		if other != "" {
			text = other + " " + text
		}
		o.heardSerpent(p.X, p.Y)
		g.hearTo(o, text, "info")
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
		g.goldNote(p, "ormen", -lost)
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
	g.moveChamber()
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

// moveChamber: när ormen har vaknat kryper den till en ny kammare, med
// samma regler som när den placerades (chamberSpots). Den gamla blir ett vanligt rum med lite guld
// kvar på golvet, och alla i tornet hör att ormen flyttar sig.
func (g *Game) moveChamber() {
	ox, oy := g.chamber()
	spots := g.chamberSpots(ox, oy)
	if len(spots) == 0 || (len(spots) == 1 && spots[0] == [2]int{ox, oy}) {
		return
	}
	n := spots[g.rng.IntN(len(spots))]

	g.chamberMoved = true // kammaren får de dörrar rummet har, inte fyra
	for _, p := range g.players {
		if p != nil && p.bot != nil {
			p.bot.hasHint = false // det de hört gäller den gamla platsen
		}
	}
	// Den gamla kammaren blir ett vanligt rum. Ett rum som ingen har sett
	// byggs först, med vanliga dörrar.
	old := &g.rooms[oy][ox]
	if old.Revealed {
		old.Kind, old.TypeID = "rum", g.roomTypeFor(ox, oy)
		old.Searched = false
		o := &g.data.Ormen
		g.spillGold(ox, oy, o.SkattMin+g.rng.IntN(max(1, o.SkattMax-o.SkattMin+1)), "ormen_guld_kvar")
	}
	g.chamberX, g.chamberY = n[0], n[1]
	r := g.build(n[0], n[1], false)
	r.Kind, r.TypeID = "kammare", "skattkammare"
	r.Event = ""
	if m := g.monsterAt(n[0], n[1]); m != nil {
		g.removeMonster(m) // varelsen flyr för ormen
	}
	for _, p := range g.players {
		if p != nil && p.Status == statusIn && p.Cellar == 0 {
			g.logTo(p, g.data.Text("ormen_flyttar"), "varn")
		}
	}
}

// breathThroughWall: den som står bredvid kammaren hör ormen andas genom
// väggen, och åt vilket håll, en gång per rum (docs/ormens-plats.md).
func (g *Game) breathThroughWall() {
	cx, cy := g.chamber()
	for _, p := range g.players {
		if p == nil || p.Status != statusIn || p.Cellar > 0 {
			continue
		}
		d, ok := g.nextTo(p, cx, cy)
		here := [2]int{p.X + 1, p.Y + 1}
		if !ok || p.breathAt == here {
			continue
		}
		p.breathAt = here
		p.heardSerpent(cx, cy)
		g.hearTo(p, g.data.Text("andetag_vagg", "riktning", g.data.Riktningar[d]), "info")
	}
}

// breathForAll: från nivå 2 hör alla inne ormens andetag en gång, åt vilket
// håll de kommer.
func (g *Game) breathForAll() {
	cx, cy := g.chamber()
	for _, p := range g.players {
		if p == nil || p.Status != statusIn || p.Cellar > 0 || (p.X == cx && p.Y == cy) {
			continue
		}
		p.heardSerpent(cx, cy)
		g.logTo(p, g.data.Text("andetag_alla", "horn", g.hornToward(p.X, p.Y, cx, cy)), "info")
	}
}
