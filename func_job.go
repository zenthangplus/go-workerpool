package workerpool

import "github.com/google/uuid"

type FuncJob struct {
	id       string
	execFunc func()
}

func NewFuncJob(execFunc func()) *FuncJob {
	return &FuncJob{
		id:       uuid.New().String(),
		execFunc: execFunc,
	}
}

func NewFuncJobWithId(id string, execFunc func()) *FuncJob {
	return &FuncJob{
		id:       id,
		execFunc: execFunc,
	}
}

func (c FuncJob) Id() string {
	return c.id
}

func (c FuncJob) Exec() {
	c.execFunc()
}
