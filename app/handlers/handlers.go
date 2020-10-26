package handlers

import (
	"context"
	"fmt"
	"github.com/Reasno/kitty/app/entity"
	"github.com/Reasno/kitty/app/msg"
	"github.com/Reasno/kitty/pkg/contract"
	kittyjwt "github.com/Reasno/kitty/pkg/jwt"
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
	log    log.Logger
	ur     UserRepository
	cr     CodeRepository
	sender contract.SmsSender
}

type CodeRepository interface {
	CheckCode(ctx context.Context, mobile, code string) (bool, error)
	AddCode(ctx context.Context, mobile string) (code string, err error)
}

type UserRepository interface {
	GetFromWechat(ctx context.Context, wechat string, device *entity.Device) (user *entity.User, err error)
	GetFromMobile(ctx context.Context, mobile string, device *entity.Device) (user *entity.User, err error)
}

func (s appService) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginReply, error) {
	if len(in.Wechat) == 0 && len(in.Mobile) == 0 {
		return nil, status.Error(codes.InvalidArgument, msg.INVALID_PARAMS)
	}
	if len(in.Mobile) != 0 && len(in.Code) == 0 {
		return nil, status.Error(codes.InvalidArgument, msg.INVALID_PARAMS)
	}
	if len(in.Mobile) != 0 && !s.verify(ctx, in.Mobile, in.Code) {
		return nil, status.Error(codes.Unauthenticated, msg.ERROR_MOBILE_NOEXIST)
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
		u, err = s.ur.GetFromWechat(ctx, in.Wechat, device)
	} else {
		u, err = s.ur.GetFromMobile(ctx, in.Mobile, device)
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
