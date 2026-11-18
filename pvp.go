package main

// Strid mellan vandrare. Samma sten-sax-påse som mot varelserna: båda väljer
// i hemlighet anfall, parering eller undanmanöver. Den som träffar rycker åt
// sig en del av motståndarens guld, och den som dödar tar allt. Portrummen är
// fredade, och varje runda med strid gör ormen lite oroligare.

import (
	"math/rand/v2"
	"strconv"
)

// Lägen för strid mellan vandrare.
const (
	pvpOff  = "av"
	pvpOn   = "pa"
	pvpHalf = "halva" // först när solen har sjunkit halvvägs
)

// Vad en vandrare får veta om strid i rummet den står i.
const (
	pvpNone     = ""       // ingen annan vandrare här
	pvpYes      = "ja"     // det går att strida
	pvpPeace    = "fredat" // portrummet är fredat
	pvpTooEarly = "for_tidigt"
	pvpDisabled = "av"
)

// pvpSkada är hur mycket liv ett slag från en annan vandrare kostar.
const pvpSkada = 1

// foesHere ger de andra vandrarna i samma rum som p.
func (g *Game) foesHere(p *Player) []*Player {
	var out []*Player
	for _, o := range g.players {
		if o != nil && o != p && o.Status == statusIn && o.Cellar == 0 && o.X == p.X && o.Y == p.Y {
			out = append(out, o)
		}
	}
	return out
}

// pvpStatus säger om p kan strida mot någon i rummet just nu.
func (g *Game) pvpStatus(p *Player) string {
	if p.Status != statusIn || p.Cellar > 0 || len(g.foesHere(p)) == 0 {
		return pvpNone
	}
	switch {
	case g.cfg.PvP == pvpOff || g.cfg.PvP == "":
		return pvpDisabled
	case g.rooms[p.Y][p.X].Kind == "ingang":
		return pvpPeace
	case g.cfg.PvP == pvpHalf && g.sun > g.cfg.SunMax/2:
		return pvpTooEarly
	}
	return pvpYes
}

// botFoe väljer vem en AI-spelare kan ge sig på. AI-spelare anfaller bara
// människor, utom när det inte finns någon människa bland vandrarna: då får
// de slåss med varandra som vilka vandrare som helst.
func (g *Game) botFoe(p *Player) *Player {
	if g.pvpStatus(p) != pvpYes {
		return nil
	}
	humans := false
	for _, o := range g.players {
		if o != nil && o.bot == nil {
			humans = true
		}
	}
	for _, o := range g.foesHere(p) {
		if o.bot == nil || !humans {
			return o
		}
	}
	return nil
}

// pickFoe väljer motståndare: den spelaren pekade ut om den står här,
// annars den första i rummet.
func (g *Game) pickFoe(p *Player, wanted int) *Player {
	foes := g.foesHere(p)
	for _, o := range foes {
		if o.ID == wanted {
			return o
		}
	}
	if len(foes) > 0 {
		return foes[0]
	}
	return nil
}

// resolveFights avgör alla strider mellan vandrare innan någon går eller tar
// skatt. Den som inte försvarar sig blir träffad om den andra anfaller, även
// om den sedan går ut ur rummet: ett slag i ryggen.
func (g *Game) resolveFights(active []*Player) {
	done := map[*Player]bool{}
	fought := false
	for _, p := range active {
		if done[p] || p.Status != statusIn || !isCombatAction(p.Action) || g.monsterAt(p.X, p.Y) != nil {
			continue
		}
		foe := g.playerByID(p.target)
		if foe == nil || foe.Status != statusIn || foe.X != p.X || foe.Y != p.Y || g.pvpStatus(p) != pvpYes {
			done[p] = true
			g.logTo(p, g.data.Text("pvp_borta"), "dim")
			continue
		}
		fought = true
		done[p] = true
		if isCombatAction(foe.Action) && foe.target == p.ID && !done[foe] {
			done[foe] = true
			g.exchange(p, foe)
		} else {
			g.ambush(p, foe)
		}
	}
	if fought {
		g.fightStirsSerpent(active)
	}
}

// exchange: båda har valt att slåss mot varandra.
func (g *Game) exchange(a, b *Player) {
	aHits := beats(a.Action, b.Action) || (a.Action == actAttack && b.Action == actAttack)
	bHits := beats(b.Action, a.Action) || (a.Action == actAttack && b.Action == actAttack)
	g.logTo(a, g.data.Text("pvp_duell", "namn", b.Name, "ditt_val", g.choiceWord(a.Action), "dess_val", g.choiceWord(b.Action)), "")
	g.logTo(b, g.data.Text("pvp_duell", "namn", a.Name, "ditt_val", g.choiceWord(b.Action), "dess_val", g.choiceWord(a.Action)), "")
	g.fightNoise(a, b)
	if !aHits && !bHits {
		g.soundTo(a, choiceSound(a.Action))
		g.soundTo(b, choiceSound(b.Action))
		g.logTo(a, g.data.Text("pvp_inget", "namn", b.Name), "dim")
		g.logTo(b, g.data.Text("pvp_inget", "namn", a.Name), "dim")
		return
	}
	// Slagen faller samtidigt: räkna bytet innan någon av dem tar något.
	aLoot, bLoot := 0, 0
	if aHits {
		aLoot = loot(b)
	}
	if bHits {
		bLoot = loot(a)
	}
	if aHits {
		g.hit(a, b, aLoot)
	}
	if bHits {
		g.hit(b, a, bLoot)
	}
}

