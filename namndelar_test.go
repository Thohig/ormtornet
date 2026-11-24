package main

import (
	"strings"
	"testing"
)

// Varje källa ger sin del, en gång, och bara till den som hittar den.
func TestPartsFromThreeSources(t *testing.T) {
	g, cs := nameGame(t, 2)
	a, b := cs[0].player, cs[1].player
	seen := map[int]bool{}
	for src := range 3 {
		seen[g.partAt[src]] = true
	}
	if len(seen) != 3 {
		t.Fatalf("delarna ligger inte på tre ställen: %v", g.partAt)
	}
	for src := range 3 {
		logs(cs[0])
		g.learnPart(a, src)
		if !a.parts[g.partAt[src]] || len(logs(cs[0])) == 0 {
			t.Fatalf("källan %s gav ingen del", partSources[src])
		}
		g.learnPart(a, src)
		if len(logs(cs[0])) != 0 {
			t.Fatal("samma del två gånger")
		}
	}
	if b.partCount() != 0 {
		t.Fatalf("den andre fick %d delar", b.partCount())
	}
	// Tecknen står i den ordning de hittades, inte i namnets.
	want := ""
	for src := range 3 {
		want += strings.ToUpper(g.serpentName[g.partAt[src]])
	}
	if st := g.stateFor(a); strings.Join(st.NameParts, "") != want {
		t.Fatalf("klienten: %v, ville ha %s", st.NameParts, want)
	}
}

// Ristningen hamnar bara i ett rum en bit in i tornet, och den som står där
// läser den.
func TestCarving(t *testing.T) {
	g, cs := nameGame(t, 1)
	p := cs[0].player
	c := g.corners()[0]
	putInRoom(g, p, c.x+1, c.y)
	if g.carvingFits(c.x+1, c.y) {
		t.Fatal("ristningen får inte ligga vid porten")
	}
	mid := g.center() - 1
	putInRoom(g, p, mid, g.center())
	if !g.carvingFits(mid, g.center()) {
		t.Fatal("ett rum nära mitten ska passa")
	}
	g.carving = carving{x: mid, y: g.center(), set: true}
	g.readCarving(p)
	if !p.parts[g.partAt[srcCarving]] {
		t.Fatal("ristningen lästes inte")
	}
}

// Finns ingen ristning på nivå 3 läggs den i ett känt rum, och den som inte
// har läst den får veta åt vilket håll den finns.
func TestCarvingOnLevelThree(t *testing.T) {
	g, cs := nameGame(t, 1)
	p := cs[0].player
	putInRoom(g, p, g.center()-1, g.center())
	g.reveal(g.center()-2, g.center())
	g.rooms[g.center()][g.center()-2].Kind = "rum"
	g.carving = carving{} // rummen ovan kan ha fått den
	logs(cs[0])
	g.level = 3
	g.partsOnLevelUp()
	if !g.carving.set || g.playerIn(g.carving.x, g.carving.y) {
		t.Fatalf("ristningen: %+v", g.carving)
	}
	if !strings.Contains(strings.Join(logs(cs[0]), "\n"), g.data.Horn[6]) {
		t.Fatal("ingen riktning till ristningen")
	}
}

// Den som träffar en varelse ser tecknet under fjällen.
func TestCreaturePart(t *testing.T) {
	g, cs := nameGame(t, 1)
	p := cs[0].player
	putInRoom(g, p, 4, 3)
	typ := &g.data.Monster.Typer[0]
	m := &Monster{ID: 50, Type: typ.ID, X: 4, Y: 3, HP: 10}
	g.monsters = append(g.monsters, m)
	sureHits(g)
	typ.Forsvar = 1 // Kraft 1 mot Försvar 1: slaget tar inte
	g.duel(p, m)
	if p.parts[g.partAt[srcCreature]] {
		t.Fatal("delen utan att träffa")
	}
	typ.Forsvar = 0
	g.duel(p, m)
	if !p.parts[g.partAt[srcCreature]] {
		t.Fatal("ingen del när vandraren träffade")
	}
}

// En näve ur kammaren ger ormens del.
func TestChamberPart(t *testing.T) {
	g, cs := nameGame(t, 1)
	p := cs[0].player
	c := g.center()
	g.reveal(c, c)
	p.X, p.Y = c, c
	g.take(p)
	if !p.parts[g.partAt[srcChamber]] {
		t.Fatal("ingen del ur kammaren")
	}
}

// Texten säger var på ormen delen sitter, inte vilken i ordningen den är.
func TestPartTextGivesPlace(t *testing.T) {
	g, _ := nameGame(t, 1)
	g.serpentName = []string{"Ka", "ru", "meth"}
	for i := range 3 {
		text := g.namePartText(srcCarving, i)
		if !strings.Contains(text, strings.ToUpper(g.serpentName[i])) {
			t.Fatalf("delen saknas: %q", text)
		}
		found := false
		for _, pl := range g.data.Avslojande.Platser[i] {
			found = found || strings.Contains(text, pl)
		}
		if !found || strings.Contains(text, "{") {
			t.Fatalf("platsen saknas: %q", text)
		}
	}
}

// Rätt delar i fel ordning: tornet säger det, försöket är förbrukat och
// ormen blir inte oroligare. En fel viskning gör ormen oroligare.
func TestWhisperWrongOrderAndWrong(t *testing.T) {
	g, cs := nameGame(t, 1)
	c := cs[0]
	g.serpentName = []string{"Ka", "ru", "meth"}
	g.level = 2
	wake := countWake(g)
	logs(c)
	g.Whisper(c, "tornet är Rukameth")
	if c.player.initiated || c.player.guessLevel != 2 || !contains(logs(c), g.data.Text("viskning_ordning")) {
		t.Fatal("fel ordning ska förbruka försöket och sägas")
	}
	if countWake(g) != wake {
		t.Fatal("fel ordning gjorde ormen oroligare")
	}
	g.level = 3
	g.Whisper(c, "tornet är en ödla")
	if c.player.initiated || countWake(g) != wake+1 {
		t.Fatalf("fel viskning: %d vakna-kort, ville ha %d", countWake(g), wake+1)
	}
	g.level = 4
	g.Whisper(c, "ka ru meth, tornet är ormen")
	if !c.player.initiated {
		t.Fatal("rätt viskning med mellanrum gick inte")
	}
}
