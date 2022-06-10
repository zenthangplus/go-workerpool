package workerpool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFuncJob_GivenAFunction_ShouldInitCorrectlyWithGeneratedId(t *testing.T) {
	job := NewFuncJob(func() {})
	assert.NotNil(t, job.id)
	assert.NotNil(t, job.execFunc)
}

func TestNewFuncJobWithId_GiveAnId_ShouldInitCorrectly(t *testing.T) {
	job1 := NewFuncJobWithId("1", func() {})
	assert.Equal(t, "1", job1.id)
	assert.NotNil(t, job1.execFunc)

	job2 := NewFuncJobWithId("2", func() {})
	assert.Equal(t, "2", job2.id)
	assert.NotNil(t, job2.execFunc)
}
