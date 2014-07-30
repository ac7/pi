package pi

import "github.com/nsf/termbox-go"

type CursorMode int

const (
	MODE_NORMAL CursorMode = iota
	MODE_SELECT
	MODE_EDIT
)

// Injected by the main executable
var Quit func()

// Represents a buffer in memory editable by a cursor
type IBuffer interface {
	// Get a line at an index. Return an empty string when index is out of range.
	Line(index int) string
	// Set a line at an index. A no-op when index is out of range.
	SetLine(index int, newValue string)

	InsertLine(int)
	DeleteLine(int)

	// The index of the top-most line that's drawn. If the user has scrolled the buffer
	// so that line 32 is shown on the topmost line of their terminal, that's what
	// Topline() would return.
	TopEdge() int
	LeftEdge() int
	// Return the number of lines in the buffer
	Len() int
	// Width does not refer to the size of the buffer, only the size of the buffer
	// on the user's screen.
	Width() int
	// See the description for Width()
	Height() int
	Filename() string

	SetTopEdge(int)

	Update()
	ForceRedraw()

	CenterOnLine(int)
	Save() error
	Close() error
	Closed() bool

	Cursor() ICursor
}

type ICursor interface {
	// Return the current mode of the cursor
	Mode() CursorMode

	// Update the cursor. This does bounds checking and moves the actual terminal cursor
	// to the correct position so the user can see where they are in the file.
	Update()

	// Handle an input event by the user
	HandleEvent(termbox.Event)
}

func Puts(x, y int, str string, fg, bg termbox.Attribute) {
	for pos, b := range str {
		termbox.SetCell(x+pos, y, b, fg, bg)
	}
}
