package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nsf/termbox-go"
)

type buffer struct {
	Filename       string
	Lines          [][]byte
	LongestLineLen int
	Cursor         *cursor
	Topline        int
	XOffset        int
}

func (buf *buffer) Width() int {
	w, _ := termbox.Size()
	return w - _LEFT_MARGIN
}

func (buf *buffer) Height() int {
	_, h := termbox.Size()
	return h - 1 // room for the status bar
}

func (buf *buffer) findLongestLine() {
	buf.LongestLineLen = 0
	for _, l := range buf.Lines {
		if len(l) > buf.LongestLineLen {
			buf.LongestLineLen = len(l)
		}
	}
}

// TODO: this is very inefficient!
func (buf *buffer) Update() {
	buf.XOffset = _LEFT_MARGIN
	if _HORIZONTAL_CENTERING {
		buf.XOffset = buf.Width()/2 - buf.LongestLineLen/2 - _LEFT_MARGIN
	}
	buf.Cursor.Update()
	for i := buf.Topline; i < buf.Topline+buf.Height(); i++ {
		if i < 0 {
			continue
		} else if i >= len(buf.Lines) {
			break
		}

		line := buf.Lines[i]
		puts(buf.XOffset-_LEFT_MARGIN, i-buf.Topline, fmt.Sprintf(fmt.Sprintf("%%%dd", _LEFT_MARGIN-1), i+1), termbox.ColorCyan, termbox.ColorBlack|termbox.AttrUnderline)
		puts(buf.XOffset, i-buf.Topline, fmt.Sprintf("%s", line), termbox.ColorWhite, termbox.ColorBlack)
	}
	statusLine("In buffer " + buf.Filename)
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
	buf := &buffer{Filename: filename}

	// we have to do this complicated split, allocate, and copy because otherwise the
	// slices bleed into each other when you edit
	lines := bytes.Split(data, []byte{'\n'})
	if len(lines) > 1 {
		lines = lines[:len(lines)-1]
	}
	buf.Lines = make([][]byte, len(lines))
	for i, l := range lines {
		buf.Lines[i] = make([]byte, len(l))
		copy(buf.Lines[i], l)
	}

	buf.findLongestLine()
	buf.Cursor = newCursor(buf)
	return buf
}

func newEmptyBuffer() *buffer {
	buf := &buffer{Lines: [][]byte{[]byte{}}}
	buf.Cursor = newCursor(buf)
	return buf
}
