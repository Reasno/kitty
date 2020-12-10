//go:generate mockery --name=CodeRepository
//go:generate mockery --name=UserRepository
//go:generate mockery --name=FileRepository

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

type appService struct {
	conf   contract.ConfigReader
	logger log.Logger
	ur     UserRepository
	cr     CodeRepository
	fr     FileRepository
	sender contract.SmsSender
	wechat wechat.Wechater
}

type tokenParam struct {
	userId                                                                 uint64
	suuid, channel, versionCode, wechat, mobile, packageName, thirdPartyId string
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
	GetAll(ctx context.Context, ids ...uint) (user []entity.User, err error)
	Save(ctx context.Context, user *entity.User) error
}

type FileRepository interface {
	UploadFromUrl(ctx context.Context, url string) (newUrl string, err error)
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
		return nil, err
	}

	// 再存一些信息
	if err := s.addChannelAndVersionInfo(ctx, in, u); err != nil {
		s.warn(err)
	}

	// Create jwt token
	tokenString, err := s.getToken(&tokenParam{uint64(u.ID), u.CommonSUUID, u.Channel, u.VersionCode, u.WechatOpenId.String, u.Mobile.String, u.PackageName, u.ThirdPartyId})
	if err != nil {
		s.warn(err)
	}

	// 拼装返回结果
	var resp = toReply(u)
	resp.Data.Token = tokenString

	return resp, nil
}

func (s appService) GetCode(ctx context.Context, in *pb.GetCodeRequest) (*pb.GenericReply, error) {
	code, err := s.cr.AddCode(ctx, in.Mobile)
	if err == repository.ErrTooFrequent {
		return nil, kerr.ResourceExhaustedErr(err, msg.ErrorTooFrequent)
	}
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorGetCode)
	}
	err = s.sender.Send(ctx, in.Mobile, code)
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorSendCode)
	}
	var resp = pb.GenericReply{
		Code: 0,
	}
	return &resp, nil
}

func (s appService) GetInfo(ctx context.Context, in *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	if in.Id == 0 {
		claim := kittyjwt.ClaimFromContext(ctx)
		in.Id = claim.UserId
	}
	u, err := s.ur.Get(ctx, uint(in.Id))
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, kerr.NotFoundErr(err, msg.ErrorRecordNotFound)
	}
	if err != nil {
		return nil, dbErr(err)
	}
	var resp = toReply(u)

	if !in.Taobao {
		resp.Data.TaobaoExtra = nil
	}

	if !in.Wechat {
		resp.Data.WechatExtra = nil
	}

	return resp, nil
}

func (s appService) Refresh(ctx context.Context, in *pb.UserRefreshRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.ClaimFromContext(ctx)
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
		u.ThirdPartyId,
	})

	if err != nil {
		err = kerr.InternalErr(err, msg.ErrorJwtFailure)
	}
	return reply, nil
}

func (s appService) UpdateInfo(ctx context.Context, in *pb.UserInfoUpdateRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.ClaimFromContext(ctx)
	u, err := s.ur.Update(ctx, uint(claim.UserId), entity.User{
		UserName:     in.UserName,
		HeadImg:      in.HeadImg,
		Gender:       int(in.Gender),
		Birthday:     in.Birthday,
		ThirdPartyId: in.ThirdPartyId,
	})
	if err != nil {
		return nil, dbErr(err)
	}

	var resp = toReply(u)
	return resp, nil

}

func (s appService) Bind(ctx context.Context, in *pb.UserBindRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.ClaimFromContext(ctx)

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
			return nil, kerr.UnauthenticatedErr(errors.Errorf("cannot verify %s with %s", in.Mobile, in.Code), msg.ErrorMobileCode)
		}
		toUpdate.Mobile = ns(in.Mobile)
	}

	// 绑定微信号
	if len(in.Wechat) > 0 {
		var wechatExtra *pb.WechatExtra
		wechatExtra, err = s.getWechatInfo(ctx, in.Wechat)
		if err != nil {
			return nil, kerr.UnauthenticatedErr(err, msg.ErrorWechatFailure)
		}
		wechatExtraBytes, err := wechatExtra.Marshal()
		if err != nil {
			return nil, kerr.InternalErr(err, msg.ErrorLogin)
		}
		toUpdate.WechatOpenId = ns(wechatExtra.OpenId)
		toUpdate.WechatUnionId = ns(wechatExtra.Unionid)
		toUpdate.WechatExtra = wechatExtraBytes
		if in.MergeInfo {
			toUpdate.UserName = wechatExtra.NickName
			toUpdate.HeadImg, _ = s.fr.UploadFromUrl(ctx, wechatExtra.Headimgurl)
			toUpdate.Gender = int(wechatExtra.Sex)
		}
	}

	// 绑定淘宝openId
	if in.TaobaoExtra != nil && len(in.TaobaoExtra.OpenId) > 0 {
		taobaoExtraBytes, err := in.TaobaoExtra.Marshal()
		if err != nil {
			return nil, kerr.InternalErr(err, msg.ErrorCorruptedData)
		}
		toUpdate.TaobaoOpenId = ns(in.TaobaoExtra.OpenId)
		toUpdate.TaobaoExtra = taobaoExtraBytes
		if in.MergeInfo {
			toUpdate.UserName = in.TaobaoExtra.Nick
			toUpdate.HeadImg, _ = s.fr.UploadFromUrl(ctx, in.TaobaoExtra.AvatarUrl)
		}
	}

	// 绑定微信openId
	if in.WechatExtra != nil && len(in.WechatExtra.OpenId) > 0 {
		wechatExtraBytes, err := in.WechatExtra.Marshal()
		if err != nil {
			return nil, kerr.InternalErr(err, msg.ErrorCorruptedData)
		}
		toUpdate.WechatOpenId = ns(in.WechatExtra.OpenId)
		toUpdate.WechatExtra = wechatExtraBytes
	}

	// 更新用户
	newUser, err := s.ur.Update(ctx, uint(claim.UserId), toUpdate)
	if errors.Is(err, repository.ErrAlreadyBind) {
		return nil, kerr.FailedPreconditionErr(err, msg.ErrorAlreadyBind)
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
		newUser.ThirdPartyId,
	})
	if err != nil {
		err = kerr.InternalErr(err, msg.ErrorJwtFailure)
	}

	return reply, err
}

