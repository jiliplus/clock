package clock

import (
	"time"
)

// Ticker represents a time.Ticker.
type Ticker struct {
	C    <-chan time.Time
	Stop func()
	// 当 ticker != nil 的时候, Ticker 代表了 real clock
	ticker *time.Ticker
	*task
}

// NewTicker returns a new Ticker containing a channel that will send the
// current time with a period specified by the duration d.
func (s *Simulator) NewTicker(d time.Duration) *Ticker {
	s.Lock()
	defer s.Unlock()
	if d <= 0 {
		panic("non-positive interval for NewTicker")
	}
	return s.newTicker(d)
}

// Tick is a convenience wrapper for NewTicker providing access to the ticking
// channel only.
func (s *Simulator) Tick(d time.Duration) <-chan time.Time {
	s.Lock()
	defer s.Unlock()
	if d <= 0 {
		return nil
	}
	return s.newTicker(d).C
}

func (s *Simulator) newTicker(d time.Duration) *Ticker {
	c := make(chan time.Time, 1)
	run := func(t *task) *task {
		// time.Tick.C 的发送逻辑是
		//   能发送，就发送
		// 不能发送，就抛弃
		// 所以，c 带有缓存，免得全部都丢弃了
		select {
		case c <- s.now:
		default:
		}
		t.deadline = t.deadline.Add(d)
		return t
	}
	t := &Ticker{
		C:    c,
		task: newTask(s.now.Add(d), run),
	}
	t.Stop = func() {
		if t.ticker != nil {
			t.ticker.Stop()
			return
		}
		s.Lock()
		s.heap.remove(t.task)
		s.Unlock()
	}
	s.accept(t.task)
	return t
}
