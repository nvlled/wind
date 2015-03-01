package wind

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/wind/size"
)

type Opt struct {
	height size.T
	width  size.T
	align  int
}

type Canvas interface {
	New(baseX, baseY, width, height int) Canvas
	Draw(x, y int, ch rune, fg, bg uint16)

	Width() int
	Height() int
	Dimension() (int, int)
}

type StringCanvas struct {
	buffer [][]rune
	baseX  int
	baseY  int
	width  int
	height int
}

func NewStringCanvas(width, height int) *StringCanvas {
	buffer := make([][]rune, height)
	for i := range buffer {
		buffer[i] = make([]rune, width)
		for j := range buffer[i] {
			buffer[i][j] = ' '
		}
	}
	return &StringCanvas{
		buffer: buffer,
		baseX:  0,
		baseY:  0,
		width:  width,
		height: height,
	}
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
	//if canvas.baseX+x < len(canvas.buffer) && canvas.baseY+y < len(canvas.buffer[0]) {
	if x < canvas.Width() && y < canvas.Height() {
		canvas.buffer[canvas.baseY+y][canvas.baseX+x] = ch
	}
}

func (canvas *StringCanvas) Width() int  { return canvas.width }
func (canvas *StringCanvas) Height() int { return canvas.height }
func (canvas *StringCanvas) Dimension() (int, int) {
	return canvas.width, canvas.height
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
		baseX:  clamp(canvas.baseX+x, 0, width),
		baseY:  clamp(canvas.baseY+y, 0, height),
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

func (canvas *TermCanvas) Width() int  { return canvas.width }
func (canvas *TermCanvas) Height() int { return canvas.height }
func (canvas *TermCanvas) Dimension() (int, int) {
	return canvas.width, canvas.height
}

type Layer interface {
	Width() size.T
	Height() size.T
	Render(canvas Canvas)
	Elements() []Layer
}

type RenderLayer func(canvas Canvas)

func (f RenderLayer) Render(canvas Canvas) { f(canvas) }
func (f RenderLayer) Width() size.T        { return size.Free }
func (f RenderLayer) Height() size.T       { return size.Free }
func (f RenderLayer) Elements() []Layer    { return nil }

func computeDimension(layer Layer, canvas Canvas) (int, int) {
	cwidth, cheight := canvas.Dimension()
	width := layer.Width().Value(cwidth)
	height := layer.Height().Value(cheight)
	return width, height
}

type hLayer struct{ elements []Layer }

func (layer hLayer) Width() size.T {
	return size.Sum(mapWidths(layer.elements))
}

func (layer hLayer) Height() size.T {
	return size.Max(mapHeights(layer.elements))
}

func (layer hLayer) Elements() []Layer {
	return layer.elements
}

func (layer hLayer) Render(canvas Canvas) {
	elements := layer.elements
	width, height := computeDimension(layer, canvas)
	widths := mapWidths(elements)
	heights := mapHeights(elements)
	x, y := 0, 0

	allocWidths := size.AllocFair(width, widths)
	allocHeights := size.AllocMax(height, heights)
	for i, elem := range elements {
		w := allocWidths[i]
		h := allocHeights[i]

		subCanvas := canvas.New(x, y, w, h)
		elem.Render(subCanvas)

		x = x + w
	}
}

type vLayer struct{ elements []Layer }

func (layer *vLayer) Width() size.T {
	return size.Max(mapWidths(layer.elements))
}

func (layer *vLayer) Height() size.T {
	return size.Sum(mapHeights(layer.elements))
}

func (layer *vLayer) Elements() []Layer {
	return layer.elements
}

func (layer *vLayer) Render(canvas Canvas) {
	width, height := computeDimension(layer, canvas)
	widths := mapWidths(layer.elements)
	heights := mapHeights(layer.elements)
	x, y := 0, 0

	allocWidths := size.AllocMax(width, widths)
	allocHeights := size.AllocFair(height, heights)
	for i, elem := range layer.elements {
		w := allocWidths[i]
		h := allocHeights[i]

		subCanvas := canvas.New(x, y, w, h)
		elem.Render(subCanvas)

		y = y + h
	}
}

type zLayer struct{ elements []Layer }

func (layer zLayer) Width() size.T {
	return size.Max(mapWidths(layer.elements))
}

func (layer zLayer) Height() size.T {
	return size.Max(mapHeights(layer.elements))
}

func (layer zLayer) Elements() []Layer {
	return layer.elements
}

func (layer zLayer) Render(canvas Canvas) {
	elements := layer.elements
	width, height := computeDimension(layer, canvas)
	widths := mapWidths(elements)
	heights := mapHeights(elements)
	x, y := 0, 0

	allocWidths := size.AllocMax(width, widths)
	allocHeights := size.AllocMax(height, heights)
	for i, elem := range elements {
		w := allocWidths[i]
		h := allocHeights[i]

		subCanvas := canvas.New(x, y, w, h)
		elem.Render(subCanvas)
	}
}

func Hlayer(elements ...Layer) Layer {
	return &hLayer{elements}
}

func Vlayer(elements ...Layer) Layer {
	return &vLayer{elements}
}

func Zlayer(elements ...Layer) Layer {
	return &zLayer{elements}
}

// meaningful only if subLayer doesn't have Free width or height
type aligner struct {
	layer Layer
	// Has effect only if (sub)layer doesn't have...
	right bool // ..free width
	down  bool // and free height
}

// Needs to have Free width and height
// to have room for aligning
func (aligner *aligner) Width() size.T {
	//return aligner.layer.Width()
	return size.Free
}

func (aligner *aligner) Height() size.T {
	return size.Free
}

func (aligner *aligner) Elements() []Layer {
	return []Layer{aligner.layer}
}

func (aligner *aligner) Render(canvas Canvas) {
	layer := aligner.layer
	if aligner.right {
		x := canvas.Width() - layer.Width().Value(canvas.Width())
		canvas = canvas.New(x, 0, canvas.Width()-x, canvas.Height())
	}
	if aligner.down {
		y := canvas.Height() - layer.Height().Value(canvas.Height())
		canvas = canvas.New(0, y, canvas.Width(), canvas.Height()-y)
	}
	layer.Render(canvas)
}

func AlignRight(layer Layer) Layer {
	return &aligner{layer, true, false}
}

func AlignDown(layer Layer) Layer {
	return &aligner{layer, false, true}
}

func AlignDownRight(layer Layer) Layer {
	return &aligner{layer, true, true}
}

type constrainer struct {
	width  size.T
	height size.T
	layer  Layer
}

func (c *constrainer) Width() size.T {
	if c.width == nil {
		return c.layer.Width()
	}
	return c.width
}

func (c *constrainer) Height() size.T {
	if c.height == nil {
		return c.layer.Height()
	}
	return c.height
}

func (c *constrainer) Elements() []Layer {
	return []Layer{c.layer}
}

func (c *constrainer) Render(canvas Canvas) {
	c.layer.Render(canvas)
}

func Dim(width, height int, layer Layer) Layer {
	w := size.Int(width)
	h := size.Int(height)
	return &constrainer{w, h, layer}
}

func DimW(width int, layer Layer) Layer {
	w := size.Int(width)
	var h size.T = nil
	return &constrainer{w, h, layer}
}

func DimH(height int, layer Layer) Layer {
	var w size.T = nil
	h := size.Int(height)
	return &constrainer{w, h, layer}
}

type Wrapper struct {
	layer Layer
}

func (wrap *Wrapper) Width() size.T {
	return wrap.layer.Width()
}

func (wrap *Wrapper) Height() size.T {
	return wrap.layer.Height()
}

func (wrap *Wrapper) Elements() []Layer {
	return []Layer{wrap.layer}
}

type borderLayer struct {
	Wrapper
	chX rune
	chY rune
}

func (bLayer *borderLayer) Render(canvas Canvas) {
	for x := 0; x < canvas.Width(); x++ {
		canvas.Draw(x, 0, bLayer.chX, 0, 0)
		canvas.Draw(x, canvas.Height()-1, bLayer.chX, 0, 0)
	}
	for y := 0; y < canvas.Height(); y++ {
		canvas.Draw(0, y, bLayer.chY, 0, 0)
		canvas.Draw(canvas.Width()-1, y, bLayer.chY, 0, 0)
	}
	canvas = canvas.New(1, 1, canvas.Width()-2, canvas.Height()-2)
	bLayer.layer.Render(canvas)
}

func BorderLayer(cx, cy rune, layer Layer) Layer {
	return &borderLayer{Wrapper{layer}, cx, cy}
}
