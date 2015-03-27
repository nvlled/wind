package wind

import (
	"github.com/nvlled/wind/size"
)

func computeDimension(layer Layer, canvas Canvas) (int, int) {
	cwidth, cheight := canvas.Dimension()
	width := layer.Width().Value(cwidth)
	height := layer.Height().Value(cheight)
	return width, height
}

func defaultSize(s size.T) size.T {
	if s == nil {
		return size.Free
	}
	return s
}

func mapWidths(frames []Layer) []size.T {
	var sizes []size.T
	for _, f := range frames {
		sizes = append(sizes, f.Width())
	}
	return sizes
}

func mapHeights(frames []Layer) []size.T {
	var sizes []size.T
	for _, f := range frames {
		sizes = append(sizes, f.Height())
	}
	return sizes
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
