package repository

import (
	"context"
	"fmt"
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/Reasno/kitty/pkg/otredis"
	"github.com/go-redis/redis/v8"
	"time"
)

const ExtraKey = "ExtraRepo"
const ExtraTTL = 30 * 24 * time.Hour

type ExtraRepo struct {
	client redis.Cmdable
	km     contract.Keyer
	ttl    time.Duration
}

func NewExtraRepo(client redis.Cmdable, keyer contract.Keyer) *ExtraRepo {
	return &ExtraRepo{
		client: client,
		km:     otredis.With(keyer, ExtraKey),
		ttl:    ExtraTTL,
	}
}

func (e *ExtraRepo) Put(ctx context.Context, id uint, name string, extra []byte) error {
	key := e.km.Key(fmt.Sprintf("%d", id), name)
	_, err := e.client.Set(ctx, key, extra, e.ttl).Result()
	return err
}

func (e *ExtraRepo) Get(ctx context.Context, id uint, name string) ([]byte, error) {
	key := e.km.Key(fmt.Sprintf("%d", id), name)
	b, err := e.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return []byte{}, nil
	}
	return []byte(b), err
}
