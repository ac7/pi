package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/nsf/termbox-go"
)

type buffer struct {
	Filename          string
	lines             []string
	highlighting      [][]termbox.Attribute
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
	for _, l := range buf.lines {
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
	for _, l := range buf.lines {
		file.WriteString(l + "\n")
	}
	StatusLine(fmt.Sprintf(`[%s] %d lines written to disk`, buf.Filename, len(buf.lines)))
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

func (buf *buffer) Line(index int) string {
	if index >= 0 && index < len(buf.lines) {
		return buf.lines[index]
	}
	return ""
}

func (buf *buffer) SetLine(index int, val string) {
	if index >= 0 && index < len(buf.lines) {
		buf.lines[index] = val
		buf.highlighting[index] = syntaxHighlight(val)
		buf.ChangedSinceWrite = true
	}
}

func (buf *buffer) DeleteLine(index int) {
	if index < 0 || index >= len(buf.lines) || len(buf.lines) <= 1 {
		return
	}
	buf.lines = append(buf.lines[:index], buf.lines[index+1:]...)
	buf.highlighting = append(buf.highlighting[:index], buf.highlighting[index+1:]...)
	buf.ChangedSinceWrite = true
}

func (buf *buffer) InsertLine(index int) {
	if index < 0 || index > len(buf.lines) {
		return
	}
	buf.lines = append(buf.lines[:index], append([]string{""}, buf.lines[index:]...)...)
	buf.highlighting = append(buf.highlighting[:index], append([][]termbox.Attribute{{}}, buf.highlighting[index:]...)...)
	buf.ChangedSinceWrite = true
}

func (buf *buffer) highlightAll() {
	buf.highlighting = make([][]termbox.Attribute, len(buf.lines))
	for i, line := range buf.lines {
		buf.highlighting[i] = syntaxHighlight(line)
	}
}

func (buf *buffer) loadData(data []byte) {
	// we have to do this complicated split, allocate, and copy because otherwise the
	// slices bleed into each other when you edit
	lines := bytes.Split(data, []byte{'\n'})
	if len(lines) > 1 {
		lines = lines[:len(lines)-1]
		StatusLine(fmt.Sprintf(`[%s] %d lines loaded`, buf.Filename, len(lines)))
	}
	buf.lines = make([]string, len(lines))
	for i, l := range lines {
		buf.lines[i] = string(l)
	}
}

func newBufferFromFile(filename string) *buffer {
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

	buf.loadData(data)
	buf.highlightAll()
	buf.findLongestLine()

	buf.Cursor = newCursor(buf)
	return buf
}

func newBufferFromStream(r io.Reader) *buffer {
	buf := &buffer{}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		buf.loadData([]byte("Could not load from stdin: " + err.Error()))
	} else {
		buf.loadData(data)
		StatusLine(fmt.Sprintf("Read %d lines from stdin", len(buf.lines)))
	}
	buf.highlightAll()
	buf.findLongestLine()

	buf.Cursor = newCursor(buf)
	return buf
}

func newEmptyBuffer() *buffer {
	buf := &buffer{
		lines:        []string{""},
		highlighting: [][]termbox.Attribute{{}},
	}
	buf.Cursor = newCursor(buf)
	return buf
}
