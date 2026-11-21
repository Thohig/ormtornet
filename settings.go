package main

// Inställningar som spelledaren kan ändra i webbläsaren. De sparas i en fil
// så att de finns kvar till nästa gång tornet startas.

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Settings struct {
	AI         int    `json:"ai"`            // datorspelare, 0–7
	Avslojande bool   `json:"avslojande"`    // låt tornet avslöja sin hemlighet; då gäller också invigningen och tornet förändras
	Franvaro   int    `json:"franvaro"`      // rundor innan AI tar över en frånvarande vandrare, 0 = aldrig
	Spelledare bool   `json:"spelledarlage"` // spelledaren spelar inte själv och ser allt
	Svarighet  string `json:"svarighet"`     // förval i data/svarighet.json, tomt för standard
	Sprak      string `json:"sprak"`         // språket i spelet och på värdsidan: "sv" eller "en"
}

// Det som inte går att ställa in. Tornets storlek, solen och turtiden
// bestäms automatiskt; klasserna och striden mellan vandrarna är alltid
// med. Flaggorna -sol, -turtid och -min finns kvar för den som provar.
const defaultMin = 2

var defaultSettings = Settings{AI: 0, Avslojande: true, Franvaro: 2, Sprak: langSv}

var errPlaying = errors.New("inställningarna kan ändras när omgången är slut")

// normalize håller värdena inom det spelet klarar av.
func (s Settings) normalize() Settings {
	// Utan spelledarläge måste en plats finnas kvar åt en människa. I
	// spelledarläget kan alla vandrarna vara AI-spelare.
	maxAI := maxPlayers - 1
	if s.Spelledare {
		maxAI = maxPlayers
	}
	s.AI = max(0, min(maxAI, s.AI))
	s.Franvaro = max(0, min(10, s.Franvaro))
	s.Sprak = langOrDefault(s.Sprak)
	return s
}

// newConfig är spelets regler innan inställningarna läggs på.
func newConfig() Config {
	return Config{Classes: true, PvP: pvpOn, MinPlayers: defaultMin}
}

// applyTo lägger inställningarna på cfg. Det som inte är en inställning,
// som solen och turtiden från flaggorna, får vara som det är.
func (s Settings) applyTo(cfg Config) Config {
	s = s.normalize()
	cfg.Difficulty = s.Svarighet
	cfg.AI = s.AI
	cfg.AbsentRounds = s.Franvaro
	cfg.GMMode = s.Spelledare
	cfg.Reveal = s.Avslojande
	cfg.Lang = s.Sprak
	return cfg
}

func settingsFrom(cfg Config) Settings {
	return Settings{Svarighet: cfg.Difficulty, AI: cfg.AI, Franvaro: cfg.AbsentRounds, Spelledare: cfg.GMMode, Avslojande: cfg.Reveal, Sprak: langOrDefault(cfg.Lang)}
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
	g.cfg = s.applyTo(g.cfg)
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
		g.sunMax = g.scaledSun()
		g.sun = g.sunMax
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
	if g.phase == phaseChoose {
		g.startGame() // klassvalet slutar direkt: den som inte har valt får ett kort
		return true
	}
	if g.phase == phaseEnd {
		g.newTower()
	}
	g.begin()
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
			g.fitTower(nil)
			if announce && len(g.data.AI.Lamnar) > 0 {
				g.logAll(fill(pick(g.aiTexts(g.data.AI.Lamnar, g.data.AI.LamnarDolt)), "namn", p.Name, "horn", g.data.Horn[p.Corner]), "dim")
			}
			return i
		}
	}
	return -1
}