// ambush: anfallaren slår mot någon som gör något annat.
func (g *Game) ambush(a, b *Player) {
	if a.Action != actAttack {
		g.logTo(a, g.data.Text("pvp_vantar", "namn", b.Name, "ditt_val", g.choiceWord(a.Action)), "dim")
		return
	}
	g.fightNoise(a, b)
	g.logTo(a, g.data.Text("pvp_oberedd", "namn", b.Name), "")
	g.logTo(b, g.data.Text("pvp_bakhall", "namn", a.Name), "varn")
	g.hit(a, b, loot(b))
}

// loot är hur mycket guld en träff river åt sig: en tredjedel, minst 1, och
// allt om slaget dödar.
func loot(victim *Player) int {
	if victim.Treasure <= 0 {
		return 0
	}
	if victim.HP <= pvpSkada {
		return victim.Treasure
	}
	return max(1, victim.Treasure/3)
}

// hit låter a träffa b, ta bytet och kanske döda b.
func (g *Game) hit(a, b *Player, gold int) {
	if b.Status != statusIn {
		return
	}
	gold = min(gold, b.Treasure)
	b.Treasure -= gold
	a.Treasure += gold
	g.soundTo(a, "hugg")
	killing := b.HP <= pvpSkada
	switch {
	case killing:
		g.logTo(a, g.data.Text("pvp_dodar", "namn", b.Name, "guld", strconv.Itoa(gold), "skatt", strconv.Itoa(a.Treasure)), "syn")
	case gold > 0:
		g.logTo(a, g.data.Text("pvp_byte", "namn", b.Name, "guld", strconv.Itoa(gold), "skatt", strconv.Itoa(a.Treasure)), "syn")
	default:
		g.logTo(a, g.data.Text("pvp_traff_tom", "namn", b.Name), "syn")
	}
	if killing {
		g.logTo(b, g.data.Text("pvp_dodad", "namn", a.Name), "varn")
	} else if gold > 0 {
		g.logTo(b, g.data.Text("pvp_bytt", "namn", a.Name, "guld", strconv.Itoa(gold), "skada", strconv.Itoa(pvpSkada)), "varn")
	} else if !killing {
		g.logTo(b, g.data.Text("pvp_traffad", "namn", a.Name, "skada", strconv.Itoa(pvpSkada)), "varn")
	}
	b.hurt("vandrare", a.Name)
	g.damage(b, pvpSkada, false)
}

// fightNoise låter resten av tornet höra striden.
func (g *Game) fightNoise(a, b *Player) {
	g.hear(a.X, a.Y, "ljud_strid", g.data.MonsterText("strid_andra", "region", g.data.region(a.X, a.Y)), "dim", a, b)
}

// fightStirsSerpent: tornet känner när det blöder inuti det. Högst ett nytt
// vakna-kort per runda.
func (g *Game) fightStirsSerpent(active []*Player) {
	for _, c := range g.data.Ormen.Kort {
		if c.Vaknar {
			g.deck = append(g.deck, c.ID)
			i, j := len(g.deck)-1, rand.IntN(len(g.deck))
			g.deck[i], g.deck[j] = g.deck[j], g.deck[i]
			break
		}
	}
	for _, p := range active {
		if p.Status == statusIn && p.Cellar == 0 && (isCombatAction(p.Action) || len(g.foesHere(p)) > 0) {
			g.logTo(p, g.data.Text("pvp_tornet_kanner"), "info")
		}
	}
}

// packHint berättar ungefär hur mycket någon bär, utan att säga hur mycket.
func (g *Game) packHint(o *Player) string {
	switch {
	case o.Treasure >= 60:
		return g.data.Text("pack_tung", "namn", o.Name)
	case o.Treasure > 0:
		return g.data.Text("pack_lite", "namn", o.Name)
	}
	return g.data.Text("pack_tom", "namn", o.Name)
}

// meetIfNew berättar för p vilka andra vandrare som står i rummet och om det
// går att strida, men bara när sällskapet har ändrats sedan förra rundan.
func (g *Game) meetIfNew(p *Player) {
	if p.Status != statusIn || p.Cellar > 0 || p.Dark > 0 {
		p.met = ""
		return
	}
	foes := g.foesHere(p)
	key := ""
	for _, o := range foes {
		key += strconv.Itoa(o.ID) + ","
	}
	if key == p.met {
		return
	}
	p.met = key
	if len(foes) == 0 {
		return
	}
	for _, o := range foes {
		g.logTo(p, g.packHint(o), "info")
	}
	switch g.pvpStatus(p) {
	case pvpPeace:
		g.logTo(p, g.data.Text("pvp_fredat"), "dim")
	case pvpTooEarly:
		g.logTo(p, g.data.Text("pvp_for_tidigt"), "dim")
	}
}

func (g *Game) playerByID(id int) *Player {
	for _, p := range g.players {
		if p != nil && p.ID == id {
			return p
		}
	}
	return nil
}
