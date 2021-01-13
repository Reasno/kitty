// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: app.proto

package kitty

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = ptypes.DynamicAny{}
)

// define the regex for a UUID once up-front
var _app_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on UserBindRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *UserBindRequest) Validate() error {
	if m == nil {
		return nil
	}

	if !_UserBindRequest_Mobile_Pattern.MatchString(m.GetMobile()) {
		return UserBindRequestValidationError{
			field:  "Mobile",
			reason: "value does not match regex pattern \"(^$|^[\\\\d]{11}$)\"",
		}
	}

	// no validation rules for Code

	// no validation rules for Wechat

	// no validation rules for OpenId

	if v, ok := interface{}(m.GetTaobaoExtra()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserBindRequestValidationError{
				field:  "TaobaoExtra",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetWechatExtra()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserBindRequestValidationError{
				field:  "WechatExtra",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for MergeInfo

	return nil
}

// UserBindRequestValidationError is the validation error returned by
// UserBindRequest.Validate if the designated constraints aren't met.
type UserBindRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserBindRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserBindRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserBindRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserBindRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserBindRequestValidationError) ErrorName() string { return "UserBindRequestValidationError" }

// Error satisfies the builtin error interface
func (e UserBindRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserBindRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserBindRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserBindRequestValidationError{}

var _UserBindRequest_Mobile_Pattern = regexp.MustCompile("(^$|^[\\d]{11}$)")

// Validate checks the field values on TaobaoExtra with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *TaobaoExtra) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for UserId

	// no validation rules for OpenSid

	// no validation rules for TopAccessToken

	// no validation rules for AvatarUrl

	// no validation rules for HavanaSsoToken

	// no validation rules for Nick

	// no validation rules for OpenId

	// no validation rules for TopAuthCode

	// no validation rules for TopExpireTime

	return nil
}

// TaobaoExtraValidationError is the validation error returned by
// TaobaoExtra.Validate if the designated constraints aren't met.
type TaobaoExtraValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TaobaoExtraValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TaobaoExtraValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TaobaoExtraValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TaobaoExtraValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TaobaoExtraValidationError) ErrorName() string { return "TaobaoExtraValidationError" }

// Error satisfies the builtin error interface
func (e TaobaoExtraValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTaobaoExtra.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TaobaoExtraValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TaobaoExtraValidationError{}

// Validate checks the field values on WechatExtra with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *WechatExtra) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for AccessToken

	// no validation rules for ExpiresIn

	// no validation rules for RefreshToken

	// no validation rules for OpenId

	// no validation rules for Scope

	// no validation rules for NickName

	// no validation rules for Sex

	// no validation rules for Province

	// no validation rules for City

	// no validation rules for Country

	// no validation rules for Headimgurl

	// no validation rules for Unionid

	return nil
}

// WechatExtraValidationError is the validation error returned by
// WechatExtra.Validate if the designated constraints aren't met.
type WechatExtraValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e WechatExtraValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e WechatExtraValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e WechatExtraValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e WechatExtraValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e WechatExtraValidationError) ErrorName() string { return "WechatExtraValidationError" }

// Error satisfies the builtin error interface
func (e WechatExtraValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sWechatExtra.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = WechatExtraValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = WechatExtraValidationError{}

// Validate checks the field values on UserRefreshRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *UserRefreshRequest) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetDevice()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserRefreshRequestValidationError{
				field:  "Device",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Channel

	if utf8.RuneCountInString(m.GetVersionCode()) < 4 {
		return UserRefreshRequestValidationError{
			field:  "VersionCode",
			reason: "value length must be at least 4 runes",
		}
	}

	if !_UserRefreshRequest_VersionCode_Pattern.MatchString(m.GetVersionCode()) {
		return UserRefreshRequestValidationError{
			field:  "VersionCode",
			reason: "value does not match regex pattern \"^[\\\\d]+$\"",
		}
	}

	return nil
}

