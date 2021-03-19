package client

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/clientv3"
)

func TestGetRawRuleSetsFromPrefix(t *testing.T) {
	if !useEtcd {
		t.Skip("test dynamic config requires etcd")
	}
	client, _ := clientv3.New(clientv3.Config{
		Endpoints: []string{"etcd-1:2379"},
	})
	repo, _ := NewRepositoryWithConfig(client, log.NewNopLogger(), RepositoryConfig{
		Prefix: "",
		Limit:  100,
	})
	v, e := repo.getRawRuleSetsFromPrefix(context.Background())
	assert.NoError(t, e)
	assert.Less(t, 100, len(v))
}
