# Golang Worker Pool

[![run tests](https://github.com/zenthangplus/go-workerpool/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/zenthangplus/go-workerpool/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/zenthangplus/go-workerpool/branch/main/graph/badge.svg)](https://codecov.io/gh/zenthangplus/go-workerpool)

Inspired from Java Thread Pool, Go WorkerPool aims to control heavy Go Routines.

## Installation

The simplest way to install the library is to run:

```shell
go get github.com/zenthangplus/go-workerpool
```

## Example

```go
package main

import (
	"fmt"
	"github.com/zenthangplus/go-workerpool"
)

func main() {
	// Init worker pool with 3 workers to run concurrently.
	pool := workerpool.NewFixedSize(3)

	// Start worker pool
	pool.Start()

	// pool.Submit will block until slot available in Pool queue. 
	// Submit an identifiable job, ID will be generated randomly (using UUID)
	pool.Submit(workerpool.NewIdentifiableJob(func() {
		// Do a heavy job
	}))
	// Use NewCustomIdentifierJob if you don't want ID to be generated randomly
	pool.Submit(workerpool.NewCustomIdentifierJob("custom-id", func() { 
		// Do a heavy job
	}))
	// or Submit a simple function without identifier
	pool.SubmitFunc(func() {// simpler way of: Submit(FuncJob(func() {})) 
		// Do a heavy job
	})

	// pool.SubmitConfidently will submit a job in confident mode, 
	// this function will return ErrPoolFull when Pool queue is full.
	err := pool.SubmitConfidently(workerpool.NewIdentifiableJob(func() {
		// Do a heavy job
	}))
	if err == workerpool.ErrPoolFull {
		fmt.Println("Pool is full")
	}
}
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/zenthangplus/go-workerpool"
)

func main() {
	// Initiate worker pool with fixed size. Eg: 3 workers to run concurrently.
	pool := workerpool.NewFixedSize(3)

	// Or initiate fixed size worker pool with custom options.
	pool = workerpool.NewFixedSize(3,
		// When you want to custom mode
		workerpool.WithMode(workerpool.FixedSize),
		
		// When you want to custom number of workers
		workerpool.WithNumberWorkers(5),
		
		// When you want to customize capacity
		workerpool.WithCapacity(6),
		
		// When you want to custom log function
		workerpool.WithLogFunc(func(msgFormat string, args ...interface{}) {
			fmt.Printf(msgFormat+"\n", args...)
		}),
	)

	// Start worker pool
	pool.Start()
	
	// Init a functional job with ID is generated randomly
	job1 := workerpool.NewIdentifiableJob(func() {})

	// init a functional job with predefined ID
	job2 := workerpool.NewCustomIdentifierJob("test-an-id", func() {})

	// Submit job in normal mode, it will block until pool has available slot.
	pool.Submit(job1)
	
	// or Submit a simple function
	pool.SubmitFunc(func() {})
	
	// Submit in confident mode, it will return ErrPoolFull when pool is full. 
	err := pool.SubmitConfidently(job2)
	if err != nil {
		fmt.Print(err)
	}
}

// CompressDirJob
// You can create a custom Job by implement `Job` interface
type CompressDirJob struct {
	directory string
}

func NewCompressDirJob(directory string) *CompressDirJob {
	return &CompressDirJob{directory: directory}
}

func (c CompressDirJob) Id() string {
	return "directory-" + c.directory
}

func (c CompressDirJob) Exec() {
	// Do compress directory
}
```