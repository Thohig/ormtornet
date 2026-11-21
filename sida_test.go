//go:build !js

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWithVersion(t *testing.T) {
	in := `<link rel="stylesheet" href="../style.css"><script src="iso/iso.js"></script>` +
		`<a href="spelledare.html">x</a><a href="./?spela">y</a><a href="https://a.b/c.js">z</a>`
	want := `<link rel="stylesheet" href="../style.css?v=abc"><script src="iso/iso.js?v=abc"></script>` +
		`<a href="spelledare.html">x</a><a href="./?spela">y</a><a href="https://a.b/c.js">z</a>`
	if got := withVersion(in, "abc"); got != want {
		t.Errorf("withVersion:\n%s\nvill ha\n%s", got, want)
	}
}

func TestWriteSiteVersions(t *testing.T) {
	t.Setenv("GITHUB_SHA", "0123456789abcdef")
	data, _, err := loadDataLang("data", langSv)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := writeSite(dir, data); err != nil {
		t.Fatal(err)
	}
	for _, page := range []string{"index.html", "torn.html", "spelledare.html"} {
		b, err := os.ReadFile(filepath.Join(dir, page))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), `texter.js?v=01234567"`) {
			t.Errorf("%s har inte versionen på texter.js", page)
		}
	}
}
