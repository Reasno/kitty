package handlers

import (
	"context"
	"fmt"
	"github.com/Reasno/kitty/app/entity"
	"github.com/Reasno/kitty/app/msg"
	"github.com/Reasno/kitty/pkg/contract"
	kittyjwt "github.com/Reasno/kitty/pkg/jwt"
	"github.com/Reasno/kitty/pkg/wechat"
	pb "github.com/Reasno/kitty/proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type appService struct {
	log      log.Logger
	ur       UserRepository
	cr       CodeRepository
	sender   contract.SmsSender
	wechat   *wechat.Transport
	uploader contract.Uploader
}

type CodeRepository interface {
	CheckCode(ctx context.Context, mobile, code string) (bool, error)
	AddCode(ctx context.Context, mobile string) (code string, err error)
	DeleteCode(ctx context.Context, mobile string) (err error)
}

type UserRepository interface {
	GetFromWechat(ctx context.Context, wechat string, device *entity.Device, wechatUser entity.User) (user *entity.User, err error)
	GetFromMobile(ctx context.Context, mobile string, device *entity.Device) (user *entity.User, err error)
}

func (s appService) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginReply, error) {
	if len(in.Wechat) == 0 && len(in.Mobile) == 0 {
		return nil, status.Error(codes.InvalidArgument, msg.INVALID_PARAMS)
	}

	var (
		u      *entity.User
		device *entity.Device
		err    error
	)

	device = &entity.Device{
		Os:        uint8(in.Device.Os),
		Imei:      in.Device.Imei,
		Idfa:      in.Device.Idfa,
		Suuid:     in.Device.Suuid,
		Oaid:      in.Device.Oaid,
		Mac:       in.Device.Mac,
		AndroidId: in.Device.AndroidId,
	}
	if len(in.Wechat) != 0 {
		u, err = s.handleWechatLogin(ctx, in.Wechat, device)
	} else {
		u, err = s.handleMobileLogin(ctx, in.Mobile, in.Code, device)
	}
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Create jwt token
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		kittyjwt.NewClaim(
			uint64(u.ID),
			viper.GetString("app_name"),
			in.Device.Suuid,
			in.Channel,
			in.VersionCode,
			in.Wechat,
			time.Hour*24*30,
		),
	)
	token.Header["kid"] = viper.GetString("security.kid")
	tokenString, err := token.SignedString(viper.GetString("security.key"))

	var resp = pb.UserLoginReply{
		Id:       uint64(u.ID),
		UserName: u.UserName,
		Wechat:   u.Wechat,
		HeadImg:  u.HeadImg,
		Gender:   pb.Gender(u.Gender),
		Birthday: u.Birthday,
		Token:    tokenString,
	}

	return &resp, nil
}

func (s appService) verify(ctx context.Context, mobile string, code string) bool {
	result, err := s.cr.CheckCode(ctx, mobile, code)
	if err != nil {
		level.Error(s.log).Log("err", err)
	}
	if result {
		err = s.cr.DeleteCode(ctx, mobile)
		if err != nil {
			s.error(err)
		}
	}
	return result
}

func (s appService) GetCode(ctx context.Context, in *pb.GetCodeRequest) (*pb.GenericReply, error) {
	code, err := s.cr.AddCode(ctx, in.Mobile)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get code")
	}
	err = s.sender.Send(ctx, in.Mobile, code)
	if err != nil {
		return nil, errors.Wrap(err, "unable to send code")
	}
	var resp = pb.GenericReply{
		Code: 0,
	}
	return &resp, nil
}

func (s appService) debug(msg string, args ...interface{}) {
	level.Debug(s.log).Log("msg", fmt.Sprintf(msg, args...))
}
func (s appService) error(err error) {
	level.Error(s.log).Log("err", err)
}

func (s appService) handleWechatLogin(ctx context.Context, wechat string, device *entity.Device) (*entity.User, error) {
	wxRes, err := s.wechat.GetWechatLoginResponse(ctx, wechat)
	if err != nil {
		return nil, errors.Wrap(err, msg.ERROR_WECHAT_FAILUER)
	}
	if wxRes.Openid == "" {
		err := errors.New(msg.ERROR_MISSING_OPENID)
		return nil, errors.Wrap(err, msg.ERROR_WECHAT_FAILUER)
	}
	wxInfo, err := s.wechat.GetWechatUserInfoResult(ctx, wxRes)
	if err != nil {
		return nil, errors.Wrap(err, msg.ERROR_WECHAT_FAILUER)
	}
	headImg, err := s.uploader.UploadFromUrl(ctx, wxInfo.Headimgurl)
	if err != nil {
		return nil, errors.Wrap(err, msg.ERROR_UPLOAD)
	}
	wechatUser := entity.User{
		UserName: wxInfo.Nickname,
		HeadImg:  headImg,
		Wechat:   wxInfo.Openid,
	}
	u, err := s.ur.GetFromWechat(ctx, wxInfo.Openid, device, wechatUser)
	if err != nil {
		return nil, errors.Wrap(err, msg.ERROR_WECHAT_FAILUER)
	}
	return u, nil
}

func (s appService) handleMobileLogin(ctx context.Context, mobile, code string, device *entity.Device) (*entity.User, error) {
	if len(code) == 0 {
		return nil, status.Error(codes.InvalidArgument, msg.INVALID_PARAMS)
	}
	if !s.verify(ctx, mobile, code) {
		return nil, status.Error(codes.Unauthenticated, msg.ERROR_MOBILE_NOEXIST)
	}
	return s.ur.GetFromMobile(ctx, mobile, device)
}

func (s appService) GetInfo(ctx context.Context, in *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	var resp pb.UserInfoReply
	return &resp, nil
}

func (s appService) UpdateInfo(ctx context.Context, in *pb.UserInfoUpdateRequest) (*pb.UserInfoReply, error) {
	userId := kittyjwt.GetClaim(ctx).Uid
	var resp pb.UserInfoReply
	resp.Id = userId
	return &resp, nil
}
