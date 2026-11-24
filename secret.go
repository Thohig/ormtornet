package main

// Den som vet tornets hemlighet kan viska den till tornet. Stämmer viskningen
// med något av svaren i data/avslojande.json blir vandraren invigd och får
// en fördel. Tornet säger aldrig vad hemligheten är: den som gissar rätt får
// bara veta att tornet har hört.

import (
	"strconv"
	"strings"
	"unicode"
)

// normalizeWhisper gör om en viskning till små bokstäver och ord åtskilda av
// ett mellanslag, så att skiljetecken och versaler inte spelar roll.
func normalizeWhisper(s string) string {
	var b strings.Builder
	space := true
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) {
			b.WriteRune(r)
			space = false
		} else if !space {
			b.WriteByte(' ')
			space = true
		}
	}
	return strings.TrimSpace(b.String())
}

// whisperMatches är sant när viskningen innehåller något av svaren som hela ord.
func whisperMatches(whisper string, answers []string) bool {
	w := " " + normalizeWhisper(whisper) + " "
	for _, a := range answers {
		if a = normalizeWhisper(a); a != "" && strings.Contains(w, " "+a+" ") {
			return true
		}
	}
	return false
}

// Invigningen: de fasta svaren räcker inte. Varje vandring har ett eget
// namn i tre delar, ormens namn, som var och en hittar i tornet
// (namndelar.go), och viskningen ska innehålla namnet. Saknar datafilerna namn räcker
// svaren.

// usesName är sant när invigningen kräver ormens namn i den här vandringen.
func (g *Game) usesName() bool {
	return len(g.serpentName) == 3
}

// pickSerpentName lottar ormens namn när en vandring börjar.
func (g *Game) pickSerpentName() {
	g.serpentName = nil
	names := g.data.Avslojande.Namn
	if len(names) == 0 {
		return
	}
	g.serpentName = names[g.rng.IntN(len(names))]
	g.placeParts()
}

// serpentNameText är namnet som det skrivs, till exempel Karumeth.
func (g *Game) serpentNameText() string {
	s := strings.ToLower(strings.Join(g.serpentName, ""))
	if s == "" {
		return ""
	}
	r := []rune(s)
	return strings.ToUpper(string(r[0])) + string(r[1:])
}

func (p *Player) partCount() int {
	n := 0
	for _, k := range p.parts {
		if k {
			n++
		}
	}
	return n
}

// namePartsFor är delarna vandraren har, med versaler, i den ordning hen
// hittade dem och tomma för dem hen saknar. Platsen i tecknen säger alltså
// inget om var delen står i namnet. Ingenting om vandraren inte har
// hittat någon del.
func (g *Game) namePartsFor(p *Player) []string {
	if !g.usesName() || p.partCount() == 0 {
		return nil
	}
	out := make([]string, 3)
	for n, i := range p.partsFound {
		if n < 3 {
			out[n] = strings.ToUpper(g.serpentName[i])
		}
	}
	return out
}

// gainPart ger vandraren del i av namnet. Falskt om hen redan hade den.
func (p *Player) gainPart(i int) bool {
	if p.parts[i] {
		return false
	}
	p.parts[i] = true
	p.partsFound = append(p.partsFound, i)
	return true
}

// whisperRight avgör om viskningen gör vandraren invigd: ett av de fasta
// svaren, eller med ormens namn namnet tillsammans med ett svar eller en
// mall som "tornet är {namn}". Namnet får skrivas ihop eller med
// mellanrum och bindestreck mellan delarna.
func (g *Game) whisperRight(text string) bool {
	a := &g.data.Avslojande
	if !g.usesName() {
		return whisperMatches(text, a.Svar)
	}
	name := normalizeWhisper(strings.Join(g.serpentName, ""))
	spaced := normalizeWhisper(strings.Join(g.serpentName, " "))
	w := " " + normalizeWhisper(text) + " "
	w = strings.ReplaceAll(w, " "+spaced+" ", " "+name+" ")
	if !strings.Contains(w, " "+name+" ") {
		return false
	}
	if whisperMatches(w, a.Svar) {
		return true
	}
	forms := make([]string, len(a.SvarMedNamn))
	for i, f := range a.SvarMedNamn {
		forms[i] = strings.ReplaceAll(f, "{namn}", name)
	}
	return whisperMatches(w, forms)
}

