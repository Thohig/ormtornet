package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// Testerna här kommer från provspelet 2026-11-22 (docs/provspel-2026-11-22.md).

// Den som gick ut tomhänt står före de döda och förlorade i resultatet.
func TestResultsEmptyHandedBeforeDead(t *testing.T) {
	g := testGame(t)
	for _, n := range []string{"Död", "Förlorad", "Tomhänt", "Rik"} {
		g.Join(testClient(), n)
	}
	status := map[string]string{"Död": statusDead, "Förlorad": statusLost, "Tomhänt": statusOut, "Rik": statusOut}
	for _, p := range g.players {
		if p != nil {
			p.Status = status[p.Name]
			if p.Name == "Rik" || p.Name == "Död" {
				p.Treasure = 30
			}
		}
	}
	var order []string
	for _, r := range g.results() {
		order = append(order, r.Name)
	}
	if got := strings.Join(order, " "); got != "Rik Tomhänt Förlorad Död" {
		t.Fatalf("ordningen blev %s", got)
	}
}

// Varje föremål har en obestämd form för texterna, på båda språken, så att
// det inte står "Sigge hittade extra veke".
func TestItemsHaveIndefiniteForm(t *testing.T) {
	for _, lang := range []string{langSv, langEn} {
		d, _, err := loadDataLang("data", lang)
		if err != nil {
			t.Fatal(err)
		}
		for _, it := range d.Foremal.Lista {
			if it.Hittad == "" {
				t.Errorf("%s: %s saknar hittad", lang, it.ID)
			}
		}
	}
}

// Svaret på ett drag mot en vägg är märkt, så att spelarsidan lägger det i
// rundan som står i bild.
func TestWallReplyIsMarked(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	room := &g.rooms[p.Y][p.X]
	wall := Dir(-1)
	for d := North; d <= West; d++ {
		if !room.Doors[d] && !(room.HasExit && room.Exit == d) {
			wall = d
			break
		}
	}
	if wall < 0 {
		t.Skip("rummet har ingen vägg")
	}
	for len(a.send) > 0 {
		<-a.send
	}
	g.Act(a, "NESW"[wall:wall+1], 0)
	found := false
	for len(a.send) > 0 {
		var m logMsg
		if json.Unmarshal(<-a.send, &m) == nil && m.T == "log" {
			found = found || m.Svar
		}
	}
	if !found {
		t.Fatal("svaret om väggen var inte märkt")
	}
}

// En runda fast säger "en runda", inte "1 rundor".
func TestStuckOneRound(t *testing.T) {
	g := testGame(t)
	if txt := g.data.Text("sam_fast_en"); txt == "" || strings.Contains(txt, "rundor") {
		t.Fatalf("sam_fast_en: %q", txt)
	}
	if txt := g.data.Text("sol_snart_en"); txt == "" || strings.Contains(txt, "rundor") {
		t.Fatalf("sol_snart_en: %q", txt)
	}
}

// Hoppar solen ner till en runda (klockfällan) säger varningen "en
// runda", inte "1 rundor".
func TestSunWarningOneRound(t *testing.T) {
	g, a, _ := startTwo(t)
	for len(a.send) > 0 {
		<-a.send
	}
	g.sun = 1
	g.sunWarning(4)
	found := false
	for len(a.send) > 0 {
		var m logMsg
		if json.Unmarshal(<-a.send, &m) == nil && m.T == "log" {
			if strings.Contains(m.Text, "1 rundor") {
				t.Fatalf("varningen säger %q", m.Text)
			}
			found = found || strings.Contains(m.Text, "En runda kvar")
		}
	}
	if !found {
		t.Fatal("varningen om en runda kom inte")
	}
}

// Med en kort sol (provspel med -sol 4) kommer ingen varning om halva
// vägen efter varningen om att solen nästan nuddar horisonten.
func TestShortSunWarningsInOrder(t *testing.T) {
	g, a, _ := startTwo(t)
	g.sunMax = 4
	var got []string
	for before := 4; before > 1; before-- {
		for len(a.send) > 0 {
			<-a.send
		}
		g.sun = before - 1
		g.sunWarning(before)
		for len(a.send) > 0 {
			var m logMsg
			if json.Unmarshal(<-a.send, &m) == nil && m.T == "log" && strings.Contains(m.Text, "Solen") {
				got = append(got, m.Text)
			}
		}
	}
	if len(got) != 1 || !strings.Contains(got[0], "nuddar") {
		t.Fatalf("varningarna: %q", got)
	}
}

