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

	err := termbox.Init()
	if err != nil {
		fmt.Println("Unable to initalize termbox:", err)
	}
	defer termbox.Close()
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)

	running := true
	for running {
		if len(buffers) == 0 {
			statusLine("No buffers loaded")
		} else {
			buffers[bufferIndex].draw()
		}
		termbox.Flush()
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			switch event.Key {
			case termbox.KeySpace:
				running = false
			}
		}
	}
}
