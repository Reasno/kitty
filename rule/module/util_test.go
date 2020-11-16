package module

import "testing"

func TestLongestCommonPrefix(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input  []string
		output string
	}{
		{
			[]string{"foooooooo", "foo"},
			"foo",
		},
		{
			[]string{"foo"},
			"foo",
		},
		{
			[]string{"foo", "bar"},
			"",
		},
		{
			[]string{"/ad/one", "/ad/another"},
			"/ad/",
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.output, func(t *testing.T) {
			o := Prefix(cc.input)
			if o != cc.output {
				t.Fatalf("want %s, got %s", cc.output, o)
			}
		})
	}
}
