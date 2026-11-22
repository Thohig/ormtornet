package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

// rooms ger rumsblocken som klienten har fått.
func rooms(c *Client) []roomMsg {
	var out []roomMsg
	for {
		select {
		case b := <-c.send:
			var m roomMsg
			if json.Unmarshal(b, &m) == nil && m.T == "rum" {
				out = append(out, m)
			}
		default:
			return out
		}
	}
}

func TestRoomBlockAtStart(t *testing.T) {
	g := testGame(t)
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	rooms(a)
	g.Ready(a)
	g.Ready(b)
	got := rooms(a)
	if len(got) != 1 || got[0].Text == "" || got[0].Har == "" || got[0].Kvar {
		t.Fatalf("vid starten väntades ett rumsblock med beskrivning, fick %+v", got)
	}
	// Början av vandringen står i blocket: porten först, solen bland det
	// man ser och den andra vandrarens port bland det man hör.
	port := g.data.Text("start_port", "horn", g.data.Horn[a.player.Corner])
	if !strings.HasPrefix(got[0].Text, port+" ") {
		t.Errorf("beskrivningen började inte med porten: %q", got[0].Text)
	}
	if !strings.Contains(got[0].Har, strconv.Itoa(g.sun)) {
		t.Errorf("solen stod inte bland det man ser: %q", got[0].Har)
	}
	if got[0].Hor != g.data.Text("start_hor_en") {
		t.Errorf("väntade att höra den andra porten, hörde %q", got[0].Hor)
	}
	// Solen sägs bara i det första blocket.
	putInRoom(g, a.player, 4, 4)
	g.tellRoom(a.player)
	if got := rooms(a); len(got) != 1 || strings.Contains(got[0].Har, strconv.Itoa(g.sun)+" ") {
		t.Errorf("solen kom med i nästa block: %+v", got)
	}
}

func TestRoomBlockTellsWhatIsHere(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	putInRoom(g, pa, 4, 4)
	putInRoom(g, pb, 4, 4)
	rooms(a)

	g.tellRoom(pa)
	got := rooms(a)
	want := g.data.Text("har_vandrare", "namn", "Bo")
	if len(got) != 1 || got[0].Kvar || got[0].Har != want {
		t.Fatalf("väntade ett block med %q, fick %+v", want, got)
	}

	g.tellRoom(pa)
	if got := rooms(a); len(got) != 0 {
		t.Fatalf("inget har ändrats, men fick %+v", got)
	}

	pb.X = 5
	g.tellRoom(pa)
	got = rooms(a)
	want = g.data.Text("har_vandrare_borta", "namn", "Bo")
	if len(got) != 1 || !got[0].Kvar || got[0].Text != "" || got[0].Har != want {
		t.Fatalf("väntade bara %q, fick %+v", want, got)
	}

	// Varelsen beskrivs med sin egen text, och samma text så länge den
	// står kvar.
	typ := g.data.Monster.Typer[0]
	g.monsters = append(g.monsters, &Monster{ID: 1, Type: typ.ID, X: 4, Y: 4, HP: 1})
	g.tellRoom(pa)
	got = rooms(a)
	want = typ.DykerUpp[1%len(typ.DykerUpp)]
	if len(got) != 1 || got[0].Har != want {
		t.Fatalf("väntade %q, fick %+v", want, got)
	}
	g.tellRoom(pa)
	if got := rooms(a); len(got) != 0 {
		t.Fatalf("varelsen står kvar, men fick %+v", got)
	}
}

// Ett tomt rum säger ingenting, och en fälla som har slagit till syns som
// ett spår utan sitt namn.
func TestRoomBlockShowsNotTells(t *testing.T) {
	g, a, _ := startTwo(t)
	pa := a.player
	r := putInRoom(g, pa, 3, 3)
	pa.told = told{}
	rooms(a)
	g.tellRoom(pa)
	got := rooms(a)
	if len(got) != 1 || got[0].Har != "" {
		t.Fatalf("ett tomt rum ska inte säga något, fick %+v", got)
	}
	trap := g.data.trap("roterande")
	r.Mark = trap.Tecken
	g.tellRoom(pa)
	got = rooms(a)
	if len(got) != 1 || got[0].Har != trap.SparEfter[0] {
		t.Fatalf("väntade spåret %q, fick %+v", trap.SparEfter[0], got)
	}
}

func TestRoomBlockWaitsInTheDark(t *testing.T) {
	g, a, _ := startTwo(t)
	pa := a.player
	putInRoom(g, pa, 4, 4)
	pa.Dark = 2
	rooms(a)
	g.tellRoom(pa)
	if got := rooms(a); len(got) != 0 {
		t.Fatalf("i mörkret ska inget block komma, fick %+v", got)
	}
}

