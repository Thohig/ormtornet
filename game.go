package main

import (
	"encoding/json"
	"math/rand/v2"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	gridSize   = 9
	center     = gridSize / 2
	maxPlayers = 4
	maxNameLen = 16
	doorChance = 0.55
)

// Spelets faser.
const (
	phaseLobby = "lobby"
	phasePlay  = "spel"
	phaseEnd   = "slut"
)

// En spelares läge.
const (
	statusIn   = "inne"
	statusOut  = "ute"
	statusLost = "forlorad"
)

// Drag en spelare kan välja under en runda.
const (
	actWait = "vanta"
	actTake = "ta"
)

type Config struct {
	SunMax     int
	TurnTime   time.Duration
	MinPlayers int
}

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
	Ready    bool
	Action   string
	Treasure int
	Status   string
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
	cfg     Config
	rooms   [gridSize][gridSize]Room
	players [maxPlayers]*Player
	nextID  int
	level   int

	phase    string
	sun      int
	round    int
	deadline time.Time
	timer    *time.Timer
}

func NewGame(d *Data, cfg Config) *Game {
	return &Game{data: d, cfg: cfg, level: 1, phase: phaseLobby, sun: cfg.SunMax}
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
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Slot     int    `json:"slot"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Online   bool   `json:"online"`
	Status   string `json:"status"`
	Ready    bool   `json:"ready"`
	Chosen   bool   `json:"chosen"`
	Treasure int    `json:"treasure"` // -1 om spelaren inte får veta
}

type roomView struct {
	X    int     `json:"x"`
	Y    int     `json:"y"`
	D    [4]bool `json:"d"`
	K    string  `json:"k"`
	Exit int     `json:"exit"`
	Name string  `json:"name"`
}

type resultView struct {
	Name     string `json:"name"`
	Slot     int    `json:"slot"`
	Treasure int    `json:"treasure"`
	Status   string `json:"status"`
	Winner   bool   `json:"winner"`
}

type stateMsg struct {
	T          string       `json:"t"`
	You        int          `json:"you"`
	Size       int          `json:"size"`
	Level      int          `json:"level"`
	Phase      string       `json:"phase"`
	Sun        int          `json:"sun"`
	SunMax     int          `json:"sunMax"`
	Round      int          `json:"round"`
	DeadlineMs int64        `json:"deadlineMs"`
	MinPlayers int          `json:"minPlayers"`
	Action     string       `json:"action"`
	Players    []playerView `json:"players"`
	Rooms      []roomView   `json:"rooms"`
	Results    []resultView `json:"results,omitempty"`
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
	if p != nil && text != "" {
		g.send(p.client, logMsg{T: "log", Text: text, C: class})
	}
}

func (g *Game) logAll(text, class string) {
	g.logOthers(nil, text, class)
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

// stateFor bygger det en viss spelare får se. Andras skatt och val av drag
// hålls hemliga tills spelet är slut.
func (g *Game) stateFor(me *Player) stateMsg {
	st := stateMsg{
		T: "state", You: me.ID, Size: gridSize, Level: g.level,
		Phase: g.phase, Sun: g.sun, SunMax: g.cfg.SunMax, Round: g.round,
		MinPlayers: g.cfg.MinPlayers, Action: me.Action,
	}
	if g.phase == phasePlay {
		st.DeadlineMs = max(0, time.Until(g.deadline).Milliseconds())
	}
	for _, p := range g.players {
		if p == nil {
			continue
		}
		pv := playerView{
			ID: p.ID, Name: p.Name, Slot: p.Slot, X: p.X, Y: p.Y,
			Online: p.client != nil, Status: p.Status, Ready: p.Ready,
			Chosen: p.Action != "", Treasure: -1,
		}
		if p == me || g.phase == phaseEnd {
			pv.Treasure = p.Treasure
		}
		st.Players = append(st.Players, pv)
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
	if g.phase == phaseEnd {
		st.Results = g.results()
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

// results sorterar spelarna: de som kom ut med skatt först, mest skatt överst.
func (g *Game) results() []resultView {
	var res []resultView
	for _, p := range g.players {
		if p == nil {
			continue
		}
		res = append(res, resultView{Name: p.Name, Slot: p.Slot, Treasure: p.Treasure, Status: p.Status})
	}
	score := func(r resultView) int {
		if r.Status == statusOut && r.Treasure > 0 {
			return r.Treasure
		}
		return -1
	}
	sort.SliceStable(res, func(i, j int) bool { return score(res[i]) > score(res[j]) })
	if len(res) > 0 && score(res[0]) > 0 {
		best := score(res[0])
		for i := range res {
			res[i].Winner = score(res[i]) == best
		}
	}
	return res
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

func (g *Game) placeAtCorner(p *Player) {
	cn := corners[p.Slot]
	p.X, p.Y = cn.x, cn.y
	p.Status = statusIn
	p.Treasure = 0
	p.Action = ""
	p.Ready = false
	g.reveal(cn.x, cn.y)
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
	if g.phase == phasePlay {
		g.send(c, logMsg{T: "error", Text: g.data.Text("spel_pagar")})
		return
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
	p := &Player{ID: g.nextID, Name: name, Slot: slot, client: c}
	g.players[slot] = p
	c.player = p
	g.placeAtCorner(p)
	horn := g.data.Horn[slot]

	g.broadcastState()
	g.soundAll("ansluter")
	g.logTo(p, g.data.Text("valkommen", "namn", name, "horn", horn), "syn")
	g.logTo(p, g.describe(&g.rooms[p.Y][p.X]), "")
	g.logTo(p, g.data.Text("lobby", "min", strconv.Itoa(g.cfg.MinPlayers)), "info")
	g.logOthers(p, g.data.Text("ansluter", "namn", name, "horn", horn), "info")
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
	if g.phase != phasePlay {
		// Utanför spelet lämnar man sin plats helt.
		g.players[p.Slot] = nil
		g.checkStart()
	} else if g.allChosen() {
		g.resolve()
		return
	}
	g.broadcastState()
}

// Ready markerar att spelaren vill börja (eller spela igen).
func (g *Game) Ready(c *Client) {
	p := c.player
	if p == nil || g.phase == phasePlay || p.Ready {
		return
	}
	p.Ready = true
	g.logOthers(p, g.data.Text("redo_andra", "namn", p.Name), "dim")
	g.soundTo(p, "val")
	if !g.checkStart() {
		g.logTo(p, g.data.Text("vantar_pa_andra"), "dim")
		g.broadcastState()
	}
}

// checkStart startar spelet när alla anslutna är redo.
func (g *Game) checkStart() bool {
	n := 0
	for _, p := range g.players {
		if p == nil || p.client == nil {
			continue
		}
		if !p.Ready {
			return false
		}
		n++
	}
	if n < g.cfg.MinPlayers {
		return false
	}
	if g.phase == phaseEnd {
		g.newTower()
	}
	g.startGame()
	return true
}

// newTower bygger ett nytt torn och ställer alla vid sina portar.
func (g *Game) newTower() {
	g.rooms = [gridSize][gridSize]Room{}
	for _, p := range g.players {
		if p != nil {
			g.placeAtCorner(p)
		}
	}
}

func (g *Game) startGame() {
	g.phase = phasePlay
	g.sun = g.cfg.SunMax
	g.round = 1
	for _, p := range g.players {
		if p != nil {
			p.Ready = false
			p.Action = ""
		}
	}
	g.soundAll("start")
	g.logAll(g.data.Text("start", "sol", strconv.Itoa(g.sun)), "syn")
	g.startRound()
}

func (g *Game) startRound() {
	g.deadline = time.Now().Add(g.cfg.TurnTime)
	round := g.round
	if g.timer != nil {
		g.timer.Stop()
	}
	g.timer = time.AfterFunc(g.cfg.TurnTime, func() {
		g.mu.Lock()
		defer g.mu.Unlock()
		if g.phase == phasePlay && g.round == round {
			g.resolve()
		}
	})
	g.broadcastState()
}

// allChosen är sant när alla anslutna spelare i tornet har valt drag.
func (g *Game) allChosen() bool {
	any := false
	for _, p := range g.players {
		if p == nil || p.Status != statusIn || p.client == nil {
			continue
		}
		if p.Action == "" {
			return false
		}
		any = true
	}
	return any
}

// Act tar emot spelarens val för rundan. Valet kan ändras tills rundan utförs.
func (g *Game) Act(c *Client, a string) {
	p := c.player
	if p == nil || g.phase != phasePlay || p.Status != statusIn {
		return
	}
	switch a {
	case "N", "E", "S", "W", actWait:
	case actTake:
		if g.rooms[p.Y][p.X].Kind != "kammare" {
			g.logTo(p, g.data.Text("ingen_skatt_har"), "dim")
			return
		}
	default:
		return
	}
	p.Action = a
	g.soundTo(p, "val")
	if g.allChosen() {
		g.resolve()
		return
	}
	g.broadcastState()
}

// resolve utför alla drag på en gång och låter solen sjunka.
func (g *Game) resolve() {
	if g.timer != nil {
		g.timer.Stop()
	}
	var active []*Player
	for _, p := range g.players {
		if p != nil && p.Status == statusIn {
			active = append(active, p)
		}
	}
	rand.Shuffle(len(active), func(i, j int) { active[i], active[j] = active[j], active[i] })
	for _, p := range active {
		a := p.Action
		if a == "" {
			a = actWait
			g.logTo(p, g.data.Text("tvekade"), "dim")
		}
		switch a {
		case actWait:
		case actTake:
			g.take(p)
		default:
			g.move(p, dirCodes[a])
		}
		p.Action = ""
	}

	g.sun--
	g.soundAll("runda")
	inside := 0
	for _, p := range g.players {
		if p != nil && p.Status == statusIn {
			inside++
		}
	}
	switch {
	case inside == 0:
		g.endGame()
	case g.sun <= 0:
		g.logAll(g.data.Text("solnedgang"), "varn")
		g.soundAll("solnedgang")
		g.endGame()
	default:
		g.sunWarning()
		g.round++
		g.startRound()
	}
}

func (g *Game) sunWarning() {
	key := ""
	switch {
	case g.sun == g.cfg.SunMax/2:
		key = "sol_halv"
	case g.sun == g.cfg.SunMax/4:
		key = "sol_kvart"
	case g.sun == 3:
		key = "sol_snart"
	}
	if key != "" {
		g.logAll(g.data.Text(key, "sol", strconv.Itoa(g.sun)), "varn")
		g.soundAll("solvarning")
	}
}

func (g *Game) endGame() {
	g.phase = phaseEnd
	if g.timer != nil {
		g.timer.Stop()
	}
	for _, p := range g.players {
		if p == nil {
			continue
		}
		p.Ready = false
		p.Action = ""
		if p.Status == statusIn {
			p.Status = statusLost
			g.logTo(p, g.data.Text("forlorad_du"), "varn")
			g.logOthers(p, g.data.Text("forlorad", "namn", p.Name), "varn")
		}
	}
	res := g.results()
	if len(res) > 0 && res[0].Winner {
		var names []string
		for _, r := range res {
			if r.Winner {
				names = append(names, r.Name)
			}
		}
		g.logAll(g.data.Text("vinnare", "namn", strings.Join(names, " och "),
			"skatt", strconv.Itoa(res[0].Treasure)), "syn")
	} else {
		g.logAll(g.data.Text("ingen_vinnare"), "varn")
	}
	g.logAll(g.data.Text("spela_igen"), "dim")
	g.soundAll("slut")
	g.broadcastState()
}

// move flyttar spelaren ett rum, eller ut genom porten.
func (g *Game) move(p *Player, d Dir) {
	room := &g.rooms[p.Y][p.X]
	if room.HasExit && room.Exit == d {
		p.Status = statusOut
		if p.Treasure > 0 {
			g.logTo(p, g.data.Text("ut_med_skatt", "skatt", strconv.Itoa(p.Treasure)), "syn")
		} else {
			g.logTo(p, g.data.Text("ut_utan_skatt"), "info")
		}
		g.soundTo(p, "ut")
		g.logOthers(p, g.data.Text("andra_ut", "namn", p.Name), "info")
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

// take ger en näve skatt i kammaren. (Ormen bryr sig inte ännu.)
func (g *Game) take(p *Player) {
	if g.rooms[p.Y][p.X].Kind != "kammare" {
		return
	}
	n := 10 + rand.IntN(41)
	p.Treasure += n
	g.soundTo(p, "skatt")
	g.logTo(p, g.data.Text("ta_skatt", "antal", strconv.Itoa(n), "skatt", strconv.Itoa(p.Treasure)), "syn")
}

type clientMsg struct {
	T    string `json:"t"`
	Name string `json:"name"`
	A    string `json:"a"`
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
		case "act":
			g.Act(c, m.A)
		case "ready":
			g.Ready(c)
		}
		g.mu.Unlock()
	}
	g.mu.Lock()
	g.Leave(c)
	g.mu.Unlock()
	conn.Close()
}
