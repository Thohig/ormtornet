//go:build !js

package main

// Tornet som en sida för GitHub Pages (docs/github-pages.md). writeSite
// skriver web/ till en mapp, med det som en vanlig webbserver behöver för
// att visa sidorna utan tornets program: bord.html där sidorna länkar till
// den, och sidornas texter som färdiga skript. Tornet själv, torn.wasm och
// wasm_exec.js, läggs dit av verktyg/bygg-sida.sh.
//
// GitHub Pages låter webbläsaren spara filerna en stund. För att en ny
// version ska slå igenom på en gång får sidornas skript, stilar och bilder
// ett versionsnummer i adressen (?v=...), och skripten hämtar det de själva
// laddar med samma nummer (ORMTORNET_V i web/texter.js).

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	return versionPages(dir, siteVersion())
}

// siteVersion är versionsnumret för bygget: början av commiten på GitHub
// (GITHUB_SHA), annars tiden.
func siteVersion() string {
	if sha := os.Getenv("GITHUB_SHA"); len(sha) >= 8 {
		return sha[:8]
	}
	return strconv.FormatInt(time.Now().Unix(), 36)
}

// localRef är en relativ adress till ett skript, en stil eller en bild.
var localRef = regexp.MustCompile(`((?:src|href)=")([^"?#:]+\.(?:js|css|png|svg))"`)

// withVersion lägger till versionsnumret på sidans egna filer.
func withVersion(html, version string) string {
	return localRef.ReplaceAllString(html, `${1}${2}?v=`+version+`"`)
}

// versionPages ger alla sidor i mappen versionsnumret.
func versionPages(dir, version string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(path, []byte(withVersion(string(b), version)), 0o644)
	})
}

func writeFile(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
