package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		fmt.Println("Unable to initalize termbox:", err)
	}
	defer termbox.Close()

	termbox.Clear(termbox.ColorWhite, termbox.ColorWhite)
	termbox.HideCursor()
	termbox.SetCell(16, 16, '@', termbox.ColorRed, termbox.ColorDefault)
	termbox.Flush()

	running := true

	for running {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			switch event.Key {
			case termbox.KeySpace:
				running = false
			}
		}
	}
}
