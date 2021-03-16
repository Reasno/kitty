package entity

import (
	"testing"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/stretchr/testify/assert"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
)

func TestAdvancedRuleItem_Calculate(t *testing.T) {
	cases := []struct {
		name    string
		yaml    string
		payload *dto.Payload
		expect  func(*testing.T, error, dto.Data)
	}{
		{
			"error",
			`
style: basic
rule:
  - if: Channel > 0
    then:
      i: 1
  - if: true
    then:
      i: 2
`,
			&dto.Payload{Channel: "foo"},
			func(t *testing.T, err error, data dto.Data) {
				assert.Error(t, err)
			},
		},
		{
			"normal",
			`
style: basic
rule:
  - if: Channel == "foo"
    then:
      i: 1
  - if: true
    then:
      i: 2
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
style: basic
rule:
  - if: Channel == "foo"
    child:
      style: advanced
      rule:
        - if: PackageName == "bar"
          child:
            style: basic
            rule:
              i: 1
        - if: true
          then:
            i: 3
  - if: true
    then:
      i: 2
`,
			&dto.Payload{Channel: "foo", PackageName: "bar"},
			func(t *testing.T, err error, data dto.Data) {
				assert.NoError(t, err)
				assert.Equal(t, 1.0, data["i"])
			},
		},
	}
	for _, cc := range cases {
		c := cc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			ar := NewAdvancedRule()
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

func TestAdvancedRule_Unmarshal(t *testing.T) {
	cases := []struct {
		name   string
		yaml   string
		expect func(*testing.T, error, *AdvancedRuleCollection)
	}{
		{
			"simple",
			`
style: advanced
rule:
  - if: true
    then:
      foo: bar
`,
			func(t *testing.T, err error, rule *AdvancedRuleCollection) {
				assert.NoError(t, err)
				assert.Equal(t, "bar", rule.items[0].then["foo"])
			},
		},
		{
			"nested",
			`
style: advanced
rule:
  - if: true
    child:
      style: advanced
      rule:
        - if: 'true'
          then:
            foo: bar
`,
			func(t *testing.T, err error, rule *AdvancedRuleCollection) {
				assert.NoError(t, err)
				assert.Equal(t, "bar", rule.items[0].child.(*AdvancedRuleCollection).items[0].then["foo"])
			},
		},
		{
			"advance-basic",
			`
style: advanced
rule:
  - if: true
    child:
      style: basic
      rule:
        foo: "some text"
`,
			func(t *testing.T, err error, rule *AdvancedRuleCollection) {
				assert.NoError(t, err)
				assert.Equal(t, "some text", rule.items[0].child.(*BasicRule).data["foo"])
			},
		},
		{
			"deep nest",
			`
style: advanced
rule:
  - if: true
    child:
      style: advanced
      rule:
        - if: true
          child:
            style: advanced
            rule:
              - if: true
                then:
                  foo: bar
`,
			func(t *testing.T, err error, rule *AdvancedRuleCollection) {
				assert.NoError(t, err)
				assert.Equal(t, "bar", rule.items[0].child.(*AdvancedRuleCollection).items[0].child.(*AdvancedRuleCollection).items[0].then["foo"])
			},
		},
		{
			"multiple element",
			`
style: advanced
rule:
  - if: true
    child:
      style: advanced
      rule:
        - if: true
          then: {}
        - if: false
          then:
            foo: bar
`,
			func(t *testing.T, err error, rule *AdvancedRuleCollection) {
				assert.NoError(t, err)
				assert.Equal(t, "bar", rule.items[0].child.(*AdvancedRuleCollection).items[1].then["foo"])
			},
		},
	}

	for _, cc := range cases {
		c := cc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			ar := NewAdvancedRule()
			k := koanf.New(".")
			err := k.Load(rawbytes.Provider([]byte(c.yaml)), yaml.Parser())
			assert.NoError(t, err)
			err = ar.Unmarshal(k)
			c.expect(t, err, ar)
		})
	}
}
