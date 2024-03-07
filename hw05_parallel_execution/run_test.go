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

	t.Run("tasks concurrency with on 1 and 4 goroutines", func(t *testing.T) {
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
				for j := 0; j < 200_000_000; j++ {
					_ = math.Sqrt(float64(j))
				}
				dur := time.Since(start)
				atomic.AddInt64(&sumTime, dur.Nanoseconds())
				return nil
			})
		}

		/*go assert.Eventually(t, func() bool {
			val := atomic.LoadInt32(&runTasksCount)
			t.Logf("Out %d\n", val)
			return val >= int32(tasksCount)
		}, 10*time.Second, 500*time.Millisecond, "expected to be executed")*/

		// start first experiment with 1 goroutine
		workersCount := 1
		maxErrorsCount := 3

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime1 := time.Since(start)
		require.NoError(t, err)
		sumTimeDuration1 := time.Duration(atomic.LoadInt64(&sumTime))
		t.Logf("Real elapsedTime on 1g = %s, sum time = %s\n", elapsedTime1.String(), sumTimeDuration1.String())
		require.Equal(t, int32(tasksCount), runTasksCount,
			"expected %d tasks to be executed, not %d", tasksCount, runTasksCount)

		// start second experiment with 4 goroutines
		workersCount = 4
		maxErrorsCount = 3
		runTasksCount = 0
		sumTime = 0

		start = time.Now()
		err2 := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime2 := time.Since(start)
		require.NoError(t, err2)
		sumTimeDuration2 := time.Duration(atomic.LoadInt64(&sumTime))
		t.Logf("Real elapsedTime on 4g = %s, sum time = %s\n", elapsedTime2.String(), sumTimeDuration1.String())
		require.Equal(t, int32(tasksCount), runTasksCount,
			"expected %d tasks to be executed, not %d", tasksCount, runTasksCount)

		// check both experiments
		require.LessOrEqual(t, int64(elapsedTime2), int64(elapsedTime1),
			"compare tasks duration on 1g and 4g")
		require.LessOrEqual(t, int64(elapsedTime2), int64(sumTimeDuration2),
			"compare real and accumulate tasks duration for 4g")
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
