package main

import (
	"bytes"

	"github.com/nsf/termbox-go"
)

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
	}

	xPos := c.x
	if xPos > len(lines[c.y]) {
		xPos = len(lines[c.y])
	}

	tabCount := bytes.Count(lines[c.y][:xPos], []byte{'\t'})
	termbox.SetCursor(xPos+_LEFT_MARGIN+tabCount*(_TAB_WIDTH-1), c.y-c.buf.topline)
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

func (c *cursor) handleKey(key rune) {
	switch key {
	case 'j':
		c.y++
	case 'k':
		c.y--
	case 'l':
		c.x++
	case 'h':
		c.x--
	case 'e':
		c.moveWord(true)
	case 'b':
		c.moveWord(false)
	}
}
