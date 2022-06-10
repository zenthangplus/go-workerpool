package workerpool

type Pool struct {
	option        *Option
	jobQueue      chan Job
	workers       []*Worker
	idleWorkers   chan int
	submittedJobs int
	assignedJobs  int
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

func NewFixedSize(numberWorkers int, optionFunc ...OptionFunc) *Pool {
	opt := Option{
		mode:          FixedSize,
		numberWorkers: numberWorkers,
		logFunc:       defaultLogFunc,
	}
	for _, optFunc := range optionFunc {
		optFunc(&opt)
	}
	return New(&opt)
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
	p.option.logFunc("Worker %d is ready for new job", worker.Id())
	p.idleWorkers <- worker.Id()
}

func (p *Pool) beforeSubmit(job Job) {
	p.option.logFunc("Job [%s] is submitted", job.Id())
	p.submittedJobs++
}

// Submit a job.
// This will block until slot available in Pool queue.
func (p *Pool) Submit(job Job) {
	p.beforeSubmit(job)
	p.jobQueue <- job
	p.option.logFunc("Job [%s] is queued", job.Id())
}

// SubmitConfidently submit a job in confidently mode.
// This will return ErrPoolFull when Pool queue is full.
func (p *Pool) SubmitConfidently(job Job) error {
	p.beforeSubmit(job)
	select {
	case p.jobQueue <- job:
		p.option.logFunc("Job [%s] is queued", job.Id())
		return nil
	default:
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
