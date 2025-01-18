package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	var resultError error

	chanTasks := make(chan Task, len(tasks))

	for _, task := range tasks {
		chanTasks <- task
	}
	close(chanTasks)

	var errPool int64
	m64 := int64(m)

	for worker := 0; worker < n; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range chanTasks {
				err := task()
				if err != nil {
					atomic.AddInt64(&errPool, 1)
				}
				if atomic.LoadInt64(&errPool) > m64 {
					mu.Lock()
					resultError = ErrErrorsLimitExceeded
					mu.Unlock()
					return
				}
			}
		}()
	}

	wg.Wait()
	return resultError
}
