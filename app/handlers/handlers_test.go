package handlers

import "testing"

func TestMobileRedact(t *testing.T) {
	cases := []struct {
		input  string
		expect string
	}{
		{
			"13799199999",
			"137****9999",
		},
		{
			"111",
			"111",
		},
		{
			"013799199999",
			"013****99999",
		},
	}
	for _, c := range cases {
		output := redact(c.input)
		if output != c.expect {
			t.Fatalf("want %s, got %s", c.expect, output)
		}
	}

}