// UserRefreshRequestValidationError is the validation error returned by
// UserRefreshRequest.Validate if the designated constraints aren't met.
type UserRefreshRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserRefreshRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserRefreshRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserRefreshRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserRefreshRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserRefreshRequestValidationError) ErrorName() string {
	return "UserRefreshRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UserRefreshRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserRefreshRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserRefreshRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserRefreshRequestValidationError{}

var _UserRefreshRequest_VersionCode_Pattern = regexp.MustCompile("^[\\d]+$")

// Validate checks the field values on UserUnbindRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *UserUnbindRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Mobile

	// no validation rules for Wechat

	// no validation rules for Taobao

	return nil
}

// UserUnbindRequestValidationError is the validation error returned by
// UserUnbindRequest.Validate if the designated constraints aren't met.
type UserUnbindRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserUnbindRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserUnbindRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserUnbindRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserUnbindRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserUnbindRequestValidationError) ErrorName() string {
	return "UserUnbindRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UserUnbindRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserUnbindRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserUnbindRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserUnbindRequestValidationError{}

// Validate checks the field values on UserLoginRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *UserLoginRequest) Validate() error {
	if m == nil {
		return nil
	}

	if !_UserLoginRequest_Mobile_Pattern.MatchString(m.GetMobile()) {
		return UserLoginRequestValidationError{
			field:  "Mobile",
			reason: "value does not match regex pattern \"(^$|^[\\\\d]{11}$)\"",
		}
	}

	// no validation rules for Code

	// no validation rules for Wechat

	if v, ok := interface{}(m.GetDevice()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserLoginRequestValidationError{
				field:  "Device",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Channel

	if utf8.RuneCountInString(m.GetVersionCode()) < 4 {
		return UserLoginRequestValidationError{
			field:  "VersionCode",
			reason: "value length must be at least 4 runes",
		}
	}

	if !_UserLoginRequest_VersionCode_Pattern.MatchString(m.GetVersionCode()) {
		return UserLoginRequestValidationError{
			field:  "VersionCode",
			reason: "value does not match regex pattern \"^[\\\\d]+$\"",
		}
	}

	if utf8.RuneCountInString(m.GetPackageName()) < 1 {
		return UserLoginRequestValidationError{
			field:  "PackageName",
			reason: "value length must be at least 1 runes",
		}
	}

	// no validation rules for ThirdPartyId

	return nil
}

// UserLoginRequestValidationError is the validation error returned by
// UserLoginRequest.Validate if the designated constraints aren't met.
type UserLoginRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserLoginRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserLoginRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserLoginRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserLoginRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserLoginRequestValidationError) ErrorName() string { return "UserLoginRequestValidationError" }

// Error satisfies the builtin error interface
func (e UserLoginRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserLoginRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserLoginRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserLoginRequestValidationError{}

var _UserLoginRequest_Mobile_Pattern = regexp.MustCompile("(^$|^[\\d]{11}$)")

var _UserLoginRequest_VersionCode_Pattern = regexp.MustCompile("^[\\d]+$")

// Validate checks the field values on Device with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Device) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Imei

	// no validation rules for Idfa

	// no validation rules for AndroidId

	// no validation rules for Suuid

	// no validation rules for Mac

	// no validation rules for Os

	// no validation rules for Oaid

	return nil
}

// DeviceValidationError is the validation error returned by Device.Validate if
// the designated constraints aren't met.
type DeviceValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeviceValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeviceValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeviceValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeviceValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeviceValidationError) ErrorName() string { return "DeviceValidationError" }

