package main

// Krönikan: vad som hände under vandringen, berättat för alla när den är
// slut. Under vandringen skrivs händelserna ner här. Texterna står i
// data/texter.json under kronika_*.

import "strconv"

// chronEvent är en händelse i krönikan, redan skriven som text.
type chronEvent struct {
	round int
	text  string
}

// chronStats räknar det som bara blir en siffra per vandrare i krönikan.
type chronStats struct {
	takes    int // nävar ur kammaren
	takeGold int // guld ur kammaren
	traps    int // fällor som slog till
	kills    int // varelser som vandraren har fällt
	searched int // rum som vandraren har sökt igenom
}

// chronicle skriver ner en händelse från den här rundan.
func (g *Game) chronicle(key string, vars ...string) {
	g.chron = append(g.chron, chronEvent{round: g.round, text: g.data.Text(key, vars...)})
}

// hurt säger vad som skadar vandraren, så att krönikan vet vad som dödade
// den om skadan blir den sista. kind är falla, varelse, vandrare, ormen,
// gift eller tornet.
func (p *Player) hurt(kind, what string) {
	p.hurtKind, p.hurtWhat = kind, what
}

// chronicleDeath skriver ner hur vandraren dog.
func (g *Game) chronicleDeath(p *Player) {
	key := "kronika_dod"
	if p.hurtKind != "" {
		key += "_" + p.hurtKind
	}
	if len(g.data.Texter[key]) == 0 {
		key = "kronika_dod"
	}
	g.chronicle(key, "namn", p.Name, "vad", p.hurtWhat)
}

// resetChronicle börjar en ny krönika inför en vandring.
func (g *Game) resetChronicle() {
	g.chron = nil
	for _, p := range g.players {
		if p != nil {
			p.stats = chronStats{}
			p.hurt("", "")
		}
	}
}

// tellChronicle berättar krönikan för alla: händelserna i tur och ordning,
// och sedan en rad för varje vandrare.
func (g *Game) tellChronicle() {
	g.logAll(g.data.Text("kronika_rubrik"), "rubrik")
	for _, e := range g.chron {
		g.logAll(g.data.Text("kronika_rad", "runda", strconv.Itoa(e.round), "text", e.text), "")
	}
	for _, p := range g.players {
		if p == nil {
			continue
		}
		s := p.stats
		g.logAll(g.data.Text("kronika_vandrare", "namn", p.Name,
			"naver", strconv.Itoa(s.takes), "kammarguld", strconv.Itoa(s.takeGold),
			"fallor", strconv.Itoa(s.traps), "varelser", strconv.Itoa(s.kills),
			"sokt", strconv.Itoa(s.searched)), "dim")
	}
}
