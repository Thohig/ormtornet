//go:build !js

package main

import (
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

//go:embed web
var webFS embed.FS

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "adress och port att lyssna på")
	dataPath := flag.String("data", "data", "mapp med rum, fällor och texter")
	sun := flag.Int("sol", 0, "antal rundor innan solen går ner (0 = automatiskt: 40, 50 eller 60 efter tornets storlek)")
	turn := flag.Duration("turtid", 0, "längsta tid att välja drag per runda (0 = automatiskt: 30–40 sekunder efter antalet människor)")
	minPlayers := flag.Int("min", defaultMin, "minsta antal spelare för att börja, datorspelare inräknade (för att prova; förval 2)")
	bots := flag.Int("ai", defaultSettings.AI, "antal datorstyrda vandrare (0-7)")
	seed := flag.Uint64("fro", 0, "frö för första vandringen; samma frö och samma drag ger samma torn (0 = slumpat)")
	difficulty := flag.String("svarighet", "", "svårighet: latt, normal, svar eller mardrom (förvalen står i data/svarighet.json)")
	settingsPath := flag.String("installningar", defaultSettingsPath(), "fil där spelledarens inställningar sparas (tom för att inte spara)")
	openPage := flag.Bool("oppna", false, "öppna introt i webbläsaren när servern har startat; det leder till värdsidan")
	awake := flag.Bool("vaken", false, "håll datorn vaken så länge servern är igång")
	autostop := flag.Duration("autostopp", 10*time.Minute, "stäng tornet när ingen spelare är ansluten och ingen sida har synts på den här datorn så här länge (0 = aldrig)")
	abandoned := flag.Duration("autostopp-max", 12*time.Hour, "stäng tornet när ingen har gjort något så här länge, även om flikar står öppna (0 = aldrig)")
	simulate := flag.Int("simulera", 0, "spela så många vandringar med bara AI-spelare (-ai, förval 4) och skriv hur det gick, utan att öppna tornet")
	check := flag.Bool("kontrollera", false, "kontrollera datafilerna och avsluta")
	site := flag.String("sida", "", "skriv sidorna för GitHub Pages till mappen och avsluta (verktyg/bygg-sida.sh lägger till tornet)")
	lang := flag.String("sprak", defaultSettings.Sprak, "språket i spelet och på värdsidan: sv eller en (texter som inte är översatta visas på svenska)")
	flag.Parse()

	// Spelledarens sparade inställningar gäller, utom det som anges med
	// flaggor: de gäller bara den här gången. Solen, turtiden och minsta
	// antalet är inga inställningar; flaggorna finns för den som provar.
	settings, err := loadSettings(*settingsPath)
	if err != nil {
		fmt.Printf("Kunde inte läsa sparade inställningar (%v). Använder standard.\n", err)
	}
	cfg := newConfig()
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "sol":
			cfg.SunMax = *sun
		case "turtid":
			if t := turn.Round(time.Second); t > 0 {
				cfg.TurnTime = max(3*time.Second, min(10*time.Minute, t))
			}
		case "min":
			cfg.MinPlayers = max(1, min(maxPlayers, *minPlayers))
		case "ai":
			settings.AI = *bots
		case "svarighet":
			settings.Svarighet = *difficulty
		case "sprak":
			settings.Sprak = *lang
		}
	})
	cfg = settings.applyTo(cfg)
	cfg.Seed = *seed

	data, src, err := loadDataLang(*dataPath, cfg.Lang)
	if err != nil {
		log.Fatalf("kunde inte läsa datafilerna: %v", err)
	}
	if *check {
		fmt.Println("Datafilerna är i ordning:", src)
		for _, l := range languages[1:] {
			if err := printTranslationStatus(*dataPath, l); err != nil {
				log.Fatalf("översättningen till %s: %v", l, err)
			}
		}
		return
	}
	if *site != "" {
		sv, _, err := loadDataLang(*dataPath, langSv)
		if err == nil {
			err = writeSite(*site, sv)
		}
		if err != nil {
			log.Fatalf("sidorna: %v", err)
		}
		fmt.Println("Sidorna är skrivna till", *site)
		return
	}
	if *simulate > 0 {
		if err := runSimulation(data, cfg, *simulate, cfg.AI, *seed, os.Stdout); err != nil {
			log.Fatalf("simuleringen: %v", err)
		}
		return
	}
	if cfg.Difficulty != "" && data.difficultyID(cfg.Difficulty) != cfg.Difficulty {
		fmt.Printf("Svårigheten %q finns inte i svarighet.json. Använder standard.\n", cfg.Difficulty)
	}
	game := NewGame(data, cfg)
	game.syncBots(false)
	addedBots := game.botCount()

	sub, _ := fs.Sub(webFS, "web")
	files := http.FileServer(http.FS(sub))
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		files.ServeHTTP(w, r)
	}))
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrade(w, r)
		if err != nil {
			return
		}
		game.serve(conn, isLocal(r))
	})

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		_, port, _ := net.SplitHostPort(*addr)
		if errors.Is(err, syscall.EADDRINUSE) && towerRunning(port) {
			fmt.Println("Ormtornet är redan igång på den här datorn.")
			fmt.Printf("Värdsidan: http://localhost:%s/spelledare\n", port)
			if *openPage {
				openBrowser("http://localhost:" + port + "/")
			}
			return
		}
		if errors.Is(err, syscall.EADDRINUSE) {
			log.Fatalf("Port %s används redan av ett annat program. Prova en annan, till exempel -addr 0.0.0.0:8081", port)
		}
		log.Fatalf("kan inte lyssna på %s: %v", *addr, err)
	}
	_, port, _ := net.SplitHostPort(ln.Addr().String())

	// Tornet stänger på samma sätt från knappen på värdsidan och från
	// Ctrl+C i terminalen: spelarna får ett meddelande och sedan avslutas servern.
	closeTower := sync.OnceFunc(func() {
		fmt.Println()
		fmt.Println("Tornet stänger ...")
		game.Close()
		time.Sleep(1500 * time.Millisecond) // låt meddelandet nå fram
		fmt.Println("Tornet är stängt.")
		os.Exit(0)
	})
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	if !signal.Ignored(syscall.SIGHUP) {
		// Startad med nohup (som från Ormtornet.app) ska den inte börja lyssna på SIGHUP.
		signal.Notify(sigs, syscall.SIGHUP)
	}
	go func() {
		<-sigs
		go closeTower()
		<-sigs // ett andra Ctrl+C stänger direkt
		os.Exit(1)
	}()
	addrs := hostAddrs{Name: "http://" + localName() + ":" + port, IPs: []string{}}
	for _, ip := range lanIPs() {
		addrs.IPs = append(addrs.IPs, "http://"+ip+":"+port)
	}
	addHostRoutes(mux, sub, game, addrs, hostHooks{
		Close: closeTower,
		Save:  func(s Settings) error { return saveSettings(*settingsPath, s) },
		Load: func(lang string) (*Data, error) {
			d, _, err := loadDataLang(*dataPath, lang)
			return d, err
		},
	})

	fmt.Println("Ormtornet är vaket.")
	fmt.Println("Texter från:", src)
	sunText, turnText := "automatiskt efter tornets storlek", "automatisk efter antalet spelare"
	if cfg.SunMax > 0 {
		sunText = fmt.Sprintf("efter %d rundor", cfg.SunMax)
	}
	if cfg.TurnTime > 0 {
		turnText = cfg.TurnTime.String()
	}
	fmt.Printf("Solen går ner %s. Turtid %s. Minst %d spelare.\n", sunText, turnText, cfg.MinPlayers)
	if s := data.difficulty(cfg.Difficulty); s != nil {
		fmt.Printf("Svårighet: %s.\n", s.Namn)
	}
	if addedBots > 0 {
		fmt.Printf("%d datorstyrda vandrare väntar vid portarna.\n", addedBots)
	}
	fmt.Println()
	fmt.Println("Spelarna öppnar:")
	fmt.Printf("  %s\n", addrs.Name)
	for _, ip := range addrs.IPs {
		fmt.Printf("  %s  (om namnet inte fungerar)\n", ip)
	}
	fmt.Println()
	fmt.Println("Värdsidan (bara på den här datorn):")
	fmt.Printf("  http://localhost:%s/spelledare\n", port)
	fmt.Println()
	if *awake {
		// caffeinate släpper av sig själv när servern avslutas.
		if err := exec.Command("caffeinate", "-i", "-w", strconv.Itoa(os.Getpid())).Start(); err == nil {
			fmt.Println("Datorn hålls vaken så länge tornet är öppet.")
		}
	}
	fmt.Println("Stäng tornet med Ctrl+C eller knappen på värdsidan.")
	if *autostop > 0 || *abandoned > 0 {
		if *autostop > 0 {
			fmt.Printf("Tornet stänger av sig självt när det har stått tomt i %s.\n", short(*autostop))
		}
		if *abandoned > 0 {
			fmt.Printf("Senast när ingen har gjort något på %s.\n", short(*abandoned))
		}
		every := 15 * time.Second
		for _, d := range []time.Duration{*autostop, *abandoned} {
			if d > 0 {
				every = min(every, d/4)
			}
		}
		go func() {
			for range time.Tick(every) {
				if game.ShouldClose(*autostop, *abandoned) {
					fmt.Print("\nIngen är kvar vid tornet.")
					closeTower()
				}
			}
		}()
	}
	fmt.Println()
	if *openPage {
		openBrowser("http://localhost:" + port + "/")
	}
	log.Fatal(http.Serve(ln, mux))
}

