//+build !wireinject

package service

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
	"glab.tagtic.cn/ad_gains/kitty/rule/entity"
	"glab.tagtic.cn/ad_gains/kitty/rule/service/mocks"
)

type mockDmpServer struct {
}

func (m mockDmpServer) UserMore(ctx context.Context, req *pb.DmpReq) (*pb.DmpResp, error) {
	return &pb.DmpResp{
		AdClick:    100,
		AdComplete: 0,
		AdDisplay:  0,
		AdCtrDev:   0,
		Register:   "2020-01-01 00:00:00",
		Score:      0,
		ScoreTotal: 0,
		BlackType:  0,
		Ext:        "",
	}, nil
}

func TestService_CalculateRules(t *testing.T) {
	cases := []struct {
		text    string
		payload dto.Payload
		result  dto.Data
	}{
		{
			`
style: advanced
rule:
  - if: true
    then: 
      foo: bar
`,
			dto.Payload{},
			dto.Data{"foo": "bar"},
		},
		{
			`
style: advanced
rule:
  - if: false
    then: 
      foo: bar
`,
			dto.Payload{},
			dto.Data{},
		},
		{
			`
style: advanced
rule:
  - if: Imei == "456"
    then: 
      foo: bar
  - if: Imei == "123"
    then:
      foo: baz
`,
			dto.Payload{
				Imei: "123",
			},
			dto.Data{
				"foo": "baz",
			},
		},
		{
			`
style: advanced
rule:
- if: Imei == "456" && Oaid == "789"
  then: 
    foo: bar
- if: Imei == "123" && Oaid == "789"
  then:
    foo: baz
- if: Imei == "123"
  then:
    foo: quz
`,
			dto.Payload{
				Imei: "123",
				Oaid: "789",
			},
			dto.Data{
				"foo": "baz",
			},
		},
		{
			`
style: advanced
enrich: true
rule:
- if: DMP.AdClick > 10
  then: 
    foo: bar
- if: true
  then:
    foo: quz
`,
			dto.Payload{
				Imei: "123",
				Oaid: "789",
			},
			dto.Data{
				"foo": "bar",
			},
		},
		{
			`
style: advanced
enrich: true
rule:
- if: HoursAgo(DMP.Register) < 1
  then: 
    foo: foo
- if: HoursAgo(DMP.Register) > 100
  then: 
    foo: bar
- if: true
  then:
    foo: quz
`,
			dto.Payload{
				Imei: "123",
				Oaid: "789",
			},
			dto.Data{
				"foo": "bar",
			},
		},
	}
	for _, c := range cases {
		cc := c
		t.Run("", func(t *testing.T) {
			repo := &mocks.Repository{}
			ser := ProvideService(log.NewNopLogger(), repo, mockDmpServer{})
			repo.On("GetCompiled", mock.Anything).Return(entity.NewRules(bytes.NewReader([]byte(cc.text)), log.NewNopLogger()))
			result, err := ser.CalculateRules(context.Background(), "", &cc.payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(result, cc.result) {
				t.Fatalf("want %#v, got %#v", cc.result, result)
			}
		})
	}
}

func TestService_UpdateRules(t *testing.T) {
	repo := &mocks.Repository{}
	ser := NewService(log.NewNopLogger(), repo)
	repo.On("SetRaw", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	repo.On("ValidateRules", mock.Anything, mock.Anything).Return(func(ruleName string, reader io.Reader) error {
		return entity.ValidateRules(reader)
	})
	err := ser.UpdateRules(context.Background(), "foo", []byte("invalid"), false)
	if err == nil {
		t.Fatal("err should not be null")
	}
	data := []byte(`
style: advanced
rule:
- if: true
  then:
    data: ok
`)
	err = ser.UpdateRules(context.Background(), "foo", data, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_Preflight(t *testing.T) {
	repo := &mocks.Repository{}
	ser := NewService(log.NewNopLogger(), repo)

	{
		repo.On("IsNewest", mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
		err := ser.Preflight(context.Background(), "foo", "fooo")
		if err != nil {
			t.Fatal("err should be null")
		}
	}

	{
		repo.On("IsNewest", mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()
		err := ser.Preflight(context.Background(), "foo", "fooo")
		if err == nil {
			t.Fatal("err should not be null")
		}
	}
}

func TestService_GetRules(t *testing.T) {
	repo := &mocks.Repository{}
	ser := NewService(log.NewNopLogger(), repo)
	{
		repo.On("GetRaw", mock.Anything, mock.Anything).Return([]byte("foo"), nil).Once()
		byt, err := ser.GetRules(context.Background(), "foo")
		assert.Nil(t, err)
		assert.Equal(t, byt, []byte("foo"))
	}
}
