package main

import (
	"encoding/json"
	"math"
	"math/rand/v2"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	defaultSize = 9  // tornets storlek i rum, när inget annat är valt
	autoSunBase = 40 // solens rundor med automatisk sol i ett torn på defaultSize
	maxPlayers  = 8
	maxNameLen  = 16
	doorChance  = 0.55
)

// Spelets faser.
const (
	phaseLobby  = "lobby"
	phaseChoose = "klassval" // efter lobbyn: vandrarna väljer klass
	phasePlay   = "spel"
	phaseEnd    = "slut"
)

// En spelares läge.
const (
	statusIn   = "inne"
	statusOut  = "ute"
	statusLost = "forlorad"
	statusDead = "dod"
	statusNew  = "ny" // kom till tornet när vandringen redan var slut; med i nästa
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
	PvP          string // strid mellan vandrare: "pa", eller "av" (tomt) i testerna
	Classes      bool   // vandrarna väljer klass; alltid i spelet, av i de flesta testerna
	SunMax       int
	TurnTime     time.Duration
	MinPlayers   int
	Difficulty   string // id för ett förval i data/svarighet.json, tomt för standard
	Lang         string // språket: "sv" eller "en" (se lang.go)
	Seed         uint64 // frö för första vandringen, 0 för ett slumpat
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

// En port: ett rum i ytterväggen med fasta dörrar inåt och en utgång.
// Portarna i hörnen har två dörrar inåt, de mitt på sidorna tre.
type corner struct {
	x, y   int
	inward []Dir
	exit   Dir
}

// corners är tornets portar, räknade från storleken: en i varje hörn, och
// i ett torn på 11 rum eller mer även en mitt på varje sida. Ordningen
// följer orden under "horn" i texter.json.
func (g *Game) corners() []corner {
	s, c := g.size-1, g.size/2
	ports := []corner{
		{0, 0, []Dir{East, South}, West},
		{s, 0, []Dir{West, South}, East},
		{0, s, []Dir{East, North}, West},
		{s, s, []Dir{West, North}, East},
	}
	if g.size >= bigTower {
		ports = append(ports,
			corner{c, 0, []Dir{South, East, West}, North},
			corner{c, s, []Dir{North, East, West}, South},
			corner{0, c, []Dir{East, North, South}, West},
			corner{s, c, []Dir{West, North, South}, East},
		)
	}
	return ports
}

// bigTower är den minsta storleken med åtta portar.
const bigTower = 11

// vandrare räknar platserna som är tagna, AI-spelare inräknade.
func (g *Game) vandrare() int {
	n := 0
	for _, p := range g.players {
		if p != nil {
			n++
		}
	}
	return n
}

// g.center() är kammarens plats, mitt i tornet.
func (g *Game) center() int { return g.size / 2 }

func newGrid[T any](size int) [][]T {
	grid := make([][]T, size)
	for y := range grid {
		grid[y] = make([]T, size)
	}
	return grid
}

// towerSize är storleken på nästa torn. Automatiskt växer tornet med antalet
// vandrare: 9 rum för upp till fyra, 11 för fem eller sex och 13 för sju
// eller åtta. Storleken går inte att välja. Udda, så att kammaren står
// mitt i.
func (g *Game) towerSize() int {
	need := defaultSize
	switch n := g.vandrare(); {
	case n > 6:
		need = 13
	case n > 4:
		need = 11
	}
	return need
}

// scaledSun är solens rundor i det här tornet: automatiskt efter tornets
// storlek (autoSun) och svårigheten, eller det -sol har gett, för den som
// provar.
func (g *Game) scaledSun() int {
	if g.cfg.SunMax == 0 {
		if f := g.data.Regler.Sol; f > 0 {
			return max(4, int(math.Round(float64(autoSun(g.size))*f)))
		}
		return autoSun(g.size)
	}
	return max(4, g.cfg.SunMax)
}

// autoSun är solens rundor med automatisk sol: 40 i ett torn på 9 rum, och
// lika många fler i ett större torn som vägen till kammaren är längre
// (50 på 11, 60 på 13).
func autoSun(size int) int {
	return autoSunBase * (size - 1) / (defaultSize - 1)
}

// turnTime är längsta tid att välja drag per runda med de människor som är
// med nu. 0 i inställningen är automatiskt (autoTurn).
func (g *Game) turnTime() time.Duration {
	if g.cfg.TurnTime > 0 {
		return g.cfg.TurnTime
	}
	return autoTurn(g.humansPresent())
}

// autoTurn är turtiden med automatisk turtid: 30 sekunder för upp till fyra
// människor, 35 för fem eller sex och 40 för sju eller åtta. Datorspelare
// behöver ingen tid. Tillägget är litet, eftersom rundan avgörs så fort alla
// har valt och ett högre tak bara låter alla vänta längre på den som dröjer.
func autoTurn(humans int) time.Duration {
	return time.Duration(30+5*(max(0, humans-3)/2)) * time.Second
}

// portIndex numrerar vandrarna i platsordning, så att var och en får en
// egen port även när platserna inte följer på varandra.
func (g *Game) assignPorts() {
	i := 0
	for _, p := range g.players {
		if p != nil {
			p.Corner = i
			i++
		}
	}
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
	Walked   bool       // någon har gått in i rummet; en port är annars bara känd för den som började där
	trail    trail      // den som senast gick ut ur rummet, för spåren
	items    []string   // föremål som någon har lämnat här
	gold     []leftGold // guld som döda vandrare bar, åt den som kommer efter
	fell     []string   // vandrarna som har dött här, för rumsblocket
	goldBy   int        // vandraren som senast tog upp guld här, 0 om ingen
	omen     string     // det ett varsel har lovat: "falla:id", "varelse:id" eller "tomt"
	Shut     bool       // porten har slutit sig (tornet förändras, nivå 4)
}

// trail är spåren efter den som senast lämnade ett rum.
type trail struct {
	by    int // vandrarens id, 0 om ingen har gått här
	round int
	dir   Dir
}

type Player struct {
	ID      int
	Name    string
	Slot    int
	Corner  int     // porten vandraren börjar vid; lottas när vandrarna inte ser varandra
	classID string  // vald klass, tom om ingen
	items   []*Item // föremål i väskan

	initiated    bool    // vandraren har viskat tornets hemlighet
	guessLevel   int     // nivån då vandraren senast viskade
	parts        [3]bool // delarna av ormens namn som vandraren har fått
	partsFound   []int   // delarna i den ordning vandraren hittade dem
	calmed       bool    // den invigde har lugnat ormen under vandringen
	X, Y         int
	client       *Client
	Ready        bool
	Action       string
	Treasure     int
	Status       string
	Seen         map[string]bool // fällor och varelser spelaren har mött
	told         told            // vad vandraren senast fick veta om sitt rum (here.go)
	roundRoom    *roomMsg        // rumsblocket den här rundan, för Din vandring
	heard        []string        // det vandraren har hört under rundan (sounds.go)
	fights       map[string]bool // strider vandraren har hört den här rundan (sounds.go)
	fightsBefore map[string]bool // och förra rundan
	trapNow      string          // fällan som slog till mot vandraren under rundan; loggen har redan sagt det (here.go)
	known        [][]bool        // rum vandraren själv har varit i, det enda den ser på kartan
	bot          *botBrain       // satt för datorstyrda vandrare

	awayRound int  // rundan då människan lämnade mitt i vandringen, 0 om den är här
	takenOver bool // en AI-spelare styr vandraren medan människan är borta
	exitAsked int  // rundan då porten frågade om vandraren ville ut tomhänt, 0 om den inte har frågat
	aiPlayed  bool // AI:n har styrt vandraren någon gång under vandringen
	vard      bool // vandraren spelas från tornets dator, av spelledaren
	still     int  // rundor i rad som en ansluten människa inte har valt något drag

	target int // vandraren som ett stridsval den här rundan riktas mot
	// attackedBy är vandraren som senast gick mot den här i strid, och
	// attackedRound rundan det hände, så att tornet kan välja motståndare.
	attackedBy, attackedRound int
	met                       string // vilka andra vandrare som stod i rummet förra gången, så att det bara sägs en gång

	HP      int
	Potions int
	Stuck   int // rundor kvar fast (fallgaller, ormgrop, skinn)
	Cellar  int // rundor kvar i källarvalven
	Poison  int // rundor kvar med gift i kroppen
	Dark    int // rundor kvar utan ljus

	stats    chronStats     // det krönikan räknar
	journal  []journeyRound // rundans kort, runda för runda, för krönikan på egen skärm
	death    string         // hur vandraren dog, för krönikan på egen skärm
	epilog   []chronLine    // epilogens rader för just den här vandraren
	hurtKind string         // vad som skadade vandraren senast, för krönikan
	hurtWhat string         // namnet på fällan, varelsen eller vandraren
	// Det som skadade vandraren och varelsen hen fällde under rundan, för
	// Din vandring (summary.go). Töms före varje runda.
	roundHurtKind, roundHurtWhat, roundKill string
	roundCards                              []gameCard   // spelkorten för rundan (spelkort.go)
	lastCard                                *gameCard    // vandringens sista kort: ute, död eller kvar när solen gick ner
	goldTakenBy                             []int        // vandrarna som tog upp guldet den döda vandraren bar
	deaths                                  []heardDeath // dödsskrik vandraren har hört (guld.go)
}

// present är sant för en spelare som deltar just nu: en människa med öppen
// webbläsare, eller en datorspelare.
func (p *Player) present() bool { return p.client != nil || p.bot != nil }

type Client struct {
	conn   msgConn
	send   chan []byte
	player *Player
	closed bool
	local  bool // ansluten från datorn som kör tornet
}

type Game struct {
	slain    map[int]int // varelser som har fällts den här vandringen, och av vem (here.go)
	hearing  bool        // rundan avgörs, och ljuden samlas till rumsblocken (sounds.go)
	mu       sync.Mutex
	base     *Data // datafilerna som de är
	nextData *Data // nyss inlästa datafiler som gäller från nästa torn
	data     *Data // datafilerna med svårighetens förval
	cfg      Config
	size     int           // tornets bredd och höjd i rum
	sunMax   int           // solens rundor i det här tornet
	turn     time.Duration // turtiden under den här vandringen
	rooms    [][]Room
	players  [maxPlayers]*Player
	nextID   int
	level    int

	phase     string
	leaveWait time.Duration // hur länge lobbyn väntar på den som har gått innan vandringen kan starta
	sun       int
	round     int
	deadline  time.Time
	timer     *time.Timer
	chooseGen int // räknas upp vid varje klassval, så att en gammal tidsgräns inte startar något

	deck      []string // ormkorten som är kvar
	deckSleep int      // antal sovkort när leken byggdes
	wakes     int      // hur många gånger ormen har vaknat

	monsters      []*Monster
	nextMonsterID int
	rng           *rand.Rand // tornets slump; samma frö ger samma torn
	seed          uint64     // fröet för den pågående vandringen
	nextSeed      uint64     // fröet nästa vandring får (seed.go)

	usedClues map[string]bool // ledtrådar som redan har delats ut

	lastActive time.Time // senast en människa var ansluten eller en sida syntes på tornets dator
	lastReal   time.Time // senast någon gjorde något eller en sida syntes
	closeAt    time.Time // när tornet stänger om ingen kommer, satt under sista minuten

	paused    bool          // rundan står still: ingen människa är ansluten
	pauseLeft time.Duration // tid som var kvar när rundan pausades

	chron  []chronEvent // krönikan över vandringen
	epilog []chronLine  // epilogen som alla fick, för krönikan på egen skärm

	graceUsed bool // porten har redan väntat på de invigda

	serpentName    []string // ormens namn i tre delar, med invigningen "ormens namn"
	partAt         [3]int   // vilken del av namnet som finns i ristningen, hos varelsen och i kammaren
	carving        carving  // rummet med ristningen
	firstInitiated int      // id för den som blev invigd först i vandringen, 0 om ingen
	clueDeals      int      // ledtrådar som har delats ut, för att varannan ska bli ett varsel
	omens          []omenNote
	portClock      int // rundor sedan nivå 4 eller sedan en port stängdes
	portNext       int // porten som ska stängas, eller -1
}

func NewGame(d *Data, cfg Config) *Game {
	cfg.Difficulty = d.difficultyID(cfg.Difficulty)
	g := &Game{base: d, data: d.withDifficulty(cfg.Difficulty), cfg: cfg, level: 1, phase: phaseLobby, sun: cfg.SunMax, usedClues: map[string]bool{}, lastActive: time.Now(), lastReal: time.Now(), leaveWait: leaveWait}
	g.size = g.towerSize()
	g.rooms = newGrid[Room](g.size)
	g.sunMax = g.scaledSun()
	g.sun = g.sunMax
	g.nextSeed = cfg.Seed
	if g.nextSeed == 0 {
		g.nextSeed = randomSeed()
	}
	g.reseed(g.nextSeed) // också portarna i lobbyn följer fröet
	return g
}

func (g *Game) inBounds(x, y int) bool { return x >= 0 && y >= 0 && x < g.size && y < g.size }

func (g *Game) cornerIndex(x, y int) int {
	for i, c := range g.corners() {
		if c.x == x && c.y == y {
			return i
		}
	}
	return -1
}

// fixedDoors ger dörrarna för rum vars form är bestämd från början.
func (g *Game) fixedDoors(x, y int) ([4]bool, bool) {
	var d [4]bool
	if x == g.center() && y == g.center() {
		return [4]bool{true, true, true, true}, true
	}
	if ci := g.cornerIndex(x, y); ci >= 0 {
		for _, dir := range g.corners()[ci].inward {
			d[dir] = true
		}
		return d, true
	}
	return d, false
}

// reveal skapar rummet om det är okänt. Dörrarna stämmer alltid med
// grannar som redan finns, så att väggarna hänger ihop.
func (g *Game) reveal(x, y int) *Room {
	return g.build(x, y, true)
}

// build skapar rummet om det är okänt. Med spawn kan en varelse dyka upp i
// det; ett rum som en varelse vandrar in i får ingen egen.
func (g *Game) build(x, y int, spawn bool) *Room {
	r := &g.rooms[y][x]
	if r.Revealed {
		return r
	}
	r.Revealed = true
	if fd, ok := g.fixedDoors(x, y); ok {
		r.Doors = fd
		if ci := g.cornerIndex(x, y); ci >= 0 {
			r.Kind, r.TypeID = "ingang", "ingang"
			r.HasExit, r.Exit = true, g.corners()[ci].exit
		} else {
			r.Kind, r.TypeID = "kammare", "skattkammare"
		}
		g.keepWayOut(x, y)
		return r
	}
	r.Kind = "rum"
	r.TypeID = g.roomTypeFor(x, y)
	var undecided []Dir
	open := 0
	for d := North; d <= West; d++ {
		nx, ny := x+dirDX[d], y+dirDY[d]
		if !g.inBounds(nx, ny) {
			continue
		}
		n := &g.rooms[ny][nx]
		if n.Revealed {
			r.Doors[d] = n.Doors[d.Opp()]
		} else if fd, ok := g.fixedDoors(nx, ny); ok {
			r.Doors[d] = fd[d.Opp()]
		} else {
			undecided = append(undecided, d)
			r.Doors[d] = g.rng.Float64() < doorChance
		}
		if r.Doors[d] {
			open++
		}
	}
	// Minst två öppningar om det går, så att gångarna sällan blir återvändsgränder.
	for open < 2 && len(undecided) > 0 {
		i := g.rng.IntN(len(undecided))
		if d := undecided[i]; !r.Doors[d] {
			r.Doors[d] = true
			open++
		}
		undecided = append(undecided[:i], undecided[i+1:]...)
	}
	g.keepWayOut(x, y)
	g.placeCarving(x, y)
	// Händelse: en varelse, en fälla eller ingenting. Har ett varsel redan
	// sagt vad som väntar blir det så.
	if r.omen != "" {
		g.fulfilOmen(r, x, y)
		return r
	}
	if !(spawn && g.spawnMonster(x, y)) && !(g.data.Regler.TryggVidPorten && g.nextToPort(x, y)) &&
		g.rng.Float64() < g.data.Regler.Handelsechans {
		if t := g.randomTrap(); t != nil {
			r.Event = t.ID
		}
	}
	return r
}

// keepWayOut ser till att ingen blir instängd. Rummen byggs först när
// någon ser in i dem, och ibland pekar alla öppningar i ett nytt rum bakåt
// in i rum som redan finns, så att den del av tornet det hör till blir
// stängd. Har den delen då ingen öppning mot ett okänt rum, och når den inte
// kammaren, får det nya rummet en dörr till mot ett okänt grannrum, eller,
// om det inte har något, ett annat rum i samma del.
func (g *Game) keepWayOut(x, y int) {
	type pos struct{ x, y int }
	seen := map[pos]bool{{x, y}: true}
	part := []pos{{x, y}}
	for i := 0; i < len(part); i++ {
		c := part[i]
		r := &g.rooms[c.y][c.x]
		if r.Kind == "kammare" {
			return
		}
		for d := North; d <= West; d++ {
			nx, ny := c.x+dirDX[d], c.y+dirDY[d]
			if !r.Doors[d] || !g.inBounds(nx, ny) {
				continue
			}
			if !g.rooms[ny][nx].Revealed {
				return // en öppning mot ett okänt rum: det finns en väg vidare
			}
			if n := (pos{nx, ny}); !seen[n] {
				seen[n] = true
				part = append(part, n)
			}
		}
	}
	// Delen är stängd. Öppna en vägg mot ett okänt grannrum, helst i det
	// nya rummet.
	for _, c := range part {
		r := &g.rooms[c.y][c.x]
		if r.Kind != "rum" {
			continue
		}
		var cand []Dir
		for d := North; d <= West; d++ {
			nx, ny := c.x+dirDX[d], c.y+dirDY[d]
			if r.Doors[d] || !g.inBounds(nx, ny) || g.rooms[ny][nx].Revealed {
				continue
			}
			if _, fixed := g.fixedDoors(nx, ny); fixed {
				continue
			}
			cand = append(cand, d)
		}
		if len(cand) > 0 {
			r.Doors[cand[g.rng.IntN(len(cand))]] = true
			return
		}
	}
}

// roomTypeFor väljer rumstyp åt ett nytt rum: en av de typer som har
// använts minst hittills i tornet, och aldrig samma som ett grannrum, så att
// vandrarna kan säga till varandra vilket rum de står i.
func (g *Game) roomTypeFor(x, y int) string {
	types := g.data.Rumstyper
	used := map[string]int{}
	for yy := range g.size {
		for xx := range g.size {
			if r := &g.rooms[yy][xx]; r.Revealed && r.Kind == "rum" {
				used[r.TypeID]++
			}
		}
	}
	near := map[string]bool{}
	for d := North; d <= West; d++ {
		if nx, ny := x+dirDX[d], y+dirDY[d]; g.inBounds(nx, ny) && g.rooms[ny][nx].Revealed {
			near[g.rooms[ny][nx].TypeID] = true
		}
	}
	var pick []string
	least := -1
	for _, t := range types {
		if near[t.ID] {
			continue
		}
		switch n := used[t.ID]; {
		case least < 0 || n < least:
			least, pick = n, []string{t.ID}
		case n == least:
			pick = append(pick, t.ID)
		}
	}
	if len(pick) == 0 { // färre typer än grannar
		return types[g.rng.IntN(len(types))].ID
	}
	return pick[g.rng.IntN(len(pick))]
}

// nextToPort: rummet ligger bredvid en port.
func (g *Game) nextToPort(x, y int) bool {
	for _, c := range g.corners() {
		if abs(c.x-x)+abs(c.y-y) == 1 {
			return true
		}
	}
	return false
}

// roomShort är rummet som spelarna får det, beskrivet och utan namn: "en
// smal gång under låga valv". Rummen har namn bara på värdsidan.
func (g *Game) roomShort(r *Room) string {
	return g.data.roomType(r.TypeID).Kort
}

// describe ger spelaren en beskrivning av rummet för tornets nuvarande nivå.
func (g *Game) describe(p *Player, r *Room, x, y int) string {
	return g.data.roomType(r.TypeID).Describe(g.level, p.Slot+x+y)
}

// Meddelanden till klienten.

type playerView struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Slot       int        `json:"slot"`
	X          int        `json:"x"`
	Y          int        `json:"y"`
	Online     bool       `json:"online"`
	Status     string     `json:"status"`
	Ready      bool       `json:"ready"`
	Chosen     bool       `json:"chosen"`
	Treasure   int        `json:"treasure"` // -1 om spelaren inte får veta
	Bot        bool       `json:"bot"`
	Hidden     bool       `json:"hidden"`                // positionen är okänd för den som tittar
	AI         bool       `json:"ai"`                    // en AI-spelare styr medan människan är borta eller står still
	Class      string     `json:"klass,omitempty"`       // klassen med liten bokstav
	ClassID    string     `json:"klassId,omitempty"`     // klassens id, för figuren på kartan snett ovanifrån
	ClassModel *ModelLook `json:"klassModell,omitempty"` // färg och storlek på figuren
}

