package main

import "github.com/nsf/termbox-go"

type cursor struct {
	x, y int
	buf  *buffer
}

func (c *cursor) update() {
	lines := c.buf.lines()

	if c.y < 0 {
		c.y = 0
	} else if c.y >= len(lines) {
		c.y = len(lines) - 1
	}

	if c.y < c.buf.topline {
		c.buf.topline -= 8
	} else if c.y >= c.buf.topline+c.buf.height() {
		c.buf.topline += 8
	}

	if c.x < 0 {
		c.x = 0
	} else if c.x > len(lines[c.y]) {
		c.x = len(lines[c.y])
	}

	termbox.SetCursor(c.x+_LEFT_MARGIN, c.y-c.buf.topline)
}

func (c *cursor) moveWord(forward bool) {
	lines := c.buf.lines()
	line := lines[c.y]

	if !forward {
		c.x--
	} else {
		c.x++
	}

	var b byte
	for b != ' ' {
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
	if forward {
		c.x--
	} else {
		c.x++
	}
}
