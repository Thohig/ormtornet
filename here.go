package main

import (
	"slices"
	"strconv"
	"strings"
)

// Rummet i krönikan: när en vandrare kommer in i ett rum får hen ett
// rumsblock med rummets namn, beskrivningen när rummet är nytt och en mening
// om vad som finns där just nu. Står vandraren kvar kommer ett nytt block
// bara när det som finns i rummet har ändrats. Blocken skickas i slutet av
// rundan, så att en fälla som flyttar vandraren hinner göra det först, och
// spelarsidan ställer dem överst i rundan.

// roomMsg är ett rumsblock till spelarsidan. Rummet har inget namn: det är
// beskrivningen och det man ser och hör.
type roomMsg struct {
	T    string `json:"t"`
	Text string `json:"text,omitempty"` // beskrivningen, bara när rummet är nytt
	Har  string `json:"har,omitempty"`  // det som finns i rummet
	Kvar bool   `json:"kvar,omitempty"` // vandraren står kvar och något har ändrats
	Hor  string `json:"hor,omitempty"`  // det vandraren hörde under rundan
}

// told är vad en vandrare senast fick veta om sitt rum.
type told struct {
	ok    bool
	x, y  int
	parts map[string]herePart // det som syntes, på nyckel
	desc  string              // beskrivningen som väntar på nästa block
	extra string              // det man också ser, bara i nästa block (atTheGate)
}

// herePart är en sak man ser i rummet: en varelse, en vandrare eller
// något annat. gone är vad som sägs när den har försvunnit, för varelser
// och vandrare.
type herePart struct {
	key, kind, name, text, gone string
}

// describeLater sparar rummets beskrivning till nästa rumsblock.
func (g *Game) describeLater(p *Player, r *Room, x, y int) {
	p.told.desc = g.describe(p, r, x, y)
}

// tellRooms skickar rumsblocken i slutet av en runda.
func (g *Game) tellRooms() {
	for _, p := range g.players {
		if p != nil {
			g.tellRoom(p)
		}
	}
}

// tellRoom skickar ett rumsblock till vandraren om hen har kommit till ett
// annat rum eller om det som syns i rummet har ändrats. Står hen kvar säger
// blocket bara det som har ändrats: det som har kommit och det som har
// gått. I mörkret och i källarvalven ser vandraren inget, och då väntar
// blocket.
func (g *Game) tellRoom(p *Player) {
	heard := p.heard
	p.heard = nil
	p.fightsBefore, p.fights = p.fights, nil
	if p.Status != statusIn || p.Cellar > 0 || p.Dark > 0 || !g.inBounds(p.X, p.Y) {
		// Utan rumsblock blir ljuden egna rader, som förut.
		for _, h := range heard {
			g.logTo(p, h, "dim")
		}
		return
	}
	all := g.hereParts(p)
	now := map[string]herePart{}
	var cur []herePart
	for _, part := range all {
		now[part.key] = part
		// Spåret efter en fälla som nyss slog till mot vandraren står inte
		// i blocket: loggen har redan sagt vad som hände, och blocket står
		// överst i rundan, före den raden.
		if p.trapNow == "" || part.key != "falla:"+p.trapNow {
			cur = append(cur, part)
		}
	}
	p.trapNow = ""
	moved := !p.told.ok || p.told.x != p.X || p.told.y != p.Y
	msg := roomMsg{T: "rum", Kvar: !moved, Hor: strings.Join(heard, " ")}
	if moved {
		msg.Text = p.told.desc
		msg.Har = strings.TrimSpace(g.sayHere(cur) + " " + p.told.extra)
	} else {
		var added []herePart
		for _, part := range cur {
			if _, ok := p.told.parts[part.key]; !ok {
				added = append(added, part)
			}
		}
		var gone []string
		for key, part := range p.told.parts {
			if _, ok := now[key]; ok || part.gone == "" {
				continue
			}
			if part.kind == "varelse" {
				// En fälld varelse är inte borta: den som fällde den har
				// redan läst det i loggen, och de andra ser den ligga där.
				id, _ := strconv.Atoi(strings.TrimPrefix(key, "varelse:"))
				if by, ok := g.slain[id]; ok {
					if by != p.ID {
						gone = append(gone, g.data.Text("har_varelse_falld", "varelse", upperFirst(part.name)))
					}
					continue
				}
			}
			gone = append(gone, part.gone)
		}
		if len(added) == 0 && len(gone) == 0 && len(heard) == 0 {
			return
		}
		slices.Sort(gone)
		msg.Har = strings.TrimSpace(strings.Join(append(gone, g.sayHere(added)), " "))
	}
	p.told = told{ok: true, x: p.X, y: p.Y, parts: now}
	if msg.Text == "" && msg.Har == "" && msg.Hor == "" {
		return // ett känt, tomt rum: blocket skulle inte säga något
	}
	g.send(p.client, msg)
	p.roundRoom = &msg
}

