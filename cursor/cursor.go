package cursor

import (
	"strings"

	"github.com/ac7/pi"
	"github.com/nsf/termbox-go"
)

type cursor struct {
	x, y      int
	cutBuffer string
	mode      pi.CursorMode
	buf       pi.IBuffer
}

func (c *cursor) Mode() pi.CursorMode        { return c.mode }
func (c *cursor) setMode(mode pi.CursorMode) { c.mode = mode }

func (c *cursor) Update() {
	lineCount := c.buf.Len()

	if c.y < 0 {
		c.y = 0
	} else if c.y >= lineCount {
		c.y = lineCount - 1
	}

	if c.y < c.buf.TopEdge() {
		c.buf.SetTopEdge(c.buf.TopEdge() - 8)
	} else if c.y >= c.buf.TopEdge()+c.buf.Height() {
		c.buf.SetTopEdge(c.buf.TopEdge() + 8)
	}

	if c.x < 0 {
		c.x = 0
	} else if c.x > len(c.buf.Line(c.y)) {
		c.x = len(c.buf.Line(c.y))
	}

	if pi.CENTER_EVERY_FRAME {
		c.buf.SetTopEdge(c.y - c.buf.Len()/2)
	}

	tabCount := strings.Count(c.buf.Line(c.y)[:c.x], "\t")
	termbox.SetCursor(c.x+tabCount*(pi.TAB_WIDTH-1)+c.buf.LeftEdge(), c.y-c.buf.TopEdge())
}

func (c *cursor) moveWord(forward bool) {
	line := c.buf.Line(c.y)

	var b byte
	for b != ' ' && b != '.' && b != ')' && b != '(' && b != '\t' {
		if forward {
			c.x++
		} else {
			c.x--
		}
		if c.x < 0 || c.x >= len(line) {
			break
		}
		b = line[c.x]
	}
}

func New(buf pi.IBuffer) pi.ICursor {
	return &cursor{
		x:    0,
		y:    0,
		buf:  buf,
		mode: pi.MODE_NORMAL,
	}
}
