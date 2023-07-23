package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error
type endTask = struct{}

func Run(tasks []Task, n, m int) error {
	die := make(chan struct{})
	workerEnd := make(chan endTask)
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
						if m == 0 || errorCount == m {
							close(die)
						}
						mu.Unlock()
					}
				} else {
					workerEnd <- endTask{}
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
