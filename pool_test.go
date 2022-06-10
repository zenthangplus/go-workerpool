package workerpool

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPool_GivenNumberWorkers_WhenNewFixedSize_ShouldInitFieldsCorrectly(t *testing.T) {
	numWorkers := 3
	pool := NewFixedSize(numWorkers)
	assert.Equal(t, FixedSize, pool.option.mode)
	assert.Equal(t, numWorkers, pool.option.numberWorkers)
	assert.Equal(t, numWorkers*defaultCapacityRatio, pool.option.capacity)
	assert.Equal(t, numWorkers*defaultCapacityRatio, pool.Capacity())
	assert.NotNil(t, pool.option.logFunc)
	assert.Equal(t, pool.option.capacity, cap(pool.jobQueue))
	assert.Len(t, pool.Workers(), 0)
	assert.Equal(t, 0, pool.SubmittedJobs())
	assert.Equal(t, 0, pool.AssignedJobs())
}

func TestPool_GivenAnOption_WhenNew_ShouldInitFieldsCorrectly(t *testing.T) {
	opt := Option{
		mode:          FlexibleSize,
		capacity:      1000,
		numberWorkers: 4,
		logFunc:       func(msgFormat string, args ...interface{}) {},
	}
	pool := New(&opt)
	assert.Equal(t, opt.mode, pool.option.mode)
	assert.Equal(t, opt.numberWorkers, pool.option.numberWorkers)
	assert.Equal(t, opt.capacity, pool.option.capacity)
	assert.Equal(t, opt.capacity, pool.Capacity())
	assert.NotNil(t, pool.option.logFunc)
	assert.Equal(t, pool.option.capacity, cap(pool.jobQueue))
	assert.Len(t, pool.Workers(), 0)
	assert.Equal(t, 0, pool.SubmittedJobs())
	assert.Equal(t, 0, pool.AssignedJobs())
}

func TestPool_GivenNegativeOption_WhenNew_ShouldInitWithFallbackValues(t *testing.T) {
	opt := Option{
		mode:          FlexibleSize,
		capacity:      -100,
		numberWorkers: -4,
		logFunc:       func(msgFormat string, args ...interface{}) {},
	}
	pool := New(&opt)
	assert.Equal(t, opt.mode, pool.option.mode)
	assert.Equal(t, defaultNumberWorkers, pool.option.numberWorkers)
	assert.Equal(t, defaultNumberWorkers*defaultCapacityRatio, pool.option.capacity)
	assert.Equal(t, defaultNumberWorkers*defaultCapacityRatio, pool.Capacity())
	assert.NotNil(t, pool.option.logFunc)
	assert.Equal(t, pool.option.capacity, cap(pool.jobQueue))
	assert.Len(t, pool.Workers(), 0)
	assert.Equal(t, 0, pool.SubmittedJobs())
	assert.Equal(t, 0, pool.AssignedJobs())
}

func startDummyPoolFixedSize(t *testing.T, numberWorkers int, options ...OptionFunc) *Pool {
	pool := NewFixedSize(numberWorkers, options...)
	pool.Start()
	assert.Equal(t, FixedSize, pool.option.mode)
	assert.Equal(t, numberWorkers, pool.option.numberWorkers)
	assert.Equal(t, numberWorkers*defaultCapacityRatio, pool.option.capacity)
	assert.Equal(t, numberWorkers*defaultCapacityRatio, pool.Capacity())
	assert.NotNil(t, pool.option.logFunc)
	assert.Len(t, pool.Workers(), numberWorkers)
	return pool
}

func TestPool_GivenAPoolFixedSize_WhenSubmitConfidentlyJob_ShouldInitWorkerCorrectly(t *testing.T) {
	numberWorkers := 2
	capacity := 2
	pool := NewFixedSize(numberWorkers, WithCapacity(capacity))
	pool.Start()
	assert.Equal(t, FixedSize, pool.option.mode)
	assert.Equal(t, numberWorkers, pool.option.numberWorkers)
	assert.Equal(t, capacity, pool.option.capacity)
	assert.Equal(t, capacity, pool.Capacity())
	assert.NotNil(t, pool.option.logFunc)
	assert.Len(t, pool.Workers(), numberWorkers)

	job1 := 0
	assert.NoError(t, pool.SubmitConfidently(NewSimpleJobWithId("1", func() {
		job1 = 1
		pool.option.logFunc("Job 1 is finished")
	})))
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 1, job1)
	assert.Equal(t, 1, pool.SubmittedJobs())
	assert.Equal(t, 1, pool.AssignedJobs())
	job2 := 0
	assert.NoError(t, pool.SubmitConfidently(NewSimpleJobWithId("2", func() {
		job2++
		time.Sleep(200 * time.Millisecond)
		pool.option.logFunc("Job 2 is finished")
		job2++
	})))
	job3 := 0
	assert.NoError(t, pool.SubmitConfidently(NewSimpleJobWithId("3", func() {
		job3++
		time.Sleep(200 * time.Millisecond)
		pool.option.logFunc("Job 3 is finished")
		job3++
	})))

	pool.option.logFunc("Wait for workers receive job 2, 3")
	time.Sleep(10 * time.Millisecond)

	// Job 5, 6 will be queued (job 4 will be in Idle holding point)
	job4 := 0
	assert.NoError(t, pool.SubmitConfidently(NewSimpleJobWithId("4", func() {
		job4++
		time.Sleep(500 * time.Millisecond)
		pool.option.logFunc("Job 4 is finished")
		job4++
	})))
	job5 := 0
	assert.NoError(t, pool.SubmitConfidently(NewSimpleJobWithId("5", func() {
		job5++
		time.Sleep(500 * time.Millisecond)
		pool.option.logFunc("Job 5 is finished")
		job5++
	})))
	job6 := 0
	assert.NoError(t, pool.SubmitConfidently(NewSimpleJobWithId("6", func() {
		job6++
		time.Sleep(500 * time.Millisecond)
		pool.option.logFunc("Job 6 is finished")
		job6++
	})))

	// Job 7 will be rejected due by queue full
	job7 := 0
	assert.ErrorIs(t, pool.SubmitConfidently(NewSimpleJobWithId("7", func() {
		job7++
		time.Sleep(500 * time.Millisecond)
		pool.option.logFunc("Job 7 is finished")
		job7++
	})), ErrPoolFull)

	pool.option.logFunc("Submitted all jobs")
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
