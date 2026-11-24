package main

// Berättelsen: när vandringen är slut berättar tornet vad som hände, som
// en gemensam historia i tre delar i stället för drag för drag. Inträdet
// säger vilka som gick in, höjdpunkterna är krönikans viktigaste händelser
// med var sitt spelkort, och slutet säger vem som kom ut. Epilogen står
// under, som förut. Vilka händelser som räknas, hur mycket de väger och
// hur många som tas med står under spelkort.berattelse i data/texter.json.
// Varje händelse berättas med en egen text, ber_ och krönikans nyckel utan
// kronika_, med {nar} först (docs/berattelsen.md).

import (
	"math/rand/v2"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// StoryData styr berättelsens höjdpunkter.
type StoryData struct {
	Antal   int          `json:"antal"`    // så många höjdpunkter som mest
	PerSlag int          `json:"per_slag"` // så många av samma slag
	Vinnare int          `json:"vinnare"`  // det här mycket mer väger en händelse om den som vann
	Slag    []StoryKind  `json:"slag"`
	Titlar  []StoryTitle `json:"titlar"` // titlarna på raden om var och en, den viktigaste först
}

// StoryTitle är en titel som den får som har mest av något (Matt), till
// exempel guld ur kammaren. Text är titeln, en nyckel bland texterna.
type StoryTitle struct {
	Matt string `json:"matt"`
	Text string `json:"text"`
}

// titleMeasures är det titlarna kan mätas i, från vandrarens siffror.
var titleMeasures = map[string]func(p *Player) int{
	"kammaren": func(p *Player) int { return p.stats.takeGold },
	"stulet":   func(p *Player) int { return p.stats.stolen },
	"traffar":  func(p *Player) int { return p.stats.hits },
	"varelser": func(p *Player) int { return p.stats.kills },
	"fallor":   func(p *Player) int { return p.stats.traps },
	"sokt":     func(p *Player) int { return p.stats.searched },
	"foremal":  func(p *Player) int { return p.stats.items },
	"ormkort":  func(p *Player) int { return len(p.stats.draws) },
	"rum": func(p *Player) int {
		n := 0
		for _, row := range p.known {
			for _, k := range row {
				if k {
					n++
				}
			}
		}
		return n
	},
}

// StoryKind är ett slags höjdpunkt: krönikans nycklar som hör dit, hur
// mycket den väger, motivet och rubriken på kortet, och en text som säger
// hur många gånger det hände, när det hände mer än en gång.
type StoryKind struct {
	ID      string   `json:"id"`
	Nycklar []string `json:"nycklar"`
	Vikt    int      `json:"vikt"`
	Forst   int      `json:"forst,omitempty"` // det här mycket mer väger den första av sitt slag
	Bild    string   `json:"bild"`
	Rubrik  string   `json:"rubrik"`
	Flera   string   `json:"flera,omitempty"` // text med {antal}, när slaget hände flera gånger
}

// storyPerson är raden om en vandrare under berättelsen: titeln, om hen
// fick någon, hens största stund och hur det slutade.
type storyPerson struct {
	Name  string `json:"namn"`
	Slot  int    `json:"slot"`
	Title string `json:"titel,omitempty"`
	Text  string `json:"text"`
	Stund string `json:"stund,omitempty"` // bara den största stunden, för slutskärmens lista
}

// storyPart är ett stycke i berättelsen: inträdet (in), en höjdpunkt
// (mitt) eller slutet (slut), med ett spelkort bredvid när det finns.
type storyPart struct {
	Del   string    `json:"del"`
	Text  string    `json:"text"`
	Round int       `json:"runda,omitempty"`
	Kort  *gameCard `json:"kort,omitempty"`
}

// storyKey är berättelsens text för krönikans nyckel key.
func storyKey(key string) string {
	return "ber_" + strings.TrimPrefix(key, "kronika_")
}

// tal skriver ett litet tal med bokstäver (tal_2 till tal_10), annars med
// siffror.
func (g *Game) tal(n int) string {
	if k := "tal_" + strconv.Itoa(n); len(g.data.Texter[k]) > 0 {
		return g.data.Text(k)
	}
	return strconv.Itoa(n)
}

// andList sätter ihop namnen: "A", "A och B", "A, B och C".
func (g *Game) andList(names []string) string {
	switch len(names) {
	case 0:
		return ""
	case 1:
		return names[0]
	}
	return g.data.Text("ber_lista", "namn", strings.Join(names[:len(names)-1], ", "), "sista", names[len(names)-1])
}

// buildStory skriver berättelsen när vandringen är slut. Den skrivs en
// gång, så att texterna inte byter variant varje gång tillståndet skickas.
func (g *Game) buildStory() {
	g.story = nil
	if in := g.storyIntro(); in != "" {
		g.story = append(g.story, storyPart{Del: "in", Text: in})
	}
	high := g.storyHighlights()
	if len(high) == 0 {
		g.story = append(g.story, storyPart{Del: "mitt", Text: g.data.Text("ber_tyst")})
	}
	g.story = append(g.story, high...)
	g.story = append(g.story, g.storyEnd())
	g.storyPeople = g.storyPersons()
	g.books = g.goldBooks()
}

// storyKind är händelsens slag i berättelsen, -1 när den inte räknas.
func (g *Game) storyKind(e chronEvent) int {
	return slices.IndexFunc(g.data.Spelkort.Berattelse.Slag, func(k StoryKind) bool { return slices.Contains(k.Nycklar, e.key) })
}

// storyTitles delar ut titlarna: var och en går till den som ensam har
// mest av sitt mått, om den inte redan har en titel. Ingen får två, och
// lika mycket ger ingen titel.
func (g *Game) storyTitles(people []*Player) map[*Player]string {
	out := map[*Player]string{}
	for _, t := range g.data.Spelkort.Berattelse.Titlar {
		measure := titleMeasures[t.Matt]
		if measure == nil {
			continue
		}
		var best *Player
		most, tie := 0, false
		for _, p := range people {
			if out[p] != "" {
				continue
			}
			switch n := measure(p); {
			case n > most:
				best, most, tie = p, n, false
			case n == most && n > 0:
				tie = true
			}
		}
		if best != nil && !tie {
			out[best] = g.data.Text(t.Text)
		}
	}
	return out
}

// storyPersons skriver raden om var och en, i placeringens ordning: titeln,
// den tyngsta händelsen som handlar om hen och hur det slutade. Var hens
// största stund att dö säger den redan hur det slutade.
func (g *Game) storyPersons() []storyPerson {
	sd := g.data.Spelkort.Berattelse
	var people []*Player
	for _, r := range g.results() {
		if p := g.playerByName(r.Name); p != nil {
			people = append(people, p)
		}
	}
	titles := g.storyTitles(people)
	var out []storyPerson
	for _, p := range people {
		best, weight := -1, 0
		for i, e := range g.chron {
			k := g.storyKind(e)
			if k < 0 || e.who != p.Name || sd.Slag[k].Vikt <= weight {
				continue
			}
			best, weight = i, sd.Slag[k].Vikt
		}
		var parts []string
		var moment string
		ended := false
		if best >= 0 {
			// Raden börjar redan med namnet: "Sigge. Fällde en grönaktig
			// skepnad." Krönikans texter börjar med namnet och verbet.
			e := g.chron[best]
			text := e.text
			if rest, ok := strings.CutPrefix(text, p.Name+" "); ok {
				text = upperFirst(rest)
			}
			parts = append(parts, text)
			moment = text
			ended = sd.Slag[g.storyKind(e)].ID == "dod"
		}
		if !ended {
			switch {
			case p.Status == statusOut && p.Treasure > 0:
				parts = append(parts, g.data.Text("ber_vem_ute", "guld", strconv.Itoa(p.Treasure)))
			case p.Status == statusOut:
				parts = append(parts, g.data.Text("ber_vem_tom"))
			case p.Status == statusDead:
				parts = append(parts, g.data.Text("ber_vem_dod"))
			default:
				parts = append(parts, g.data.Text("ber_vem_kvar"))
			}
		}
		out = append(out, storyPerson{Name: p.Name, Slot: p.Slot, Title: titles[p], Text: strings.Join(parts, " "), Stund: moment})
	}
	return out
}

// storyIntro: vilka som gick in, med sina klasser, och hur länge
// vandringen varade.
func (g *Game) storyIntro() string {
	var who []string
	for _, p := range g.players {
		if p == nil || p.Status == statusNew {
			continue
		}
		if cl := g.classOf(p); cl != nil {
			who = append(who, g.data.Text("ber_in_vem", "klass", cl.Kort, "namn", p.Name))
		} else {
			who = append(who, p.Name)
		}
	}
	if len(who) == 0 {
		return ""
	}
	return upperFirst(g.data.Text("ber_in", "antal", g.tal(len(who)), "vilka", g.andList(who), "rundor", strconv.Itoa(g.round)))
}

// storyHighlights väljer krönikans viktigaste händelser: varje händelse
// väger sitt slags vikt, mer om den är den första av sitt slag eller
// handlar om den som vann. De tyngsta tas med, högst PerSlag av varje och
// Antal i allt, och berättas i den ordning de hände, med när de hände
// räknat från den förra.
func (g *Game) storyHighlights() []storyPart {
	sd := g.data.Spelkort.Berattelse
	winners := map[string]bool{}
	for _, r := range g.results() {
		if r.Winner {
			winners[r.Name] = true
		}
	}
	type pick struct {
		kind, order, weight int
	}
	var all []pick
	count := map[int]int{}
	for i, e := range g.chron {
		k := g.storyKind(e)
		if k < 0 {
			continue
		}
		w := sd.Slag[k].Vikt
		if count[k] == 0 {
			w += sd.Slag[k].Forst
		}
		if e.who != "" && winners[e.who] {
			w += sd.Vinnare
		}
		count[k]++
		all = append(all, pick{k, i, w})
	}
	sort.SliceStable(all, func(a, b int) bool { return all[a].weight > all[b].weight })
	var chosen []pick
	per := map[int]int{}
	for _, p := range all {
		if len(chosen) >= sd.Antal {
			break
		}
		if sd.PerSlag > 0 && per[p.kind] >= sd.PerSlag {
			continue
		}
		per[p.kind]++
		chosen = append(chosen, p)
	}
	sort.SliceStable(chosen, func(a, b int) bool { return chosen[a].order < chosen[b].order })

	var out []storyPart
	told := map[int]bool{} // slagen vars antal redan är berättat
	said := map[string]bool{}
	last := 0
	for _, p := range chosen {
		e := g.chron[p.order]
		k := sd.Slag[p.kind]
		// Samma slags händelse berättas inte med samma ord två gånger: av
		// textens varianter tas, i slumpad ordning, den första som inte
		// redan är sagd. Det som jämförs är texten utan namnen och talen.
		nar := g.when(e.round, last)
		vars := append([]string{"nar", nar}, e.vars...)
		var text, shape string
		list := g.data.Texter[storyKey(e.key)]
		for _, i := range rand.Perm(len(list)) {
			text = upperFirst(fill(list[i], vars...))
			shape = strings.Replace(text, upperFirst(nar), "", 1)
			for j := 1; j < len(e.vars); j += 2 {
				if e.vars[j] != "" {
					shape = strings.ReplaceAll(shape, e.vars[j], "")
				}
			}
			if !said[shape] {
				break
			}
		}
		said[shape] = true
		if burit := varOf(e.vars, "burit"); burit != "" {
			text += " " + g.data.Text("ber_dod_guld", "guld", burit)
		}
		if k.Flera != "" && count[p.kind] > 1 && !told[p.kind] {
			told[p.kind] = true
			text += " " + g.data.Text(k.Flera, "antal", g.tal(count[p.kind]))
		}
		bild := k.Bild
		if e.bild != "" {
			bild = e.bild
		}
		c := gameCard{Slag: k.ID, Bild: bild, Rubrik: upperFirst(g.data.Text(k.Rubrik)), Fallen: k.ID == "varelse"}
		out = append(out, storyPart{Del: "mitt", Text: text, Round: e.round, Kort: &c})
		last = e.round
	}
	return out
}

// when säger när något hände, räknat från rundan last då det förra
// hände (0 för det första).
func (g *Game) when(round, last int) string {
	switch {
	case round == g.round && last > 0 && round-last > 1:
		return g.data.Text("ber_nar_sist")
	case last == 0:
		return g.data.Text("ber_nar_runda", "runda", strconv.Itoa(round))
	case round == last:
		return g.data.Text("ber_nar_samma")
	case round == last+1:
		return g.data.Text("ber_nar_efter")
	case round-last > 10:
		return g.data.Text("ber_nar_langt")
	}
	return g.data.Text("ber_nar_senare", "antal", g.tal(round-last))
}

// varOf ger värdet för nyckeln key bland textens variabler.
func varOf(vars []string, key string) string {
	for i := 0; i+1 < len(vars); i += 2 {
		if vars[i] == key {
			return vars[i+1]
		}
	}
	return ""
}

// storyEnd: vem som bar ut mest, vilka andra som kom ut, vilka som aldrig
// kom ut och vilka som var kvar när solen gick ner. Vinnaren, eller tornet
// om ingen vann, får kortet.
func (g *Game) storyEnd() storyPart {
	res := g.results()
	var win, out, empty, dead, lost []string
	gold := 0
	for _, r := range res {
		switch {
		case r.Winner:
			win = append(win, r.Name)
			gold = r.Treasure
		case r.Status == statusOut && r.Treasure > 0:
			out = append(out, r.Name)
		case r.Status == statusOut:
			empty = append(empty, r.Name)
		case r.Status == statusDead:
			dead = append(dead, r.Name)
		default:
			lost = append(lost, r.Name)
		}
	}
	var parts []string
	var c gameCard
	if len(win) > 0 {
		parts = append(parts, g.data.Text("ber_slut_vinnare", "namn", g.andList(win), "guld", strconv.Itoa(gold)))
		c = gameCard{Slag: "vinnare", Bild: "ut", Rubrik: upperFirst(g.data.Text("stund_vinnare"))}
	} else {
		parts = append(parts, g.data.Text("ber_slut_ingen"))
		c = gameCard{Slag: "tornet", Bild: "sol", Rubrik: upperFirst(g.data.Text("ber_tornet_vann"))}
	}
	add := func(names []string, key string) {
		if len(names) > 0 {
			parts = append(parts, g.data.Text(key, "namn", g.andList(names)))
		}
	}
	add(out, "ber_slut_ute")
	add(empty, "ber_slut_tom")
	add(lost, "ber_slut_kvar")
	add(dead, "ber_slut_doda")
	return storyPart{Del: "slut", Text: strings.Join(parts, " "), Kort: &c}
}
