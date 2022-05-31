package workerpool

import "github.com/google/uuid"

type SimpleJob struct {
	id       string
	execFunc func()
}

func NewSimpleJob(execFunc func()) *SimpleJob {
	return &SimpleJob{
		id:       uuid.New().String(),
		execFunc: execFunc,
	}
}

func NewSimpleJobWithId(id string, execFunc func()) *SimpleJob {
	return &SimpleJob{
		id:       id,
		execFunc: execFunc,
	}
}

func (c SimpleJob) Id() string {
	return c.id
}

func (c SimpleJob) Exec() {
	c.execFunc()
}
