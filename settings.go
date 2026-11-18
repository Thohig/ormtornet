package main

// Inställningar som spelledaren kan ändra i webbläsaren. De sparas i en fil
// så att de finns kvar till nästa gång tornet startas.

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Settings struct {
	AI         int    `json:"ai"`            // datorspelare, 0–3
	Min        int    `json:"min"`           // minsta antal vandrare, datorspelare inräknade
	Sol        int    `json:"sol"`           // rundor innan solen går ner
	Turtid     int    `json:"turtid"`        // sekunder per runda
	Avslojande bool   `json:"avslojande"`    // låt tornet avslöja sin hemlighet
	Franvaro   int    `json:"franvaro"`      // rundor innan AI tar över en frånvarande vandrare, 0 = aldrig
	Strid      string `json:"strid"`         // strid mellan vandrare: "av", "pa" eller "halva"
	Spelledare bool   `json:"spelledarlage"` // spelledaren spelar inte själv och ser allt
	Svarighet  string `json:"svarighet"`     // förval i data/svarighet.json, tomt för standard
	Synlighet  string `json:"synlighet"`     // vilka vandrare och rum man ser: "alla", "moten" eller "egen"
	Klasser    bool   `json:"klasser"`       // vandrarna väljer klass
	Storlek    int    `json:"storlek"`       // tornets storlek i rum: 9, 11 eller 13, 0 för automatisk
}

const (
	sightAll  = "alla"
	sightMeet = "moten" // andra vandrare syns bara i samma rum
	sightOwn  = "egen"  // som "moten", och var och en ser bara rum den själv har varit i
)

var defaultSettings = Settings{AI: 0, Min: 2, Sol: 40, Turtid: 30, Avslojande: true, Franvaro: 2, Strid: pvpOn, Synlighet: sightOwn, Klasser: true}

var errPlaying = errors.New("inställningarna kan ändras när omgången är slut")

// normalize håller värdena inom det spelet klarar av.
func (s Settings) normalize() Settings {
	// Utan spelledarläge måste en plats finnas kvar åt en människa. I
	// spelledarläget kan alla fyra vandrarna vara AI-spelare.
	maxAI := maxPlayers - 1
	if s.Spelledare {
		maxAI = maxPlayers
	}
	s.AI = max(0, min(maxAI, s.AI))
	s.Min = max(1, min(maxPlayers, s.Min))
	s.Sol = max(4, min(200, s.Sol))
	s.Turtid = max(3, min(600, s.Turtid))
	s.Franvaro = max(0, min(10, s.Franvaro))
	if s.Strid != pvpOff && s.Strid != pvpOn && s.Strid != pvpHalf {
		s.Strid = pvpOn
	}
	if s.Storlek != 9 && s.Storlek != 11 && s.Storlek != 13 {
		s.Storlek = 0
	}
	if s.Synlighet != sightAll && s.Synlighet != sightMeet && s.Synlighet != sightOwn {
		s.Synlighet = sightOwn
	}
	return s
}

func (s Settings) config() Config {
	s = s.normalize()
	return Config{Difficulty: s.Svarighet, AI: s.AI, AbsentRounds: s.Franvaro, PvP: s.Strid, Sight: s.Synlighet, Classes: s.Klasser, Size: s.Storlek, GMMode: s.Spelledare, MinPlayers: s.Min, SunMax: s.Sol, TurnTime: time.Duration(s.Turtid) * time.Second, Reveal: s.Avslojande}
}

func settingsFrom(cfg Config) Settings {
	return Settings{Svarighet: cfg.Difficulty, AI: cfg.AI, Franvaro: cfg.AbsentRounds, Strid: cfg.PvP, Synlighet: cfg.Sight, Klasser: cfg.Classes, Storlek: cfg.Size, Spelledare: cfg.GMMode, Min: cfg.MinPlayers, Sol: cfg.SunMax, Turtid: int(cfg.TurnTime / time.Second), Avslojande: cfg.Reveal}
}

// defaultSettingsPath ger ~/Library/Application Support/Ormtornet/installningar.json.
func defaultSettingsPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "Ormtornet", "installningar.json")
}

