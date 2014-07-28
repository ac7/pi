package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

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

func drawStatusLine(buf *buffer) {
	w, h := termbox.Size()
	var modeString string
	switch buf.Cursor.mode {
	case _MODE_NORMAL:
		modeString = "normal"
	case _MODE_EDIT:
		modeString = "edit"
	}

	puts(0, h-2, fmt.Sprintf(fmt.Sprintf("  %%-%ds%%s", w-14), _statusLine, modeString), termbox.ColorGreen, termbox.ColorBlack)
}
