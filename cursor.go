package main

import "github.com/nsf/termbox-go"

type cursor struct {
	x, y int
}

func (c *cursor) update(buf *buffer) {
	lines := buf.lines()

	if c.y < 0 {
		c.y = 0
	} else if c.y >= len(lines) {
		c.y = len(lines) - 1
	}
	if c.x < 0 {
		c.x = 0
	} else if c.x > len(lines[c.y]) {
		c.x = len(lines[c.y])
	}

	termbox.SetCursor(c.x+_LEFT_MARGIN, c.y)
}
