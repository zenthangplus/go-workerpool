package workerpool

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPool_GivenNumberWorkers_WhenNewFixedSize_ShouldInitFieldsCorrectly(t *testing.T) {
	numWorkers := 3
	pool := NewFixedSize(numWorkers)
	assert.Equal(t, FixedSize, pool.option.Mode)
	assert.Equal(t, numWorkers, pool.option.NumberWorkers)
	assert.Equal(t, numWorkers*defaultCapacityRatio, pool.option.Capacity)
	assert.Equal(t, numWorkers*defaultCapacityRatio, pool.Capacity())
	assert.NotNil(t, pool.option.LogFunc)
	assert.Equal(t, pool.option.Capacity, cap(pool.jobQueue))
	assert.Len(t, pool.Workers(), 0)
	assert.Equal(t, 0, pool.SubmittedJobs())
	assert.Equal(t, 0, pool.AssignedJobs())
}

func TestPool_GivenAnOption_WhenNew_ShouldInitFieldsCorrectly(t *testing.T) {
	opt := Option{
		Mode:          FlexibleSize,
		Capacity:      1000,
		NumberWorkers: 4,
		LogFunc:       func(msgFormat string, args ...interface{}) {},
	}
	pool := New(&opt)
	assert.Equal(t, opt.Mode, pool.option.Mode)
	assert.Equal(t, opt.NumberWorkers, pool.option.NumberWorkers)
	assert.Equal(t, opt.Capacity, pool.option.Capacity)
	assert.Equal(t, opt.Capacity, pool.Capacity())
	assert.NotNil(t, pool.option.LogFunc)
	assert.Equal(t, pool.option.Capacity, cap(pool.jobQueue))
	assert.Len(t, pool.Workers(), 0)
	assert.Equal(t, 0, pool.SubmittedJobs())
	assert.Equal(t, 0, pool.AssignedJobs())
}

func TestPool_GivenNegativeOption_WhenNew_ShouldInitWithFallbackValues(t *testing.T) {
	opt := Option{
		Mode:          FlexibleSize,
		Capacity:      -100,
		NumberWorkers: -4,
		LogFunc:       func(msgFormat string, args ...interface{}) {},
	}
	pool := New(&opt)
	assert.Equal(t, opt.Mode, pool.option.Mode)
	assert.Equal(t, defaultNumberWorkers, pool.option.NumberWorkers)
	assert.Equal(t, defaultNumberWorkers*defaultCapacityRatio, pool.option.Capacity)
	assert.Equal(t, defaultNumberWorkers*defaultCapacityRatio, pool.Capacity())
	assert.NotNil(t, pool.option.LogFunc)
	assert.Equal(t, pool.option.Capacity, cap(pool.jobQueue))
	assert.Len(t, pool.Workers(), 0)
	assert.Equal(t, 0, pool.SubmittedJobs())
	assert.Equal(t, 0, pool.AssignedJobs())
}

func startDummyPoolFixedSize(t *testing.T, numberWorkers int, capacity int) *Pool {
	pool := NewFixedSize(numberWorkers, WithCapacity(capacity))
	pool.Start()
	assert.Equal(t, FixedSize, pool.option.Mode)
	assert.Equal(t, numberWorkers, pool.option.NumberWorkers)
	assert.Equal(t, capacity, pool.option.Capacity)
	assert.Equal(t, capacity, pool.Capacity())
	assert.NotNil(t, pool.option.LogFunc)
	assert.Len(t, pool.Workers(), numberWorkers)
	return pool
}

