package main

// Rundans sammanfattning: efter varje runda får var och en ett kort över
// vad som hände just henom, byggt av skillnaden mellan före och efter.
// Det farliga står först. Texterna står i data/texter.json under "sam_".

import "slices"

// before är det som behövs för att se vad rundan ändrade.
type before struct {
	x, y, hp, treasure, potions int
	stuck, poison, dark, cellar int
	status, action              string
	items                       []string
	draws                       int    // ormkort dragna före rundan
	foe                         string // vandraren som stridsvalet riktades mot
}

type summaryLine struct {
	Text string `json:"text"`
	C    string `json:"c,omitempty"` // "varn", "gold" eller tomt
	// Ormkort är id:t på ormkortet vandraren drog, för bilden bredvid raden.
	Ormkort string `json:"ormkort,omitempty"`
}

type summaryMsg struct {
	T     string        `json:"t"`
	Round int           `json:"round"`
	Lines []summaryLine `json:"lines"`
}

// snapshot sparar allas läge innan rundan avgörs.
func (g *Game) snapshot() map[*Player]before {
	m := map[*Player]before{}
	for _, p := range g.players {
		if p == nil || p.client == nil || p.Status != statusIn {
			continue
		}
		var items []string
		for _, it := range p.items {
			items = append(items, it.ID)
		}
		m[p] = before{
			x: p.X, y: p.Y, hp: p.HP, treasure: p.Treasure, potions: p.Potions,
			stuck: p.Stuck, poison: p.Poison, dark: p.Dark, cellar: p.Cellar,
			status: p.Status, action: p.Action, items: items, draws: len(p.stats.draws),
			foe: g.foeName(p),
		}
	}
	return m
}

// foeName är namnet på vandraren som p:s stridsval riktas mot, eller tomt
// om p slåss mot en varelse eller inte slåss alls.
func (g *Game) foeName(p *Player) string {
	if !isCombatAction(p.Action) || g.monsterAt(p.X, p.Y) != nil || g.pvpStatus(p) != pvpYes {
		return ""
	}
	if foe := g.playerByID(p.target); foe != nil && foe != p && foe.Status == statusIn && foe.X == p.X && foe.Y == p.Y {
		return foe.Name
	}
	return ""
}

// summarize skickar rundans sammanfattning till var och en som var i tornet.
func (g *Game) summarize(snap map[*Player]before, wakesBefore int) {
	for p, b := range snap {
		if p.bot != nil && !p.takenOver && p.client == nil {
			continue // en datorspelare behöver ingen krönika på skärmen
		}
		lines := g.summaryFor(p, b, g.wakes > wakesBefore)
		room := p.roundRoom
		p.roundRoom = nil
		if len(lines) > 0 || room != nil {
			// Lines är aldrig nil, så att spelarsidan alltid får en lista.
			jr := journeyRound{Round: g.round, Lines: append([]summaryLine{}, lines...)}
			if room != nil {
				jr.Room, jr.Here, jr.Heard = room.Text, room.Har, room.Hor
			}
			p.journal = append(p.journal, jr)
		}
		if p.client != nil && len(lines) > 0 {
			g.send(p.client, summaryMsg{T: "runda", Round: g.round, Lines: lines})
		}
	}
}

