package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type (
	Task       func() error
	taskChanel = chan Task
)

type Worker struct {
	mu sync.Mutex
	wg sync.WaitGroup
	tc taskChanel
}

func (r *Worker) worker(errorCount *int) {
	isDie := false
	for t := range r.tc {
		if t == nil {
			return
		}
		err := t()
		r.mu.Lock()
		if *errorCount <= 0 {
			isDie = true
		}
		if err != nil {
			*errorCount--
		}
		r.mu.Unlock()
		if isDie {
			return
		}
	}
}

func (r *Worker) Run(tasks []Task, n int, m int) error {
	errorCount := m
	r.tc = make(taskChanel, len(tasks))
	r.wg.Add(n)

	r.runWorker(&errorCount, n)

	for _, t := range tasks {
		r.tc <- t
	}

	close(r.tc)
	r.wg.Wait()
	if errorCount <= 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func (r *Worker) runWorker(errorCount *int, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func() {
			defer r.wg.Done()
			r.worker(errorCount)
		}()
	}
}
