package workerpool

import (
	"time"
)

type Pool struct {
	option        *Option
	jobQueue      chan Job
	workers       []*Worker
	idleWorkers   chan int
	hungAt        *time.Time
	hungTimes     int
	submittedJobs int
	assignedJobs  int
}

func NewFixedSize(numberWorkers int, optionFunc ...OptionFunc) *Pool {
	opt := Option{
		mode:          FixedSize,
		numberWorkers: numberWorkers,
		capacity:      numberWorkers * defaultCapacityRatio,
		logFunc:       defaultLogFunc,
	}
	for _, optFunc := range optionFunc {
		optFunc(&opt)
	}
	return New(&opt)
}

func New(option *Option) *Pool {
	makeDefaultOption(option)
	return &Pool{
		jobQueue:    make(chan Job, option.capacity),
		workers:     make([]*Worker, 0),
		idleWorkers: make(chan int, option.numberWorkers),
		option:      option,
	}
}

func (p *Pool) Start() {
	for i := 1; i <= p.option.numberWorkers; i++ {
		worker := NewWorker(i)
		p.option.logFunc("Worker %d initialed", worker.id)
		p.workers = append(p.workers, worker)
		p.idleWorkers <- i
	}
	go func() {
		for job := range p.jobQueue {
			// Idle holding point
			idleWorkerId := <-p.idleWorkers
			worker := p.workers[idleWorkerId-1]
			go p.dispatch(job, worker)
		}
	}()
}

func (p *Pool) dispatch(job Job, worker *Worker) {
	p.option.logFunc("Worker %d got a job [%s]", worker.Id(), job.Id())
	p.assignedJobs++
	worker.Run(job)
	p.idleWorkers <- worker.Id()
	p.option.logFunc("Worker %d is ready for new job", worker.Id())
}

func (p *Pool) Submit(job Job) error {
	p.option.logFunc("Job [%s] is submitted", job.Id())
	p.submittedJobs++
	select {
	case p.jobQueue <- job:
		p.option.logFunc("Job [%s] is queued", job.Id())
		p.hungAt = nil
		return nil
	default:
		now := time.Now().UTC()
		p.hungAt = &now
		p.hungTimes++
		p.option.logFunc("Job [%s] is rejected due by pool full", job.Id())
		return ErrPoolFull
	}
}

func (p Pool) Capacity() int {
	return p.option.capacity
}

func (p Pool) Workers() []*Worker {
	return p.workers
}

func (p Pool) SubmittedJobs() int {
	return p.submittedJobs
}

func (p Pool) AssignedJobs() int {
	return p.assignedJobs
}

func (p Pool) HungAt() *time.Time {
	return p.hungAt
}

func (p Pool) HungDuration() time.Duration {
	if p.hungAt == nil {
		return time.Duration(0)
	}
	return time.Now().UTC().Sub(*p.hungAt)
}

func (p Pool) HungTimes() int {
	return p.hungTimes
}
