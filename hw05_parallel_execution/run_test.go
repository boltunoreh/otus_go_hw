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

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestCustomRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	tests := []struct {
		name              string
		successTasksCount int
		errorTasksCount   int
		workersCount      int
		maxErrorsCount    int
	}{
		{name: "tasks len < n", successTasksCount: 0, errorTasksCount: 10, workersCount: 20, maxErrorsCount: 5},
		{name: "tasks len == n", successTasksCount: 0, errorTasksCount: 10, workersCount: 10, maxErrorsCount: 10},
		{name: "tasks len > n", successTasksCount: 0, errorTasksCount: 10, workersCount: 5, maxErrorsCount: 10},
		{name: "error count > task count", successTasksCount: 0, errorTasksCount: 10, workersCount: 20, maxErrorsCount: 100},
		{name: "zero m", successTasksCount: 0, errorTasksCount: 10, workersCount: 20, maxErrorsCount: 0},
		{name: "negative m", successTasksCount: 0, errorTasksCount: 10, workersCount: 20, maxErrorsCount: -10},

		{name: "tasks len < n", successTasksCount: 10, errorTasksCount: 10, workersCount: 20, maxErrorsCount: 10},
		{name: "tasks len == n", successTasksCount: 10, errorTasksCount: 10, workersCount: 10, maxErrorsCount: 10},
		{name: "tasks len > n", successTasksCount: 10, errorTasksCount: 10, workersCount: 5, maxErrorsCount: 10},
		{name: "error count > task count", successTasksCount: 10, errorTasksCount: 10, workersCount: 20, maxErrorsCount: 100},
		{name: "zero m", successTasksCount: 10, errorTasksCount: 10, workersCount: 20, maxErrorsCount: 0},
		{name: "negative m", successTasksCount: 10, errorTasksCount: 10, workersCount: 20, maxErrorsCount: -10},

		{
			name:              "random n m",
			successTasksCount: rand.Intn(100),
			errorTasksCount:   rand.Intn(100),
			workersCount:      rand.Intn(100),
			maxErrorsCount:    rand.Intn(100),
		},
		{
			name:              "random n m",
			successTasksCount: rand.Intn(100),
			errorTasksCount:   rand.Intn(100),
			workersCount:      rand.Intn(100),
			maxErrorsCount:    rand.Intn(100),
		},
		{
			name:              "random n m",
			successTasksCount: rand.Intn(100),
			errorTasksCount:   rand.Intn(100),
			workersCount:      rand.Intn(100),
			maxErrorsCount:    rand.Intn(100),
		},
		{
			name:              "random n m",
			successTasksCount: rand.Intn(100),
			errorTasksCount:   rand.Intn(100),
			workersCount:      rand.Intn(100),
			maxErrorsCount:    rand.Intn(100),
		},
		{
			name:              "random n m",
			successTasksCount: rand.Intn(100),
			errorTasksCount:   rand.Intn(100),
			workersCount:      rand.Intn(100),
			maxErrorsCount:    rand.Intn(100),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tasks := make([]Task, 0, test.successTasksCount)

			var runTasksCount int32
			var runErrorTasksCount int32
			var runSuccessTasksCount int32

			for i := 0; i < test.errorTasksCount; i++ {
				err := fmt.Errorf("error from task %d", i)
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runErrorTasksCount, 1)
					atomic.AddInt32(&runTasksCount, 1)
					return err
				})
			}

			for i := 0; i < test.successTasksCount; i++ {
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runSuccessTasksCount, 1)
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}

			err := Run(tasks, test.workersCount, test.maxErrorsCount)

			var maxTasks int32
			if test.maxErrorsCount >= test.errorTasksCount {
				maxTasks = int32(test.workersCount + test.errorTasksCount)
			} else {
				maxTasks = int32(test.workersCount + test.maxErrorsCount)
			}

			var expectedError error
			if test.maxErrorsCount <= test.errorTasksCount {
				expectedError = ErrErrorsLimitExceeded
			}

			require.Truef(t, errors.Is(err, expectedError), "actual err - %v", err)
			require.LessOrEqual(t, runErrorTasksCount, maxTasks, "extra tasks were started")
		})
	}
}
