package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nsf/termbox-go"
)

const _LEFT_MARGIN = 4

type buffer struct {
	filename string
	data     []byte
	curs     *cursor
}

func (b *buffer) lines() [][]byte {
	lines := bytes.Split(b.data, []byte{'\n'})
	return lines
}

// TODO: this is very inefficient!
func (b *buffer) draw() {
	for i, line := range b.lines() {
		puts(0, i, fmt.Sprintf(fmt.Sprintf("%%%dd", _LEFT_MARGIN-1), i), termbox.ColorCyan, termbox.ColorBlack|termbox.AttrUnderline)
		puts(_LEFT_MARGIN, i, fmt.Sprintf("%s", line), termbox.ColorWhite, termbox.ColorBlack)
	}
	b.curs.update(b)
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
	return &buffer{
		filename,
		data,
		&cursor{0, 0},
	}
}

func newEmptyBuffer() *buffer {
	return &buffer{
		"",
		[]byte{},
		&cursor{0, 0},
	}
}
