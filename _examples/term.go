package main

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/wind"
)

func createLayer() wind.Layer {
	return wind.Vlayer(
		wind.Border('-', '|', wind.TextLine("This is a text")),
		wind.Line('^'),
		wind.Hlayer(
			wind.Vlayer(
				wind.SetColor(uint16(term.ColorRed), 0, wind.Text("burn this mudder flaundering text")),
				wind.Line('^'),
				wind.Line('*'),
				wind.Line('&'),
				wind.TextLine("you swelling fork"),
				wind.TextLine("keel yourself"),
				wind.SetColor(uint16(term.ColorBlue), 0, wind.Text("cool text")),
			),
			wind.Vlayer(
				wind.Border('.', '.', wind.Size(-1, 5, wind.Text("some text\nwith words\nthat says something"))),
				// another case with a surprising result, constrast (1) with (2)
				// (1)
				wind.Border(
					'+', '+',
					wind.Size(10, 5, wind.NoExpand(wind.Size(5, 3, wind.CharBlock('>')))),
				),
				wind.Line('─'),
				// (2)
				wind.Border(
					'+', '+',
					wind.Size(10, 5, wind.Size(5, 3, wind.CharBlock('>'))),
				),
			),
		),
		wind.Border('x', '+', wind.Border(' ', ' ', wind.SizeW(-1, wind.Text("Saying something pointless to see if something doesn't work\nAlso saying more things to see if doesn't work again\nLastly saying something to just because")))),
	)
}

func main() {
	term.Init()

	canvas := wind.NewTermCanvas()
	layer := createLayer()
	layer.Render(canvas)

	term.Sync()

	term.PollEvent()
	term.Close()
}
