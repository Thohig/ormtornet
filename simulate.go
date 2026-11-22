package main

// Simulering: många vandringar med bara AI-spelare, utan webbläsare och
// utan väntetider. Den visar hur det brukar gå, så att siffrorna i data/
// kan justeras: hur många som kommer ut med skatt, vad de dör av, hur ofta
// ormen vaknar och hur det går för varje klass.
//
//	./ormtornet -simulera 500 -ai 4 -svarighet svar

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"sort"
	"strings"
	"time"
)

// simWalker är hur det gick för en AI-spelare under en vandring.
type simWalker struct {
	class    string
	pers     string // AI-spelarens personlighet, till exempel "djärv"
	status   string
	treasure int
	death    string // vad som dödade vandraren, till exempel "varelse: blekkryparen"
	kills    int
	traps    int
	takes    int
	won      bool // hade mest skatt av dem som kom ut
}

// simResult är en hel vandring.
type simResult struct {
	rounds  int
	wakes   int
	size    int
	walkers []simWalker
}

// simulateOne spelar en vandring med n AI-spelare och fröet seed.
func simulateOne(d *Data, cfg Config, n int, seed uint64) (simResult, error) {
	cfg.Seed = seed
	cfg.AI = n
	cfg.MinPlayers = 1
	cfg.TurnTime = time.Hour
	g := NewGame(d, cfg)
	// Tornets tidtagare väntar på låset och hinner aldrig gå medan vandringen
	// spelas här.
	g.mu.Lock()
	defer g.mu.Unlock()
	if added := g.AddBots(n); added != n {
		return simResult{}, fmt.Errorf("fick bara plats för %d AI-spelare", added)
	}
	simDealClasses(g, seed)
	g.begin()
	if g.phase == phaseChoose {
		g.startGame() // AI-spelarna har valt; vänta inte på tiden
	}
	if g.phase != phasePlay {
		return simResult{}, fmt.Errorf("vandringen började inte (fasen är %q)", g.phase)
	}
	for guard := 0; g.phase == phasePlay; guard++ {
		if guard > 20*g.sunMax {
			return simResult{}, fmt.Errorf("frö %d: vandringen tog aldrig slut", seed)
		}
		round := g.round
		for _, p := range g.players {
			if p != nil && p.bot != nil {
				g.botAct(p, round)
			}
		}
		if g.phase == phasePlay && g.round == round {
			g.resolve() // som när turtiden går ut
		}
	}
	res := simResult{rounds: g.round, wakes: g.wakes, size: g.size}
	winners := map[string]bool{}
	for _, r := range g.results() {
		if r.Winner {
			winners[r.Name] = true
		}
	}
	for _, p := range g.players {
		if p == nil {
			continue
		}
		w := simWalker{won: winners[p.Name], status: p.Status, treasure: p.Treasure,
			kills: p.stats.kills, traps: p.stats.traps, takes: p.stats.takes}
		if c := g.classOf(p); c != nil {
			w.class = c.Kort
		}
		if p.bot != nil && p.bot.pers != nil {
			w.pers = p.bot.pers.Namn
		}
		if p.Status == statusDead {
			w.death = simCause(p.hurtKind)
			if p.hurtKind != "vandrare" && p.hurtWhat != "" {
				w.death += ": " + p.hurtWhat
			}
		}
		res.walkers = append(res.walkers, w)
	}
	return res, nil
}

// simDealClasses ger AI-spelarna klasserna i tur och ordning efter fröet, i
// stället för efter personligheten. Så spelas varje klass lika ofta, och
// lika ofta av varje personlighet, och klasserna kan jämföras med varandra.
// Varje AI-spelare får en kopia av sin personlighet som bara vill ha den
// klassen; allt annat i personligheten är som förut. Med samma frö blir
// det samma klasser, så att -simulera 1 -fro fungerar.
func simDealClasses(g *Game, seed uint64) {
	k := len(g.data.Klasser)
	if !g.cfg.Classes || k == 0 {
		return
	}
	var bots []*Player
	for _, p := range g.players {
		if p != nil && p.bot != nil && p.bot.pers != nil {
			bots = append(bots, p)
		}
	}
	// I den ordning de kom in, som är den ordning personligheterna turas om.
	sort.Slice(bots, func(i, j int) bool { return bots[i].ID < bots[j].ID })
	for i, p := range bots {
		pers := *p.bot.pers
		pers.Klasser = []string{g.data.Klasser[(int(seed%uint64(k))+i)%k].ID}
		p.bot.pers = &pers
	}
}

