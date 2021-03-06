package wind

import (
	"github.com/nvlled/wind/size"
	"strings"
)

type Canvas interface {
	New(baseX, baseY, width, height int) Canvas
	Draw(x, y int, ch rune, fg, bg uint16)
	DrawText(x, y int, s string, fg, bg uint16)
	Clear()

	Width() int
	Height() int
	Dimension() (int, int)
	Base() (int, int)
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
		rect: rect{
			x:      0,
			y:      0,
			width:  width,
			height: height,
		},
		buffer: buffer,
	}
}

// Invoke termbox.Init() before creating TermCanvas
func NewTermCanvas() Canvas {
	canvas := &FullTermCanvas{}
	canvas.Dimension()
	return canvas
}

func ChangeDefaultColor(fg, bg uint16, canvas Canvas) Canvas {
	return &ColorCanvas{
		fg:     fg,
		bg:     bg,
		canvas: canvas,
	}
}

type Layer interface {
	Width() size.T
	Height() size.T
	Render(canvas Canvas)
}

type TabLayer interface {
	Layer
	Name(name string, elem Layer) Layer
	SetElements(elements ...Layer) TabLayer
	ShowName(name string) TabLayer
	ShowIndex(index int) TabLayer
	Hide() TabLayer
}

type SizedLayer interface {
	Layer
	SetSize(w, h int) SizedLayer

	AdaptWidth() SizedLayer
	AdaptHeight() SizedLayer

	FreeWidth() SizedLayer
	FreeHeight() SizedLayer

	InheritHeight() SizedLayer
	InheritWidth() SizedLayer
}

type RenderLayer func(canvas Canvas)

type Defer func() Layer

func Text(s string) Layer {
	lines := strings.Split(s, "\n")
	h := len(lines)
	w := 0
	for _, line := range lines {
		length := len(line)
		if length > w {
			w = length
		}
	}
	return Size(w, h, RenderLayer(func(canvas Canvas) {
		for y, line := range lines {
			for x, ch := range []rune(line) {
				canvas.Draw(x, y, ch, 0, 0)
			}
		}
	}))
}

func CharBlock(ch rune) Layer {
	return RenderLayer(func(canvas Canvas) {
		w, h := canvas.Dimension()
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				canvas.Draw(x, y, ch, 0, 0)
			}
		}
	})
}

func TextLine(s string) Layer {
	return SizeH(1, RenderLayer(func(canvas Canvas) {
		x := 0
		for _, ch := range []rune(s) {
			if ch == '\n' {
				ch = '↵'
			}
			canvas.Draw(x, 0, ch, 0, 0)
			x++
		}
	}))
}

func Cache(element Layer) Layer {
	if subLayer, ok := element.(listLayer); ok {
		return &cacheLayer{
			subLayer:     subLayer,
			width:        size.Free,
			height:       size.Free,
			allocWidths:  nil,
			allocHeights: nil,
		}
	}
	return element
}

func Hlayer(elements ...Layer) Layer {
	return Cache(&hLayer{elements})
}

func Vlayer(elements ...Layer) Layer {
	return Cache(&vLayer{elements})
}

func Zlayer(elements ...Layer) Layer {
	return Cache(&zLayer{elements})
}

func NoExpand(layer Layer) Layer {
	return &aligner{layer, false, false}
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

// TODO: rename to FreeSize
func Free(layer Layer) Layer {
	return &constrainer{size.Free, size.Free, layer}
}

func Size(width, height int, layer Layer) SizedLayer {
	w := size.Int(width)
	h := size.Int(height)
	return &constrainer{w, h, layer}
}

func SizeW(width int, layer Layer) SizedLayer {
	w := size.Int(width)
	var h size.T = nil
	return &constrainer{w, h, layer}
}

func SizeH(height int, layer Layer) SizedLayer {
	var w size.T = nil
	h := size.Int(height)
	return &constrainer{w, h, layer}
}

func Border(cx, cy rune, layer Layer) Layer {
	return &borderLayer{layer, cx, cy}
}

func LineH(ch rune) Layer {
	return SizeH(1, CharBlock(ch)).AdaptWidth()
}

func LineV(ch rune) Layer {
	return SizeW(1, CharBlock(ch)).AdaptHeight()
}

func TapRender(layer Layer, render func(layer Layer, canvas Canvas)) Layer {
	return &Wrapper{
		layer: layer,
		renderer: func(canvas Canvas) {
			render(layer, canvas)
		},
	}
}

func SetColor(fg, bg uint16, layer Layer) Layer {
	return TapRender(layer, func(layer Layer, canvas Canvas) {
		canvas = ChangeDefaultColor(fg, bg, canvas)
		layer.Render(canvas)
	})
}

// ref must not be a subLayer
// or else tortoise all the way down
func SyncSize(ref Layer, layer Layer) Layer {
	return &syncer{ref, layer, true, true}
}

func SyncSizeW(ref Layer, layer Layer) Layer {
	return &syncer{ref, layer, true, false}
}

func SyncSizeH(ref Layer, layer Layer) Layer {
	return &syncer{ref, layer, false, true}
}

// nil comparison may fail,
// deferred function must make sure
// that nil is returned.
func Either(left Defer, right Layer) Layer {
	return RenderLayer(func(canvas Canvas) {
		layer := left()
		if layer == nil {
			layer = right
		}
		layer.Render(canvas)
	})
}

func Tab() TabLayer {
	return &tabLayer{
		namedElements: make(map[string]Layer),
	}
}

func ClearCache(layer Layer) {
	if cache, ok := layer.(*cacheLayer); ok {
		cache.clear()
	}
}

func PreRender(layer Layer, w, h int) {
	layer.Render(&nilCanvas{
		rect: rect{width: w, height: h},
	})
}
