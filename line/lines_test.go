package lines

import (
	"testing"

	"github.com/ac7/pi/interfaces"
)

func TestDisplayWidth(t *testing.T) {
	cases := []struct {
		s        string
		expected int
	}{
		{`	hi`, pi.TabWidth + 2},
		{`	`, pi.TabWidth},
		{`word`, 4},
		{`word	`, 4 + pi.TabWidth},
	}

	for _, c := range cases {
		result := DisplayWidth(c.s)
		if result != c.expected {
			t.Errorf("Recieved: %d, expected: %d", result, c.expected)
		}
	}
}
