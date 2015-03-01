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
		SizeH(3, text("some text here\nand more text here\nand here")),
		Size(20, 2, spikes),
		SizeH(5, text("you just spent eternities\nworking on this crap\n")),
	)
	layer2 := Border('―', '|', Vlayer(
		layer1,
		SizeH(1, Hlayer(spikes)),
		Zlayer(
			stars,
			SizeH(3, text("wasn't I just going to make text games...\nhow did I come to this")),
			//Size(4, 3, doughs),
			AlignDownRight(Size(10, 3, doughs)),
			AlignDown(SizeH(1, text("oh well... text aligned at the bottom here"))),
			//Zlayer(AlignRight(Size(4, 3, doughs))),
		),
	))
	layer3 := AlignDown(Zlayer(
		Size(22, 22, stars),
		Size(21, 21, doughs),
		Size(20, 20, spikes),
		Size(19, 19, stars),
		Size(18, 18, doughs),
	))
	layer4 := Hlayer(
		Size(30, -1, layer3),
		Size(-1, -1, layer2),
	)
	//println(layer1, layer2, layer3, layer4)
	canvas := NewStringCanvas(100, 25)
	layer4.Render(canvas)
	println(canvas.String())
}

func TestStringCanvas2(t *testing.T) {
	canvas := NewStringCanvas(100, 10)
	layer := Border('+', 'x',
		Zlayer(
			AlignDownRight(Hlayer(
				Size(5, 5, stars),
				Size(8, 6, doughs),
				Size(30, 10, spikes),
			)),
			AlignDown(Size(20, 3, Border('-', '|', text("bordered text")))),
		))
	layer.Render(canvas)
	println(canvas.String())
}

func printInts(xs []int) {
	for _, x := range xs {
		print(x, " ")
	}
	println()
}
