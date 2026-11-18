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
	"math/rand/v2"
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
}

// monsterMove är det drag den döda spelaren har valt för en varelse.
type monsterMove struct {
	id  int
	dir Dir
}

type monsterView struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Glyph string `json:"glyph"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	HP    int    `json:"hp"`
	MaxHP int    `json:"maxHp"`
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
	if rand.Float64() >= g.data.Monster.Chans {
		return false
	}
	w := map[string]int{}
	for _, t := range g.data.Monster.Typer {
		w[t.ID] = t.Vikt
	}
	t := g.data.monsterType(weighted(w))
	if t == nil {
		return false
	}
	g.nextMonsterID++
	g.monsters = append(g.monsters, &Monster{ID: g.nextMonsterID, Type: t.ID, X: x, Y: y, HP: t.Liv})
	return true
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
		return []string{actAttack, actParry, actDodge}[rand.IntN(3)]
	}
	return weighted(t.Vanor)
}

// monsterVars ger namn och val som texterna kan använda.
func (g *Game) monsterVars(m *Monster, extra ...string) []string {
	name := g.data.monsterType(m.Type).Namn
	return append([]string{"Monster", name, "monster", lowerFirst(name)}, extra...)
}

func (g *Game) choiceWord(a string) string {
	if w := g.data.Monster.Val[a]; w != "" {
		return w
	}
	return a
}

// announceMonster berättar för spelaren att en varelse finns i rummet.
func (g *Game) announceMonster(p *Player) {
	m := g.monsterAt(p.X, p.Y)
	if m == nil || p.Status != statusIn {
		return
	}
	t := g.data.monsterType(m.Type)
	p.Seen[m.Type] = true
	g.soundTo(p, "monster")
	g.logTo(p, t.Namn+"!", "varn")
	g.logTo(p, pick(t.DykerUpp), "")
	g.explainNow(p, m.Type)
}

// duel utför en stridsrunda mellan spelaren och varelsen.
func (g *Game) duel(p *Player, m *Monster, a string) {
	t := g.data.monsterType(m.Type)
	mc := g.monsterChoice(m)
	g.logTo(p, g.data.MonsterText("duell", g.monsterVars(m,
		"ditt_val", g.choiceWord(a), "dess_val", g.choiceWord(mc))...), "")
	g.logOthers(p, g.data.MonsterText("strid_andra", "region", g.data.region(p.X, p.Y)), "dim")

	playerHits := beats(a, mc) || (a == actAttack && mc == actAttack)
	monsterHits := beats(mc, a) || (a == actAttack && mc == actAttack)
	switch {
	case playerHits && monsterHits:
		g.soundTo(p, "hugg")
		g.logTo(p, g.data.MonsterText("bada_traffar", g.monsterVars(m, "skada", strconv.Itoa(t.Skada))...), "varn")
	case playerHits:
		g.soundTo(p, "hugg")
		g.logTo(p, g.data.MonsterText("du_traffar", g.monsterVars(m)...), "syn")
	case monsterHits:
		g.logTo(p, g.data.MonsterText("monster_traffar", g.monsterVars(m, "skada", strconv.Itoa(t.Skada))...), "varn")
	default:
		g.soundTo(p, choiceSound(a))
		g.logTo(p, g.data.MonsterText("inget", g.monsterVars(m)...), "dim")
	}
	if monsterHits {
		g.damage(p, t.Skada, false)
	}
	if playerHits {
		m.HP -= g.data.Monster.Spelarskada
		if m.HP <= 0 {
			g.defeat(p, m)
		}
	}
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
	mc := g.monsterChoice(m)
	g.logTo(p, g.data.MonsterText(key, g.monsterVars(m, "dess_val", g.choiceWord(mc))...), "varn")
	if mc == actAttack {
		g.logTo(p, g.data.MonsterText("monster_traffar", g.monsterVars(m, "skada", strconv.Itoa(t.Skada))...), "varn")
		g.damage(p, t.Skada, false)
	}
}

func (g *Game) defeat(p *Player, m *Monster) {
	t := g.data.monsterType(m.Type)
	g.removeMonster(m)
	p.Treasure += t.Guld
	g.soundTo(p, "monster_dor")
	g.logTo(p, pick(t.Dor), "")
	g.logTo(p, g.data.MonsterText("seger", g.monsterVars(m,
		"guld", strconv.Itoa(t.Guld), "skatt", strconv.Itoa(p.Treasure))...), "syn")
}

// isController är sant för den döda spelare som styr varelserna.
func (g *Game) isController(p *Player) bool {
	return p != nil && g.controller != 0 && p.ID == g.controller && p.Status == statusDead
}

// makeController låter den första som dör ta över varelserna.
func (g *Game) makeController(p *Player) {
	if g.controller != 0 {
		return
	}
	g.controller = p.ID
	g.soundTo(p, "tjanare")
	g.logTo(p, g.data.MonsterText("tjanare"), "syn")
	g.logOthers(p, g.data.MonsterText("tjanare_andra", "namn", p.Name), "varn")
}

// canMonsterGo säger om varelsen kan gå åt hållet d. Varelser går bara genom
// öppningar till kända, vanliga rum där ingen annan varelse finns.
func (g *Game) canMonsterGo(m *Monster, d Dir) bool {
	room := &g.rooms[m.Y][m.X]
	if !room.Doors[d] {
		return false
	}
	nx, ny := m.X+dirDX[d], m.Y+dirDY[d]
	if !inBounds(nx, ny) {
		return false
	}
	n := &g.rooms[ny][nx]
	return n.Revealed && n.Kind == "rum" && g.monsterAt(nx, ny) == nil
}

// actController tar emot den döda spelarens val. a är en riktning eller "vanta".
func (g *Game) actController(p *Player, a string, id int) {
	if a == actWait {
		p.Action = actWait
		g.pendingMove = nil
	} else {
		d, ok := dirCodes[a]
		m := g.monsterByID(id)
		if !ok || m == nil {
			return
		}
		if !g.canMonsterGo(m, d) {
			g.logTo(p, g.data.MonsterText("kan_inte_flytta"), "dim")
			return
		}
		p.Action = "m" + strconv.Itoa(id) + a
		g.pendingMove = &monsterMove{id: id, dir: d}
	}
	g.soundTo(p, "val")
	if g.allChosen() {
		g.resolve()
		return
	}
	g.broadcastState()
}

// moveMonsters utför den döda spelarens drag i början av rundan.
func (g *Game) moveMonsters() {
	mv := g.pendingMove
	g.pendingMove = nil
	if mv == nil {
		return
	}
	m := g.monsterByID(mv.id)
	if m == nil || !g.canMonsterGo(m, mv.dir) {
		return
	}
	m.X += dirDX[mv.dir]
	m.Y += dirDY[mv.dir]
	for _, p := range g.players {
		if p != nil && p.Status == statusIn && p.Cellar == 0 && p.X == m.X && p.Y == m.Y {
			g.soundTo(p, "monster")
			g.logTo(p, g.data.MonsterText("monster_kommer", g.monsterVars(m)...), "varn")
		}
	}
}

func (g *Game) monsterViews(me *Player, blind bool) []monsterView {
	var out []monsterView
	for _, m := range g.monsters {
		if !g.rooms[m.Y][m.X].Revealed || (blind && (m.X != me.X || m.Y != me.Y)) {
			continue
		}
		t := g.data.monsterType(m.Type)
		out = append(out, monsterView{
			ID: m.ID, Name: t.Namn, Glyph: t.Tecken, X: m.X, Y: m.Y, HP: m.HP, MaxHP: t.Liv,
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
		// Flykt: varelsen får ett fritt drag, sedan går spelaren.
		g.freeAttack(p, m, "fly")
		return p.Status != statusIn
	default:
		g.freeAttack(p, m, "inte_beredd")
		return true
	}
}
