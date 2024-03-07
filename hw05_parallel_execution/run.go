package hw05parallelexecution

import (
	"errors"
	"sync"
)

// ErrErrorsLimitExceeded is error about exceeding the allowed number of errors.
var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

// Task is type of runnable task.
type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving maxErrors errors from tasks.
func Run(tasks []Task, n, maxErrors int) error {
	if maxErrors < 0 {
		maxErrors = 0
	}
	errorsCount := 0
	var wg sync.WaitGroup
	taskCh := make(chan func() error)
	doneWithMaxErrorsCh := make(chan bool, n)
	mutex := sync.RWMutex{}
	isErrorsLimitExceeded := false
	worker := func() {
		defer wg.Done()
		for task := range taskCh {
			err := task()
			if err != nil {
				mutex.Lock()
				errorsCount++
				mutex.Unlock()
			}
			mutex.RLock()
			if errorsCount > maxErrors {
				mutex.RUnlock()
				doneWithMaxErrorsCh <- true
				return
			}
			mutex.RUnlock()
		}
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker()
	}

	for _, task := range tasks {
		select {
		case taskCh <- task:
		case <-doneWithMaxErrorsCh:
			isErrorsLimitExceeded = true
		}
		if isErrorsLimitExceeded {
			break
		}
	}
	close(taskCh)
	wg.Wait()
	if isErrorsLimitExceeded {
		return ErrErrorsLimitExceeded
	}
	return nil
}
