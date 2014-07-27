package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

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
		buf.Update()
		termbox.Flush()

		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			switch event.Ch {
			case 0:
				switch event.Key {
				case termbox.KeySpace:
					running = false
				}
			default:
			}
		}
		buf.Cursor.HandleEvent(event)
	}
}
