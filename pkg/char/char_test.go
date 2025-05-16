package char

import (
	"testing"
	"unicode"
)

func Test(t *testing.T) {
	input := rune('️')
	input = '1'
	input = '\n'

	t.Logf("%x", input)

	t.Log("IsPrint", unicode.IsPrint(input))
	t.Log("IsControl", unicode.IsControl(input))
	t.Log("IsGraphic", unicode.IsGraphic(input))
	t.Log("IsSpace", unicode.IsSpace(input))
	t.Log("IsMark", unicode.IsMark(input))

}

func TestFilter(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "换行符",
			in:   "\n",
			out:  "",
		},
		{
			name: "空白字符",
			in:   " ",
			out:  " ",
		},
		{
			name: "未知的字符",
			in:   "️",
			out:  "",
		},
		{
			name: "未知的字符",
			in:   "\u203C",
			out:  "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if out := Filter(test.in); out != test.out {
				t.Errorf("Filter(%s) = %s, want %s", test.in, out, test.out)
			}
		})
	}
}
