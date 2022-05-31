package workerpool

import "time"

type Worker struct {
	id           int
	workedJobs   int
	runningJobId string
	busiedAt     *time.Time
	idledAt      *time.Time
}

func NewWorker(id int) *Worker {
	now := time.Now().UTC()
	return &Worker{
		id:      id,
		idledAt: &now,
	}
}

func (w *Worker) Run(job Job) {
	w.beforeJob(job)
	job.Exec()
	w.afterJob()
}

func (w *Worker) beforeJob(job Job) {
	now := time.Now().UTC()
	w.busiedAt = &now
	w.idledAt = nil
	w.workedJobs++
	w.runningJobId = job.Id()
}

func (w *Worker) afterJob() {
	w.busiedAt = nil
	now := time.Now().UTC()
	w.idledAt = &now
	w.runningJobId = ""
}

func (w Worker) Id() int {
	return w.id
}

func (w Worker) BusiedAt() *time.Time {
	return w.busiedAt
}

func (w Worker) BusiedDuration() time.Duration {
	if w.busiedAt == nil {
		return time.Duration(0)
	}
	return time.Now().UTC().Sub(*w.busiedAt)
}

func (w Worker) IdledAt() *time.Time {
	return w.idledAt
}

func (w Worker) IdledDuration() time.Duration {
	if w.idledAt == nil {
		return time.Duration(0)
	}
	return time.Now().UTC().Sub(*w.idledAt)
}

func (w Worker) RunningJobId() string {
	return w.runningJobId
}

func (w Worker) WorkedJobs() int {
	return w.workedJobs
}
