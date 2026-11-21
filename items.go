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
	Kort         string   `json:"kort"` // med liten bokstav, för texterna
	Beskrivning  string   `json:"beskrivning"`
	Vikt         int      `json:"vikt"`
	Anvandningar int      `json:"anvandningar"` // 0 = verkar hela vandringen
	Verkar       []string `json:"verkar"`
	Slut         []string `json:"slut"`
	Aktiv        bool     `json:"aktiv"`     // används med 1–3 i stället för att verka av sig självt
	Riktning     bool     `json:"riktning"`  // pekar åt ett håll, som väljs med pilarna
	Forbannad    bool     `json:"forbannad"` // ser värdefullt ut, men har en baksida
	Varde        int      `json:"varde"`     // guld det är värt när man kommer ut
	Forbannelse  []string `json:"forbannelse"`
	Brev         bool     `json:"brev"` // läses direkt och bärs inte
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
	Aktiv       bool   `json:"aktiv,omitempty"`
	Riktning    bool   `json:"riktning,omitempty"`
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
	if id == "krita" && g.is(p, "kartritare") {
		uses += int(g.classVal(p, "extra_krita", 3))
	}
	p.items = append(p.items, &Item{ID: id, uses: uses})
	p.stats.items++
	g.chronicle("kronika_foremal", "namn", p.Name, "foremal", lowerFirst(t.Namn))
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
		g.logTo(p, g.data.Text("pilgrim_forbannad", "foremal", lowerFirst(t.Namn)), "info")
		return
	}
	if t.Brev {
		g.give(p, id)
		return
	}
	if g.give(p, id) {
		g.logTo(p, g.data.Text("foremal_hittar", "foremal", lowerFirst(t.Namn), "beskrivning", t.Beskrivning), "syn")
		return
	}
	room.items = append(room.items, id)
	g.logTo(p, g.data.Text("foremal_full", "foremal", lowerFirst(t.Namn), "max", strconv.Itoa(g.data.maxItems())), "info")
}

