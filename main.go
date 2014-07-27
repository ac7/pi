package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

func statusLine(info string) {
	_, h := termbox.Size()
	puts(0, h-1, ">status: "+info, termbox.ColorBlue, termbox.ColorDefault)
}

func main() {
	buffers := []*buffer{}
	bufferIndex := 0
	for _, arg := range os.Args[1:] {
		buffers = append(buffers, newBuffer(arg))
	}

	if len(buffers) == 0 {
		buffers = []*buffer{
			newEmptyBuffer(),
		}
	}

	err := termbox.Init()
	if err != nil {
		fmt.Println("Unable to initalize termbox:", err)
	}
	defer termbox.Close()

	running := true
	for running {
		buf := buffers[bufferIndex]

		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		buf.draw()
		termbox.Flush()

		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			switch event.Ch {
			case 0:
				switch event.Key {
				case termbox.KeySpace:
					running = false
				}

			case 'j':
				buf.curs.y++
			case 'k':
				buf.curs.y--
			case 'l':
				buf.curs.x++
			case 'h':
				buf.curs.x--
			case 'e':
				buf.curs.moveWord(true)
			case 'b':
				buf.curs.moveWord(false)
			}
		}
	}
}
