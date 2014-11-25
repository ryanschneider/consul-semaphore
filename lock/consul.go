/*
	Port of https://github.com/coreos/locksmith/blob/master/lock/etcd.go
	from etcd to Consul
*/

package lock

import (
	"encoding/json"
	api "github.com/armon/consul-api"
)

// ConsulLockClient is a wrapper around the consul-api client
// that provides simple primitives to operate on a named semaphore
// stored as a Consul KV.
type ConsulLockClient struct {
	Path   string
	client api.Client
	kv     api.KV
}

func (c *ConsulLockClient) Init() (err error) {
	sem := newSemaphore()
	b, err := json.Marshal(sem)
	if err != nil {
		return err
	}

	client, _ := api.NewClient(api.DefaultConfig())
	c.client = *client
	c.kv = *client.KV()

	pair, _, err := c.kv.Get(c.Path, nil)
	if err != nil {
		return err
	}

	if pair == nil {
		p := &api.KVPair{Key: c.Path, Value: b}
		_, err := c.kv.Put(p, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
