package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"glab.tagtic.cn/ad_gains/kitty/app/msg"
	"glab.tagtic.cn/ad_gains/kitty/app/repository"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	kittyjwt "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
	"glab.tagtic.cn/ad_gains/kitty/pkg/wechat"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

var taobaoExtraKey struct{}
var wechatExtraKey struct{}

type appService struct {
	conf     contract.ConfigReader
	logger   log.Logger
	ur       UserRepository
	cr       CodeRepository
	er       ExtraRepository
	sender   contract.SmsSender
	wechat   wechat.Wechater
	uploader contract.Uploader
}

type tokenParam struct {
	userId                                                   uint64
	suuid, channel, versionCode, wechat, mobile, packageName string
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

type ExtraRepository interface {
	Put(ctx context.Context, id uint, name string, extra []byte) error
	Get(ctx context.Context, id uint, name string) ([]byte, error)
	Del(ctx context.Context, id uint, name string) error
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
	u, err = s.loginFrom(ctx, in, device)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = s.addChannelAndVersionInfo(ctx, in, u)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Create jwt token
	tokenString, err := s.getToken(&tokenParam{uint64(u.ID), u.CommonSUUID, u.Channel, u.VersionCode, u.WechatOpenId.String, u.Mobile.String, u.PackageName})
	if err != nil {
		return nil, kerr.InternalErr(errors.Wrap(err, msg.ErrorJwtFailure))
	}

	// 拼装返回结果
	var resp = toReply(u)
	resp.Data.Token = tokenString

	s.persistExtra(ctx)
	s.decorateResponse(ctx, resp.Data)

	return resp, nil
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

func (s appService) GetInfo(ctx context.Context, in *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	if in.Id == 0 {
		claim := kittyjwt.GetClaim(ctx)
		in.Id = claim.UserId
	}
	u, err := s.ur.Get(ctx, uint(in.Id))
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, kerr.NotFoundErr(errors.Wrap(err, msg.ErrorRecordNotFound))
	}
	if err != nil {
		return nil, dbErr(err)
	}
	var resp = toReply(u)

	if in.Taobao {
		resp.Data.TaobaoExtra = s.getTaobaoExtra(ctx, uint(in.Id))
	}

	if in.Wechat {
		resp.Data.WechatExtra = s.getWechatExtra(ctx, uint(in.Id))
	}

	return resp, nil
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
		return nil, dbErr(err)
	}

	u.CommonSUUID = in.Device.Suuid
	u.Channel = in.Channel
	u.VersionCode = in.VersionCode
	u.AddNewDevice(device)

	if err := s.ur.Save(ctx, u); err != nil {
		return nil, dbErr(err)
	}

	reply := toReply(u)
	reply.Data.Token, err = s.getToken(&tokenParam{
		uint64(u.ID),
		u.CommonSUUID,
		u.Channel,
		u.VersionCode,
		u.WechatOpenId.String,
		u.Mobile.String,
		u.PackageName,
	})

	if err != nil {
		err = kerr.InternalErr(errors.Wrap(err, msg.ErrorJwtFailure))
	}
	s.decorateResponse(ctx, reply.Data)
	return reply, nil
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

	var resp = toReply(u)
	s.decorateResponse(ctx, resp.Data)
	return resp, nil
}

func (s appService) Bind(ctx context.Context, in *pb.UserBindRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.GetClaim(ctx)

	var (
		toUpdate entity.User
		err      error
	)

	// 兼容旧接口
	if len(in.OpenId) > 0 {
		in.WechatExtra = &pb.WechatExtra{OpenId: in.OpenId}
	}

	// 绑定手机号
	if len(in.Mobile) > 0 && len(in.Code) > 0 {
		if ok, err := s.verify(ctx, in.Mobile, in.Code); err != nil {
			return nil, dbErr(err)
		} else if !ok {
			return nil, kerr.UnauthorizedErr(errors.New(msg.ErrorMobileCode))
		}
		toUpdate = entity.User{Mobile: ns(in.Mobile)}
	}

	// 绑定微信号
	if len(in.Wechat) > 0 {
		var wechatExtra *pb.WechatExtra
		wechatExtra, err = s.getWechatInfo(ctx, in.Wechat)
		if err != nil {
			return nil, kerr.UnauthorizedErr(err)
		}
		ctx = context.WithValue(ctx, wechatExtraKey, wechatExtra)
		toUpdate = entity.User{
			WechatOpenId:  ns(wechatExtra.OpenId),
			WechatUnionId: ns(wechatExtra.Unionid),
		}
	}

	// 绑定淘宝openId
	if in.TaobaoExtra != nil && len(in.TaobaoExtra.OpenId) > 0 {
		ctx = context.WithValue(ctx, taobaoExtraKey, in.TaobaoExtra)
		toUpdate = entity.User{
			TaobaoOpenId: ns(in.TaobaoExtra.OpenId),
		}
	}

	// 绑定微信openId
	if in.WechatExtra != nil && len(in.WechatExtra.OpenId) > 0 {
		ctx = context.WithValue(ctx, wechatExtraKey, in.WechatExtra)
		toUpdate = entity.User{
			WechatOpenId: ns(in.WechatExtra.OpenId),
		}
	}

	// 更新用户
	newUser, err := s.ur.Update(ctx, uint(claim.UserId), toUpdate)
	if errors.Is(err, repository.ErrAlreadyBind) {
		return nil, kerr.FailedPreconditionErr(errors.Wrap(err, msg.ErrorAlreadyBind))
	}
	if err != nil {
		return nil, dbErr(err)
	}

	// 获取Token
	reply := toReply(newUser)
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

	// 组装数据
	s.persistExtra(ctx)
	s.decorateResponse(ctx, reply.Data)
	return reply, err
}

