package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"

	"github.com/ac7/pi"
	"github.com/ac7/pi/buffer"
	"github.com/ac7/pi/prompt"
	"github.com/ac7/pi/status"
)

var running = true
var buffers []pi.IBuffer

func init() {
	pi.Quit = func() {
		for _, buf := range buffers {
			err := buf.Close()
			if err != nil {
				status.Set(fmt.Sprintf("[%s] has unsaved changes", buf.Filename()))
				return
			}
		}
		running = false
	}
}

func main() {
	buffers = []pi.IBuffer{}
	bufferIndex := 0
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			buffers = append(buffers, buffer.NewFromFile(arg))
		}
	} else {
		buffers = append(buffers, buffer.NewFromStream(os.Stdin))
	}

	err := termbox.Init()
	if err != nil {
		fmt.Println("Unable to initalize termbox:", err)
	}
	defer termbox.Close()

	for running {
		if len(buffers) == 0 {
			buffers = []pi.IBuffer{
				buffer.NewEmpty(),
			}
			bufferIndex = 0
		}

		buf := buffers[bufferIndex]
		if buf.Closed() {
			buffers = buffers[:bufferIndex+copy(buffers[bufferIndex:], buffers[bufferIndex+1:])]
			bufferIndex--
			if bufferIndex < 0 {
				bufferIndex = 0
			}
			continue
		}

		buf.Update()
		status.Draw(buf)
		termbox.Flush()

		event := termbox.PollEvent()
		if event.Key == termbox.KeyCtrlQ {
			pi.Quit()
			break
		}
		if buf.Cursor().Mode() == pi.MODE_NORMAL {
			switch event.Ch {
			case '{':
				bufferIndex--
				if bufferIndex < 0 {
					bufferIndex = len(buffers) - 1
				}
				status.Set(fmt.Sprintf(`Switched backward to file [%s]`, buffers[bufferIndex].Filename))
			case '}':
				bufferIndex++
				if bufferIndex >= len(buffers) {
					bufferIndex = 0
				}
				status.Set(fmt.Sprintf(`Switched forward to file [%s]`, buffers[bufferIndex].Filename))
			case ';':
				if filename, ok := prompt.Ask("Filename:"); ok {
					buffers = append(buffers, buffer.NewFromFile(filename))
					bufferIndex = len(buffers) - 1
				}
			default:
				buf.Cursor().HandleEvent(event)
			}
		} else {
			buf.Cursor().HandleEvent(event)
		}
	}
}
