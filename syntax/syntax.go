package syntax

import (
	"strings"

	"github.com/ac7/pi"
	"github.com/nsf/termbox-go"
)

var tokens = map[string]termbox.Attribute{
	// keywords
	"break":       termbox.ColorBlue,
	"case":        termbox.ColorBlue,
	"continue":    termbox.ColorBlue,
	"copy":        termbox.ColorBlue,
	"defer":       termbox.ColorBlue,
	"else":        termbox.ColorBlue,
	"fallthrough": termbox.ColorBlue,
	"for":         termbox.ColorBlue,
	"func":        termbox.ColorBlue,
	"if":          termbox.ColorBlue,
	"import":      termbox.ColorBlue,
	"interface":   termbox.ColorBlue,
	"len":         termbox.ColorBlue,
	"package":     termbox.ColorBlue,
	"return":      termbox.ColorBlue,
	"struct":      termbox.ColorBlue,
	"switch":      termbox.ColorBlue,
	"type":        termbox.ColorBlue,

	// anything related to variable allocation
	":=":   termbox.ColorGreen,
	"[]":   termbox.ColorGreen,
	"make": termbox.ColorGreen,
	"new":  termbox.ColorGreen,
	"var":  termbox.ColorGreen,

	// literals
	"false": termbox.ColorCyan,
	"true":  termbox.ColorCyan,
	"0":     termbox.ColorCyan,
	"1":     termbox.ColorCyan,
	"2":     termbox.ColorCyan,
	"3":     termbox.ColorCyan,
	"4":     termbox.ColorCyan,
	"5":     termbox.ColorCyan,
	"6":     termbox.ColorCyan,
	"7":     termbox.ColorCyan,
	"8":     termbox.ColorCyan,
	"9":     termbox.ColorCyan,

	// types
	"bool":   termbox.ColorMagenta,
	"byte":   termbox.ColorMagenta,
	"error":  termbox.ColorMagenta,
	"int":    termbox.ColorMagenta,
	"map":    termbox.ColorMagenta,
	"rune":   termbox.ColorMagenta,
	"string": termbox.ColorMagenta,

	// common variable names
	"err": termbox.ColorRed,
	"ok":  termbox.ColorRed,
}

func Highlighting(line string) []termbox.Attribute {
	attributes := make([]termbox.Attribute, len(line))
	for i := range attributes {
		attributes[i] = termbox.ColorDefault
	}
	if !pi.SYNTAX_HIGHLIGHTING {
		return attributes
	}

	comment := strings.Index(line, "//")
	if comment != -1 {
		for i := comment; i < len(line); i++ {
			attributes[i] = termbox.ColorYellow
		}
		line = line[:comment]
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
