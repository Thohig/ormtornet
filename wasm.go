//go:build js && wasm

package main

// Tornet i webbläsaren (docs/github-pages.md). Samma spel som programmet i
// main.go, men utan nätverk och filer: skalet web/torn.html laddar
// torn.wasm, och sidorna i webbläsaren når tornet genom skalet.
//
// Go ger skalet det globala objektet Tornet:
//
//	Tornet.fraga(metod, vag, huvud, kropp, klar, lokal)
//	                                              en förfrågan som /intro
//	                                              eller /spelledare/status;
//	                                              klar(status, typ, kropp).
//	                                              lokal är sant när den
//	                                              kommer från spelledarens
//	                                              egen webbläsare
//	Tornet.anslut(id, lokal)                      en spelarsida ansluter
//	Tornet.taEmot(id, text)                       spelarsidan skickar
//	Tornet.kopplaNer(id)                          spelarsidan är borta
//
// och anropar skalets tornetSkickar(id, text), tornetKopplarNer(id),
// tornetRedo() och tornetStangt().

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"syscall/js"
	"testing/fstest"
	"time"
)

const settingsKey = "ormtornet-installningar"

func main() {
	settings := loadSettingsJS()
	cfg := settings.config()
	data, _, err := loadDataLang("", cfg.Lang)
	if err != nil {
		js.Global().Call("tornetFel", err.Error())
		return
	}
	game := NewGame(data, cfg)
	game.syncBots(false)

	var closed bool
	var closedMu sync.Mutex
	isClosed := func() bool {
		closedMu.Lock()
		defer closedMu.Unlock()
		return closed
	}
	closeTower := sync.OnceFunc(func() {
		game.Close()
		time.Sleep(1500 * time.Millisecond) // låt meddelandet nå fram
		closedMu.Lock()
		closed = true
		closedMu.Unlock()
		js.Global().Call("tornetStangt")
	})

	addrs := hostAddrs{Name: js.Global().Get("ORMTORNET_ADRESS").String(), IPs: []string{}}
	mux := http.NewServeMux()
	// Sidorna själva kommer från Pages; tornet behöver inga egna filer.
	addHostRoutes(mux, fstest.MapFS{}, game, addrs, hostHooks{
		Close: closeTower,
		Save:  saveSettingsJS,
		Load: func(lang string) (*Data, error) {
			d, _, err := loadDataLang("", lang)
			return d, err
		},
	})

	var connsMu sync.Mutex
	conns := map[string]*pageConn{}

	api := js.Global().Get("Object").New()
	api.Set("fraga", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 5 {
			return nil
		}
		method, path, body, done := args[0].String(), args[1].String(), args[3].String(), args[4]
		local := len(args) > 5 && args[5].Truthy()
		header := http.Header{}
		if h := args[2]; h.Type() == js.TypeObject {
			keys := js.Global().Get("Object").Call("keys", h)
			for i := 0; i < keys.Length(); i++ {
				k := keys.Index(i).String()
				header.Set(k, h.Get(k).String())
			}
		}
		go func() {
			if isClosed() {
				done.Invoke(503, "text/plain", "")
				return
			}
			status, ctype, out := serveRequest(mux, method, path, header, body, local)
			done.Invoke(status, ctype, out)
		}()
		return nil
	}))
	api.Set("anslut", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 2 || isClosed() {
			return false
		}
		id, local := args[0].String(), args[1].Truthy()
		c := newPageConn(id)
		connsMu.Lock()
		if old := conns[id]; old != nil {
			old.Close()
		}
		conns[id] = c
		connsMu.Unlock()
		go func() {
			game.serve(c, local)
			connsMu.Lock()
			if conns[id] == c {
				delete(conns, id)
			}
			connsMu.Unlock()
		}()
		return true
	}))
	api.Set("taEmot", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 2 {
			return nil
		}
		connsMu.Lock()
		c := conns[args[0].String()]
		connsMu.Unlock()
		if c != nil {
			c.receive([]byte(args[1].String()))
		}
		return nil
	}))
	api.Set("kopplaNer", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return nil
		}
		connsMu.Lock()
		c := conns[args[0].String()]
		connsMu.Unlock()
		if c != nil {
			c.Close()
		}
		return nil
	}))
	js.Global().Set("Tornet", api)
	js.Global().Call("tornetRedo")
	select {}
}

