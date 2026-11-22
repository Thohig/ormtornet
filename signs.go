package main

// Tecknen från kartan på krönikans rader och på rundans kort: samma bilder
// som i teckenförklaringen, så att det man läser känns igen från kartan.
// Spelarsidan ritar bilden med kartans egen funktion och visar tecknet som
// text där ingen bild går att rita.

// mapSign är ett tecken bredvid en rad. Key är samma nyckel som i
// teckenförklaringen (varelse:<tecken>, falla:<tecken>, kammare, port,
// port_stangd, vandrare) eller foremal:<id>. ID och Model säger vilken
// modell som ska ritas, Slot vilken färg en vandrare har, och Name är
// förklaringen när man pekar på tecknet. Utan namn används tecken_* i
// data/texter.json, som i teckenförklaringen.
type mapSign struct {
	Key   string     `json:"key"`
	Glyph string     `json:"glyph,omitempty"`
	ID    string     `json:"id,omitempty"`
	Model *ModelLook `json:"modell,omitempty"`
	Slot  *int       `json:"slot,omitempty"`
	Name  string     `json:"namn,omitempty"`
}

func monsterSign(t *MonsterType) mapSign {
	return mapSign{Key: "varelse:" + t.Tecken, Glyph: t.Tecken, ID: t.ID, Model: t.Modell, Name: t.Namn}
}

func trapSign(t *Trap) mapSign {
	return mapSign{Key: "falla:" + t.Tecken, Glyph: t.Tecken, ID: t.ID, Name: t.Namn}
}

func itemSign(t *ItemType) mapSign {
	return mapSign{Key: "foremal:" + t.ID, Glyph: t.Tecken, Name: t.Namn}
}

func playerSign(p *Player) mapSign {
	slot := p.Slot
	return mapSign{Key: "vandrare", Glyph: "●", Slot: &slot, Name: p.Name}
}

// placeSign är ett tecken för en plats i tornet: kammare, port eller
// port_stangd.
func placeSign(key string) mapSign {
	glyph := map[string]string{"kammare": "¤", "port": "˄", "port_stangd": "█"}[key]
	return mapSign{Key: key, Glyph: glyph}
}

// signs gör en lista av tecken, så att anropen blir korta.
func signs(s ...mapSign) []mapSign { return s }
