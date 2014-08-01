package lines

import (
	"testing"

	"github.com/ac7/pi"
)

func TestDisplayWidth(t *testing.T) {
	cases := []struct {
		s        string
		expected int
	}{
		{`	hi`, pi.TAB_WIDTH + 2},
		{`	`, pi.TAB_WIDTH},
		{`word`, 4},
		{`word	`, 4 + pi.TAB_WIDTH},
	}

	for _, c := range cases {
		result := DisplayWidth(c.s)
		if result != c.expected {
			t.Errorf("Recieved: %d, expected: %d", result, c.expected)
		}
	}
}