func TestJourneyRoundHasRoom(t *testing.T) {
	g, a, b := startTwo(t)
	pa := a.player
	putInRoom(g, pa, 4, 4)
	putInRoom(g, b.player, 2, 2)
	g.describeLater(pa, &g.rooms[4][4], 4, 4)
	step(g, a, b, actWait, actWait)
	if len(pa.journal) == 0 {
		t.Fatal("ingen runda i vandringen")
	}
	last := pa.journal[len(pa.journal)-1]
	if last.Room == "" {
		t.Fatalf("rundan i vandringen saknar rummet: %+v", last)
	}
}

// Det vandraren hör under rundan står sist i rumsblocket och inte som egna
// rader; står hen kvar kommer ett block bara med ljuden.
func TestRoomBlockCarriesWhatIsHeard(t *testing.T) {
	g, a, _ := startTwo(t)
	pa := a.player
	putInRoom(g, pa, 3, 3)
	pa.told = told{}
	g.tellRoom(pa)
	rooms(a)
	g.hearing = true
	g.hearTo(pa, "Stridsljud från norr.", "dim")
	g.hearTo(pa, "Stridsljud från norr.", "dim")
	g.hearing = false
	if l := logs(a); len(l) != 0 {
		t.Fatalf("ljudet kom som egna rader: %v", l)
	}
	g.tellRoom(pa)
	got := rooms(a)
	if len(got) != 1 || !got[0].Kvar || got[0].Hor != "Stridsljud från norr." {
		t.Fatalf("väntade ett block med ljudet en gång, fick %+v", got)
	}
}

// Rummen har inga namn för spelarna: rumsblocket och texterna har bara
// beskrivningen, och varje rum har en kort beskrivning på båda språken.
func TestRoomsHaveNoNamesForPlayers(t *testing.T) {
	g := testGame(t)
	types := append([]RoomType{g.data.Ingang, g.data.Skattkammare}, g.data.Rumstyper...)
	for _, rt := range types {
		if rt.Kort == "" {
			t.Errorf("rummet %q har ingen kort beskrivning", rt.ID)
		}
	}
	a, b := testClient(), testClient()
	g.Join(a, "Alva")
	g.Join(b, "Bo")
	g.Ready(a)
	g.Ready(b)
	for _, raw := range drainAll(a) {
		for _, rt := range types {
			if rt.Namn != "" && strings.Contains(raw, rt.Namn) {
				t.Errorf("spelarsidan fick rummets namn %q: %s", rt.Namn, raw)
			}
		}
	}
}

// drainAll ger allt klienten har fått, som text.
func drainAll(c *Client) []string {
	var out []string
	for {
		select {
		case b := <-c.send:
			out = append(out, string(b))
		default:
			return out
		}
	}
}

// Fällorna har inga namn för spelarna: den som går i en fälla får texten
// om vad som händer, men aldrig fällans namn.
func TestTrapsHaveNoNamesForPlayers(t *testing.T) {
	for _, trap := range testGame(t).data.Fallor {
		g, a, _ := startTwo(t)
		if trap.Kort == "" {
			t.Errorf("fällan %q har ingen kort beskrivning", trap.ID)
		}
		p := a.player
		p.HP = 99
		putInRoom(g, p, 4, 4)
		drainAll(a)
		g.trigger(p, g.data.trap(trap.ID), 0)
		g.broadcastState()
		for _, raw := range drainAll(a) {
			// Fällans egna texter får använda vanliga ord som "mörkret".
			own := false
			for _, text := range trap.Texter {
				own = own || strings.Contains(raw, text)
			}
			if own {
				continue
			}
			if strings.Contains(raw, trap.Namn) || strings.Contains(raw, lowerFirst(trap.Namn)) {
				t.Errorf("spelarsidan fick fällans namn %q: %s", trap.Namn, raw)
			}
		}
	}
}

// Varelserna har inga namn för spelarna: varken i rumsblocket, striden,
// när de dör eller i tillståndet.
func TestMonstersHaveNoNamesForPlayers(t *testing.T) {
	for _, mt := range testGame(t).data.Monster.Typer {
		g, a, b := startTwo(t)
		if mt.Kort == "" || mt.KortBestamd == "" {
			t.Errorf("varelsen %q saknar kort beskrivning", mt.ID)
		}
		p := a.player
		p.HP = 99
		putInRoom(g, p, 4, 2)
		putInRoom(g, b.player, 1, 1)
		m := &Monster{ID: 77, Type: mt.ID, X: 4, Y: 2, HP: 1}
		g.monsters = append(g.monsters, m)
		drainAll(a)
		g.tellRoom(p)
		for range 6 {
			if g.monsterAt(4, 2) == nil {
				break
			}
			g.Act(a, actAttack, m.ID)
			g.Act(b, actWait, 0)
		}
		for _, raw := range drainAll(a) {
			if strings.Contains(raw, mt.Namn) || strings.Contains(raw, lowerFirst(mt.Namn)) {
				t.Errorf("spelarsidan fick varelsens namn %q: %s", mt.Namn, raw)
			}
		}
	}
}

