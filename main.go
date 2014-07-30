package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

var running = true

var buffers []*buffer

func quit() {
	for _, buf := range buffers {
		err := buf.Close()
		if err != nil {
			StatusLine(fmt.Sprintf("[%s] has unsaved changes", buf.Filename))
			return
		}
	}
	running = false
}

func main() {
	buffers = []*buffer{}
	bufferIndex := 0
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			buffers = append(buffers, newBufferFromFile(arg))
		}
	} else {
		buffers = append(buffers, newBufferFromStream(os.Stdin))
	}

	err := termbox.Init()
	if err != nil {
		fmt.Println("Unable to initalize termbox:", err)
	}
	defer termbox.Close()

	for running {
		if len(buffers) == 0 {
			buffers = []*buffer{
				newEmptyBuffer(),
			}
			bufferIndex = 0
		}

		buf := buffers[bufferIndex]
		if buf.Closed {
			buffers = buffers[:bufferIndex+copy(buffers[bufferIndex:], buffers[bufferIndex+1:])]
			bufferIndex--
			if bufferIndex < 0 {
				bufferIndex = 0
			}
			continue
		}

		buf.Update()
		drawStatusLine(buf)
		termbox.Flush()

		event := termbox.PollEvent()
		if event.Key == termbox.KeyCtrlQ {
			running = false
			break
		}
		if buf.Cursor.mode == _MODE_NORMAL {
			switch event.Ch {
			case '{':
				bufferIndex--
				if bufferIndex < 0 {
					bufferIndex = len(buffers) - 1
				}
				StatusLine(fmt.Sprintf(`Switched backward to file [%s]`, buffers[bufferIndex].Filename))
			case '}':
				bufferIndex++
				if bufferIndex >= len(buffers) {
					bufferIndex = 0
				}
				StatusLine(fmt.Sprintf(`Switched forward to file [%s]`, buffers[bufferIndex].Filename))
			default:
				buf.Cursor.HandleEvent(event)
			}
		} else {
			buf.Cursor.HandleEvent(event)
		}
	}
}
