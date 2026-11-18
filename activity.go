package main

// Autostopp: tornet stänger av sig självt när ingen längre är där, så att en
// bortglömd server inte håller datorn vaken.
//
// Två regler gäller:
//   - Ingen spelare är ansluten och ingen sida är synlig på tornets dator
//     (introt eller värdsidan): tornet stänger efter autostopp, 10 min.
//   - Ingen har gjort något alls, inte ens med en flik som står öppen men
//     gömd: tornet stänger efter autostopp-max, 12 timmar.
//
// En minut innan får alla som är anslutna en varning.

import "time"

const closeWarning = time.Minute

// seen säger att någon är där på riktigt: en spelare som gör något, eller en
// sida som syns på skärmen. Anropas med låset taget.
func (g *Game) seen() {
	now := time.Now()
	g.lastActive = now
	g.lastReal = now
	g.closeAt = time.Time{}
}

// Touch är seen för den som inte redan har låset.
func (g *Game) Touch() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.seen()
}

// ShouldClose säger om tornet ska stänga nu. empty är tiden utan spelare och
// synliga sidor, abandoned tiden utan någon aktivitet alls (0 = regeln gäller
// inte). En minut innan varnas alla som är anslutna.
func (g *Game) ShouldClose(empty, abandoned time.Duration) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := time.Now()
	if g.humansPresent() > 0 {
		g.lastActive = now
	}
	left := time.Duration(-1)
	consider := func(limit time.Duration, since time.Time) {
		if limit <= 0 {
			return
		}
		if l := limit - now.Sub(since); left < 0 || l < left {
			left = max(0, l)
		}
	}
	consider(empty, g.lastActive)
	consider(abandoned, g.lastReal)
	switch {
	case left < 0:
		return false
	case left == 0:
		return true
	case left <= closeWarning:
		if g.closeAt.IsZero() {
			g.closeAt = now.Add(left)
			g.logAll(g.data.Text("stanger_snart"), "varn")
		}
	default:
		g.closeAt = time.Time{}
	}
	return false
}
