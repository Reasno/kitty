package client

import (
	"context"
	"flag"
	"regexp"
	"testing"

	kconf "glab.tagtic.cn/ad_gains/kitty/pkg/config"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
	repository2 "glab.tagtic.cn/ad_gains/kitty/rule/repository"
	"go.etcd.io/etcd/clientv3"
)

var useEtcd bool

func init() {
	flag.BoolVar(&useEtcd, "etcd", false, "use local etcd for testing")
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
	client.Put(context.Background(), repository2.OtherConfigPathPrefix+"/bar-kitty-testing", `
style: basic
rule:
  foo: baz
  orientation_events:
    - id: 5
      another_name: "str"
`)
	cases := []struct {
		name    string
		engine  *RuleEngine
		asserts func(t *testing.T, e *RuleEngine)
	}{
		{
			"normal",
			func() *RuleEngine {
				dynConf, err := NewRuleEngine(WithClient(client), Rule("kitty-testing"))
				if err != nil {
					t.Fatal(err)
				}
				return dynConf
			}(),
			func(t *testing.T, e *RuleEngine) {

			},
		},
		{
			"prefix",
			func() *RuleEngine {
				dynConf, err := NewRuleEngine(WithClient(client), WithRulePrefix("kitty"))
				if err != nil {
					t.Fatal(err)
				}
				return dynConf
			}(),
			func(t *testing.T, e *RuleEngine) {

			},
		},
		{
			"exp",
			func() *RuleEngine {
				exp := regexp.MustCompile(".*tty.*")
				dynConf, err := NewRuleEngine(WithClient(client), WithRuleRegexp(exp))
				if err != nil {
					t.Fatal(err)
				}
				return dynConf
			}(),
			func(t *testing.T, e *RuleEngine) {
			},
		},
		{
			"bar",
			func() *RuleEngine {
				exp := regexp.MustCompile(".*tty.*")
				dynConf, err := NewRuleEngine(WithClient(client), WithRuleRegexp(exp))
				if err != nil {
					t.Fatal(err)
				}
				return dynConf
			}(),
			func(t *testing.T, dynConf *RuleEngine) {
				reader, err := dynConf.Of("kitty-testing").Payload(&dto.Payload{
					PackageName: "com.foo.bar",
				})
				if err != nil {
					t.Fatal(err)
				}
				if reader.String("foo") != "baz" {
					t.Fatalf("want %s, got %s", "baz", reader.String("foo"))
				}
				reader, err = dynConf.Of("kitty-testing").Payload(&dto.Payload{
					PackageName: "com.foo.bar2",
				})
				if err != nil {
					t.Fatal(err)
				}
				if reader.String("foo") != "bar" {
					t.Fatalf("want %s, got %s", "bar", reader.String("foo"))
				}
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			dynConf := c.engine
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
			c.asserts(t, dynConf)
		})
	}
}

func TestClientIntegration(t *testing.T) {
	if !useEtcd {
		t.Skip("test dynamic config requires etcd")
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		t.Fatal(err)
	}
	client.Put(context.Background(), repository2.OtherConfigPathPrefix+"/kitty2-testing", `
style: advanced
enrich: true
rule:
- if: HoursAgo(DMP.Register) > 1
  then: 
    foo: bar
- if: true
  then:
    foo: quz
`)
	dynConf, err := NewRuleEngine(WithClient(client), Rule("kitty2-testing"), WithDMPAddr("dmp-grpc.dev.tagtic.cn:443"), WithEnv(kconf.Env("local")))
	if err != nil {
		t.Fatal(err)
	}
	conf, err := dynConf.Of("kitty2-testing").Payload(&dto.Payload{
		Channel:     "walk",
		Suuid:       "DoNews1a674f54-6889-4798-bddf-1cb5ca5c6164",
		PackageName: "com.walk.qnjb",
		DMP:         pb.DmpResp{},
		Context:     context.Background(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if conf.String("foo") != "bar" {
		t.Fatalf("want %s, got %s", "bar", conf.String("foo"))
	}
}
