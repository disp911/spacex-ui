package docsite

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"html"
	"image"
	"image/color"
	"image/png"
	"math"
	"strings"
)

// The site's pictures are drawn here from the theme and the project: the
// logo mark, the picture on the front page, the diagram of what the project
// does and the icons of its features. Inline pictures take their colours
// from the style sheet, so they follow the light and dark schemes.

var marks = []string{"hex", "rings", "layers", "monogram", "diamond", "orbit", "chevrons", "quarters"}

// mark draws the logo in a 32×32 box with the given colours.
func mark(kind, letter, a, b string) string {
	var s string
	switch kind {
	case "hex":
		s = `<path fill="` + a + `" d="M16 1.5 28.6 8.75v14.5L16 30.5 3.4 23.25V8.75z"/><path fill="none" stroke="#fff" stroke-width="2.5" d="M16 9.5 21.6 12.75v6.5L16 22.5l-5.6-3.25v-6.5z"/>`
	case "rings":
		s = `<circle cx="12" cy="16" r="9.5" fill="` + a + `"/><circle cx="20" cy="16" r="9.5" fill="` + b + `" fill-opacity=".85"/>`
	case "layers":
		s = `<rect x="3" y="4" width="26" height="6.5" rx="3" fill="` + a + `"/><rect x="7" y="12.75" width="22" height="6.5" rx="3" fill="` + b + `"/><rect x="11" y="21.5" width="18" height="6.5" rx="3" fill="` + a + `" fill-opacity=".6"/>`
	case "monogram":
		s = `<rect x="2" y="2" width="28" height="28" rx="8" fill="` + a + `"/><text x="16" y="22.5" fill="#fff" font-family="system-ui, -apple-system, Segoe UI, sans-serif" font-size="18" font-weight="700" text-anchor="middle">` + html.EscapeString(letter) + `</text>`
	case "diamond":
		s = `<path fill="` + a + `" d="M16 1 31 16 16 31 1 16z"/><path fill="` + b + `" d="M16 9 23 16 16 23 9 16z"/>`
	case "orbit":
		s = `<circle cx="16" cy="16" r="8" fill="` + a + `"/><ellipse cx="16" cy="16" rx="14" ry="6" fill="none" stroke="` + b + `" stroke-width="2.5" transform="rotate(-30 16 16)"/><circle cx="27.5" cy="9.5" r="2.5" fill="` + b + `"/>`
	case "chevrons":
		s = `<path fill="none" stroke="` + a + `" stroke-width="4.5" stroke-linecap="round" stroke-linejoin="round" d="M6 8l8 8-8 8"/><path fill="none" stroke="` + b + `" stroke-width="4.5" stroke-linecap="round" stroke-linejoin="round" d="M16 8l8 8-8 8"/>`
	default: // quarters
		s = `<path fill="` + a + `" d="M2 15A13 13 0 0 1 15 2v13z"/><path fill="` + b + `" d="M17 2a13 13 0 0 1 13 13H17z"/><path fill="` + b + `" d="M2 17h13v13A13 13 0 0 1 2 17z"/><path fill="` + a + `" d="M17 17h13a13 13 0 0 1-13 13z"/>`
	}
	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" aria-hidden="true">` + s + `</svg>`
}

