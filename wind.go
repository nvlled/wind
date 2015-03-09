package wind

import (
	"github.com/nvlled/wind/size"
)

func (f RenderLayer) Render(canvas Canvas) { f(canvas) }
func (f RenderLayer) Width() size.T        { return size.Free }
func (f RenderLayer) Height() size.T       { return size.Free }

type blank struct{}

func (_ blank) Width() size.T        { return size.Const(0) }
func (_ blank) Height() size.T       { return size.Const(0) }
func (_ blank) Render(canvas Canvas) {}

func wrapNil(layer Layer) Layer {
	if layer == nil {
		return blank{}
	}
	return layer
}

func (fn Defer) Width() size.T        { return wrapNil(fn()).Width() }
func (fn Defer) Height() size.T       { return wrapNil(fn()).Height() }
func (fn Defer) Render(canvas Canvas) { wrapNil(fn()).Render(canvas) }

func computeDimension(layer Layer, canvas Canvas) (int, int) {
	cwidth, cheight := canvas.Dimension()
	width := layer.Width().Value(cwidth)
	height := layer.Height().Value(cheight)
	return width, height
}

type hLayer struct {
	elements []Layer
}

func (layer hLayer) Width() size.T {
	return size.Sum(mapWidths(layer.elements))
}

func (layer hLayer) Height() size.T {
	return size.Max(mapHeights(layer.elements))
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

type aligner struct {
	layer Layer
	right bool
	down  bool
}

// (sub)layer needs to have a size
// smaller than the one allocated to the aligner.
// Otherwise there will be no noticeable effect.
// Returning Free as the size seems to work, but has
// an unexpected result.

// I'm still thinking of what the best approach
// for this, but the alternative options would be:
//   SizeW(10, AlignRight(...))
// or
//   Free(AlignRight())
// return aligner.layer.Width()

func (aligner *aligner) Width() size.T {
	return size.Free
}

func (aligner *aligner) Height() size.T {
	return size.Free
}

func (aligner *aligner) Render(canvas Canvas) {
	x, y := 0, 0
	layer := aligner.layer
	w, h := computeDimension(layer, canvas)

	if aligner.right {
		x = canvas.Width() - w
	}
	if aligner.down {
		y = canvas.Height() - h
	}

	canvas = canvas.New(x, y, w, h)
	layer.Render(canvas)
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

func (c *constrainer) Render(canvas Canvas) {
	c.layer.Render(canvas)
}

type Wrapper struct {
	layer    Layer
	renderer func(canvas Canvas)
}

func (wrap *Wrapper) Width() size.T {
	return wrap.layer.Width()
}

func (wrap *Wrapper) Height() size.T {
	return wrap.layer.Height()
}

func (wrap *Wrapper) Render(canvas Canvas) {
	if wrap.renderer != nil {
		wrap.renderer(canvas)
	} else {
		wrap.layer.Render(canvas)
	}
}

type borderLayer struct {
	layer Layer
	chX   rune
	chY   rune
}

func (bLayer *borderLayer) Width() size.T {
	return bLayer.layer.Width().Add(size.Const(2))
}

func (bLayer *borderLayer) Height() size.T {
	return bLayer.layer.Height().Add(size.Const(2))
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

// ref must not be a subLayer
// or else tortoise all the way down

type syncer struct {
	ref        Layer
	layer      Layer
	syncWidth  bool
	syncHeight bool
}

func (s *syncer) Width() size.T {
	if s.ref != nil && s.syncWidth {
		return s.ref.Width()
	}
	return s.layer.Width()
}

func (s *syncer) Height() size.T {
	if s.ref != nil && s.syncHeight {
		return s.ref.Height()
	}
	return s.layer.Height()
}

func (s *syncer) Render(canvas Canvas) {
	s.layer.Render(canvas)
}
