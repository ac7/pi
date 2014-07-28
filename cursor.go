package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

type cursor struct {
	x, y      int
	mode      mode
	cutBuffer []byte
	buf       *buffer
}

func (c *cursor) SetMode(mode mode) {
	c.mode = mode
	c.buf.ChangedSinceWrite = true
}

func (c *cursor) Update() {
	lines := c.buf.Lines

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

	if _CENTER_EVERY_FRAME {
		c.buf.Topline = c.y - c.buf.Height()/2
	}

	xPos := c.x
	if xPos > len(lines[c.y]) {
		xPos = len(lines[c.y])
	}

	tabCount := strings.Count(lines[c.y][:xPos], "\t")
	termbox.SetCursor(xPos+tabCount*(_TAB_WIDTH-1)+c.buf.XOffset, c.y-c.buf.Topline)
}

func (c *cursor) moveWord(forward bool) {
	line := c.buf.Lines[c.y]

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

func (c *cursor) DeleteLine() {
	if len(c.buf.Lines) < 2 {
		return
	}
	c.cutBuffer = make([]byte, len(c.buf.Lines[c.y]))
	copy(c.cutBuffer, c.buf.Lines[c.y])
	c.buf.Lines = append(c.buf.Lines[:c.y], c.buf.Lines[c.y+1:]...)
	c.y--
	c.buf.ChangedSinceWrite = true
}

func (c *cursor) InsertLine() {
	c.y++
	c.buf.Lines = append(c.buf.Lines[:c.y], append([]string{""}, c.buf.Lines[c.y:]...)...)
	c.x = 0
	c.buf.ChangedSinceWrite = true
}

func newCursor(buf *buffer) *cursor {
	return &cursor{
		x:    0,
		y:    0,
		buf:  buf,
		mode: _MODE_NORMAL,
	}
}
