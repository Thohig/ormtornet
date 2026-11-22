package main

import (
	"strings"
	"testing"
)

// fightReady ger två vandrare i samma vanliga rum, med strid påslagen.
func fightReady(t *testing.T) (*Game, *Client, *Client) {
	t.Helper()
	g, a, b := twoPlaying(t, 0)
	g.cfg.PvP = pvpOn
	putInRoom(g, a.player, 4, 1)
	putInRoom(g, b.player, 4, 1)
	g.deck = []string{"sover", "sover", "sover"}
	return g, a, b
}

func wakeCards(g *Game) int {
	n := 0
	for _, id := range g.deck {
		if c := g.data.card(id); c != nil && c.Vaknar {
			n++
		}
	}
	return n
}

func TestPvPExchangeStealsAThird(t *testing.T) {
	g, a, b := fightReady(t)
	alva, bo := a.player, b.player
	bo.Treasure = 30
	hp := bo.HP
	g.Act(a, actAttack, -bo.ID)
	g.Act(b, actDodge, -alva.ID) // anfall slår undanmanöver
	if bo.HP != hp-pvpSkada || bo.Treasure != 20 || alva.Treasure != 10 {
		t.Fatalf("Bo %d liv %d guld, Alva %d guld", bo.HP, bo.Treasure, alva.Treasure)
	}
	if wakeCards(g) != 1 {
		t.Fatalf("striden lade %d vakna-kort i leken, ville ha 1", wakeCards(g))
	}
}

func TestPvPAmbushAndKill(t *testing.T) {
	g, a, b := fightReady(t)
	alva, bo := a.player, b.player
	bo.Treasure = 50
	bo.HP = 1
	g.Act(a, actAttack, -bo.ID)
	g.Act(b, actWait, 0) // Bo försvarar sig inte
	if bo.Status != statusDead || bo.Treasure != 0 || alva.Treasure != 50 {
		t.Fatalf("Bo %s med %d guld, Alva %d guld", bo.Status, bo.Treasure, alva.Treasure)
	}
}

func TestPvPParryAgainstNoAttack(t *testing.T) {
	g, a, b := fightReady(t)
	bo := b.player
	hp := bo.HP
	g.Act(a, actParry, -bo.ID)
	g.Act(b, actWait, 0)
	if bo.HP != hp {
		t.Fatal("en parering skadade någon som inte anföll")
	}
}

func TestPvPRules(t *testing.T) {
	g, a, _ := fightReady(t)
	alva := a.player
	try := func() bool {
		alva.Action = ""
		g.actPlayer(alva, actAttack, 0)
		return alva.Action == actAttack
	}
	if !try() {
		t.Fatal("kunde inte anfalla med strid påslagen")
	}
	g.cfg.PvP = pvpOff
	if try() {
		t.Fatal("anföll fast strid var avstängd")
	}
	g.cfg.PvP = pvpOn
	g.rooms[alva.Y][alva.X].Kind = "ingang"
	if try() {
		t.Fatal("anföll i ett fredat portrum")
	}
}

func TestBotsFightEachOtherOnlyWithoutHumans(t *testing.T) {
	g := testGame(t)
	g.cfg.PvP = pvpOn
	g.cfg.AI = 2
	g.AddBots(2)
	var bots []*Player
	for _, p := range g.players {
		if p != nil && p.bot != nil {
			bots = append(bots, p)
			p.Status = statusIn
			putInRoom(g, p, 4, 1)
		}
	}
	h := testClient()
	g.Join(h, "Alva")
	if g.botFoe(bots[0]) != nil {
		t.Fatal("en AI-spelare gav sig på en annan AI-spelare fast en människa var med")
	}
	g.Leave(h)
	// I lobbyn byggs tornet om när någon går, så ställ AI-spelarna i rummet igen.
	for _, p := range bots {
		p.Status = statusIn
		putInRoom(g, p, 4, 1)
	}
	if g.botFoe(bots[0]) != bots[1] {
		t.Fatal("utan människor ska AI-spelarna kunna slåss med varandra")
	}
}

