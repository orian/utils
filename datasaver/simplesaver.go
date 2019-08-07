package datasaver

import (
	"github.com/sirupsen/logrus"

	"time"
)

type Saver interface {
	Save(r interface{})
}

type RefreshPolicy interface {
	Refresh(encoder Encoder) bool
}

type ConstantIntervalRefresh struct {
	Interval   time.Duration
	lastChange time.Time
}

func NewConstantIntervalRefresh(interval time.Duration) *ConstantIntervalRefresh {
	return &ConstantIntervalRefresh{interval, time.Now()}
}

func (c *ConstantIntervalRefresh) Refresh(encoder Encoder) bool {
	t := time.Now()
	if t.Sub(c.lastChange) > c.Interval {
		c.lastChange = t
	}
	return encoder.Created().Before(c.lastChange)
}

type ExactTimeOfDayRefresh struct {
	nextChange time.Time
	prevChange time.Time
}

func NewExactUtcTimeOfDayRefresh(hour, minute int) *ExactTimeOfDayRefresh {
	r := &ExactTimeOfDayRefresh{}
	t := time.Now().UTC()
	if h, m, _ := t.Clock(); h > hour || h == hour && m >= minute {
		t = t.AddDate(0, 0, 1)
	}
	y, mo, d := t.Date()
	r.nextChange = time.Date(y, mo, d, hour, minute, 0, 0, time.UTC)
	r.prevChange = r.nextChange.AddDate(0, 0, -1)
	return r
}

func (c *ExactTimeOfDayRefresh) Refresh(encoder Encoder) bool {
	t := time.Now()
	if t.After(c.nextChange) {
		c.prevChange = c.nextChange
		c.nextChange = c.nextChange.AddDate(0, 0, 1)
	}
	return encoder.Created().Before(c.prevChange)
}

type RefreshCond func(encoder Encoder) bool

func (r RefreshCond) Refresh(encoder Encoder) bool {
	return r(encoder)
}

type BoolOrRefresh struct {
	Policies []RefreshPolicy
}

func (b BoolOrRefresh) Refresh(encoder Encoder) bool {
	for _, p := range b.Policies {
		if p.Refresh(encoder) {
			return true
		}
	}
	return false
}

type SimpleSaver struct {
	r        *RotateFile
	policy   RefreshPolicy
	ticker   *time.Ticker
	q        chan interface{}
	finished chan bool
}

func NewSimpleSaver(r *RotateFile, policy RefreshPolicy, queueSize int) *SimpleSaver {
	return &SimpleSaver{r, policy, time.NewTicker(time.Second),
		make(chan interface{}, queueSize), make(chan bool, 1)}
}

func (s *SimpleSaver) Save(r interface{}) {
	s.q <- r
}

func (s *SimpleSaver) GracefulStop() {
	s.ticker.Stop()
	close(s.q)
	<-s.finished
}

func (s *SimpleSaver) save(v interface{}) error {
	return s.r.Encode(v)
}

func (s *SimpleSaver) Start() {
	var err error
	for {
		select {
		case <-s.ticker.C:
			if s.policy.Refresh(s.r.enc) {
				if err = s.r.Refresh(); err != nil {
					logrus.Errorf("refreshing: %s", err)
				}
			}
		case r, ok := <-s.q:
			if !ok {
				goto end
			}
			if err := s.save(r); err != nil {
				logrus.Errorf("writing record: %s", err)
			}
		}
	}
end:
	if err = s.r.Close(); err != nil {
		logrus.Errorf("closing file: %s", err)
	}
	s.finished <- true
}
