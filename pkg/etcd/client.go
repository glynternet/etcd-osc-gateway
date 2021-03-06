package etcd

import (
	"context"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

// Client is an etcd KV client
type Client struct {
	clientv3.KV
}

// Put puts the key and value into the etcd store
func (c Client) Put(ctx context.Context, key, value string) error {
	_, err := c.KV.Put(ctx, key, value)
	return errors.Wrap(err, "getting response from store")
}
