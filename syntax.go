package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

var tokens = []struct {
	literal string
	attr    termbox.Attribute
}{
	{"package", termbox.ColorBlue},
	{"for", termbox.ColorBlue},
	{"len", termbox.ColorBlue},
	{"return", termbox.ColorBlue},
	{"copy", termbox.ColorBlue},
	{"make", termbox.ColorBlue},
	{"new", termbox.ColorBlue},

	{"var", termbox.ColorBlue},
	{":=", termbox.ColorBlue},
	{"func", termbox.ColorBlue},
	{"struct", termbox.ColorBlue},

	{"true", termbox.ColorRed},
	{"false", termbox.ColorRed},

	{"string", termbox.ColorMagenta},
	{"int", termbox.ColorMagenta},
	{"byte", termbox.ColorMagenta},
	{"error", termbox.ColorMagenta},
}

func syntaxHighlight(line string) []termbox.Attribute {
	attr := make([]termbox.Attribute, len(line))
	for i := range attr {
		attr[i] = termbox.ColorDefault
	}

	if !_SYNTAX_HIGHLIGHTING {
		return attr
	}

	cutoff := 0
	for _, token := range tokens {
		index := 0
		for {
			index = strings.Index(line, token.literal)
			if index < 0 {
				break
			}

			cutoff += index
			line = line[index+len(token.literal):]

			for i := range token.literal {
				attr[i+cutoff] = token.attr
			}
			cutoff += len(token.literal)
		}
	}

	return attr
}