// icons are 24×24 line drawings stroked with the current colour.
var icons = map[string]string{
	"layers":   `<path d="M12 3 21 8l-9 5-9-5z"/><path d="m3 12.5 9 5 9-5"/><path d="m3 17 9 5 9-5"/>`,
	"lock":     `<rect x="5" y="11" width="14" height="10" rx="2"/><path d="M8 11V8a4 4 0 0 1 8 0v3"/>`,
	"cloud":    `<path d="M7 19h10.5a4.5 4.5 0 0 0 .6-8.96A6.5 6.5 0 0 0 5.6 9.4 4.8 4.8 0 0 0 7 19z"/>`,
	"bolt":     `<path d="M13 2 4 14h7l-1 8 9-12h-7z"/>`,
	"clock":    `<circle cx="12" cy="12" r="9"/><path d="M12 7v5l3.5 2"/>`,
	"shield":   `<path d="M12 3 20 6v6c0 4.8-3.4 8-8 9.5C7.4 20 4 16.8 4 12V6z"/><path d="m8.5 12 2.5 2.5 4.5-5"/>`,
	"terminal": `<rect x="3" y="4" width="18" height="16" rx="2"/><path d="m7 9 3 3-3 3"/><path d="M13 15h4"/>`,
	"filter":   `<path d="M3 5h18l-7 8.5V19l-4 2v-7.5z"/>`,
	"database": `<ellipse cx="12" cy="5.5" rx="8" ry="3"/><path d="M4 5.5v13c0 1.7 3.6 3 8 3s8-1.3 8-3v-13"/><path d="M4 12c0 1.7 3.6 3 8 3s8-1.3 8-3"/>`,
	"plug":     `<path d="M9 2v5"/><path d="M15 2v5"/><path d="M6 7h12v4a6 6 0 0 1-12 0z"/><path d="M12 17v5"/>`,
	"refresh":  `<path d="M20 11A8 8 0 0 0 6.3 6.3L4 8.5"/><path d="M4 4v4.5h4.5"/><path d="M4 13a8 8 0 0 0 13.7 4.7l2.3-2.2"/><path d="M20 20v-4.5h-4.5"/>`,
	"code":     `<path d="m8 7-5 5 5 5"/><path d="m16 7 5 5-5 5"/><path d="m13.5 4-3 16"/>`,
	"chart":    `<path d="M3 20h18"/><path d="M6 20v-7"/><path d="M11 20V6"/><path d="M16 20v-10"/><path d="M21 20V3"/>`,
	"globe":    `<circle cx="12" cy="12" r="9"/><path d="M3 12h18"/><path d="M12 3c2.5 2.7 3.8 5.7 3.8 9s-1.3 6.3-3.8 9c-2.5-2.7-3.8-5.7-3.8-9S9.5 5.7 12 3z"/>`,
	"bell":     `<path d="M6 16v-5a6 6 0 0 1 12 0v5l2 2H4z"/><path d="M10 21h4"/>`,
}

func icon(name string) string {
	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">` + icons[name] + `</svg>`
}

var illustrations = []string{"terminal", "flow", "stack", "graph", "terminal"}

// artStyle colours the inline pictures from the style sheet.
const artStyle = `<style>.f-a{fill:var(--accent)}.f-s{fill:var(--accent-soft)}.f-b{fill:var(--bg)}.f-2{fill:var(--bg2)}.f-l{fill:var(--border)}.f-t{fill:var(--text)}.f-m{fill:var(--muted)}.s-a{stroke:var(--accent)}.s-l{stroke:var(--border)}.s-m{stroke:var(--muted)}.f-c{fill:var(--code-bg)}.f-ct{fill:var(--code-text)}.f-cd{fill:var(--code-dim)}</style>`

// illustration is the picture beside or under the front page's headline.
func illustration(r *rnd, kind string, p *project) string {
	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 560 400" role="img" aria-label="` + html.EscapeString(p.Name) + `">` + artStyle)
	switch kind {
	case "terminal":
		terminalArt(&b, r, p)
	case "flow":
		flowArt(&b, p, 560, 152)
	case "stack":
		stackArt(&b, r)
	default:
		graphArt(&b, r)
	}
	b.WriteString(`</svg>`)
	return b.String()
}

func terminalArt(b *strings.Builder, r *rnd, p *project) {
	// Shapes behind the window.
	fmt.Fprintf(b, `<circle class="f-s" cx="%d" cy="%d" r="%d"/>`, r.between(420, 500), r.between(60, 110), r.between(70, 110))
	fmt.Fprintf(b, `<rect class="f-s" x="%d" y="%d" width="140" height="140" rx="28" transform="rotate(%d 90 330)"/>`, r.between(10, 40), r.between(250, 280), r.between(8, 24))
	// The window fits its lines and sits in the middle.
	h := 100 + 30*len(p.kind.demo)
	top := 200 - h/2
	fmt.Fprintf(b, `<rect class="f-c" x="40" y="%d" width="480" height="%d" rx="14"/>`, top, h)
	for i, c := range []string{"#ff5f57", "#febc2e", "#28c840"} {
		fmt.Fprintf(b, `<circle cx="%d" cy="%d" r="6" fill="%s"/>`, 66+i*20, top+24, c)
	}
	y := top + 68
	for _, line := range p.kind.demo {
		line = p.text(r, line)
		cls := "f-cd"
		if strings.HasPrefix(line, "$ ") {
			cls = "f-ct"
		}
		fmt.Fprintf(b, `<text class="%s" x="66" y="%d" font-family="ui-monospace, Menlo, Consolas, monospace" font-size="15">%s</text>`, cls, y, html.EscapeString(clip(line, 50)))
		y += 30
	}
	fmt.Fprintf(b, `<rect class="f-a" x="66" y="%d" width="10" height="18"/>`, y-14)
}

func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

