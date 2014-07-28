package main

import (
	"bytes"

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
	// vi-ish keybindings
	switch c.mode {
	case _MODE_NORMAL, _MODE_SELECT:
		switch event.Ch {
		// movement
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

		// save/quit
		case 'w':
			c.buf.Save()
		case 'Z':
			err := c.buf.Save()
			if err == nil {
				c.buf.Close()
			}
		case 'Q':
			quit()

		// make edits
		case 'i':
			c.SetMode(_MODE_EDIT)
		case 'a':
			c.SetMode(_MODE_EDIT)
			c.x++
		case 'A':
			c.SetMode(_MODE_EDIT)
			c.x = len(c.buf.Lines[c.y])
		case 'I':
			c.SetMode(_MODE_EDIT)
			c.x = 0
		case 'D':
			c.DeleteLine()
			c.y++
		case 'p':
			c.InsertLine()
			c.buf.Lines[c.y] = make([]byte, len(c.cutBuffer))
			copy(c.buf.Lines[c.y], c.cutBuffer)
		case 'O':
			c.y--
			fallthrough
		case 'o':
			c.InsertLine()
			c.SetMode(_MODE_EDIT)
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
				} else {
					c.DeleteLine()
				}
				return
			case termbox.KeyCtrlC, termbox.KeyEsc:
				c.SetMode(_MODE_NORMAL)
				c.buf.findLongestLine()
				return
			case termbox.KeyEnter:
				c.InsertLine()
				return
			default:
				return
			}
		}

		// insert the byte into the middle of the line
		if c.x == len(c.buf.Lines[c.y]) {
			// shortcut case: when we're already at the end of the line
			c.buf.Lines[c.y] = append(c.buf.Lines[c.y], byte(ch))
		} else {
			c.buf.Lines[c.y] = append(c.buf.Lines[c.y], 0)
			copy(c.buf.Lines[c.y][c.x+1:], c.buf.Lines[c.y][c.x:])
			c.buf.Lines[c.y][c.x] = byte(ch)
		}
		c.x++
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
	c.buf.Lines = append(c.buf.Lines[:c.y], append([][]byte{[]byte{}}, c.buf.Lines[c.y:]...)...)
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
