package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Reasno/kitty/app/entity"
	"github.com/Reasno/kitty/app/msg"
	"github.com/Reasno/kitty/app/repository"
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/Reasno/kitty/pkg/kerr"
	kittyjwt "github.com/Reasno/kitty/pkg/kjwt"
	"github.com/Reasno/kitty/pkg/wechat"
	pb "github.com/Reasno/kitty/proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

var taobaoExtraKey struct{}
var wechatExtraKey struct{}

type appService struct {
	conf     contract.ConfigReader
	log      log.Logger
	ur       UserRepository
	cr       CodeRepository
	er       ExtraRepository
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
	GetFromWechat(ctx context.Context, packageName, wechat string, device *entity.Device, wechatUser entity.User) (user *entity.User, err error)
	GetFromMobile(ctx context.Context, packageName, mobile string, device *entity.Device) (user *entity.User, err error)
	GetFromDevice(ctx context.Context, packageName, suuid string, device *entity.Device) (user *entity.User, err error)
	Update(ctx context.Context, id uint, user entity.User) (newUser *entity.User, err error)
	Get(ctx context.Context, id uint) (user *entity.User, err error)
	Save(ctx context.Context, user *entity.User) error
}

type FileRepository interface {
	UploadFromUrl(ctx context.Context, oldUrl string) (newUrl string, err error)
}

type ExtraRepository interface {
	Put(ctx context.Context, id uint, name string, extra []byte) error
	Get(ctx context.Context, id uint, name string) ([]byte, error)
}

func (s appService) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserInfoReply, error) {
	var (
		u           *entity.User
		device      *entity.Device
		wechatExtra *pb.WechatExtra
		err         error
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
		u, err = s.handleMobileLogin(ctx, in.PackageName, in.Mobile, in.Code, device)
	} else if len(in.Wechat) != 0 {
		u, wechatExtra, err = s.handleWechatLogin(ctx, in.PackageName, in.Wechat, device)
		ctx = context.WithValue(ctx, wechatExtraKey, wechatExtra)
	} else {
		u, err = s.handleDeviceLogin(ctx, in.PackageName, device.Suuid, device)
	}
	if err != nil {
		return nil, err
	}

	err = s.addUserSourceInfo(ctx, in, u)
	if err != nil {
		return nil, err
	}

	// Create jwt token
	tokenString, err := s.getToken(&tokenParam{uint64(u.ID), in.Device.Suuid, in.Channel, in.VersionCode, u.WechatOpenId.String, in.Mobile, in.PackageName})
	if err != nil {
		return nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorJwtFailure))
	}

	// 拼装返回结果
	var resp = u.ToReply()

	resp.Data.Token = tokenString

	s.decorate(ctx, resp.Data)
	return resp, nil
}

func (s appService) addUserSourceInfo(ctx context.Context, in *pb.UserLoginRequest, u *entity.User) error {
	var (
		err      error
		hasExtra bool
	)
	if in.ThirdPartyId != "" && in.ThirdPartyId != u.ThirdPartyId {
		u.ThirdPartyId = in.ThirdPartyId
		hasExtra = true
	}
	// 任何情况下Channel不得更新
	if in.Channel != "" && u.Channel == "" {
		u.Channel = in.Channel
		hasExtra = true
	}
	if in.VersionCode != "" && in.VersionCode != u.VersionCode {
		u.VersionCode = in.VersionCode
		hasExtra = true
	}

	if hasExtra {
		err = s.ur.Save(ctx, u)
		if err != nil {
			return kerr.InternalErr(errors.Wrap(err, msg.ErrorDatabaseFailure))
		}
	}
	return nil
}

type tokenParam struct {
	userId                                                   uint64
	suuid, channel, versionCode, wechat, mobile, packageName string
}

func (s appService) getToken(param *tokenParam) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		kittyjwt.NewClaim(
			param.userId,
			s.conf.String("name"),
			param.suuid, param.channel, param.versionCode, param.wechat, param.mobile, param.packageName,
			time.Hour*24*30,
		),
	)
	token.Header["kid"] = s.conf.String("security.kid")
	return token.SignedString([]byte(s.conf.String("security.key")))
}

