package main

import (
	"strings"
	"testing"
)

// Ett varsel bestämmer vad som väntar, och det blir så när någon går in.
func TestOmenComesTrue(t *testing.T) {
	for _, c := range []struct {
		monster, trap float64
		kind          string
	}{{0, 1, "falla"}, {1, 0, "varelse"}, {0, 0, "tomt"}} {
		g, a, _ := startTwo(t)
		g.data.Monster.Chans, g.data.Regler.Handelsechans = c.monster, c.trap
		g.level = 2
		logs(a)
		g.dealOmens()
		if len(g.omens) != 2 {
			t.Fatalf("%s: %d varsel, ville ha ett var", c.kind, len(g.omens))
		}
		o := g.omens[0]
		r := &g.rooms[o.Y][o.X]
		if !strings.HasPrefix(o.What, c.kind) || r.Revealed {
			t.Fatalf("%s: varslet %q, rummet byggt %v", c.kind, o.What, r.Revealed)
		}
		if _, fixed := g.fixedDoors(o.X, o.Y); fixed {
			t.Fatal("varsel om en port eller kammaren")
		}
		if len(logs(a)) == 0 {
			t.Fatal("Alva fick inget varsel")
		}
		// Chanserna ändras efteråt, men rummet blir som varslet sa.
		g.data.Monster.Chans, g.data.Regler.Handelsechans = 0, 0
		g.build(o.X, o.Y, true)
		_, id, _ := strings.Cut(o.What, ":")
		switch c.kind {
		case "falla":
			if r.Event != id {
				t.Fatalf("fällan %q, varslet sa %q", r.Event, id)
			}
		case "varelse":
			if m := g.monsterAt(o.X, o.Y); m == nil || m.Type != id {
				t.Fatalf("ingen %s i rummet", id)
			}
		case "tomt":
			if r.Event != "" || g.monsterAt(o.X, o.Y) != nil {
				t.Fatal("något i ett rum som skulle vara tomt")
			}
		}
	}
}
