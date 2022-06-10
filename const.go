package workerpool

import (
	"errors"
	"fmt"
	"time"
)

const (
	FixedSize Mode = iota
	FlexibleSize
)

const (
	defaultCapacityRatio = 20
	defaultNumberWorkers = 10
)

var defaultLogFunc = func(msgFormat string, args ...interface{}) {
	fmt.Printf(time.Now().Format("2006-01-02T15:04:05.999Z")+": "+msgFormat+"\n", args...)
}

var ErrPoolFull = errors.New("pool is full")