// simCause gör skadans slag läsbart i utskriften, till exempel "fälla"
// i stället för nyckeln "falla".
func simCause(kind string) string {
	switch kind {
	case "":
		return "okänt"
	case "falla":
		return "fälla"
	case "vandrare":
		return "en annan vandrare"
	}
	return kind
}

// runSimulation spelar count vandringar och skriver en sammanfattning till
// out. Vandring i får fröet seed+i, så att en konstig vandring kan spelas
// upp igen med -fro.
func runSimulation(d *Data, cfg Config, count, bots int, seed uint64, out io.Writer) error {
	if bots <= 0 {
		bots = 4
	}
	bots = min(bots, maxPlayers)
	if seed == 0 {
		seed = 1 + rand.Uint64()%1000000
	}
	// AI-spelarna ska tänka direkt, och fröet för varje vandring ska inte
	// skrivas ut.
	d.AI.BetanketidMs = []int{3600000, 3600000}
	logOut := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(logOut)

	started := time.Now()
	var results []simResult
	for i := 0; i < count; i++ {
		r, err := simulateOne(d, cfg, bots, seed+uint64(i))
		if err != nil {
			return err
		}
		results = append(results, r)
	}
	name := "Normal"
	if s := d.difficulty(cfg.Difficulty); s != nil {
		name = s.Namn
	}
	fmt.Fprintf(out, "%d vandringar med %d AI-spelare, svårighet %s, torn %d × %d.\n",
		count, bots, name, results[0].size, results[0].size)
	fmt.Fprintf(out, "Frön %d till %d (kör om en med -simulera 1 -fro och fröet). Tog %s.\n\n",
		seed, seed+uint64(count-1), time.Since(started).Round(time.Millisecond))
	writeSimSummary(out, results)
	return nil
}

// simTally räknar ihop vandrare.
type simTally struct {
	walkers, rich, poor, dead, lost int
	gold                            int // skatt hos dem som kom ut
	kills, traps, wins              int
}

func (t *simTally) add(w simWalker) {
	t.walkers++
	if w.won {
		t.wins++
	}
	t.kills += w.kills
	t.traps += w.traps
	switch w.status {
	case statusOut:
		if w.treasure > 0 {
			t.rich++
			t.gold += w.treasure
		} else {
			t.poor++
		}
	case statusDead:
		t.dead++
	default:
		t.lost++
	}
}

func percent(n, of int) string {
	if of == 0 {
		return "-"
	}
	return fmt.Sprintf("%d %%", (n*100+of/2)/of)
}

func average(n, of int) string {
	if of == 0 {
		return "-"
	}
	return strings.Replace(fmt.Sprintf("%.1f", float64(n)/float64(of)), ".", ",", 1)
}

