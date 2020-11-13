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

func getConf() contract.ConfigReader {
	conf := &mc.ConfigReader{}
	conf.On("String", "name").Return("foo", nil)
	conf.On("String", "security.kid").Return("foo", nil)
	conf.On("String", "security.key").Return("foo", nil)
	return conf
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
				conf:   getConf(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					ur.On("GetFromMobile", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, packageName, mobile string, device *entity.Device) *entity.User {
						return &entity.User{Mobile: ns(mobile), PackageName: packageName}
					}, nil).Once()
					//cr.On("DeleteCode", mock.Anything, mock.Anything).Return(nil)
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					cr.On("CheckCode", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
					cr.On("DeleteCode", mock.Anything, mock.Anything).Return(nil)
					return cr
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
				conf:   getConf(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					ur.On("GetFromWechat", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, packageName, wechat string, device *entity.Device, wechatUser entity.User) *entity.User {
						return &entity.User{Mobile: ns("000"), WechatOpenId: ns(wechat), WechatExtra: wechatUser.WechatExtra}
					}, nil).Once()
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					return cr
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
					Wechat: "bar",
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
				conf:   getConf(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					ur.On("GetFromDevice", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, packageName, suuid string, device *entity.Device) *entity.User {
						return &entity.User{Mobile: ns("000"), CommonSUUID: "123"}
					}, nil).Once()
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					return cr
				})(),
				sender: &mc.SmsSender{},
				wechat: (func() wechat.Wechater {
					m := &wm.Wechater{}
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
				conf:   getConf(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					wc, _ := (&pb.WechatExtra{OpenId: "bar"}).Marshal()
					taobao, _ := (&pb.TaobaoExtra{OpenId: "baz"}).Marshal()
					ur.On("GetFromDevice", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, packageName, suuid string, device *entity.Device) *entity.User {
						return &entity.User{Mobile: ns("000"), CommonSUUID: "123", WechatExtra: wc, TaobaoExtra: taobao, WechatOpenId: ns("bar"), TaobaoOpenId: ns("baz")}
					}, nil).Once()
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					cr.On("CheckCode", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
					cr.On("DeleteCode", mock.Anything, mock.Anything).Return(nil)
					return cr
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
					Wechat: "bar",
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
			if out.Data.Wechat != cc.out.Data.Wechat {
				t.Fatalf("want %s, got %s", cc.out.Data.Wechat, out.Data.Wechat)
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
			if out.Data.TaobaoExtra.OpenId != cc.out.Data.TaobaoExtra.OpenId {
				t.Fatalf("want %s, got %s", cc.out.Data.TaobaoExtra.OpenId, out.Data.TaobaoExtra.OpenId)
			}
		})
	}
}

