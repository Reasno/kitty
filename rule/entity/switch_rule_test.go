package entity

import (
	"testing"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/stretchr/testify/assert"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
)

func TestSwitchRule_Calculate(t *testing.T) {
	cases := []struct {
		name    string
		yaml    string
		payload *dto.Payload
		expect  func(*testing.T, error, dto.Data)
	}{
		{
			"normal",
			`
style: switch
by: Channel
rule:
  - case: foo
    style: basic
    rule:
      i: 1
  - case: bar
    style: basic
    rule:
      i: 2
default:
  style: basic
  rule:
    i: 3
`,
			&dto.Payload{Channel: "foo"},
			func(t *testing.T, err error, data dto.Data) {
				assert.NoError(t, err)
				assert.Equal(t, 1.0, data["i"])
			},
		},
		{
			"normal",
			`
style: switch
by: Channel
rule:
  - case: foo
    style: basic
    rule:
      i: 1
  - case: bar
    style: basic
    rule:
      i: 2
default:
  style: basic
  rule:
    i: 3
`,
			&dto.Payload{Channel: "bar"},
			func(t *testing.T, err error, data dto.Data) {
				assert.NoError(t, err)
				assert.Equal(t, 2.0, data["i"])
			},
		},
		{
			"normal",
			`
style: switch
by: Channel
rule:
  - case: foo
    style: basic
    rule:
      i: 1
  - case: bar
    style: basic
    rule:
      i: 2
default:
  style: basic
  rule:
    i: 3
`,
			&dto.Payload{Channel: "baz"},
			func(t *testing.T, err error, data dto.Data) {
				assert.NoError(t, err)
				assert.Equal(t, 3.0, data["i"])
			},
		},
		{
			"repeat",
			`
style: switch
by: Channel
rule:
  - case: foo
    style: basic
    rule:
      i: 1
  - case: foo
    style: basic
    rule:
      i: 2
default:
  style: basic
  rule:
    i: 3
`,
			&dto.Payload{Channel: "foo"},
			func(t *testing.T, err error, data dto.Data) {
				assert.NoError(t, err)
				assert.Equal(t, 1.0, data["i"])
			},
		},
		{
			"nest",
			`
style: switch
by: Channel
rule:
  - case: foo
    style: advanced
    rule:
      - if: PackageName == "foo"
        then:
          i: 1
      - if: true
        then:
          i: 4
  - case: foo
    style: basic
    rule:
      i: 2
default:
  style: basic
  rule:
    i: 3
`,
			&dto.Payload{Channel: "foo"},
			func(t *testing.T, err error, data dto.Data) {
				assert.NoError(t, err)
				assert.Equal(t, 4.0, data["i"])
			},
		},
		{
			"nest",
			`
style: switch
by: Channel
rule:
  - case: foo
    style: advanced
    rule:
      - if: PackageName == "foo"
        then:
          i: 1
      - if: true
        then:
          i: 4
  - case: foo
    style: basic
    rule:
      i: 2
default:
  style: advanced
  rule:
    - if: false
      then:
        i: 5
    - if: true
      then:
        i: 6
`,
			&dto.Payload{Channel: ""},
			func(t *testing.T, err error, data dto.Data) {
				assert.NoError(t, err)
				assert.Equal(t, 6.0, data["i"])
			},
		},
	}

	for _, cc := range cases {
		c := cc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			ar := NewSwitchRule()
			k := koanf.New(".")
			err := k.Load(rawbytes.Provider([]byte(c.yaml)), yaml.Parser())
			assert.NoError(t, err)
			err = ar.Unmarshal(k)
			assert.NoError(t, err)
			err = ar.Compile()
			assert.NoError(t, err)
			result, err := ar.Calculate(c.payload)
			c.expect(t, err, result)
		})
	}

}
