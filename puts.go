package main

import "github.com/nsf/termbox-go"

func puts(x, y int, str string, fg, bg termbox.Attribute) {
	pos := 0
	for _, b := range str {
		if b == '\t' {
			pos += 8
		} else {
			pos++
		}
		termbox.SetCell(x+pos, y, b, fg, bg)
	}
}
