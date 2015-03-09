package wind

import (
	term "github.com/nsf/termbox-go"
)

type StringCanvas struct {
	buffer [][]rune
	baseX  int
	baseY  int
	width  int
	height int
}

func (canvas *StringCanvas) New(x, y, width, height int) Canvas {
	return &StringCanvas{
		buffer: canvas.buffer,
		baseX:  canvas.baseX + x,
		baseY:  canvas.baseY + y,
		width:  width,
		height: height,
	}
}

func (canvas *StringCanvas) Draw(x, y int, ch rune, _, _ uint16) {
	if x < canvas.Width() && y < canvas.Height() {
		canvas.buffer[canvas.baseY+y][canvas.baseX+x] = ch
	}
}

func (canvas *StringCanvas) DrawText(x, y int, s string, _, _ uint16) {
	for i, ch := range []rune(s) {
		canvas.Draw(x+i, y, ch, 0, 0)
	}
}

func (canvas *StringCanvas) Clear() {
	for x := 0; x < canvas.width; x++ {
		for y := 0; y < canvas.height; y++ {
			canvas.Draw(x, y, ' ', 0, 0)
		}
	}
}

func (canvas *StringCanvas) Width() int  { return canvas.width }
func (canvas *StringCanvas) Height() int { return canvas.height }
func (canvas *StringCanvas) Dimension() (int, int) {
	return canvas.width, canvas.height
}

func (canvas *StringCanvas) Base() (int, int) {
	return canvas.baseX, canvas.baseY
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
	baseX  int
	baseY  int
	width  int
	height int
}

func (canvas *TermCanvas) New(x, y, width, height int) Canvas {
	return &TermCanvas{
		baseX:  canvas.baseX + x,
		baseY:  canvas.baseY + y,
		width:  clamp(width, 0, canvas.width),
		height: clamp(height, 0, canvas.height),
	}
}

func (canvas *TermCanvas) Draw(x, y int, ch rune, fg, bg uint16) {
	if x >= 0 && x <= canvas.width &&
		y >= 0 && y <= canvas.height {
		term.SetCell(canvas.baseX+x, canvas.baseY+y,
			ch, term.Attribute(fg), term.Attribute(bg))
	}
}

func (canvas *TermCanvas) Clear() {
	for x := 0; x < canvas.width; x++ {
		for y := 0; y < canvas.height; y++ {
			canvas.Draw(x, y, ' ', 0, 0)
		}
	}
}

func (canvas *TermCanvas) DrawText(x, y int, s string, fg, bg uint16) {
	for i, ch := range []rune(s) {
		canvas.Draw(x+i, y, ch, fg, bg)
	}
}

func (canvas *TermCanvas) Width() int  { return canvas.width }
func (canvas *TermCanvas) Height() int { return canvas.height }
func (canvas *TermCanvas) Dimension() (int, int) {
	return canvas.width, canvas.height
}

func (canvas *TermCanvas) Base() (int, int) {
	return canvas.baseX, canvas.baseY
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
	return ccanvas.Width()
}

func (ccanvas *ColorCanvas) Height() int {
	return ccanvas.Height()
}

func (ccanvas *ColorCanvas) Dimension() (int, int) {
	return ccanvas.Dimension()
}

func (ccanvas *ColorCanvas) Base() (int, int) {
	return ccanvas.Base()
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
