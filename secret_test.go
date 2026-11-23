package main

import (
	"strings"
	"testing"
)

func TestWhisperMatching(t *testing.T) {
	answers := []string{"tornet är ormen", "inuti ormen"}
	for _, c := range []struct {
		w  string
		ok bool
	}{
		{"Tornet ÄR ormen!", true},
		{"jag tror att tornet är ormen", true},
		{"vi är inuti ormen", true},
		{"ormen sover i tornet", false},
		{"tornet är ormens", false},
		{"", false},
	} {
		if got := whisperMatches(c.w, answers); got != c.ok {
			t.Errorf("%q: %v", c.w, got)
		}
	}
}

func TestWhisperInitiates(t *testing.T) {
	g, a, b := startTwo(t)
	pa := a.player
	g.serpentName = []string{"Ka", "ru", "meth"}
	g.Whisper(a, "Karumeth, tornet är ormen")
	if pa.initiated {
		t.Fatal("invigd på nivå 1")
	}
	g.level = 2
	logs(b)
	g.Whisper(a, "tornet är ett torn")
	if pa.initiated || pa.guessLevel != 2 {
		t.Fatal("fel viskning gick igenom")
	}
	g.Whisper(a, "Karumeth, tornet är ormen")
	if pa.initiated {
		t.Fatal("två viskningar på samma nivå")
	}
	g.level = 3
	g.Whisper(a, "Karumeth, tornet är ormen.")
	if !pa.initiated {
		t.Fatal("rätt viskning gick inte igenom")
	}
	for _, line := range logs(b) {
		if strings.Contains(line, "ormen") || strings.Contains(line, "Karumeth") || strings.Contains(line, pa.Name) {
			t.Fatalf("Bo fick veta något: %q", line)
		}
	}
	if st := g.stateFor(pa); !st.Initiated || st.Heart == "" || st.CanWhisper {
		t.Fatalf("klienten: %+v", st.Initiated)
	}
}

func TestHeartbeatPointsToChamber(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.initiated = true
	putInRoom(g, p, 1, 1)
	if want := g.data.Text("hjartslag", "horn", "sydöst"); g.heartbeat(p) != want {
		t.Fatalf("%q", g.heartbeat(p))
	}
	p.X, p.Y = g.center(), g.center()
	if g.heartbeat(p) != g.data.Text("hjartslag_har") {
		t.Fatal("inget dunk under kammaren")
	}
}

func TestPortWaitsForInitiated(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	pa.initiated = true
	putInRoom(g, pa, 4, 3)
	putInRoom(g, pb, 2, 3)
	g.sun = 1
	g.resolve()
	if g.phase != phasePlay || pa.Status != statusIn || pb.Status != statusLost {
		t.Fatalf("fas %s, Alva %s, Bo %s", g.phase, pa.Status, pb.Status)
	}
	g.resolve()
	if g.phase != phaseEnd || pa.Status != statusLost {
		t.Fatal("porten väntade mer än en runda")
	}
}

func TestInitiatedBenefits(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.initiated = true
	g.data.Avslojande.Invigd.Falla = 1
	// Fällan i ett nytt rum: den invigde kliver förbi.
	putInRoom(g, p, 4, 3)
	g.rooms[3][4].Doors[East] = true
	g.rooms[3][5] = Room{}
	g.reveal(5, 3).Event = "pilar"
	g.rooms[3][5].Revealed = false
	g.rooms[3][5].Event = "pilar"
	hp := p.HP
	g.move(p, East, 0)
	if p.HP != hp || g.rooms[3][5].Mark == "" {
		t.Fatalf("liv %d av %d, tecken %q", p.HP, hp, g.rooms[3][5].Mark)
	}

	// Lugna ormen i kammaren, en gång.
	p.X, p.Y = g.center(), g.center()
	g.reveal(p.X, p.Y)
	before := g.wakeCardsLeft()
	if !g.canWhisper(p) {
		t.Fatal("kan inte viska i kammaren")
	}
	g.Whisper(a, "sov")
	if g.wakeCardsLeft() != before-1 || !p.calmed || g.canWhisper(p) {
		t.Fatalf("vakna-kort %d, före %d", g.wakeCardsLeft(), before)
	}

	// Guld när den invigde kommer ut.
	cn := g.corners()[0]
	p.X, p.Y, p.Treasure = cn.x, cn.y, 100
	g.reveal(cn.x, cn.y)
	g.move(p, cn.exit, 0)
	if p.Status != statusOut || p.Treasure != 125 {
		t.Fatalf("%s med %d", p.Status, p.Treasure)
	}
}

