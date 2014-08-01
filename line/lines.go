package lines

import "github.com/ac7/pi"

func DisplayWidth(line string) int {
	width := len(line)
	for _, c := range line {
		if c == '\t' {
			width += pi.TAB_WIDTH - 1
		}
	}
	return width
}
