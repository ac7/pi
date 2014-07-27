package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nsf/termbox-go"
)

type buffer struct {
	Filename string
	Lines    [][]byte
	Cursor   *cursor
	Topline  int
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
func (b *buffer) Update() {
	b.Cursor.Update()
	for i := b.Topline; i < b.Topline+b.Height(); i++ {
		if i < 0 {
			continue
		} else if i >= len(b.Lines) {
			break
		}

		line := b.Lines[i]
		puts(0, i-b.Topline, fmt.Sprintf(fmt.Sprintf("%%%dd", _LEFT_MARGIN-1), i+1), termbox.ColorCyan, termbox.ColorBlack|termbox.AttrUnderline)
		puts(_LEFT_MARGIN, i-b.Topline, fmt.Sprintf("%s", line), termbox.ColorWhite, termbox.ColorBlack)
	}
	statusLine("In buffer " + b.Filename)
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
	buf.Lines = bytes.Split(data, []byte{'\n'})
	if len(buf.Lines) > 1 {
		buf.Lines = buf.Lines[:len(buf.Lines)-1]
	}
	buf.Cursor = newCursor(buf)
	return buf
}

func newEmptyBuffer() *buffer {
	buf := &buffer{Lines: [][]byte{[]byte{}}}
	buf.Cursor = newCursor(buf)
	return buf
}
