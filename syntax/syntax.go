package syntax

import (
	"strings"

	"github.com/nsf/termbox-go"

	"github.com/ac7/pi"
)

var tokens = map[string]termbox.Attribute{
	"copy":    termbox.ColorBlue,
	"for":     termbox.ColorBlue,
	"len":     termbox.ColorBlue,
	"package": termbox.ColorBlue,
	"return":  termbox.ColorBlue,

	":=":     termbox.ColorGreen,
	"func":   termbox.ColorGreen,
	"make":   termbox.ColorGreen,
	"new":    termbox.ColorGreen,
	"struct": termbox.ColorGreen,
	"var":    termbox.ColorGreen,

	"false": termbox.ColorRed,
	"true":  termbox.ColorRed,

	"bool":   termbox.ColorMagenta,
	"byte":   termbox.ColorMagenta,
	"error":  termbox.ColorMagenta,
	"int":    termbox.ColorMagenta,
	"string": termbox.ColorMagenta,
}

func Highlighting(line string) []termbox.Attribute {
	attributes := make([]termbox.Attribute, len(line))
	for i := range attributes {
		attributes[i] = termbox.ColorDefault
	}
	if !pi.SYNTAX_HIGHLIGHTING {
		return attributes
	}

	for literal, attribute := range tokens {
		line := line
		cutoff := 0
		for {
			index := strings.Index(line, literal)
			if index < 0 {
				break
			}

			cutoff += index
			line = line[index+len(literal):]
			for i := range literal {
				attributes[i+cutoff] = attribute
			}
			cutoff += len(literal)
		}
	}

	return attributes
}
