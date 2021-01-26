package listener

import (
	"context"

	appevent "glab.tagtic.cn/ad_gains/kitty/app/event"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/event"
)

type EventBus interface {
	Emit(ctx context.Context, event string, tenant *config.Tenant) error
}

type UserCreated struct {
	Bus EventBus
}

func (u UserCreated) Listen() []contract.Event {
	return event.Of(appevent.UserCreated{})
}

func (u UserCreated) Process(event contract.Event) error {
	data := event.Data().(appevent.UserCreated)
	claim := config.Tenant{
		PackageName: data.PackageName,
		UserId:      data.Id,
		Suuid:       data.Suuid,
		Channel:     data.Channel,
		VersionCode: data.VersionCode,
	}
	_ = u.Bus.Emit(event.Context(), "new_user", &claim)
	return nil
}
