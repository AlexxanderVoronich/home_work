package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math"
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
		maxErrorsCount := 0

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks concurrency with three allowed errors", func(t *testing.T) {
		tasksCount := 60
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime int64

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				start := time.Now()
				atomic.AddInt32(&runTasksCount, 1)
				val := atomic.LoadInt32(&runTasksCount)
				if val == 3 || val == 13 || val == 27 {
					wrappedErr := fmt.Errorf("wrap error: %w", err)
					return wrappedErr
				}
				for j := 0; j < 1_000_000_000; j++ {
					_ = math.Sqrt(float64(j))
				}
				dur := time.Since(start)
				// fmt.Printf("i=%d :Part duration = %s\n", i, dur.String())
				atomic.AddInt64(&sumTime, dur.Nanoseconds())
				// fmt.Printf("i=%d :Whole time = %d\n", i, sum)
				return nil
			})
		}

		workersCount := 4
		maxErrorsCount := 3

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)
		/*require.Eventually(t, func() bool {
			return runTasksCount == int32(tasksCount)
		}, 5*time.Second, time.Second, "expected %d tasks to be executed, not %d", tasksCount, runTasksCount)
		*/
		sumTimeDuration := time.Duration(atomic.LoadInt64(&sumTime))
		fmt.Printf("Real elapsedTime = %s, sum time = %s\n", elapsedTime.String(), sumTimeDuration.String())
		require.Equal(t, int32(tasksCount), runTasksCount,
			"expected %d tasks to be executed, not %d", tasksCount, runTasksCount)
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTimeDuration), "tasks were run sequentially?")
	})

	t.Run("if the max number of errors is negative, than consider it to be zero", func(t *testing.T) {
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
		maxErrorsCount := -3
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount), "extra tasks were started")
	})
}
