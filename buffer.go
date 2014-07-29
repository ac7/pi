package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nsf/termbox-go"
)

type buffer struct {
	Filename          string
	Lines             []string
	Highlighting      [][]termbox.Attribute
	LongestLineLen    int
	Cursor            *cursor
	Topline           int
	XOffset           int
	Closed            bool
	ChangedSinceWrite bool
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

func (buf *buffer) Save() error {
	file, err := os.Create(buf.Filename)
	if err != nil {
		StatusLine(fmt.Sprintf(`Could not open file "%s" for writing: %s`, buf.Filename, err))
		return err
	}
	defer file.Close()
	for _, l := range buf.Lines {
		file.WriteString(l + "\n")
	}
	StatusLine(fmt.Sprintf(`[%s] %d lines written to disk`, buf.Filename, len(buf.Lines)))
	buf.ChangedSinceWrite = false
	return nil
}

func (buf *buffer) Close() error {
	if buf.ChangedSinceWrite {
		return fmt.Errorf("Unsaved changes")
	}
	buf.Closed = true
	return nil
}

func (buf *buffer) CenterOnCursor() {
	buf.Topline = buf.Cursor.y - buf.Height()/2
}

func newBuffer(filename string) *buffer {
	var data []byte
	file, err := os.Open(filename)
	if err != nil {
		StatusLine(fmt.Sprintf(`[%s] New file loaded`, filename))
	} else {
		defer file.Close()
		data, err = ioutil.ReadAll(file)
		if err != nil {
			StatusLine(fmt.Sprintf(`Unable to read from file "%s"`, filename))
			data = []byte{}
		}
	}
	buf := &buffer{Filename: filename}

	// we have to do this complicated split, allocate, and copy because otherwise the
	// slices bleed into each other when you edit
	lines := bytes.Split(data, []byte{'\n'})
	if len(lines) > 1 {
		lines = lines[:len(lines)-1]
		StatusLine(fmt.Sprintf(`[%s] %d lines loaded`, filename, len(lines)))
	}
	buf.Lines = make([]string, len(lines))
	for i, l := range lines {
		buf.Lines[i] = string(l)
	}

	buf.findLongestLine()
	buf.Cursor = newCursor(buf)
	return buf
}

func newEmptyBuffer() *buffer {
	buf := &buffer{Lines: []string{""}}
	buf.Cursor = newCursor(buf)
	return buf
}
