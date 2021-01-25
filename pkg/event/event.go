package event

import (
	"context"
	"fmt"

	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type Evt struct {
	body interface{}
	ctx  context.Context
}

func (e Evt) Data() interface{} {
	return e.body
}

func (e Evt) Type() string {
	return fmt.Sprintf("%T", e.body)
}

func (e Evt) Context() context.Context {
	return e.ctx
}

func NewEvent(ctx context.Context, i interface{}) Evt {
	return Evt{
		body: i,
		ctx:  ctx,
	}
}

func Of(i ...interface{}) []contract.Event {
	var out []contract.Event
	for _, ii := range i {
		out = append(out, NewEvent(context.Background(), ii))
	}
	return out
}
