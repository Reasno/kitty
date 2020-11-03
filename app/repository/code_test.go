package repository

import (
	"context"
	"flag"
	"github.com/Reasno/kitty/pkg/config"
	"github.com/Reasno/kitty/pkg/otredis"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

var cr *CodeRepo
var useRedis bool

func init() {
	flag.BoolVar(&useRedis, "redis", false, "use local redis for testing")
}

func TestCodeRepoFrequencyLimit(t *testing.T) {
	if !useRedis {
		return
	}
	client := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs: []string{"127.0.0.1:6379"},
		})
	client.FlushAll(context.Background())
	cr = &CodeRepo{
		client,
		otredis.NewKeyManager(":", "test"),
		2 * time.Second,
		time.Second,
		config.Env("testing"),
	}
	_, err := cr.AddCode(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = cr.AddCode(context.Background(), "1")
	if err != ErrTooFrequent {
		t.Fatal("should receive ErrTooFrequent")
	}
	_, err = cr.AddCode(context.Background(), "2")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	_, err = cr.AddCode(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	_, err = cr.AddCode(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
}
