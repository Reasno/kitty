package client

import (
	"context"
	"flag"
	"fmt"
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
		Endpoints: []string{"etcd-1:2379", "etcd-2:2379", "etcd-3:2379"},
		Context:   context.Background(),
	})
	if err != nil {
		t.Fatal(err)
	}
	dynConf, err := NewRuleEngine(WithClient(client), Rule("kitty-testing"))
	if err != nil {
		t.Fatal(err)
	}
	reader, err := dynConf.Of("kitty-testing").Payload(&rule.Payload{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(reader.Get("foo"))
}
