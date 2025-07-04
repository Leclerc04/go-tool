package concurrent

import (
	"strings"
	"sync"
)

// MultiError can contain multiple erro
type MultiError []error

func (m MultiError) Error() string {
	var msgs []string
	for _, e := range m {
		msgs = append(msgs, e.Error())
	}
	return strings.Join(msgs, "\n")
}

// Pool helps to limit concurrency.
type Pool struct {
	queue chan func() error

	mu  sync.Mutex
	err MultiError
}

// NewPool creates the pool and returns the destructor.
func NewPool(numWorkers, maxQueueSize int) (*Pool, func() error) {
	p := &Pool{
		queue: make(chan func() error, maxQueueSize),
	}
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for f := range p.queue {
				err := f()
				if err != nil {
					func() {
						p.mu.Lock()
						defer p.mu.Unlock()
						p.err = append(p.err, err)
					}()
				}
			}
		}()
	}
	return p, func() error {
		close(p.queue)
		wg.Wait()
		if len(p.err) == 0 {
			// Return the error interface nil.
			return nil
		}
		return p.err
	}
}

// Run dispatces the given closure to the queue.
func (p *Pool) Run(f func() error) {
	p.queue <- f
}