func (g *Game) summaryFor(p *Player, b before, woke bool) []summaryLine {
	var danger, rest []summaryLine
	t := func(key string, vars ...string) summaryLine {
		return summaryLine{Text: g.data.Text(key, vars...)}
	}
	add := func(l summaryLine, c string) {
		l.C = c
		if c == "varn" {
			danger = append(danger, l)
		} else {
			rest = append(rest, l)
		}
	}
	room := ""
	if g.inBounds(p.X, p.Y) {
		room = g.roomShort(&g.rooms[p.Y][p.X])
	}
	moved := p.X != b.x || p.Y != b.y

	// Vad vandraren gjorde, eller vad som hände med henom.
	switch {
	case p.Status == statusOut:
		add(t("sam_ute", "skatt", itoa(p.Treasure)), "gold")
	case p.Status == statusDead && b.status != statusDead:
		add(t("sam_dod"), "varn")
	case p.Cellar > 0 && b.cellar == 0:
		add(t("sam_kallare"), "varn")
	case b.cellar > 0 && p.Cellar == 0:
		add(t("sam_upp", "rum", room), "")
	case b.cellar > 0:
		add(t("sam_i_kallaren"), "")
	default:
		// Vart vandraren gick står i rumsblocket. Kortet säger det bara
		// när steget inte gick som tänkt eller i mörkret, där inget
		// rumsblock kommer. Striden står kvar, eftersom tornet väljer
		// motståndaren, och att man inte valde något.
		_, step := dirCodes[b.action]
		switch {
		case step && (!moved || p.Dark > 0), b.action == "", isCombatAction(b.action):
			add(g.actionLine(p, b, moved, room, t), "")
		}
	}
	if p.Status == statusIn || p.Status == statusDead {
		if d := p.HP - b.hp; d < 0 && p.Status != statusDead {
			add(t("sam_liv_minus", "antal", itoa(-d), "liv", itoa(max(0, p.HP))), "varn")
		} else if d > 0 {
			add(t("sam_liv_plus", "antal", itoa(d), "liv", itoa(p.HP)), "")
		}
	}
	// Ormkortet vandraren drog i kammaren, före guldet.
	for _, id := range p.stats.draws[min(b.draws, len(p.stats.draws)):] {
		if c := g.data.card(id); c != nil {
			l := t("sam_ormkort", "kort", c.Namn)
			l.Ormkort = c.ID
			class := ""
			if c.Vaknar {
				class = "varn"
			}
			add(l, class)
		}
	}
	if p.Status != statusOut {
		if d := p.Treasure - b.treasure; d > 0 {
			add(t("sam_guld_plus", "antal", itoa(d), "skatt", itoa(p.Treasure)), "gold")
		} else if d < 0 {
			add(t("sam_guld_minus", "antal", itoa(-d), "skatt", itoa(p.Treasure)), "varn")
		}
	}
	if p.Potions > b.potions {
		add(t("sam_dryck", "antal", itoa(p.Potions-b.potions)), "")
	}
	for _, it := range p.items {
		if !slices.Contains(b.items, it.ID) {
			if ty := g.data.itemType(it.ID); ty != nil {
				add(t("sam_foremal", "foremal", ty.Kort), "gold")
			}
		}
	}
	if p.Status == statusIn {
		if p.Stuck > 0 && b.stuck == 0 {
			add(t("sam_fast", "antal", itoa(p.Stuck)), "varn")
		}
		if p.Poison > 0 && b.poison == 0 {
			add(t("sam_gift"), "varn")
		}
		if p.Dark > 0 && b.dark == 0 {
			add(t("sam_morker"), "varn")
		}
	}
	if woke {
		add(t("sam_ormen"), "varn")
	}
	return append(danger, rest...)
}

// actionLine säger vad vandraren valde och hur det gick.
func (g *Game) actionLine(p *Player, b before, moved bool, room string, t func(string, ...string) summaryLine) summaryLine {
	a := b.action
	if d, ok := dirCodes[a]; ok {
		dir := g.data.Riktningar[d] // texterna böjer själva, till exempel "{riktning}ut"
		switch {
		case moved && p.Dark > 0:
			return t("sam_gick_morker", "riktning", dir)
		case moved:
			return t("sam_gick", "riktning", dir, "rum", room)
		default:
			return t("sam_stod_kvar", "riktning", dir)
		}
	}
	if moved {
		return t("sam_hamnade", "rum", room)
	}
	switch a {
	case "":
		return t("sam_inget")
	case actAttack, actParry, actDodge:
		if b.foe != "" {
			return t("sam_"+a+"_vandrare", "namn", b.foe)
		}
		return t("sam_" + a)
	case actWait, actTake, actSearch, actDrink:
		return t("sam_" + a)
	}
	return t("sam_vanta")
}
