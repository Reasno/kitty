package rule

import (
	"context"
	"flag"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/magiconair/properties/assert"
	"go.etcd.io/etcd/clientv3"
)

var useEtcd bool

func init() {
	flag.BoolVar(&useEtcd, "etcd", false, "use local mysql for testing")
}

func TestRepository_WatchConfigUpdate(t *testing.T) {
	if !useEtcd {
		t.Skip("test dynamic config requires etcd")
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		t.Fatal(err)
	}
	client.Put(context.Background(), CentralConfigPath, `
style: basic
rule:
  list:
    - name: 商业化平台
      icon: home2
      children:
        - name: 用户体系
          path: /kitty
          id: user
`)
	client.Put(context.Background(), OtherConfigPathPrefix+"/kitty-testing", `
style: basic
rule:
  foo: bar
`)
	client.Put(context.Background(), OtherConfigPathPrefix+"/egg-testing", `
style: basic
rule:
  foo: qux
`)
	repo, err := NewRepository(client, log.NewNopLogger())
	if err != nil {
		t.Fatal(err)
	}
	repo.updateChan = make(chan struct{})
	watchCxt, cancel := context.WithCancel(context.Background())
	defer cancel()
	go repo.WatchConfigUpdate(watchCxt)

	type caseList []struct {
		name    string
		dbKey   string
		dataFoo interface{}
	}

	cases := caseList{
		{
			"central-config",
			CentralConfigPath,
			nil,
		},
		{
			"kitty-testing",
			OtherConfigPathPrefix + "/kitty-testing",
			"bar",
		},
	}
	for _, c := range cases {
		assert.Equal(t, repo.containers[c.name].DbKey, c.dbKey)
		assert.Equal(t, repo.containers[c.name].Name, c.name)
		assert.Equal(t, repo.containers[c.name].RuleSet[0].Then["foo"], c.dataFoo)
	}

	_, err = client.Delete(context.Background(), "/monetization/kitty-testing")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Put(context.Background(), "/monetization/kitty-testing", `
style: basic
rule:
  foo: baz`)
	if err != nil {
		t.Fatal(err)
	}

	<-repo.updateChan

	cases = caseList{
		{
			"central-config",
			CentralConfigPath,
			nil,
		},
		{
			"kitty-testing",
			OtherConfigPathPrefix + "/kitty-testing",
			"baz",
		},
	}
	for _, c := range cases {
		assert.Equal(t, repo.containers[c.name].DbKey, c.dbKey)
		assert.Equal(t, repo.containers[c.name].Name, c.name)
		assert.Equal(t, repo.containers[c.name].RuleSet[0].Then["foo"], c.dataFoo)
	}

	client.Put(context.Background(), CentralConfigPath, `
style: basic
rule:
  list:
    - name: 商业化平台
      icon: home2
      children:
        - name: 用户体系
          path: /kitty
          id: user
        - name: 积分体系
          path: /score
          id: score
    - name: 活动
      icon: material
      children:
        - name: 砸金蛋
          path: /egg
          id: egg
        - name: 惊喜福利砸中你
          path: /surprise
          id: surprise
`)
	<-repo.updateChan
	_, ok := repo.containers["egg-local"]
	if !ok {
		t.Fatal("egg should exist")
	}
	_, ok = repo.containers["score-local"]
	if !ok {
		t.Fatal("score should exist")
	}
	_, ok = repo.containers["surprise-local"]
	if !ok {
		t.Fatal("surprise should exist")
	}
}
