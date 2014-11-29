package semaphore

import (
	"log"

	api "github.com/armon/consul-api"
	lock "github.com/ryanschneider/consul-semaphore/lock"
)

// Semaphore represents a Consul-backed semaphore.
// Based off of the etcd semaphore used in CoreOS' Locksmith,
// Semaphore can be used to coordiate a set of workers around
// a Consul KV.  For exmple, restarting services in a consul-template
// action in a controlled manner.
type Semaphore struct {
	Path string
	lock *lock.Lock
}

// New creates and returns a new Semaphore.
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

// SetMax sets the maximum number of concurrent holders of a Semaphore.
// If the max is raised, multiple Acquirers may be signalled.  If the max is
// lowered below the current number of holders, no one will be signalled until
// the number of holders drops below max.
func (s *Semaphore) SetMax(max uint) (oldMax uint, err error) {
	_, iOldMax, err := s.lock.SetMax(int(max))
	if err != nil {
		return 0, err
	}

	oldMax = uint(iOldMax)
	return oldMax, nil
}

// Acquire acquires a portion of the Semaphore, optionally waiting if the
// Semaphore is currently maxed out.
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

// Releases releases a portion of the Semaphore.
// Releasing allows waiting Acquirers to be signalled.
// Note: In a highly contentious Semaphore, there may be CheckAndSet (CAS)
// errors writing to the semaphore.  These are handled inside Release, which
// may lead to Release blocking while it attempts to cleanly write to the KV.
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
