package main

import "strings"

// Rummet i krönikan: när en vandrare kommer in i ett rum får hen ett
// rumsblock med rummets namn, beskrivningen när rummet är nytt och en mening
// om vad som finns där just nu. Står vandraren kvar kommer ett nytt block
// bara när det som finns i rummet har ändrats. Blocken skickas i slutet av
// rundan, så att en fälla som flyttar vandraren hinner göra det först, och
// spelarsidan ställer dem överst i rundan.

// roomMsg är ett rumsblock till spelarsidan.
type roomMsg struct {
	T    string `json:"t"`
	Namn string `json:"namn"`
	Text string `json:"text,omitempty"` // beskrivningen, bara när rummet är nytt
	Har  string `json:"har,omitempty"`  // det som finns i rummet
	Kvar bool   `json:"kvar,omitempty"` // vandraren står kvar och något har ändrats
}

// told är vad en vandrare senast fick veta om sitt rum.
type told struct {
	ok   bool
	x, y int
	har  string
	desc string // beskrivningen som väntar på nästa block
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
// annat rum eller om det som finns i rummet har ändrats. I mörkret och i
// källarvalven ser vandraren inget, och då väntar blocket.
func (g *Game) tellRoom(p *Player) {
	if p.Status != statusIn || p.Cellar > 0 || p.Dark > 0 || !g.inBounds(p.X, p.Y) {
		return
	}
	r := &g.rooms[p.Y][p.X]
	har := g.hereText(p)
	moved := !p.told.ok || p.told.x != p.X || p.told.y != p.Y
	if !moved && har == p.told.har {
		return
	}
	msg := roomMsg{T: "rum", Namn: g.roomName(r), Har: har, Kvar: !moved}
	if moved {
		msg.Text = p.told.desc
	}
	g.send(p.client, msg)
	p.told = told{ok: true, x: p.X, y: p.Y, har: har}
}

// forgetRoom gör att nästa rumsblock börjar om, som när vandraren ställs vid
// en port.
func (p *Player) forgetRoom() {
	p.told = told{}
}

// hereText säger i ord vad som finns i vandrarens rum: varelsen, de andra
// vandrarna, skattkammaren, porten, en fälla som har slagit till och
// föremål som någon har lämnat. Texterna är har_* i data/texter.json.
func (g *Game) hereText(p *Player) string {
	r := &g.rooms[p.Y][p.X]
	var parts []string
	add := func(key string, vars ...string) {
		if s := g.data.Text(key, vars...); s != "" {
			parts = append(parts, s)
		}
	}
	if m := g.monsterAt(p.X, p.Y); m != nil {
		if t := g.data.monsterType(m.Type); t != nil {
			add("har_varelse", "varelse", t.Namn)
		}
	}
	var others []string
	for _, o := range g.players {
		if o != nil && o != p && o.X == p.X && o.Y == p.Y && g.sees(p, o) {
			others = append(others, o.Name)
		}
	}
	switch len(others) {
	case 0:
	case 1:
		add("har_vandrare", "namn", others[0])
	default:
		add("har_vandrare_flera", "namn", strings.Join(others[:len(others)-1], ", "), "sista", others[len(others)-1])
	}
	if r.Kind == "kammare" {
		add("har_kammare")
	}
	if r.HasExit {
		if r.Shut {
			add("har_port_stangd")
		} else {
			add("har_port")
		}
	}
	if r.Mark != "" {
		if t := g.data.trapByMark(r.Mark); t != nil {
			add("har_falla", "falla", t.Namn)
		}
	}
	for _, id := range r.items {
		if t := g.data.itemType(id); t != nil {
			add("har_foremal", "foremal", lowerFirst(t.Namn))
		}
	}
	if len(parts) == 0 {
		add("har_tomt")
	}
	return strings.Join(parts, " ")
}
