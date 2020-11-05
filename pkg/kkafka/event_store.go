package kkafka

import (
	"context"
)

type EventStore struct {
	factory *KafkaProducerFactory
}

func (e *EventStore) Emit(ctx context.Context, event string) {
	//claim := kjwt.GetClaim(ctx)
	//dto := &Message{
	//	Timestamp:   time.Now().UTC().Format(time.RFC3339),
	//	Suuid:       claim.,
	//	VersionCode: "",
	//	Channel:     "",
	//	Event:       "",
	//	UserId:      "",
	//	PackageName: "",
	//}
}
