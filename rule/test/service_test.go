//+build !wireinject

package test

import (
	"bytes"
	"context"
	"github.com/Reasno/kitty/rule"
	"github.com/Reasno/kitty/rule/mocks"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestGetCompiled(t *testing.T) {
	cases := []struct {
		text    string
		payload rule.Payload
		result  rule.Data
	}{
		{
			`
style: advanced
rules:
  - if: true
    then: 
      foo: bar
`,
			rule.Payload{},
			rule.Data{"foo": "bar"},
		},
		{
			`
style: advanced
rules:
  - if: false
    then: 
      foo: bar
`,
			rule.Payload{},
			rule.Data{},
		},
		{
			`
style: advanced
rules:
  - if: Imei == "456"
    then: 
      foo: bar
  - if: Imei == "123"
    then:
      foo: baz
`,
			rule.Payload{
				Imei: "123",
			},
			rule.Data{
				"foo": "baz",
			},
		},
		{
			`
style: advanced
rules:
- if: Imei == "456" && Oaid = "789"
  then: 
    foo: bar
- if: Imei == "123" && Oaid == "789"
  then:
    foo: baz
- if: Imei == "123"
  then:
    foo: quz
`,
			rule.Payload{
				Imei: "123",
				Oaid: "789",
			},
			rule.Data{
				"foo": "baz",
			},
		},
	}
	for _, c := range cases {
		cc := c
		t.Run("", func(t *testing.T) {
			repo := &mocks.Repository{}
			ser := rule.NewService(log.NewNopLogger(), repo)
			repo.On("GetCompiled", mock.Anything).Return(rule.NewRules(bytes.NewReader([]byte(cc.text)), log.NewNopLogger()))
			result, err := ser.CalculateRules(context.Background(), "", &cc.payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(result, cc.result) {
				t.Fatalf("want %v, got %v", cc.result, result)
			}
		})
	}
}

func TestSet(t *testing.T) {
	repo := &mocks.Repository{}
	ser := rule.NewService(log.NewNopLogger(), repo)
	repo.On("SetRaw", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	err := ser.UpdateRules(context.Background(), "foo", []byte("invalid"), false)
	if err == nil {
		t.Fatal("err should not be null")
	}
	data := []byte(`
- if: true
  then:
    data: ok
`)
	err = ser.UpdateRules(context.Background(), "foo", data, false)
	if err != nil {
		t.Fatal(err)
	}
}
