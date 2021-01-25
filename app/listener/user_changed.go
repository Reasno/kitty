package listener

import (
	"context"

	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/event"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

type UserBus interface {
	Emit(ctx context.Context, info contract.Marshaller) error
}

type UserChanged struct {
	Bus UserBus
}

func (u UserChanged) Listen() []contract.Event {
	return event.Of(&pb.UserInfoDetail{})
}

func (u UserChanged) Process(event contract.Event) error {
	data := event.Data().(*pb.UserInfoDetail)
	_ = u.Bus.Emit(event.Context(), data)
	return nil
}
