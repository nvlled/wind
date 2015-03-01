package wind

import (
	//"github.com/nvlled/wind/size"
	"testing"
)

//           1         2
//  12345678901234567890
// 1*****^^^^^^^--------
// 2*****^^^^^^^--------
// 3*****^^^^^^^--------
// 4*****---------------
// 5--------------------

func makeRenderLayer(ch rune) Layer {
	return RenderLayer(func(canvas Canvas) {
		w, h := canvas.Dimension()
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				canvas.Draw(x, y, ch, 0, 0)
			}
		}
	})
}

func text(s string) Layer {
	return RenderLayer(func(canvas Canvas) {
		x := 0
		y := 0
		for _, ch := range []rune(s) {
			if ch == '\n' {
				x = 0
				y++
			} else {
				canvas.Draw(x, y, ch, 0, 0)
				x++
			}
		}
	})
}

var blanks = makeRenderLayer('_')
var spikes = makeRenderLayer('^')
var doughs = makeRenderLayer('$')
var stars = makeRenderLayer('*')

//---------------------------------------------
//|a | b  | c |  d |  e | f |                 |
//|  |    |   |    |    |   |                 |
//|  |    |   |    |    |   |                 |
//|-----------+-------------+-----------------|
//|    j              i                       |
//|                                           |
//---------------------------------------------
// a: free
// b: min(4)
// c: max(5)

// d: const(4)
// e: max(5)
// f: free

// j.width(): free + min(4) + max(5) = range(4, 5)
// i.width(): const(4) + max(5) + free =

//           1         2
//  12345678901234567890
//  ****^^^^^___________
//  ****^^^^^___________
//  ****^^^^^___________
//  ____^^^^^___________
//  ____^^^^^___________
//   a   b
// a: free
// b: const(9)

// a + b = min(9)
// compute(8, min(9)) = 8
func TestStringCanvas(t *testing.T) {
	layer1 := Vlayer(
		Dim(-1, 3, text("some text here\nand more text here\nand here")),
		Dim(20, 2, spikes),
		Dim(-1, 5, text("you just spent eternities\nworking on this crap\n")),
	)
	layer2 := BorderLayer('―', '|', Vlayer(
		layer1,
		Dim(-1, 1, Hlayer(spikes)),
		Zlayer(
			Dim(-1, -1, stars),
			Dim(-1, 3, text("wasn't I just going to make text games...\nhow did I come to this")),
			//Dim(4, 3, doughs),
			AlignDownRight(Dim(10, 3, doughs)),
			AlignDown(Dim(-1, 1, text("oh well... text aligned at the bottom here"))),
			//Zlayer(AlignRight(Dim(4, 3, doughs))),
		),
	))
	layer3 := AlignDown(Zlayer(
		Dim(22, 22, stars),
		Dim(21, 21, doughs),
		Dim(20, 20, spikes),
		Dim(19, 19, stars),
		Dim(18, 18, doughs),
	))
	layer4 := Hlayer(
		Dim(30, -1, layer3),
		Dim(-1, -1, layer2),
	)
	//println(layer1, layer2, layer3, layer4)
	canvas := NewStringCanvas(100, 25)
	layer4.Render(canvas)
	println(canvas.String())
}

func TestStringCanvas2(t *testing.T) {
	canvas := NewStringCanvas(100, 5)
	layer := Hlayer(AlignDownRight(Hlayer(
		Dim(5, 5, stars),
		Dim(8, 6, doughs),
		Dim(30, 10, spikes),
	)))
	layer.Render(canvas)
	println(canvas.String())
}

func printInts(xs []int) {
	for _, x := range xs {
		print(x, " ")
	}
	println()
}
