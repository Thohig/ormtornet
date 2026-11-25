package main

// Monster och strider (docs/bara-anfall.md). Det enda stridsvalet är
// anfall: då slår både vandraren och varelsen, och var och en träffar med
// sin chans (hitChance). Gör vandraren något annat anfaller varelsen med
// chansen i sitt lynne. Den som väntar håller garden och är svårare att
// träffa; den som söker eller dricker är lättare.

import (
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const actAttack = "anfall"

type Monster struct {
	ID   int
	Type string
	X, Y int
	HP   int
	// did är varelsens val under rundan och hit att den blev träffad, för
	// figuren på kartan.
	did string
	hit bool
	// named är rundan då varje vandrare (på ID) först fick höra varelsen
	// beskriven i striden. Rundorna efter räcker den korta formen.
	named map[int]int
}

type monsterView struct {
	ID    int        `json:"id"`
	Name  string     `json:"name"`
	Glyph string     `json:"glyph"`
	Type  string     `json:"typ"`               // varelsens id i monster.json, för modellen på kartan snett ovanifrån
	Model *ModelLook `json:"modell,omitempty"`  // färg och storlek på modellen
	Did   string     `json:"gjorde,omitempty"`  // varelsens val under rundan, för figuren
	Hit   bool       `json:"traffad,omitempty"` // varelsen blev träffad under rundan
	X     int        `json:"x"`
	Y     int        `json:"y"`
	HP    int        `json:"hp"`
	MaxHP int        `json:"maxHp"`
}

func isCombatAction(a string) bool {
	return a == actAttack
}

// hitChance är chansen i procent att den med kraften träffar den med
// försvaret, med mod (garden eller oberedd) lagt till.
func (g *Game) hitChance(kraft, forsvar, mod int) int {
	r := g.data.Monster.Traff
	c := r.Bas + r.PerSteg*(kraft-forsvar) + mod
	return min(r.Max, max(r.Min, c))
}

// roll är sant med chansen pct i procent.
func (g *Game) roll(pct int) bool {
	return g.rng.IntN(100) < pct
}

// kraft och forsvar är vandrarens värden i strid: klassens, eller
// värdena för den som inte har någon klass.
// Föremålen i väskan lägger till sitt.
func (g *Game) kraft(p *Player) int {
	n := g.data.Monster.Traff.Kraft
	if c := g.classOf(p); c != nil {
		n = c.Kraft
	}
	k, _ := g.itemBonus(p)
	return max(0, n+k)
}

func (g *Game) forsvar(p *Player) int {
	n := g.data.Monster.Traff.Forsvar
	if c := g.classOf(p); c != nil {
		n = c.Forsvar
	}
	_, f := g.itemBonus(p)
	return max(0, n+f)
}

// fightView säger i ord hur lätt vandraren och motståndaren träffar
// varandra när de anfaller: "latt", "ibland" eller "sallan".
type fightView struct {
	Du  string `json:"du,omitempty"`
	Den string `json:"den"`
}

// fightOutlook är hur striden i vandrarens rum ser ut, mot varelsen eller
// mot den vandrare tornet väljer. Talen syns aldrig, bara orden.
func (g *Game) fightOutlook(p *Player) *fightView {
	if g.phase != phasePlay || p.Status != statusIn || p.Cellar > 0 {
		return nil
	}
	if m := g.monsterAt(p.X, p.Y); m != nil {
		t := g.data.monsterType(m.Type)
		if t == nil {
			return nil
		}
		return &fightView{Du: hitWord(g.hitChance(g.kraft(p), t.Forsvar, 0)), Den: hitWord(g.hitChance(t.Kraft, g.forsvar(p), 0))}
	}
	if g.pvpStatus(p) != pvpYes {
		return nil
	}
	foe := g.pickFoe(p, 0)
	if foe == nil {
		return nil
	}
	du := hitWord(g.hitChance(g.kraft(p), g.forsvar(foe), 0))
	if g.is(p, "pilgrim") {
		du = "" // pilgrimen lyfter aldrig vapen mot en vandrare
	}
	return &fightView{Du: du, Den: hitWord(g.hitChance(g.kraft(foe), g.forsvar(p), 0))}
}

// hitWord gör en träffchans till ett ord.
func hitWord(pct int) string {
	switch {
	case pct >= 65:
		return "latt"
	case pct >= 35:
		return "ibland"
	}
	return "sallan"
}

// guardMod är hur mycket lättare eller svårare det är att träffa den som
// gjorde a i stället för att anfalla: garden när hen väntar, oberedd när
// hen söker, dricker eller tar skatt. Den som går eller inte valde något
// har inget tillägg.
func (g *Game) guardMod(a string) int {
	switch a {
	case actWait:
		return g.data.Monster.Traff.Garden
	case actSearch, actDrink, actTake:
		return g.data.Monster.Traff.Oberedd
	}
	return 0
}

func lowerFirst(s string) string {
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

func upperFirst(s string) string {
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func (g *Game) monsterAt(x, y int) *Monster {
	for _, m := range g.monsters {
		if m.X == x && m.Y == y {
			return m
		}
	}
	return nil
}

func (g *Game) monsterByID(id int) *Monster {
	for _, m := range g.monsters {
		if m.ID == id {
			return m
		}
	}
	return nil
}

// spawnMonster kanske placerar en varelse i ett nytt rum.
func (g *Game) spawnMonster(x, y int) bool {
	if g.rng.Float64() >= g.data.Monster.Chans {
		return false
	}
	w := map[string]int{}
	for _, t := range g.data.Monster.Typer {
		w[t.ID] = t.Vikt
	}
	t := g.data.monsterType(weighted(g.rng, w))
	if t == nil {
		return false
	}
	g.addMonster(t, x, y)
	return true
}

// addMonster sätter en varelse av typen t i rummet (x, y).
func (g *Game) addMonster(t *MonsterType, x, y int) {
	g.nextMonsterID++
	g.monsters = append(g.monsters, &Monster{ID: g.nextMonsterID, Type: t.ID, X: x, Y: y, HP: t.Liv})
}

func (g *Game) removeMonster(m *Monster) {
	for i, o := range g.monsters {
		if o == m {
			g.monsters = append(g.monsters[:i], g.monsters[i+1:]...)
			return
		}
	}
}

// lynne är chansen i procent att varelsen anfaller p en runda då p inte
// anfaller. Den invigde och den som bär Myntet med hål gör den mer lynnig,
// och fjällmanteln lugnar den.
func (g *Game) lynne(m *Monster, p *Player) int {
	t := g.data.monsterType(m.Type)
	if t == nil {
		return 50
	}
	l := t.Lynne
	if p.initiated {
		l += g.data.Avslojande.Invigd.ExtraLynne
	}
	if p.has("mynt_hal") != nil {
		g.curse(p, "mynt_hal")
		l += 25
	}
	l += g.knifeScent(p)
	l += int(g.classVal(p, "extra_lynne", 0)) // gravplundraren luktar död
	if p.has("fjallmantel") != nil {
		// Fjällmanteln: varelserna tar bäraren för något av sitt eget slag.
		l /= 2
	}
	return min(100, l)
}

// monsterViews: lyktbäraren ser också varelser i grannrummen.
func (g *Game) seesMonsterAt(me *Player, x, y int) bool {
	if g.knowsRoom(me, x, y) {
		return true
	}
	if (g.is(me, "lyktbarare") || me.has("glodmask") != nil) && g.rooms[y][x].Revealed {
		_, ok := g.nextTo(me, x, y)
		return ok
	}
	return false
}

// monsterVars ger varelsen och val som texterna kan använda. Varelsen har
// inget namn för spelarna: {monster} är "det bleka mångbenta". Första
// rundan p slåss mot den står hela beskrivningen, sedan den korta formen
// (kort_igen), så att samma långa beskrivning inte upprepas runda efter
// runda.
func (g *Game) monsterVars(m *Monster, p *Player, extra ...string) []string {
	t := g.data.monsterType(m.Type)
	the := t.KortBestamd
	if m.named == nil {
		m.named = map[int]int{}
	}
	if r, ok := m.named[p.ID]; !ok {
		m.named[p.ID] = g.round
	} else if r < g.round && t.KortIgen != "" {
		the = t.KortIgen
	}
	return append([]string{"Monster", upperFirst(the), "monster", the}, extra...)
}

// announceMonster säger med ett ljud att en varelse finns i rummet. Vad
// spelaren ser står i rumsblocket.
func (g *Game) announceMonster(p *Player) {
	m := g.monsterAt(p.X, p.Y)
	if m == nil || p.Status != statusIn {
		return
	}
	p.Seen[m.Type] = true
	g.soundTo(p, "monster")
	// Spelkortet: det man ser när varelsen dyker upp, samma text som i
	// rumsblocket. I mörkret ser man den inte.
	if t := g.data.monsterType(m.Type); t != nil && p.Dark == 0 {
		g.learn(p, "varelse", m.Type)
		text := ""
		if len(t.DykerUpp) > 0 {
			text = t.DykerUpp[m.ID%len(t.DykerUpp)]
		}
		g.monsterCard(p, t, t.Kort, text, false)
	}
	// Varelsen beskrivs i rumsblocket i slutet av rundan (here.go), med sin
	// egen text och utan utrop.
	g.explainNow(p, m.Type)
	g.knife(p, m)
}

// duel utför en stridsrunda mellan spelaren och varelsen: spelaren
// anfaller, varelsen slår tillbaka, och var och en träffar med sin chans.
func (g *Game) duel(p *Player, m *Monster) {
	t := g.data.monsterType(m.Type)
	m.did = actAttack
	g.hearFight(p.X, p.Y, p)

	playerHits := g.roll(g.hitChance(g.kraft(p), t.Forsvar, 0))
	monsterHits := g.roll(g.hitChance(t.Kraft, g.forsvar(p), 0))
	var result string // hur rundan gick, för spelkortet
	switch {
	case playerHits && monsterHits:
		g.soundTo(p, "hugg")
		result = g.data.MonsterText("bada_traffar", g.monsterVars(m, p, "skada", strconv.Itoa(t.Skada))...)
		g.logTo(p, result, "varn")
	case playerHits:
		g.soundTo(p, "hugg")
		result = g.data.MonsterText("du_traffar", g.monsterVars(m, p)...)
		g.logTo(p, result, "syn")
	case monsterHits:
		result = g.data.MonsterText("monster_traffar", g.monsterVars(m, p, "skada", strconv.Itoa(t.Skada))...)
		g.logTo(p, result, "varn")
	default:
		g.soundTo(p, "parera")
		result = g.data.MonsterText("inget", g.monsterVars(m, p)...)
		g.logTo(p, result, "dim")
	}
	if monsterHits {
		p.hurt("varelse", t.Kort)
		g.shieldWear(p)
		g.damage(p, t.Skada, false)
	}
	if playerHits {
		g.learnPart(p, srcCreature) // under fjällen där vapnet träffade
		m.hit = true
		m.HP -= g.data.Monster.Spelarskada
		if m.HP <= 0 {
			g.defeat(p, m) // kortet säger att varelsen föll
			return
		}
	}
	g.monsterCard(p, t, g.fightName(m, p), result, false)
}

// fightName är varelsens korta beskrivning i striden, som i texterna:
// "det bleka mångbenta", och från andra rundan "det mångbenta".
func (g *Game) fightName(m *Monster, p *Player) string {
	vars := g.monsterVars(m, p)
	return vars[3]
}

// monsterCard lägger varelsens spelkort, med den korta beskrivningen som
// rubrik. fallen är en varelse som vandraren fällde.
func (g *Game) monsterCard(p *Player, t *MonsterType, rubrik, text string, fallen bool) {
	g.addCard(p, gameCard{Slag: "varelse", Bild: cardImage(t.Bild, t.ID), Rubrik: rubrik, Text: text, Fallen: fallen})
}

// freeAttack: spelaren gör något annat än att anfalla, och varelsen
// anfaller kanske, med chansen i sitt lynne. key är texternas början
// (inte_beredd, garden eller fly) och mod hur mycket lättare eller svårare
// spelaren är att träffa. always är kartritaren som flyr: slaget mot
// ryggen tar alltid.
func (g *Game) freeAttack(p *Player, m *Monster, key string, mod int, always bool) {
	t := g.data.monsterType(m.Type)
	if !g.roll(g.lynne(m, p)) {
		g.logTo(p, g.data.MonsterText(key+"_lugn", g.monsterVars(m, p)...), "dim")
		return
	}
	m.did = actAttack
	if !always && !g.roll(g.hitChance(t.Kraft, g.forsvar(p), mod)) {
		g.soundTo(p, "parera")
		text := g.data.MonsterText(key+"_miss", g.monsterVars(m, p)...)
		g.logTo(p, text, "dim")
		g.monsterCard(p, t, g.fightName(m, p), text, false)
		return
	}
	if always {
		key := "klass_kartritare_fly"
		if !g.is(p, "kartritare") {
			key = "brynja_fly"
		}
		g.logTo(p, g.data.Text(key), "varn")
	}
	text := g.data.MonsterText(key+"_traff", g.monsterVars(m, p, "skada", strconv.Itoa(t.Skada))...)
	g.logTo(p, text, "varn")
	g.monsterCard(p, t, g.fightName(m, p), text, false)
	p.hurt("varelse", t.Kort)
	g.shieldWear(p)
	g.damage(p, t.Skada, false)
}

func (g *Game) defeat(p *Player, m *Monster) {
	t := g.data.monsterType(m.Type)
	g.removeMonster(m)
	g.fallen = append(g.fallen, m)
	if g.slain != nil {
		g.slain[m.ID] = p.ID
	}
	if p.Status != statusIn {
		// Båda träffade och vandraren föll: varelsen faller också, och dess
		// guld hamnar bland den dödes guld i rummet.
		g.addLeftGold(p, t.Guld)
		p.stats.kills++
		return
	}
	p.Treasure += t.Guld
	g.goldNote(p, "varelse", t.Guld, "varelse", t.Kort)
	p.stats.kills++
	p.roundKill = t.KortBestamd
	g.knifeRemembers(p)
	g.chronicle("kronika_varelse", "namn", p.Name, "varelse", t.Kort, "bild", cardImage(t.Bild, t.ID))
	g.soundTo(p, "monster_dor")
	dor := pick(t.Dor)
	g.logTo(p, dor, "")
	g.monsterCard(p, t, t.KortBestamd, dor, true)
	g.logTo(p, g.data.MonsterText("seger", g.monsterVars(m, p,
		"guld", strconv.Itoa(t.Guld), "skatt", strconv.Itoa(p.Treasure))...), "syn")
}

// canMonsterGo säger om varelsen kan gå åt hållet d. Varelser går genom
// öppningar till vanliga rum där ingen annan varelse finns, också till rum
// som ingen har hittat än. Portrummen och kammaren går de aldrig in i.
func (g *Game) canMonsterGo(m *Monster, d Dir) bool {
	room := &g.rooms[m.Y][m.X]
	if !room.Doors[d] {
		return false
	}
	nx, ny := m.X+dirDX[d], m.Y+dirDY[d]
	if !g.inBounds(nx, ny) || g.monsterAt(nx, ny) != nil {
		return false
	}
	n := &g.rooms[ny][nx]
	if !n.Revealed {
		_, fixed := g.fixedDoors(nx, ny)
		return !fixed
	}
	return n.Kind == "rum" && n.blocked == 0
}

// stepMonster flyttar varelsen ett rum åt hållet d. Finns rummet inte än
// byggs det, men det räknas inte som hittat förrän en vandrare går in.
func (g *Game) stepMonster(m *Monster, d Dir) {
	m.X += dirDX[d]
	m.Y += dirDY[d]
	g.build(m.X, m.Y, false)
}

// playerIn är sant när en vandrare står i rummet (x, y).
func (g *Game) playerIn(x, y int) bool {
	for _, p := range g.players {
		if p != nil && p.Status == statusIn && p.X == x && p.Y == y {
			return true
		}
	}
	return false
}

// wanderMonsters låter varelserna vandra av sig själva i början av rundan.
// En varelse som står hos en vandrare stannar och strider, och ingen
// vandrar in till en vandrare.
func (g *Game) wanderMonsters() {
	order := slices.Clone(g.monsters)
	g.rng.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })
	for _, m := range order {
		if g.playerIn(m.X, m.Y) {
			continue
		}
		chance := g.data.Monster.VandraChans
		if t := g.data.monsterType(m.Type); t != nil && t.Vandrar != nil {
			chance = *t.Vandrar
		}
		if g.rng.Float64() >= chance {
			continue
		}
		var dirs []Dir
		for d := North; d <= West; d++ {
			if g.canMonsterGo(m, d) && !g.playerIn(m.X+dirDX[d], m.Y+dirDY[d]) {
				dirs = append(dirs, d)
			}
		}
		if len(dirs) == 0 {
			continue
		}
		g.stepMonster(m, dirs[g.rng.IntN(len(dirs))])
		g.hear(m.X, m.Y, "ljud_varelse", "", "dim")
	}
}

func (g *Game) monsterViews(me *Player, blind bool) []monsterView {
	return g.viewsOf(g.monsters, me, blind)
}

// fallenViews är varelserna som föll under rundan och som me ser, så att
// kartan kan låta dem falla. Rundan efter är de borta.
func (g *Game) fallenViews(me *Player, blind bool) []monsterView {
	return g.viewsOf(g.fallen, me, blind)
}

func (g *Game) viewsOf(ms []*Monster, me *Player, blind bool) []monsterView {
	var out []monsterView
	for _, m := range ms {
		if !g.seesMonsterAt(me, m.X, m.Y) || (blind && (m.X != me.X || m.Y != me.Y)) {
			continue
		}
		t := g.data.monsterType(m.Type)
		out = append(out, monsterView{
			ID: m.ID, Name: t.KortBestamd, Glyph: t.Tecken, Type: m.Type, Model: t.Modell, X: m.X, Y: m.Y, HP: m.HP, MaxHP: t.Liv,
			Did: m.did, Hit: m.hit,
		})
	}
	return out
}

// clearMoves glömmer det figurerna gjorde förra rundan, innan en ny runda
// avgörs.
func (g *Game) clearMoves() {
	for _, p := range g.players {
		if p != nil {
			p.did, p.hit = "", false
		}
	}
	for _, m := range g.monsters {
		m.did, m.hit = "", false
	}
	g.fallen = nil
}

// combatTurn hanterar spelarens drag när en varelse står i samma rum.
// Den returnerar false om draget ska utföras som vanligt.
func (g *Game) combatTurn(p *Player, a string) bool {
	m := g.monsterAt(p.X, p.Y)
	if m == nil {
		return false
	}
	if !isCombatAction(a) {
		g.thirst(p)
		if p.Status != statusIn {
			return true
		}
	}
	switch {
	case isCombatAction(a):
		g.duel(p, m)
		return true
	case a == actDrink:
		g.drink(p)
		g.freeAttack(p, m, "inte_beredd", g.guardMod(a), false)
		return true
	case a == actSearch:
		// Man kan söka medan varelsen står där, men den kanske anfaller.
		g.freeAttack(p, m, "inte_beredd", g.guardMod(a), false)
		if p.Status == statusIn {
			g.search(p)
		}
		return true
	case a == actWait && p.Action != "":
		// Den som väntar håller garden. Den som inte valde något står
		// bara still.
		g.freeAttack(p, m, "garden", g.guardMod(a), false)
		return true
	case strings.Contains("NESW", a) && len(a) == 1:
		// Flykt: varelsen får ett fritt drag, sedan går spelaren. Med
		// rökpulver hinner den inte slå.
		if g.smoke(p, dirCodes[a]) {
			return false
		}
		g.freeAttack(p, m, "fly", 0, g.slowFeet(p))
		return p.Status != statusIn
	default:
		g.freeAttack(p, m, "inte_beredd", g.guardMod(p.Action), false)
		return true
	}
}