func TestBind(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		app  appService
		in   pb.UserBindRequest
		out  pb.UserInfoReply
	}{
		{
			"绑定淘宝",
			appService{
				conf:   getConf(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					ur.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, id uint, user entity.User) *entity.User {
						return &user
					}, nil).Once()
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					return cr
				})(),
				sender: &mc.SmsSender{},
				wechat: (func() wechat.Wechater {
					m := &wm.Wechater{}
					return m
				})(),
			},
			pb.UserBindRequest{
				TaobaoExtra: &pb.TaobaoExtra{
					OpenId: "foo",
				},
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					TaobaoExtra: &pb.TaobaoExtra{
						OpenId: "foo",
					},
					WechatExtra: &pb.WechatExtra{},
				},
			},
		},
		{
			"绑定微信",
			appService{
				conf:   getConf(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					ur.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, id uint, user entity.User) *entity.User {
						return &user
					}, nil).Once()
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					return cr
				})(),
				sender: &mc.SmsSender{},
				wechat: (func() wechat.Wechater {
					m := &wm.Wechater{}
					m.On("GetLoginResponse", mock.Anything, mock.Anything).Return(&wechat.WxLoginResult{
						Openid: "bar",
					}, nil)
					m.On("GetUserInfoResult", mock.Anything, mock.Anything).Return(&wechat.WxUserInfoResult{
						Openid: "bar",
					}, nil)
					return m
				})(),
			},
			pb.UserBindRequest{
				Wechat: "foo",
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					Wechat:      "bar",
					TaobaoExtra: &pb.TaobaoExtra{},
					WechatExtra: &pb.WechatExtra{
						OpenId: "bar",
					},
				},
			},
		},
		{
			"绑定微信OpenId",
			appService{
				conf:   getConf(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					ur.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, id uint, user entity.User) *entity.User {
						return &user
					}, nil).Once()
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					return cr
				})(),
				sender: &mc.SmsSender{},
				wechat: (func() wechat.Wechater {
					m := &wm.Wechater{}
					return m
				})(),
			},
			pb.UserBindRequest{
				WechatExtra: &pb.WechatExtra{OpenId: "bar"},
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					Wechat:      "bar",
					TaobaoExtra: &pb.TaobaoExtra{},
					WechatExtra: &pb.WechatExtra{
						OpenId: "bar",
					},
				},
			},
		},
		{
			"绑定手机",
			appService{
				conf:   getConf(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					ur.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, id uint, user entity.User) *entity.User {
						return &user
					}, nil).Once()
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					cr.On("CheckCode", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
					cr.On("DeleteCode", mock.Anything, mock.Anything).Return(nil)
					return cr
				})(),
				sender: &mc.SmsSender{},
				wechat: (func() wechat.Wechater {
					m := &wm.Wechater{}
					return m
				})(),
			},
			pb.UserBindRequest{
				Mobile: "000",
				Code:   "666666",
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					Mobile:      "000",
					TaobaoExtra: &pb.TaobaoExtra{},
					WechatExtra: &pb.WechatExtra{},
				},
			},
		},
		{
			"同时绑定多个属性",
			appService{
				conf:   getConf(),
				logger: log.NewNopLogger(),
				ur: (func() UserRepository {
					ur := &mocks.UserRepository{}
					ur.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, id uint, user entity.User) *entity.User {
						return &user
					}, nil).Once()
					return ur
				})(),
				cr: (func() CodeRepository {
					cr := &mocks.CodeRepository{}
					cr.On("CheckCode", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
					cr.On("DeleteCode", mock.Anything, mock.Anything).Return(nil)
					return cr
				})(),
				sender: &mc.SmsSender{},
				wechat: (func() wechat.Wechater {
					m := &wm.Wechater{}
					m.On("GetLoginResponse", mock.Anything, mock.Anything).Return(&wechat.WxLoginResult{
						Openid: "bar",
					}, nil)
					m.On("GetUserInfoResult", mock.Anything, mock.Anything).Return(&wechat.WxUserInfoResult{
						Openid: "bar",
					}, nil)
					return m
				})(),
			},
			pb.UserBindRequest{
				Mobile: "000",
				Code:   "666666",
				TaobaoExtra: &pb.TaobaoExtra{
					OpenId: "baz",
				},
				Wechat: "foo",
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					Mobile: "000",
					Wechat: "bar",
					TaobaoExtra: &pb.TaobaoExtra{
						OpenId: "baz",
					},
					WechatExtra: &pb.WechatExtra{
						OpenId: "bar",
					},
				},
			},
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			out, err := cc.app.Bind(context.Background(), &cc.in)
			if err != nil {
				t.Fatal(err)
			}
			if out.Code != cc.out.Code {
				t.Fatalf("want %d, got %d", cc.out.Code, out.Code)
			}
			if out.Data.Wechat != cc.out.Data.Wechat {
				t.Fatalf("want %s, got %s", cc.out.Data.Wechat, out.Data.Wechat)
			}
			if out.Data.WechatExtra.OpenId != cc.out.Data.WechatExtra.OpenId {
				t.Fatalf("want %s, got %s", cc.out.Data.WechatExtra.OpenId, out.Data.WechatExtra.OpenId)
			}
			if out.Data.Mobile != redact(cc.out.Data.Mobile) {
				t.Fatalf("want %s, got %s", redact(cc.out.Data.Mobile), out.Data.Mobile)
			}
			if out.Data.TaobaoExtra.OpenId != cc.out.Data.TaobaoExtra.OpenId {
				t.Fatalf("want %s, got %s", cc.out.Data.TaobaoExtra.OpenId, out.Data.TaobaoExtra.OpenId)
			}
		})
	}
}

