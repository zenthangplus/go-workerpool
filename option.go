package workerpool

type Mode int
type LogFunc func(msgFormat string, args ...interface{})

type Option struct {
	Mode          Mode
	Capacity      int
	NumberWorkers int
	LogFunc       LogFunc
}

type OptionFunc func(opt *Option)

func WithMode(mode Mode) OptionFunc {
	return func(opt *Option) {
		opt.Mode = mode
	}
}

func WithNumberWorkers(numberWorkers int) OptionFunc {
	return func(opt *Option) {
		opt.NumberWorkers = numberWorkers
	}
}

func WithCapacity(capacity int) OptionFunc {
	return func(opt *Option) {
		opt.Capacity = capacity
	}
}

func WithLogFunc(logFunc LogFunc) OptionFunc {
	return func(opt *Option) {
		opt.LogFunc = logFunc
	}
}

func makeDefaultOption(option *Option) {
	if option.LogFunc == nil {
		option.LogFunc = defaultLogFunc
	}
	if option.Mode != FixedSize && option.Mode != FlexibleSize {
		option.Mode = FixedSize
		option.LogFunc("Invalid pool mode, fallback to FixedSize")
	}
	if option.NumberWorkers <= 0 {
		option.NumberWorkers = defaultNumberWorkers
	}
	if option.Capacity <= 0 {
		option.Capacity = option.NumberWorkers * defaultCapacityRatio
	}
}
