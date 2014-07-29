package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

func puts(x, y int, str string, fg, bg termbox.Attribute) {
	for pos, b := range str {
		termbox.SetCell(x+pos, y, b, fg, bg)
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

	puts(2, h-2, fmt.Sprintf(fmt.Sprintf("%%-%ds%%s", w-14), _statusLine, modeString), termbox.ColorRed|termbox.AttrBold, termbox.ColorWhite)
}