// short skriver en tid kort, som 10m eller 12h, utan nollor på slutet.
func short(d time.Duration) string {
	s := d.String()
	for _, z := range []string{"0s", "0m"} {
		if len(s) > 2 && s[len(s)-2:] == z && (s[len(s)-3] < '0' || s[len(s)-3] > '9') {
			s = s[:len(s)-2]
		}
	}
	return s
}

// towerRunning är sant om det redan är ett Ormtornet som svarar på porten.
func towerRunning(port string) bool {
	c := http.Client{Timeout: 2 * time.Second}
	resp, err := c.Get("http://127.0.0.1:" + port + "/spelledare/status")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	var v hostView
	return resp.StatusCode == http.StatusOK && json.NewDecoder(resp.Body).Decode(&v) == nil && v.MaxPlayers > 0
}

func openBrowser(url string) {
	if err := exec.Command("open", url).Start(); err != nil {
		fmt.Println("Kunde inte öppna webbläsaren. Gå till", url)
	}
}

// localName ger Macens Bonjour-namn, t.ex. "Din-Mac.local".
func localName() string {
	if out, err := exec.Command("scutil", "--get", "LocalHostName").Output(); err == nil {
		if n := strings.TrimSpace(string(out)); n != "" {
			return n + ".local"
		}
	}
	h, _ := os.Hostname()
	if h == "" {
		return "localhost"
	}
	if !strings.HasSuffix(h, ".local") {
		h += ".local"
	}
	return h
}

func lanIPs() []string {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, a := range addrs {
		if n, ok := a.(*net.IPNet); ok && !n.IP.IsLoopback() && n.IP.To4() != nil && n.IP.IsPrivate() {
			ips = append(ips, n.IP.String())
		}
	}
	return ips
}
