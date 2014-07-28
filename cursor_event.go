package main

import "github.com/nsf/termbox-go"

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
