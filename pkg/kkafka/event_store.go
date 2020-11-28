package kkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
)

type EventStore struct {
	Factory *KafkaFactory
	Topic   string
	MW      Middleware
	once    sync.Once
	handler Handler
}

func (e *EventStore) Emit(ctx context.Context, event string, claim *kjwt.Claim) error {
	e.once.Do(func() {
		e.handler = e.Factory.MakeHandler(e.Topic)
	})
	dto := &Message{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Suuid:       claim.Suuid,
		VersionCode: claim.VersionCode,
		Channel:     claim.Channel,
		Event:       event,
		UserId:      fmt.Sprintf("%d", claim.UserId),
		PackageName: claim.PackageName,
		AppKey:      "appwangzhuan",
	}
	b, err := json.Marshal(dto)
	if err != nil {
		return errors.Wrap(err, "unable to marshal dto")
	}

	return e.MW(e.handler).Handle(ctx, kafka.Message{
		Value: b,
	})
}
