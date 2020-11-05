package kkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Reasno/kitty/pkg/kjwt"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type EventStore struct {
	Factory *KafkaProducerFactory
	Topic string
}

func (e *EventStore) Emit(ctx context.Context, event string) error {
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

	return e.Factory.Writer(e.Topic).WriteMessages(ctx, kafka.Message{
		Value: b,
	})
}
