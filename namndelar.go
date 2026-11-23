package main

// Ormens namn i tornet. Namnets tre delar finns på tre ställen, och var och
// en hittar dem själv:
//
//   - ristningen: tecken i väggen i ett rum en bit in i tornet,
//   - varelsen: ett tecken under fjällen på en varelse man träffar,
//   - kammaren: ett ord som ormen väser i sömnen när man tar en näve.
//
// Vilken del som finns var lottas varje vandring. Ingen behöver de andra,
// och ingen får veta vad någon annan har hittat. Texterna står under
// delar i data/avslojande.json.

import "strings"

// Källorna, i samma ordning som g.partAt.
const (
	srcCarving = iota
	srcCreature
	srcChamber
)

var partSources = [3]string{"ristning", "varelse", "kammare"}

// carving är rummet med ristningen, när det har bestämts.
type carving struct {
	x, y int
	set  bool
}

// placeParts lottar vilken del av namnet som finns var, när vandringen
// börjar. Ristningens rum bestäms först när tornet bygger ett rum som
// passar (placeCarving).
func (g *Game) placeParts() {
	g.carving = carving{}
	perm := g.rng.Perm(3)
	for src := range 3 {
		g.partAt[src] = perm[src]
	}
}

// carvingFits är sant för ett vanligt rum som ligger en bit in i tornet,
// minst en tredjedel av bredden från varje port.
func (g *Game) carvingFits(x, y int) bool {
	if g.rooms[y][x].Kind != "rum" {
		return false
	}
	far := max(3, g.size/3)
	for _, c := range g.corners() {
		if abs(x-c.x)+abs(y-c.y) < far {
			return false
		}
	}
	return true
}

// placeCarving körs när tornet bygger ett nytt rum. Passar rummet kan
// ristningen hamna där. Den hamnar bara i rum som byggs när någon ser in i
// dem eller går dit, så att den alltid går att nå.
func (g *Game) placeCarving(x, y int) {
	if !g.usesName() || g.carving.set || !g.carvingFits(x, y) {
		return
	}
	if g.rng.Float64() < g.data.Avslojande.RistningChans {
		g.carving = carving{x: x, y: y, set: true}
	}
}

// forceCarving lägger ristningen i ett rum som redan finns, om inget nytt
// rum har passat när nivå 3 kommer: det som ligger längst från portarna och
// där ingen står.
func (g *Game) forceCarving() {
	if !g.usesName() || g.carving.set {
		return
	}
	best, bx, by := -1, 0, 0
	for y := range g.size {
		for x := range g.size {
			r := &g.rooms[y][x]
			if !r.Revealed || r.Kind != "rum" || g.playerIn(x, y) {
				continue
			}
			d := g.size * 4
			for _, c := range g.corners() {
				d = min(d, abs(x-c.x)+abs(y-c.y))
			}
			if d > best {
				best, bx, by = d, x, y
			}
		}
	}
	if best >= 0 {
		g.carving = carving{x: bx, y: by, set: true}
	}
}

// learnPart ger vandraren delen från källan src, om hen inte redan har den.
func (g *Game) learnPart(p *Player, src int) {
	if !g.usesName() || p.bot != nil || p.Status != statusIn {
		return
	}
	i := g.partAt[src]
	if !p.gainPart(i) {
		return
	}
	text := g.namePartText(src, i)
	g.logTo(p, text, "rost")
	g.addCard(p, gameCard{Slag: "namn", Bild: "namndel", Rubrik: g.data.Text("kort_namndel"), Text: text})
}

// readCarving: den som står i ristningens rum läser den.
func (g *Game) readCarving(p *Player) {
	if g.carving.set && p.X == g.carving.x && p.Y == g.carving.y && p.Cellar == 0 {
		g.learnPart(p, srcCarving)
	}
}

// namePartText är texten när vandraren hittar del i av namnet i källan src.
// Texten säger inte rakt ut var i namnet delen står, bara var på ormen:
// vid huvudet, i buken eller i svansen.
func (g *Game) namePartText(src, i int) string {
	a := &g.data.Avslojande
	plats := ""
	if i < len(a.Platser) {
		plats = pick(a.Platser[i])
	}
	return strings.NewReplacer("{plats}", plats, "{del}", strings.ToUpper(g.serpentName[i])).Replace(pick(a.Delar[partSources[src]]))
}

// partsOnLevelUp: från nivå 3 hör den som inte har läst ristningen åt
// vilket håll den finns, och Spårfinnaren redan från nivå 2. På nivå 4 får
// den som har hittat något veta vems namnet är.
func (g *Game) partsOnLevelUp() {
	if !g.usesName() {
		return
	}
	if g.level >= 3 {
		g.forceCarving()
	}
	for _, p := range g.players {
		if p == nil || p.bot != nil || p.Status != statusIn {
			continue
		}
		if g.carving.set && !p.parts[g.partAt[srcCarving]] && (g.level >= 3 || g.is(p, "sparfinnare")) {
			g.tellCarving(p)
		}
		if g.level == 4 && p.partCount() > 0 {
			g.logTo(p, g.data.Text("namn_tornet"), "rost")
		}
	}
}

// tellCarving säger åt vilket håll ristningen finns.
func (g *Game) tellCarving(p *Player) {
	if p.X == g.carving.x && p.Y == g.carving.y {
		return
	}
	g.logTo(p, g.data.Text("ristning_riktning", "horn", g.hornToward(p.X, p.Y, g.carving.x, g.carving.y)), "info")
}
