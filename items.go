package main

// Föremål: saker vandrarna hittar när de söker igenom ett rum och bär med
// sig, högst några stycken. Namn, texter och siffror står i
// data/foremal.json. Här står vad varje föremål gör, styrt av dess id.

import (
	"math/rand/v2"
	"strconv"
)

type ItemType struct {
	ID           string   `json:"id"`
	Namn         string   `json:"namn"`
	Tecken       string   `json:"tecken"` // står framför namnet i väskan
	Bild         string   `json:"bild"`   // spelkortets motiv (web/kort.js); utan bild gäller id:t
	Kort         string   `json:"kort"`   // med liten bokstav, för texterna
	Hittad       string   `json:"hittad"` // obestämd form för texterna: "en extra veke"
	Beskrivning  string   `json:"beskrivning"`
	Vikt         int      `json:"vikt"`
	Anvandningar int      `json:"anvandningar"` // 0 = verkar hela vandringen
	Verkar       []string `json:"verkar"`
	Slut         []string `json:"slut"`
	Forbannad    bool     `json:"forbannad"` // ser värdefullt ut, men har en baksida
	Varde        int      `json:"varde"`     // guld det är värt när man kommer ut
	Forbannelse  []string `json:"forbannelse"`
	Brev         bool     `json:"brev"`    // läses direkt och bärs inte
	Kraft        int      `json:"kraft"`   // läggs till vandrarens Kraft så länge hen bär föremålet
	Forsvar      int      `json:"forsvar"` // läggs till vandrarens Försvar
}

// found är föremålet i obestämd form, som texterna säger att någon hittar.
func (t *ItemType) found() string {
	if t.Hittad != "" {
		return t.Hittad
	}
	return lowerFirst(t.Namn)
}

type ItemData struct {
	MaxBara int                 `json:"max_bara"`
	Lista   []ItemType          `json:"lista"`
	Brev    map[string][]string `json:"brev"` // brev från tidigare vandrare, efter avslöjandets nivå
}

// Item är ett föremål i en väska. uses är hur många gånger det verkar till,
// 0 för föremål som verkar hela vandringen.
type Item struct {
	ID     string
	uses   int
	cursed bool // förbannelsen har visat sig
}

type itemView struct {
	ID          string `json:"id"`
	Tecken      string `json:"tecken,omitempty"`
	Namn        string `json:"namn"`
	Beskrivning string `json:"beskrivning"`
	Kvar        int    `json:"kvar,omitempty"`
}

func (d *Data) itemType(id string) *ItemType {
	for i := range d.Foremal.Lista {
		if d.Foremal.Lista[i].ID == id {
			return &d.Foremal.Lista[i]
		}
	}
	return nil
}

func (d *Data) maxItems() int {
	if d.Foremal.MaxBara > 0 {
		return d.Foremal.MaxBara
	}
	return 3
}

// randomItem drar ett föremål efter vikterna, eller "" om det inte finns några.
func (d *Data) randomItem(rng *rand.Rand) string {
	w := map[string]int{}
	for _, t := range d.Foremal.Lista {
		if t.Vikt > 0 {
			w[t.ID] = t.Vikt
		}
	}
	if len(w) == 0 {
		return ""
	}
	return weighted(rng, w)
}

// has ger föremålet id i vandrarens väska, eller nil.
func (p *Player) has(id string) *Item {
	for _, it := range p.items {
		if it.ID == id {
			return it
		}
	}
	return nil
}

// give lägger föremålet i väskan. Den är falsk när väskan redan är full.
func (g *Game) give(p *Player, id string) bool {
	t := g.data.itemType(id)
	if t == nil {
		return false
	}
	if t.Brev {
		g.readLetter(p)
		return true
	}
	if len(p.items) >= g.data.maxItems() {
		return false
	}
	uses := t.Anvandningar
	p.items = append(p.items, &Item{ID: id, uses: uses})
	p.stats.items++
	g.chronicle("kronika_foremal", "namn", p.Name, "foremal", t.found(), "bild", cardImage(t.Bild, t.ID))
	return true
}

