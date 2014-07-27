package main

import (
	"bytes"

	"github.com/nsf/termbox-go"
)

type cursor struct {
	x, y int
	buf  *buffer
}

func (c *cursor) Update() {
	lines := c.buf.Lines()

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
	}

	xPos := c.x
	if xPos > len(lines[c.y]) {
		xPos = len(lines[c.y])
	}

	tabCount := bytes.Count(lines[c.y][:xPos], []byte{'\t'})
	termbox.SetCursor(xPos+_LEFT_MARGIN+tabCount*(_TAB_WIDTH-1), c.y-c.buf.Topline)
}

func (c *cursor) moveWord(forward bool) {
	lines := c.buf.Lines()
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

func (c *cursor) HandleEvent(event termbox.Event) {
	if event.Ch == 0 {
		return
	}
	switch event.Ch {
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
