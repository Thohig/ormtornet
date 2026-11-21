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
// namn i tre delar, ormens namn, som delas ut till människorna i tornet,
// och viskningen ska innehålla namnet. Saknar datafilerna namn räcker
// svaren.

// usesName är sant när invigningen kräver ormens namn i den här vandringen.
func (g *Game) usesName() bool {
	return g.cfg.Reveal && len(g.serpentName) == 3
}

// pickSerpentName lottar ormens namn när en vandring börjar.
func (g *Game) pickSerpentName() {
	g.serpentName = nil
	names := g.data.Avslojande.Namn
	if !g.cfg.Reveal || len(names) == 0 {
		return
	}
	g.serpentName = names[g.rng.IntN(len(names))]
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

// namePartsFor är delarna vandraren har, med versaler, och tomma där den
// inte vet. Ingenting om vandraren inte har fått någon del.
func (g *Game) namePartsFor(p *Player) []string {
	if !g.usesName() || p.partCount() == 0 {
		return nil
	}
	out := make([]string, 3)
	for i, k := range p.parts {
		if k {
			out[i] = strings.ToUpper(g.serpentName[i])
		}
	}
	return out
}

// dealNameParts ger människorna i tornet delar av namnet, så att var och en
// har lika många delar som nivån minus ett (högst tre). Varje del som delas
// ut är den som färst i gruppen har, så att gruppen tillsammans vet så
// mycket som möjligt och måste prata för att få ihop namnet.
func (g *Game) dealNameParts() {
	if !g.usesName() {
		return
	}
	want := min(g.level-1, 3)
	var group []*Player
	for _, p := range g.players {
		if p != nil && p.bot == nil && p.Status == statusIn {
			group = append(group, p)
		}
	}
	g.rng.Shuffle(len(group), func(i, j int) { group[i], group[j] = group[j], group[i] })
	for range 3 {
		for _, p := range group {
			if p.partCount() >= want {
				continue
			}
			best, bestN := -1, 0
			for i := range 3 {
				if p.parts[i] {
					continue
				}
				n := 0
				for _, o := range group {
					if o.parts[i] {
						n++
					}
				}
				if best < 0 || n < bestN {
					best, bestN = i, n
				}
			}
			p.parts[best] = true
			g.logTo(p, g.namePartText(best), "info")
		}
	}
}

// namePartText är texten när vandraren får del i av namnet.
func (g *Game) namePartText(i int) string {
	a := &g.data.Avslojande
	list := a.Delar[strconv.Itoa(g.level)]
	if len(list) == 0 {
		list = a.Delar["4"]
	}
	ord := ""
	if i < len(a.Ordning) {
		ord = a.Ordning[i]
	}
	return strings.NewReplacer("{ordning}", ord, "{del}", strings.ToUpper(g.serpentName[i])).Replace(pick(list))
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

// Whisper tar emot en viskning från en vandrare.
func (g *Game) Whisper(c *Client, text string) {
	p := c.player
	if p == nil || p.bot != nil || g.phase != phasePlay || p.Status != statusIn || !g.cfg.Reveal {
		return
	}
	if p.initiated {
		// I kammaren kan den invigde lugna ormen, en gång per vandring.
		switch {
		case g.rooms[p.Y][p.X].Kind != "kammare":
			g.logTo(p, g.data.Text("viskning_redan"), "dim")
		case p.calmed:
			g.logTo(p, g.data.Text("invigd_lugnat"), "dim")
		default:
			p.calmed = true
			g.calmSerpent()
			g.logTo(p, g.data.Text("invigd_lugnar"), "rost")
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
		g.logTo(p, g.data.Text("viskning_tidigt"), "dim")
		return
	}
	if p.guessLevel >= g.level {
		g.logTo(p, g.data.Text("viskning_vanta"), "dim")
		return
	}
	p.guessLevel = g.level
	p.whispered = text
	if !g.whisperRight(text) {
		g.logTo(p, g.data.Text("viskning_fel"), "dim")
		return
	}
	p.initiated = true
	g.soundTo(p, "rost")
	g.logTo(p, g.data.Text("invigd"), "rost")
	g.logOthers(p, g.data.Text("invigd_andra"), "dim")
	g.chronicle("kronika_invigd", "namn", p.Name, "runda", strconv.Itoa(g.round))
	g.checkMercy()
}

// takeNameParts: den som dödar en vandrare hittar de delar av namnet som
// offret hade och dödaren saknade.
func (g *Game) takeNameParts(killer, victim *Player) {
	if !g.usesName() || killer.bot != nil {
		return
	}
	var got []string
	for i := range 3 {
		if victim.parts[i] && !killer.parts[i] {
			killer.parts[i] = true
			got = append(got, strings.ToUpper(g.serpentName[i]))
		}
	}
	if len(got) == 0 {
		return
	}
	g.logTo(killer, g.data.Text("namndel_byte", "offer", victim.Name, "delar", strings.Join(got, ", ")), "rost")
	g.chronicle("kronika_namndel", "namn", killer.Name, "offer", victim.Name)
}

// checkMercy ger tornets nåd när alla människor som är kvar inne är
// invigda: en andel av vakna-korten tas bort, en gång per vandring.
func (g *Game) checkMercy() {
	if !g.usesName() || g.mercyGiven || g.phase != phasePlay {
		return
	}
	var inside []*Player
	for _, p := range g.players {
		if p != nil && p.bot == nil && p.Status == statusIn {
			if !p.initiated {
				return
			}
			inside = append(inside, p)
		}
	}
	if len(inside) == 0 {
		return
	}
	wake := 0
	for _, id := range g.deck {
		if c := g.data.card(id); c != nil && c.Vaknar {
			wake++
		}
	}
	n := max(1, int(float64(wake)*g.data.Avslojande.NadAndel))
	removed := 0
	for i := 0; i < len(g.deck) && removed < n; {
		if c := g.data.card(g.deck[i]); c != nil && c.Vaknar {
			g.deck = append(g.deck[:i], g.deck[i+1:]...)
			removed++
			continue
		}
		i++
	}
	g.mercyGiven = true
	g.mercy = n
	for _, p := range inside {
		g.logTo(p, g.data.Text("nad"), "rost")
		g.soundTo(p, "rost")
	}
	g.chronicle("kronika_nad")
}

func (g *Game) anyInitiatedInside() bool {
	for _, p := range g.players {
		if p != nil && p.Status == statusIn && p.initiated {
			return true
		}
	}
	return false
}

// heartbeat är vad den invigde hör: åt vilket håll kammaren ligger.
func (g *Game) heartbeat(p *Player) string {
	if !p.initiated || p.Status != statusIn || p.Cellar > 0 {
		return ""
	}
	c := g.center()
	if p.X == c && p.Y == c {
		return g.data.Text("hjartslag_har")
	}
	idx := map[[2]int]int{{-1, -1}: 0, {1, -1}: 1, {-1, 1}: 2, {1, 1}: 3, {0, -1}: 4, {0, 1}: 5, {-1, 0}: 6, {1, 0}: 7}[[2]int{sign(c - p.X), sign(c - p.Y)}]
	return g.data.Text("hjartslag", "horn", g.data.Horn[idx])
}

// portGrace: när solen går ner håller porten sig öppen en runda till för
// de invigda som är kvar inne. De andra är förlorade. Sant om vandringen
// ska fortsätta en runda.
func (g *Game) portGrace() bool {
	if g.graceUsed {
		return false
	}
	if !g.anyInitiatedInside() {
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
	if !g.cfg.Reveal || p.bot != nil || p.Status != statusIn || p.Cellar > 0 || g.phase != phasePlay {
		return false
	}
	if p.initiated {
		return !p.calmed && g.rooms[p.Y][p.X].Kind == "kammare"
	}
	return g.level >= max(1, g.data.Avslojande.ViskaFranNiva) && p.guessLevel < g.level
}

// calmSerpent tar ett vakna-kort ur leken, eller lägger till ett sovkort om
// det inte finns något vakna-kort kvar.
func (g *Game) calmSerpent() string {
	for i, id := range g.deck {
		if c := g.data.card(id); c != nil && c.Vaknar {
			g.deck = append(g.deck[:i], g.deck[i+1:]...)
			return "vakna"
		}
	}
	for _, c := range g.data.Ormen.Kort {
		if !c.Vaknar {
			g.deck = append(g.deck, c.ID)
			return "sov"
		}
	}
	return ""
}
