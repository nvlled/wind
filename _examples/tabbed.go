package main

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/wind"
)

var tabElem1 = wind.SetColor(uint16(term.ColorRed), 0, wind.CharBlock('1'))
var tabElem2 = wind.SetColor(uint16(term.ColorBlue), 0, wind.CharBlock('2'))
var tabElem3 = wind.SetColor(uint16(term.ColorGreen), 0, wind.CharBlock('3'))

func main() {
	term.Init()

	// TODO: show a tab title that changes with the content

	canvas := wind.NewTermCanvas()
	tab := wind.Tab()
	layer := wind.Vlayer(
		wind.Text("Tabbed"),
		wind.Line('─'),
		tab.SetElements(
			tab.Name("ones", tabElem1),
			tab.Name("twos", tabElem2),
			tabElem3, // allows optional naming
		),
	)

	for {
		layer.Render(canvas)
		term.Flush()
		e := term.PollEvent()
		if e.Key == 0 {
			switch e.Ch {
			case '1':
				tab.ShowName("ones")
			case '2':
				tab.ShowName("twos")
			case '3':
				tab.ShowIndex(2)
			case 'q':
				goto exit
			}
		} else {
			switch e.Key {
			case term.KeyArrowLeft:
				tab.ShowName("blocks")
			case term.KeyArrowRight:
			}
		}
	}

exit:
	term.Close()
}
