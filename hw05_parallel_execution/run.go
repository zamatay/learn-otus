package hw05parallelexecution

import (
	"errors"
	"log"
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
		log.Printf("gorutine begin work %v\n", i)
		for {
			select {
			case <-die:
				log.Printf("die %v\n", i)
				return
			default:
				err := t()
				if err != nil {
					log.Printf("error %v\n", i)
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
							log.Printf("error count is m = %v\n", errorCount)
							close(die)
						}
					}
				} else {
					log.Printf("end %v\n", i)
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
	log.Printf("exit from loop \n")
	for i := 1; i < n; i++ {
		<-workerEnd
	}
	return nil
}
