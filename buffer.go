package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nsf/termbox-go"
)

type buffer struct {
	filename    string
	unchanged   bool
	data        []byte
	cachedLines [][]byte
	curs        *cursor
	topline     int
}

func (b *buffer) lines() [][]byte {
	if !b.unchanged {
		b.cachedLines = bytes.Split(b.data, []byte{'\n'})
		b.unchanged = true
	}
	return b.cachedLines
}

func (b *buffer) width() int {
	w, _ := termbox.Size()
	return w - _LEFT_MARGIN
}

func (b *buffer) height() int {
	_, h := termbox.Size()
	return h - 1 // room for the status bar
}

// TODO: this is very inefficient!
func (b *buffer) draw() {
	b.curs.update()
	lines := b.lines()
	for i := b.topline; i < b.topline+b.height(); i++ {
		if i < 0 {
			continue
		} else if i >= len(lines) {
			break
		}

		line := lines[i]
		puts(0, i-b.topline, fmt.Sprintf(fmt.Sprintf("%%%dd", _LEFT_MARGIN-1), i+1), termbox.ColorCyan, termbox.ColorBlack|termbox.AttrUnderline)
		puts(_LEFT_MARGIN, i-b.topline, fmt.Sprintf("%s", line), termbox.ColorWhite, termbox.ColorBlack)
	}
	statusLine("In buffer " + b.filename)
}

func newBuffer(filename string) *buffer {
	var data []byte
	file, err := os.Open(filename)
	if err != nil {
		data = []byte("Unable to open file " + filename)
		filename = ""
	} else {
		defer file.Close()
		data, err = ioutil.ReadAll(file)
		if err != nil {
			data = []byte("Unable to read from file " + filename)
			filename = ""
		}
	}
	buf := &buffer{
		filename: filename,
		data:     data,
	}
	buf.curs = &cursor{0, 0, buf}
	return buf
}

func newEmptyBuffer() *buffer {
	buf := &buffer{
		data: []byte{},
	}
	buf.curs = &cursor{0, 0, buf}
	return buf
}
