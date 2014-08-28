package prompt

import (
	"fmt"

	"github.com/ac7/pi/interfaces"
	"github.com/nsf/termbox-go"
)

func drawQuery(query string, partialAnswer string) (cursorX, cursorY int) {
	w, h := termbox.Size()
	y := h - 2
	str := query + " " + partialAnswer

	pi.Puts(0, y, fmt.Sprintf(fmt.Sprintf("%%-%ds", w), str), termbox.ColorBlack, termbox.ColorDefault)

	cursorX, cursorY = len(str)+1, y
	return
}

// Ask requests for a string from the user given a certain query. The execution of the
// program is halted while the user is typing. The ok return value signifies that the user
// has confirmed their answer by pressing the enter key.
func Ask(query string) (answer string, ok bool) {
	for {
		drawQuery(query, answer)
		termbox.Flush()
		event := termbox.PollEvent()
		switch event.Key {
		case pi.KillKey:
			pi.Quit()
		case termbox.KeyEsc, termbox.KeyCtrlC:
			return "", false
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			if len(answer) > 0 {
				answer = answer[:len(answer)-1]
			}
		case termbox.KeyEnter:
			return answer, true
		case 0:
			answer += string(event.Ch)
		}
	}
}
