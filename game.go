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
	statusDead = "dod"
)

// Drag en spelare kan välja under en runda.
const (
	actWait   = "vanta"
	actTake   = "ta"
	actSearch = "sok"
	actDrink  = "drick"
)

type Config struct {
	AI           int    // önskat antal datorspelare
	AbsentRounds int    // rundor innan AI tar över en frånvarande vandrare, 0 = aldrig
	PvP          string // strid mellan vandrare: "av", "pa" eller "halva"
	GMMode       bool   // spelledarläge: spelledaren spelar inte och ser allt
	Reveal       bool   // false stänger av avslöjandet
	SunMax       int
	TurnTime     time.Duration
	MinPlayers   int
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
	Event    string // fälla som väntar på den som först går in
	Mark     string // tecken som visar att en fälla har slagit till
	Searched bool
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
	Seen     map[string]bool // fällor och varelser spelaren har mött
	bot      *botBrain       // satt för datorstyrda vandrare

	awayRound int  // rundan då människan lämnade mitt i vandringen, 0 om den är här
	takenOver bool // en AI-spelare styr vandraren medan människan är borta
	aiPlayed  bool // AI:n har styrt vandraren någon gång under vandringen

	target int    // vandraren som ett stridsval den här rundan riktas mot
	met    string // vilka andra vandrare som stod i rummet förra gången, så att det bara sägs en gång

	HP      int
	Potions int
	Stuck   int // rundor kvar fast (fallgaller, ormgrop, skinn)
	Cellar  int // rundor kvar i källarvalven
	Poison  int // rundor kvar med gift i kroppen
	Dark    int // rundor kvar utan ljus
}

// present är sant för en spelare som deltar just nu: en människa med öppen
// webbläsare, eller en datorspelare.
func (p *Player) present() bool { return p.client != nil || p.bot != nil }

type Client struct {
	conn   *wsConn
	send   chan []byte
	player *Player
	closed bool
	local  bool // ansluten från datorn som kör tornet
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

	deck      []string // ormkorten som är kvar
	deckSleep int      // antal sovkort när leken byggdes
	wakes     int      // hur många gånger ormen har vaknat

	monsters      []*Monster
	nextMonsterID int
	controller    int          // id för den döda spelare som styr varelserna
	pendingMove   *monsterMove // varelsens drag för rundan

	usedClues map[string]bool // ledtrådar som redan har delats ut

	lastActive time.Time // senast en människa var ansluten eller spelledarens sida frågade

	gmLog []gmEntry // allt vandrarna har fått läsa, för spelledarläget

	paused    bool          // spelledaren har pausat rundan
	pauseLeft time.Duration // tid som var kvar när rundan pausades
}

func NewGame(d *Data, cfg Config) *Game {
	return &Game{data: d, cfg: cfg, level: 1, phase: phaseLobby, sun: cfg.SunMax, usedClues: map[string]bool{}, lastActive: time.Now()}
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
	// Händelse: en varelse, en fälla eller ingenting.
	if !g.spawnMonster(x, y) && rand.Float64() < g.data.Regler.Handelsechans {
		if t := g.data.randomTrap(); t != nil {
			r.Event = t.ID
		}
	}
	return r
}

func (g *Game) roomName(r *Room) string {
	return g.data.roomType(r.TypeID).Namn
}

// describe ger spelaren en beskrivning av rummet för tornets nuvarande nivå.
func (g *Game) describe(p *Player, r *Room, x, y int) string {
	return g.data.roomType(r.TypeID).Describe(g.level, p.Slot+x+y)
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
	Bot      bool   `json:"bot"`
	Hidden   bool   `json:"hidden"` // positionen är okänd för den som tittar
	AI       bool   `json:"ai"`     // en AI-spelare styr medan människan är borta
}

type roomView struct {
	X    int     `json:"x"`
	Y    int     `json:"y"`
	D    [4]bool `json:"d"`
	K    string  `json:"k"`
	Exit int     `json:"exit"`
	Name string  `json:"name"`
	Mark string  `json:"mark,omitempty"`
	S    bool    `json:"s,omitempty"` // genomsökt
}

