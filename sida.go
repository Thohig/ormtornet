//go:build !js

package main

// Tornet som en sida för GitHub Pages (docs/github-pages.md). writeSite
// skriver web/ till en mapp, med det som en vanlig webbserver behöver för
// att visa sidorna utan tornets program: bord.html där sidorna länkar till
// den, och sidornas texter som färdiga skript. Tornet själv, torn.wasm och
// wasm_exec.js, läggs dit av verktyg/bygg-sida.sh.

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func writeSite(dir string, data *Data) error {
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		return err
	}
	if err := os.CopyFS(dir, sub); err != nil {
		return err
	}
	// Spelledarbordet heter spelledare/bord(.html) i tornets program, och
	// sidorna länkar dit med relativa adresser.
	board, err := fs.ReadFile(sub, "bord.html")
	if err != nil {
		return err
	}
	if err := writeFile(filepath.Join(dir, "spelledare", "bord.html"), board); err != nil {
		return err
	}
	// Texterna på svenska. Byter tornet språk hämtar sidorna nya texter
	// från tornet. ORMTORNET_SIDA säger åt web/natet.js att sidorna kommer
	// från Pages och att tornet körs i webbläsaren.
	for page, texts := range map[string]map[string]string{"spelarsidan": data.Spelarsidan, "vardsidan": data.Vardsidan} {
		b, err := json.Marshal(map[string]any{"sprak": data.Lang, "texter": texts})
		if err != nil {
			return err
		}
		js := fmt.Sprintf("window.UI = %s;\nwindow.ORMTORNET_SIDA = true;\n", b)
		if err := writeFile(filepath.Join(dir, "texter", page+".js"), []byte(js)); err != nil {
			return err
		}
	}
	return nil
}

func writeFile(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