func (s appService) verify(ctx context.Context, mobile string, code string) (bool, error) {
	result, err := s.cr.CheckCode(ctx, mobile, code)
	if err != nil {
		return false, err
	}
	err = s.cr.DeleteCode(ctx, mobile)
	s.warn(err)
	return result, nil
}

func (s appService) GetCode(ctx context.Context, in *pb.GetCodeRequest) (*pb.GenericReply, error) {
	code, err := s.cr.AddCode(ctx, in.Mobile)
	if err == repository.ErrTooFrequent {
		return nil, kerr.ResourceExhaustedErr(err)
	}
	if err != nil {
		return nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorGetCode))
	}
	err = s.sender.Send(ctx, in.Mobile, code)
	if err != nil {
		return nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorSendCode))
	}
	var resp = pb.GenericReply{
		Code: 0,
	}
	return &resp, nil
}

func (s appService) debug(err error) {
	if err != nil {
		level.Debug(s.log).Log("err", err)
	}
}

func (s appService) infof(msg string, args ...interface{}) {
	level.Info(s.log).Log("msg", fmt.Sprintf(msg, args))
}

func (s appService) error(err error) {
	if err != nil {
		level.Error(s.log).Log("err", err)
	}
}
func (s appService) warn(err error) {
	if err != nil {
		level.Warn(s.log).Log("err", err)
	}
}

func (s appService) getWechatInfo(ctx context.Context, wechat string) (*pb.WechatExtra, error) {
	wxRes, err := s.wechat.GetWechatLoginResponse(ctx, wechat)
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorWechatFailure)
	}
	if wxRes.Openid == "" {
		return nil, errors.New(msg.ErrorMissingOpenid)
	}
	wxInfo, err := s.wechat.GetWechatUserInfoResult(ctx, wxRes)
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorWechatFailure)
	}
	// side effect: save extra wechat info
	infoPb := &pb.WechatExtra{
		AccessToken:  wxRes.AccessToken,
		ExpiresIn:    wxRes.ExpiresIn,
		RefreshToken: wxRes.RefreshToken,
		OpenId:       wxRes.Openid,
		Scope:        wxRes.Scope,
		NickName:     wxInfo.Nickname,
		Sex:          int32(wxInfo.Sex),
		Province:     wxInfo.Province,
		City:         wxInfo.City,
		Country:      wxInfo.Country,
		Headimgurl:   wxInfo.Headimgurl,
		Privilege:    wxInfo.Privilege,
		Unionid:      wxInfo.Unionid,
	}
	b, err := infoPb.Marshal()
	if err != nil {
		s.warn(err)
	}
	userId := kittyjwt.GetClaim(ctx).UserId
	err = s.er.Put(ctx, uint(userId), pb.Extra_WECHAT_EXTRA.String(), b)
	if err != nil {
		s.warn(err)
	}
	return infoPb, nil
}

func (s appService) handleWechatLogin(ctx context.Context, packageName, wechat string, device *entity.Device) (*entity.User, *pb.WechatExtra, error) {
	wxInfo, err := s.getWechatInfo(ctx, wechat)
	if err != nil {
		return nil, nil, kerr.UnauthorizedErr(err)
	}

	headImg, err := s.fr.UploadFromUrl(ctx, wxInfo.Headimgurl)
	s.warn(err)

	wechatUser := entity.User{
		UserName:      wxInfo.NickName,
		HeadImg:       headImg,
		WechatOpenId:  ns(wxInfo.OpenId),
		WechatUnionId: ns(wxInfo.Unionid),
	}

	u, err := s.ur.GetFromWechat(ctx, packageName, wxInfo.OpenId, device, wechatUser)
	if err != nil {
		return nil, nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorDatabaseFailure))
	}
	s.infof(msg.WxSuccess, u.ID)
	return u, wxInfo, nil
}

