package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

func (buf *buffer) Update() {
	if _HORIZONTAL_CENTERING {
		buf.XOffset = buf.Width()/2 - buf.LongestLineLen/2 - _LEFT_MARGIN
	} else {
		buf.XOffset = _LEFT_MARGIN
	}
	buf.Cursor.Update()
}

func (buf *buffer) Redraw() {
	termbox.Clear(termbox.ColorBlack, termbox.ColorWhite)

	for i := buf.topline; i < buf.topline+buf.Height()-1; i++ {
		buf.DrawLine(i)
	}
}

func (buf *buffer) DrawLine(i int) {
	if i < 0 {
		return
	} else if i >= len(buf.lines) {
		return
	}

	line := buf.lines[i]

	// line number
	puts(buf.XOffset-_LEFT_MARGIN, i-buf.topline, fmt.Sprintf(fmt.Sprintf("%%%dd", _LEFT_MARGIN-1), i+1), termbox.ColorCyan, termbox.ColorWhite)

	// actual line
	pos := 0
	for x, c := range line {
		termbox.SetCell(pos+buf.XOffset, i-buf.topline, c,
			buf.highlighting[i][x], termbox.ColorWhite)

		if c == '\t' {
			pos += _TAB_WIDTH
		} else {
			pos++
		}
	}
}
