package main

// Spelledarläget: den som leder tornet spelar inte själv, men ser allt som
// finns i tornet och kan starta en vandring även när alla vandrare är
// AI-spelare. Ingen kan ansluta som vandrare från spelledarens dator.

import "time"

// gmEntry är en rad i spelledarens logg: en text och vilka som fick den.
type gmEntry struct {
	Round int      `json:"round"`
	Text  string   `json:"text"`
	Class string   `json:"c,omitempty"`
	To    []string `json:"to"`
}

const gmLogMax = 400

// recordLog sparar det en vandrare får läsa, så att spelledaren kan se allt.
// Samma text till flera i följd blir en rad med alla mottagarna.
func (g *Game) recordLog(p *Player, text, class string) {
	if n := len(g.gmLog); n > 0 {
		last := &g.gmLog[n-1]
		if last.Text == text && last.Class == class && last.Round == g.round {
			for _, name := range last.To {
				if name == p.Name {
					return
				}
			}
			last.To = append(last.To, p.Name)
			return
		}
	}
	g.gmLog = append(g.gmLog, gmEntry{Round: g.round, Text: text, Class: class, To: []string{p.Name}})
	if len(g.gmLog) > gmLogMax {
		g.gmLog = append([]gmEntry(nil), g.gmLog[len(g.gmLog)-gmLogMax:]...)
	}
}

