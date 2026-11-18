package main

// Svårighet: förval som ändrar några av reglerna i datafilerna. Förvalen
// står i data/svarighet.json. Ett värde som inte står med i ett förval
// tas från de vanliga datafilerna.

import "math"

type Difficulty struct {
	Standard string           `json:"standard"`
	Nivaer   []DifficultyStep `json:"nivaer"`
}

type DifficultyStep struct {
	ID           string   `json:"id"`
	Namn         string   `json:"namn"`
	Beskrivning  string   `json:"beskrivning"`
	Liv          *int     `json:"liv"`           // liv för varje vandrare
	Fallchans    *float64 `json:"fallchans"`     // chans att ett nytt rum har en fälla (regler.handelsechans)
	Varelsechans *float64 `json:"varelsechans"`  // chans att ett nytt rum har en varelse (monster.chans)
	Sovkort      *int     `json:"sovkort"`       // antal kort i leken där ormen inte vaknar
	Vaknakort    *int     `json:"vaknakort"`     // antal kort där ormen vaknar
	Ormskada     *int     `json:"ormskada"`      // skada när ormen vaknar
	SkattForlust *float64 `json:"skatt_forlust"` // andel av skatten som går förlorad när ormen vaknar
}

// difficultyView är det värdsidan behöver för att visa ett förval.
type difficultyView struct {
	ID          string `json:"id"`
	Namn        string `json:"namn"`
	Beskrivning string `json:"beskrivning"`
}

func (d *Data) difficulties() []difficultyView {
	out := []difficultyView{}
	for _, s := range d.Svarighet.Nivaer {
		out = append(out, difficultyView{ID: s.ID, Namn: s.Namn, Beskrivning: s.Beskrivning})
	}
	return out
}

// difficulty ger förvalet med id, eller standardförvalet om id inte finns.
func (d *Data) difficulty(id string) *DifficultyStep {
	for _, want := range []string{id, d.Svarighet.Standard} {
		for i := range d.Svarighet.Nivaer {
			if d.Svarighet.Nivaer[i].ID == want {
				return &d.Svarighet.Nivaer[i]
			}
		}
	}
	return nil
}

// difficultyID ger id för förvalet som gäller för id: id självt om det
// finns, annars standardförvalet, och tomt om det inte finns några förval.
func (d *Data) difficultyID(id string) string {
	if s := d.difficulty(id); s != nil {
		return s.ID
	}
	return ""
}

// withDifficulty ger en kopia av datan med förvalet id. Originalet ändras inte.
func (d *Data) withDifficulty(id string) *Data {
	s := d.difficulty(id)
	if s == nil || (s.Liv == nil && s.Fallchans == nil && s.Varelsechans == nil &&
		s.Sovkort == nil && s.Vaknakort == nil && s.Ormskada == nil && s.SkattForlust == nil) {
		return d
	}
	c := *d
	if s.Liv != nil {
		c.Regler.Liv = max(1, *s.Liv)
	}
	if s.Fallchans != nil {
		c.Regler.Handelsechans = clamp01(*s.Fallchans)
	}
	if s.Varelsechans != nil {
		c.Monster.Chans = clamp01(*s.Varelsechans)
	}
	if s.Ormskada != nil {
		c.Ormen.Skada = max(0, *s.Ormskada)
	}
	if s.SkattForlust != nil {
		c.Ormen.SkattForlust = clamp01(*s.SkattForlust)
	}
	if s.Sovkort != nil || s.Vaknakort != nil {
		c.Ormen.Kort = append([]Card(nil), d.Ormen.Kort...)
		if s.Sovkort != nil {
			spreadCards(c.Ormen.Kort, false, max(1, *s.Sovkort))
		}
		if s.Vaknakort != nil {
			spreadCards(c.Ormen.Kort, true, max(1, *s.Vaknakort))
		}
	}
	return &c
}

// spreadCards sätter antalet kort av ena slaget (vaknar eller inte) till
// total, fördelat på korten i samma proportioner som i datafilen.
func spreadCards(cards []Card, wakes bool, total int) {
	before, largest := 0, -1
	for i, c := range cards {
		if c.Vaknar != wakes {
			continue
		}
		before += c.Antal
		if largest < 0 || c.Antal > cards[largest].Antal {
			largest = i
		}
	}
	if largest < 0 {
		return
	}
	sum := 0
	for i := range cards {
		if cards[i].Vaknar != wakes {
			continue
		}
		if before > 0 {
			cards[i].Antal = int(math.Round(float64(cards[i].Antal) * float64(total) / float64(before)))
		} else {
			cards[i].Antal = 0
		}
		sum += cards[i].Antal
	}
	// Avrundningen kan ge ett kort för mycket eller för lite.
	cards[largest].Antal = max(0, cards[largest].Antal+total-sum)
}

func clamp01(f float64) float64 {
	return math.Max(0, math.Min(1, f))
}
