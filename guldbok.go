package main

// Guldboken: hur var och en fick sitt guld. Varje gång en vandrares guld
// ändras skrivs det ner (goldNote), med rundan, hur mycket och vad som
// hände. När vandringen är slut blir det en berättelse för varje vandrare,
// som går att läsa genom att välja vandraren på slutskärmen, och raderna
// runda för runda bakom en knapp (docs/guldboken.md). Texterna heter gb_
// i data/texter.json.

import (
	"strconv"
	"strings"
)

// goldEntry är en ändring av en vandrares guld: amount är plus eller
// minus, total vad vandraren bar efteråt.
type goldEntry struct {
	round  int
	kind   string
	amount int
	total  int
	vars   []string
}

// goldNote skriver ner att p:s guld ändrades med amount. Den anropas efter
// ändringen, så att p.Treasure är summan efteråt.
func (g *Game) goldNote(p *Player, kind string, amount int, vars ...string) {
	if p == nil || amount == 0 {
		return
	}
	p.goldLog = append(p.goldLog, goldEntry{round: g.round, kind: kind, amount: amount, total: p.Treasure, vars: vars})
}

// goldRow är en rad i guldboken runda för runda.
type goldRow struct {
	Round  int    `json:"runda"`
	Text   string `json:"text"`
	Amount int    `json:"guld"`
	Total  int    `json:"summa"`
}

// goldBook är en vandrares guldbok på slutskärmen: berättelsen i stycken
// och raderna under.
type goldBook struct {
	Name  string    `json:"namn"`
	Story []string  `json:"berattelse"`
	Rows  []goldRow `json:"rader"`
}

// goldVars är radens variabler med guldet som ett positivt tal.
func goldVars(e goldEntry, extra ...string) []string {
	n := e.amount
	if n < 0 {
		n = -n
	}
	return append(append([]string{"guld", strconv.Itoa(n), "summa", strconv.Itoa(e.total)}, e.vars...), extra...)
}

// goldBooks skriver allas guldbok, i placeringens ordning.
func (g *Game) goldBooks() []goldBook {
	var out []goldBook
	for _, r := range g.results() {
		p := g.playerByName(r.Name)
		if p == nil {
			continue
		}
		book := goldBook{Name: p.Name, Story: g.goldStory(p, r), Rows: []goldRow{}}
		for _, e := range p.goldLog {
			book.Rows = append(book.Rows, goldRow{Round: e.round, Text: g.data.Text("gb_rad_"+e.kind, goldVars(e)...), Amount: e.amount, Total: e.total})
		}
		out = append(out, book)
	}
	return out
}

// goldRun är hur många rundor som högst får gå mellan två ändringar av
// samma slag för att de ska berättas som en.
const goldRun = 3

// goldAtGate är det som läggs till vid porten; det berättas i slutet.
func goldAtGate(kind string) bool {
	return kind == "foremal" || kind == "invigd" || kind == "invigd_forst" || kind == "ute_igen"
}

// goldStory är berättelsen om hur p fick sitt guld: vem hen var och hur
// det gick, vägen dit i den ordning det hände, och slutet.
func (g *Game) goldStory(p *Player, r resultView) []string {
	t := g.data.Text
	name := p.Name
	if cl := g.classOf(p); cl != nil {
		name = g.data.Text("ber_in_vem", "klass", cl.Kort, "namn", p.Name)
	}
	gold := strconv.Itoa(p.Treasure)
	var open string
	switch {
	case p.Status == statusOut && r.Winner:
		open = t("gb_in_vinnare", "namn", name, "guld", gold)
	case p.Status == statusOut && p.Treasure > 0:
		open = t("gb_in_ute", "namn", name, "guld", gold)
	case p.Status == statusOut:
		open = t("gb_in_tom", "namn", name)
	case p.Status == statusDead:
		open = t("gb_in_dod", "namn", name)
	default:
		open = t("gb_in_kvar", "namn", name)
	}
	story := []string{upperFirst(open)}

	// Vägen: ändringarna i tur och ordning, med när de hände räknat från
	// den förra. Flera av samma slag i rad blir en mening.
	var road []goldEntry
	for _, e := range p.goldLog {
		if !goldAtGate(e.kind) {
			road = append(road, e)
		}
	}
	if len(road) == 0 {
		story = append(story, t("gb_aldrig"))
	} else {
		biggest := -1
		gains := 0
		for i, e := range road {
			if e.amount > 0 {
				gains++
				if biggest < 0 || e.amount > road[biggest].amount {
					biggest = i
				}
			}
		}
		var sentences []string
		last := 0
		for i := 0; i < len(road); {
			e := road[i]
			j := i + 1
			sum := e.amount
			for j < len(road) && road[j].kind == e.kind && road[j].round-road[j-1].round <= goldRun && e.kind != "slag" && e.kind != "slagen" {
				sum += road[j].amount
				j++
			}
			nar := g.when(e.round, last)
			var s string
			if j-i > 1 {
				end := road[j-1]
				n := sum
				if n < 0 {
					n = -n
				}
				s = t("gb_flera_"+e.kind, "nar", nar, "antal", g.tal(j-i), "guld", strconv.Itoa(n), "summa", strconv.Itoa(end.total))
				last = end.round
				if biggest >= i && biggest < j && gains > j-i {
					s += " " + t("gb_storst_av", "guld", strconv.Itoa(road[biggest].amount))
				}
			} else {
				s = t("gb_"+e.kind, append([]string{"nar", nar}, goldVars(e)...)...)
				last = e.round
				if i == biggest && gains > 1 {
					s += " " + t("gb_storst")
				}
			}
			s = upperFirst(s)
			if road[j-1].amount < 0 && road[j-1].total > 0 {
				s += " " + t("gb_kvar", "summa", strconv.Itoa(road[j-1].total))
			}
			sentences = append(sentences, s)
			i = j
		}
		story = append(story, strings.Join(sentences, " "))
	}

	// Slutet: det porten lade till, vart guldet tog vägen för den som dog
	// och placeringen, för den som kom ut med guld.
	var end []string
	for _, e := range p.goldLog {
		if goldAtGate(e.kind) {
			end = append(end, t("gb_"+e.kind, goldVars(e)...))
		}
	}
	for _, id := range p.goldTakenBy {
		if o := g.playerByID(id); o != nil {
			key := "gb_guld_ute"
			if o.Status != statusOut {
				key = "gb_guld_inte_ute"
			}
			end = append(end, t(key, "namn", o.Name))
		}
	}
	// Den som var kvar när solen gick ner står som förlorad i slutet.
	if (p.Status == statusIn || p.Status == statusLost) && p.Treasure > 0 {
		end = append(end, t("gb_kvar_slut", "guld", gold))
	}
	if p.Status == statusOut && p.Treasure > 0 && len(end) > 0 {
		end = append(end, t("gb_summa", "guld", gold))
	}
	if r.Place > 0 && p.Status == statusOut && p.Treasure > 0 {
		end = append(end, t("gb_plats", "plats", strconv.Itoa(r.Place), "antal", strconv.Itoa(len(g.results()))))
	}
	if len(end) > 0 {
		story = append(story, strings.Join(end, " "))
	}
	return story
}
