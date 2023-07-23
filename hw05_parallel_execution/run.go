package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	die := make(chan struct{})
	workerEnd := make(chan struct{})
	mu := sync.Mutex{}
	errorCount := 0

	w := func(t Task, i int) {
		for {
			select {
			case <-die:
				return
			default:
				err := t()
				if err != nil {
					select {
					case _, ok := <-die:
						if !ok {
							return
						}
					default:
						mu.Lock()
						errorCount++
						mu.Unlock()
						if m == 0 || errorCount == m {
							close(die)
						}
					}
				} else {
					workerEnd <- struct{}{}
					return
				}
			}
		}
	}
	index := 0
	for i, t := range tasks {
		go w(t, i)
		index++
		if index == n {
			select {
			case <-workerEnd:
				index--
			case <-die:
				return ErrErrorsLimitExceeded
			}
		}
	}
	for i := 1; i < n; i++ {
		<-workerEnd
	}
	return nil
}