// loadSettings läser sparade inställningar. Saknas filen används standard.
func loadSettings(path string) (Settings, error) {
	s := defaultSettings
	if path == "" {
		return s, nil
	}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return s, err
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return defaultSettings, err
	}
	return s.normalize(), nil
}

func saveSettings(path string, s Settings) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Settings ger tornets nuvarande inställningar.
func (g *Game) Settings() Settings {
	g.mu.Lock()
	defer g.mu.Unlock()
	return settingsFrom(g.cfg)
}

// ApplySettings byter inställningar mellan omgångarna. Datorspelare kommer
// och går direkt; resten gäller från nästa omgång.
func (g *Game) ApplySettings(s Settings) (Settings, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.phase == phasePlay {
		return settingsFrom(g.cfg), errPlaying
	}
	wasGM := g.cfg.GMMode
	g.cfg = s.config()
	g.cfg.Difficulty = g.base.difficultyID(g.cfg.Difficulty)
	g.data = g.base.withDifficulty(g.cfg.Difficulty)
	if g.phase == phaseLobby {
		for _, p := range g.players {
			if p != nil {
				p.HP = g.maxHP(p)
			}
		}
	}
	if g.cfg.GMMode != wasGM {
		if g.cfg.GMMode {
			g.kickLocal()
		}
		g.gmAnnounce(g.cfg.GMMode)
	}
	if g.phase == phaseLobby {
		g.sun = g.cfg.SunMax
	}
	if g.cfg.Classes && g.phase != phasePlay {
		g.offerAll()
	}
	if g.phase == phaseLobby && g.towerSize() != g.size {
		g.newTower() // ett torn av den nya storleken
	}
	g.syncBots(true)
	if !g.checkStart() {
		g.broadcastState()
	}
	return settingsFrom(g.cfg), nil
}

// ForceStart låter spelledaren börja utan att vänta på att alla är redo
// eller att tillräckligt många har kommit. Minst en människa måste vara med,
// utom i spelledarläget: där räcker det med AI-spelare.
func (g *Game) ForceStart() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.phase == phasePlay {
		return false
	}
	if g.humansPresent() == 0 && !(g.cfg.GMMode && g.botCount() > 0) {
		return false
	}
	if g.phase == phaseEnd {
		g.newTower()
	}
	g.startGame()
	return true
}

func (g *Game) humansPresent() int {
	n := 0
	for _, p := range g.players {
		if p != nil && p.bot == nil && p.present() {
			n++
		}
	}
	return n
}

func (g *Game) botCount() int {
	n := 0
	for _, p := range g.players {
		if p != nil && p.bot != nil && !p.takenOver {
			n++
		}
	}
	return n
}

// syncBots lägger till eller tar bort datorspelare så att det blir så många
// som spelledaren vill ha, men aldrig fler än det finns plats för bredvid
// människorna. Görs bara utanför en pågående omgång.
func (g *Game) syncBots(announce bool) {
	if g.phase == phasePlay {
		return
	}
	want := min(g.cfg.AI, maxPlayers-g.humansPresent())
	for g.botCount() > want {
		if g.removeBot(announce) < 0 {
			break
		}
	}
	for n := g.botCount(); n < want; n = g.botCount() {
		before := g.players
		if g.AddBots(1) == 0 {
			break
		}
		if !announce {
			continue
		}
		for i, p := range g.players {
			if p != nil && before[i] == nil {
				g.logAll(fill(pick(g.aiTexts(g.data.AI.Ansluter, g.data.AI.AnsluterDolt)), "namn", p.Name, "horn", g.data.Horn[p.Corner]), "info")
				g.soundAll("ansluter")
			}
		}
	}
}

// removeBot tar bort den datorspelare som kom in sist och ger tillbaka dess
// plats, eller -1 om det inte finns någon.
func (g *Game) removeBot(announce bool) int {
	if g.phase == phasePlay {
		return -1
	}
	for i, p := range g.players {
		if p != nil && p.bot != nil {
			g.players[i] = nil
			if announce && len(g.data.AI.Lamnar) > 0 {
				g.logAll(fill(pick(g.aiTexts(g.data.AI.Lamnar, g.data.AI.LamnarDolt)), "namn", p.Name, "horn", g.data.Horn[p.Corner]), "dim")
			}
			return i
		}
	}
	return -1
}