// retellRoom skickar rummet vandraren står i en gång till, när hen kommer
// tillbaka mitt i vandringen efter att ha laddat om sidan: loggen börjar
// om, och utan blocket vet hen inte var hen är.
func (g *Game) retellRoom(p *Player) {
	if g.phase != phasePlay || p.Status != statusIn || p.Cellar > 0 || p.Dark > 0 || !g.inBounds(p.X, p.Y) {
		return
	}
	r := &g.rooms[p.Y][p.X]
	msg := roomMsg{T: "rum", Text: g.describe(p, r, p.X, p.Y), Har: g.sayHere(g.hereParts(p))}
	g.send(p.client, msg)
}

// atTheGate gör början av vandringen till en del av det första rumsblocket:
// vid vilken port vandraren kliver in står först i beskrivningen, solen
// hör till det man ser, och portarna där de andra går in är det man hör.
func (g *Game) atTheGate(p *Player) {
	p.told.desc = strings.TrimSpace(g.data.Text("start_port", "horn", g.data.Horn[p.Corner]) + " " + p.told.desc)
	p.told.extra = g.data.Text("start", "sol", strconv.Itoa(g.sun))
	others := 0
	for _, o := range g.players {
		if o != nil && o != p {
			others++
		}
	}
	switch {
	case others == 1:
		p.heard = append(p.heard, g.data.Text("start_hor_en"))
	case others > 1:
		p.heard = append(p.heard, g.data.Text("start_hor_flera"))
	}
}

// forgetRoom gör att nästa rumsblock börjar om, som när vandraren ställs vid
// en port.
func (p *Player) forgetRoom() {
	p.told = told{}
}

// hereParts är det vandraren ser i sitt rum, beskrivet och inte förklarat:
// varelsens egen text, de andra vandrarna, guldet i kammaren, porten,
// spåret efter en fälla och föremål som någon har lämnat. Ett tomt rum
// säger ingenting. Texterna är har_* i data/texter.json och spar_efter i
// data/fallor.json.
func (g *Game) hereParts(p *Player) []herePart {
	r := &g.rooms[p.Y][p.X]
	var parts []herePart
	other := func(key, text string) {
		if text != "" {
			parts = append(parts, herePart{key: key, kind: "annat", text: text})
		}
	}
	if m := g.monsterAt(p.X, p.Y); m != nil {
		if t := g.data.monsterType(m.Type); t != nil && len(t.DykerUpp) > 0 {
			// Samma text så länge varelsen står kvar, så att blocket inte
			// tror att något har ändrats.
			parts = append(parts, herePart{key: "varelse:" + strconv.Itoa(m.ID), kind: "varelse", name: t.KortBestamd,
				text: t.DykerUpp[m.ID%len(t.DykerUpp)], gone: g.data.Text("har_varelse_borta", "varelse", upperFirst(t.KortBestamd))})
		}
	}
	for _, o := range g.players {
		if o != nil && o != p && o.X == p.X && o.Y == p.Y && g.sees(p, o) {
			parts = append(parts, herePart{key: "vandrare:" + strconv.Itoa(o.ID), kind: "vandrare", name: o.Name,
				gone: g.data.Text("har_vandrare_borta", "namn", o.Name)})
		}
	}
	if r.Kind == "kammare" {
		other("kammare", g.data.Text("har_kammare"))
	}
	if r.HasExit {
		if r.Shut {
			other("port_stangd", g.data.Text("har_port_stangd"))
		} else {
			other("port", g.data.Text("har_port"))
		}
	}
	if r.Mark != "" {
		if t := g.data.trapByMark(r.Mark); t != nil && len(t.SparEfter) > 0 {
			other("falla:"+t.ID, t.SparEfter[0])
		}
	}
	// Där någon har dött syns det, och att guldet är taget.
	for _, name := range r.fell {
		other("fall:"+name, g.data.Text("har_fall", "namn", name))
	}
	if r.goldBy != 0 && r.goldBy != p.ID && len(r.gold) == 0 {
		// Spårfinnaren läser i dammet vem som tog mynten.
		if o := g.playerByID(r.goldBy); o != nil && g.is(p, "sparfinnare") {
			other("guld_taget", g.data.Text("har_guld_taget_namn", "namn", o.Name))
		} else {
			other("guld_taget", g.data.Text("har_guld_taget"))
		}
	}
	for i, id := range r.items {
		if t := g.data.itemType(id); t != nil {
			other("foremal:"+strconv.Itoa(i)+":"+id, g.data.Text("har_foremal", "foremal", t.found()))
		}
	}
	return parts
}

// sayHere skriver delarna som löpande text: varelsen först, sedan
// vandrarna i en mening och sist resten.
func (g *Game) sayHere(parts []herePart) string {
	var out, names []string
	for _, part := range parts {
		if part.kind == "varelse" {
			out = append(out, part.text)
		}
	}
	for _, part := range parts {
		if part.kind == "vandrare" {
			names = append(names, part.name)
		}
	}
	switch len(names) {
	case 0:
	case 1:
		out = append(out, g.data.Text("har_vandrare", "namn", names[0]))
	default:
		out = append(out, g.data.Text("har_vandrare_flera", "namn", strings.Join(names[:len(names)-1], ", "), "sista", names[len(names)-1]))
	}
	for _, part := range parts {
		if part.kind == "annat" {
			out = append(out, part.text)
		}
	}
	return strings.Join(out, " ")
}