type resultView struct {
	Name     string `json:"name"`
	Slot     int    `json:"slot"`
	Treasure int    `json:"treasure"`
	Status   string `json:"status"`
	Winner   bool   `json:"winner"`
	AI       bool   `json:"ai"` // AI:n spelade en del av vandringen
}

type stateMsg struct {
	T          string        `json:"t"`
	You        int           `json:"you"`
	Size       int           `json:"size"`
	Level      int           `json:"level"`
	Reveal     bool          `json:"reveal"`
	Phase      string        `json:"phase"`
	Sun        int           `json:"sun"`
	SunMax     int           `json:"sunMax"`
	Round      int           `json:"round"`
	DeadlineMs int64         `json:"deadlineMs"`
	MinPlayers int           `json:"minPlayers"`
	Action     string        `json:"action"`
	HP         int           `json:"hp"`
	MaxHP      int           `json:"maxHp"`
	Potions    int           `json:"potions"`
	Stuck      int           `json:"stuck"`
	Cellar     int           `json:"cellar"`
	Poison     int           `json:"poison"`
	Dark       int           `json:"dark"`
	Sleep      int           `json:"sleep"` // -1 utanför kammaren
	Controller bool          `json:"controller"`
	Monsters   []monsterView `json:"monsters"`
	PvP        string        `json:"pvp"` // om det går att strida mot andra vandrare här
	Paused     bool          `json:"paused,omitempty"`
	SleepText  string        `json:"sleepText,omitempty"`
	Players    []playerView  `json:"players"`
	Rooms      []roomView    `json:"rooms"`
	Results    []resultView  `json:"results,omitempty"`
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
		g.recordLog(p, text, class)
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
		T: "state", You: me.ID, Size: gridSize, Level: g.level, Reveal: g.cfg.Reveal,
		Phase: g.phase, Sun: g.sun, SunMax: g.cfg.SunMax, Round: g.round,
		MinPlayers: g.cfg.MinPlayers, Action: me.Action,
		HP: me.HP, MaxHP: g.data.Regler.Liv, Potions: me.Potions,
		Stuck: me.Stuck, Cellar: me.Cellar, Poison: me.Poison, Dark: me.Dark,
	}
	// I mörker ser man bara sig själv och rummet man står i.
	blind := g.phase == phasePlay && me.Status == statusIn && me.Dark > 0
	st.Sleep = -1
	if g.phase == phasePlay && me.Status == statusIn && me.Cellar == 0 && g.rooms[me.Y][me.X].Kind == "kammare" {
		st.Sleep = g.sleepLevel()
		st.SleepText = g.data.Ormen.Somn[min(st.Sleep, len(g.data.Ormen.Somn)-1)]
	}
	if g.phase == phasePlay {
		st.DeadlineMs = max(0, time.Until(g.deadline).Milliseconds())
		if g.paused {
			st.DeadlineMs = g.pauseLeft.Milliseconds()
			st.Paused = true
		}
	}
	for _, p := range g.players {
		if p == nil {
			continue
		}
		pv := playerView{
			ID: p.ID, Name: p.Name, Slot: p.Slot, X: p.X, Y: p.Y,
			Online: p.present(), Bot: p.bot != nil, Status: p.Status, Ready: p.Ready,
			Chosen: p.Action != "", Treasure: -1, AI: p.takenOver,
		}
		if p == me || g.phase == phaseEnd {
			pv.Treasure = p.Treasure
		}
		if p.Cellar > 0 || (blind && p != me) {
			pv.Hidden = true
			pv.X, pv.Y = -1, -1
		}
		st.Players = append(st.Players, pv)
	}
	for y := range gridSize {
		for x := range gridSize {
			r := &g.rooms[y][x]
			if !r.Revealed || (blind && (x != me.X || y != me.Y)) {
				continue
			}
			rv := roomView{X: x, Y: y, D: r.Doors, K: r.Kind, Exit: -1, Name: g.roomName(r), Mark: r.Mark, S: r.Searched}
			if r.HasExit {
				rv.Exit = int(r.Exit)
			}
			st.Rooms = append(st.Rooms, rv)
		}
	}
	st.Monsters = g.monsterViews(me, blind)
	if g.phase == phasePlay {
		st.PvP = g.pvpStatus(me)
	}
	st.Controller = g.isController(me) && g.phase == phasePlay
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
		res = append(res, resultView{Name: p.Name, Slot: p.Slot, Treasure: p.Treasure, Status: p.Status, AI: p.aiPlayed})
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
	p.Ready = p.bot != nil // datorspelare är alltid redo
	p.HP = g.data.Regler.Liv
	p.Seen = map[string]bool{}
	p.Potions, p.Stuck, p.Cellar, p.Poison, p.Dark = 0, 0, 0, 0, 0
	g.reveal(cn.x, cn.y)
}