func TestUnbind(t *testing.T) {
	t.Parallel()
	app := appService{
		conf:   getConf(),
		logger: log.NewNopLogger(),
		ur: (func() UserRepository {
			ur := &mocks.UserRepository{}
			ur.On("Save", mock.Anything, mock.Anything).Return(nil)
			ur.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, id uint) *entity.User {
				taobao, _ := (&pb.TaobaoExtra{OpenId: "1"}).Marshal()
				wc, _ := (&pb.WechatExtra{OpenId: "1"}).Marshal()
				return &entity.User{
					UserName:      "1",
					WechatOpenId:  ns("1"),
					WechatUnionId: ns("1"),
					Mobile:        ns("1"),
					TaobaoOpenId:  ns("1"),
					WechatExtra:   wc,
					TaobaoExtra:   taobao,
				}
			}, nil)
			return ur
		})(),
		cr: (func() CodeRepository {
			cr := &mocks.CodeRepository{}
			return cr
		})(),
		sender: &mc.SmsSender{},
		wechat: (func() wechat.Wechater {
			m := &wm.Wechater{}
			return m
		})(),
	}
	cases := []struct {
		name string
		app  appService
		in   pb.UserUnbindRequest
		out  pb.UserInfoReply
	}{
		{
			"解绑淘宝",
			app,
			pb.UserUnbindRequest{
				Taobao: true,
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					Mobile: "1",
					Wechat: "1",
					WechatExtra: &pb.WechatExtra{
						OpenId: "1",
					},
					TaobaoExtra: &pb.TaobaoExtra{},
				},
			},
		},
		{
			"解绑手机",
			app,
			pb.UserUnbindRequest{
				Mobile: true,
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					Wechat: "1",
					WechatExtra: &pb.WechatExtra{
						OpenId: "1",
					},
					TaobaoExtra: &pb.TaobaoExtra{
						OpenId: "1",
					},
				},
			},
		},
		{
			"解绑微信",
			app,
			pb.UserUnbindRequest{
				Wechat: true,
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					Mobile:      "1",
					WechatExtra: &pb.WechatExtra{},
					TaobaoExtra: &pb.TaobaoExtra{
						OpenId: "1",
					},
				},
			},
		},
		{
			"全解绑",
			app,
			pb.UserUnbindRequest{
				Wechat: true,
				Mobile: true,
				Taobao: true,
			},
			pb.UserInfoReply{
				Code: 0,
				Data: &pb.UserInfo{
					WechatExtra: &pb.WechatExtra{},
					TaobaoExtra: &pb.TaobaoExtra{},
				},
			},
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			out, err := cc.app.Unbind(context.Background(), &cc.in)
			if err != nil {
				t.Fatal(err)
			}
			if out.Code != cc.out.Code {
				t.Fatalf("want %d, got %d", cc.out.Code, out.Code)
			}
			if out.Data.Wechat != cc.out.Data.Wechat {
				t.Fatalf("want %s, got %s", cc.out.Data.Wechat, out.Data.Wechat)
			}
			if out.Data.WechatExtra.OpenId != cc.out.Data.WechatExtra.OpenId {
				t.Fatalf("want %s, got %s", cc.out.Data.WechatExtra.OpenId, out.Data.WechatExtra.OpenId)
			}
			if out.Data.Mobile != redact(cc.out.Data.Mobile) {
				t.Fatalf("want %s, got %s", redact(cc.out.Data.Mobile), out.Data.Mobile)
			}
			if out.Data.TaobaoExtra.OpenId != cc.out.Data.TaobaoExtra.OpenId {
				t.Fatalf("want %s, got %s", cc.out.Data.TaobaoExtra.OpenId, out.Data.TaobaoExtra.OpenId)
			}
		})
	}
}

func TestWireType(t *testing.T) {
	byt, err := (&pb.WechatExtra{OpenId: "1"}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var out pb.WechatExtra
	err = out.Unmarshal(byt)
	if err != nil {
		t.Fatal(err)
	}
	if out.OpenId != "1" {
		t.Fatalf("want 1 got %s", out.OpenId)
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
