package status

import (
	"fmt"

	"github.com/ac7/pi/interfaces"
	"github.com/nsf/termbox-go"
)

var statusLine = ""

func Set(info string) {
	statusLine = "status: " + info
}

func Draw(buf pi.IBuffer) {
	w, h := termbox.Size()
	var modeString string
	switch buf.Cursor().Mode() {
	case pi.MODE_NORMAL:
		modeString = "normal"
	case pi.MODE_EDIT:
		modeString = "edit"
	}

	pi.Puts(2, h-2, fmt.Sprintf(fmt.Sprintf("%%-%ds%%s", w-14), statusLine, modeString), termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
}
