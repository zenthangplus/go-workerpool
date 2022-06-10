package workerpool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithMode(t *testing.T) {
	opt := &Option{Mode: 0}
	WithMode(5)(opt)
	assert.Equal(t, Mode(5), opt.Mode)
}

func TestWithNumberWorkers(t *testing.T) {
	opt := &Option{NumberWorkers: 2}
	WithNumberWorkers(5)(opt)
	assert.Equal(t, 5, opt.NumberWorkers)
}

func TestWithCapacity(t *testing.T) {
	opt := &Option{Capacity: 2}
	WithCapacity(5)(opt)
	assert.Equal(t, 5, opt.Capacity)
}

func TestWithLogFunc(t *testing.T) {
	opt := &Option{LogFunc: nil}
	WithLogFunc(func(msgFormat string, args ...interface{}) {})(opt)
	assert.NotNil(t, opt.LogFunc)
}

func Test_makeDefaultOption(t *testing.T) {
	opt := Option{}
	makeDefaultOption(&opt)
	assert.Equal(t, FixedSize, opt.Mode)
	assert.Equal(t, defaultNumberWorkers, opt.NumberWorkers)
	assert.Equal(t, defaultNumberWorkers*defaultCapacityRatio, opt.Capacity)
	assert.NotNil(t, opt.LogFunc)
}