// findItem: vandraren hittar ett föremål när den söker. Är väskan full blir
// föremålet kvar i rummet, åt den som kommer efter.
func (g *Game) findItem(p *Player) {
	id := g.data.randomItem(g.rng)
	t := g.data.itemType(id)
	if t == nil {
		g.soundTo(p, "inget")
		g.logTo(p, g.data.SearchText("inget"), "dim")
		return
	}
	g.soundTo(p, "hitta")
	room := &g.rooms[p.Y][p.X]
	if t.Forbannad && g.is(p, "pilgrim") {
		// Pilgrimen känner att något är fel och låter det ligga.
		room.items = append(room.items, id)
		g.logTo(p, g.data.Text("pilgrim_forbannad", "foremal", t.found()), "info")
		return
	}
	if t.Brev {
		g.give(p, id)
		return
	}
	if g.give(p, id) {
		g.logTo(p, g.data.Text("foremal_hittar", "foremal", t.found(), "beskrivning", t.Beskrivning), "syn")
		g.itemCard(p, t)
		return
	}
	room.items = append(room.items, id)
	g.logTo(p, g.data.Text("foremal_full", "foremal", t.found(), "max", strconv.Itoa(g.data.maxItems())), "info")
}

// pickUp: den som kommer in i ett rum där något ligger tar det, om det får
// plats i väskan.
func (g *Game) pickUp(p *Player) {
	if p.Status != statusIn || p.Cellar > 0 {
		return
	}
	g.pickUpGold(p)
	room := &g.rooms[p.Y][p.X]
	var left []string
	for _, id := range room.items {
		t := g.data.itemType(id)
		if t == nil {
			continue
		}
		if t.Forbannad && g.is(p, "pilgrim") {
			left = append(left, id)
			g.logTo(p, g.data.Text("pilgrim_forbannad", "foremal", t.found()), "dim")
			continue
		}
		if g.give(p, id) {
			g.logTo(p, g.data.Text("foremal_ligger", "foremal", t.found(), "beskrivning", t.Beskrivning), "syn")
			g.itemCard(p, t)
		} else {
			left = append(left, id)
			g.logTo(p, g.data.Text("foremal_ligger_full", "foremal", t.found()), "dim")
		}
	}
	room.items = left
}

// useItem förbrukar en användning av föremålet id, om vandraren har det, och
// berättar vad som händer. Föremål som verkar hela vandringen förbrukas inte.
func (g *Game) useItem(p *Player, id string, vars ...string) bool {
	it := p.has(id)
	if it == nil {
		return false
	}
	t := g.data.itemType(id)
	if t != nil {
		if text := pick(t.Verkar); text != "" {
			g.logTo(p, fill(text, vars...), "info")
		}
	}
	if it.uses > 0 {
		it.uses--
		if it.uses == 0 {
			if t != nil {
				if text := pick(t.Slut); text != "" {
					g.logTo(p, text, "dim")
				}
			}
			g.dropItem(p, it)
		}
	}
	return true
}

func (g *Game) dropItem(p *Player, it *Item) {
	for i, o := range p.items {
		if o == it {
			p.items = append(p.items[:i], p.items[i+1:]...)
			return
		}
	}
}

// hidden halverar en förlust för den som bär väskan med dubbel botten.
func (g *Game) hiddenShare(p *Player, loss int) int {
	if loss <= 0 || p.has("vaska") == nil {
		return loss
	}
	kept := loss / 2
	g.useItem(p, "vaska", "antal", strconv.Itoa(kept))
	return loss - kept
}

// shieldWear: skölden ger Försvar tills den har tagit tre träffar från
// varelser.
func (g *Game) shieldWear(p *Player) {
	if p.has("skold") != nil {
		g.useItem(p, "skold")
	}
}

// itemBonus är det föremålen i väskan lägger till vandrarens Kraft och
// Försvar.
func (g *Game) itemBonus(p *Player) (kraft, forsvar int) {
	for _, it := range p.items {
		if t := g.data.itemType(it.ID); t != nil {
			kraft += t.Kraft
			forsvar += t.Forsvar
		}
	}
	return kraft, forsvar
}

// slowFeet: den som flyr blir alltid träffad av slaget mot ryggen, om hen
// är kartritaren eller bär brynjan.
func (g *Game) slowFeet(p *Player) bool {
	return g.is(p, "kartritare") || p.has("brynja") != nil
}

// thirst: svärdet som törstar tar ett liv varje runda i strid då
// vandraren inte anfaller.
func (g *Game) thirst(p *Player) {
	if p.has("tarst_svard") == nil || p.Status != statusIn {
		return
	}
	g.curse(p, "tarst_svard")
	t := g.data.itemType("tarst_svard")
	g.logTo(p, g.data.Text("tarst_tar"), "varn")
	p.hurt("foremal", t.Kort)
	g.damage(p, 1, false)
}

