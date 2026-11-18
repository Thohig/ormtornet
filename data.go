package main

// Läser in rum och texter från data/texter.json. Själva innehållet ligger
// i JSON-filen så att det kan ändras utan att röra koden.

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
)

//go:embed data/texter.json
var embeddedData []byte

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
}

// loadData läser filen på disk om den finns, annars den inbyggda kopian.
func loadData(path string) (*Data, string, error) {
	raw, err := os.ReadFile(path)
	src := path
	if errors.Is(err, os.ErrNotExist) {
		raw, src = embeddedData, "inbyggd kopia"
	} else if err != nil {
		return nil, "", err
	}
	var d Data
	if err := json.Unmarshal(raw, &d); err != nil {
		return nil, "", fmt.Errorf("%s: %w", src, err)
	}
	if len(d.Rumstyper) == 0 {
		return nil, "", fmt.Errorf("%s: inga rumstyper", src)
	}
	if len(d.Horn) != 4 || len(d.Riktningar) != 4 || len(d.Regioner) != 3 {
		return nil, "", fmt.Errorf("%s: horn och riktningar ska ha 4 ord, regioner 3 rader", src)
	}
	return &d, src, nil
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
// Finns ingen text för nivån används närmaste lägre nivå.
func (rt *RoomType) Describe(level int) string {
	for l := level; l >= 1; l-- {
		if list := rt.Beskrivningar[fmt.Sprint(l)]; len(list) > 0 {
			return pick(list)
		}
	}
	return ""
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
