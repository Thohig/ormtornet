package main

// Läser in rum, fällor och texter från datamappen. Själva innehållet ligger
// i JSON-filerna så att det kan ändras utan att röra koden.

import (
	"embed"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed data/*.json data/en/*.json
var embeddedData embed.FS

// ModelLook är färgen och storleken på en modell på kartan snett ovanifrån
// (klient/src/modeller.ts). Formen står i koden; utan modell används
// modellens egen färg och storlek.
type ModelLook struct {
	Farg    string  `json:"farg,omitempty"`    // "#rrggbb": varelsens kropp, eller det klassen bär
	Storlek float64 `json:"storlek,omitempty"` // 1 är modellens egen storlek
}

// Trap är en fälla. Vad den gör styrs av id:t i koden; siffror och texter
// kommer från datafilen.
type Trap struct {
	ID     string   `json:"id"`
	Namn   string   `json:"namn"` // bara på värdsidan; spelarna får inga namn
	Kort   string   `json:"kort"` // det spelarna får i stället: "ett galler av spetsiga järnstänger"
	Vikt   int      `json:"vikt"`
	Skada  int      `json:"skada"`
	Rundor int      `json:"rundor"`
	Sol    int      `json:"sol"`
	Tecken string   `json:"tecken"`
	Texter []string `json:"texter"`
	Andra  []string `json:"andra"`
	// AndraDolt ersätter Andra när vandrarna bara ser varandra vid möten.
	AndraDolt []string `json:"andra_dolt"`
	// SparEfter är det man ser i rummet efter att fällan har slagit till,
	// i rumsblocket (here.go). Fällans namn sägs inte där.
	SparEfter []string `json:"spar_efter"`
}

type Rules struct {
	Liv           int     `json:"liv"`
	Handelsechans float64 `json:"handelsechans"`
	DryckHelar    int     `json:"dryck_helar"`
	SokSkattMin   int     `json:"sok_skatt_min"`
	SokSkattMax   int     `json:"sok_skatt_max"`
	SparRundor    int     `json:"spar_rundor"` // hur många rundor spåren efter en vandrare syns
	// TryggVidPorten: rummen bredvid en port har aldrig en fälla, så att
	// ingen förlorar liv redan på första steget in.
	TryggVidPorten bool `json:"trygg_vid_porten"`
}

type Search struct {
	Vikter map[string]int      `json:"vikter"`
	Texter map[string][]string `json:"texter"`
}

// Serpent styr ormen i skattkammaren och leken med ormkort.
type Serpent struct {
	SkattMin           int      `json:"skatt_min"`
	SkattMax           int      `json:"skatt_max"`
	Skada              int      `json:"skada"`
	SkattForlust       float64  `json:"skatt_forlust"`
	FarreEfterVakning  int      `json:"farre_efter_vakning"`
	MinstSovkort       int      `json:"minst_sovkort"`
	Kort               []Card   `json:"kort"`
	VaknarAndra        []string `json:"vaknar_andra"`
	SovkortPerVandrare int      `json:"sovkort_per_vandrare_over_fyra"` // fler sovkort när fler än fyra vandrar
	KastasUt           []string `json:"kastas_ut"`
	SomnarOm           []string `json:"somnar_om"`
	Somn               []string `json:"somn"`
}

// Monsters styr varelserna i tornet och striderna mot dem.
type Monsters struct {
	Chans       float64             `json:"chans"`
	VandraChans float64             `json:"vandra_chans"` // chans per runda att en varelse vandrar av sig själv
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
	Modell   *ModelLook     `json:"modell,omitempty"` // färg och storlek på kartan snett ovanifrån
	Vanor    map[string]int `json:"vanor"`
	Vandrar  *float64       `json:"vandrar,omitempty"` // egen chans att vandra, annars vandra_chans
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
	Namn   string   `json:"namn"` // kortets namn när det vänds upp, till exempel Ormen sover
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
	Namn          string              `json:"namn"` // bara på värdsidan; spelarna får inga namn
	Kort          string              `json:"kort"` // det spelarna får i stället: "en smal gång under låga valv"
	Beskrivningar map[string][]string `json:"beskrivningar"`
	Modell        *RoomLook           `json:"modell,omitempty"`
}

// RoomLook är hur ett rum ser ut på kartan snett ovanifrån: inredningen,
// med namn som klient/src/inredning.ts känner till (furnitureNames).
type RoomLook struct {
	Inredning []string `json:"inredning,omitempty"`
}

// furnitureNames är inredningen kartan kan rita, i klient/src/inredning.ts.
var furnitureNames = []string{"pelare", "pelare_syd", "brunn", "bank", "spjut", "altare", "trappa", "hylla", "tunnor", "sackar",
	"lyktor", "stenblock", "assja", "stad", "bokhylla", "bord", "stenhog"}

type Data struct {
	Rumstyper        []RoomType          `json:"rumstyper"`
	Ingang           RoomType            `json:"ingang"`
	Skattkammare     RoomType            `json:"skattkammare"`
	Horn             []string            `json:"horn"`
	Riktningar       []string            `json:"riktningar"`
	Regioner         [][]string          `json:"regioner"`
	Texter           map[string][]string `json:"texter"`
	Regler           Rules               `json:"regler"`
	Fallor           []Trap              `json:"fallor"`
	Sok              Search              `json:"sok"`
	Ormen            Serpent             `json:"ormen"`
	Monster          Monsters            `json:"monster"`
	Avslojande       Reveal              `json:"avslojande"`
	AI               AIData              `json:"ai"`
	Svarighet        Difficulty          `json:"svarighet"`
	Intro            IntroData           `json:"intro"`
	Klasser          []Class             `json:"klasser"`
	KlassvalSekunder int                 `json:"klassval_sekunder"` // hur länge klassvalet får ta
	Foremal          ItemData            `json:"foremal"`
	Vandrarnamn      NameData            `json:"vandrarnamn"` // namnen tornet ger vandrarna (namn.json)
	// Spelarsidan och Vardsidan är sidornas egna texter, på nyckel. Sidorna
	// hämtar dem från /texter/ (se web/texter.js).
	Spelarsidan map[string]string `json:"spelarsidan"`
	Vardsidan   map[string]string `json:"vardsidan"`
	Lang        string            `json:"-"` // språket texterna är på
}

// IntroData styr startsidan.
type IntroData struct {
	Vanta int `json:"vanta"` // sekunder innan "Gå in i tornet" visas första gången
}

// Filerna i datamappen. Alla läses in i samma struktur; "texter" slås ihop.
var dataFiles = []string{"texter.json", "fallor.json", "ormen.json", "monster.json", "avslojande.json", "ai.json", "svarighet.json", "klasser.json", "foremal.json", "namn.json", "spelarsidan.json", "vardsidan.json"}

// loadData läser filerna i mappen dir på svenska. Saknas en fil används den
// inbyggda kopian.
func loadData(dir string) (*Data, string, error) {
	return loadDataLang(dir, langSv)
}

// readDataFile läser den svenska datafilen name ur dir, eller den inbyggda
// kopian om den saknas. Med dir tom används bara de inbyggda filerna, som i
// webbläsaren där det inte finns några filer att läsa.
func readDataFile(dir, name string) ([]byte, string, error) {
	if dir == "" {
		raw, err := embeddedData.ReadFile("data/" + name)
		return raw, name + " (inbyggd)", err
	}
	path := filepath.Join(dir, name)
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		raw, err = embeddedData.ReadFile("data/" + name)
		path = name + " (inbyggd)"
	}
	return raw, path, err
}

// loadDataLang läser datafilerna och lägger översättningen till lang över
// dem (se lang.go). Text som inte är översatt står kvar på svenska.
func loadDataLang(dir, lang string) (*Data, string, error) {
	d := &Data{Texter: map[string][]string{}, Lang: langOrDefault(lang)}
	var sources []string
	for _, name := range dataFiles {
		raw, src, err := readDataFile(dir, name)
		if err != nil {
			return nil, "", err
		}
		tr, trSrc, err := readTranslation(dir, langOrDefault(lang), name)
		if err != nil {
			return nil, "", err
		}
		if tr != nil {
			rep, err := checkTranslation(raw, tr)
			if err == nil && len(rep.problems) > 0 {
				err = errors.New(strings.Join(rep.problems, "; "))
			}
			if err == nil {
				raw, err = translateJSON(raw, tr)
			}
			if err != nil {
				return nil, "", fmt.Errorf("%s: %w", trSrc, err)
			}
			src += " + " + trSrc
		}
		var part Data
		if err := decodeStrict(raw, &part); err != nil {
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
	if len(d.Horn) != maxPlayers || len(d.Riktningar) != 4 || len(d.Regioner) != 3 {
		return nil, "", fmt.Errorf("horn ska ha 8 ord (fyra hörn och fyra sidor), riktningar 4 ord och regioner 3 rader")
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
	if len(d.Vandrarnamn.Fornamn) < maxPlayers || len(d.Vandrarnamn.Tillnamn) == 0 {
		return nil, "", fmt.Errorf("namn.json: minst %d förnamn och ett tillnamn behövs", maxPlayers)
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
	if problems := d.validate(); len(problems) > 0 {
		return nil, "", fmt.Errorf("fel i datafilerna:\n  %s", strings.Join(problems, "\n  "))
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
	if len(part.Vandrarnamn.Fornamn) > 0 {
		d.Vandrarnamn.Fornamn = part.Vandrarnamn.Fornamn
	}
	if len(part.Vandrarnamn.Tillnamn) > 0 {
		d.Vandrarnamn.Tillnamn = part.Vandrarnamn.Tillnamn
	}
	if len(part.Klasser) > 0 {
		d.Klasser = part.Klasser
	}
	if part.KlassvalSekunder > 0 {
		d.KlassvalSekunder = part.KlassvalSekunder
	}
	if len(part.Foremal.Lista) > 0 {
		d.Foremal = part.Foremal
	}
	if len(part.Regioner) > 0 {
		d.Regioner = part.Regioner
	}
	if len(part.Spelarsidan) > 0 {
		d.Spelarsidan = part.Spelarsidan
	}
	if len(part.Vardsidan) > 0 {
		d.Vardsidan = part.Vardsidan
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
func weighted(rng *rand.Rand, weights map[string]int) string {
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
	n := rng.IntN(total)
	for _, k := range keys {
		n -= weights[k]
		if n < 0 {
			return k
		}
	}
	return keys[len(keys)-1]
}

func (d *Data) randomTrap(rng *rand.Rand) *Trap {
	w := map[string]int{}
	for _, t := range d.Fallor {
		w[t.ID] = t.Vikt
	}
	return d.trap(weighted(rng, w))
}

func (d *Data) trap(id string) *Trap {
	for i := range d.Fallor {
		if d.Fallor[i].ID == id {
			return &d.Fallor[i]
		}
	}
	return nil
}

// trapByMark ger fällan med tecknet mark, eller nil.
func (d *Data) trapByMark(mark string) *Trap {
	for i := range d.Fallor {
		if mark != "" && d.Fallor[i].Tecken == mark {
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

// HostText är en av värdsidans texter, till exempel ett besked eller ett
// fel, med platshållarna ifyllda.
func (d *Data) HostText(key string, vars ...string) string {
	s, ok := d.Vardsidan[key]
	if !ok {
		return "[" + key + "]"
	}
	return fill(s, vars...)
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
	GaHemVidLiv int       `json:"ga_hem_vid_liv"` // går mot porten vid så här lite liv, räknat på liven i datafilerna
	Klasser     []string  `json:"klasser"`        // klasserna personligheten helst väljer, i tur och ordning
	DrickVidLiv int       `json:"drick_vid_liv"`  // dricker vid så här lite liv, räknat på samma sätt
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
	EpilogInvigd     []string            `json:"epilog_invigd"`   // till den som viskade hemligheten
	Svar             []string            `json:"svar"`            // viskningar som visar att man vet hemligheten
	ViskaFranNiva    int                 `json:"viska_fran_niva"` // första nivån då man kan viska
	Invigd           Initiate            `json:"invigd"`

	// Med invigningen "ormens namn" (se secret.go).
	Namn        [][]string          `json:"namn"`          // namnen, tre delar vardera
	Ordning     []string            `json:"ordning"`       // första, andra och tredje
	Delar       map[string][]string `json:"delar"`         // texterna när en del delas ut, per nivå
	SvarMedNamn []string            `json:"svar_med_namn"` // mallar med {namn} som räknas som svar
	EpilogNamn  []string            `json:"epilog_namn"`   // epilogen säger namnet
	NadAndel    float64             `json:"nad_andel"`     // andelen vakna-kort som tornets nåd tar bort
	Varsel      Omens               `json:"varsel"`        // varslen om rum som ingen har varit i
	Forandring  Change              `json:"forandring"`    // hur tornet förändras med nivåerna
}

// Change är hur tornet förändras med avslöjandets nivåer (se change.go).
type Change struct {
	Fallvikter    map[string]map[string]int `json:"fallvikter"`      // nivå -> fälla -> vikt
	PortStangsVar int                       `json:"port_stangs_var"` // rundor mellan portarna som stängs på nivå 4
}

// Omens är texterna för varslen (se omens.go).
type Omens struct {
	Fallor   map[string][]string `json:"fallor"`   // per fälla
	Varelser map[string][]string `json:"varelser"` // per varelse
	Tomt     []string            `json:"tomt"`     // ett rum utan något
	Avstand  []string            `json:"avstand"`  // två och tre rum bort
}

// Initiate är fördelarna och svagheten för den som har viskat hemligheten.
type Initiate struct {
	Guld        float64 `json:"guld"`         // andel extra guld när den invigde kommer ut
	Falla       float64 `json:"falla"`        // chansen att känna igen en fälla och kliva förbi
	ExtraAnfall int     `json:"extra_anfall"` // varelserna anfaller oftare
	Vacker      float64 `json:"vacker"`       // chansen att en näve drar ett vakna-kort
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

// region ger delen av tornet (till exempel "nordvästra") för rummet (x, y)
// i ett torn med storleken size.
func (d *Data) region(x, y, size int) string {
	part := func(v int) int {
		switch {
		case v < size/3:
			return 0
		case v >= size-size/3:
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
