//go:build !js

package main

import (
	"regexp"
	"slices"
	"testing"
)

// web/kort.js ritar alla motiv servern känner (cardImages), och inga fler.
func TestCardImagesDrawn(t *testing.T) {
	raw, err := webFS.ReadFile("web/kort.js")
	if err != nil {
		t.Fatal(err)
	}
	var drawn []string
	for _, m := range regexp.MustCompile(`BILDER\.([a-z_]+) =`).FindAllStringSubmatch(string(raw), -1) {
		drawn = append(drawn, m[1])
	}
	for _, b := range cardImages {
		if !slices.Contains(drawn, b) {
			t.Errorf("web/kort.js ritar inte %q", b)
		}
	}
	for _, b := range drawn {
		if !slices.Contains(cardImages, b) {
			t.Errorf("web/kort.js ritar %q, som inte står i cardImages", b)
		}
	}
}

// Rundans kort kommer i ordningen ur data: det viktigaste först.
func TestRoundCardsInOrder(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	g.addCard(p, gameCard{Slag: "dryck", Bild: "fackla", Rubrik: "en dryck"})
	g.addCard(p, gameCard{Slag: "falla", Bild: "pilar", Rubrik: "pilar ur små hål i väggen"})
	g.addCard(p, gameCard{Slag: "okand", Bild: "fackla", Rubrik: "x"})
	cards := g.roundCardsOf(p)
	if len(cards) != 3 || cards[0].Slag != "falla" || cards[1].Slag != "dryck" || cards[2].Slag != "okand" {
		t.Fatalf("ordningen: %+v", cards)
	}
	if cards[0].Rubrik != "Pilar ur små hål i väggen" {
		t.Errorf("rubriken ska börja med stor bokstav: %q", cards[0].Rubrik)
	}
}

// En fälla som slår till ger ett spelkort med den korta beskrivningen som
// rubrik, och en fälla som amuletten tar ett svagt kort. Namnet syns aldrig.
func TestTrapCard(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	putInRoom(g, p, 4, 4)
	tr := g.data.trap("pilar")
	g.trigger(p, tr, 0)
	cards := g.roundCardsOf(p)
	if len(cards) != 1 || cards[0].Slag != "falla" || cards[0].Bild != "pilar" || cards[0].Svag || cards[0].Text == "" {
		t.Fatalf("fällans kort: %+v", cards)
	}
	if cards[0].Rubrik == tr.Namn || cards[0].Rubrik != upperFirst(tr.Kort) {
		t.Errorf("rubriken: %q", cards[0].Rubrik)
	}
	p.roundCards = nil
	p.HP = 9
	g.give(p, "amulett")
	g.trigger(p, tr, 0)
	cards = g.roundCardsOf(p)
	if len(cards) != 1 || !cards[0].Svag {
		t.Fatalf("amuletten: %+v", cards)
	}
}