// wrongOrder är sant när viskningen har namnets tre delar, men i en annan
// ordning än namnets, ihop eller med mellanrum.
func (g *Game) wrongOrder(text string) bool {
	if !g.usesName() {
		return false
	}
	w := " " + normalizeWhisper(text) + " "
	for _, perm := range [][3]int{{0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}} {
		parts := []string{g.serpentName[perm[0]], g.serpentName[perm[1]], g.serpentName[perm[2]]}
		if strings.Join(parts, "") == strings.Join(g.serpentName, "") {
			continue // två delar är lika: samma namn
		}
		for _, sep := range []string{"", " "} {
			if strings.Contains(w, " "+normalizeWhisper(strings.Join(parts, sep))+" ") {
				return true
			}
		}
	}
	return false
}

// onlyName är sant när viskningen är ormens namn och inget mer. Det är rätt
// namn men ingen mening om det, och då förbrukas inte försöket.
func (g *Game) onlyName(text string) bool {
	if !g.usesName() {
		return false
	}
	w := normalizeWhisper(text)
	return w == normalizeWhisper(strings.Join(g.serpentName, "")) ||
		w == normalizeWhisper(strings.Join(g.serpentName, " "))
}

// Whisper tar emot en viskning från en vandrare.
func (g *Game) Whisper(c *Client, text string) {
	p := c.player
	if p == nil || p.bot != nil || g.phase != phasePlay || p.Status != statusIn {
		return
	}
	if p.initiated {
		// I kammaren kan den invigde lugna ormen, en gång per vandring.
		switch {
		case g.rooms[p.Y][p.X].Kind != "kammare":
			g.answerTo(p, g.data.Text("viskning_redan"), "dim")
		case p.calmed:
			g.answerTo(p, g.data.Text("invigd_lugnat"), "dim")
		default:
			p.calmed = true
			p.safeDraw = true
			g.answerTo(p, g.data.Text("invigd_lugnar"), "rost")
		}
		return
	}
	if len([]rune(text)) > 80 {
		text = string([]rune(text)[:80])
	}
	if strings.TrimSpace(text) == "" {
		return
	}
	minLevel := max(1, g.data.Avslojande.ViskaFranNiva)
	if g.level < minLevel {
		g.answerTo(p, g.data.Text("viskning_tidigt"), "dim")
		return
	}
	if p.guessLevel >= g.level {
		g.answerTo(p, g.data.Text("viskning_vanta"), "dim")
		return
	}
	if g.onlyName(text) {
		// Rätt namn men ingenting om det: tornet känner igen namnet, och
		// vandraren får försöka igen.
		g.answerTo(p, g.data.Text("viskning_bara_namn"), "rost")
		return
	}
	p.guessLevel = g.level
	switch {
	case g.whisperRight(text):
	case g.wrongOrder(text):
		// Rätt tecken i fel ordning: tornet säger det, men försöket är
		// förbrukat, så att ingen kan pröva sig igenom ordningarna.
		g.answerTo(p, g.data.Text("viskning_ordning"), "rost")
		return
	default:
		// En fel viskning gör ormen oroligare: ett vakna-kort till.
		g.stirSerpent()
		g.answerTo(p, g.data.Text("viskning_fel"), "rost")
		return
	}
	p.initiated = true
	g.soundTo(p, "rost")
	said := g.data.Text("invigd")
	g.answerTo(p, said, "rost")
	g.sendCard(p, gameCard{Slag: "viska", Bild: "invigd", Rubrik: g.data.Text("kort_invigd"), Text: said})
	g.logOthers(p, g.data.Text("invigd_andra"), "dim")
	if g.firstInitiated == 0 {
		// Den som hinner först får tornets gunst: mer guld när hen kommer ut.
		g.firstInitiated = p.ID
		g.logTo(p, g.data.Text("invigd_forst"), "rost")
		g.chronicle("kronika_invigd_forst", "namn", p.Name, "runda", strconv.Itoa(g.round))
	} else {
		g.chronicle("kronika_invigd", "namn", p.Name, "runda", strconv.Itoa(g.round))
	}
}

