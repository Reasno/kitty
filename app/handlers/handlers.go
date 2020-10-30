package handlers

import (
	"context"
	"database/sql"
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
	Save(ctx context.Context, user *entity.User) error
}

type FileRepository interface {
	UploadFromUrl(ctx context.Context, oldUrl string) (newUrl string, err error)
}

func (s appService) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserInfoReply, error) {
	// TODO: 如果用户已经登陆了 就不能再登陆了，需要执行bind
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
	if len(in.Mobile) != 0 {
		u, err = s.handleMobileLogin(ctx, in.Mobile, in.Code, device)
	} else if len(in.Wechat) != 0 {
		u, err = s.handleWechatLogin(ctx, in.Wechat, device)
	} else {
		u, err = s.handleDeviceLogin(ctx, device.Suuid, device)
	}
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorLogin)
	}

	// Create jwt token
	tokenString, err := s.getToken(uint64(u.ID), in.Device.Suuid, in.Channel, in.VersionCode, u.WechatOpenId.String, u.Mobile.String)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create jwt token")
	}

	var resp = u.ToReply()
	resp.Data.Token = tokenString
	return resp, nil
}

func (s appService) getToken(userId uint64, suuid, channel, versionCode, wechat, mobile string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		kittyjwt.NewClaim(
			userId,
			s.conf.GetString("name"),
			suuid, channel, versionCode, wechat, mobile,
			s.conf.GetString("packageName"),
			time.Hour*24*30,
		),
	)
	token.Header["kid"] = s.conf.GetString("security.kid")
	return token.SignedString([]byte(s.conf.GetString("security.key")))
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

func (s appService) getWechatInfo(ctx context.Context, wechat string) (*wechat.WxUserInfoResult, error) {
	wxRes, err := s.wechat.GetWechatLoginResponse(ctx, wechat)
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorWechatFailure)
	}
	if wxRes.Openid == "" {
		err := errors.New(msg.ErrorMissingOpenid)
		return nil, errors.Wrap(err, msg.ErrorWechatFailure)
	}
	wxInfo, err := s.wechat.GetWechatUserInfoResult(ctx, wxRes)
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorWechatFailure)
	}
	return wxInfo, nil
}

func (s appService) handleWechatLogin(ctx context.Context, wechat string, device *entity.Device) (*entity.User, error) {
	wxInfo, err := s.getWechatInfo(ctx, wechat)
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorWechatFailure)
	}
	headImg, err := s.fr.UploadFromUrl(ctx, wxInfo.Headimgurl)
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorUpload)
	}
	wechatUser := entity.User{
		UserName:      wxInfo.Nickname,
		HeadImg:       headImg,
		WechatOpenId:  ns(wxInfo.Openid),
		WechatUnionId: ns(wxInfo.Unionid),
	}
	u, err := s.ur.GetFromWechat(ctx, wxInfo.Openid, device, wechatUser)
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorWechatFailure)
	}
	return u, nil
}

func (s appService) handleMobileLogin(ctx context.Context, mobile, code string, device *entity.Device) (*entity.User, error) {
	if len(code) == 0 {
		return nil, status.Error(codes.InvalidArgument, msg.InvalidParams)
	}
	if !s.verify(ctx, mobile, code) {
		return nil, status.Error(codes.Unauthenticated, msg.ErrorMobileCode)
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
			return nil, status.Error(codes.Unauthenticated, msg.ErrorNeedLogin)
		}
		in.Id = claim.UserId
	}
	u, err := s.ur.Get(ctx, uint(in.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return u.ToReply(), nil
}

func (s appService) UpdateInfo(ctx context.Context, in *pb.UserInfoUpdateRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.GetClaim(ctx)
	if claim == nil {
		return nil, status.Error(codes.Unauthenticated, msg.ErrorNeedLogin)
	}
	u, err := s.ur.Update(ctx, uint(claim.UserId), entity.User{
		UserName: in.UserName,
		HeadImg:  in.HeadImg,
		Gender:   int(in.Gender),
		Birthday: in.Birthday,
	})
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorDatabaseFailure)
	}

	return u.ToReply(), nil
}