// compass är vad kompassen visar: riktningen mot närmaste port, som ett av
// orden under "horn" i texter.json.
func (g *Game) compass(p *Player) string {
	if p.has("kompass") == nil || p.Status != statusIn || p.Cellar > 0 {
		return ""
	}
	best, bx, by := -1, 0, 0
	for _, cn := range g.corners() {
		d := abs(cn.x-p.X) + abs(cn.y-p.Y)
		if best < 0 || d < best {
			best, bx, by = d, cn.x, cn.y
		}
	}
	if best == 0 {
		return g.data.Text("kompass_har")
	}
	sx, sy := sign(bx-p.X), sign(by-p.Y)
	// Samma ordning som "horn": nordväst, nordöst, sydväst, sydöst, norr, söder, väster, öster.
	idx := map[[2]int]int{{-1, -1}: 0, {1, -1}: 1, {-1, 1}: 2, {1, 1}: 3, {0, -1}: 4, {0, 1}: 5, {-1, 0}: 6, {1, 0}: 7}[[2]int{sx, sy}]
	return g.data.Text("kompass_pekar", "horn", g.data.Horn[idx])
}

func sign(v int) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}

func (g *Game) itemViews(p *Player) []itemView {
	var out []itemView
	for _, it := range p.items {
		if t := g.data.itemType(it.ID); t != nil {
			out = append(out, itemView{ID: t.ID, Tecken: t.Tecken, Namn: t.Namn, Beskrivning: t.Beskrivning, Kvar: it.uses})
		}
	}
	return out
}

// itemIDs är id för vandrarens föremål, för figuren på spelledarens karta.
func (g *Game) itemIDs(p *Player) []string {
	var out []string
	for _, it := range p.items {
		out = append(out, it.ID)
	}
	return out
}

// Föremålen som förr användes med 1–3 verkar nu av sig själva, som de
// andra: rökpulvret när vandraren lämnar en strid, spegelskärvan när hen
// kommer in i ett nytt rum, kastkniven när en varelse dyker upp och
// timglaset när solen är nästan nere. Kritan togs bort 2026-11-22.

// neighbor ger rummet åt hållet d, om det finns en dörr dit.
func (g *Game) neighbor(p *Player, d Dir) (int, int, bool) {
	if !g.rooms[p.Y][p.X].Doors[d] {
		return 0, 0, false
	}
	x, y := p.X+dirDX[d], p.Y+dirDY[d]
	return x, y, g.inBounds(x, y)
}

// smoke: den som flyr ur en strid med rökpulver i väskan kommer undan utan
// att varelsen hinner slå. Svarar true när pulvret verkade.
func (g *Game) smoke(p *Player, d Dir) bool {
	if p.has("rokpulver") == nil {
		return false
	}
	if _, _, ok := g.neighbor(p, d); !ok {
		return false
	}
	return g.useItem(p, "rokpulver", "riktning", g.data.Riktningar[d])
}

// mirror: i ett nytt rum visar spegelskärvan ett grannrum bakom en dörr,
// ett som vandraren inte har sett.
func (g *Game) mirror(p *Player) {
	if p.has("spegel") == nil || p.Dark > 0 {
		return
	}
	for _, d := range g.rng.Perm(4) {
		x, y, ok := g.neighbor(p, Dir(d))
		if !ok || p.known[y][x] {
			continue
		}
		g.useItem(p, "spegel", "riktning", g.data.Riktningar[d])
		r := g.reveal(x, y)
		p.known[y][x] = true
		g.logTo(p, g.data.Text("spegel_ser", "rum", g.roomShort(r)), "info")
		if r.Event != "" {
			if trap := g.data.trap(r.Event); trap != nil {
				g.logTo(p, g.data.Text("spegel_falla", "falla", trap.Kort), "varn")
			}
		}
		if mm := g.monsterAt(x, y); mm != nil {
			g.logTo(p, g.data.Text("spegel_varelse", "varelse", g.data.monsterType(mm.Type).Kort), "varn")
		}
		for _, o := range g.players {
			if o != nil && o != p && o.Status == statusIn && o.Cellar == 0 && o.X == x && o.Y == y {
				g.logTo(p, g.data.Text("spegel_vandrare", "namn", o.Name), "info")
			}
		}
		return
	}
}

