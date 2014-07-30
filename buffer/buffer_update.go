package buffer

import (
	"fmt"

	"github.com/ac7/pi"
	"github.com/nsf/termbox-go"
)

func (buf *buffer) Update() {
	if pi.HORIZONTAL_CENTERING {
		buf.leftEdge = buf.Width()/2 - buf.longestLineLen/2 - pi.LEFT_MARGIN
	} else {
		buf.leftEdge = pi.LEFT_MARGIN
	}
	buf.cursor.Update()
}

func (buf *buffer) ForceRedraw() {
	termbox.Clear(termbox.ColorBlack, termbox.ColorWhite)

	for i := buf.topEdge; i < buf.topEdge+buf.Height()-1; i++ {
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
	pi.Puts(buf.leftEdge-pi.LEFT_MARGIN, i-buf.topEdge, fmt.Sprintf(fmt.Sprintf("%%%dd", pi.LEFT_MARGIN-1), i+1), termbox.ColorCyan, termbox.ColorWhite)

	// actual line
	pos := 0
	for x, c := range line {
		termbox.SetCell(pos+buf.leftEdge, i-buf.topEdge, c,
			buf.highlighting[i][x], termbox.ColorWhite)

		if c == '\t' {
			pos += pi.TAB_WIDTH
		} else {
			pos++
		}
	}
}
