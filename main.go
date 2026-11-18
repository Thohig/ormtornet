package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

//go:embed web
var webFS embed.FS

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "adress och port att lyssna på")
	dataPath := flag.String("data", "data/texter.json", "fil med rum och texter")
	flag.Parse()

	data, src, err := loadData(*dataPath)
	if err != nil {
		log.Fatalf("kunde inte läsa texterna: %v", err)
	}
	game := NewGame(data)

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
		game.serve(conn)
	})

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("kan inte lyssna på %s: %v", *addr, err)
	}
	_, port, _ := net.SplitHostPort(ln.Addr().String())

	fmt.Println("Ormtornet är vaket.")
	fmt.Println("Texter från:", src)
	fmt.Println()
	fmt.Println("Spelarna öppnar:")
	fmt.Printf("  http://%s:%s\n", localName(), port)
	for _, ip := range lanIPs() {
		fmt.Printf("  http://%s:%s  (om namnet inte fungerar)\n", ip, port)
	}
	fmt.Println()
	log.Fatal(http.Serve(ln, mux))
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