// chronTexts är krönikans händelser som text.
func chronTexts(g *Game) []string {
	var out []string
	for _, e := range g.chron {
		out = append(out, e.text)
	}
	return out
}

func TestPvPHitsInChronicle(t *testing.T) {
	g, a, b := fightReady(t)
	alva, bo := a.player, b.player
	bo.Treasure = 30
	g.Act(a, actAttack, -bo.ID)
	g.Act(b, actDodge, -alva.ID)
	want := g.data.Text("kronika_slag_byte", "namn", alva.Name, "offer", bo.Name, "guld", "10")
	if got := chronTexts(g); len(got) != 1 || got[0] != want {
		t.Fatalf("krönikan efter duellen: %q, ville ha %q", got, want)
	}

	g.chron = nil
	alva.Treasure = 0
	g.Act(b, actAttack, -alva.ID)
	g.Act(a, actWait, 0)
	want = g.data.Text("kronika_bakhall", "namn", bo.Name, "offer", alva.Name)
	if got := chronTexts(g); len(got) != 1 || got[0] != want {
		t.Fatalf("krönikan efter bakhållet: %q, ville ha %q", got, want)
	}
}

func TestPvPKillOnlyOnceInChronicle(t *testing.T) {
	g, a, b := fightReady(t)
	bo := b.player
	bo.Treasure = 50
	bo.HP = 1
	g.Act(a, actAttack, -bo.ID)
	g.Act(b, actWait, 0)
	for _, s := range chronTexts(g) {
		if s == g.data.Text("kronika_bakhall_byte", "namn", a.player.Name, "offer", bo.Name, "guld", "50") {
			t.Fatalf("slaget som dödade stod också som ett slag: %q", chronTexts(g))
		}
	}
}

func TestPvPStatsLine(t *testing.T) {
	g, a, b := fightReady(t)
	alva, bo := a.player, b.player
	bo.Treasure = 30
	g.Act(a, actAttack, -bo.ID)
	g.Act(b, actWait, 0)
	if alva.stats.hits != 1 || alva.stats.stolen != 10 {
		t.Fatalf("Alva: %d träffar, %d stulet", alva.stats.hits, alva.stats.stolen)
	}
	strid := g.data.Text("kronika_vandrare_strid", "traffar", "1", "stulet", "10")
	if !strings.HasSuffix(g.statsLine(alva), strid) {
		t.Fatalf("Alvas rad saknar striderna: %q", g.statsLine(alva))
	}
	if strings.Contains(g.statsLine(bo), strid) {
		t.Fatalf("Bos rad har strider fast Bo inte slog: %q", g.statsLine(bo))
	}
}

// Tornet väljer motståndare: den som anföll en förra rundan, annars den som
// kom först in i tornet. En människa kan inte peka ut någon själv, och den
// som bär mest guld väljs inte för det, eftersom skatten är hemlig.
func TestTowerPicksTheFoe(t *testing.T) {
	g := testGame(t)
	a, b, c := testClient(), testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	g.Join(c, "Cia")
	g.Ready(a)
	g.Ready(b)
	g.Ready(c)
	if g.phase != phasePlay {
		t.Fatal("vandringen startade inte")
	}
	g.cfg.PvP = pvpOn
	alva, bo, cia := a.player, b.player, c.player
	for _, p := range []*Player{alva, bo, cia} {
		putInRoom(g, p, 4, 1)
	}
	cia.Treasure = 99

	if f := g.pickFoe(bo, -cia.ID); f != alva {
		t.Fatalf("Bo fick %s, ville ha Alva (först in i tornet)", f.Name)
	}
	bo.attackedBy, bo.attackedRound = cia.ID, g.round
	if f := g.pickFoe(bo, 0); f != cia {
		t.Fatalf("Bo fick %s, ville ha Cia som anföll", f.Name)
	}
	bo.attackedRound = g.round - 2
	if f := g.pickFoe(bo, 0); f != alva {
		t.Fatalf("Bo fick %s, ville ha Alva när anfallet var länge sedan", f.Name)
	}
}