type roomView struct {
	X     int      `json:"x"`
	Y     int      `json:"y"`
	D     [4]bool  `json:"d"`
	K     string   `json:"k"`
	Exit  int      `json:"exit"`
	Mark  string   `json:"mark,omitempty"`
	MarkI string   `json:"marki,omitempty"`  // fällans id, för modellen på kartan snett ovanifrån
	S     bool     `json:"s,omitempty"`      // genomsökt
	F     bool     `json:"f,omitempty"`      // här ligger ett föremål eller guld (bara för spårfinnaren)
	Fell  bool     `json:"dod,omitempty"`    // här har någon dött (bara för kartritaren)
	In    []string `json:"in,omitempty"`     // inredningen på kartan snett ovanifrån
	Shut  bool     `json:"stangd,omitempty"` // porten har slutit sig
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
	Sprak      string        `json:"sprak"` // texternas språk; byts det hämtar sidan nya texter
	You        int           `json:"you"`
	Size       int           `json:"size"`
	Level      int           `json:"level"`
	Reveal     bool          `json:"reveal"` // alltid sant; kartan i klient/ läser det
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
	Monsters   []monsterView `json:"monsters"`
	PvP        string        `json:"pvp"` // om det går att strida mot andra vandrare här
	// Motstandare är vandraren tornet väljer om du slåss här (pickFoe).
	Motstandare string       `json:"motstandare,omitempty"`
	Paused      bool         `json:"paused,omitempty"`
	Vard        bool         `json:"vard,omitempty"` // spelar från tornets dator: spelledaren
	SleepText   string       `json:"sleepText,omitempty"`
	WakeCards   *int         `json:"vaknaKort,omitempty"` // vakna-korten kvar i leken, bara för ormtjusaren
	Class       *classView   `json:"klass,omitempty"`
	Items       []itemView   `json:"foremal,omitempty"`
	MaxItems    int          `json:"maxForemal,omitempty"` // så många föremål man kan bära
	Compass     string       `json:"kompass,omitempty"`
	Heart       string       `json:"hjartslag,omitempty"`
	Initiated   bool         `json:"invigd,omitempty"`
	CanWhisper  bool         `json:"kanViska,omitempty"`
	NameParts   []string     `json:"namndelar,omitempty"` // delarna av ormens namn, tomma där vandraren inte vet
	Offer       []*classView `json:"erbjudande,omitempty"`
	Players     []playerView `json:"players"`
	Rooms       []roomView   `json:"rooms"`
	Results     []resultView `json:"results,omitempty"`
	Journey     *journeyView `json:"vandring,omitempty"` // den egna vandringen, för den som har dött och för alla i slutet
	Chronicle   []chronLine  `json:"kronika,omitempty"`  // hela krönikan, när vandringen är slut
	Epilogue    []chronLine  `json:"epilog,omitempty"`   // epilogen, när vandringen är slut
	Moments     []momentView `json:"stunder,omitempty"`  // vandringens största stunder, när vandringen är slut
}

