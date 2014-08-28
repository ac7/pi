package lines

import "github.com/ac7/pi/interfaces"

// DisplayWidth returns the width of the line as it would be drawn on the screen.
// Automatically adjust for tab characters.
func DisplayWidth(line string) int {
	width := len(line)
	for _, c := range line {
		if c == '\t' {
			width += pi.TabWidth - 1
		}
	}
	return width
}
