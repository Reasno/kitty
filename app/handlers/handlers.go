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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type appService struct {
	conf     contract.ConfigReader
	log      log.Logger
	ur       UserRepository
	cr       CodeRepository
	sender   contract.SmsSender
	wechat   *wechat.Transport
	uploader contract.Uploader
	fr       FileRepository
}

type CodeRepository interface {
	CheckCode(ctx context.Context, mobile, code string) (bool, error)
	AddCode(ctx context.Context, mobile string) (code string, err error)
	DeleteCode(ctx context.Context, mobile string) (err error)
}

type UserRepository interface {
	GetFromWechat(ctx context.Context, wechat string, device *entity.Device, wechatUser entity.User) (user *entity.User, err error)
	GetFromMobile(ctx context.Context, mobile string, device *entity.Device) (user *entity.User, err error)
	GetFromDevice(ctx context.Context, suuid string, device *entity.Device) (user *entity.User, err error)
	Update(ctx context.Context, id uint, user entity.User) (newUser *entity.User, err error)
	Get(ctx context.Context, id uint) (user *entity.User, err error)
}

type FileRepository interface {
	UploadFromUrl(ctx context.Context, oldUrl string) (newUrl string, err error)
}

func (s appService) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserInfoReply, error) {
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
	if len(in.Wechat) == 0 && len(in.Mobile) == 0 {
		u, err = s.handleDeviceLogin(ctx, device.Suuid, device)
	} else if len(in.Wechat) != 0 {
		u, err = s.handleWechatLogin(ctx, in.Wechat, device)
	} else {
		u, err = s.handleMobileLogin(ctx, in.Mobile, in.Code, device)
	}
	if err != nil {
		return nil, errors.Wrap(err, "login failed")
	}

	// Create jwt token
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		kittyjwt.NewClaim(
			uint64(u.ID),
			s.conf.GetString("name"),
			in.Device.Suuid,
			in.Channel,
			in.VersionCode,
			in.Wechat,
			time.Hour*24*30,
		),
	)
	token.Header["kid"] = s.conf.GetString("security.kid")
	tokenString, err := token.SignedString([]byte(s.conf.GetString("security.key")))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create jwt token")
	}

	var resp = pb.UserInfoReply{
		Code: 0,
		Data: &pb.UserInfo{
			Id:       uint64(u.ID),
			UserName: u.UserName,
			Wechat:   u.WechatOpenId,
			HeadImg:  u.HeadImg,
			Gender:   pb.Gender(u.Gender),
			Birthday: u.Birthday,
			Token:    tokenString,
		},
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
	headImg, err := s.fr.UploadFromUrl(ctx, wxInfo.Headimgurl)
	if err != nil {
		return nil, errors.Wrap(err, msg.ERROR_UPLOAD)
	}
	wechatUser := entity.User{
		UserName:      wxInfo.Nickname,
		HeadImg:       headImg,
		WechatOpenId:  wxInfo.Openid,
		WechatUnionId: wxInfo.Unionid,
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

func (s appService) handleDeviceLogin(ctx context.Context, suuid string, device *entity.Device) (*entity.User, error) {
	return s.ur.GetFromDevice(ctx, suuid, device)
}

func (s appService) GetInfo(ctx context.Context, in *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	if in.Id == 0 {
		claim := kittyjwt.GetClaim(ctx)
		if claim == nil {
			return nil, status.Error(codes.Unauthenticated, msg.ERROR_NEED_LOGIN)
		}
		in.Id = claim.Uid
	}
	u, err := s.ur.Get(ctx, uint(in.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	var resp = pb.UserInfoReply{
		Code:    0,
		Message: "",
		Data: &pb.UserInfo{
			Id:       uint64(u.ID),
			UserName: u.UserName,
			Wechat:   u.WechatOpenId,
			HeadImg:  u.HeadImg,
			Gender:   pb.Gender(u.Gender),
			Birthday: u.Birthday,
		},
	}
	return &resp, nil
}

func (s appService) UpdateInfo(ctx context.Context, in *pb.UserInfoUpdateRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.GetClaim(ctx)
	if claim == nil {
		return nil, status.Error(codes.Unauthenticated, msg.ERROR_NEED_LOGIN)
	}
	u, err := s.ur.Update(ctx, uint(claim.Uid), entity.User{
		UserName: in.UserName,
		HeadImg:  in.HeadImg,
		Gender:   int(in.Gender),
		Birthday: in.Birthday,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	var resp = pb.UserInfoReply{
		Code:    0,
		Message: "",
		Data: &pb.UserInfo{
			Id:       uint64(u.ID),
			UserName: u.UserName,
			Wechat:   u.WechatOpenId,
			HeadImg:  u.HeadImg,
			Gender:   pb.Gender(u.Gender),
			Birthday: u.Birthday,
		},
	}
	return &resp, nil
}
