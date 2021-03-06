/*
	Port of https://github.com/coreos/locksmith/blob/master/lock/etcd.go
	from etcd to Consul
*/

package lock

import (
	"encoding/json"
	"errors"

	api "github.com/armon/consul-api"
)

var (
	CheckAndSetFailedErr = errors.New("Someone else modified the semaphore")
)

// ConsulLockClient is a wrapper around the consul-api client
// that provides simple primitives to operate on a named semaphore
// stored as a Consul KV.
type ConsulLockClient struct {
	Path   string
	client *api.Client
}

func NewConsulLockClient(apiClient *api.Client) (client *ConsulLockClient, err error) {
	client = &ConsulLockClient{client: apiClient}
	err = nil
	return
}

func (c *ConsulLockClient) SetPath(path string) error {
	c.Path = path
	return nil
}

func (c *ConsulLockClient) Init() (err error) {
	if c.Path == "" {
		return errors.New("cannot initialize semaphore without a path")
	}

	sem := newSemaphore()
	b, err := json.Marshal(sem)
	if err != nil {
		return err
	}

	kv := c.client.KV()

	pair, _, err := kv.Get(c.Path, nil)
	if err != nil {
		return err
	}

	if pair == nil {
		p := &api.KVPair{Key: c.Path, Value: b}
		_, err := kv.Put(p, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ConsulLockClient) Get() (sem *Semaphore, err error) {
	kv := c.client.KV()
	qo := &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
	}
	pair, _, err := kv.Get(c.Path, qo)
	if err != nil {
		return nil, err
	}

	sem = &Semaphore{}
	err = json.Unmarshal([]byte(pair.Value), sem)
	if err != nil {
		return nil, err
	}

	sem.Index = pair.ModifyIndex

	return sem, nil
}

func (c *ConsulLockClient) Set(sem *Semaphore) (err error) {
	if sem == nil {
		return errors.New("cannot set nil semaphore")
	}
	b, err := json.Marshal(sem)
	if err != nil {
		return err
	}

	pair := &api.KVPair{Key: c.Path, Value: b}
	pair.ModifyIndex = sem.Index

	kv := c.client.KV()

	written, _, err := kv.CAS(pair, nil)
	if err != nil {
		return err
	}

	if written != true {
		// TODO: Should I handle CAS errors in .Set()?
		// Makes the API cleaner, but makes Set potentially blocking
		return CheckAndSetFailedErr
	}

	return nil
}

func (c *ConsulLockClient) Watch(sem *Semaphore) (changed bool, err error) {

	kv := c.client.KV()
	qo := &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
		WaitIndex:         sem.Index,
	}
	pair, meta, err := kv.Get(c.Path, qo)
	if err != nil {
		return true, err
	}

	if meta == nil {
		return true, errors.New("Returned with nil metadata, see https://github.com/hashicorp/consul-template/issues/72")
	}

	changed = meta.LastIndex != sem.Index

	// NOTE: modifies input argument..
	err = json.Unmarshal([]byte(pair.Value), sem)
	if err != nil {
		return true, err
	}

	sem.Index = pair.ModifyIndex
	return
}
