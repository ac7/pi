package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

var running = true

func quit() {
	running = false
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

	buffers[bufferIndex].CenterOnCursor()
	for running {
		buf := buffers[bufferIndex]

		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		buf.Update()
		drawStatusLine()
		termbox.Flush()

		event := termbox.PollEvent()
		buf.Cursor.HandleEvent(event)
	}
}