func writeSimSummary(out io.Writer, results []simResult) {
	var all simTally
	rounds, wakes, takes := 0, 0, 0
	byClass := map[string]*simTally{}
	byPers := map[string]*simTally{}
	byClassPers := map[[2]string]*simTally{}
	var persNames []string
	deaths := map[string]int{}
	for _, r := range results {
		rounds += r.rounds
		wakes += r.wakes
		for _, w := range r.walkers {
			all.add(w)
			takes += w.takes
			c := w.class
			if c == "" {
				c = "utan klass"
			}
			if byClass[c] == nil {
				byClass[c] = &simTally{}
			}
			byClass[c].add(w)
			if w.pers != "" {
				if byPers[w.pers] == nil {
					byPers[w.pers] = &simTally{}
					persNames = append(persNames, w.pers)
				}
				byPers[w.pers].add(w)
				key := [2]string{c, w.pers}
				if byClassPers[key] == nil {
					byClassPers[key] = &simTally{}
				}
				byClassPers[key].add(w)
			}
			if w.death != "" {
				deaths[w.death]++
			}
		}
	}
	n := len(results)
	fmt.Fprintln(out, "Vandrarna")
	rows := [][2]string{
		{"Kom ut med skatt", percent(all.rich, all.walkers)},
		{"Kom ut utan skatt", percent(all.poor, all.walkers)},
		{"Dog", percent(all.dead, all.walkers)},
		{"Blev kvar i mörkret", percent(all.lost, all.walkers)},
		{"Skatt hos dem med skatt", average(all.gold, all.rich)},
		{"Nävar ur kammaren", average(takes, all.walkers)},
		{"Fällda varelser", average(all.kills, all.walkers)},
		{"Fällor som slog till", average(all.traps, all.walkers)},
	}
	for _, r := range rows {
		fmt.Fprintf(out, "  %-26s %8s\n", r[0], r[1])
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Vandringarna")
	fmt.Fprintf(out, "  %-26s %8s\n", "Rundor", average(rounds, n))
	fmt.Fprintf(out, "  %-26s %8s\n", "Ormen vaknade", average(wakes, n))
	fmt.Fprintln(out)

	if len(deaths) > 0 {
		fmt.Fprintln(out, "Dödsorsaker")
		type cause struct {
			what string
			n    int
		}
		var list []cause
		for k, v := range deaths {
			list = append(list, cause{k, v})
		}
		sort.Slice(list, func(i, j int) bool {
			if list[i].n != list[j].n {
				return list[i].n > list[j].n
			}
			return list[i].what < list[j].what
		})
		for _, c := range list {
			fmt.Fprintf(out, "  %-30s %4s\n", c.what, percent(c.n, all.dead))
		}
		fmt.Fprintln(out)
	}

	sort.Strings(persNames)
	if len(byPers) > 0 {
		fmt.Fprintln(out, "Personligheterna")
		writeSimHeader(out, nil)
		for _, k := range persNames {
			writeSimRow(out, k, byPers[k], nil)
		}
		fmt.Fprintln(out)
	}

	// Klasserna, med andelen som kom ut med skatt för varje personlighet
	// längst till höger.
	fmt.Fprintln(out, "Klasserna")
	writeSimHeader(out, persNames)
	var names []string
	for k := range byClass {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, k := range names {
		var extra []string
		for _, pn := range persNames {
			t := byClassPers[[2]string{k, pn}]
			if t == nil {
				t = &simTally{}
			}
			extra = append(extra, percent(t.rich, t.walkers))
		}
		writeSimRow(out, k, byClass[k], extra)
	}
	for _, what := range []string{"guld per vandrare", "vann"} {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Klass och strategi: "+what)
		fmt.Fprintf(out, "  %-14s", "")
		for _, pn := range persNames {
			fmt.Fprintf(out, " %13s", pn)
		}
		fmt.Fprintln(out)
		for _, k := range names {
			fmt.Fprintf(out, "  %-14s", k)
			for _, pn := range persNames {
				t := byClassPers[[2]string{k, pn}]
				if t == nil {
					t = &simTally{}
				}
				v := average(t.gold, t.walkers)
				if what == "vann" {
					v = percent(t.wins, t.walkers)
				}
				fmt.Fprintf(out, " %13s", v)
			}
			fmt.Fprintln(out)
		}
	}
}

func writeSimHeader(out io.Writer, extra []string) {
	fmt.Fprintf(out, "  %-14s %8s %10s %6s %6s %8s %6s %8s", "", "vandrare", "med skatt", "dog", "skatt", "varelser", "vann", "guld/v")
	for _, e := range extra {
		fmt.Fprintf(out, " %10s", e)
	}
	fmt.Fprintln(out)
}

func writeSimRow(out io.Writer, name string, t *simTally, extra []string) {
	fmt.Fprintf(out, "  %-14s %8d %10s %6s %6s %8s %6s %8s", name, t.walkers,
		percent(t.rich, t.walkers), percent(t.dead, t.walkers),
		average(t.gold, t.rich), average(t.kills, t.walkers), percent(t.wins, t.walkers), average(t.gold, t.walkers))
	for _, e := range extra {
		fmt.Fprintf(out, " %10s", e)
	}
	fmt.Fprintln(out)
}