func (g *Game) Join(c *Client, rawName string) {
	if c.player != nil {
		return
	}
	if c.local && g.cfg.GMMode {
		g.send(c, logMsg{T: "error", Text: g.data.Text("spelledare_kan_inte")})
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
			if p.client != nil || (p.bot != nil && !p.takenOver) {
				g.send(c, logMsg{T: "error", Text: g.data.Text("namn_upptaget")})
				return
			}
			if p.takenOver {
				g.handBack(p, c)
				return
			}
			p.awayRound = 0
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
			if !p.present() {
				slot = i
				break
			}
		}
	}
	if slot < 0 {
		// Utanför spelet kliver en datorspelare åt sidan för en människa.
		slot = g.removeBot(true)
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
	g.logTo(p, g.describe(p, &g.rooms[p.Y][p.X], p.X, p.Y), "")
	for _, b := range g.players {
		if b != nil && b.bot != nil {
			g.logTo(p, fill(pick(g.data.AI.Ansluter), "namn", b.Name, "horn", g.data.Horn[b.Slot]), "info")
		}
	}
	if g.cfg.GMMode {
		g.logTo(p, g.data.Text("spelledare_pa"), "info")
	}
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
	if g.phase == phasePlay {
		p.awayRound = g.round
	}
	g.logOthers(p, g.data.Text("lamnar", "namn", p.Name), "dim")
	if g.phase != phasePlay {
		// Utanför spelet lämnar man sin plats helt, och en datorspelare
		// kan ta den om spelledaren vill ha fler.
		g.players[p.Slot] = nil
		g.syncBots(true)
		g.checkStart()
	} else if g.allChosen() && !g.paused {
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

// checkStart startar spelet när alla anslutna människor är redo.
// Datorspelare är alltid redo, men minst en människa måste vara med.
func (g *Game) checkStart() bool {
	n, humans := 0, 0
	for _, p := range g.players {
		if p == nil || !p.present() {
			continue
		}
		if p.bot == nil {
			if !p.Ready {
				return false
			}
			humans++
		}
		n++
	}
	if humans == 0 || n < g.cfg.MinPlayers {
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
	g.monsters = nil
	for _, p := range g.players {
		if p != nil {
			g.placeAtCorner(p)
		}
	}
}

func (g *Game) startGame() {
	g.paused = false
	// De som lämnade förra vandringen och inte kom tillbaka är inte med längre.
	if g.dropAbsent() > 0 {
		g.syncBots(true) // AI-spelare kan ta de lediga platserna
	}
	g.phase = phasePlay
	g.sun = g.cfg.SunMax
	g.round = 1
	g.wakes = 0
	g.buildDeck()
	g.controller = 0
	g.pendingMove = nil
	g.level = 1
	g.usedClues = map[string]bool{}
	for _, p := range g.players {
		if p != nil {
			p.Ready = p.bot != nil
			p.Action = ""
			p.awayRound = 0
			p.aiPlayed = false
		}
	}
	g.soundAll("start")
	g.logAll(g.data.Text("start", "sol", strconv.Itoa(g.sun)), "syn")
	g.startRound()
}

func (g *Game) startRound() {
	wait := g.cfg.TurnTime
	if g.onlyCellar() {
		// De anslutna ligger alla i källaren och kan inte välja: gå vidare snabbt.
		wait = min(wait, 2*time.Second)
	}
	g.armTimer(wait)
	g.autoTakeOver()
	g.scheduleBots()
	g.broadcastState()
}

// armTimer sätter rundans tidsgräns. När tiden är slut avgörs rundan, om
// inte spelledaren har pausat.
func (g *Game) armTimer(wait time.Duration) {
	g.deadline = time.Now().Add(wait)
	round := g.round
	if g.timer != nil {
		g.timer.Stop()
	}
	g.timer = time.AfterFunc(wait, func() {
		g.mu.Lock()
		defer g.mu.Unlock()
		if g.phase == phasePlay && g.round == round && !g.paused {
			g.resolve()
		}
	})
}

// onlyCellar är sant om alla anslutna spelare i tornet ligger i källaren.
func (g *Game) onlyCellar() bool {
	cellar := false
	for _, p := range g.players {
		if p == nil || p.Status != statusIn || !p.present() {
			continue
		}
		if p.Cellar == 0 {
			return false
		}
		cellar = true
	}
	return cellar
}

// allChosen är sant när alla anslutna spelare i tornet har valt drag.
func (g *Game) allChosen() bool {
	any := false
	for _, p := range g.players {
		if p == nil || !p.present() {
			continue
		}
		if !g.isController(p) && (p.Status != statusIn || p.Cellar > 0) {
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
func (g *Game) Act(c *Client, a string, monsterID int) {
	g.actPlayer(c.player, a, monsterID)
}

// actPlayer tar emot ett val från en människa eller en datorspelare.
func (g *Game) actPlayer(p *Player, a string, monsterID int) {
	if p != nil && g.phase == phasePlay && g.isController(p) {
		g.actController(p, a, monsterID)
		return
	}
	if p == nil || g.phase != phasePlay || p.Status != statusIn || p.Cellar > 0 {
		return
	}
	room := &g.rooms[p.Y][p.X]
	switch a {
	case actAttack, actParry, actDodge:
		if g.monsterAt(p.X, p.Y) == nil {
			// Ingen varelse: kanske en annan vandrare. Ett negativt id pekar
			// ut vandraren man vill slåss mot.
			switch g.pvpStatus(p) {
			case pvpYes:
				p.target = g.pickFoe(p, -monsterID).ID
			case pvpPeace:
				g.logTo(p, g.data.Text("pvp_fredat"), "dim")
				return
			case pvpTooEarly:
				g.logTo(p, g.data.Text("pvp_for_tidigt"), "dim")
				return
			default:
				g.logTo(p, g.data.MonsterText("inget_monster_har"), "dim")
				return
			}
		}
	case "N", "E", "S", "W", actWait:
	case actTake:
		if room.Kind != "kammare" {
			g.logTo(p, g.data.Text("ingen_skatt_har"), "dim")
			return
		}
	case actSearch:
		if room.Kind != "rum" {
			g.logTo(p, g.data.SearchText("inte_har"), "dim")
			return
		}
		if room.Searched {
			g.logTo(p, g.data.SearchText("redan"), "dim")
			return
		}
	case actDrink:
		if p.Potions <= 0 {
			g.logTo(p, g.data.Text("ingen_dryck"), "dim")
			return
		}
	default:
		return
	}
	p.Action = a
	g.soundTo(p, "val")
	if g.allChosen() && !g.paused {
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
	sunBefore := g.sun
	g.moveMonsters()
	for _, p := range g.players {
		if g.isController(p) {
			p.Action = ""
		}
	}
	var active []*Player
	for _, p := range g.players {
		if p != nil && p.Status == statusIn {
			active = append(active, p)
		}
	}
	rand.Shuffle(len(active), func(i, j int) { active[i], active[j] = active[j], active[i] })
	g.resolveFights(active)
	for _, p := range active {
		if p.Status == statusIn {
			g.playTurn(p)
		}
		p.Action = ""
	}
	for _, p := range g.players {
		if p != nil {
			g.meetIfNew(p)
		}
	}

	g.sun--
	g.soundAll("runda")
	inside := 0
	for _, p := range g.players {
		if p != nil && p.Status == statusIn {
			inside++
		}
	}
	if g.sun < 0 {
		g.sun = 0
	}
	switch {
	case inside == 0:
		g.endGame()
	case g.sun <= 0:
		g.logAll(g.data.Text("solnedgang"), "varn")
		g.soundAll("solnedgang")
		g.endGame()
	default:
		g.sunWarning(sunBefore)
		g.updateLevel()
		g.round++
		g.startRound()
	}
}

// sunWarning varnar när solen passerar halva, en fjärdedel och tre kvar.
// Klockfällan kan få solen att hoppa över en gräns, därför jämförs före och efter.
func (g *Game) sunWarning(before int) {
	passed := func(limit int) bool { return before > limit && g.sun <= limit }
	key := ""
	switch {
	case passed(3):
		key = "sol_snart"
	case passed(g.cfg.SunMax / 4):
		key = "sol_kvart"
	case passed(g.cfg.SunMax / 2):
		key = "sol_halv"
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
	g.releaseTakenOver()
	for _, p := range g.players {
		if p == nil {
			continue
		}
		p.Ready = p.bot != nil
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
	g.epilogue()
	g.logAll(g.data.Text("spela_igen"), "dim")
	g.soundAll("slut")
	g.broadcastState()
}

// playTurn utför en spelares drag för rundan, och räknar ner det som
// pågår (gift, mörker, fastklämd, källaren).
func (g *Game) playTurn(p *Player) {
	if p.Cellar > 0 {
		p.Cellar--
		if p.Cellar == 0 {
			g.leaveCellar(p)
		} else {
			g.logTo(p, g.data.Text("kallare"), "dim")
		}
		return
	}
	wasStuck, wasPoisoned, wasDark := p.Stuck > 0, p.Poison > 0, p.Dark > 0
	a := p.Action
	if a == "" {
		a = actWait
		if p.client != nil {
			g.logTo(p, g.data.Text("tvekade"), "dim")
		}
	}
	switch {
	case g.combatTurn(p, a):
	case isCombatAction(a):
		// Striden mot en annan vandrare är redan avgjord i resolveFights.
	case a == actWait:
	case a == actTake:
		g.take(p)
	case a == actSearch:
		g.search(p)
	case a == actDrink:
		g.drink(p)
	case p.Status != statusIn:
	default:
		if p.Stuck > 0 {
			g.logTo(p, g.data.Text("fast"), "dim")
		} else {
			g.move(p, dirCodes[a], 0)
		}
	}
	if p.Status != statusIn || p.Cellar > 0 {
		return
	}
	if wasStuck && p.Stuck > 0 {
		p.Stuck--
	}
	if wasPoisoned && p.Poison > 0 {
		p.Poison--
		g.logTo(p, g.data.Text("gift", "liv", strconv.Itoa(max(0, p.HP-1))), "varn")
		g.damage(p, 1, false)
		if p.Poison == 0 && p.Status == statusIn {
			g.logTo(p, g.data.Text("gift_slut"), "dim")
		}
	}
	if wasDark && p.Dark > 0 {
		p.Dark--
		if p.Dark == 0 {
			g.logTo(p, g.data.Text("morker_slut"), "info")
		}
	}
}

// move flyttar spelaren ett rum, eller ut genom porten. depth hindrar att
// roterande rum kastar runt spelaren i all oändlighet.
func (g *Game) move(p *Player, d Dir, depth int) {
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
		if p.Dark > 0 {
			g.logTo(p, g.data.Text("morker_nytt"), "dim")
		} else {
			g.logTo(p, name+".", "rubrik")
			g.logTo(p, g.describe(p, target, nx, ny), "")
		}
		g.logOthers(p, g.data.Text("nytt_rum_andra",
			"namn", p.Name, "region", g.data.region(nx, ny)), "info")
		if target.Event != "" {
			ev := target.Event
			target.Event = ""
			g.trigger(p, g.data.trap(ev), depth)
		}
	} else {
		g.soundTo(p, "steg")
		g.logTo(p, g.data.Text("kant_rum", "rum", name), "dim")
	}
	// En fälla kan ha flyttat spelaren; varelsen gäller rummet spelaren står i nu.
	if p.X == nx && p.Y == ny {
		g.announceMonster(p)
	}
}

type clientMsg struct {
	T    string `json:"t"`
	Name string `json:"name"`
	A    string `json:"a"`
	M    int    `json:"m"`
}

// serve sköter en websocket-anslutning tills den stängs.
func (g *Game) serve(conn *wsConn, local bool) {
	c := &Client{conn: conn, send: make(chan []byte, 64), local: local}
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
			g.Act(c, m.A, m.M)
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