func (s appService) handleMobileLogin(ctx context.Context, packageName, mobile, code string, device *entity.Device) (*entity.User, error) {
	if len(code) == 0 {
		return nil, kerr.InvalidArgumentErr(errors.New(msg.InvalidParams))
	}
	if ok, err := s.verify(ctx, mobile, code); err != nil {
		return nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorDatabaseFailure))
	} else if !ok {
		return nil, kerr.UnauthorizedErr(errors.New(msg.ErrorMobileCode))
	}
	u, err := s.ur.GetFromMobile(ctx, packageName, mobile, device)
	if err != nil {
		return nil, dbErr(err)
	}
	s.infof(msg.MobileSuccess, u.ID)
	return u, nil
}

func (s appService) handleDeviceLogin(ctx context.Context, packageName, suuid string, device *entity.Device) (*entity.User, error) {
	u, err := s.ur.GetFromDevice(ctx, packageName, suuid, device)
	if err != nil {
		return nil, dbErr(err)
	}
	s.infof(msg.DeviceSuccess, u.ID)
	return u, nil
}

func dbErr(err error) kerr.ServerError {
	return kerr.InternalErr(errors.Wrap(err, msg.ErrorDatabaseFailure))
}

func (s appService) GetInfo(ctx context.Context, in *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	if in.Id == 0 {
		claim := kittyjwt.GetClaim(ctx)
		in.Id = claim.UserId
	}
	u, err := s.ur.Get(ctx, uint(in.Id))
	if err != nil {
		return nil, kerr.NotFoundErr(errors.Wrap(err, msg.ErrorUserNotFound))
	}
	var resp = u.ToReply()

	if in.Taobao {
		resp.Data.TaobaoExtra = s.getTaobaoExtra(ctx, uint(in.Id))
	}

	if in.Wechat {
		resp.Data.WechatExtra = s.getWechatExtra(ctx, uint(in.Id))
	}

	return resp, nil
}

func (s appService) UpdateInfo(ctx context.Context, in *pb.UserInfoUpdateRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.GetClaim(ctx)
	u, err := s.ur.Update(ctx, uint(claim.UserId), entity.User{
		UserName:     in.UserName,
		HeadImg:      in.HeadImg,
		Gender:       int(in.Gender),
		Birthday:     in.Birthday,
		ThirdPartyId: in.ThirdPartyId,
	})
	if err != nil {
		return nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorDatabaseFailure))
	}

	var resp = u.ToReply()
	s.decorate(ctx, resp.Data)
	return resp, nil
}

func (s appService) Bind(ctx context.Context, in *pb.UserBindRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.GetClaim(ctx)

	var (
		toUpdate    entity.User
		wechatExtra *pb.WechatExtra
		taobaoExtra *pb.TaobaoExtra
		err         error
	)

	// 绑定手机号
	if len(in.Mobile) > 0 && len(in.Code) > 0 {
		if ok, err := s.verify(ctx, in.Mobile, in.Code); err != nil {
			return nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorDatabaseFailure))
		} else if !ok {
			return nil, kerr.UnauthorizedErr(errors.New(msg.ErrorMobileCode))
		}
		toUpdate = entity.User{Mobile: ns(in.Mobile)}
	}

	// 绑定微信号
	if len(in.Wechat) > 0 {
		wechatExtra, err = s.getWechatInfo(ctx, in.Wechat)
		ctx = context.WithValue(ctx, wechatExtraKey, wechatExtra)
		if err != nil {
			return nil, kerr.UnauthorizedErr(err)
		}
		toUpdate = entity.User{
			WechatOpenId:  ns(wechatExtra.OpenId),
			WechatUnionId: ns(wechatExtra.Unionid),
		}
	}

	// 绑定淘宝openId
	if in.TaobaoExtra != nil && len(in.TaobaoExtra.OpenId) > 0 {
		taobaoExtra = in.TaobaoExtra
		ctx = context.WithValue(ctx, taobaoExtraKey, taobaoExtra)
		toUpdate = entity.User{
			TaobaoOpenId: ns(in.TaobaoExtra.OpenId),
		}
		extra, err := in.TaobaoExtra.Marshal()
		if err != nil {
			s.warn(err)
		}
		err = s.er.Put(ctx, uint(claim.UserId), pb.Extra_TAOBAO_EXTRA.String(), extra)
		if err != nil {
			s.warn(err)
		}
	}

	// 绑定微信openId
	if len(in.OpenId) > 0 {
		toUpdate = entity.User{
			WechatOpenId: ns(in.OpenId),
		}
	}

	// 更新用户
	newUser, err := s.ur.Update(ctx, uint(claim.UserId), toUpdate)
	if err != nil {
		return nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorDatabaseFailure))
	}
	reply := newUser.ToReply()
	reply.Data.Token, err = s.getToken(&tokenParam{
		uint64(newUser.ID),
		newUser.CommonSUUID,
		newUser.Channel,
		newUser.VersionCode,
		newUser.WechatOpenId.String,
		newUser.Mobile.String,
		newUser.PackageName,
	})
	if err != nil {
		err = kerr.InternalErr(errors.Wrap(err, msg.ErrorJwtFailure))
	}

	s.decorate(ctx, reply.Data)
	return reply, err
}

