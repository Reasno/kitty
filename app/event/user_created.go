package event

import pb "glab.tagtic.cn/ad_gains/kitty/proto"

type UserCreated struct {
	*pb.UserInfoDetail
}
