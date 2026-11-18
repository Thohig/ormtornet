// appikon ritar Ormtornets ikon (samma motiv som web/ikon.svg) som PNG-filer
// i alla storlekar en Mac-app behöver, med genomskinliga hörn.
//
//	go run ./verktyg/appikon <mapp.iconset>
//
// Gör sedan om mappen till en .icns med: iconutil -c icns <mapp.iconset>
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

var (
	bg     = color.NRGBA{0x0b, 0x09, 0x06, 0xff}
	sun    = color.NRGBA{0xff, 0x9a, 0x5a, 0xff}
	stone  = color.NRGBA{0xd8, 0xc3, 0x91, 0xff}
	ground = color.NRGBA{0x7d, 0x6d, 0x4c, 0xff}
	star   = color.NRGBA{0xff, 0xe7, 0xb0, 0xff}
	none   = color.NRGBA{}
)

func in(x, y, x0, y0, x1, y1 float64) bool { return x >= x0 && x < x1 && y >= y0 && y < y1 }

func inRounded(x, y, size, r float64) bool {
	if x < 0 || y < 0 || x >= size || y >= size {
		return false
	}
	cx := math.Max(r, math.Min(size-r, x))
	cy := math.Max(r, math.Min(size-r, y))
	return (x-cx)*(x-cx)+(y-cy)*(y-cy) <= r*r
}

func inCircle(x, y, cx, cy, r float64) bool { return (x-cx)*(x-cx)+(y-cy)*(y-cy) <= r*r }

// paint ger färgen i en punkt i ikonens koordinater (0–120), som i SVG-filen.
func paint(x, y float64) color.NRGBA {
	if !inRounded(x, y, 120, 22) {
		return none
	}
	c := bg
	if y <= 86 && inCircle(x, y, 60, 86, 32) {
		c = sun
	}
	if in(x, y, 26, 74, 94, 77) || in(x, y, 26, 64, 94, 66) {
		c = bg
	}
	// Tornet med sina tinnar.
	if in(x, y, 50, 34, 55, 44) || in(x, y, 55, 39, 58, 44) || in(x, y, 58, 34, 62, 44) ||
		in(x, y, 62, 39, 65, 44) || in(x, y, 65, 34, 70, 44) || in(x, y, 52, 43, 68, 86) {
		c = stone
	}
	if in(x, y, 58.5, 52, 61.5, 58) {
		c = bg
	}
	if in(x, y, 18, 86, 102, 88.5) {
		c = ground
	}
	if inCircle(x, y, 30, 30, 1.3) || inCircle(x, y, 88, 26, 1.3) || inCircle(x, y, 96, 44, 1) {
		c = star
	}
	return c
}

// render ritar ikonen i size×size pixlar. Motivet tar upp 80 % av ytan, som
// andra Mac-ikoner, och varje pixel är medelvärdet av 4×4 punkter.
func render(size int) *image.NRGBA {
	const ss = 4
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	scale := 120 / (float64(size) * 0.8)
	off := float64(size) * 0.1
	for py := range size {
		for px := range size {
			var r, g, b, a float64
			for sy := range ss {
				for sx := range ss {
					x := (float64(px) + (float64(sx)+0.5)/ss - off) * scale
					y := (float64(py) + (float64(sy)+0.5)/ss - off) * scale
					c := paint(x, y)
					fa := float64(c.A) / 255
					r += float64(c.R) * fa
					g += float64(c.G) * fa
					b += float64(c.B) * fa
					a += fa
				}
			}
			if a == 0 {
				continue
			}
			img.SetNRGBA(px, py, color.NRGBA{
				uint8(math.Round(r / a)), uint8(math.Round(g / a)), uint8(math.Round(b / a)),
				uint8(math.Round(a / (ss * ss) * 255)),
			})
		}
	}
	return img
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "användning: go run ./verktyg/appikon <mapp.iconset>")
		os.Exit(2)
	}
	dir := os.Args[1]
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, s := range []int{16, 32, 128, 256, 512} {
		for _, k := range []int{1, 2} {
			name := fmt.Sprintf("icon_%dx%d.png", s, s)
			if k == 2 {
				name = fmt.Sprintf("icon_%dx%d@2x.png", s, s)
			}
			f, err := os.Create(filepath.Join(dir, name))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if err := png.Encode(f, render(s*k)); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			f.Close()
		}
	}
}
