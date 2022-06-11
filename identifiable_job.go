package workerpool

import "github.com/google/uuid"

type IdentifiableJob struct {
	id       string
	execFunc func()
}

func NewIdentifiableJob(execFunc func()) *IdentifiableJob {
	return &IdentifiableJob{
		id:       uuid.New().String(),
		execFunc: execFunc,
	}
}

func NewCustomIdentifierJob(id string, execFunc func()) *IdentifiableJob {
	return &IdentifiableJob{
		id:       id,
		execFunc: execFunc,
	}
}

func (c IdentifiableJob) Id() string {
	return c.id
}

func (c IdentifiableJob) Exec() {
	c.execFunc()
}