// gmPlayer är allt spelledaren ser om en vandrare.
type gmPlayer struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Slot     int    `json:"slot"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Status   string `json:"status"`
	Online   bool   `json:"online"`
	Bot      bool   `json:"bot"`
	AI       bool   `json:"ai"`
	HP       int    `json:"hp"`
	Treasure int    `json:"treasure"`
	Potions  int    `json:"potions"`
	Stuck    int    `json:"stuck"`
	Cellar   int    `json:"cellar"`
	Poison   int    `json:"poison"`
	Dark     int    `json:"dark"`
	Action   string `json:"action"` // valt drag, med motståndare eller varelse
	Room     string `json:"room"`
}

type gmSerpent struct {
	Sleep      int `json:"sleep"` // 0 djupt, 1 oroligt, 2 nära att vakna
	SleepCards int `json:"sleepCards"`
	WakeCards  int `json:"wakeCards"`
	Wakes      int `json:"wakes"`
	Risk       int `json:"risk"` // chansen i procent att ormen vaknar vid nästa näve
}

// gmView är spelledarens översikt.
type gmView struct {
	Map     stateMsg   `json:"map"` // samma form som spelarnas karta, men allt syns
	MaxHP   int        `json:"maxHp"`
	Players []gmPlayer `json:"players"`
	Serpent gmSerpent  `json:"serpent"`
	Level   int        `json:"level"`
	Held    bool       `json:"held"`     // spelledaren håller kvar nivån
	SunLvl  int        `json:"sunLevel"` // nivån solen ger just nu
	Clues   int        `json:"clues"`    // ledtrådar som har delats ut
	Traps   int        `json:"traps"`    // fällor som väntar i kända rum
	Unknown int        `json:"unknown"`  // rum som ännu inte finns
	Log     []gmEntry  `json:"log"`

	Paused       bool        `json:"paused"`
	TimeLeft     int         `json:"timeLeft"` // sekunder kvar av rundan
	Rooms        []gmRoomRef `json:"rooms"`    // rum där spelledaren kan lägga fällor och varelser
	TrapTypes    []gmCatalog `json:"trapTypes"`
	MonsterTypes []gmCatalog `json:"monsterTypes"`
}

type gmRoomRef struct {
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Name    string `json:"name"`
	Trap    bool   `json:"trap"`
	Monster bool   `json:"monster"`
}

// GMView bygger översikten. Tornet växer fram när någon går in i ett nytt
// rum, så rum som ingen har varit i finns inte än och kan inte visas.
func (g *Game) GMView(visible bool) gmView {
	g.mu.Lock()
	defer g.mu.Unlock()
	if visible {
		g.seen()
	}
	m := stateMsg{
		T: "state", You: 0, Size: gridSize, Level: g.level, Reveal: g.cfg.Reveal,
		Phase: g.phase, Sun: g.sun, SunMax: g.cfg.SunMax, Round: g.round,
		Sleep: -1,
	}
	v := gmView{MaxHP: g.data.Regler.Liv, Level: g.level, Clues: len(g.usedClues), Players: []gmPlayer{}}
	for y := range gridSize {
		for x := range gridSize {
			r := &g.rooms[y][x]
			if !r.Revealed {
				v.Unknown++
				continue
			}
			rv := roomView{X: x, Y: y, D: r.Doors, K: r.Kind, Exit: -1, Name: g.roomName(r), Mark: r.Mark, S: r.Searched}
			if r.HasExit {
				rv.Exit = int(r.Exit)
			}
			if r.Event != "" {
				// En fälla som väntar visas med sitt tecken, för spelledaren.
				v.Traps++
				if t := g.data.trap(r.Event); t != nil {
					rv.Mark = t.Tecken
				}
			}
			m.Rooms = append(m.Rooms, rv)
			if r.Kind == "rum" {
				v.Rooms = append(v.Rooms, gmRoomRef{X: x, Y: y, Name: g.roomName(r), Trap: r.Event != "", Monster: g.monsterAt(x, y) != nil})
			}
		}
	}
	for _, p := range g.players {
		if p == nil {
			continue
		}
		pv := playerView{ID: p.ID, Name: p.Name, Slot: p.Slot, X: p.X, Y: p.Y, Online: p.present(),
			Bot: p.bot != nil, Status: p.Status, Ready: p.Ready, Chosen: p.Action != "", Treasure: p.Treasure, AI: p.takenOver}
		if p.Cellar > 0 {
			pv.Hidden = true
		}
		m.Players = append(m.Players, pv)
		room := ""
		if p.Cellar > 0 {
			room = "Källarvalven"
		} else if inBounds(p.X, p.Y) && g.rooms[p.Y][p.X].Revealed {
			room = g.roomName(&g.rooms[p.Y][p.X])
		}
		v.Players = append(v.Players, gmPlayer{
			ID: p.ID, Name: p.Name, Slot: p.Slot, X: p.X, Y: p.Y, Status: p.Status, Online: p.present(),
			Bot: p.bot != nil && !p.takenOver, AI: p.takenOver, HP: p.HP, Treasure: p.Treasure, Potions: p.Potions,
			Stuck: p.Stuck, Cellar: p.Cellar, Poison: p.Poison, Dark: p.Dark, Action: g.actionText(p), Room: room,
		})
	}
	m.Monsters = g.monsterViews(nil, false)
	v.Map = m
	for _, id := range g.deck {
		if c := g.data.card(id); c != nil && c.Vaknar {
			v.Serpent.WakeCards++
		} else {
			v.Serpent.SleepCards++
		}
	}
	v.Serpent.Risk = g.wakeRisk()
	v.Serpent.Sleep = g.sleepLevel()
	v.Held = g.levelHeld
	v.SunLvl = g.levelFor(g.sun)
	v.Serpent.Wakes = g.wakes
	v.Log = append([]gmEntry(nil), g.gmLog...)
	v.Paused = g.paused
	if g.phase == phasePlay {
		left := time.Until(g.deadline)
		if g.paused {
			left = g.pauseLeft
		}
		v.TimeLeft = int(max(0, left.Round(time.Second)/time.Second))
	}
	v.TrapTypes, v.MonsterTypes = g.gmCatalogs()
	return v
}

// actionText beskriver ett valt drag för spelledaren.
func (g *Game) actionText(p *Player) string {
	a := p.Action
	switch {
	case a == "":
		return ""
	case isCombatAction(a):
		if m := g.monsterAt(p.X, p.Y); m != nil {
			return g.choiceWord(a) + " mot " + lowerFirst(g.data.monsterType(m.Type).Namn)
		}
		if foe := g.playerByID(p.target); foe != nil {
			return g.choiceWord(a) + " mot " + foe.Name
		}
		return g.choiceWord(a)
	case a[0] == 'm':
		return "styr en varelse"
	}
	words := map[string]string{"N": "går norrut", "E": "går österut", "S": "går söderut", "W": "går västerut",
		actWait: "väntar", actTake: "tar skatt", actSearch: "söker", actDrink: "dricker"}
	if w, ok := words[a]; ok {
		return w
	}
	return a
}

// kickLocal tar bort vandrare som spelar från spelledarens egen dator när
// spelledarläget slås på. De får veta varför och försöker inte ansluta igen.
func (g *Game) kickLocal() {
	for i, p := range g.players {
		if p == nil || p.client == nil || !p.client.local {
			continue
		}
		g.send(p.client, logMsg{T: "spelledare", Text: g.data.Text("spelledare_kan_inte")})
		p.client.player = nil
		g.players[i] = nil
	}
}

// gmAnnounce berättar för alla att tornet har, eller inte längre har, en spelledare.
func (g *Game) gmAnnounce(on bool) {
	key := "spelledare_av"
	if on {
		key = "spelledare_pa"
	}
	g.logAll(g.data.Text(key), "info")
}
