package main

import (
	"bytes"

	"github.com/nsf/termbox-go"
)

type cursor struct {
	x, y int
	mode mode
	buf  *buffer
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

	tabCount := bytes.Count(lines[c.y][:xPos], []byte{'\t'})
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

func (c *cursor) HandleEvent(event termbox.Event) {
	switch c.mode {
	case _MODE_NORMAL, _MODE_SELECT:
		switch event.Ch {
		// vi-ish keybindings
		case 'j':
			c.y++
		case 'k':
			c.y--
		case 'l':
			c.x++
		case 'h':
			if c.x >= len(c.buf.Lines[c.y]) {
				c.x = len(c.buf.Lines[c.y])
			}
			c.x--
		case 'e':
			if c.x < 0 {
				c.x = 0
			}
			c.moveWord(true)
		case 'b':
			if c.x >= len(c.buf.Lines[c.y]) {
				c.x = len(c.buf.Lines[c.y]) - 1
			}
			c.moveWord(false)
		case 'g':
			c.y = 0
		case 'G':
			c.y = len(c.buf.Lines)
		case 'z':
			c.buf.CenterOnCursor()
		case 'Z':
			c.buf.Save()
		case 'Q':
			quit()
		case 'i':
			c.mode = _MODE_EDIT
		case 'a':
			c.mode = _MODE_EDIT
			c.x++
		case 'A':
			c.mode = _MODE_EDIT
			c.x = len(c.buf.Lines[c.y])
		case 'I':
			c.mode = _MODE_EDIT
			c.x = 0
		case 'O':
			c.y--
			fallthrough
		case 'o':
			c.InsertLine()
			c.mode = _MODE_EDIT
		}
	case _MODE_EDIT:
		ch := event.Ch
		if event.Key != 0 {
			switch event.Key {
			case termbox.KeySpace:
				ch = ' '
			case termbox.KeyTab:
				ch = '\t'
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if c.x > 0 {
					c.x--
					copy(c.buf.Lines[c.y][c.x:], c.buf.Lines[c.y][c.x+1:])
					c.buf.Lines[c.y][len(c.buf.Lines[c.y])-1] = 0
					c.buf.Lines[c.y] = c.buf.Lines[c.y][:len(c.buf.Lines[c.y])-1]
				}
				return
			case termbox.KeyCtrlC, termbox.KeyEsc:
				c.mode = _MODE_NORMAL
				c.buf.findLongestLine()
				return
			case termbox.KeyEnter:
				c.InsertLine()
				return
			default:
				return
			}
		}
		if c.x == len(c.buf.Lines[c.y]) {
			c.buf.Lines[c.y] = append(c.buf.Lines[c.y], byte(ch))
			c.x++
		} else {
			c.buf.Lines[c.y] = append(c.buf.Lines[c.y], 0)
			copy(c.buf.Lines[c.y][c.x+1:], c.buf.Lines[c.y][c.x:])
			c.buf.Lines[c.y][c.x] = byte(ch)
			c.x++
		}
	}
}

func (c *cursor) InsertLine() {
	c.y++
	c.buf.Lines = append(c.buf.Lines[:c.y], append([][]byte{[]byte{}}, c.buf.Lines[c.y:]...)...)
	c.x = 0
}

func newCursor(buf *buffer) *cursor {
	return &cursor{
		x:    0,
		y:    0,
		buf:  buf,
		mode: _MODE_NORMAL,
	}
}
