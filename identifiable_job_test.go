package workerpool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifiableJob_GivenAFunction_ShouldInitCorrectlyWithGeneratedId(t *testing.T) {
	job := NewIdentifiableJob(func() {})
	assert.NotNil(t, job.id)
	assert.NotNil(t, job.execFunc)
}

func TestNewCustomIdentifierJob_GiveAnId_ShouldInitCorrectly(t *testing.T) {
	job1 := NewCustomIdentifierJob("1", func() {})
	assert.Equal(t, "1", job1.id)
	assert.NotNil(t, job1.execFunc)

	job2 := NewCustomIdentifierJob("2", func() {})
	assert.Equal(t, "2", job2.id)
	assert.NotNil(t, job2.execFunc)
}
