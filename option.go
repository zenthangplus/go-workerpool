package workerpool

import (
	"fmt"
)

type Mode int
type LogFunc func(msgFormat string, args ...interface{})

type Option struct {
	mode          Mode
	capacity      int
	numberWorkers int
	logFunc       LogFunc
}

type OptionFunc func(opt *Option)

func WithMode(mode Mode) OptionFunc {
	return func(opt *Option) {
		opt.mode = mode
	}
}

func WithNumberWorkers(numberWorkers int) OptionFunc {
	return func(opt *Option) {
		opt.numberWorkers = numberWorkers
	}
}

func WithCapacity(capacity int) OptionFunc {
	return func(opt *Option) {
		opt.capacity = capacity
	}
}

func WithLogFunc(logFunc LogFunc) OptionFunc {
	return func(opt *Option) {
		opt.logFunc = logFunc
	}
}

func makeDefaultOption(option *Option) {
	if option.logFunc == nil {
		option.logFunc = func(msgFormat string, args ...interface{}) {
			fmt.Printf(msgFormat, args...)
		}
	}
	if option.mode != FixedSize && option.mode != FlexibleSize {
		option.mode = FixedSize
		option.logFunc("Invalid pool mode, fallback to FixedSize")
	}
	if option.numberWorkers <= 0 {
		option.numberWorkers = defaultNumberWorkers
	}
	if option.capacity <= 0 {
		option.capacity = option.numberWorkers * defaultCapacityRatio
	}
}
