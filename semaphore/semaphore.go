package semaphore

import (
	api "github.com/armon/consul-api"
	lock "github.com/ryanschneider/consul-semaphore/lock"
	"log"
)

type Semaphore struct {
	Path string
	lock *lock.Lock
}

func New(path string, holder string) (s *Semaphore, err error) {
	apiClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	client, err := lock.NewConsulLockClient(apiClient)
	if err != nil {
		return nil, err
	}

	lock, err := lock.New(path, holder, client)
	if err != nil {
		return nil, err
	}

	return &Semaphore{path, lock}, nil
}

func (s *Semaphore) SetMax(max uint) (oldMax uint, err error) {
	_, iOldMax, err := s.lock.SetMax(int(max))
	if err != nil {
		return 0, err
	}

	oldMax = uint(iOldMax)
	return oldMax, nil
}

func (s *Semaphore) Acquire(wait bool) (err error) {
	for {
		err = s.lock.Lock()
		if err == nil {
			return nil
		}

		// only go again if we are waiting
		if !wait {
			return err
		}

		_, isExhausted := err.(lock.SemaphoreExhaustedErr)
		casFailed := (err == lock.CheckAndSetFailedErr)

		switch {
		case isExhausted:
			log.Printf("Exhausted, trying again")
		case casFailed:
			log.Printf("CAS failed, trying again")
		default:
			return err
		}

		changed, err := s.lock.Watch()
		if err != nil {
			return err
		}

		log.Printf("Watch woke up, changed: %v", changed)
		if changed {
			// TODO: add some random sleep here to avoid
			// too many CAS errors on thundering herd
			continue
		}
	}

	return err
}

func (s *Semaphore) Release() (err error) {
	for {
		err = s.lock.Unlock()
		if err == lock.CheckAndSetFailedErr {
			// TODO: add some sleep here to avoid
			// too many CAS errors on thundering herd
			continue
		}
		return
	}
}