func (s appService) Unbind(ctx context.Context, in *pb.UserUnbindRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.GetClaim(ctx)
	user, err := s.ur.Get(ctx, uint(claim.UserId))
	if err != nil {
		return nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorDatabaseFailure))
	}
	if in.Mobile {
		user.Mobile = sql.NullString{}
	}
	if in.Wechat {
		user.WechatUnionId = sql.NullString{}
		user.WechatOpenId = sql.NullString{}
	}
	if in.Taobao {
		user.TaobaoOpenId = sql.NullString{}
	}
	err = s.ur.Save(ctx, user)
	if err != nil {
		return nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorDatabaseFailure))
	}
	var resp = user.ToReply()
	s.decorate(ctx, resp.Data)
	return resp, nil
}

func ns(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func (s appService) Refresh(ctx context.Context, in *pb.UserRefreshRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.GetClaim(ctx)
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
		return nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorDatabaseFailure))
	}

	u.CommonSUUID = in.Device.Suuid
	u.Channel = in.Channel
	u.VersionCode = in.VersionCode
	u.AddNewDevice(device)

	if err := s.ur.Save(ctx, u); err != nil {
		return nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorDatabaseFailure))
	}

	reply := u.ToReply()
	reply.Data.Token, err = s.getToken(&tokenParam{
		uint64(u.ID),
		u.CommonSUUID,
		u.Channel,
		u.VersionCode,
		u.WechatOpenId.String,
		u.Mobile.String,
		u.PackageName,
	})
	s.decorate(ctx, reply.Data)
	if err != nil {
		err = kerr.InternalErr(errors.Wrap(err, msg.ErrorJwtFailure))
	}
	return reply, nil
}

func (s appService) getWechatExtra(ctx context.Context, id uint) *pb.WechatExtra {
	var extra pb.WechatExtra

	if extra, ok := ctx.Value(wechatExtraKey).(*pb.WechatExtra); ok {
		return extra
	}

	b, err := s.er.Get(ctx, uint(id), pb.Extra_WECHAT_EXTRA.String())
	if err != nil {
		s.warn(err)
	}
	extra.Unmarshal(b)
	return &extra
}

func (s appService) getTaobaoExtra(ctx context.Context, id uint) *pb.TaobaoExtra {
	var extra pb.TaobaoExtra

	if extra, ok := ctx.Value(taobaoExtraKey).(*pb.TaobaoExtra); ok {
		return extra
	}

	b, err := s.er.Get(ctx, uint(id), pb.Extra_TAOBAO_EXTRA.String())
	if err != nil {
		s.warn(err)
	}
	extra.Unmarshal(b)
	return &extra
}

func (s appService) decorate(ctx context.Context, data *pb.UserInfo) {
	data.TaobaoExtra = s.getTaobaoExtra(ctx, uint(data.Id))
	data.WechatExtra = s.getWechatExtra(ctx, uint(data.Id))
}
