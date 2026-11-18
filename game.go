package main

import (
	"encoding/json"
	"math/rand/v2"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	gridSize     = 9
	center       = gridSize / 2
	maxPlayers   = 4
	maxNameLen   = 16
	moveCooldown = 120 * time.Millisecond
	doorChance   = 0.55
)

type Dir int

const (
	North Dir = iota
	East
	South
	West
)

var (
	dirDX    = [4]int{0, 1, 0, -1}
	dirDY    = [4]int{-1, 0, 1, 0}
	dirCodes = map[string]Dir{"N": North, "E": East, "S": South, "W": West}
)

func (d Dir) Opp() Dir { return (d + 2) % 4 }

// Varje spelarplats hör till ett hörn. Hörnrummen har två fasta dörrar inåt
// och en utgång genom ytterväggen.
type corner struct {
	x, y   int
	inward [2]Dir
	exit   Dir
}

var corners = [maxPlayers]corner{
	{0, 0, [2]Dir{East, South}, West},
	{gridSize - 1, 0, [2]Dir{West, South}, East},
	{0, gridSize - 1, [2]Dir{East, North}, West},
	{gridSize - 1, gridSize - 1, [2]Dir{West, North}, East},
}

type Room struct {
	Revealed bool
	Doors    [4]bool
	Exit     Dir
	HasExit  bool
	Kind     string // "rum", "ingang" eller "kammare"
	TypeID   string
}

type Player struct {
	ID       int
	Name     string
	Slot     int
	X, Y     int
	client   *Client
	lastMove time.Time
}

type Client struct {
	conn   *wsConn
	send   chan []byte
	player *Player
	closed bool
}

type Game struct {
	mu      sync.Mutex
	data    *Data
	rooms   [gridSize][gridSize]Room
	players [maxPlayers]*Player
	nextID  int
	level   int
}

func NewGame(d *Data) *Game {
	return &Game{data: d, level: 1}
}

func inBounds(x, y int) bool { return x >= 0 && y >= 0 && x < gridSize && y < gridSize }

func cornerIndex(x, y int) int {
	for i, c := range corners {
		if c.x == x && c.y == y {
			return i
		}
	}
	return -1
}

// fixedDoors ger dörrarna för rum vars form är bestämd från början.
func fixedDoors(x, y int) ([4]bool, bool) {
	var d [4]bool
	if x == center && y == center {
		return [4]bool{true, true, true, true}, true
	}
	if ci := cornerIndex(x, y); ci >= 0 {
		for _, dir := range corners[ci].inward {
			d[dir] = true
		}
		return d, true
	}
	return d, false
}

// reveal skapar rummet om det är okänt. Dörrarna stämmer alltid med
// grannar som redan finns, så att väggarna hänger ihop.
func (g *Game) reveal(x, y int) *Room {
	r := &g.rooms[y][x]
	if r.Revealed {
		return r
	}
	r.Revealed = true
	if fd, ok := fixedDoors(x, y); ok {
		r.Doors = fd
		if ci := cornerIndex(x, y); ci >= 0 {
			r.Kind, r.TypeID = "ingang", "ingang"
			r.HasExit, r.Exit = true, corners[ci].exit
		} else {
			r.Kind, r.TypeID = "kammare", "skattkammare"
		}
		return r
	}
	r.Kind = "rum"
	r.TypeID = g.data.Rumstyper[rand.IntN(len(g.data.Rumstyper))].ID
	var undecided []Dir
	open := 0
	for d := North; d <= West; d++ {
		nx, ny := x+dirDX[d], y+dirDY[d]
		if !inBounds(nx, ny) {
			continue
		}
		n := &g.rooms[ny][nx]
		if n.Revealed {
			r.Doors[d] = n.Doors[d.Opp()]
		} else if fd, ok := fixedDoors(nx, ny); ok {
			r.Doors[d] = fd[d.Opp()]
		} else {
			undecided = append(undecided, d)
			r.Doors[d] = rand.Float64() < doorChance
		}
		if r.Doors[d] {
			open++
		}
	}
	// Minst två öppningar om det går, så att gångarna sällan blir återvändsgränder.
	for open < 2 && len(undecided) > 0 {
		i := rand.IntN(len(undecided))
		if d := undecided[i]; !r.Doors[d] {
			r.Doors[d] = true
			open++
		}
		undecided = append(undecided[:i], undecided[i+1:]...)
	}
	return r
}

