package pi

import "github.com/nsf/termbox-go"

func Puts(x, y int, str string, fg, bg termbox.Attribute) {
	for pos, b := range str {
		termbox.SetCell(x+pos, y, b, fg, bg)
	}
}
