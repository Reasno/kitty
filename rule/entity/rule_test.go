package entity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRules(t *testing.T) {
	cases := []struct {
		name    string
		rule    string
		asserts func(t *testing.T, err error)
	}{
		{
			"invalid",
			`
style: advanced
rule:
  - if: Channel > 0
    then:
      sms: 1
`,
			func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateRules(strings.NewReader(c.rule))
			c.asserts(t, err)
		})
	}
}
