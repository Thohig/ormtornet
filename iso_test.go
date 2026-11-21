//go:build !js

package main

import (
	"strings"
	"testing"
)

// Kartan snett ovanifrån byggs från klient/ till web/iso/iso.js, som checkas
// in. Testet ser till att skriptet finns i programmet och har IsoKarta.
func TestIsoScriptEmbedded(t *testing.T) {
	raw, err := webFS.ReadFile("web/iso/iso.js")
	if err != nil {
		t.Fatalf("web/iso/iso.js saknas; kör npm install och npm run bygg i klient/: %v", err)
	}
	if !strings.Contains(string(raw), "IsoKarta") {
		t.Error("web/iso/iso.js har inte IsoKarta; bygg om med npm run bygg i klient/")
	}
}
