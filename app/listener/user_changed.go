package listener

import (
	"context"

	appevent "glab.tagtic.cn/ad_gains/kitty/app/event"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/event"
)

type UserBus interface {
	Emit(ctx context.Context, info contract.Marshaller) error
}

type UserChanged struct {
	Bus UserBus
}

func (u UserChanged) Listen() []contract.Event {
	return event.Of(appevent.UserChanged{})
}

func (u UserChanged) Process(event contract.Event) error {
	data := event.Data().(appevent.UserChanged)
	_ = u.Bus.Emit(event.Context(), data)
	return nil
}
