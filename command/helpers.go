package command

import (
	api "github.com/armon/consul-api"
	"github.com/ryanschneider/consul-semaphore/lock"
)

func getClient() (client *lock.ConsulLockClient, err error) {
	apiClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	client, err = lock.NewConsulLockClient(apiClient)
	return
}
