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

const KillKey = termbox.KeyCtrlQ

// Represents a buffer in memory editable by a cursor
type IBuffer interface {
	IViewport

	// Each buffer has a cursor
	Cursor() ICursor

	// Get a line at an index. Return an empty string when index is out of range.
	Line(index int) string
	// Set a line at an index. A no-op when index is out of range.
	SetLine(index int, newValue string)

	InsertLine(int)
	DeleteLine(int)

	// Return the number of lines in the buffer
	Len() int
	Filename() string

	Update()

	Save() error
	Close() error
	Closed() bool
}

type IViewport interface {
	// Width does not refer to the size of the buffer, only the size of the buffer
	// on the user's screen.
	Width() int
	// See the description for Width()
	Height() int

	ForceRedraw()
	CenterOnLine(int)

	// The index of the top-most line that's drawn. If the user has scrolled the buffer
	// so that line 32 is shown on the topmost line of their terminal, that's what
	// Topline() would return.
	TopEdge() int
	LeftEdge() int
	SetTopEdge(int)
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

const (
	LEFT_MARGIN          = 4
	TAB_WIDTH            = 8
	CENTER_EVERY_FRAME   = false
	HORIZONTAL_CENTERING = false
	SYNTAX_HIGHLIGHTING  = true
)
