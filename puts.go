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

var _statusLine = ""

func StatusLine(info string) {
	_statusLine = "status: " + info
}

func drawStatusLine() {
	_, h := termbox.Size()
	puts(0, h-1, _statusLine, termbox.ColorBlue, termbox.ColorDefault)
}