func TestInitiatedWakesSerpentMoreEasily(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	p.initiated = true
	g.data.Avslojande.Invigd.Vacker = 1
	p.X, p.Y = g.center(), g.center()
	g.reveal(p.X, p.Y)
	wakes := g.wakes
	g.take(p)
	if g.wakes != wakes+1 {
		t.Fatal("ormen vaknade inte")
	}
}

// nameGame startar en vandring med n människor.
func nameGame(t *testing.T, n int) (*Game, []*Client) {
	t.Helper()
	g := testGame(t)
	g.cfg.MinPlayers = 1
	var cs []*Client
	for i := range n {
		c := testClient()
		g.Join(c, string(rune('A'+i))+"lva")
		cs = append(cs, c)
	}
	if !g.ForceStart() || len(g.serpentName) != 3 {
		t.Fatalf("vandringen startade inte, eller inget namn: %v", g.serpentName)
	}
	return g, cs
}

func TestNamePartsDealt(t *testing.T) {
	for _, n := range []int{1, 2, 8} {
		g, cs := nameGame(t, n)
		for level := 2; level <= 4; level++ {
			g.level = level
			g.dealNameParts()
			for _, c := range cs {
				p := c.player
				if got, want := p.partCount(), min(level-1, 3); got != want {
					t.Fatalf("%d människor, nivå %d: %s har %d delar, ville ha %d", n, level, p.Name, got, want)
				}
				if st := g.stateFor(p); len(st.NameParts) != 3 {
					t.Fatalf("klienten fick inte delarna: %v", st.NameParts)
				}
			}
		}
	}
}

func TestNameWhisper(t *testing.T) {
	g, _ := nameGame(t, 1)
	g.serpentName = []string{"Ka", "ru", "meth"}
	for _, c := range []struct {
		w  string
		ok bool
	}{
		{"Karumeth, tornet är ormen", true},
		{"tornet är Karumeth!", true},
		{"tornet är ka-ru-meth", true},
		{"tornet är KA RU METH", true},
		{"tornet är ormen", false},
		{"Karumeth", false},
		{"tornet är Karumet", false},
		{"tornet är karumeths", false},
	} {
		if got := g.whisperRight(c.w); got != c.ok {
			t.Errorf("%q: %v", c.w, got)
		}
	}
	// Saknar datafilerna namn räcker svaret.
	g.serpentName = nil
	if !g.whisperRight("tornet är ormen") {
		t.Fatal("svaret gick inte igenom utan namn")
	}
	// Epilogen säger namnet.
	g.serpentName = []string{"Ka", "ru", "meth"}
	if g.serpentNameText() != "Karumeth" {
		t.Fatalf("namnet skrivs %q", g.serpentNameText())
	}
}

func TestKillTakesNameParts(t *testing.T) {
	g, cs := nameGame(t, 2)
	a, b := cs[0].player, cs[1].player
	a.parts = [3]bool{true, false, false}
	b.parts = [3]bool{false, true, true}
	g.takeNameParts(a, b)
	if a.parts != [3]bool{true, true, true} {
		t.Fatalf("dödaren fick %v", a.parts)
	}
	// Utan namn finns inga delar att ta.
	g2, cs2 := nameGame(t, 2)
	g2.serpentName = nil
	c, d := cs2[0].player, cs2[1].player
	d.parts = [3]bool{true, true, true}
	g2.takeNameParts(c, d)
	if c.partCount() != 0 {
		t.Fatal("delar fast vandringen saknar namn")
	}
}