// pickUp: den som kommer in i ett rum där något ligger tar det, om det får
// plats i väskan.
func (g *Game) pickUp(p *Player) {
	if p.Status != statusIn || p.Cellar > 0 {
		return
	}
	room := &g.rooms[p.Y][p.X]
	var left []string
	for _, id := range room.items {
		t := g.data.itemType(id)
		if t == nil {
			continue
		}
		if t.Forbannad && g.is(p, "pilgrim") {
			left = append(left, id)
			g.logTo(p, g.data.Text("pilgrim_forbannad", "foremal", lowerFirst(t.Namn)), "dim")
			continue
		}
		if g.give(p, id) {
			g.logTo(p, g.data.Text("foremal_ligger", "foremal", lowerFirst(t.Namn), "beskrivning", t.Beskrivning), "syn")
		} else {
			left = append(left, id)
			g.logTo(p, g.data.Text("foremal_ligger_full", "foremal", lowerFirst(t.Namn)), "dim")
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

// shieldDamage: skölden tar en skada av varje träff från en varelse.
func (g *Game) shieldDamage(p *Player, dmg int) int {
	if dmg <= 0 || p.has("skold") == nil {
		return dmg
	}
	g.useItem(p, "skold")
	return dmg - 1
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
			out = append(out, itemView{ID: t.ID, Tecken: t.Tecken, Namn: t.Namn, Beskrivning: t.Beskrivning, Kvar: it.uses, Aktiv: t.Aktiv, Riktning: t.Riktning})
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

// itemNames är vandrarens föremål för spelledarens översikt.
func (g *Game) itemNames(p *Player) []string {
	var out []string
	for _, it := range p.items {
		if t := g.data.itemType(it.ID); t != nil {
			out = append(out, t.Kort)
		}
	}
	return out
}

// Föremål som vandraren använder själv. Draget skrivs "f" och platsen i
// väskan (1–3), och för föremål som pekar åt ett håll en riktning efter:
// "f2E" kastar det andra föremålet österut.

func isItemAction(a string) bool {
	return len(a) >= 2 && a[0] == 'f' && a[1] >= '1' && a[1] <= '9'
}

// parseItemAction ger föremålet och riktningen i ett drag, eller nil.
func (g *Game) parseItemAction(p *Player, a string) (*Item, *ItemType, Dir, bool) {
	if !isItemAction(a) {
		return nil, nil, 0, false
	}
	slot := int(a[1] - '1')
	if slot >= len(p.items) {
		return nil, nil, 0, false
	}
	it := p.items[slot]
	t := g.data.itemType(it.ID)
	if t == nil || !t.Aktiv {
		return nil, nil, 0, false
	}
	if !t.Riktning {
		return it, t, 0, len(a) == 2
	}
	if len(a) != 3 {
		return nil, nil, 0, false
	}
	d, ok := dirCodes[a[2:]]
	return it, t, d, ok
}

// neighbor ger rummet åt hållet d, om det finns en dörr dit.
func (g *Game) neighbor(p *Player, d Dir) (int, int, bool) {
	if !g.rooms[p.Y][p.X].Doors[d] {
		return 0, 0, false
	}
	x, y := p.X+dirDX[d], p.Y+dirDY[d]
	return x, y, g.inBounds(x, y)
}

// useActive utför ett föremål som vandraren har valt att använda.
func (g *Game) useActive(p *Player, a string) {
	it, t, d, ok := g.parseItemAction(p, a)
	if !ok {
		return
	}
	m := g.monsterAt(p.X, p.Y)
	if m != nil && t.ID != "rokpulver" {
		// Den som rotar i väskan mitt i en strid ger varelsen ett fritt drag.
		g.freeAttack(p, m, "inte_beredd")
		if p.Status != statusIn {
			return
		}
	}
	riktning := g.data.Riktningar[d]
	switch t.ID {
	case "rokpulver":
		if _, _, ok := g.neighbor(p, d); !ok {
			g.logTo(p, g.data.Text("vagg", "riktning", riktning), "dim")
			return
		}
		g.useItem(p, it.ID, "riktning", riktning)
		g.move(p, d, 0)
	case "spegel":
		x, y, ok := g.neighbor(p, d)
		if !ok {
			g.logTo(p, g.data.Text("spegel_vagg", "riktning", riktning), "dim")
			return
		}
		g.useItem(p, it.ID, "riktning", riktning)
		r := g.reveal(x, y)
		p.known[y][x] = true
		g.logTo(p, g.data.Text("spegel_ser", "rum", g.roomName(r)), "info")
		if r.Event != "" {
			if trap := g.data.trap(r.Event); trap != nil {
				g.logTo(p, g.data.Text("spegel_falla", "falla", lowerFirst(trap.Namn)), "varn")
			}
		}
		if mm := g.monsterAt(x, y); mm != nil {
			g.logTo(p, g.data.Text("spegel_varelse", "varelse", lowerFirst(g.data.monsterType(mm.Type).Namn)), "varn")
		}
		for _, o := range g.players {
			if o != nil && o != p && o.Status == statusIn && o.Cellar == 0 && o.X == x && o.Y == y {
				g.logTo(p, g.data.Text("spegel_vandrare", "namn", o.Name), "info")
			}
		}
	case "kastkniv":
		x, y, ok := g.neighbor(p, d)
		if !ok {
			g.logTo(p, g.data.Text("vagg", "riktning", riktning), "dim")
			return
		}
		g.useItem(p, it.ID, "riktning", riktning)
		mm := g.monsterAt(x, y)
		if mm == nil {
			// Kniven blir liggande, åt den som hittar den.
			g.logTo(p, g.data.Text("kniv_miss", "riktning", riktning), "dim")
			g.rooms[y][x].items = append(g.rooms[y][x].items, "kastkniv")
			return
		}
		mt := g.data.monsterType(mm.Type)
		g.logTo(p, g.data.Text("kniv_traff", "varelse", lowerFirst(mt.Namn)), "syn")
		mm.HP -= g.data.Monster.Spelarskada
		if mm.HP <= 0 {
			g.defeat(p, mm)
		}
	case "krita":
		g.useItem(p, it.ID, "riktning", riktning)
		r := &g.rooms[p.Y][p.X]
		r.chalk, r.hasChalk = d, true
	case "timglas":
		g.useItem(p, it.ID)
		g.sun++ // rundan räknar ner solen igen; den står still
		g.logOthers(p, g.data.Text("timglas_alla"), "info")
	}
}

// noticeChalk: den som kommer in i ett rum med en pil i krita ser den.
func (g *Game) noticeChalk(p *Player) {
	r := &g.rooms[p.Y][p.X]
	if r.hasChalk && p.Dark == 0 {
		g.logTo(p, g.data.Text("krita_ser", "riktning", g.data.Riktningar[r.chalk]), "info")
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
	level := g.level
	if !g.cfg.Reveal {
		level = 1
	}
	for l := level; l >= 1; l-- {
		if list := g.data.Foremal.Brev[strconv.Itoa(l)]; len(list) > 0 {
			g.logTo(p, g.data.Text("brev_hittar"), "syn")
			g.logTo(p, pick(list), "rost")
			p.stats.items++
			return
		}
	}
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
			if text := pick(t.Forbannelse); text != "" {
				g.logTo(p, text, "varn")
			}
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
