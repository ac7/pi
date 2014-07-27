package main

import "github.com/nsf/termbox-go"

func puts(x, y int, str string, fg, bg termbox.Attribute) {
	for i, b := range str {
		termbox.SetCell(x+i, y, b, fg, bg)
	}
}
