package main

// Reprisen: vandringen spelad igen på slutskärmen (docs/repris.md). Tornet
// skriver ner var alla står när vandringen börjar och efter varje runda,
// och när vandringen är slut får alla tornet som det blev, med allas vägar
// runda för runda. Klienten ritar den platt och vriden som minikartan
// (web/repris.js). Inget i reprisen säger något om ormens namn.

// replayStep är var en vandrare var efter en runda: rummet, hur det gick
// (inne, ute, dod, forlorad), guldet och om hen låg i valven.
type replayStep struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Status string `json:"s"`
	Gold   int    `json:"g"`
	Cellar bool   `json:"v,omitempty"`
}

// replayPerson är en vandrare i reprisen, i samma ordning som stegen.
type replayPerson struct {
	Name string `json:"namn"`
	Slot int    `json:"slot"`
}

// replayRoom är ett rum som någon har sett, med dörrarna och porten.
type replayRoom struct {
	X    int     `json:"x"`
	Y    int     `json:"y"`
	D    [4]bool `json:"d"`
	K    string  `json:"k"`
	Exit int     `json:"exit"`
	Fell bool    `json:"dod,omitempty"`
}

// replayView är reprisen som skickas i slutet. Rundor[0] är var alla stod
// när vandringen började, och Rundor[n] var de stod efter runda n; den
// sista är efter att solen gick ner eller den siste kom ut.
type replayView struct {
	Size   int            `json:"storlek"`
	Rooms  []replayRoom   `json:"rum"`
	People []replayPerson `json:"vandrare"`
	Frames [][]replayStep `json:"rundor"`
}

// replayStart börjar en ny repris med dem som är med i vandringen.
func (g *Game) replayStart() {
	g.replayWho = nil
	g.replayFrames = nil
	g.replayDone = nil
	for _, p := range g.players {
		if p != nil && p.Status == statusIn {
			g.replayWho = append(g.replayWho, p)
		}
	}
	g.replayNote()
}

// replayNote skriver ner var alla står just nu.
func (g *Game) replayNote() {
	if len(g.replayWho) == 0 {
		return
	}
	frame := make([]replayStep, len(g.replayWho))
	for i, p := range g.replayWho {
		frame[i] = replayStep{X: p.X, Y: p.Y, Status: p.Status, Gold: p.Treasure, Cellar: p.Cellar > 0}
	}
	g.replayFrames = append(g.replayFrames, frame)
}

// replayEnd skriver ner hur det slutade och bygger reprisen som alla får.
func (g *Game) replayEnd() {
	if len(g.replayWho) == 0 {
		g.replayDone = nil
		return
	}
	g.replayNote()
	v := &replayView{Size: g.size, Frames: g.replayFrames}
	for _, p := range g.replayWho {
		v.People = append(v.People, replayPerson{Name: p.Name, Slot: p.Slot})
	}
	for y := range g.size {
		for x := range g.size {
			r := &g.rooms[y][x]
			if !r.Revealed {
				continue
			}
			rr := replayRoom{X: x, Y: y, D: r.Doors, K: r.Kind, Exit: -1, Fell: len(r.fell) > 0}
			if r.HasExit {
				rr.Exit = int(r.Exit)
			}
			v.Rooms = append(v.Rooms, rr)
		}
	}
	g.replayDone = v
}
