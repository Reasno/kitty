package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"time"
)

type CodeRepo struct {
	client redis.Cmdable
}

func NewCodeRepo(cmdable redis.Cmdable) *CodeRepo {
	return &CodeRepo{cmdable}
}

func (c *CodeRepo) AddCode(ctx context.Context, mobile string) (code string, err error) {
	n := rand.Intn(1_000_000)
	code = pad(n)
	_, err = c.client.Set(ctx, "CodeRepo:"+mobile, code, 15*time.Minute).Result()
	if err != nil {
		return "", errors.Wrap(err, "cannot persist code in redis")
	}
	return code, nil
}

func (c *CodeRepo) CheckCode(ctx context.Context, mobile, code string) (bool, error) {
	value, err := c.client.Get(ctx, "CodeRepo:"+mobile).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "cannot query code in redis")
	}
	return value == code, nil
}

func (c *CodeRepo) DeleteCode(ctx context.Context, mobile string) (err error) {
	_, err = c.client.Del(ctx, "CodeRepo:"+mobile).Result()
	return err
}


func pad(n int) string {
	s := strconv.Itoa(n)
	for len(s) < 6 {
		s = "0" + s
	}
	return s
}