// flowArt draws what the project does as three boxes joined by arrows, the
// boxes' tops at y.
func flowArt(b *strings.Builder, p *project, w, y int) {
	bw, bh := w*29/100, 96
	gap := (w - 3*bw) / 4
	mid := y + bh/2
	fmt.Fprintf(b, `<defs><marker id="ar" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto"><path class="f-a" d="M0 0 10 5 0 10z"/></marker></defs>`)
	labels := p.kind.flow
	for i := range 3 {
		x := gap + i*(bw+gap)
		cls, text := "f-b", "f-t"
		if i == 1 {
			cls, text = "f-a", "f-b"
		}
		fmt.Fprintf(b, `<rect class="%s s-l" x="%d" y="%d" width="%d" height="%d" rx="12" stroke-width="1.5"/>`, cls, x, y, bw, bh)
		lines := wrapLabel(labels[i], 15)
		if i == 1 {
			fmt.Fprintf(b, `<text class="%s" x="%d" y="%d" font-size="17" font-weight="700" text-anchor="middle">%s</text>`, text, x+bw/2, mid-4-9*(len(lines)-1), html.EscapeString(p.Name))
			for n, l := range lines {
				fmt.Fprintf(b, `<text class="%s" x="%d" y="%d" font-size="13" text-anchor="middle" opacity=".85">%s</text>`, text, x+bw/2, mid+18-9*(len(lines)-1)+n*17, html.EscapeString(l))
			}
		} else {
			for n, l := range lines {
				fmt.Fprintf(b, `<text class="%s" x="%d" y="%d" font-size="14" font-weight="600" text-anchor="middle">%s</text>`, text, x+bw/2, mid+5-9*(len(lines)-1)+n*18, html.EscapeString(l))
			}
		}
		if i < 2 {
			fmt.Fprintf(b, `<path class="s-a" d="M%d %dH%d" stroke-width="2" fill="none" marker-end="url(#ar)"/>`, x+bw+6, mid, x+bw+gap-6)
		}
	}
	// Dotted traces above the boxes, as data going by.
	for i := range 3 {
		x := gap + i*(bw+gap) + bw/2
		fmt.Fprintf(b, `<path class="s-l" d="M%d %dV%d" stroke-width="2" stroke-dasharray="2 6" stroke-linecap="round"/>`, x, y-48, y-8)
		fmt.Fprintf(b, `<circle class="f-s s-a" cx="%d" cy="%d" r="7" stroke-width="1.5"/>`, x, y-56)
	}
}

