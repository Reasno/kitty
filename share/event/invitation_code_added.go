package event

import (
	"time"
)

type InvitationCodeAdded struct {
	InviteeId   uint64    // 受邀者
	InviterId   uint64    // 邀请人
	InviteCode  string    // 邀请码
	PackageName string    // 应用包名
	DateTime    time.Time // 时间
	Channel     string    // 渠道
}
