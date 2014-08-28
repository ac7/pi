package buffer

import (
	"fmt"

	"github.com/ac7/pi/interfaces"
	"github.com/nsf/termbox-go"
)

func (buf *buffer) Update() {
	if pi.HorizontalCentering {
		oldLeftEdge := buf.leftEdge
		buf.leftEdge = buf.Width()/2 - buf.longestLineLen/2 - pi.LeftMargin
		if oldLeftEdge != buf.leftEdge {
			buf.ForceRedraw()
		}
	} else {
		buf.leftEdge = pi.LeftMargin
	}
	buf.cursor.Update()
	buf.ForceRedraw()
}

func (buf *buffer) ForceRedraw() {
	termbox.Clear(termbox.ColorBlack, termbox.ColorDefault)

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
	pi.Puts(buf.leftEdge-pi.LeftMargin, i-buf.topEdge, fmt.Sprintf(fmt.Sprintf("%%%dd", pi.LeftMargin-1), i+1), termbox.ColorCyan, termbox.ColorDefault)

	// actual line
	pos := 0
	for x, c := range line {
		termbox.SetCell(pos+buf.leftEdge, i-buf.topEdge, c,
			buf.highlighting[i][x], termbox.ColorDefault)

		if c == '\t' {
			pos += pi.TabWidth
		} else {
			pos++
		}
	}
}
