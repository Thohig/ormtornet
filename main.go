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
	sun := flag.Int("sol", defaultSettings.Sol, "antal rundor innan solen går ner")
	turn := flag.Duration("turtid", time.Duration(defaultSettings.Turtid)*time.Second, "längsta tid att välja drag per runda")
	minPlayers := flag.Int("min", defaultSettings.Min, "minsta antal spelare för att börja, datorspelare inräknade")
	bots := flag.Int("ai", defaultSettings.AI, "antal datorstyrda vandrare (0-3, eller 4 i spelledarläget)")
	difficulty := flag.String("svarighet", "", "svårighet: latt, normal, svar eller mardrom (förvalen står i data/svarighet.json)")
	reveal := flag.Bool("avslojande", defaultSettings.Avslojande, "låt tornet avslöja sin hemlighet (false för den som redan känner till den)")
	sight := flag.String("synlighet", defaultSettings.Synlighet, "vilka vandrare och rum man ser: egen (bara dem man möter och egna rum), moten (bara dem man möter) eller alla")
	classes := flag.Bool("klasser", defaultSettings.Klasser, "vandrarna väljer klass i lobbyn (false för att spela utan klasser)")
	size := flag.Int("storlek", defaultSettings.Storlek, "tornets storlek i rum: 9, 11 eller 13 (0 = automatisk)")
	settingsPath := flag.String("installningar", defaultSettingsPath(), "fil där spelledarens inställningar sparas (tom för att inte spara)")
	openPage := flag.Bool("oppna", false, "öppna introt i webbläsaren när servern har startat; det leder till värdsidan")
	awake := flag.Bool("vaken", false, "håll datorn vaken så länge servern är igång")
	autostop := flag.Duration("autostopp", 10*time.Minute, "stäng tornet när ingen spelare är ansluten och ingen sida har synts på den här datorn så här länge (0 = aldrig)")
	abandoned := flag.Duration("autostopp-max", 12*time.Hour, "stäng tornet när ingen har gjort något så här länge, även om flikar står öppna (0 = aldrig)")
	flag.Parse()

	// Spelledarens sparade inställningar gäller, utom det som anges med
	// flaggor: de gäller bara den här gången.
	settings, err := loadSettings(*settingsPath)
	if err != nil {
		fmt.Printf("Kunde inte läsa sparade inställningar (%v). Använder standard.\n", err)
	}
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "sol":
			settings.Sol = *sun
		case "turtid":
			settings.Turtid = int(turn.Round(time.Second) / time.Second)
		case "min":
			settings.Min = *minPlayers
		case "ai":
			settings.AI = *bots
		case "avslojande":
			settings.Avslojande = *reveal
		case "svarighet":
			settings.Svarighet = *difficulty
		case "synlighet":
			settings.Synlighet = *sight
		case "klasser":
			settings.Klasser = *classes
		case "storlek":
			settings.Storlek = *size
		}
	})
	cfg := settings.config()

	data, src, err := loadData(*dataPath)
	if err != nil {
		log.Fatalf("kunde inte läsa texterna: %v", err)
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
	})

	fmt.Println("Ormtornet är vaket.")
	fmt.Println("Texter från:", src)
	fmt.Printf("Solen går ner efter %d rundor. Turtid %s. Minst %d spelare.\n", cfg.SunMax, cfg.TurnTime, cfg.MinPlayers)
	if s := data.difficulty(cfg.Difficulty); s != nil {
		fmt.Printf("Svårighet: %s.\n", s.Namn)
	}
	if !cfg.Reveal {
		fmt.Println("Avslöjandet är avstängt.")
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

// localName ger Macens Bonjour-namn, t.ex. "Stefans-MacBook.local".
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