// serveRequest låter mux svara på en förfrågan från en sida i webbläsaren.
// Förfrågningar från spelledarens egen webbläsare (local) räknas som från
// tornets dator (isLocal). De som kommer över WebRTC från spelarna får en
// adress som inte är det.
func serveRequest(mux *http.ServeMux, method, path string, header http.Header, body string, local bool) (int, string, string) {
	r, err := http.NewRequest(method, "http://tornet"+path, strings.NewReader(body))
	if err != nil {
		return http.StatusBadRequest, "text/plain", err.Error()
	}
	r.Header = header
	r.RemoteAddr = "192.0.2.1:1" // en adress som bara finns i exempel (RFC 5737)
	if local {
		r.RemoteAddr = "127.0.0.1:1"
	}
	w := &memResponse{header: http.Header{}, status: http.StatusOK}
	mux.ServeHTTP(w, r)
	return w.status, w.header.Get("Content-Type"), w.body.String()
}

// memResponse samlar ett svar i minnet.
type memResponse struct {
	header http.Header
	status int
	wrote  bool
	body   bytes.Buffer
}

func (w *memResponse) Header() http.Header { return w.header }

func (w *memResponse) WriteHeader(code int) {
	if !w.wrote {
		w.status, w.wrote = code, true
	}
}

func (w *memResponse) Write(b []byte) (int, error) {
	w.wrote = true
	return w.body.Write(b)
}

// pageConn är en spelarsidas anslutning genom skalet.
type pageConn struct {
	id     string
	inbox  chan []byte
	done   chan struct{}
	closer sync.Once
}

func newPageConn(id string) *pageConn {
	return &pageConn{id: id, inbox: make(chan []byte, 64), done: make(chan struct{})}
}

// receive tar emot ett meddelande från sidan. Den anropas från JavaScript och
// får inte vänta; hänger tornet inte med stängs anslutningen, som en
// websocket som inte hinner läsas.
func (c *pageConn) receive(b []byte) {
	if len(b) > wsMaxMessage {
		c.Close()
		return
	}
	select {
	case <-c.done:
	case c.inbox <- b:
	default:
		c.Close()
	}
}

func (c *pageConn) ReadMessage() ([]byte, error) {
	select {
	case b := <-c.inbox:
		return b, nil
	case <-c.done:
		return nil, io.EOF
	}
}

func (c *pageConn) WriteText(b []byte) error {
	select {
	case <-c.done:
		return errors.New("anslutningen är stängd")
	default:
	}
	js.Global().Call("tornetSkickar", c.id, string(b))
	return nil
}

func (c *pageConn) Close() error {
	c.closer.Do(func() {
		close(c.done)
		js.Global().Call("tornetKopplarNer", c.id)
	})
	return nil
}

// Inställningarna sparas i webbläsaren i stället för i en fil.
func loadSettingsJS() Settings {
	s := defaultSettings
	v := localStorage("getItem", settingsKey)
	if v.Type() != js.TypeString {
		return s
	}
	if json.Unmarshal([]byte(v.String()), &s) != nil {
		return defaultSettings
	}
	return s.normalize()
}

func saveSettingsJS(s Settings) error {
	b, _ := json.Marshal(s)
	if v := localStorage("setItem", settingsKey, string(b)); v.Type() == js.TypeString && v.String() == "fel" {
		return errors.New("webbläsaren sparar inte")
	}
	return nil
}

// localStorage anropar webbläsarens localStorage, som kan saknas eller vara
// avstängd. Svarar strängen "fel" när anropet inte gick.
func localStorage(method string, args ...any) (v js.Value) {
	defer func() {
		if recover() != nil {
			v = js.ValueOf("fel")
		}
	}()
	return js.Global().Get("localStorage").Call(method, args...)
}
