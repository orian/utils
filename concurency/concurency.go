package concurency

import (
	"github.com/robfig/cron"

	"sync/atomic"
)

type SingletonJob struct {
	f         cron.Job
	isRunning int32
}

func NewSingletonJob(j cron.Job) *SingletonJob {
	return &SingletonJob{j, 0}
}

func (s *SingletonJob) Run() {
	if atomic.CompareAndSwapInt32(&s.isRunning, 0, 1) {
		defer atomic.CompareAndSwapInt32(&s.isRunning, 1, 0)
		s.f.Run()
	}
}
