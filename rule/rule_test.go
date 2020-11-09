package rule

import (
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
)

func TestNewRules(t *testing.T) {
	cases := []struct {
		raw      string
		expected []Rule
	}{{
		`
style: basic
rule:
  foo: bar
`,
		[]Rule{
			{
				If: "true",
				Then: Data{
					"foo": "bar",
				},
			},
		},
	},
		{
			`
style: advanced
rule:
  - if: "true"
    then:
      foo: bar
`,
			[]Rule{
				{
					If: "true",
					Then: Data{
						"foo": "bar",
					},
				},
			},
		},
	}
	for _, c := range cases {
		cc := c
		t.Run("", func(t *testing.T) {
			rules := NewRules(strings.NewReader(cc.raw), log.NewNopLogger())
			for i, r := range rules {
				if r.If != cc.expected[i].If {
					t.Fatalf("want %s, got %s", cc.expected[i].If, r.If)
				}
				if len(r.Then) != len(cc.expected[i].Then) {
					t.Fatalf("want %+v, got %+v", cc.expected[i].Then, r.Then)
				}
			}
		})
	}

}
