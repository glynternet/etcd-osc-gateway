package etcd

import (
	"context"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/client"
)

// Client is an etcd KV client
type Client struct {
	client.KeysAPI
}

// Put puts the key and value into the etcd store
func (c Client) Put(ctx context.Context, key, value string) error {
	_, err := c.KeysAPI.Set(ctx, key, value, nil)
	return errors.Wrap(err, "getting response from store")
}
