package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type (
	Task       func() error
	tackChanel = chan Task
)

var (
	mu sync.Mutex
	wg sync.WaitGroup
	tc tackChanel
)

func worker(errorCount *int) {
	isDie := false
	for t := range tc {
		if t == nil {
			return
		}
		err := t()
		mu.Lock()
		if *errorCount <= 0 {
			isDie = true
		}
		if err != nil {
			*errorCount--
		}
		mu.Unlock()
		if isDie {
			return
		}
	}
}

func Run(tasks []Task, n int, m int) error {
	errorCount := m
	tc = make(tackChanel, len(tasks))
	wg.Add(n)

	runWorker(&errorCount, n)

	for _, t := range tasks {
		tc <- t
	}

	close(tc)
	wg.Wait()
	if errorCount <= 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func runWorker(errorCount *int, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			worker(errorCount)
		}()
	}
}
