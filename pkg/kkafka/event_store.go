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
	Factory *KafkaProducerFactory
	Topic   string
	once    sync.Once
	writer  *kafka.Writer
}

func (e *EventStore) Emit(ctx context.Context, event string) error {
	e.once.Do(func() {
		e.writer = e.Factory.Writer(e.Topic)
	})
	claim := kjwt.GetClaim(ctx)
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

	return e.writer.WriteMessages(ctx, kafka.Message{
		Value: b,
	})
}
