package main

// Krönikan: vad som hände under vandringen, berättat för alla när den är
// slut. Under vandringen skrivs händelserna ner här. Texterna står i
// data/texter.json under kronika_*.

import (
	"strconv"
	"strings"
)

// chronEvent är en händelse i krönikan, redan skriven som text.
type chronEvent struct {
	round int
	text  string
}

// chronStats räknar det som bara blir en siffra per vandrare i krönikan.
type chronStats struct {
	takes    int      // nävar ur kammaren
	takeGold int      // guld ur kammaren
	traps    int      // fällor som slog till
	kills    int      // varelser som vandraren har fällt
	searched int      // rum som vandraren har sökt igenom
	items    int      // föremål som vandraren har hittat
	hits     int      // slag som träffade en annan vandrare
	stolen   int      // guld som vandraren har tagit från andra vandrare
	draws    []string // ormkorten vandraren har dragit, i tur och ordning
}

// cardTag är ett ormkort i krönikan: id:t för bilden och namnet.
type cardTag struct {
	ID   string `json:"id"`
	Namn string `json:"namn"`
}

// cardTags gör vandrarens dragna kort till krönikans kort.
func (g *Game) cardTags(p *Player) []cardTag {
	var out []cardTag
	for _, id := range p.stats.draws {
		if c := g.data.card(id); c != nil {
			out = append(out, cardTag{ID: c.ID, Namn: c.Namn})
		}
	}
	return out
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
	p.death = g.data.Text("kronika_rad", "runda", strconv.Itoa(g.round), "text", g.chron[len(g.chron)-1].text)
}

// resetChronicle börjar en ny krönika inför en vandring.
func (g *Game) resetChronicle() {
	g.chron = nil
	g.epilog = nil
	for _, p := range g.players {
		if p != nil {
			p.stats = chronStats{}
			p.hurt("", "")
			p.journal, p.death, p.epilog = nil, "", nil
		}
	}
}

// chronLine är en rad i krönikan, med sin klass i loggen.
type chronLine struct {
	Text  string    `json:"text"`
	Class string    `json:"c,omitempty"`
	Kort  []cardTag `json:"kort,omitempty"` // ormkorten vandraren drog, på raden med siffror
}

// chronicleLines är krönikan: rubriken, händelserna i tur och ordning och
// sedan en rad för varje vandrare.
func (g *Game) chronicleLines() []chronLine {
	lines := []chronLine{{Text: g.data.Text("kronika_rubrik"), Class: "rubrik"}}
	for _, e := range g.chron {
		lines = append(lines, chronLine{Text: g.data.Text("kronika_rad", "runda", strconv.Itoa(e.round), "text", e.text)})
	}
	// De invigda som kom ut lyfts fram, bredvid den som vann på guld.
	var knew []string
	for _, p := range g.players {
		if p != nil && p.initiated && p.Status == statusOut {
			knew = append(knew, p.Name)
		}
	}
	if len(knew) > 0 {
		lines = append(lines, chronLine{Text: g.data.Text("kronika_invigda_ute", "namn", strings.Join(knew, ", ")), Class: "rost"})
	}
	for _, p := range g.players {
		if p != nil && p.Status != statusNew {
			lines = append(lines, chronLine{Text: g.statsLine(p), Class: "dim", Kort: g.cardTags(p)})
		}
	}
	return lines
}

// statsLine är vandrarens rad med siffror i krönikan.
func (g *Game) statsLine(p *Player) string {
	s := p.stats
	line := g.data.Text("kronika_vandrare", "namn", p.Name,
		"naver", strconv.Itoa(s.takes), "kammarguld", strconv.Itoa(s.takeGold),
		"fallor", strconv.Itoa(s.traps), "varelser", strconv.Itoa(s.kills),
		"sokt", strconv.Itoa(s.searched))
	// Striderna mot andra vandrare står bara med hos den som har slagits.
	if s.hits > 0 {
		line += g.data.Text("kronika_vandrare_strid",
			"traffar", strconv.Itoa(s.hits), "stulet", strconv.Itoa(s.stolen))
	}
	return line
}

// journeyRound är det som hände vandraren under en runda: rundans kort.
type journeyRound struct {
	Round int           `json:"runda"`
	Lines []summaryLine `json:"rader"`
}

// journeyView är vandrarens egen vandring, som den som har dött ser på en
// egen skärm: rundans kort runda för runda, hur hen dog och siffrorna. Den
// berättar inget om de andra utom hur många som är kvar, så att den döda
// inte kan avslöja något för de levande. När vandringen är slut får alla
// sin vandring, bredvid hela krönikan.
type journeyView struct {
	Rounds []journeyRound `json:"rundor"`
	Death  string         `json:"dod,omitempty"`
	Stats  string         `json:"siffror"`
	Cards  []cardTag      `json:"kort,omitempty"` // ormkorten vandraren har dragit
	In     int            `json:"inne"`           // vandrare kvar i tornet
	Out    int            `json:"ute"`            // vandrare som har gått ut
}

func (g *Game) journeyOf(me *Player) *journeyView {
	if me == nil {
		return nil
	}
	v := &journeyView{Rounds: me.journal, Death: me.death, Stats: g.statsLine(me), Cards: g.cardTags(me)}
	if v.Rounds == nil {
		v.Rounds = []journeyRound{}
	}
	for _, p := range g.players {
		if p == nil {
			continue
		}
		switch p.Status {
		case statusIn:
			v.In++
		case statusOut:
			v.Out++
		}
	}
	return v
}
