package ptime

import "time"

type SimpleTimer struct {
	t0 time.Time
}

func (s *SimpleTimer) Duration() time.Duration {
	return time.Now().Sub(s.t0)
}

func (s *SimpleTimer) Reset() time.Duration {
	t := time.Now()
	t, s.t0 = s.t0, t
	return s.t0.Sub(t)
}

func NewTimer() *SimpleTimer {
	return &SimpleTimer{time.Now()}
}
