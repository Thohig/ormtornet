package main

// Läser in rum, fällor och texter från datamappen. Själva innehållet ligger
// i JSON-filerna så att det kan ändras utan att röra koden.

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed data/*.json
var embeddedData embed.FS

// Trap är en fälla. Vad den gör styrs av id:t i koden; siffror och texter
// kommer från datafilen.
type Trap struct {
	ID     string   `json:"id"`
	Namn   string   `json:"namn"`
	Vikt   int      `json:"vikt"`
	Skada  int      `json:"skada"`
	Rundor int      `json:"rundor"`
	Sol    int      `json:"sol"`
	Tecken string   `json:"tecken"`
	Texter []string `json:"texter"`
	Andra  []string `json:"andra"`
	// AndraDolt ersätter Andra när vandrarna bara ser varandra vid möten.
	AndraDolt []string `json:"andra_dolt"`
}

type Rules struct {
	Liv           int     `json:"liv"`
	Handelsechans float64 `json:"handelsechans"`
	DryckHelar    int     `json:"dryck_helar"`
	SokSkattMin   int     `json:"sok_skatt_min"`
	SokSkattMax   int     `json:"sok_skatt_max"`
	SparRundor    int     `json:"spar_rundor"` // hur många rundor spåren efter en vandrare syns
}

type Search struct {
	Vikter map[string]int      `json:"vikter"`
	Texter map[string][]string `json:"texter"`
}

// Serpent styr ormen i skattkammaren och leken med ormkort.
type Serpent struct {
	SkattMin          int      `json:"skatt_min"`
	SkattMax          int      `json:"skatt_max"`
	Skada             int      `json:"skada"`
	SkattForlust      float64  `json:"skatt_forlust"`
	FarreEfterVakning int      `json:"farre_efter_vakning"`
	MinstSovkort      int      `json:"minst_sovkort"`
	Kort              []Card   `json:"kort"`
	VaknarAndra       []string `json:"vaknar_andra"`
	KastasUt          []string `json:"kastas_ut"`
	SomnarOm          []string `json:"somnar_om"`
	Somn              []string `json:"somn"`
}

// Monsters styr varelserna i tornet och striderna mot dem.
type Monsters struct {
	Chans       float64             `json:"chans"`
	Spelarskada int                 `json:"spelarskada"`
	Val         map[string]string   `json:"val"`
	Typer       []MonsterType       `json:"typer"`
	Texter      map[string][]string `json:"texter"`
}

type MonsterType struct {
	ID       string         `json:"id"`
	Namn     string         `json:"namn"`
	Vikt     int            `json:"vikt"`
	Liv      int            `json:"liv"`
	Skada    int            `json:"skada"`
	Guld     int            `json:"guld"`
	Tecken   string         `json:"tecken"`
	Vanor    map[string]int `json:"vanor"`
	DykerUpp []string       `json:"dyker_upp"`
	Dor      []string       `json:"dor"`
}

func (d *Data) monsterType(id string) *MonsterType {
	for i := range d.Monster.Typer {
		if d.Monster.Typer[i].ID == id {
			return &d.Monster.Typer[i]
		}
	}
	return nil
}

// MonsterText väljer en text för striderna.
func (d *Data) MonsterText(key string, vars ...string) string {
	return fill(pick(d.Monster.Texter[key]), vars...)
}

type Card struct {
	ID     string   `json:"id"`
	Antal  int      `json:"antal"`
	Vaknar bool     `json:"vaknar"`
	Texter []string `json:"texter"`
	Andra  []string `json:"andra"`
}

func (d *Data) card(id string) *Card {
	for i := range d.Ormen.Kort {
		if d.Ormen.Kort[i].ID == id {
			return &d.Ormen.Kort[i]
		}
	}
	return nil
}

type RoomType struct {
	ID            string              `json:"id"`
	Namn          string              `json:"namn"`
	Beskrivningar map[string][]string `json:"beskrivningar"`
}

type Data struct {
	Rumstyper    []RoomType          `json:"rumstyper"`
	Ingang       RoomType            `json:"ingang"`
	Skattkammare RoomType            `json:"skattkammare"`
	Horn         []string            `json:"horn"`
	Riktningar   []string            `json:"riktningar"`
	Regioner     [][]string          `json:"regioner"`
	Texter       map[string][]string `json:"texter"`
	Regler       Rules               `json:"regler"`
	Fallor       []Trap              `json:"fallor"`
	Sok          Search              `json:"sok"`
	Ormen        Serpent             `json:"ormen"`
	Monster      Monsters            `json:"monster"`
	Avslojande   Reveal              `json:"avslojande"`
	AI           AIData              `json:"ai"`
	Svarighet    Difficulty          `json:"svarighet"`
	Intro        IntroData           `json:"intro"`
	Klasser      []Class             `json:"klasser"`
}

// IntroData styr startsidan.
type IntroData struct {
	Vanta int `json:"vanta"` // sekunder innan "Gå in i tornet" visas första gången
}

// Filerna i datamappen. Alla läses in i samma struktur; "texter" slås ihop.
var dataFiles = []string{"texter.json", "fallor.json", "ormen.json", "monster.json", "avslojande.json", "ai.json", "svarighet.json", "klasser.json"}

// loadData läser filerna i mappen dir. Saknas en fil används den inbyggda kopian.
func loadData(dir string) (*Data, string, error) {
	d := &Data{Texter: map[string][]string{}}
	var sources []string
	for _, name := range dataFiles {
		path := filepath.Join(dir, name)
		raw, err := os.ReadFile(path)
		src := path
		if errors.Is(err, os.ErrNotExist) {
			raw, err = embeddedData.ReadFile("data/" + name)
			src = name + " (inbyggd)"
		}
		if err != nil {
			return nil, "", err
		}
		var part Data
		if err := json.Unmarshal(raw, &part); err != nil {
			return nil, "", fmt.Errorf("%s: %w", src, err)
		}
		texts := part.Texter
		part.Texter = nil
		mergeData(d, &part)
		for k, v := range texts {
			d.Texter[k] = v
		}
		sources = append(sources, src)
	}
	if len(d.Rumstyper) == 0 {
		return nil, "", fmt.Errorf("inga rumstyper")
	}
	if len(d.Horn) != 4 || len(d.Riktningar) != 4 || len(d.Regioner) != 3 {
		return nil, "", fmt.Errorf("horn och riktningar ska ha 4 ord, regioner 3 rader")
	}
	if d.Regler.Liv <= 0 {
		d.Regler.Liv = 6
	}
	if len(d.Avslojande.Granser) != 3 {
		return nil, "", fmt.Errorf("avslojande.json: granser ska ha tre värden")
	}
	if d.Avslojande.LedtradVarNRunda <= 0 {
		d.Avslojande.LedtradVarNRunda = 4
	}
	if len(d.Monster.Typer) == 0 {
		return nil, "", fmt.Errorf("monster.json: inga monster")
	}
	if d.Monster.Spelarskada <= 0 {
		d.Monster.Spelarskada = 1
	}
	if len(d.Ormen.Kort) == 0 || len(d.Ormen.Somn) == 0 {
		return nil, "", fmt.Errorf("ormen.json: kort och sömntexter saknas")
	}
	return d, strings.Join(sources, ", "), nil
}

// mergeData kopierar de fält som är satta i part till d.
func mergeData(d, part *Data) {
	if len(part.Rumstyper) > 0 {
		d.Rumstyper = part.Rumstyper
	}
	if part.Ingang.ID != "" {
		d.Ingang = part.Ingang
	}
	if part.Skattkammare.ID != "" {
		d.Skattkammare = part.Skattkammare
	}
	if len(part.Horn) > 0 {
		d.Horn = part.Horn
	}
	if len(part.Riktningar) > 0 {
		d.Riktningar = part.Riktningar
	}
	if part.Intro.Vanta > 0 {
		d.Intro = part.Intro
	}
	if len(part.Klasser) > 0 {
		d.Klasser = part.Klasser
	}
	if len(part.Regioner) > 0 {
		d.Regioner = part.Regioner
	}
	if part.Regler.Liv > 0 {
		d.Regler = part.Regler
	}
	if len(part.Fallor) > 0 {
		d.Fallor = part.Fallor
	}
	if len(part.Sok.Vikter) > 0 {
		d.Sok = part.Sok
	}
	if len(part.Ormen.Kort) > 0 {
		d.Ormen = part.Ormen
	}
	if len(part.Monster.Typer) > 0 {
		d.Monster = part.Monster
	}
	if len(part.Avslojande.Granser) > 0 {
		d.Avslojande = part.Avslojande
	}
	if len(part.AI.Personligheter) > 0 {
		d.AI = part.AI
	}
	if len(part.Svarighet.Nivaer) > 0 {
		d.Svarighet = part.Svarighet
	}
}

// weighted väljer en nyckel efter vikt.
func weighted(weights map[string]int) string {
	keys := make([]string, 0, len(weights))
	total := 0
	for k, w := range weights {
		if w > 0 {
			keys = append(keys, k)
			total += w
		}
	}
	if total == 0 {
		return ""
	}
	sort.Strings(keys)
	n := rand.IntN(total)
	for _, k := range keys {
		n -= weights[k]
		if n < 0 {
			return k
		}
	}
	return keys[len(keys)-1]
}

func (d *Data) randomTrap() *Trap {
	w := map[string]int{}
	for _, t := range d.Fallor {
		w[t.ID] = t.Vikt
	}
	return d.trap(weighted(w))
}

func (d *Data) trap(id string) *Trap {
	for i := range d.Fallor {
		if d.Fallor[i].ID == id {
			return &d.Fallor[i]
		}
	}
	return nil
}

// SearchText väljer en text för ett sökresultat.
func (d *Data) SearchText(key string, vars ...string) string {
	return fill(pick(d.Sok.Texter[key]), vars...)
}

func pick(list []string) string {
	if len(list) == 0 {
		return ""
	}
	return list[rand.IntN(len(list))]
}

// fill byter ut {nyckel} mot värden. vars är par: nyckel, värde, nyckel, värde ...
func fill(s string, vars ...string) string {
	for i := 0; i+1 < len(vars); i += 2 {
		s = strings.ReplaceAll(s, "{"+vars[i]+"}", vars[i+1])
	}
	return s
}

// Text väljer en slumpad variant av texten med namnet key.
func (d *Data) Text(key string, vars ...string) string {
	list := d.Texter[key]
	if len(list) == 0 {
		return "[" + key + "]"
	}
	return fill(pick(list), vars...)
}

// Describe väljer en beskrivning för den nivå tornet befinner sig på.
// variant väljer bland texterna, så att olika spelare kan få olika texter
// om samma rum. Finns ingen text för nivån används närmaste lägre nivå.
func (rt *RoomType) Describe(level, variant int) string {
	for l := level; l >= 1; l-- {
		if list := rt.Beskrivningar[fmt.Sprint(l)]; len(list) > 0 {
			return list[((variant%len(list))+len(list))%len(list)]
		}
	}
	return ""
}

// AIData beskriver de datorstyrda vandrarna.
type AIData struct {
	Namn           []string      `json:"namn"`
	Personligheter []Personality `json:"personligheter"`
	BetanketidMs   []int         `json:"betanketid_ms"`
	Ansluter       []string      `json:"ansluter"`
	Lamnar         []string      `json:"lamnar"`
	AnsluterDolt   []string      `json:"ansluter_dolt"` // utan port, när vandrarna inte ser varandra
	LamnarDolt     []string      `json:"lamnar_dolt"`
	TarOver        []string      `json:"tar_over"`
	Tillbaka       []string      `json:"tillbaka"`
}

// Personality styr hur modig en datorspelare är.
type Personality struct {
	ID          string    `json:"id"`
	Namn        string    `json:"namn"`
	Marginal    int       `json:"marginal"`       // extra rundor att ha kvar när den går mot porten
	SkattMal    int       `json:"skatt_mal"`      // går ut när den bär så här mycket
	TaVidSomn   []float64 `json:"ta_vid_somn"`    // chans att ta skatt per sömnnivå
	SokChans    float64   `json:"sok_chans"`      // chans att söka igenom ett rum
	FlyMarginal int       `json:"fly_marginal"`   // flyr när ett slag till (plus så här mycket) kan döda
	GaHemVidLiv int       `json:"ga_hem_vid_liv"` // går mot porten vid så här lite liv
	Klasser     []string  `json:"klasser"`        // klasserna personligheten helst väljer, i tur och ordning
	DrickVidLiv int       `json:"drick_vid_liv"`  // dricker vid så här lite liv
	Stridslust  float64   `json:"stridslust"`     // chans att slåss med en vandrare i samma rum
}

// Reveal styr hur tornet avslöjas i fyra nivåer.
type Reveal struct {
	Granser          []float64           `json:"granser"`
	LedtradVarNRunda int                 `json:"ledtrad_var_n_runda"`
	Nivabyte         map[string][]string `json:"nivabyte"`
	Ledtradar        map[string][]string `json:"ledtradar"`
	Forklaringar     map[string]string   `json:"forklaringar"`
	ForklaringRubrik string              `json:"forklaring_rubrik"`
	Epilog           []string            `json:"epilog"`
	EpilogUte        []string            `json:"epilog_ute"`
	EpilogForlorad   []string            `json:"epilog_forlorad"`
	EpilogDod        []string            `json:"epilog_dod"`
}

func (d *Data) roomType(id string) *RoomType {
	switch id {
	case "ingang":
		return &d.Ingang
	case "skattkammare":
		return &d.Skattkammare
	}
	for i := range d.Rumstyper {
		if d.Rumstyper[i].ID == id {
			return &d.Rumstyper[i]
		}
	}
	return &d.Rumstyper[0]
}

func (d *Data) region(x, y int) string {
	part := func(v int) int {
		switch {
		case v < gridSize/3:
			return 0
		case v >= gridSize-gridSize/3:
			return 2
		}
		return 1
	}
	row := d.Regioner[part(y)]
	if len(row) != 3 {
		return ""
	}
	return row[part(x)]
}