func (s appService) Bind(ctx context.Context, in *pb.UserBindRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.GetClaim(ctx)
	if claim == nil {
		return nil, status.Error(codes.Unauthenticated, msg.ErrorNeedLogin)
	}

	var (
		toUpdate entity.User
	)

	// 绑定手机号
	if len(in.Mobile) > 0 && len(in.Code) > 0 {
		if !s.verify(ctx, in.Mobile, in.Code) {
			return nil, status.Error(codes.Unauthenticated, msg.ErrorMobileCode)
		}
		toUpdate = entity.User{Mobile: ns(in.Mobile)}
	}

	// 绑定微信号
	if len(in.Wechat) > 0 {
		wxInfo, err := s.getWechatInfo(ctx, in.Wechat)
		if err != nil {
			return nil, errors.Wrap(err, msg.ErrorWechatFailure)
		}
		toUpdate = entity.User{
			WechatOpenId:  ns(wxInfo.Openid),
			WechatUnionId: ns(wxInfo.Unionid),
		}
	}

	if len(in.OpenId) > 0 {
		toUpdate = entity.User{
			WechatOpenId: ns(in.OpenId),
		}
	}

	// 更新用户
	newUser, err := s.ur.Update(ctx, uint(claim.UserId), toUpdate)
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorDatabaseFailure)
	}
	reply := newUser.ToReply()
	reply.Data.Token, err = s.getToken(
		uint64(newUser.ID),
		newUser.CommonSUUID,
		newUser.Channel,
		newUser.VersionCode,
		newUser.WechatOpenId.String,
		newUser.Mobile.String,
	)
	if err != nil {
		err = errors.Wrap(err, msg.ErrorJwtFailure)
	}

	return reply, err
}

func (s appService) Unbind(ctx context.Context, in *pb.UserUnbindRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.GetClaim(ctx)
	if claim == nil {
		return nil, status.Error(codes.Unauthenticated, msg.ErrorNeedLogin)
	}
	user, err := s.ur.Get(ctx, uint(claim.UserId))
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorDatabaseFailure)
	}
	if in.Mobile {
		user.Mobile = sql.NullString{}
	}
	if in.Wechat {
		user.WechatUnionId = sql.NullString{}
		user.WechatOpenId = sql.NullString{}
	}
	err = s.ur.Save(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorDatabaseFailure)
	}
	return user.ToReply(), nil
}

func ns(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func (s appService) Refresh(ctx context.Context, in *pb.UserRefreshRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.GetClaim(ctx)
	if claim == nil {
		return nil, status.Error(codes.Unauthenticated, msg.ErrorNeedLogin)
	}
	device := &entity.Device{
		Os:        uint8(in.Device.Os),
		Imei:      in.Device.Imei,
		Idfa:      in.Device.Idfa,
		Suuid:     in.Device.Suuid,
		Oaid:      in.Device.Oaid,
		Mac:       in.Device.Mac,
		AndroidId: in.Device.AndroidId,
	}
	u, err := s.ur.Get(ctx, uint(claim.UserId))
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorDatabaseFailure)
	}

	u.CommonSUUID = in.Device.Suuid
	u.Channel = in.Channel
	u.VersionCode = in.VersionCode
	u.AddNewDevice(device)

	if err := s.ur.Save(ctx, u); err != nil {
		return nil, errors.Wrap(err, msg.ErrorDatabaseFailure)
	}

	reply := u.ToReply()
	reply.Data.Token, err = s.getToken(
		uint64(u.ID),
		u.CommonSUUID,
		u.Channel,
		u.VersionCode,
		u.WechatOpenId.String,
		u.Mobile.String,
	)
	if err != nil {
		err = errors.Wrap(err, msg.ErrorJwtFailure)
	}
	return reply, nil
}
