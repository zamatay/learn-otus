package hw05parallelexecution

import (
	"errors"
	"log"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error
type tackChanel = chan Task

var mu sync.Mutex
var wg = sync.WaitGroup{}

func worker(tc tackChanel, errorCount *int, i int) {
	log.Printf("run worker %d", i)
	isDie := false
	for t := range tc {
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
	tCh := make(tackChanel, len(tasks))
	wg.Add(n)

	runWorker(&tCh, &errorCount, n)

	for _, t := range tasks {
		tCh <- t
	}

	close(tCh)
	wg.Wait()
	if errorCount < 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func runWorker(ch *tackChanel, errorCount *int, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(i int) {
			defer wg.Done()
			worker(*ch, errorCount, i)
		}(i)
	}
}
