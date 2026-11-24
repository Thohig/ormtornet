package main

// Strid mellan vandrare. Samma sten-sax-påse som mot varelserna: båda väljer
// i hemlighet anfall, parering eller undanmanöver. Den som träffar rycker åt
// sig en del av motståndarens guld, och den som dödar tar allt. Portrummen är
// fredade, och varje runda med strid gör ormen lite oroligare.

import (
	"strconv"
)

// Lägen för strid mellan vandrare.
const (
	pvpOff = "av" // bara i testerna
	pvpOn  = "pa"
)

// Vad en vandrare får veta om strid i rummet den står i.
const (
	pvpNone     = ""       // ingen annan vandrare här
	pvpYes      = "ja"     // det går att strida
	pvpPeace    = "fredat" // portrummet är fredat
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

// pickFoe väljer motståndare. Spelarna väljer inte själva: tornet tar den
// som anföll vandraren förra rundan, annars den som kom först in i tornet.
// Skatten avgör inte, eftersom den är hemlig. Datorspelarna kan peka ut
// någon (wanted), som de alltid har gjort.
func (g *Game) pickFoe(p *Player, wanted int) *Player {
	foes := g.foesHere(p)
	if p.bot != nil {
		for _, o := range foes {
			if o.ID == wanted {
				return o
			}
		}
	}
	if p.attackedRound >= g.round-1 {
		for _, o := range foes {
			if o.ID == p.attackedBy {
				return o
			}
		}
	}
	var best *Player
	for _, o := range foes {
		if best == nil || o.ID < best.ID {
			best = o
		}
	}
	return best
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
		foe.attackedBy, foe.attackedRound = p.ID, g.round
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

// exchange: båda anfaller varandra, och var och en träffar med sin chans.
// Pilgrimen lyfter aldrig vapen mot en vandrare.
func (g *Game) exchange(a, b *Player) {
	aHits := !g.is(a, "pilgrim") && g.roll(g.hitChance(g.kraft(a), g.forsvar(b), 0))
	bHits := !g.is(b, "pilgrim") && g.roll(g.hitChance(g.kraft(b), g.forsvar(a), 0))
	g.logTo(a, g.data.Text("pvp_duell", "namn", b.Name), "")
	g.logTo(b, g.data.Text("pvp_duell", "namn", a.Name), "")
	g.fightNoise(a, b)
	if !aHits && !bHits {
		g.soundTo(a, "parera")
		g.soundTo(b, "parera")
		g.logTo(a, g.data.Text("pvp_inget", "namn", b.Name), "dim")
		g.logTo(b, g.data.Text("pvp_inget", "namn", a.Name), "dim")
		return
	}
	// Slagen faller samtidigt: räkna bytet innan någon av dem tar något.
	aLoot, bLoot := 0, 0
	if aHits {
		aLoot = g.loot(a, b)
	}
	if bHits {
		bLoot = g.loot(b, a)
	}
	if aHits {
		g.hit(a, b, aLoot, false)
	}
	if bHits {
		g.hit(b, a, bLoot, false)
	}
}

// ambush: anfallaren slår mot någon som gör något annat. Den som håller
// garden är svårare att träffa, den som söker eller dricker lättare, och
// kartritaren som flyr blir alltid träffad.
func (g *Game) ambush(a, b *Player) {
	if g.is(a, "pilgrim") {
		g.logTo(a, g.data.Text("klass_pilgrim_fred", "namn", b.Name), "dim")
		return
	}
	g.fightNoise(a, b)
	_, fleeing := dirCodes[b.Action]
	guard := b.Action == actWait
	always := fleeing && g.slowFeet(b)
	if !always && !g.roll(g.hitChance(g.kraft(a), g.forsvar(b), g.guardMod(b.Action))) {
		g.soundTo(a, "parera")
		if guard {
			g.logTo(a, g.data.Text("pvp_garden_miss", "namn", b.Name), "dim")
			g.logTo(b, g.data.Text("pvp_garden_hall", "namn", a.Name), "dim")
		} else {
			g.logTo(a, g.data.Text("pvp_miss", "namn", b.Name), "dim")
			g.logTo(b, g.data.Text("pvp_missad", "namn", a.Name), "varn")
		}
		return
	}
	if guard {
		g.logTo(a, g.data.Text("pvp_garden_traff", "namn", b.Name), "")
		g.logTo(b, g.data.Text("pvp_garden_traffad", "namn", a.Name), "varn")
	} else {
		g.logTo(a, g.data.Text("pvp_oberedd", "namn", b.Name), "")
		g.logTo(b, g.data.Text("pvp_bakhall", "namn", a.Name), "varn")
	}
	g.hit(a, b, g.loot(a, b), true)
}

// loot är hur mycket guld en träff river åt sig: en tredjedel (hälften för
// tjuven), minst 1, och allt om slaget dödar.
func (g *Game) loot(attacker, victim *Player) int {
	if victim.Treasure <= 0 {
		return 0
	}
	if victim.HP <= pvpSkada {
		return victim.Treasure
	}
	share := g.classVal(attacker, "byte", 1.0/3)
	return max(1, int(float64(victim.Treasure)*share))
}

// hit låter a träffa b, ta bytet och kanske döda b. bakhall säger att b
// inte försvarade sig, för krönikan.
func (g *Game) hit(a, b *Player, gold int, bakhall bool) {
	if b.Status != statusIn {
		return
	}
	gold = g.hiddenShare(b, min(gold, b.Treasure))
	b.Treasure -= gold
	a.Treasure += gold
	a.stats.hits++
	a.stats.stolen += gold
	g.soundTo(a, "hugg")
	killing := b.HP <= pvpSkada
	// Siffrorna, liv och guld, står på spelkortet.
	var aText, bText string
	switch {
	case killing:
		aText = g.data.Text("pvp_dodar", "namn", b.Name, "guld", strconv.Itoa(gold))
	case gold > 0:
		aText = g.data.Text("pvp_byte", "namn", b.Name, "guld", strconv.Itoa(gold))
	default:
		aText = g.data.Text("pvp_traff_tom", "namn", b.Name)
	}
	g.logTo(a, aText, "syn")
	g.addCard(a, gameCard{Slag: "vandrare", Bild: "vandrare", Rubrik: b.Name, Text: aText})
	if killing {
		g.takeNameParts(a, b)
		bText = g.data.Text("pvp_dodad", "namn", a.Name)
	} else if gold > 0 {
		bText = g.data.Text("pvp_bytt", "namn", a.Name, "guld", strconv.Itoa(gold))
	} else {
		bText = g.data.Text("pvp_traffad", "namn", a.Name)
	}
	g.logTo(b, bText, "varn")
	g.addCard(b, gameCard{Slag: "vandrare", Bild: "vandrare", Rubrik: a.Name, Text: bText})
	if !killing {
		g.chronicleHit(a, b, gold, bakhall)
	}
	b.hurt("vandrare", a.Name)
	g.damage(b, pvpSkada, false)
}

// chronicleHit skriver ner ett slag mellan vandrare som inte dödade. Ett
// slag som dödar står redan i krönikan som en död.
func (g *Game) chronicleHit(a, b *Player, gold int, bakhall bool) {
	key := "kronika_slag"
	if bakhall {
		key = "kronika_bakhall"
	}
	if gold > 0 {
		key += "_byte"
	}
	g.chronicle(key, "namn", a.Name, "offer", b.Name, "guld", strconv.Itoa(gold))
}

// fightNoise låter resten av tornet höra striden.
func (g *Game) fightNoise(a, b *Player) {
	g.hearFight(a.X, a.Y, a, b)
}

// fightStirsSerpent: tornet känner när det blöder inuti det. Högst ett nytt
// vakna-kort per runda.
func (g *Game) fightStirsSerpent(active []*Player) {
	g.stirSerpent()
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
	g.shareMap(p, foes)
	for _, o := range foes {
		g.logTo(p, g.packHint(o), "info")
	}
	if g.pvpStatus(p) == pvpPeace {
		g.logTo(p, g.data.Text("pvp_fredat"), "dim")
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
