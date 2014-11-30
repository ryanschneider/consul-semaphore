package semaphore

import (
	"log"
	"math/rand"
	"time"

	api "github.com/armon/consul-api"
	lock "github.com/ryanschneider/consul-semaphore/lock"
)

// Semaphore represents a Consul-backed semaphore.
// Based off of the etcd semaphore used in CoreOS' Locksmith,
// Semaphore can be used to coordiate a set of workers around
// a Consul KV.  For exmple, restarting services in a consul-template
// action in a controlled manner.
type Semaphore struct {
	Path   string
	Holder string
	lock   *lock.Lock
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

	return &Semaphore{path, holder, lock}, nil
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
		log.Printf("Holder %v: acquiring..", s.Holder)
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
			log.Printf("Holder %v: Semaphore exhausted, trying again", s.Holder)
		case casFailed:
			log.Printf("Holder %v: CheckAndSet failed, trying again", s.Holder)
		default:
			return err
		}

		changed, err := s.lock.Watch()
		if err != nil {
			return err
		}

		log.Printf("Holder %v: Watch woke up, changed: %v", s.Holder, changed)
		if changed {
			// Sleep here to avoid too many CAS errors on thundering herd
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			time.Sleep(time.Duration(r.Intn(1000)) * time.Millisecond)
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
			// Sleep here to avoid too many CAS errors on thundering herd
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			time.Sleep(time.Duration(r.Intn(1000)) * time.Millisecond)
			continue
		}
		return
	}
}
