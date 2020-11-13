package client

import (
	"context"
	"flag"
	"testing"

	"glab.tagtic.cn/ad_gains/kitty/rule"
	"go.etcd.io/etcd/clientv3"
)

var useEtcd bool

func init() {
	flag.BoolVar(&useEtcd, "etcd", false, "use local mysql for testing")
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
	client.Put(context.Background(), rule.OtherConfigPathPrefix+"/kitty-testing", `
style: basic
rule:
  foo: bar
`)
	dynConf, err := NewRuleEngine(WithClient(client), Rule("kitty-testing"))
	if err != nil {
		t.Fatal(err)
	}
	reader, err := dynConf.Of("kitty-testing").Payload(&rule.Payload{})
	if err != nil {
		t.Fatal(err)
	}
	if reader.String("foo") != "bar" {
		t.Fatalf("want %s, got %s", "foo", reader.String("foo"))
	}
}