// takeNameParts: den som dödar en vandrare hittar de delar av namnet som
// offret hade och dödaren saknade.
func (g *Game) takeNameParts(killer, victim *Player) {
	if !g.usesName() || killer.bot != nil {
		return
	}
	var got []string
	for i := range 3 {
		if victim.parts[i] && killer.gainPart(i) {
			// Stenen visar var på ormen tecknet sitter, som där offret fann det.
			got = append(got, strings.ToUpper(g.serpentName[i])+" "+g.data.Text("namndel_plats", "plats", pick(g.data.Avslojande.Platser[i])))
		}
	}
	if len(got) == 0 {
		return
	}
	text := g.data.Text("namndel_byte", "offer", victim.Name, "delar", strings.Join(got, ", "))
	g.logTo(killer, text, "rost")
	g.addCard(killer, gameCard{Slag: "namn", Bild: "namndel", Rubrik: g.data.Text("kort_namndel"), Text: text})
	g.chronicle("kronika_namndel", "namn", killer.Name, "offer", victim.Name)
}

func (g *Game) anyGracedInside() bool {
	for _, p := range g.players {
		if p != nil && p.Status == statusIn && g.graced(p) {
			return true
		}
	}
	return false
}

// graced är sant för den som porten väntar på när solen har gått ner: den
// invigde, och den som bär timglaset.
func (g *Game) graced(p *Player) bool {
	return p.initiated || p.has("timglas") != nil
}

// heartbeat är vad den invigde hör: åt vilket håll kammaren ligger.
func (g *Game) heartbeat(p *Player) string {
	if !p.initiated || p.Status != statusIn || p.Cellar > 0 {
		return ""
	}
	cx, cy := g.chamber()
	if p.X == cx && p.Y == cy {
		return g.data.Text("hjartslag_har")
	}
	idx := map[[2]int]int{{-1, -1}: 0, {1, -1}: 1, {-1, 1}: 2, {1, 1}: 3, {0, -1}: 4, {0, 1}: 5, {-1, 0}: 6, {1, 0}: 7}[[2]int{sign(cx - p.X), sign(cy - p.Y)}]
	return g.data.Text("hjartslag", "horn", g.data.Horn[idx])
}

// portGrace: när solen går ner håller porten sig öppen en runda till för
// de invigda och den som bär timglaset. De andra är förlorade. Sant om vandringen
// ska fortsätta en runda.
func (g *Game) portGrace() bool {
	if g.graceUsed {
		return false
	}
	if !g.anyGracedInside() {
		return false
	}
	g.graceUsed = true
	for _, p := range g.players {
		if p == nil || p.Status != statusIn {
			continue
		}
		if p.initiated {
			g.logTo(p, g.data.Text("porten_vantar"), "rost")
			continue
		}
		if g.useItem(p, "timglas") {
			continue // timglaset säger själv att solen står kvar för en
		}
		p.Status = statusLost
		g.chronicle("kronika_forlorad", "namn", p.Name)
		g.logTo(p, g.data.Text("forlorad_du"), "varn")
		g.logOthers(p, g.data.Text("forlorad", "namn", p.Name), "varn")
	}
	return true
}

// canWhisper är sant när vandraren kan viska nu: för att bli invigd, eller
// som invigd i kammaren för att lugna ormen.
func (g *Game) canWhisper(p *Player) bool {
	if p.bot != nil || p.Status != statusIn || p.Cellar > 0 || g.phase != phasePlay {
		return false
	}
	if p.initiated {
		return !p.calmed && g.rooms[p.Y][p.X].Kind == "kammare"
	}
	return g.level >= max(1, g.data.Avslojande.ViskaFranNiva) && p.guessLevel < g.level
}
