//go:build !js

package main

import (
	"regexp"
	"slices"
	"strings"
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

// En varelse som dyker upp ger ett kort med den korta beskrivningen och
// samma text som rumsblocket. En varelse man fäller ger ett överkorsat kort.
func TestMonsterCards(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	putInRoom(g, p, 4, 4)
	typ := g.data.monsterType("krypare")
	m := &Monster{ID: 3, Type: typ.ID, X: 4, Y: 4, HP: 1}
	g.monsters = []*Monster{m}
	g.announceMonster(p)
	cards := g.roundCardsOf(p)
	if len(cards) != 1 || cards[0].Slag != "varelse" || cards[0].Bild != "krypare" || cards[0].Rubrik != upperFirst(typ.Kort) || cards[0].Text != typ.DykerUpp[3%len(typ.DykerUpp)] {
		t.Fatalf("när den dyker upp: %+v", cards)
	}
	p.roundCards = nil
	g.defeat(p, m)
	cards = g.roundCardsOf(p)
	if len(cards) != 1 || !cards[0].Fallen || cards[0].Rubrik != upperFirst(typ.KortBestamd) {
		t.Fatalf("när den faller: %+v", cards)
	}
	// I mörkret syns varelsen inte.
	p.roundCards = nil
	p.Dark = 2
	g.monsters = []*Monster{{ID: 4, Type: typ.ID, X: 4, Y: 4, HP: 1}}
	g.announceMonster(p)
	if cards := g.roundCardsOf(p); len(cards) != 0 {
		t.Fatalf("i mörkret: %+v", cards)
	}
}

// Ett slag mellan vandrare ger båda ett kort med den andras namn, och
// texten säger inte siffrorna; de står på kortets rader.
func TestPvpCards(t *testing.T) {
	g, a, b := startTwo(t)
	pa, pb := a.player, b.player
	putInRoom(g, pa, 4, 4)
	putInRoom(g, pb, 4, 4)
	pa.HP, pb.HP, pb.Treasure = 9, 9, 30
	g.hit(pa, pb, 10, false)
	ca, cb := g.roundCardsOf(pa), g.roundCardsOf(pb)
	if len(ca) != 1 || ca[0].Slag != "vandrare" || ca[0].Rubrik != pb.Name {
		t.Fatalf("den som slog: %+v", ca)
	}
	if len(cb) != 1 || cb[0].Rubrik != pa.Name || strings.Contains(cb[0].Text, "liv") {
		t.Fatalf("den som blev slagen: %+v", cb)
	}
}

// Ett föremål man tar ger ett kort med namnet och vad det gör, och det
// står inte i loggen. Ett förbannat föremål som visar sin baksida ger ett
// förbannat kort. En dryck man dricker ger ett kort utan siffror.
func TestItemCards(t *testing.T) {
	g, a, _ := startTwo(t)
	p := a.player
	r := putInRoom(g, p, 4, 4)
	r.items = []string{"kompass"}
	g.pickUp(p)
	cards := g.roundCardsOf(p)
	typ := g.data.itemType("kompass")
	if len(cards) != 1 || cards[0].Slag != "foremal" || cards[0].Rubrik != typ.Namn || cards[0].Text != typ.Beskrivning {
		t.Fatalf("föremålet: %+v", cards)
	}
	p.roundCards = nil
	g.give(p, "svart_parla")
	g.curse(p, "svart_parla")
	g.curse(p, "svart_parla") // baksidan visas bara en gång
	cards = g.roundCardsOf(p)
	if len(cards) != 1 || cards[0].Slag != "forbannelse" || !cards[0].Forbannad || cards[0].Bild != "svart_parla" {
		t.Fatalf("förbannelsen: %+v", cards)
	}
	p.roundCards = nil
	p.Potions, p.HP = 1, 2
	g.drink(p)
	cards = g.roundCardsOf(p)
	if len(cards) != 1 || cards[0].Slag != "dryck" || strings.ContainsAny(cards[0].Text, "0123456789") {
		t.Fatalf("drycken: %+v", cards)
	}
}
