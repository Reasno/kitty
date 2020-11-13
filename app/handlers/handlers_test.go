package handlers

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/mock"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"glab.tagtic.cn/ad_gains/kitty/app/handlers/mocks"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	mc "glab.tagtic.cn/ad_gains/kitty/pkg/contract/mocks"
	"glab.tagtic.cn/ad_gains/kitty/pkg/wechat"
	wm "glab.tagtic.cn/ad_gains/kitty/pkg/wechat/mocks"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

func getHandler() appService {
	return appService{
		conf:   &mc.ConfigReader{},
		logger: log.NewNopLogger(),
		ur:     &mocks.UserRepository{},
		cr:     &mocks.CodeRepository{},
		er:     &mocks.ExtraRepository{},
		sender: &mc.SmsSender{},
		wechat: &wm.Wechater{},
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		service appService
		in      pb.UserLoginRequest
		out     pb.UserInfoReply
	}{
		{
			"手机登陆",
			appService{
				conf: (func() contract.ConfigReader {
					conf := &mc.ConfigReader{}
					conf.On("String", "name").Return("foo", nil)
					conf.On("String", "security.kid").Return("foo", nil)
					conf.On("String", "security.key").Return("foo", nil)
					return conf
				})(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					ur.On("GetFromMobile", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, packageName, mobile string, device *entity.Device) *entity.User {
						return &entity.User{Mobile: ns(mobile), PackageName: packageName}
					}, nil)
					//cr.On("DeleteCode", mock.Anything, mock.Anything).Return(nil)
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					cr.On("CheckCode", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
					cr.On("DeleteCode", mock.Anything, mock.Anything).Return(nil)
					return cr
				})(),
				er: (func() ExtraRepository {
					er := &mocks.ExtraRepository{}
					er.On("Get", mock.Anything, mock.Anything, mock.Anything).Return([]byte{}, nil)
					er.AssertNotCalled(t, "Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
					return er
				})(),
				sender: &mc.SmsSender{},
				wechat: &wm.Wechater{},
			},

			pb.UserLoginRequest{
				Mobile: "13699179983",
				Code:   "666666",
				Wechat: "",
				Device: &pb.Device{},
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					Mobile:      "13699179983",
					WechatExtra: &pb.WechatExtra{},
					TaobaoExtra: &pb.TaobaoExtra{
						OpenId: "",
					},
				},
			},
		},
		{
			"微信登陆",
			appService{
				conf: (func() contract.ConfigReader {
					conf := &mc.ConfigReader{}
					conf.On("String", "name").Return("foo", nil)
					conf.On("String", "security.kid").Return("foo", nil)
					conf.On("String", "security.key").Return("foo", nil)
					return conf
				})(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					ur.On("GetFromWechat", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, packageName, wechat string, device *entity.Device, wechatUser entity.User) *entity.User {
						return &entity.User{Mobile: ns("000")}
					}, nil)
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					cr.On("CheckCode", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
					cr.On("DeleteCode", mock.Anything, mock.Anything).Return(nil)
					return cr
				})(),
				er: (func() ExtraRepository {
					er := &mocks.ExtraRepository{}
					er.On("Get", mock.Anything, mock.Anything, mock.Anything).Return([]byte{}, nil)
					er.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
					return er
				})(),
				sender: &mc.SmsSender{},
				wechat: (func() wechat.Wechater {
					m := &wm.Wechater{}
					m.On("GetLoginResponse", mock.Anything, mock.Anything).Return(&wechat.WxLoginResult{
						AccessToken: "foo",
						Openid:      "bar",
					}, nil)
					m.On("GetUserInfoResult", mock.Anything, mock.Anything).Return(&wechat.WxUserInfoResult{
						Openid: "bar",
					}, nil)
					return m
				})(),
			},

			pb.UserLoginRequest{
				Wechat: "fff",
				Device: &pb.Device{},
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					Mobile: "000",
					WechatExtra: &pb.WechatExtra{
						OpenId: "bar",
					},
					TaobaoExtra: &pb.TaobaoExtra{
						OpenId: "",
					},
				},
			},
		},
		{
			"设备登陆",
			appService{
				conf: (func() contract.ConfigReader {
					conf := &mc.ConfigReader{}
					conf.On("String", "name").Return("foo", nil)
					conf.On("String", "security.kid").Return("foo", nil)
					conf.On("String", "security.key").Return("foo", nil)
					return conf
				})(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					ur.On("GetFromDevice", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, packageName, suuid string, device *entity.Device) *entity.User {
						return &entity.User{Mobile: ns("000"), CommonSUUID: "123"}
					}, nil)
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					cr.On("CheckCode", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
					cr.On("DeleteCode", mock.Anything, mock.Anything).Return(nil)
					return cr
				})(),
				er: (func() ExtraRepository {
					er := &mocks.ExtraRepository{}
					er.On("Get", mock.Anything, mock.Anything, mock.Anything).Return([]byte{}, nil)
					er.AssertNotCalled(t, "Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
					return er
				})(),
				sender: &mc.SmsSender{},
				wechat: (func() wechat.Wechater {
					m := &wm.Wechater{}
					m.On("GetLoginResponse", mock.Anything, mock.Anything).Return(&wechat.WxLoginResult{
						AccessToken: "foo",
						Openid:      "bar",
					}, nil)
					m.On("GetUserInfoResult", mock.Anything, mock.Anything).Return(&wechat.WxUserInfoResult{
						Openid: "bar",
					}, nil)
					return m
				})(),
			},

			pb.UserLoginRequest{
				Device: &pb.Device{
					Suuid: "123",
				},
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					Mobile: "000",
					WechatExtra: &pb.WechatExtra{
						OpenId: "",
					},
					TaobaoExtra: &pb.TaobaoExtra{
						OpenId: "",
					},
				},
			},
		},
		{
			"设备登陆时关联已绑定的淘宝号和微信号",
			appService{
				conf: (func() contract.ConfigReader {
					conf := &mc.ConfigReader{}
					conf.On("String", "name").Return("foo", nil)
					conf.On("String", "security.kid").Return("foo", nil)
					conf.On("String", "security.key").Return("foo", nil)
					return conf
				})(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					ur.On("GetFromDevice", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, packageName, suuid string, device *entity.Device) *entity.User {
						return &entity.User{Mobile: ns("000"), CommonSUUID: "123"}
					}, nil)
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					cr.On("CheckCode", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
					cr.On("DeleteCode", mock.Anything, mock.Anything).Return(nil)
					return cr
				})(),
				er: (func() ExtraRepository {
					er := &mocks.ExtraRepository{}
					taobao, _ := (&pb.TaobaoExtra{OpenId: "baz"}).Marshal()
					wc, _ := (&pb.WechatExtra{OpenId: "bar"}).Marshal()
					er.On("Get", mock.Anything, mock.Anything, pb.Extra_WECHAT_EXTRA.String()).Return(wc, nil)
					er.On("Get", mock.Anything, mock.Anything, pb.Extra_TAOBAO_EXTRA.String()).Return(taobao, nil)
					er.AssertNotCalled(t, "Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
					return er
				})(),
				sender: &mc.SmsSender{},
				wechat: (func() wechat.Wechater {
					m := &wm.Wechater{}
					m.On("GetLoginResponse", mock.Anything, mock.Anything).Return(&wechat.WxLoginResult{
						AccessToken: "foo",
						Openid:      "bar",
					}, nil)
					m.On("GetUserInfoResult", mock.Anything, mock.Anything).Return(&wechat.WxUserInfoResult{
						Openid: "bar",
					}, nil)
					return m
				})(),
			},

			pb.UserLoginRequest{
				Device: &pb.Device{
					Suuid: "123",
				},
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					Mobile: "000",
					WechatExtra: &pb.WechatExtra{
						OpenId: "bar",
					},
					TaobaoExtra: &pb.TaobaoExtra{
						OpenId: "baz",
					},
				},
			},
		},
	}
	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			out, err := cc.service.Login(context.Background(), &cc.in)
			if err != nil {
				t.Fatal(err)
			}
			if out.Code != cc.out.Code {
				t.Fatalf("want %d, got %d", cc.out.Code, out.Code)
			}
			if out.Data.WechatExtra.OpenId != cc.out.Data.WechatExtra.OpenId {
				t.Fatalf("want %s, got %s", cc.out.Data.WechatExtra.OpenId, out.Data.WechatExtra.OpenId)
			}
			if out.Data.Mobile != redact(cc.out.Data.Mobile) {
				t.Fatalf("want %s, got %s", redact(cc.out.Data.Mobile), out.Data.Mobile)
			}
			if out.Data.TaobaoExtra.OpenId != redact(cc.out.Data.TaobaoExtra.OpenId) {
				t.Fatalf("want %s, got %s", redact(cc.out.Data.Mobile), out.Data.Mobile)
			}

		})
	}
}

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
