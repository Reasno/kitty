package client

import (
	"context"
	"flag"
	"testing"

	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
	repository2 "glab.tagtic.cn/ad_gains/kitty/rule/repository"
	"go.etcd.io/etcd/clientv3"
)

var useEtcd bool

func init() {
	flag.BoolVar(&useEtcd, "etcd", false, "use local mysql for testing")
}

type OrientationEvent struct {
	Id          int
	AnotherName string `koanf:"another_name"`
}

func TestClient(t *testing.T) {
	if !useEtcd {
		t.Skip("test dynamic config requires etcd")
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		t.Fatal(err)
	}
	client.Put(context.Background(), repository2.OtherConfigPathPrefix+"/kitty-testing", `
style: basic
rule:
  foo: bar
  orientation_events:
    - id: 5
      another_name: "str"
`)
	dynConf, err := NewRuleEngine(WithClient(client), Rule("kitty-testing"))
	if err != nil {
		t.Fatal(err)
	}
	reader, err := dynConf.Of("kitty-testing").Payload(&dto.Payload{})
	if err != nil {
		t.Fatal(err)
	}
	if reader.String("foo") != "bar" {
		t.Fatalf("want %s, got %s", "foo", reader.String("foo"))
	}

	var sh []OrientationEvent
	err = reader.Unmarshal("orientation_events", &sh)
	if err != nil {
		t.Fatal(err)
	}
	if sh[0].AnotherName != "str" {
		t.Fatalf("want %s, got %s", "str", sh[0].AnotherName)
	}
}
