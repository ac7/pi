package syntax

import (
	"reflect"
	"testing"

	"github.com/nsf/termbox-go"
)

func TestHighlighting(t *testing.T) {
	input := "var x bool = true"
	expected := []termbox.Attribute{
		// var
		termbox.ColorGreen,
		termbox.ColorGreen,
		termbox.ColorGreen,
		termbox.ColorDefault,

		// x
		termbox.ColorDefault,
		termbox.ColorDefault,

		// bool
		termbox.ColorMagenta,
		termbox.ColorMagenta,
		termbox.ColorMagenta,
		termbox.ColorMagenta,
		termbox.ColorDefault,

		// =
		termbox.ColorDefault,
		termbox.ColorDefault,

		// true
		termbox.ColorRed,
		termbox.ColorRed,
		termbox.ColorRed,
		termbox.ColorRed,
	}

	attributes := Highlighting(input)
	if !reflect.DeepEqual(expected, attributes) {
		t.Errorf("Unexpected syntax highlighting\nRecieved: %v\nExpected: %v",
			attributes, expected)
	}
}
