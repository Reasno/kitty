package consumer

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kkafka"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

type EventReceiver struct {
	AppName contract.AppName
	Conf    contract.ConfigReader
	Manager InvitationManager
	Factory *kkafka.KafkaFactory
}

type InvitationManager interface {
	CompleteStep(ctx context.Context, apprenticeId uint64, eventName string) error
}

func (er *EventReceiver) handleSign(ctx context.Context, msg kafka.Message) error {
	var signEvent pb.SignEvent
	err := signEvent.Unmarshal(msg.Value)
	if err != nil {
		return err
	}
	ctx = withTenant(ctx, &signEvent)
	return er.Manager.CompleteStep(ctx, signEvent.UserId, signEvent.EventName)
}

func (er *EventReceiver) handleTask(ctx context.Context, msg kafka.Message) error {
	var taskEvent pb.TaskEvent
	err := taskEvent.Unmarshal(msg.Value)
	if err != nil {
		return err
	}
	ctx = withTenant(ctx, &taskEvent)
	return er.Manager.CompleteStep(ctx, taskEvent.UserId, taskEvent.EventName)
}

func (er *EventReceiver) SubscribeTask(ctx context.Context) error {
	groupId := strings.Join([]string{"EventReceiver", er.AppName.String()}, "-")
	return er.Factory.Reader(
		er.Conf.String("kafka.taskEventBus"),
		kkafka.WithGroup(groupId),
	).Subscribe(ctx, er.handleTask)
}

func (er *EventReceiver) SubscribeCheckin(ctx context.Context) error {
	groupId := strings.Join([]string{"EventReceiver", er.AppName.String()}, "-")
	return er.Factory.Reader(
		er.Conf.String("kafka.signEventBus"),
		kkafka.WithGroup(groupId),
	).Subscribe(ctx, er.handleSign)
}

type Tenanter interface {
	GetChannel() string
	GetUserId() uint64
	GetPackageName() string
}

func withTenant(ctx context.Context, t Tenanter) context.Context {
	return context.WithValue(ctx, config.TenantKey, &config.Tenant{
		Channel:     t.GetChannel(),
		UserId:      t.GetUserId(),
		PackageName: t.GetPackageName(),
	})
}