// knife: när en varelse dyker upp kastar den som har en kastkniv den
// först, och varelsen tar en skada innan striden börjar.
func (g *Game) knife(p *Player, m *Monster) {
	if p.has("kastkniv") == nil {
		return
	}
	mt := g.data.monsterType(m.Type)
	g.useItem(p, "kastkniv", "varelse", mt.KortBestamd)
	g.logTo(p, g.data.Text("kniv_traff", "varelse", mt.KortBestamd), "syn")
	m.HP -= g.data.Monster.Spelarskada
	if m.HP <= 0 {
		g.defeat(p, m)
	}
}

// flute: den sövande flöjten låter den som drar ett vakna-kort dra om.
func (g *Game) flute(p *Player, i int) int {
	if c := g.data.card(g.deck[i]); c == nil || !c.Vaknar || !g.useItem(p, "flojt") {
		return i
	}
	return g.rng.IntN(len(g.deck))
}

// readLetter: ett brev från en tidigare vandrare läses direkt. Texten följer
// avslöjandets nivå och säger aldrig hemligheten rakt ut.
func (g *Game) readLetter(p *Player) {
	for l := g.level; l >= 1; l-- {
		if list := g.data.Foremal.Brev[strconv.Itoa(l)]; len(list) > 0 {
			g.logTo(p, g.data.Text("brev_hittar"), "syn")
			text := pick(list)
			g.logTo(p, text, "rost")
			if tip := g.letterTip(p); tip != "" {
				g.logTo(p, tip, "info")
				text += " " + tip
			}
			if t := g.data.itemType("brev"); t != nil {
				g.addCard(p, gameCard{Slag: "brev", Bild: cardImage(t.Bild, t.ID), Rubrik: t.Namn, Text: text})
			}
			p.stats.items++
			return
		}
	}
}

// letterTip är det brevet pekar mot i just det här tornet: ristningen, om
// den finns och vandraren inte har läst den, annars närmaste varelse, och
// annars kammaren.
func (g *Game) letterTip(p *Player) string {
	if g.carving.set && g.usesName() && !p.parts[g.partAt[srcCarving]] && (p.X != g.carving.x || p.Y != g.carving.y) {
		return g.data.Text("brev_tips_ristning", "horn", g.hornToward(p.X, p.Y, g.carving.x, g.carving.y))
	}
	var near *Monster
	best := 0
	for _, m := range g.monsters {
		d := abs(m.X-p.X) + abs(m.Y-p.Y)
		if d > 0 && (near == nil || d < best) {
			near, best = m, d
		}
	}
	if near != nil {
		return g.data.Text("brev_tips_varelse", "horn", g.hornToward(p.X, p.Y, near.X, near.Y))
	}
	cx, cy := g.chamber()
	if cx == p.X && cy == p.Y {
		return ""
	}
	return g.data.Text("brev_tips_kammare", "horn", g.hornToward(p.X, p.Y, cx, cy))
}

// itemCard lägger föremålets spelkort: namnet och vad det gör, som inte
// står i loggen.
func (g *Game) itemCard(p *Player, t *ItemType) {
	g.addCard(p, gameCard{Slag: "foremal", Bild: cardImage(t.Bild, t.ID), Rubrik: t.Namn, Text: t.Beskrivning})
}

// curse: ett förbannat föremål visar sin baksida första gången den verkar.
func (g *Game) curse(p *Player, id string) bool {
	it := p.has(id)
	if it == nil {
		return false
	}
	if !it.cursed {
		it.cursed = true
		if t := g.data.itemType(id); t != nil {
			text := pick(t.Forbannelse)
			if text != "" {
				g.logTo(p, text, "varn")
			}
			g.addCard(p, gameCard{Slag: "forbannelse", Bild: cardImage(t.Bild, t.ID), Rubrik: t.Namn, Text: text, Forbannad: true})
		}
	}
	return true
}

// itemGold är guldet föremålen i väskan är värda när vandraren kommer ut.
func (g *Game) itemGold(p *Player) int {
	n := 0
	for _, it := range p.items {
		if t := g.data.itemType(it.ID); t != nil {
			n += t.Varde
		}
	}
	return n
}

// dropAll: den som dör tappar sina föremål i rummet, åt den som kommer efter.
func (g *Game) dropAll(p *Player) {
	if len(p.items) == 0 || !g.inBounds(p.X, p.Y) {
		return
	}
	r := &g.rooms[p.Y][p.X]
	for _, it := range p.items {
		r.items = append(r.items, it.ID)
	}
	p.items = nil
}