func TestPool_GivenAPoolFixedSize_WhenSubmitJob_ShouldRunCorrectly(t *testing.T) {
	pool := startDummyPoolFixedSize(t, 2, 2)
	job1 := 0
	go pool.Submit(NewFuncJobWithId("1", func() {
		job1 = 1
		pool.option.LogFunc("Job 1 is finished")
	}))
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 1, job1)
	assert.Equal(t, 1, pool.SubmittedJobs())
	assert.Equal(t, 1, pool.AssignedJobs())
	assert.Len(t, pool.jobQueue, 0)

	// Job 2, 3 will be assigned to workers (1, 2)
	job2 := 0
	pool.Submit(NewFuncJobWithId("2", func() {
		job2++
		time.Sleep(200 * time.Millisecond)
		pool.option.LogFunc("Job 2 is finished")
		job2++
	}))
	job3 := 0
	pool.Submit(NewFuncJobWithId("3", func() {
		job3++
		time.Sleep(200 * time.Millisecond)
		pool.option.LogFunc("Job 3 is finished")
		job3++
	}))
	// Job 4, 5 will be added to job queue
	// When all worker are busy, job 4 or 5 will be wait in Idle holding point (See Pool.Start function)
	job4 := 0
	pool.Submit(NewFuncJobWithId("4", func() {
		job4++
		time.Sleep(200 * time.Millisecond)
		pool.option.LogFunc("Job 4 is finished")
		job4++
	}))
	job5 := 0
	pool.Submit(NewFuncJobWithId("5", func() {
		job5++
		time.Sleep(200 * time.Millisecond)
		pool.option.LogFunc("Job 5 is finished")
		job5++
	}))
	// Job 6 will be added to the queue even though capacity=2 due by Idle holding point (See Pool.Start function)
	job6 := 0
	pool.Submit(NewFuncJobWithId("6", func() {
		job6++
		time.Sleep(200 * time.Millisecond)
		pool.option.LogFunc("Job 6 is finished")
		job6++
	}))

	// Job 7 will be hang due by job queue's capacity is full now
	job7 := 0
	job7InQueue := false
	go func() {
		pool.Submit(NewFuncJobWithId("7", func() {
			job7++
			time.Sleep(200 * time.Millisecond)
			pool.option.LogFunc("Job 7 is finished")
			job7++
		}))
		job7InQueue = true
	}()
	pool.option.LogFunc("Submitted all jobs")

	time.Sleep(50 * time.Millisecond)
	// Job 2, 3 is processing
	// Job 4 in Idle holding point, job 5, 6 in queue
	// Job 7 is hanged
	assert.Equal(t, 1, job2)
	assert.Equal(t, 1, job3)
	assert.Equal(t, 0, job4)
	assert.Equal(t, 0, job5)
	assert.Equal(t, 0, job6)
	assert.Equal(t, 0, job7)
	assert.False(t, job7InQueue)
	assert.Len(t, pool.jobQueue, 2)
	assert.Equal(t, 7, pool.SubmittedJobs())
	assert.Equal(t, 3, pool.AssignedJobs())

	time.Sleep(200 * time.Millisecond)
	// Job 2, 3 is finished
	// Job 4, 5 is processing
	// Job 6 in Idle holding point
	// Job 7 in queue
	assert.Equal(t, 2, job2)
	assert.Equal(t, 2, job3)
	assert.Equal(t, 1, job4)
	assert.Equal(t, 1, job5)
	assert.Equal(t, 0, job6)
	assert.Equal(t, 0, job7)
	assert.True(t, job7InQueue)
	assert.Len(t, pool.jobQueue, 1)
	assert.Equal(t, 7, pool.SubmittedJobs())
	assert.Equal(t, 5, pool.AssignedJobs())

	time.Sleep(200 * time.Millisecond)
	// Job 2, 3, 4, 5 is finished
	// Job 6, 7 is processing
	assert.Equal(t, 2, job2)
	assert.Equal(t, 2, job3)
	assert.Equal(t, 2, job4)
	assert.Equal(t, 2, job5)
	assert.Equal(t, 1, job6)
	assert.Equal(t, 1, job7)
	assert.Len(t, pool.jobQueue, 0)
	assert.Equal(t, 7, pool.SubmittedJobs())
	assert.Equal(t, 7, pool.AssignedJobs())

	time.Sleep(200 * time.Millisecond)
	// All jobs are finished
	assert.Equal(t, 2, job2)
	assert.Equal(t, 2, job3)
	assert.Equal(t, 2, job4)
	assert.Equal(t, 2, job5)
	assert.Equal(t, 2, job6)
	assert.Equal(t, 2, job7)
	assert.Len(t, pool.jobQueue, 0)
	assert.Equal(t, 7, pool.SubmittedJobs())
	assert.Equal(t, 7, pool.AssignedJobs())
}

