package buffer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ac7/pi"
	"github.com/ac7/pi/cursor"
	"github.com/ac7/pi/status"
	"github.com/nsf/termbox-go"
)

func (buf *buffer) loadData(data []byte) {
	// we have to do this complicated split, allocate, and copy because otherwise the
	// slices bleed into each other when you edit
	lines := bytes.Split(data, []byte{'\n'})
	if len(lines) > 1 {
		lines = lines[:len(lines)-1]
		status.Set(fmt.Sprintf(`[%s] %d lines loaded`, buf.filename, len(lines)))
	}
	buf.lines = make([]string, len(lines))
	for i, l := range lines {
		buf.lines[i] = string(l)
	}
}

func NewFromFile(filename string) pi.IBuffer {
	var data []byte
	file, err := os.Open(filename)
	if err != nil {
		status.Set(fmt.Sprintf(`[%s] New file loaded`, filename))
	} else {
		defer file.Close()
		data, err = ioutil.ReadAll(file)
		if err != nil {
			status.Set(fmt.Sprintf(`Unable to read from file "%s"`, filename))
			data = []byte{}
		}
	}
	buf := &buffer{filename: filename}

	buf.loadData(data)
	buf.highlightAll()
	buf.findLongestLine()

	buf.cursor = cursor.New(buf)
	return buf
}

func NewFromStream(r io.Reader) pi.IBuffer {
	buf := &buffer{}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		buf.loadData([]byte("Could not load from stdin: " + err.Error()))
	} else {
		buf.loadData(data)
		status.Set(fmt.Sprintf("Read %d lines from stdin", len(buf.lines)))
	}
	buf.highlightAll()
	buf.findLongestLine()

	buf.cursor = cursor.New(buf)
	return buf
}

func NewEmpty() pi.IBuffer {
	buf := &buffer{
		lines:        []string{""},
		highlighting: [][]termbox.Attribute{{}},
	}
	buf.cursor = cursor.New(buf)
	return buf
}
