package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	var errorCount int
	mu := sync.Mutex{}
	tasksCh := make(chan Task, len(tasks))

	for _, task := range tasks {
		tasksCh <- task
	}
	close(tasksCh)

	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			for task := range tasksCh {
				taskError := task()
				if taskError != nil {
					mu.Lock()
					errorCount++
					if errorCount >= m {
						mu.Unlock()
						break
					}
					mu.Unlock()
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()

	if errorCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
