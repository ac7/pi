package cursor

import (
	"github.com/ac7/pi/interfaces"
	"github.com/nsf/termbox-go"
)

func (c *cursor) HandleEvent(event termbox.Event) {
	if event.Type == termbox.EventResize {
		c.buf.ForceRedraw()
		return
	} else if event.Type != termbox.EventKey {
		return
	}

	// vi-ish keybindings
	switch c.mode {
	case pi.ModeNormal, pi.ModeSelect:
		if event.Key == termbox.KeyCtrlL {
			c.buf.ForceRedraw()
			return
		}
		switch event.Ch {
		// movement
		case 'j':
			c.y++
		case 'k':
			c.y--
		case 'l':
			c.x++
		case 'h':
			if c.x >= len(c.buf.Line(c.y)) {
				c.x = len(c.buf.Line(c.y))
			}
			c.x--
		case 'e':
			if c.x < 0 {
				c.x = 0
			}
			c.moveWord(true)
		case 'b':
			if c.x >= len(c.buf.Line(c.y)) {
				c.x = len(c.buf.Line(c.y)) - 1
			}
			c.moveWord(false)
		case 'g':
			c.y = 0
			c.buf.CenterOnLine(c.y)
		case 'G':
			c.y = c.buf.Len()
			c.buf.CenterOnLine(c.y)
		case 'z':
			c.buf.CenterOnLine(c.y)

		// save/quit
		case 'w':
			c.buf.Save()
		case 'Z':
			err := c.buf.Save()
			if err == nil {
				c.buf.Close()
			}
		case 'q':
			pi.Quit()

		// make edits
		case 'i':
			c.setMode(pi.ModeEdit)
		case 'a':
			c.setMode(pi.ModeEdit)
			c.x++
		case 'A':
			c.setMode(pi.ModeEdit)
			c.x = len(c.buf.Line(c.y))
		case 'I':
			c.setMode(pi.ModeEdit)
			c.x = 0
		case 'd':
			c.cutBuffer = c.buf.Line(c.y)
			c.buf.DeleteLine(c.y)
		case 'x':
			line := c.buf.Line(c.y)
			if c.x < len(line) {
				c.buf.SetLine(c.y, line[:c.x]+line[c.x+1:])
			}
		case 'p':
			c.y++
			c.buf.InsertLine(c.y)
			c.buf.SetLine(c.y, c.cutBuffer)
		case 'O':
			c.y--
			fallthrough
		case 'o':
			c.y++
			c.buf.InsertLine(c.y)
			c.setMode(pi.ModeEdit)
		}
	case pi.ModeEdit:
		ch := event.Ch
		if event.Key != 0 {
			switch event.Key {
			case termbox.KeySpace:
				ch = ' '
			case termbox.KeyTab:
				ch = '\t'
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				line := c.buf.Line(c.y)
				if c.x > 0 {
					c.x--
					c.buf.SetLine(c.y, line[:c.x]+line[c.x+1:])
				} else if c.y > 0 {
					c.buf.DeleteLine(c.y)
					c.y--
					aboveLine := c.buf.Line(c.y)
					c.buf.SetLine(c.y, c.buf.Line(c.y)+line)
					c.x = len(aboveLine)
				}
				return
			case termbox.KeyCtrlC, termbox.KeyEsc:
				c.setMode(pi.ModeNormal)
				return
			case termbox.KeyEnter:
				line := c.buf.Line(c.y)
				c.buf.InsertLine(c.y)
				c.buf.SetLine(c.y, line[:c.x])
				c.y++
				c.buf.SetLine(c.y, line[c.x:])
				c.x = 0
				return
			default:
				return
			}
		}

		// insert the byte into the middle of the line
		if c.x == len(c.buf.Line(c.y)) {
			// shortcut case: when we're already at the end of the line
			c.buf.SetLine(c.y, c.buf.Line(c.y)+string(ch))
		} else {
			c.buf.SetLine(c.y, c.buf.Line(c.y)[:c.x]+string(ch)+c.buf.Line(c.y)[c.x:])
		}
		c.x++
	}
}