func (s appService) Unbind(ctx context.Context, in *pb.UserUnbindRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.GetClaim(ctx)
	user, err := s.ur.Get(ctx, uint(claim.UserId))
	if err != nil {
		return nil, dbErr(err)
	}
	if in.Mobile {
		user.Mobile = sql.NullString{}
	}
	if in.Wechat {
		user.WechatUnionId = sql.NullString{}
		user.WechatOpenId = sql.NullString{}
		s.warn(s.er.Del(ctx, user.ID, pb.Extra_WECHAT_EXTRA.String()))
	}
	if in.Taobao {
		user.TaobaoOpenId = sql.NullString{}
		s.warn(s.er.Del(ctx, user.ID, pb.Extra_TAOBAO_EXTRA.String()))
	}
	err = s.ur.Save(ctx, user)
	if err != nil {
		return nil, dbErr(err)
	}

	var resp = toReply(user)
	s.decorateResponse(ctx, resp.Data)
	return resp, nil
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

func (s appService) debug(err error) {
	if err != nil {
		level.Debug(s.logger).Log("err", err)
	}
}

func (s appService) error(err error) {
	if err != nil {
		level.Error(s.logger).Log("err", err)
	}
}

func (s appService) warn(err error) {
	if err != nil {
		level.Warn(s.logger).Log("err", err)
	}
}
func (s appService) getWechatInfo(ctx context.Context, wechat string) (*pb.WechatExtra, error) {
	wxRes, err := s.wechat.GetLoginResponse(ctx, wechat)
	if err != nil {
		return nil, errors.Wrap(err, msg.ErrorWechatFailure)
	}
	if wxRes.Openid == "" {
		return nil, errors.New(msg.ErrorMissingOpenid)
	}
	wxInfo, err := s.wechat.GetUserInfoResult(ctx, wxRes)
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
	return infoPb, nil
}

func (s appService) handleWechatLogin(ctx context.Context, packageName, wechat string, device *entity.Device) (*entity.User, *pb.WechatExtra, error) {
	wxInfo, err := s.getWechatInfo(ctx, wechat)
	if err != nil {
		return nil, nil, kerr.UnauthorizedErr(err)
	}

	wechatUser := entity.User{
		UserName:      wxInfo.NickName,
		HeadImg:       wxInfo.Headimgurl,
		WechatOpenId:  ns(wxInfo.OpenId),
		WechatUnionId: ns(wxInfo.Unionid),
	}

	u, err := s.ur.GetFromWechat(ctx, packageName, wxInfo.OpenId, device, wechatUser)
	if err != nil {
		return nil, nil, dbErr(err)
	}
	level.Info(s.logger).Log("msg", fmt.Sprintf(msg.WxSuccess, u.ID), "suuid", device.Suuid, "userId", u.ID, "packageName", packageName)
	return u, wxInfo, nil
}

func (s appService) handleMobileLogin(ctx context.Context, packageName, mobile, code string, device *entity.Device) (*entity.User, error) {
	if len(code) == 0 {
		return nil, kerr.InvalidArgumentErr(errors.New(msg.InvalidParams))
	}
	if ok, err := s.verify(ctx, mobile, code); err != nil {
		return nil, err
	} else if !ok {
		return nil, kerr.UnauthorizedErr(errors.New(msg.ErrorMobileCode))
	}
	u, err := s.ur.GetFromMobile(ctx, packageName, mobile, device)
	if err != nil {
		return nil, dbErr(err)
	}
	level.Info(s.logger).Log("msg", fmt.Sprintf(msg.MobileSuccess, u.ID), "suuid", device.Suuid, "userId", u.ID, "packageName", packageName)
	return u, nil
}

func (s appService) handleDeviceLogin(ctx context.Context, packageName, suuid string, device *entity.Device) (*entity.User, error) {
	u, err := s.ur.GetFromDevice(ctx, packageName, suuid, device)
	if err != nil {
		return nil, dbErr(err)
	}
	level.Info(s.logger).Log("msg", fmt.Sprintf(msg.DeviceSuccess, u.ID), "suuid", device.Suuid, "userId", u.ID, "packageName", packageName)
	return u, nil
}

func (s appService) loginFrom(ctx context.Context, in *pb.UserLoginRequest, device *entity.Device) (*entity.User, error) {

	if len(in.Mobile) != 0 {
		return s.handleMobileLogin(ctx, in.PackageName, in.Mobile, in.Code, device)
	}

	if len(in.Wechat) != 0 {
		u, wechatExtra, err := s.handleWechatLogin(ctx, in.PackageName, in.Wechat, device)
		ctx = context.WithValue(ctx, wechatExtraKey, wechatExtra)
		return u, err
	}

	return s.handleDeviceLogin(ctx, in.PackageName, device.Suuid, device)
}

func (s appService) addChannelAndVersionInfo(ctx context.Context, in *pb.UserLoginRequest, u *entity.User) error {
	var (
		err      error
		hasExtra bool
	)
	if in.ThirdPartyId != "" && in.ThirdPartyId != u.ThirdPartyId {
		u.ThirdPartyId = in.ThirdPartyId
		hasExtra = true
	}

	if in.Channel != "" && u.Channel != in.Channel {
		u.Channel = in.Channel
		hasExtra = true
	}
	if in.VersionCode != "" && in.VersionCode != u.VersionCode {
		u.VersionCode = in.VersionCode
		hasExtra = true
	}
	if in.InviteCode != "" && u.InviteCode == "" {
		u.InviteCode = in.InviteCode
		hasExtra = true
	}

	if hasExtra {
		err = s.ur.Save(ctx, u)
		if err != nil {
			return dbErr(err)
		}
	}
	return nil
}

func (s appService) verify(ctx context.Context, mobile string, code string) (bool, error) {
	result, err := s.cr.CheckCode(ctx, mobile, code)
	if err != nil {
		return false, dbErr(err)
	}
	err = s.cr.DeleteCode(ctx, mobile)
	s.warn(err)
	return result, nil
}

func (s appService) getWechatExtra(ctx context.Context, id uint) *pb.WechatExtra {
	var extra pb.WechatExtra

	if extra, ok := ctx.Value(wechatExtraKey).(*pb.WechatExtra); ok {
		return extra
	}

	b, err := s.er.Get(ctx, id, pb.Extra_WECHAT_EXTRA.String())
	s.warn(err)
	err = extra.Unmarshal(b)
	s.warn(err)

	return &extra
}

func (s appService) getTaobaoExtra(ctx context.Context, id uint) *pb.TaobaoExtra {
	var extra pb.TaobaoExtra

	if extra, ok := ctx.Value(taobaoExtraKey).(*pb.TaobaoExtra); ok {
		return extra
	}

	b, err := s.er.Get(ctx, id, pb.Extra_TAOBAO_EXTRA.String())
	s.warn(err)
	err = extra.Unmarshal(b)
	s.warn(err)

	return &extra
}

func (s appService) decorateResponse(ctx context.Context, data *pb.UserInfo) {
	data.TaobaoExtra = s.getTaobaoExtra(ctx, uint(data.Id))
	data.WechatExtra = s.getWechatExtra(ctx, uint(data.Id))
	// 如果不是用户本人，则隐去手机号部分内容
	if data.Id != kittyjwt.GetClaim(ctx).UserId {
		data.Mobile = redact(data.Mobile)
	}
}

func (s appService) persistExtra(ctx context.Context) {
	s.persistTaobaoExtra(ctx)
	s.persistWechatExtra(ctx)
}

func (s appService) persistTaobaoExtra(ctx context.Context) {
	claim := kittyjwt.GetClaim(ctx)
	extra, ok := ctx.Value(taobaoExtraKey).(*pb.TaobaoExtra)
	if !ok {
		return
	}
	b, err := extra.Marshal()
	s.warn(err)

	err = s.er.Put(ctx, uint(claim.UserId), pb.Extra_TAOBAO_EXTRA.String(), b)
	s.warn(err)
}

func (s appService) persistWechatExtra(ctx context.Context) {
	claim := kittyjwt.GetClaim(ctx)
	extra, ok := ctx.Value(wechatExtraKey).(*pb.WechatExtra)
	if !ok {
		return
	}
	b, err := extra.Marshal()
	s.warn(err)

	err = s.er.Put(ctx, uint(claim.UserId), pb.Extra_WECHAT_EXTRA.String(), b)
	s.warn(err)
}

func dbErr(err error) kerr.ServerError {
	return kerr.InternalErr(errors.Wrap(err, msg.ErrorDatabaseFailure))
}

func ns(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func redact(mobile string) string {
	if len(mobile) >= 11 {
		mobile = mobile[:3] + "****" + mobile[7:]
	}
	return mobile
}

func toReply(user *entity.User) *pb.UserInfoReply {
	return &pb.UserInfoReply{
		Code:    0,
		Message: "",
		Data: &pb.UserInfo{
			Id:           uint64(user.ID),
			UserName:     user.UserName,
			Wechat:       user.WechatOpenId.String,
			HeadImg:      user.HeadImg,
			Gender:       pb.Gender(user.Gender),
			Birthday:     user.Birthday,
			ThirdPartyId: user.ThirdPartyId,
			Mobile:       user.Mobile.String,
			IsNew:        user.IsNew,
		},
	}
}