// Spåret efter en fälla som nyss slog till mot vandraren står inte i
// blocket, eftersom blocket står före loggens rad om fällan.
func TestRoomBlockSkipsTrapThatJustHit(t *testing.T) {
	g, a, _ := startTwo(t)
	pa := a.player
	pa.HP = 99
	putInRoom(g, pa, 3, 3)
	pa.told = told{}
	trap := g.data.trap("roterande")
	g.trigger(pa, trap, 0)
	rooms(a)
	g.tellRoom(pa)
	for _, m := range rooms(a) {
		if strings.Contains(m.Har, trap.SparEfter[0]) {
			t.Fatalf("spåret stod i blocket samma runda: %+v", m)
		}
	}
	// Nästa runda har inget ändrats, och spåret sägs inte i efterhand.
	g.tellRoom(pa)
	if got := rooms(a); len(got) != 0 {
		t.Fatalf("spåret kom i efterhand: %+v", got)
	}
}

// En fälld varelse är inte "borta": den som fällde den får inget, de
// andra i rummet ser den ligga där.
func TestRoomBlockSlainMonster(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	putInRoom(g, pa, 4, 4)
	putInRoom(g, pb, 4, 4)
	m := &Monster{ID: 5, Type: g.data.Monster.Typer[0].ID, X: 4, Y: 4, HP: 1}
	g.monsters = append(g.monsters, m)
	g.tellRoom(pa)
	g.tellRoom(pb)
	rooms(a)
	rooms(b)
	g.defeat(pa, m)
	g.tellRoom(pa)
	g.tellRoom(pb)
	mt := g.data.monsterType(m.Type)
	borta := g.data.Text("har_varelse_borta", "varelse", upperFirst(mt.KortBestamd))
	for _, got := range rooms(a) {
		if strings.Contains(got.Har, borta) {
			t.Errorf("den som fällde varelsen fick %q", got.Har)
		}
	}
	want := g.data.Text("har_varelse_falld", "varelse", upperFirst(mt.KortBestamd))
	got := rooms(b)
	if len(got) != 1 || got[0].Har != want {
		t.Errorf("den andra väntades se %q, fick %+v", want, got)
	}
}

// En strid som pågår hörs bara första rundan, inte runda efter runda.
func TestOngoingFightHeardOnce(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	putInRoom(g, pa, 1, 1)
	putInRoom(g, pb, 5, 5)
	g.tellRoom(pa)
	rooms(a)
	heard := func() string {
		g.hearing = true
		g.hearFight(5, 5, pb)
		g.hearing = false
		g.tellRoom(pa)
		got := rooms(a)
		if len(got) == 0 {
			return ""
		}
		return got[0].Hor
	}
	if heard() == "" {
		t.Fatal("striden hördes inte första rundan")
	}
	if h := heard(); h != "" {
		t.Fatalf("samma strid hördes igen: %q", h)
	}
	// En runda tyst, sedan börjar det om och hörs igen.
	g.tellRoom(pa)
	if heard() == "" {
		t.Fatal("en ny strid hördes inte")
	}
}

// Varelsen beskrivs helt första stridsrundan och kort rundorna efter.
func TestMonsterShortAfterFirstRound(t *testing.T) {
	g, a, _ := startTwo(t)
	for _, mt := range g.data.Monster.Typer {
		if mt.KortIgen == "" {
			t.Errorf("varelsen %q saknar kort_igen", mt.ID)
		}
	}
	m := &Monster{ID: 9, Type: g.data.Monster.Typer[0].ID}
	mt := g.data.monsterType(m.Type)
	if v := g.monsterVars(m, a.player); v[3] != mt.KortBestamd {
		t.Fatalf("första rundan väntades %q, fick %q", mt.KortBestamd, v[3])
	}
	if v := g.monsterVars(m, a.player); v[3] != mt.KortBestamd {
		t.Fatalf("samma runda väntades %q, fick %q", mt.KortBestamd, v[3])
	}
	g.round++
	if v := g.monsterVars(m, a.player); v[3] != mt.KortIgen {
		t.Fatalf("rundan efter väntades %q, fick %q", mt.KortIgen, v[3])
	}
}
