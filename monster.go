package main

// Monster och strider. En strid är som sten-sax-påse: båda sidor väljer i
// hemlighet anfall, parering eller undanmanöver.
//
//	parering slår anfall (slaget avleds och besvaras)
//	anfall slår undanmanöver (den som viker undan blir träffad)
//	undanmanöver slår parering (man glider runt skölden)
//
// Väljer båda anfall träffar båda. Annars händer ingenting vid lika val.

import (
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	actAttack = "anfall"
	actParry  = "parera"
	actDodge  = "undan"
)

type Monster struct {
	ID   int
	Type string
	X, Y int
	HP   int
	// named är rundan då varje vandrare (på ID) först fick höra varelsen
	// beskriven i striden. Rundorna efter räcker den korta formen.
	named map[int]int
}

type monsterView struct {
	ID    int        `json:"id"`
	Name  string     `json:"name"`
	Glyph string     `json:"glyph"`
	Type  string     `json:"typ"`              // varelsens id i monster.json, för modellen på kartan snett ovanifrån
	Model *ModelLook `json:"modell,omitempty"` // färg och storlek på modellen
	X     int        `json:"x"`
	Y     int        `json:"y"`
	HP    int        `json:"hp"`
	MaxHP int        `json:"maxHp"`
}

func isCombatAction(a string) bool {
	return a == actAttack || a == actParry || a == actDodge
}