// wrapLabel breaks a label longer than n characters at the space nearest
// its middle.
func wrapLabel(s string, n int) []string {
	if len(s) <= n {
		return []string{s}
	}
	best := -1
	for i, c := range s {
		if c == ' ' && (best < 0 || abs(i-len(s)/2) < abs(best-len(s)/2)) {
			best = i
		}
	}
	if best < 0 {
		return []string{s}
	}
	return []string{s[:best], s[best+1:]}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
func stackArt(b *strings.Builder, r *rnd) {
	n := r.between(3, 4)
	cx, w, h := 280, 220, 110
	top := 200 - (n-1)*34
	b.WriteString(`<g stroke-linejoin="round">`)
	for i := n - 1; i >= 0; i-- {
		y := top + i*68
		face := fmt.Sprintf("M%d %dL%d %dL%d %dL%d %dZ", cx, y-h/2, cx+w, y, cx, y+h/2, cx-w, y)
		side := fmt.Sprintf("M%d %dL%d %dL%d %dL%d %dZ", cx-w, y, cx, y+h/2, cx, y+h/2+18, cx-w, y+18)
		side2 := fmt.Sprintf("M%d %dL%d %dL%d %dL%d %dZ", cx+w, y, cx, y+h/2, cx, y+h/2+18, cx+w, y+18)
		cls := "f-2"
		if i == 0 {
			cls = "f-a"
		}
		// The sides are opaque, so nothing below shows through them.
		fmt.Fprintf(b, `<path class="f-l" d="%s"/><path class="f-l" d="%s"/><path class="f-m" opacity=".3" d="%s"/><path class="%s s-l" stroke-width="1.5" d="%s"/>`, side, side2, side2, cls, face)
		if i == 0 {
			for j := range 3 {
				fmt.Fprintf(b, `<path class="f-b" opacity=".85" d="M%d %dl36 18-36 18-36-18z"/>`, cx-70+j*70, y-18+j*0)
			}
		}
	}
	b.WriteString(`</g>`)
	for range 6 {
		fmt.Fprintf(b, `<circle class="f-a" opacity=".5" cx="%d" cy="%d" r="%d"/>`, r.between(40, 520), r.between(20, 80), r.between(3, 6))
	}
}

func graphArt(b *strings.Builder, r *rnd) {
	type pt struct{ x, y int }
	var pts []pt
	for row := range 4 {
		for col := range 5 {
			if r.chance(0.8) {
				pts = append(pts, pt{70 + col*105 + r.between(-22, 22), 60 + row*92 + r.between(-18, 18)})
			}
		}
	}
	for i, a := range pts {
		for _, c := range pts[i+1:] {
			d := math.Hypot(float64(a.x-c.x), float64(a.y-c.y))
			if d < 150 && r.chance(0.6) {
				cls := "s-l"
				if r.chance(0.25) {
					cls = "s-a"
				}
				fmt.Fprintf(b, `<path class="%s" d="M%d %dL%d %d" stroke-width="2"/>`, cls, a.x, a.y, c.x, c.y)
			}
		}
	}
	for i, a := range pts {
		switch {
		case i%5 == 2:
			fmt.Fprintf(b, `<circle class="f-a" cx="%d" cy="%d" r="16"/><circle class="f-b" cx="%d" cy="%d" r="6"/>`, a.x, a.y, a.x, a.y)
		case i%3 == 0:
			fmt.Fprintf(b, `<rect class="f-b s-a" x="%d" y="%d" width="26" height="26" rx="7" stroke-width="2"/>`, a.x-13, a.y-13)
		default:
			fmt.Fprintf(b, `<circle class="f-s s-a" cx="%d" cy="%d" r="10" stroke-width="2"/>`, a.x, a.y)
		}
	}
}

// diagram is the flow drawing on the introduction page.
func diagram(p *project) string {
	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 190" role="img" aria-label="` + html.EscapeString(strings.Join(p.kind.flow[:], " → ")) + `">` + artStyle)
	flowArt(&b, p, 640, 74)
	b.WriteString(`</svg>`)
	return b.String()
}

// favicon is an icon file with one 32×32 picture: the accent colour in a
// rounded square with a ring, the closest the logo marks come to at that
// size.
func favicon(hue, sat int) []byte {
	const n = 32
	img := image.NewNRGBA(image.Rect(0, 0, n, n))
	accent := hslColor(float64(hue), float64(sat)/100, .44)
	for y := range n {
		for x := range n {
			px, py := float64(x)+.5, float64(y)+.5
			// Distance to the rounded square, then to the ring.
			qx, qy := math.Abs(px-16)-15+7, math.Abs(py-16)-15+7
			box := math.Hypot(math.Max(qx, 0), math.Max(qy, 0)) + math.Min(math.Max(qx, qy), 0) - 7
			ring := math.Abs(math.Hypot(px-16, py-16)-7.5) - 2.5
			a := clamp01(.5 - box)
			w := clamp01(.5 - ring)
			c := blend(accent, color.NRGBA{255, 255, 255, 255}, w)
			c.A = uint8(a * 255)
			img.SetNRGBA(x, y, c)
		}
	}
	var pic bytes.Buffer
	png.Encode(&pic, img)
	var out bytes.Buffer
	binary.Write(&out, binary.LittleEndian, [3]uint16{0, 1, 1})
	binary.Write(&out, binary.LittleEndian, struct {
		W, H, Colors, Reserved uint8
		Planes, Bits           uint16
		Size, Offset           uint32
	}{n, n, 0, 0, 1, 32, uint32(pic.Len()), 22})
	out.Write(pic.Bytes())
	return out.Bytes()
}

func clamp01(v float64) float64 { return math.Max(0, math.Min(1, v)) }

func blend(a, b color.NRGBA, t float64) color.NRGBA {
	mix := func(x, y uint8) uint8 { return uint8(float64(x)*(1-t) + float64(y)*t + .5) }
	return color.NRGBA{mix(a.R, b.R), mix(a.G, b.G), mix(a.B, b.B), 255}
}

func hslColor(h, s, l float64) color.NRGBA {
	c := (1 - math.Abs(2*l-1)) * s
	hp := math.Mod(h, 360) / 60
	x := c * (1 - math.Abs(math.Mod(hp, 2)-1))
	var r, g, b float64
	switch int(hp) {
	case 0:
		r, g = c, x
	case 1:
		r, g = x, c
	case 2:
		g, b = c, x
	case 3:
		g, b = x, c
	case 4:
		r, b = x, c
	default:
		r, b = c, x
	}
	m := l - c/2
	to := func(v float64) uint8 { return uint8(math.Round((v + m) * 255)) }
	return color.NRGBA{to(r), to(g), to(b), 255}
}