// Error satisfies the builtin error interface
func (e DeviceValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDevice.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeviceValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeviceValidationError{}

// Validate checks the field values on UserInfo with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *UserInfo) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Id

	// no validation rules for UserName

	// no validation rules for Wechat

	// no validation rules for HeadImg

	// no validation rules for Gender

	// no validation rules for Birthday

	// no validation rules for Token

	// no validation rules for ThirdPartyId

	// no validation rules for IsNew

	if v, ok := interface{}(m.GetWechatExtra()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserInfoValidationError{
				field:  "WechatExtra",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetTaobaoExtra()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserInfoValidationError{
				field:  "TaobaoExtra",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Mobile

	// no validation rules for InviteCode

	// no validation rules for IsDeleted

	// no validation rules for IsInvited

	// no validation rules for Suuid

	// no validation rules for CreatedAt

	return nil
}

// UserInfoValidationError is the validation error returned by
// UserInfo.Validate if the designated constraints aren't met.
type UserInfoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserInfoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserInfoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserInfoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserInfoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserInfoValidationError) ErrorName() string { return "UserInfoValidationError" }

// Error satisfies the builtin error interface
func (e UserInfoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserInfo.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserInfoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserInfoValidationError{}

// Validate checks the field values on UserInfoReply with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *UserInfoReply) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Code

	// no validation rules for Message

	if v, ok := interface{}(m.GetData()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserInfoReplyValidationError{
				field:  "Data",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Msg

	return nil
}

// UserInfoReplyValidationError is the validation error returned by
// UserInfoReply.Validate if the designated constraints aren't met.
type UserInfoReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserInfoReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserInfoReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserInfoReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserInfoReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserInfoReplyValidationError) ErrorName() string { return "UserInfoReplyValidationError" }

// Error satisfies the builtin error interface
func (e UserInfoReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserInfoReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserInfoReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserInfoReplyValidationError{}

// Validate checks the field values on UserInfoBatchReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *UserInfoBatchReply) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Code

	for idx, item := range m.GetData() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UserInfoBatchReplyValidationError{
					field:  fmt.Sprintf("Data[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for Msg

	// no validation rules for Count

	return nil
}

// UserInfoBatchReplyValidationError is the validation error returned by
// UserInfoBatchReply.Validate if the designated constraints aren't met.
type UserInfoBatchReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserInfoBatchReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserInfoBatchReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserInfoBatchReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserInfoBatchReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserInfoBatchReplyValidationError) ErrorName() string {
	return "UserInfoBatchReplyValidationError"
}

// Error satisfies the builtin error interface
func (e UserInfoBatchReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserInfoBatchReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserInfoBatchReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserInfoBatchReplyValidationError{}

// Validate checks the field values on GetCodeRequest with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *GetCodeRequest) Validate() error {
	if m == nil {
		return nil
	}

	if !_GetCodeRequest_Mobile_Pattern.MatchString(m.GetMobile()) {
		return GetCodeRequestValidationError{
			field:  "Mobile",
			reason: "value does not match regex pattern \"\\\\d{11}\"",
		}
	}

	// no validation rules for PackageName

	return nil
}

// GetCodeRequestValidationError is the validation error returned by
// GetCodeRequest.Validate if the designated constraints aren't met.
type GetCodeRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCodeRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCodeRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCodeRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCodeRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCodeRequestValidationError) ErrorName() string { return "GetCodeRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetCodeRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCodeRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCodeRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCodeRequestValidationError{}

var _GetCodeRequest_Mobile_Pattern = regexp.MustCompile("\\d{11}")

// Validate checks the field values on UserInfoBatchRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *UserInfoBatchRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for PackageName

	// no validation rules for After

	// no validation rules for Before

	// no validation rules for Mobile

	// no validation rules for Name

	// no validation rules for PerPage

	// no validation rules for Page

	return nil
}

// UserInfoBatchRequestValidationError is the validation error returned by
// UserInfoBatchRequest.Validate if the designated constraints aren't met.
type UserInfoBatchRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserInfoBatchRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserInfoBatchRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserInfoBatchRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserInfoBatchRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserInfoBatchRequestValidationError) ErrorName() string {
	return "UserInfoBatchRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UserInfoBatchRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserInfoBatchRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserInfoBatchRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserInfoBatchRequestValidationError{}

// Validate checks the field values on UserInfoRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *UserInfoRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Id

	// no validation rules for Wechat

	// no validation rules for Taobao

	return nil
}

// UserInfoRequestValidationError is the validation error returned by
// UserInfoRequest.Validate if the designated constraints aren't met.
type UserInfoRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserInfoRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserInfoRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserInfoRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserInfoRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserInfoRequestValidationError) ErrorName() string { return "UserInfoRequestValidationError" }

// Error satisfies the builtin error interface
func (e UserInfoRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserInfoRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserInfoRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserInfoRequestValidationError{}

// Validate checks the field values on UserInfoUpdateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *UserInfoUpdateRequest) Validate() error {
	if m == nil {
		return nil
	}

	if l := utf8.RuneCountInString(m.GetUserName()); l < 2 || l > 10 {
		return UserInfoUpdateRequestValidationError{
			field:  "UserName",
			reason: "value length must be between 2 and 10 runes, inclusive",
		}
	}

	if !_UserInfoUpdateRequest_HeadImg_Pattern.MatchString(m.GetHeadImg()) {
		return UserInfoUpdateRequestValidationError{
			field:  "HeadImg",
			reason: "value does not match regex pattern \"^(|https?://.*)$\"",
		}
	}

	// no validation rules for Gender

	if !_UserInfoUpdateRequest_Birthday_Pattern.MatchString(m.GetBirthday()) {
		return UserInfoUpdateRequestValidationError{
			field:  "Birthday",
			reason: "value does not match regex pattern \"^(|\\\\d{4}-\\\\d{1,2}-\\\\d{1,2})$\"",
		}
	}

	// no validation rules for ThirdPartyId

	return nil
}

// UserInfoUpdateRequestValidationError is the validation error returned by
// UserInfoUpdateRequest.Validate if the designated constraints aren't met.
type UserInfoUpdateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserInfoUpdateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserInfoUpdateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserInfoUpdateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserInfoUpdateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserInfoUpdateRequestValidationError) ErrorName() string {
	return "UserInfoUpdateRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UserInfoUpdateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserInfoUpdateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserInfoUpdateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserInfoUpdateRequestValidationError{}

var _UserInfoUpdateRequest_HeadImg_Pattern = regexp.MustCompile("^(|https?://.*)$")

var _UserInfoUpdateRequest_Birthday_Pattern = regexp.MustCompile("^(|\\d{4}-\\d{1,2}-\\d{1,2})$")

// Validate checks the field values on EmptyRequest with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *EmptyRequest) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// EmptyRequestValidationError is the validation error returned by
// EmptyRequest.Validate if the designated constraints aren't met.
type EmptyRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EmptyRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EmptyRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EmptyRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EmptyRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EmptyRequestValidationError) ErrorName() string { return "EmptyRequestValidationError" }

// Error satisfies the builtin error interface
func (e EmptyRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEmptyRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EmptyRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EmptyRequestValidationError{}

// Validate checks the field values on GenericReply with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *GenericReply) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Code

	// no validation rules for Message

	// no validation rules for Msg

	return nil
}

// GenericReplyValidationError is the validation error returned by
// GenericReply.Validate if the designated constraints aren't met.
type GenericReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GenericReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GenericReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GenericReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GenericReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GenericReplyValidationError) ErrorName() string { return "GenericReplyValidationError" }

// Error satisfies the builtin error interface
func (e GenericReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGenericReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GenericReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GenericReplyValidationError{}

// Validate checks the field values on UserSoftDeleteRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *UserSoftDeleteRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Id

	return nil
}

// UserSoftDeleteRequestValidationError is the validation error returned by
// UserSoftDeleteRequest.Validate if the designated constraints aren't met.
type UserSoftDeleteRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserSoftDeleteRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserSoftDeleteRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserSoftDeleteRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserSoftDeleteRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserSoftDeleteRequestValidationError) ErrorName() string {
	return "UserSoftDeleteRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UserSoftDeleteRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserSoftDeleteRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserSoftDeleteRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserSoftDeleteRequestValidationError{}
