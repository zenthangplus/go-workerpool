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
		jobQueue:    make(chan Job, option.Capacity),
		workers:     make([]*Worker, 0),
		idleWorkers: make(chan int, option.NumberWorkers),
		option:      option,
	}
}

func NewFixedSize(numberWorkers int, optionFunc ...OptionFunc) *Pool {
	opt := Option{
		Mode:          FixedSize,
		NumberWorkers: numberWorkers,
		LogFunc:       defaultLogFunc,
	}
	for _, optFunc := range optionFunc {
		optFunc(&opt)
	}
	return New(&opt)
}

func (p *Pool) Start() {
	for i := 1; i <= p.option.NumberWorkers; i++ {
		worker := NewWorker(i)
		p.option.LogFunc("Worker %d initialed", worker.id)
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
	p.option.LogFunc("Worker %d got a job [%s]", worker.Id(), job.Id())
	p.assignedJobs++
	worker.Run(job)
	p.option.LogFunc("Worker %d is ready for new job", worker.Id())
	p.idleWorkers <- worker.Id()
}

func (p *Pool) beforeSubmit(job Job) {
	p.option.LogFunc("Job [%s] is submitted", job.Id())
	p.submittedJobs++
}

// Submit a job.
// This will block until slot available in Pool queue.
func (p *Pool) Submit(job Job) {
	p.beforeSubmit(job)
	p.jobQueue <- job
	p.option.LogFunc("Job [%s] is queued", job.Id())
}

// SubmitConfidently submit a job in confidently mode.
// This will return ErrPoolFull when Pool queue is full.
func (p *Pool) SubmitConfidently(job Job) error {
	p.beforeSubmit(job)
	select {
	case p.jobQueue <- job:
		p.option.LogFunc("Job [%s] is queued", job.Id())
		return nil
	default:
		p.option.LogFunc("Job [%s] is rejected due by pool full", job.Id())
		return ErrPoolFull
	}
}

func (p Pool) Capacity() int {
	return p.option.Capacity
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