func countWake(g *Game) int {
	n := 0
	for _, id := range g.deck {
		if c := g.data.card(id); c != nil && c.Vaknar {
			n++
		}
	}
	return n
}

func TestFirstInitiatedGetsMoreGold(t *testing.T) {
	g, cs := nameGame(t, 2)
	g.serpentName = []string{"Ka", "ru", "meth"}
	g.level = 2
	a, b := cs[0], cs[1]
	g.Whisper(a, "tornet är Karumeth")
	g.Whisper(b, "tornet är Karumeth")
	if !a.player.initiated || !b.player.initiated || g.firstInitiated != a.player.ID {
		t.Fatalf("invigda %v %v, först %d", a.player.initiated, b.player.initiated, g.firstInitiated)
	}
	heard := func(c *Client) bool {
		return contains(logs(c), g.data.Text("invigd_forst"))
	}
	if !heard(a) || heard(b) {
		t.Fatal("bara den förste ska få höra att hen var först")
	}
	// Båda går ut genom porten de står vid, med 100 guld.
	gold := func(c *Client) int {
		p := c.player
		p.Treasure = 100
		room := &g.rooms[p.Y][p.X]
		if !room.HasExit {
			t.Fatalf("%s står inte vid en port", p.Name)
		}
		g.move(p, room.Exit, 0)
		if p.Status != statusOut {
			t.Fatalf("%s kom inte ut", p.Name)
		}
		return p.Treasure
	}
	inv := g.data.Avslojande.Invigd
	if got, want := gold(a), 100+int(100*(inv.Guld+inv.ForstGuld)); got != want {
		t.Fatalf("den förste fick %d, ville ha %d", got, want)
	}
	if got, want := gold(b), 100+int(100*inv.Guld); got != want {
		t.Fatalf("den andre fick %d, ville ha %d", got, want)
	}
}

// Bara namnet förbrukar inte försöket: tornet känner igen det och frågar
// vad man vill säga om det.
func TestWhisperOnlyNameKeepsTheTry(t *testing.T) {
	g, a, _ := startTwo(t)
	pa := a.player
	g.serpentName = []string{"Ka", "ru", "meth"}
	g.level = 2
	for _, w := range []string{"Karumeth", "ka ru meth!", "KARUMETH."} {
		logs(a)
		g.Whisper(a, w)
		if pa.initiated || pa.guessLevel != 0 || !g.canWhisper(pa) {
			t.Fatalf("%q: invigd %v, försök %d", w, pa.initiated, pa.guessLevel)
		}
		if !contains(logs(a), g.data.Text("viskning_bara_namn")) {
			t.Fatalf("%q gav inget svar om namnet", w)
		}
	}
	g.Whisper(a, "tornet är Karumeth")
	if !pa.initiated {
		t.Fatal("rätt viskning efter namnet gick inte igenom")
	}
}

// Var och en drar ledtrådar ur sin egen lek: alla får till slut alla
// texterna på nivån, oberoende av de andra.
func TestCluesFromOwnDeck(t *testing.T) {
	g, cs := nameGame(t, 3)
	g.level = 2
	list := g.data.Avslojande.Ledtradar["2"]
	for _, c := range cs {
		logs(c)
	}
	for range list {
		g.dealClues(g.data.Avslojande.Ledtradar, "ledtrad", "info")
	}
	for _, c := range cs {
		got := logs(c)
		for _, want := range list {
			if !contains(got, want) {
				t.Fatalf("%s fick inte %q", c.player.Name, want)
			}
		}
	}
}
