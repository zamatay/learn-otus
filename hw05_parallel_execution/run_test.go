package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func getTask(runTasksCount *int32, taskSleep time.Duration, i int, withError bool) func() error {
	return func() error {
		var err error
		if withError {
			err = fmt.Errorf("error from task %d", i)
		} else {
			err = nil
		}
		time.Sleep(taskSleep)
		atomic.AddInt32(runTasksCount, 1)
		return err
	}
}

//nolint:funlen,gocognit
func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			tasks = append(tasks, getTask(&runTasksCount, taskSleep, i, true))
		}

		workersCount := 10
		maxErrorsCount := 23
		w := Worker{}
		err := w.Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v, (%d)", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, getTask(&runTasksCount, taskSleep, i, false))
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		w := Worker{}
		err := w.Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks without m = 0 errors", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				err := fmt.Errorf("error from task %d", i)
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 5
		maxErrorsCount := 0

		w := Worker{}
		err := w.Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "Not errors")
	})

	t.Run("tasks without m = 0 errors", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, getTask(&runTasksCount, taskSleep, i, true))
		}

		workersCount := 5
		maxErrorsCount := 0

		start := time.Now()
		w := Worker{}
		err := w.Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("nil task", func(t *testing.T) {
		tasksCount := 2
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			if i == 1 {
				tasks = append(tasks, nil)
			} else {
				tasks = append(tasks, getTask(&runTasksCount, taskSleep, i, false))
			}
		}

		workersCount := 2
		maxErrorsCount := 1
		w := Worker{}
		err := w.Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err, "actual err - %v", err)
	})

	t.Run("count error and max error is equal", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			if i == 0 || i == 1 {
				tasks = append(tasks, getTask(&runTasksCount, taskSleep, i, true))
			} else {
				tasks = append(tasks, getTask(&runTasksCount, taskSleep, i, false))
			}
		}

		workersCount := 2
		maxErrorsCount := 2

		w := Worker{}
		err := w.Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "OK")
	})

	t.Run("count error and max error is equal", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			if i == 0 || i == 1 {
				tasks = append(tasks, getTask(&runTasksCount, taskSleep, i, true))
			} else {
				tasks = append(tasks, getTask(&runTasksCount, taskSleep, i, false))
			}
		}

		for i := 0; i < 5; i++ {
			go func() {
				workersCount := 2
				maxErrorsCount := 2
				w := Worker{}
				w.Run(tasks, workersCount, maxErrorsCount)
			}()
		}
	})
}
