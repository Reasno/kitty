package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kkafka"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

func MakeDataInfoEndpoint(store kkafka.DataStore) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UserInfo)
		return nil, store.Emit(ctx, req)
	}
}

func encodeUserInfoRequest(_ context.Context, msg *kafka.Message, request interface{}) error {
	marshaller := request.(contract.Marshaller)
	byt, err := marshaller.Marshal()
	if err != nil {
		return err
	}
	msg.Value = byt
	return nil
}

type Event struct {
	Name   string
	Tenant *config.Tenant
}

func encodeEventRequest(_ context.Context, msg *kafka.Message, request interface{}) error {
	req := request.(*Event)
	dto := &kkafka.Message{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Suuid:       req.Tenant.Suuid,
		VersionCode: req.Tenant.VersionCode,
		Channel:     req.Tenant.Channel,
		Event:       req.Name,
		UserId:      fmt.Sprintf("%d", req.Tenant.UserId),
		PackageName: req.Tenant.PackageName,
		AppKey:      "appwangzhuan",
	}
	b, err := json.Marshal(dto)
	if err != nil {
		return errors.Wrap(err, "unable to marshal dto")
	}
	msg.Value = b
	return nil
}

type DataStore struct {
	e endpoint.Endpoint
}

func NewDataStore(topic string, factory *kkafka.KafkaFactory, options []kkafka.PublisherOption, mw endpoint.Middleware) *DataStore {
	pub := kkafka.NewPublisher(
		factory.MakeHandler(topic),
		encodeUserInfoRequest,
		options...,
	).Endpoint()

	pub = mw(pub)
	return &DataStore{pub}
}

func (d *DataStore) Emit(ctx context.Context, user contract.Marshaller) error {
	_, err := d.e(ctx, user)
	return err
}

type EventStore struct {
	e endpoint.Endpoint
}

func NewEventStore(topic string, factory *kkafka.KafkaFactory, options []kkafka.PublisherOption, mw endpoint.Middleware) *EventStore {
	pub := kkafka.NewPublisher(
		factory.MakeHandler(topic),
		encodeEventRequest,
		options...,
	).Endpoint()
	return &EventStore{mw(pub)}
}

func (e *EventStore) Emit(ctx context.Context, event string, tenant *config.Tenant) error {
	_, err := e.e(ctx, &Event{Name: event, Tenant: tenant})
	return err
}
