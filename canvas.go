package wind

import (
	term "github.com/nsf/termbox-go"
)

type rect struct {
	x      int
	y      int
	width  int
	height int
}

func (r rect) Width() int            { return r.width }
func (r rect) Height() int           { return r.height }
func (r rect) Base() (int, int)      { return r.x, r.y }
func (r rect) Dimension() (int, int) { return r.width, r.height }

func (r rect) subRect(x, y, w, h int) rect {
	return rect{
		x:      r.x + x,
		y:      r.y + y,
		width:  clamp(w, 0, r.width),
		height: clamp(h, 0, r.height),
	}
}

func (canvas *StringCanvas) New(x, y, width, height int) Canvas {
	return &StringCanvas{
		buffer: canvas.buffer,
		rect:   canvas.rect.subRect(x, y, width, height),
	}
}

func (canvas *StringCanvas) Draw(x, y int, ch rune, _, _ uint16) {
	if x < canvas.Width() && y < canvas.Height() {
		baseX, baseY := canvas.Base()
		canvas.buffer[baseY+y][baseX+x] = ch
	}
}

func (canvas *StringCanvas) DrawText(x, y int, s string, _, _ uint16) {
	for i, ch := range []rune(s) {
		canvas.Draw(x+i, y, ch, 0, 0)
	}
}

func (canvas *StringCanvas) Clear() {
	for x := 0; x < canvas.Width(); x++ {
		for y := 0; y < canvas.Height(); y++ {
			canvas.Draw(x, y, ' ', 0, 0)
		}
	}
}

func (canvas *StringCanvas) String() string {
	s := ""
	w, h := canvas.Dimension()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			s += string(canvas.buffer[y][x])
		}
		s += "\n"
	}
	return s
}

type TermCanvas struct {
	rect
}

func (canvas *TermCanvas) New(x, y, width, height int) Canvas {
	return &TermCanvas{
		rect: canvas.rect.subRect(x, y, width, height),
	}
}

func (canvas *TermCanvas) Draw(x, y int, ch rune, fg, bg uint16) {
	baseX, baseY := canvas.Base()
	w, h := canvas.Dimension()
	if x >= 0 && x <= w &&
		y >= 0 && y <= h {
		term.SetCell(baseX+x, baseY+y,
			ch, term.Attribute(fg), term.Attribute(bg))
	}
}

func (canvas *TermCanvas) Clear() {
	w, h := canvas.Dimension()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			canvas.Draw(x, y, ' ', 0, 0)
		}
	}
}

func (canvas *TermCanvas) DrawText(x, y int, s string, fg, bg uint16) {
	for i, ch := range []rune(s) {
		canvas.Draw(x+i, y, ch, fg, bg)
	}
}

type FullTermCanvas struct {
	TermCanvas
}

func (canvas *FullTermCanvas) Width() int {
	// embedding doesn't work like inheritance
	width, _ := term.Size()
	canvas.width = width
	return width
}

func (canvas *FullTermCanvas) Height() int {
	_, height := term.Size()
	canvas.height = height
	return height
}

func (canvas *FullTermCanvas) Dimension() (int, int) {
	w, h := term.Size()
	canvas.width = w
	canvas.height = h
	return w, h
}

type ColorCanvas struct {
	fg     uint16
	bg     uint16
	canvas Canvas
}

func (ccanvas *ColorCanvas) New(x, y, width, height int) Canvas {
	return &ColorCanvas{
		fg:     ccanvas.fg,
		bg:     ccanvas.bg,
		canvas: ccanvas.canvas.New(x, y, width, height),
	}
}

func (ccanvas *ColorCanvas) Width() int {
	return ccanvas.canvas.Width()
}

func (ccanvas *ColorCanvas) Height() int {
	return ccanvas.canvas.Height()
}

func (ccanvas *ColorCanvas) Dimension() (int, int) {
	return ccanvas.canvas.Dimension()
}

func (ccanvas *ColorCanvas) Base() (int, int) {
	return ccanvas.canvas.Base()
}

func (ccanvas *ColorCanvas) Clear() {
	for x := 0; x < ccanvas.Width(); x++ {
		for y := 0; y < ccanvas.Height(); y++ {
			ccanvas.Draw(x, y, ' ', ccanvas.fg, ccanvas.bg)
		}
	}
}

func (ccanvas *ColorCanvas) Draw(x, y int, ch rune, fg, bg uint16) {
	if fg == 0 {
		fg = ccanvas.fg
	}
	if bg == 0 {
		bg = ccanvas.bg
	}
	ccanvas.canvas.Draw(x, y, ch, fg, bg)
}

func (ccanvas *ColorCanvas) DrawText(x, y int, s string, fg, bg uint16) {
	if fg == 0 {
		fg = ccanvas.fg
	}
	if bg == 0 {
		bg = ccanvas.bg
	}
	ccanvas.canvas.DrawText(x, y, s, fg, bg)
}
