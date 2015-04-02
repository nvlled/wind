package wind

import (
	"github.com/nvlled/wind/size"
)

// TODO: Rename Vlayer to Ylayer, Hlayer to Xlayer
//		 because it sounds more cool-er
//		 and consistent with Zlayer

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

type listLayer interface {
	Layer
	Elements() []Layer
	AllocSizes(w, h int) ([]int, []int)
	RenderAlloc(canvas Canvas, widths, heights []int)
}

type cacheLayer struct {
	subLayer     listLayer
	width        size.T
	height       size.T
	allocWidths  []int
	allocHeights []int
	sizeCached   bool
	allocCached  bool
}

func (layer *cacheLayer) cacheSize() {
	subLayer := layer.subLayer
	if !layer.sizeCached {
		layer.width = subLayer.Width()
		layer.height = subLayer.Height()
		layer.sizeCached = true
	}
}

func (layer *cacheLayer) cacheAlloc(w, h int) ([]int, []int) {
	if !layer.allocCached {
		subLayer := layer.subLayer
		widths, heights := subLayer.AllocSizes(w, h)
		layer.allocWidths = widths
		layer.allocHeights = heights
		layer.allocCached = true
	}
	return layer.allocWidths, layer.allocHeights
}

func (layer *cacheLayer) Width() size.T {
	layer.cacheSize()
	layer.sizeCached = true
	return layer.width
}

func (layer *cacheLayer) Height() size.T {
	layer.cacheSize()
	return layer.height
}

func (layer *cacheLayer) Render(canvas Canvas) {
	subLayer := layer.subLayer
	w, h := computeDimension(layer, canvas)
	widths, heights := layer.cacheAlloc(w, h)
	subLayer.RenderAlloc(canvas, widths, heights)
}

func renderListLayer(layer listLayer, canvas Canvas) {
	w, h := computeDimension(layer, canvas)
	widths, heights := layer.AllocSizes(w, h)
	layer.RenderAlloc(canvas, widths, heights)
}

type hLayer struct{ elements []Layer }

func (layer hLayer) Elements() []Layer {
	return layer.elements
}

func (layer hLayer) Width() size.T {
	return size.Sum(mapWidths(layer.elements))
}

func (layer hLayer) Height() size.T {
	return size.Max(mapHeights(layer.elements))
}

func (layer hLayer) AllocSizes(w, h int) ([]int, []int) {
	widths := size.AllocFair(w, mapWidths(layer.elements))
	heights := size.AllocMax(h, mapHeights(layer.elements))
	return widths, heights
}

func (layer hLayer) RenderAlloc(canvas Canvas, widths, heights []int) {
	elements := layer.elements
	x, y := 0, 0

	for i, elem := range elements {
		w := widths[i]
		h := heights[i]

		subCanvas := canvas.New(x, y, w, h)
		elem.Render(subCanvas)

		x = x + w
	}
}

func (layer hLayer) Render(canvas Canvas) { renderListLayer(layer, canvas) }

type vLayer struct{ elements []Layer }

func (layer *vLayer) Elements() []Layer {
	return layer.elements
}

func (layer *vLayer) Width() size.T {
	return size.Max(mapWidths(layer.elements))
}

func (layer *vLayer) Height() size.T {
	return size.Sum(mapHeights(layer.elements))
}

func (layer *vLayer) AllocSizes(w, h int) ([]int, []int) {
	widths := size.AllocMax(w, mapWidths(layer.elements))
	heights := size.AllocFair(h, mapHeights(layer.elements))
	return widths, heights
}

func (layer *vLayer) RenderAlloc(canvas Canvas, widths, heights []int) {
	x, y := 0, 0
	for i, elem := range layer.elements {
		w := widths[i]
		h := heights[i]

		subCanvas := canvas.New(x, y, w, h)
		elem.Render(subCanvas)

		y = y + h
	}
}

func (layer *vLayer) Render(canvas Canvas) { renderListLayer(layer, canvas) }

type zLayer struct{ elements []Layer }

func (layer *zLayer) Elements() []Layer {
	return layer.elements
}

func (layer zLayer) Width() size.T {
	return size.Max(mapWidths(layer.elements))
}

func (layer zLayer) Height() size.T {
	return size.Max(mapHeights(layer.elements))
}

func (layer zLayer) AllocSizes(w, h int) ([]int, []int) {
	widths := size.AllocMax(w, mapWidths(layer.elements))
	heights := size.AllocMax(h, mapHeights(layer.elements))
	return widths, heights
}

func (layer zLayer) RenderAlloc(canvas Canvas, widths, heights []int) {
	x, y := 0, 0
	for i, elem := range layer.elements {
		w := widths[i]
		h := heights[i]

		subCanvas := canvas.New(x, y, w, h)
		elem.Render(subCanvas)
	}
}

func (layer *zLayer) Render(canvas Canvas) { renderListLayer(layer, canvas) }

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
