package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// TODO: this is very inefficient!
func (buf *buffer) Update() {
	buf.XOffset = _LEFT_MARGIN
	if _HORIZONTAL_CENTERING {
		buf.XOffset = buf.Width()/2 - buf.LongestLineLen/2 - _LEFT_MARGIN
	}
	buf.Cursor.Update()

	for i := buf.Topline; i < buf.Topline+buf.Height()-1; i++ {
		if i < 0 {
			continue
		} else if i >= len(buf.lines) {
			break
		}

		line := buf.lines[i]

		// line number
		puts(buf.XOffset-_LEFT_MARGIN, i-buf.Topline, fmt.Sprintf(fmt.Sprintf("%%%dd", _LEFT_MARGIN-1), i+1), termbox.ColorCyan, termbox.ColorWhite)

		// actual line
		pos := 0
		for x, c := range line {
			termbox.SetCell(pos+buf.XOffset, i-buf.Topline, c,
				buf.highlighting[i][x], termbox.ColorWhite)

			if c == '\t' {
				pos += _TAB_WIDTH
			} else {
				pos++
			}
		}
	}
}
