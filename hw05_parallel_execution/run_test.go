package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks with zero m", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32
		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}
		workersCount := 10
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		mu := sync.Mutex{}

		var runTasksCount int32
		var currentTaskCount int32
		var tempInt int32
		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&currentTaskCount, 1)
				defer atomic.AddInt32(&currentTaskCount, -1)

				mu.Lock()
				for tempInt = 0; tempInt < 100000; {
					tempInt++
				}
				mu.Unlock()

				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1
		var err error
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = Run(tasks, workersCount, maxErrorsCount)
		}()

		require.Eventually(t, func() bool {
			return atomic.LoadInt32(&currentTaskCount) >= int32(3)
		}, time.Second*5, time.Millisecond, "tasks were run sequentially?")
		wg.Wait()
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.NoError(t, err)
	})
}
