package buffer

import (
	"fmt"
	"os"

	"github.com/ac7/pi/interfaces"
	"github.com/ac7/pi/status"
	"github.com/ac7/pi/syntax"
	"github.com/nsf/termbox-go"
)

type buffer struct {
	filename          string
	lines             []string
	highlighting      [][]termbox.Attribute
	longestLineLen    int
	cursor            pi.ICursor
	topEdge           int
	leftEdge          int
	closed            bool
	changedSinceWrite bool
}

func (buf *buffer) Width() int {
	w, _ := termbox.Size()
	return w - pi.LeftMargin
}

func (buf *buffer) Height() int {
	_, h := termbox.Size()
	return h - 1 // room for the status bar
}

func (buf *buffer) findLongestLine() {
	buf.longestLineLen = 0
	for _, l := range buf.lines {
		if len(l) > buf.longestLineLen {
			buf.longestLineLen = len(l)
		}
	}
}

// getter methods
func (buf *buffer) Len() int         { return len(buf.lines) }
func (buf *buffer) TopEdge() int     { return buf.topEdge }
func (buf *buffer) LeftEdge() int    { return buf.leftEdge }
func (buf *buffer) Filename() string { return buf.filename }
func (buf *buffer) Closed() bool     { return buf.closed }

func (buf *buffer) SetTopEdge(line int) {
	buf.topEdge = line
	buf.ForceRedraw()
}

func (buf *buffer) Save() error {
	file, err := os.Create(buf.filename)
	if err != nil {
		status.Set(fmt.Sprintf(`Could not open file "%s" for writing: %s`, buf.filename, err))
		return err
	}
	defer file.Close()
	for _, l := range buf.lines {
		file.WriteString(l + "\n")
	}
	status.Set(fmt.Sprintf(`[%s] %d lines written to disk`, buf.filename, len(buf.lines)))
	buf.changedSinceWrite = false
	return nil
}

func (buf *buffer) Close() error {
	if buf.changedSinceWrite {
		return fmt.Errorf("Unsaved changes")
	}
	buf.closed = true
	return nil
}

func (buf *buffer) CenterOnLine(line int) {
	buf.SetTopEdge(line - buf.Height()/2)
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
		buf.highlighting[index] = syntax.Highlighting(val)
		buf.findLongestLine()
		buf.changedSinceWrite = true
		buf.DrawLine(index)
	}
}

func (buf *buffer) DeleteLine(index int) {
	if index < 0 || index >= len(buf.lines) || len(buf.lines) <= 1 {
		return
	}
	buf.lines = append(buf.lines[:index], buf.lines[index+1:]...)
	buf.highlighting = append(buf.highlighting[:index], buf.highlighting[index+1:]...)
	buf.ForceRedraw()
	buf.changedSinceWrite = true
}

func (buf *buffer) InsertLine(index int) {
	if index < 0 || index > len(buf.lines) {
		return
	}
	buf.lines = append(buf.lines[:index], append([]string{""}, buf.lines[index:]...)...)
	buf.highlighting = append(buf.highlighting[:index], append([][]termbox.Attribute{{}}, buf.highlighting[index:]...)...)
	buf.ForceRedraw()
	buf.changedSinceWrite = true
}

func (buf *buffer) Cursor() pi.ICursor {
	return buf.cursor
}

func (buf *buffer) highlightAll() {
	buf.highlighting = make([][]termbox.Attribute, len(buf.lines))
	for i, line := range buf.lines {
		buf.highlighting[i] = syntax.Highlighting(line)
	}
}
