package main

import "github.com/nsf/termbox-go"

func puts(x, y int, str string, fg, bg termbox.Attribute) {
	pos := 0
	for _, b := range str {
		termbox.SetCell(x+pos, y, b, fg, bg)
		if b == '\t' {
			pos += _TAB_WIDTH
		} else {
			pos++
		}
	}
}

func statusLine(info string) {
	_, h := termbox.Size()
	puts(0, h-1, ">status: "+info, termbox.ColorBlue, termbox.ColorDefault)
}