func (s appService) Unbind(ctx context.Context, in *pb.UserUnbindRequest) (*pb.UserInfoReply, error) {
	claim := kittyjwt.ClaimFromContext(ctx)
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
		user.WechatExtra = nil
	}
	if in.Taobao {
		user.TaobaoOpenId = sql.NullString{}
		user.TaobaoExtra = nil
	}

	err = s.ur.Save(ctx, user)
	if err != nil {
		return nil, dbErr(err)
	}
	var resp = toReply(user)
	return resp, nil
}

func (s appService) getToken(param *tokenParam) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		kittyjwt.NewClaim(
			param.userId,
			s.conf.String("name"),
			param.suuid, param.channel, param.versionCode, param.wechat, param.mobile, param.packageName, param.thirdPartyId,
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

func (s appService) handleWechatLogin(ctx context.Context, packageName, wechat string, device *entity.Device) (*entity.User, error) {
	wxInfo, err := s.getWechatInfo(ctx, wechat)
	if err != nil {
		return nil, kerr.UnauthenticatedErr(err, msg.ErrorWechatFailure)
	}
	extra, err := wxInfo.Marshal()
	if err != nil {
		return nil, kerr.InternalErr(err, msg.ErrorCorruptedData)
	}

	wechatUser := entity.User{
		UserName:      wxInfo.NickName,
		HeadImg:       wxInfo.Headimgurl,
		WechatOpenId:  ns(wxInfo.OpenId),
		WechatUnionId: ns(wxInfo.Unionid),
		WechatExtra:   extra,
	}

	u, err := s.ur.GetFromWechat(ctx, packageName, wxInfo.OpenId, device, wechatUser)
	if err != nil {
		return nil, dbErr(err)
	}
	level.Info(s.logger).Log("msg", fmt.Sprintf(msg.WxSuccess, u.ID), "suuid", device.Suuid, "userId", u.ID, "packageName", packageName)
	return u, nil
}

func (s appService) handleMobileLogin(ctx context.Context, packageName, mobile, code string, device *entity.Device) (*entity.User, error) {
	if len(code) == 0 {
		return nil, kerr.InvalidArgumentErr(errors.New("code cannot be 0"), msg.InvalidParams)
	}
	if ok, err := s.verify(ctx, mobile, code); err != nil {
		return nil, err
	} else if !ok {
		return nil, kerr.UnauthenticatedErr(errors.Errorf("cannot verify %s with %s", mobile, code), msg.ErrorMobileCode)
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
		return s.handleWechatLogin(ctx, in.PackageName, in.Wechat, device)
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

func dbErr(err error) kerr.ServerError {
	return kerr.InternalErr(err, msg.ErrorDatabaseFailure)
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
	var wechatExtra pb.WechatExtra
	_ = wechatExtra.Unmarshal(user.WechatExtra)
	var taobaoExtra pb.TaobaoExtra
	_ = taobaoExtra.Unmarshal(user.TaobaoExtra)
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
			Mobile:       redact(user.Mobile.String),
			IsNew:        user.IsNew,
			WechatExtra:  &wechatExtra,
			TaobaoExtra:  &taobaoExtra,
		},
	}
}

func (s appService) GetInfoBatch(ctx context.Context, in *pb.UserInfoBatchRequest) (*pb.UserInfoBatchReply, error) {
	var args []uint
	for _, v := range in.Id {
		args = append(args, uint(v))
	}
	users, err := s.ur.GetAll(ctx, args...)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, kerr.NotFoundErr(err, msg.ErrorRecordNotFound)
	}
	if err != nil {
		return nil, dbErr(err)
	}
	var resp = pb.UserInfoBatchReply{
		Code: 0,
		Data: []*pb.UserInfo{},
	}

	for _, v := range users {
		tmp := toReply(&v).Data
		if !in.Taobao {
			tmp.TaobaoExtra = nil
		}
		if !in.Wechat {
			tmp.WechatExtra = nil
		}
		resp.Data = append(resp.Data, tmp)
	}

	return &resp, nil
}