func (g *Game) roomName(r *Room) string {
	return g.data.roomType(r.TypeID).Namn
}

func (g *Game) describe(r *Room) string {
	return g.data.roomType(r.TypeID).Describe(g.level)
}

// Meddelanden till klienten.

type playerView struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Slot   int    `json:"slot"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Online bool   `json:"online"`
}

type roomView struct {
	X    int     `json:"x"`
	Y    int     `json:"y"`
	D    [4]bool `json:"d"`
	K    string  `json:"k"`
	Exit int     `json:"exit"`
	Name string  `json:"name"`
}

type stateMsg struct {
	T       string       `json:"t"`
	You     int          `json:"you"`
	Size    int          `json:"size"`
	Level   int          `json:"level"`
	Players []playerView `json:"players"`
	Rooms   []roomView   `json:"rooms"`
}

type logMsg struct {
	T    string `json:"t"`
	Text string `json:"text"`
	C    string `json:"c,omitempty"`
}

type soundMsg struct {
	T string `json:"t"`
	S string `json:"s"`
}

func (g *Game) send(c *Client, v any) {
	if c == nil || c.closed {
		return
	}
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	select {
	case c.send <- b:
	default:
		// Klienten hänger inte med. Stäng så att den kan ansluta på nytt.
		c.conn.Close()
	}
}

func (g *Game) logTo(p *Player, text, class string) {
	if p != nil {
		g.send(p.client, logMsg{T: "log", Text: text, C: class})
	}
}

func (g *Game) logOthers(except *Player, text, class string) {
	for _, p := range g.players {
		if p != nil && p != except {
			g.logTo(p, text, class)
		}
	}
}

func (g *Game) soundTo(p *Player, s string) {
	if p != nil {
		g.send(p.client, soundMsg{T: "sound", S: s})
	}
}

func (g *Game) soundAll(s string) {
	for _, p := range g.players {
		g.soundTo(p, s)
	}
}

// stateFor bygger det en viss spelare får se.
func (g *Game) stateFor(me *Player) stateMsg {
	st := stateMsg{T: "state", You: me.ID, Size: gridSize, Level: g.level}
	for _, p := range g.players {
		if p == nil {
			continue
		}
		st.Players = append(st.Players, playerView{
			ID: p.ID, Name: p.Name, Slot: p.Slot, X: p.X, Y: p.Y, Online: p.client != nil,
		})
	}
	for y := range gridSize {
		for x := range gridSize {
			r := &g.rooms[y][x]
			if !r.Revealed {
				continue
			}
			rv := roomView{X: x, Y: y, D: r.Doors, K: r.Kind, Exit: -1, Name: g.roomName(r)}
			if r.HasExit {
				rv.Exit = int(r.Exit)
			}
			st.Rooms = append(st.Rooms, rv)
		}
	}
	return st
}

func (g *Game) broadcastState() {
	for _, p := range g.players {
		if p != nil && p.client != nil {
			g.send(p.client, g.stateFor(p))
		}
	}
}

func cleanName(s string) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || r == '<' || r == '>' {
			return -1
		}
		return r
	}, s)
	s = strings.Join(strings.Fields(s), " ")
	for utf8.RuneCountInString(s) > maxNameLen {
		_, size := utf8.DecodeLastRuneInString(s)
		s = s[:len(s)-size]
	}
	return strings.TrimSpace(s)
}

func (g *Game) Join(c *Client, rawName string) {
	if c.player != nil {
		return
	}
	name := cleanName(rawName)
	if name == "" {
		g.send(c, logMsg{T: "error", Text: g.data.Text("tomt_namn")})
		return
	}
	// Återanslutning: samma namn tar tillbaka sin vandrare.
	for _, p := range g.players {
		if p != nil && strings.EqualFold(p.Name, name) {
			if p.client != nil {
				g.send(c, logMsg{T: "error", Text: g.data.Text("namn_upptaget")})
				return
			}
			p.client, c.player = c, p
			g.broadcastState()
			g.logTo(p, g.data.Text("tillbaka"), "syn")
			g.logOthers(p, g.data.Text("atervander", "namn", p.Name), "info")
			return
		}
	}
	slot := -1
	for i, p := range g.players {
		if p == nil {
			slot = i
			break
		}
	}
	if slot < 0 {
		// Alla platser tagna: en vandrare som lämnat spelet får ge plats.
		for i, p := range g.players {
			if p.client == nil {
				slot = i
				break
			}
		}
	}
	if slot < 0 {
		g.send(c, logMsg{T: "error", Text: g.data.Text("fullt")})
		return
	}
	g.nextID++
	cn := corners[slot]
	p := &Player{ID: g.nextID, Name: name, Slot: slot, X: cn.x, Y: cn.y, client: c}
	g.players[slot] = p
	c.player = p
	room := g.reveal(cn.x, cn.y)
	horn := g.data.Horn[slot]

	g.broadcastState()
	g.soundAll("ansluter")
	g.logTo(p, g.data.Text("valkommen", "namn", name, "horn", horn), "syn")
	g.logTo(p, g.describe(room), "")
	g.logOthers(p, g.data.Text("ansluter", "namn", name, "horn", horn), "info")
}

func (g *Game) Move(c *Client, code string) {
	p := c.player
	if p == nil {
		return
	}
	d, ok := dirCodes[code]
	if !ok {
		return
	}
	now := time.Now()
	if now.Sub(p.lastMove) < moveCooldown {
		return
	}
	p.lastMove = now
	room := &g.rooms[p.Y][p.X]
	if room.HasExit && room.Exit == d {
		g.logTo(p, g.data.Text("utgang"), "info")
		return
	}
	if !room.Doors[d] {
		g.soundTo(p, "vagg")
		g.logTo(p, g.data.Text("vagg", "riktning", g.data.Riktningar[d]), "dim")
		return
	}
	nx, ny := p.X+dirDX[d], p.Y+dirDY[d]
	target := &g.rooms[ny][nx]
	isNew := !target.Revealed
	g.reveal(nx, ny)
	p.X, p.Y = nx, ny
	g.broadcastState()

	name := g.roomName(target)
	if isNew {
		g.soundAll("nyttrum")
		g.logTo(p, name+".", "rubrik")
		g.logTo(p, g.describe(target), "")
		g.logOthers(p, g.data.Text("nytt_rum_andra",
			"namn", p.Name, "region", g.data.region(nx, ny)), "info")
	} else {
		g.soundTo(p, "steg")
		g.logTo(p, g.data.Text("kant_rum", "rum", name), "dim")
	}
}

func (g *Game) Leave(c *Client) {
	c.closed = true
	close(c.send)
	p := c.player
	if p == nil || p.client != c {
		return
	}
	p.client = nil
	g.logOthers(p, g.data.Text("lamnar", "namn", p.Name), "dim")
	g.broadcastState()
}

type clientMsg struct {
	T    string `json:"t"`
	Name string `json:"name"`
	Dir  string `json:"dir"`
}

// serve sköter en websocket-anslutning tills den stängs.
func (g *Game) serve(conn *wsConn) {
	c := &Client{conn: conn, send: make(chan []byte, 64)}
	go func() {
		for b := range c.send {
			if err := conn.WriteText(b); err != nil {
				conn.Close()
				for range c.send {
				}
				return
			}
		}
	}()
	for {
		raw, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var m clientMsg
		if json.Unmarshal(raw, &m) != nil {
			continue
		}
		g.mu.Lock()
		switch m.T {
		case "join":
			g.Join(c, m.Name)
		case "move":
			g.Move(c, m.Dir)
		}
		g.mu.Unlock()
	}
	g.mu.Lock()
	g.Leave(c)
	g.mu.Unlock()
	conn.Close()
}