// beats säger om val a vinner över val b.
func beats(a, b string) bool {
	return (a == actParry && b == actAttack) ||
		(a == actAttack && b == actDodge) ||
		(a == actDodge && b == actParry)
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

// monsterChoice låter varelsen välja efter sina vanor.
func (g *Game) monsterChoice(m *Monster) string {
	t := g.data.monsterType(m.Type)
	if t == nil || len(t.Vanor) == 0 {
		return []string{actAttack, actParry, actDodge}[g.rng.IntN(3)]
	}
	return weighted(g.rng, t.Vanor)
}

// monsterChoiceVs är varelsens val mot en viss vandrare. Ormtjusaren lockar
// fram fler anfall.
func (g *Game) monsterChoiceVs(m *Monster, p *Player) string {
	extra := int(g.classVal(p, "extra_anfall", 0))
	if p.initiated {
		extra += g.data.Avslojande.Invigd.ExtraAnfall
	}
	if p.has("mynt_hal") != nil {
		g.curse(p, "mynt_hal")
		extra += 2
	}
	t := g.data.monsterType(m.Type)
	if extra <= 0 || t == nil || len(t.Vanor) == 0 {
		return g.monsterChoice(m)
	}
	w := map[string]int{}
	for k, v := range t.Vanor {
		w[k] = v
	}
	w[actAttack] += extra
	return weighted(g.rng, w)
}

// monsterViews: lyktbäraren ser också varelser i grannrummen.
func (g *Game) seesMonsterAt(me *Player, x, y int) bool {
	if g.knowsRoom(me, x, y) {
		return true
	}
	if g.is(me, "lyktbarare") && g.rooms[y][x].Revealed {
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

func (g *Game) choiceWord(a string) string {
	if w := g.data.Monster.Val[a]; w != "" {
		return w
	}
	return a
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

// duel utför en stridsrunda mellan spelaren och varelsen.
func (g *Game) duel(p *Player, m *Monster, a string) {
	t := g.data.monsterType(m.Type)
	mc := g.monsterChoiceVs(m, p)
	if p.bot != nil {
		p.bot.saw(m.Type, mc)
	}
	g.logTo(p, g.data.MonsterText("duell", g.monsterVars(m, p,
		"ditt_val", g.choiceWord(a), "dess_val", g.choiceWord(mc))...), "")
	g.hearFight(p.X, p.Y, p)

	playerHits := beats(a, mc) || (a == actAttack && mc == actAttack)
	monsterHits := beats(mc, a) || (a == actAttack && mc == actAttack)
	if a == actAttack && mc == actAttack && g.is(p, "legoknekt") {
		monsterHits = false
		g.logTo(p, g.data.Text("klass_legoknekt_skold"), "info")
	}
	if a == actDodge && g.is(p, "kartritare") {
		playerHits, monsterHits = false, true
		g.logTo(p, g.data.Text("klass_kartritare_undan"), "varn")
	}
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
		// Texten följer det du valde, så att den inte säger att du vek undan
		// när du anföll.
		key := "monster_traffar_" + a
		if len(g.data.Monster.Texter[key]) == 0 {
			key = "monster_traffar"
		}
		result = g.data.MonsterText(key, g.monsterVars(m, p, "skada", strconv.Itoa(t.Skada))...)
		g.logTo(p, result, "varn")
	default:
		g.soundTo(p, choiceSound(a))
		result = g.data.MonsterText("inget", g.monsterVars(m, p)...)
		g.logTo(p, result, "dim")
	}
	if monsterHits {
		p.hurt("varelse", t.Kort)
		g.damage(p, g.shieldDamage(p, t.Skada), false)
	}
	if playerHits {
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

func choiceSound(a string) string {
	if a == actParry {
		return "parera"
	}
	return "undan"
}

// freeAttack: spelaren gör något annat än att slåss, och varelsen får ett
// fritt drag. Bara ett anfall träffar.
func (g *Game) freeAttack(p *Player, m *Monster, key string) {
	t := g.data.monsterType(m.Type)
	mc := g.monsterChoiceVs(m, p)
	g.logTo(p, g.data.MonsterText(key, g.monsterVars(m, p, "dess_val", g.choiceWord(mc))...), "varn")
	if mc == actAttack {
		text := g.data.MonsterText("monster_traffar", g.monsterVars(m, p, "skada", strconv.Itoa(t.Skada))...)
		g.logTo(p, text, "varn")
		g.monsterCard(p, t, g.fightName(m, p), text, false)
		p.hurt("varelse", t.Kort)
		g.damage(p, g.shieldDamage(p, t.Skada), false)
	}
}

func (g *Game) defeat(p *Player, m *Monster) {
	t := g.data.monsterType(m.Type)
	g.removeMonster(m)
	if g.slain != nil {
		g.slain[m.ID] = p.ID
	}
	p.Treasure += t.Guld
	p.stats.kills++
	p.roundKill = t.KortBestamd
	g.chronicle("kronika_varelse", "namn", p.Name, "varelse", t.Kort)
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
	return n.Kind == "rum"
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
	var out []monsterView
	for _, m := range g.monsters {
		if !g.seesMonsterAt(me, m.X, m.Y) || (blind && (m.X != me.X || m.Y != me.Y)) {
			continue
		}
		t := g.data.monsterType(m.Type)
		out = append(out, monsterView{
			ID: m.ID, Name: t.KortBestamd, Glyph: t.Tecken, Type: m.Type, Model: t.Modell, X: m.X, Y: m.Y, HP: m.HP, MaxHP: t.Liv,
		})
	}
	return out
}

// combatTurn hanterar spelarens drag när en varelse står i samma rum.
// Den returnerar false om draget ska utföras som vanligt.
func (g *Game) combatTurn(p *Player, a string) bool {
	m := g.monsterAt(p.X, p.Y)
	if m == nil {
		return false
	}
	switch {
	case isCombatAction(a):
		g.duel(p, m, a)
		return true
	case a == actDrink:
		g.drink(p)
		g.freeAttack(p, m, "inte_beredd")
		return true
	case strings.Contains("NESW", a) && len(a) == 1:
		// Flykt: varelsen får ett fritt drag, sedan går spelaren. Med
		// rökpulver hinner den inte slå.
		if g.smoke(p, dirCodes[a]) {
			return false
		}
		g.freeAttack(p, m, "fly")
		return p.Status != statusIn
	default:
		g.freeAttack(p, m, "inte_beredd")
		return true
	}
}
