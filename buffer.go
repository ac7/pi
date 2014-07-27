package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/nsf/termbox-go"
)

type buffer struct {
	filename string
	data     []byte
}

// TODO: this is very inefficient!
func (b *buffer) draw() {
	lines := bytes.Split(b.data, []byte{'\n'})
	for i, line := range lines {
		puts(0, i, string(line), termbox.ColorGreen, termbox.ColorDefault)
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
	return &buffer{
		filename,
		data,
	}
}