func TestPool_GivenAPoolFixedSize_WhenSubmitConfidentlyJob_ShouldRunCorrectly(t *testing.T) {
	pool := startDummyPoolFixedSize(t, 2, 2)
	job1 := 0
	assert.NoError(t, pool.SubmitConfidently(NewFuncJobWithId("1", func() {
		job1 = 1
		pool.option.LogFunc("Job 1 is finished")
	})))
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 1, job1)
	assert.Equal(t, 1, pool.SubmittedJobs())
	assert.Equal(t, 1, pool.AssignedJobs())

	job2 := 0
	assert.NoError(t, pool.SubmitConfidently(NewFuncJobWithId("2", func() {
		job2++
		time.Sleep(200 * time.Millisecond)
		pool.option.LogFunc("Job 2 is finished")
		job2++
	})))
	job3 := 0
	assert.NoError(t, pool.SubmitConfidently(NewFuncJobWithId("3", func() {
		job3++
		time.Sleep(200 * time.Millisecond)
		pool.option.LogFunc("Job 3 is finished")
		job3++
	})))

	pool.option.LogFunc("Wait for workers receive job 2, 3")
	time.Sleep(10 * time.Millisecond)

	// Job 5, 6 will be queued (job 4 will be in Idle holding point)
	job4 := 0
	assert.NoError(t, pool.SubmitConfidently(NewFuncJobWithId("4", func() {
		job4++
		time.Sleep(500 * time.Millisecond)
		pool.option.LogFunc("Job 4 is finished")
		job4++
	})))
	job5 := 0
	assert.NoError(t, pool.SubmitConfidently(NewFuncJobWithId("5", func() {
		job5++
		time.Sleep(500 * time.Millisecond)
		pool.option.LogFunc("Job 5 is finished")
		job5++
	})))
	job6 := 0
	assert.NoError(t, pool.SubmitConfidently(NewFuncJobWithId("6", func() {
		job6++
		time.Sleep(500 * time.Millisecond)
		pool.option.LogFunc("Job 6 is finished")
		job6++
	})))

	// Job 7 will be rejected due by queue full
	job7 := 0
	assert.ErrorIs(t, pool.SubmitConfidently(NewFuncJobWithId("7", func() {
		job7++
		time.Sleep(500 * time.Millisecond)
		pool.option.LogFunc("Job 7 is finished")
		job7++
	})), ErrPoolFull)

	pool.option.LogFunc("Submitted all jobs")
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, job2)
	assert.Equal(t, 1, job3)
	assert.Equal(t, 0, job4)
	assert.Equal(t, 0, job5)
	assert.Equal(t, 0, job6)
	assert.Equal(t, 0, job7)
	assert.Equal(t, 7, pool.SubmittedJobs())
	assert.Equal(t, 3, pool.AssignedJobs())

	time.Sleep(150 * time.Millisecond)
	assert.Equal(t, 2, job2)
	assert.Equal(t, 2, job3)
	assert.Equal(t, 1, job4)
	assert.Equal(t, 1, job5)
	assert.Equal(t, 0, job6)
	assert.Equal(t, 0, job7)
	assert.Equal(t, 7, pool.SubmittedJobs())
	assert.Equal(t, 5, pool.AssignedJobs())

	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 2, job2)
	assert.Equal(t, 2, job3)
	assert.Equal(t, 2, job4)
	assert.Equal(t, 2, job5)
	assert.Equal(t, 1, job6)
	assert.Equal(t, 0, job7)
	assert.Equal(t, 7, pool.SubmittedJobs())
	assert.Equal(t, 6, pool.AssignedJobs())

	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 2, job2)
	assert.Equal(t, 2, job3)
	assert.Equal(t, 2, job4)
	assert.Equal(t, 2, job5)
	assert.Equal(t, 2, job6)
	assert.Equal(t, 0, job7)
	assert.Equal(t, 7, pool.SubmittedJobs())
	assert.Equal(t, 6, pool.AssignedJobs())
}
