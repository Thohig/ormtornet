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

// Whisper tar emot en viskning från en vandrare.
func (g *Game) Whisper(c *Client, text string) {
	p := c.player
	if p == nil || p.bot != nil || g.phase != phasePlay || p.Status != statusIn || !g.cfg.Reveal {
		return
	}
	if p.initiated {
		g.logTo(p, g.data.Text("viskning_redan"), "dim")
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
	if !whisperMatches(text, g.data.Avslojande.Svar) {
		g.logTo(p, g.data.Text("viskning_fel"), "dim")
		return
	}
	p.initiated = true
	g.soundTo(p, "rost")
	g.logTo(p, g.data.Text("invigd"), "rost")
	g.logOthers(p, g.data.Text("invigd_andra"), "dim")
	g.chronicle("kronika_invigd", "namn", p.Name, "runda", strconv.Itoa(g.round))
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