// Den som laddar om sidan i slutet står kvar i resultatet och kommer
// tillbaka till sin egen vandring, inte som en nykomling.
func TestReloadOnEndScreenKeepsResult(t *testing.T) {
	g, a, b := startTwo(t)
	g.leaveWait = 0
	pa := a.player
	pa.X, pa.Y = 0, 0
	g.endGame()
	g.Leave(a)
	named := func() bool {
		for _, r := range g.results() {
			if r.Name == "Alva" {
				return true
			}
		}
		return false
	}
	if !named() {
		t.Fatal("Alva försvann ur resultatet när sidan laddades om")
	}
	// Bo trycker Enter medan Alva är borta: det nya tornet väntar inte på
	// den som har gått, men Alva står kvar tills tornet byggs.
	a2 := testClient()
	g.Join(a2, "Alva")
	if a2.player != pa || pa.Status == statusNew {
		t.Fatalf("Alva kom tillbaka som en ny vandrare: %+v", a2.player)
	}
	if st := g.stateFor(pa); st.Journey == nil {
		t.Fatal("Alva fick inte tillbaka sin vandring")
	}
	g.Ready(a2)
	g.Ready(b)
	if g.phase == phaseEnd {
		t.Fatal("nästa torn startade inte")
	}
}

// Den som lämnar i slutet och inte kommer tillbaka är inte med i nästa torn.
func TestLeftOnEndScreenDroppedFromNextTower(t *testing.T) {
	g, a, b := startTwo(t)
	g.leaveWait = 0
	g.cfg.MinPlayers = 1
	g.endGame()
	g.Leave(a)
	g.Ready(b)
	for _, p := range g.players {
		if p != nil && p.Name == "Alva" {
			t.Fatalf("Alva är kvar i fas %s", g.phase)
		}
	}
}

// Den som laddar om sidan mitt i vandringen får rummet hen står i igen.
func TestReloadMidGameRetellsRoom(t *testing.T) {
	g, a, _ := startTwo(t)
	pa := a.player
	g.Leave(a)
	a2 := testClient()
	g.Join(a2, "Alva")
	if a2.player != pa {
		t.Fatal("Alva fick inte tillbaka sin vandrare")
	}
	found := false
	for len(a2.send) > 0 {
		var m roomMsg
		if json.Unmarshal(<-a2.send, &m) == nil && m.T == "rum" && m.Text != "" {
			found = true
		}
	}
	if !found {
		t.Fatal("rummet skickades inte igen")
	}
}

// I mörkret ser man inte dörrarna eller porten i rummet man står i, så
// kartan visar inga tangenter; ett steg mot en vägg kostar rundan.
func TestDarkHidesDoors(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	room := &g.rooms[p.Y][p.X]
	if !room.HasExit {
		t.Fatal("vandraren står inte i portrummet")
	}
	p.Dark = 2
	st := g.stateFor(p)
	if len(st.Rooms) != 1 {
		t.Fatalf("i mörkret ska bara det egna rummet synas, inte %d", len(st.Rooms))
	}
	rv := st.Rooms[0]
	if rv.D != [4]bool{} || rv.Exit != -1 {
		t.Fatalf("dörrarna eller porten syns i mörkret: %+v", rv)
	}
	p.Dark = 0
	if rv := g.stateFor(p).Rooms; len(rv) == 0 || rv[0].Exit < 0 {
		t.Fatal("porten syns inte i ljuset")
	}
}

// Dricker, flyr eller väntar man medan en varelse inte anfaller står det
// att den låter en vara, och man förlorar inget.
func TestFreeAttackCalmWhenNotAttacking(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	typ := g.data.Monster.Typer[0]
	calm, hit := 0, 0
	for i := 0; i < 60; i++ {
		for len(a.send) > 0 {
			<-a.send
		}
		m := &Monster{ID: 90 + i, Type: typ.ID, X: p.X, Y: p.Y, HP: 2}
		p.HP = 50
		key := []string{"inte_beredd", "fly", "garden"}[i%3]
		g.freeAttack(p, m, key, 0, false)
		for len(a.send) > 0 {
			var msg logMsg
			if json.Unmarshal(<-a.send, &msg) != nil || msg.T != "log" {
				continue
			}
			if strings.Contains(msg.Text, "försöker") {
				t.Fatalf("fritt drag: %q", msg.Text)
			}
			if strings.Contains(msg.Text, "låter dig vara") || strings.Contains(msg.Text, "följer inte efter") || strings.Contains(msg.Text, "väntar också") || strings.Contains(msg.Text, "inte rör sig") {
				calm++
				if p.HP != 50 {
					t.Fatalf("skadad fast varelsen inte anföll: %d", p.HP)
				}
			}
		}
		if p.HP < 50 {
			hit++
		}
	}
	if calm == 0 || hit == 0 {
		t.Fatalf("lugna %d, träffar %d", calm, hit)
	}
}

// Kartritaren får inte höra att kartan bara visar de egna rummen.
func TestCartographerIntro(t *testing.T) {
	g, a, b := startTwo(t)
	g.cfg.Classes = true
	a.player.classID = "kartritare"
	b.player.classID = "legoknekt"
	if k := g.visibilityKey(a.player); k != "synlighet_kartritare" {
		t.Fatalf("kartritaren: %s", k)
	}
	if k := g.visibilityKey(b.player); k != "synlighet_egen" {
		t.Fatalf("legoknekten: %s", k)
	}
}
