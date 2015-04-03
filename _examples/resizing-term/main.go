package main

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/wind"
)

func createLayer() wind.Layer {
	char := wind.CharBlock
	border := func(layer wind.Layer) wind.Layer {
		return wind.Border('─', '│', layer)
	}
	aaa := border(char('a'))
	bbb := border(char('v'))
	ccc := border(char('c'))
	xxx := border(char('x'))
	yyy := border(char('y'))
	zzz := border(char('z'))
	return wind.Hlayer(
		wind.Vlayer(
			wind.Zlayer(
				aaa,
				wind.SetColor(0, uint16(term.ColorYellow),
					wind.Text("Press 'q' to quit")),
			),
			wind.Vlayer(
				wind.Zlayer(
					bbb,
					// TODO: the text doesn't resize quite right
					wind.SetColor(0, uint16(term.ColorCyan),
						wind.Text("Try resizing the terminal window\nand see if it does work")),
				),
			),
			ccc,
		),
		wind.SetColor(0, uint16(term.ColorBlue), wind.SizeW(20, xxx)),
		wind.SetColor(0, uint16(term.ColorRed), yyy),
		wind.SetColor(0, uint16(term.ColorGreen), zzz),
	)
}

func main() {
	term.Init()

	canvas := wind.NewTermCanvas()
	layer := createLayer()

	for {
		layer.Render(canvas)
		term.Flush()

		e := term.PollEvent()
		if e.Ch == 'q' {
			break
		} else if e.Type == term.EventResize {
			term.Clear(0, 0)
			wind.ClearCache(layer)
		}
	}

	term.Close()
}
