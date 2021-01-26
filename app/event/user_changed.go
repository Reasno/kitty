package event

import pb "glab.tagtic.cn/ad_gains/kitty/proto"

type UserChanged struct {
	*pb.UserInfoDetail
}