type logMsg struct {
	T    string `json:"t"`
	Text string `json:"text"`
	C    string `json:"c,omitempty"`
	// Svar är sant för svaret på ett drag som inte blev av, som en vägg.
	// Spelarsidan lägger det i rundan som står i bild i stället för att
	// börja en ny, så att rummet man står i inte försvinner ur bild.
	Svar bool `json:"svar,omitempty"`
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

// answerTo svarar på ett drag som inte blev av (se logMsg.Svar).
func (g *Game) answerTo(p *Player, text, class string) {
	if p != nil && text != "" {
		g.send(p.client, logMsg{T: "log", Text: text, C: class, Svar: true})
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

// sees är sant om me får veta var o är. Andra vandrare syns bara i samma rum
// (och nära lyktbäraren), och när vandringen är slut syns alla igen.
func (g *Game) sees(me, o *Player) bool {
	if g.phase == phaseEnd || g.nearLantern(me, o) {
		return true
	}
	return me.Status == statusIn && o.Status == statusIn && me.Cellar == 0 &&
		o.X == me.X && o.Y == me.Y
}

// hiddenText väljer den vaga varianten av en text (nyckeln med "_dolt"),
// som inte säger var den andra är, om den finns i datafilerna.
func (g *Game) hiddenText(key string, vars ...string) string {
	if len(g.data.Texter[key+"_dolt"]) > 0 {
		key += "_dolt"
	}
	return g.data.Text(key, vars...)
}

// stateFor bygger det en viss spelare får se. Andras skatt och val av drag
// hålls hemliga tills spelet är slut.
func (g *Game) stateFor(me *Player) stateMsg {
	g.noteRoom(me)
	st := stateMsg{
		T: "state", Sprak: g.base.Lang, You: me.ID, Size: g.size, Level: g.level, Reveal: true,
		Phase: g.phase, Sun: g.sun, SunMax: g.sunMax, Round: g.round,
		MinPlayers: g.cfg.MinPlayers, Action: me.Action,
		HP: me.HP, MaxHP: g.maxHP(me), Potions: me.Potions, MaxItems: g.data.maxItems(),
		Stuck: me.Stuck, Cellar: me.Cellar, Poison: me.Poison, Dark: me.Dark,
		Vard: me.client != nil && me.client.local,
	}
	// I mörker ser man bara sig själv och rummet man står i, utan dess
	// dörrar och port.
	blind := g.phase == phasePlay && me.Status == statusIn && me.Dark > 0
	st.Sleep = -1
	if g.phase == phasePlay && me.Status == statusIn && me.Cellar == 0 && g.rooms[me.Y][me.X].Kind == "kammare" {
		st.Sleep = g.sleepLevel()
		st.SleepText = g.data.Ormen.Somn[min(st.Sleep, len(g.data.Ormen.Somn)-1)]
		if g.is(me, "ormtjusare") {
			n := g.wakeCardsLeft()
			st.SleepText += " " + g.data.Text("klass_ormtjusare_kort", "antal", itoa(n))
			st.WakeCards = &n
		}
	}
	st.Class = viewOfClass(g.classOf(me))
	st.Items = g.itemViews(me)
	if g.phase == phasePlay && me.Status == statusIn {
		g.curse(me, "rostkrona") // kronan syns redan första rundan
	}
	if g.phase == phasePlay {
		st.Compass = g.compass(me)
		st.Heart = g.heartbeat(me)
		st.Initiated = me.initiated
		st.CanWhisper = g.canWhisper(me)
		st.NameParts = g.namePartsFor(me)
	}
	if g.phase == phaseChoose {
		st.Offer = g.classCards(me)
		st.DeadlineMs = max(0, time.Until(g.deadline).Milliseconds())
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
			Chosen: p.Action != "", Treasure: -1, AI: p.takenOver, Class: g.classLabel(p), ClassID: g.classID(p), ClassModel: g.classModel(p),
		}
		if p == me || g.phase == phaseEnd {
			pv.Treasure = p.Treasure
		}
		if p.Cellar > 0 || (blind && p != me) || (p != me && !g.sees(me, p)) {
			pv.Hidden = true
			pv.X, pv.Y = -1, -1
		}
		st.Players = append(st.Players, pv)
	}
	for y := range g.size {
		for x := range g.size {
			r := &g.rooms[y][x]
			if !g.knowsRoom(me, x, y) || (blind && (x != me.X || y != me.Y)) {
				continue
			}
			rv := roomView{X: x, Y: y, D: r.Doors, K: r.Kind, Exit: -1, Mark: r.Mark, S: r.Searched}
			if t := g.data.trapByMark(r.Mark); t != nil {
				rv.MarkI = t.ID
			}
			if t := g.data.roomType(r.TypeID); r.Kind == "rum" && t.Modell != nil {
				rv.In = t.Modell.Inredning
			}
			rv.F = (len(r.items) > 0 || len(r.gold) > 0) && g.is(me, "sparfinnare")
			rv.Fell = len(r.fell) > 0 && g.is(me, "kartritare")
			if r.HasExit {
				rv.Exit = int(r.Exit)
				rv.Shut = r.Shut
			}
			if blind {
				// I mörkret ser man inte dörrarna eller porten och får pröva
				// sig fram (actPlayer): kartan ritar hela murar och inga
				// tangenter vid dörrarna.
				rv.D = [4]bool{}
				rv.Exit, rv.Shut = -1, false
			}
			st.Rooms = append(st.Rooms, rv)
		}
	}
	st.Monsters = g.monsterViews(me, blind)
	if g.phase == phasePlay {
		st.PvP = g.pvpStatus(me)
		if st.PvP == pvpYes && g.monsterAt(me.X, me.Y) == nil {
			if f := g.pickFoe(me, 0); f != nil {
				st.Motstandare = f.Name
			}
		}
	}
	if g.phase == phaseEnd {
		st.Results = g.results()
	}
	if (g.phase == phasePlay && me.Status == statusDead) || (g.phase == phaseEnd && me.Status != statusNew) {
		st.Journey = g.journeyOf(me)
	}
	if g.phase == phaseEnd {
		st.Chronicle = g.chronicleLines()
		st.Epilogue = append(slices.Clone(g.epilog), me.epilog...)
		st.Moments = g.moments()
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

// results sorterar spelarna: de som kom ut med skatt först, mest skatt överst,
// sedan de som kom ut tomhänta och sist de döda och förlorade.
func (g *Game) results() []resultView {
	var res []resultView
	for _, p := range g.players {
		if p == nil || p.Status == statusNew {
			continue
		}
		res = append(res, resultView{Name: p.Name, Slot: p.Slot, Treasure: p.Treasure, Status: p.Status, AI: p.aiPlayed})
	}
	score := func(r resultView) int {
		if r.Status == statusOut {
			return r.Treasure // tomhänt ger 0, före de döda och förlorade
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
	cn := g.corners()[p.Corner]
	p.X, p.Y = cn.x, cn.y
	p.Status = statusIn
	p.Treasure = 0
	p.Action = ""
	p.Ready = p.bot != nil // datorspelare är alltid redo
	p.HP = g.maxHP(p)
	p.Seen = map[string]bool{}
	p.Potions, p.Stuck, p.Cellar, p.Poison, p.Dark = 0, 0, 0, 0, 0
	p.known = newGrid[bool](g.size)
	p.items = nil
	p.initiated, p.guessLevel, p.calmed = false, 0, false
	p.parts, p.partsFound = [3]bool{}, nil
	p.forgetRoom()
	g.revealPorts()
	g.noteRoom(p)
}

// noteRoom kommer ihåg att vandraren har varit i rummet den står i.
func (g *Game) noteRoom(p *Player) {
	if p.Status == statusIn && p.Cellar == 0 && g.inBounds(p.X, p.Y) {
		p.known[p.Y][p.X] = true
		g.peekNeighbors(p)
	}
}

// knowsRoom är sant om me ser rummet på sin karta. Man ser bara rum man
// själv har varit i, och av portarna bara den man kom in genom; kartan
// växer därifrån. Kartritaren ser alla rum som någon har gått in i. Rum som någon bara
// har sett in i, och de andras portar, är okända tills någon går in i dem: från början känner var och en bara
// sitt eget startrum. Den döda spelaren som styr varelserna ser hela
// tornet, och när vandringen är slut syns allt.
func (g *Game) knowsRoom(me *Player, x, y int) bool {
	r := &g.rooms[y][x]
	if !r.Revealed {
		return false
	}
	if me == nil || g.phase == phaseEnd {
		return true
	}
	if me.known != nil && me.known[y][x] {
		return true
	}
	if !r.Walked {
		return false // en annans port, eller ett rum som någon bara har sett in i
	}
	return g.is(me, "kartritare") // kartritaren ser alla rum som någon har hittat
}

func (g *Game) Join(c *Client, rawName string) {
	if c.player != nil {
		return
	}
	// Spelarsidan skickar namnet tornet gav den förra gången, eller inget.
	// Utan namn lottar tornet ett (names.go).
	name := cleanName(rawName)
	if name == "" {
		name = g.randomName()
	}
	// Återanslutning: samma namn tar tillbaka sin vandrare. Är namnet redan
	// upptaget, som när två flikar i samma webbläsare öppnar tornet, får den
	// nya ett eget namn.
	for _, p := range g.players {
		if p != nil && strings.EqualFold(p.Name, name) {
			if p.client != nil || (p.bot != nil && !p.takenOver) {
				name = g.randomName()
				break
			}
			if p.takenOver {
				g.handBack(p, c)
				return
			}
			p.awayRound, p.still = 0, 0
			p.client, c.player = c, p
			p.vard = c.local
			g.checkAwait()
			g.broadcastState()
			g.logTo(p, g.data.Text("tillbaka"), "syn")
			g.retellRoom(p)
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
	p.vard = c.local
	g.fitTower(p)
	if g.phase == phaseEnd {
		// Vandringen är redan slut: den nya väntar på nästa torn och står
		// inte i resultatet eller krönikan.
		p.Status = statusNew
		g.broadcastState()
		g.soundAll("ansluter")
		g.logTo(p, g.data.Text("namn_givet", "namn", name), "rost")
		g.logTo(p, g.data.Text("sen_ankomst", "namn", name), "syn")
		g.logOthers(p, g.hiddenText("ansluter", "namn", name, "horn", g.data.Horn[slot]), "info")
		return
	}
	horn := g.data.Horn[slot]
	if g.phase == phaseChoose {
		g.logTo(p, g.data.Text("klass_val"), "info")
	}

	g.broadcastState()
	g.soundAll("ansluter")
	g.logTo(p, g.data.Text("namn_givet", "namn", name), "rost")
	g.logTo(p, g.hiddenText("valkommen", "namn", name, "horn", horn), "syn")
	// Portrummet beskrivs först när vandringen börjar (atTheGate i here.go).
	for _, b := range g.players {
		if b != nil && b.bot != nil {
			g.logTo(p, fill(pick(g.aiTexts(g.data.AI.Ansluter, g.data.AI.AnsluterDolt)), "namn", b.Name, "horn", g.data.Horn[b.Corner]), "info")
		}
	}
	g.logTo(p, g.data.Text("lobby", "min", strconv.Itoa(g.cfg.MinPlayers)), "info")
	g.logOthers(p, g.hiddenText("ansluter", "namn", name, "horn", horn), "info")
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
	g.checkAwait()
	if g.phase == phaseChoose {
		// Under klassvalet behåller man sin plats, som under vandringen, så
		// att den som laddar om sidan kommer tillbaka. Klassvalet väntar på
		// hen tills tiden är slut; sedan följer hen med som frånvarande.
		// AI-spelarna kan behöva välja nu när hen inte längre väljer.
		if g.checkChosen() {
			return
		}
	} else if g.phase == phaseEnd {
		// I slutet står man kvar i resultatet och krönikan, så att den som
		// laddar om sidan, eller vars telefon somnar, kommer tillbaka till
		// sin egen vandring. Den som inte kommer tillbaka tas bort när nästa
		// torn byggs (dropAbsent). Det nya tornet väntar en stund på hen.
		g.startAfterLeave()
	} else if g.phase != phasePlay {
		// I lobbyn lämnar man sin plats helt, och en datorspelare kan ta den
		// om spelledaren vill ha fler. Vandringen startar inte genast även om
		// alla som är kvar är redo: den som laddar om sidan ska hinna tillbaka.
		g.players[p.Slot] = nil
		g.fitTower(nil)
		g.syncBots(true)
		g.startAfterLeave()
	} else if g.allChosen() && !g.paused {
		g.resolve()
		return
	}
	g.broadcastState()
}

// leaveWait är hur länge lobbyn och slutet väntar på den som har gått innan
// vandringen kan starta utan hen.
const leaveWait = 10 * time.Second

// startAfterLeave startar vandringen när leaveWait har gått, om alla som är
// kvar då är redo och ingen vandring har börjat.
func (g *Game) startAfterLeave() {
	if g.leaveWait <= 0 {
		g.checkStart()
		return
	}
	time.AfterFunc(g.leaveWait, func() {
		g.mu.Lock()
		defer g.mu.Unlock()
		if g.phase == phaseLobby || g.phase == phaseEnd {
			g.checkStart()
		}
	})
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
	if g.phase == phaseChoose {
		return g.checkChosen()
	}
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
	if g.phase == phaseEnd && g.cfg.Classes {
		// Nytt torn: tillbaka till lobbyn, där alla får välja klass igen.
		g.toLobby()
		return true
	}
	if g.phase == phaseEnd {
		g.towerForNext()
	}
	g.begin()
	return true
}

// begin går från lobbyn till klassvalet, eller direkt in i vandringen när
// klasserna är avstängda.
func (g *Game) begin() {
	// De som lämnade förra vandringen och inte kom tillbaka är inte med längre.
	// Det görs här och inte när vandringen startar, så att den som tappar
	// anslutningen under klassvalet får behålla sin plats.
	if g.dropAbsent() > 0 {
		g.syncBots(true) // AI-spelare kan ta de lediga platserna
	}
	if g.cfg.Classes && len(g.data.Klasser) > 0 {
		g.startChoosing()
		return
	}
	g.startGame()
}

// toLobby bygger ett nytt torn efter en vandring och låter alla vänta i
// lobbyn. När alla är redo igen väljer de nya klasser.
func (g *Game) toLobby() {
	if g.nextData != nil {
		g.useData(g.nextData)
	}
	g.phase = phaseLobby
	g.level = 1 // tornet i lobbyn har inte börjat avslöja något
	if g.dropAbsent() > 0 {
		g.syncBots(true)
	}
	g.resetClasses()
	for _, p := range g.players {
		if p != nil {
			p.Ready = false // newTower behåller den som var redo
		}
	}
	g.towerForNext()
	g.sun = g.sunMax
	for _, p := range g.players {
		if p != nil && p.client != nil {
			g.logTo(p, g.data.Text("nytt_torn", "min", strconv.Itoa(g.cfg.MinPlayers)), "syn")
		}
	}
	g.broadcastState()
}

// newTower bygger ett nytt torn och ställer alla vid sina portar.
func (g *Game) newTower() {
	g.buildTower(g.towerSize())
}

// buildTower bygger ett torn med size × size rum. Storleken följer annars
// antalet vandrare (towerSize); testerna väljer den själva.
func (g *Game) buildTower(size int) {
	g.size = size
	g.rooms = newGrid[Room](g.size)
	g.sunMax = g.scaledSun()
	if g.phase != phasePlay {
		g.sun = g.sunMax
	}
	g.monsters = nil
	g.assignPorts()
	for _, p := range g.players {
		if p != nil {
			ready := p.Ready // den som redan är redo förblir det
			g.placeAtCorner(p)
			p.Ready = ready || p.bot != nil
		}
	}
	g.revealPorts()
}

// revealPorts låter portrummen finnas från början. Var och en känner bara
// sin egen tills någon går in i en annans (se knowsRoom).
func (g *Game) revealPorts() {
	for _, cn := range g.corners() {
		g.reveal(cn.x, cn.y)
	}
}

// drawCorners lottar vid vilken port var och en börjar, så att ordningen i
// listan och färgen inte avslöjar något. Var och en får veta sin egen port.
func (g *Game) drawCorners() {
	order := g.rng.Perm(len(g.corners()))
	i := 0
	for _, p := range g.players {
		if p == nil {
			continue
		}
		p.Corner = order[i]
		i++
		cn := g.corners()[p.Corner]
		p.X, p.Y = cn.x, cn.y
		p.known = newGrid[bool](g.size)
		g.reveal(cn.x, cn.y)
	}
}

func (g *Game) startGame() {
	g.paused = false
	g.seedWandering()
	g.ensureClasses()
	if g.towerSize() != g.size {
		g.newTower() // fler eller färre vandrare än när tornet byggdes
	}
	g.phase = phasePlay
	g.sunMax = g.scaledSun()
	g.sun = g.sunMax
	g.turn = g.turnTime() // ligger fast under vandringen
	g.round = 1
	g.wakes = 0
	g.firstInitiated = 0
	g.clueDeals, g.omens = 0, nil
	g.portClock, g.portNext = 0, -1
	g.buildDeck()
	g.level = 1
	g.graceUsed = false
	g.usedClues = map[string]bool{}
	g.pickSerpentName()
	g.resetChronicle()
	g.slain = map[int]int{}
	for _, p := range g.players {
		if p != nil {
			p.Ready = p.bot != nil
			p.Action = ""
			p.HP = g.maxHP(p) // svårigheten kan ha ändrats sedan tornet byggdes
			p.Potions = g.startPotions(p)
			p.awayRound, p.exitAsked, p.still = 0, 0, 0
			if p.bot == nil && p.client == nil {
				p.awayRound = 1 // lämnade under klassvalet: borta från första rundan
			}
			p.aiPlayed = false
		}
	}
	g.soundAll("start")
	if s := g.data.difficulty(g.cfg.Difficulty); s != nil {
		// Livet efter klassen, så att legoknekten får höra att hen har ett extra.
		for _, p := range g.players {
			if p != nil {
				g.logTo(p, g.data.Text("svarighet", "namn", s.Namn, "liv", strconv.Itoa(g.maxHP(p))), "dim")
			}
		}
	}
	for _, p := range g.players {
		if p != nil {
			g.logTo(p, g.data.Text(g.visibilityKey(p)), "dim")
		}
	}
	g.revealPorts()
	g.drawCorners()
	for _, p := range g.players {
		if p != nil {
			g.noteRoom(p) // efter klassvalet: kartritaren ser in i grannrummen
			p.forgetRoom()
			g.describeLater(p, &g.rooms[p.Y][p.X], p.X, p.Y)
			g.atTheGate(p)
			g.tellRoom(p)
		}
	}
	g.startRound()
}

func (g *Game) startRound() {
	wait := g.turn
	if wait == 0 {
		wait = g.turnTime()
	}
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
		if p.Status != statusIn || p.Cellar > 0 {
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
// Styr tornet vandraren för att människan stod still, får hen tillbaka den.
func (g *Game) Act(c *Client, a string, monsterID int) {
	if p := c.player; p != nil && p.takenOver && p.client == c {
		g.reclaim(p)
	}
	g.actPlayer(c.player, a, monsterID)
}

// actPlayer tar emot ett val från en människa eller en datorspelare.
func (g *Game) actPlayer(p *Player, a string, monsterID int) {
	if p == nil || g.phase != phasePlay || p.Status != statusIn || p.Cellar > 0 {
		return
	}
	room := &g.rooms[p.Y][p.X]
	// Sök tar en näve i kammaren: där finns inget annat att söka efter, och
	// spelarna behöver bara en tangent för båda.
	if a == actSearch && room.Kind == "kammare" {
		a = actTake
	}
	switch a {
	case actAttack, actParry, actDodge:
		if g.monsterAt(p.X, p.Y) == nil {
			// Ingen varelse: kanske en annan vandrare. Tornet väljer vem;
			// en datorspelare kan peka ut någon med ett negativt id.
			switch g.pvpStatus(p) {
			case pvpYes:
				p.target = g.pickFoe(p, -monsterID).ID
			case pvpPeace:
				g.answerTo(p, g.data.Text("pvp_fredat"), "dim")
				return
			default:
				g.answerTo(p, g.data.MonsterText("inget_monster_har"), "dim")
				return
			}
		}
	case "N", "E", "S", "W", actWait:
		// Mot en vägg: avvisas direkt, som ett rum som redan är genomsökt,
		// så att ett feltryck inte kostar rundan. I mörkret ser man inte
		// dörrarna och får pröva sig fram som förut.
		if d, ok := dirCodes[a]; ok && p.Dark == 0 && !(room.HasExit && room.Exit == d) && !room.Doors[d] {
			g.soundTo(p, "vagg")
			g.answerTo(p, g.data.Text("vagg", "riktning", g.data.Riktningar[d]), "dim")
			return
		}
		if g.askEmptyExit(p, room, a) {
			return
		}
	case actTake:
		if room.Kind != "kammare" {
			g.answerTo(p, g.data.Text("ingen_skatt_har"), "dim")
			return
		}
	case actSearch:
		if room.Kind != "rum" {
			g.answerTo(p, g.data.SearchText("inte_har"), "dim")
			return
		}
		if room.Searched {
			g.answerTo(p, g.data.SearchText("redan"), "dim")
			return
		}
	case actDrink:
		if p.Potions <= 0 {
			g.answerTo(p, g.data.Text("ingen_dryck"), "dim")
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

// askEmptyExit: en människa som går mot en öppen port utan guld får först en
// fråga, så att ingen råkar gå ut tomhänt med ett enda tryck. Går hen mot
// porten igen samma eller nästa runda går hen ut. Svarar true när draget blev en fråga.
func (g *Game) askEmptyExit(p *Player, room *Room, a string) bool {
	d, ok := dirCodes[a]
	if !ok || p.bot != nil || p.Stuck > 0 || !room.HasExit || room.Exit != d || room.Shut {
		return false
	}
	// Frågan gäller också nästa runda, om rundan hann ta slut före svaret.
	if p.Treasure+g.itemGold(p) > 0 || (p.exitAsked > 0 && g.round-p.exitAsked <= 1) {
		return false
	}
	p.exitAsked = g.round
	g.answerTo(p, g.data.Text("port_tomhant"), "varn")
	return true
}

// resolve utför alla drag på en gång och låter solen sjunka.
func (g *Game) resolve() {
	if g.timer != nil {
		g.timer.Stop()
	}
	sunBefore := g.sun
	snap, wakes := g.snapshot(), g.wakes
	g.hearing = true // ljuden samlas till rumsblocken
	g.wanderMonsters()
	var active []*Player
	for _, p := range g.players {
		if p != nil && p.Status == statusIn {
			active = append(active, p)
		}
	}
	g.countStill(active)
	g.rng.Shuffle(len(active), func(i, j int) { active[i], active[j] = active[j], active[i] })
	g.resolveFights(active)
	for _, p := range active {
		if p.Status == statusIn {
			g.playTurn(p)
		}
		p.Action = ""
	}
	// Den som står kvar där någon har dött tar guldet som ligger där.
	for _, p := range active {
		g.pickUpGold(p)
	}
	for _, p := range g.players {
		if p != nil {
			g.meetIfNew(p)
			g.readCarving(p)
		}
	}
	g.tellRooms()
	g.hearing = false

	g.sun--
	if g.sun < 0 {
		g.sun = 0
	}
	g.hourglass()
	g.summarize(snap, wakes)
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
	case g.sun <= 0 && !g.graceUsed && g.anyInitiatedInside():
		text := g.data.Text("solnedgang")
		g.logAll(text, "varn")
		g.bigCardInside(gameCard{Slag: "sol", Bild: "sol", Rubrik: g.data.Text("kort_sol"), Text: text})
		g.soundAll("solnedgang")
		g.portGrace()
		g.round++
		g.startRound()
	case g.sun <= 0:
		if !g.graceUsed {
			g.logAll(g.data.Text("solnedgang"), "varn")
			g.soundAll("solnedgang")
		}
		g.endGame()
	default:
		g.sunWarning(sunBefore)
		g.updateLevel()
		g.closePorts()
		g.round++
		g.startRound()
	}
}

// sunWarning varnar när solen passerar halva, en fjärdedel och tre kvar.
// Klockfällan kan få solen att hoppa över en gräns, därför jämförs före och efter.
// Med en kort sol hamnar halva och en fjärdedel vid eller under tre rundor;
// de varningarna hoppas då över, så att "halva vägen" inte kommer efter
// "nuddar nästan horisonten".
func (g *Game) sunWarning(before int) {
	passed := func(limit int) bool { return before > limit && g.sun <= limit }
	early := func(limit int) bool { return limit > 3 && passed(limit) }
	key := ""
	switch {
	case passed(3):
		key = "sol_snart"
	case early(g.sunMax / 4):
		key = "sol_kvart"
	case early(g.sunMax / 2):
		key = "sol_halv"
	}
	if key == "sol_snart" && g.sun == 1 {
		key = "sol_snart_en" // klockfällan kan hoppa ända ner till en runda
	}
	if key != "" {
		text := g.data.Text(key, "sol", strconv.Itoa(g.sun))
		g.logAll(text, "varn")
		g.soundAll("solvarning")
		if key == "sol_snart" {
			g.bigCardInside(gameCard{Slag: "sol", Bild: "sol", Rubrik: g.data.Text("kort_sol"), Text: text})
		}
	}
}

func (g *Game) endGame() {
	g.phase = phaseEnd
	g.paused = false
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
			g.chronicle("kronika_forlorad", "namn", p.Name)
			text := g.data.Text("forlorad_du")
			g.logTo(p, text, "varn")
			g.endCard(p, gameCard{Slag: "forlorad", Bild: "sol", Rubrik: g.data.Text("kort_forlorad"), Text: text, Stor: true})
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
	g.goldEpilogue()
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
		p.hurt("gift", "")
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
	if room.HasExit && room.Exit == d && room.Shut {
		g.soundTo(p, "vagg")
		g.logTo(p, g.data.Text("port_stangd_har"), "varn")
		return
	}
	if room.HasExit && room.Exit == d || room.Doors[d] {
		room.trail = trail{by: p.ID, round: g.round, dir: d}
	}
	if room.HasExit && room.Exit == d {
		p.Status = statusOut
		if gold := g.itemGold(p); gold > 0 {
			p.Treasure += gold
			g.logTo(p, g.data.Text("foremal_guld", "antal", strconv.Itoa(gold)), "syn")
		}
		if p.initiated && p.Treasure > 0 {
			share := g.data.Avslojande.Invigd.Guld
			if p.ID == g.firstInitiated {
				share += g.data.Avslojande.Invigd.ForstGuld
			}
			if bonus := int(float64(p.Treasure) * share); bonus > 0 {
				p.Treasure += bonus
				g.logTo(p, g.data.Text("invigd_guld", "antal", strconv.Itoa(bonus)), "rost")
			}
		}
		if p.Treasure > 0 {
			g.chronicle("kronika_ut", "namn", p.Name, "guld", strconv.Itoa(p.Treasure))
		} else {
			g.chronicle("kronika_ut_tom", "namn", p.Name)
		}
		var text string
		if p.Treasure > 0 {
			text = g.data.Text("ut_med_skatt", "skatt", strconv.Itoa(p.Treasure))
			g.logTo(p, text, "syn")
		} else {
			text = g.data.Text("ut_utan_skatt")
			g.logTo(p, text, "info")
		}
		g.endCard(p, gameCard{Slag: "ut", Bild: "ut", Rubrik: g.data.Text("kort_ut"), Text: text, Stor: true})
		g.soundTo(p, "ut")
		g.logOthers(p, g.data.Text("andra_ut", "namn", p.Name), "info")
		g.hear(p.X, p.Y, "ljud_port", "", "dim", p)
		return
	}
	if !room.Doors[d] {
		g.soundTo(p, "vagg")
		g.logTo(p, g.data.Text("vagg", "riktning", g.data.Riktningar[d]), "dim")
		return
	}
	nx, ny := p.X+dirDX[d], p.Y+dirDY[d]
	target := &g.rooms[ny][nx]
	isNew := !target.Walked // ett rum som någon bara har sett in i är fortfarande nytt
	g.reveal(nx, ny)
	target.Walked = true
	p.X, p.Y = nx, ny
	g.noteRoom(p)

	if isNew {
		g.newRoomSound(p)
		if p.Dark > 0 {
			g.logTo(p, g.data.Text("morker_nytt"), "dim")
		} else {
			g.describeLater(p, target, nx, ny)
			g.mirror(p)
		}
		if target.Event != "" {
			ev := target.Event
			target.Event = ""
			t := g.data.trap(ev)
			switch {
			case t != nil && g.is(p, "sparfinnare") && g.rng.Float64() < g.classVal(p, "se_falla", 0):
				// Spårfinnaren ser fällan i tid. Tecknet visar att den finns.
				target.Mark = t.Tecken
				text := g.data.Text("klass_sparfinnare_falla", "falla", t.Kort)
				g.logTo(p, text, "info")
				g.trapCard(p, t, text, true)
			case t != nil && p.initiated && g.rng.Float64() < g.data.Avslojande.Invigd.Falla:
				// Den invigde känner igen fällan för vad den är.
				target.Mark = t.Tecken
				text := g.data.Text("invigd_falla", "falla", t.Kort)
				g.logTo(p, text, "rost")
				g.trapCard(p, t, text, true)
			default:
				g.trigger(p, g.data.trap(ev), depth)
			}
		}
	} else {
		g.soundTo(p, "steg")
		g.noticeTrail(p) // i ett nytt rum har ingen gått förut
	}
	// En fälla kan ha flyttat spelaren; varelsen gäller rummet spelaren står i nu.
	if p.X == nx && p.Y == ny {
		g.pickUp(p)
		g.announceMonster(p)
		g.hornListen(p)
	}
}

type clientMsg struct {
	T    string `json:"t"`
	Name string `json:"name"`
	A    string `json:"a"`
	M    int    `json:"m"`
}

// serve sköter en websocket-anslutning tills den stängs.
// msgConn är en anslutning till en spelarsida: en websocket när tornet är ett
// program (ws.go), eller en kanal i webbläsaren när tornet körs där (wasm.go).
type msgConn interface {
	ReadMessage() ([]byte, error)
	WriteText([]byte) error
	Close() error
}

func (g *Game) serve(conn msgConn, local bool) {
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
		g.seen() // en spelare gör något, eller säger att sidan syns (puls)
		switch m.T {
		case "join":
			g.Join(c, m.Name)
		case "act":
			g.Act(c, m.A, m.M)
		case "ready":
			g.Ready(c)
		case "klass":
			g.ChooseClass(c, m.Name)
		case "viska":
			g.Whisper(c, m.Name)
			g.broadcastState() // invigd, nåden och om det går att viska igen syns direkt
		}
		g.mu.Unlock()
	}
	g.mu.Lock()
	g.Leave(c)
	g.mu.Unlock()
	conn.Close()
}

// aiTexts väljer de vaga texterna om datorspelare, som inte säger vid vilken
// port de kom, om de finns.
func (g *Game) aiTexts(list, hidden []string) []string {
	if len(hidden) > 0 {
		return hidden
	}
	return list
}

// fitTower håller tornet i lobbyn lagom stort och ger alla en egen port
// när någon kommer eller går. Under en vandring får den som kommer sin port
// utan att någon annan flyttas.
func (g *Game) fitTower(p *Player) {
	if g.phase == phaseLobby {
		g.newTower()
		return
	}
	if p != nil {
		used := map[int]bool{}
		for _, o := range g.players {
			if o != nil && o != p {
				used[o.Corner] = true
			}
		}
		for i := range g.corners() {
			if !used[i] {
				p.Corner = i
				break
			}
		}
		g.placeAtCorner(p)
	}
}

// newRoomSound spelar ljudet för ett nytt rum, bara för den som gick in.
func (g *Game) newRoomSound(p *Player) {
	g.soundTo(p, "nyttrum")
}

// visibilityKey är texten om synligheten när vandringen börjar. Kartritaren
// ser alla rum som någon har hittat, så hen får inte höra att kartan bara
// visar de egna rummen.
func (g *Game) visibilityKey(p *Player) string {
	if g.is(p, "kartritare") {
		return "synlighet_kartritare"
	}
	return "synlighet_egen"
}
