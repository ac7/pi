package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nsf/termbox-go"
)

type buffer struct {
	filename string
	data     []byte

	cachedLines [][]byte
	cacheValid  bool

	Curs    *cursor
	Topline int
}

func (b *buffer) Lines() [][]byte {
	if !b.cacheValid {
		b.cachedLines = bytes.Split(b.data, []byte{'\n'})
		if len(b.cachedLines) > 1 {
			b.cachedLines = b.cachedLines[:len(b.cachedLines)-1]
		}
		b.cacheValid = true
	}
	return b.cachedLines
}

func (b *buffer) Width() int {
	w, _ := termbox.Size()
	return w - _LEFT_MARGIN
}

func (b *buffer) Height() int {
	_, h := termbox.Size()
	return h - 1 // room for the status bar
}

// TODO: this is very inefficient!
func (b *buffer) Draw() {
	b.Curs.Update()
	lines := b.Lines()
	for i := b.Topline; i < b.Topline+b.Height(); i++ {
		if i < 0 {
			continue
		} else if i >= len(lines) {
			break
		}

		line := lines[i]
		puts(0, i-b.Topline, fmt.Sprintf(fmt.Sprintf("%%%dd", _LEFT_MARGIN-1), i+1), termbox.ColorCyan, termbox.ColorBlack|termbox.AttrUnderline)
		puts(_LEFT_MARGIN, i-b.Topline, fmt.Sprintf("%s", line), termbox.ColorWhite, termbox.ColorBlack)
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
	buf.Curs = &cursor{0, 0, buf}
	return buf
}

func newEmptyBuffer() *buffer {
	buf := &buffer{
		data: []byte{},
	}
	buf.Curs = &cursor{0, 0, buf}
	return buf
}
