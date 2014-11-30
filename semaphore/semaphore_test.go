package semaphore

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"testing"
)

// TODO: Mock out consul-api so these aren't integration tests
// Will be harder since api.Client is a struct rather than an interface.

func TestSemaphore(t *testing.T) {
	const (
		path = "tests/integration/semaphore/TestSemaphore"
	)

	sem, err := New(path, "holder")
	if err != nil {
		t.Error(err)
	}

	oldMax, err := sem.SetMax(2)
	if err != nil {
		t.Error(err)
	}

	oldMax, err = sem.SetMax(1)
	if err != nil {
		t.Error(err)
	}

	if oldMax != 2 {
		t.Error("SetMax did not return expected oldMax")
	}

	err = sem.Acquire(false)
	if err != nil {
		t.Error(err)
	}

	err = sem.Release()
	if err != nil {
		t.Error(err)
	}
}

func TestAcquireWait(t *testing.T) {
	const (
		path = "tests/integration/semaphore/TestAcquireWait"
	)

	sem1, err := New(path, "1")
	if err != nil {
		t.Error(err)
	}

	sem2, err := New(path, "2")
	if err != nil {
		t.Error(err)
	}

	err = sem1.Acquire(false)
	if err != nil {
		t.Error(err)
	}

	go func() {
		err = sem2.Acquire(true)
		if err != nil {
			t.Error(err)
		}

		sem2.Release()
		if err != nil {
			t.Error(err)
		}
	}()

	err = sem1.Release()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkContention(b *testing.B) {
	const (
		path  = "tests/integration/semaphore/BenchmarkContention"
		count = 20
		max   = 1
	)

	log.SetOutput(ioutil.Discard)

	for x := 0; x < b.N; x++ {

		sem, err := New(path, "INITIAL")
		if err != nil {
			b.Error(err)
		}

		_, err = sem.SetMax(max)
		if err != nil {
			b.Error(err)
		}

		sems := make([]*Semaphore, 0, count)
		for i := 0; i < count; i++ {
			s, err := New(path, fmt.Sprintf("CONTENDER-%v", i))
			if err != nil {
				b.Error(err)
			}
			sems = append(sems, s)
		}

		err = sem.Acquire(false)
		if err != nil {
			b.Error(fmt.Sprintf("Error acquiring for holder %v: %v", sem.Holder, err))
		}

		b.Log(fmt.Sprintf("Holder %v acquired", sem.Holder))

		wg := sync.WaitGroup{}

		for _, s := range sems {
			wg.Add(1)

			go func(s *Semaphore) {
				defer wg.Done()

				ea := s.Acquire(true)
				if ea != nil {
					b.Error(fmt.Sprintf("Error acquiring for holder %v: %v", s.Holder, ea))
					return
				}

				b.Log(fmt.Sprintf("Holder %v acquired, releasing", s.Holder))

				er := s.Release()
				if er != nil {
					b.Error(fmt.Sprintf("Error releasing for holder %v: %v", s.Holder, er))
					return
				}
			}(s)
		}

		b.Log(fmt.Sprintf("Holder %v ...", sem.Holder))

		wg.Add(1)
		go func() {
			defer wg.Done()

			b.Log(fmt.Sprintf("Holder %v releasing", sem.Holder))

			er := sem.Release()
			if er != nil {
				b.Error(fmt.Sprintf("Error releasing for holder %v: %v", sem.Holder, er))
				return
			}
		}()

		/*
			go func() {
				time.Sleep(15 * time.Second)
				panic(fmt.Sprintf("Still waiting %v!", wg))
			}()
		*/

		wg.Wait()
	}
}
