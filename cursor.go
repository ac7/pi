package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

type cursor struct {
	x, y      int
	mode      mode
	cutBuffer string
	buf       *buffer
}

func (c *cursor) SetMode(mode mode) {
	c.mode = mode
}

func (c *cursor) Update() {
	lines := c.buf.lines

	if c.y < 0 {
		c.y = 0
	} else if c.y >= len(lines) {
		c.y = len(lines) - 1
	}

	if c.y < c.buf.Topline {
		c.buf.Topline -= 8
	} else if c.y >= c.buf.Topline+c.buf.Height() {
		c.buf.Topline += 8
	}

	if c.x < 0 {
		c.x = 0
	} else if c.x > len(lines[c.y]) {
		c.x = len(lines[c.y])
	}

	if _CENTER_EVERY_FRAME {
		c.buf.Topline = c.y - c.buf.Height()/2
	}

	tabCount := strings.Count(lines[c.y][:c.x], "\t")
	termbox.SetCursor(c.x+tabCount*(_TAB_WIDTH-1)+c.buf.XOffset, c.y-c.buf.Topline)
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

func newCursor(buf *buffer) *cursor {
	return &cursor{
		x:    0,
		y:    0,
		buf:  buf,
		mode: _MODE_NORMAL,
	}
}
