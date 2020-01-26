package clock

import "time"

// Timer 替代 time.Timer.
type Timer struct {
	C <-chan time.Time
	// 当 timer != nil 的时候, Timer 代表了 real clock
	timer *time.Timer
	*task

	// Stop prevents the Timer from firing.
	// It returns true if the call stops the timer, false if the timer has already expired or been stopped.
	Stop func() bool
	// Reset changes the timer to expire after duration d.
	// It returns true if the timer had been active, false if the timer had expired or been stopped.
	//
	// A negative or zero duration fires the timer immediately.
	Reset func(d time.Duration) bool
}

// Sleep pauses the current goroutine for at least the duration d.
//
// A negative or zero duration causes Sleep to return immediately.
func (s *Simulator) Sleep(d time.Duration) {
	<-s.After(d)
}

// After waits for the duration to elapse and then sends the current time on
// the returned channel.
//
// A negative or zero duration fires the underlying timer immediately.
func (s *Simulator) After(d time.Duration) <-chan time.Time {
	return s.NewTimer(d).C
}

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine.
// It returns a Timer that can be used to cancel the call using its Stop method.
//
// A negative or zero duration fires the timer immediately.
func (s *Simulator) AfterFunc(d time.Duration, f func()) *Timer {
	s.Lock()
	defer s.Unlock()
	return s.newTimerFunc(s.now.Add(d), f)
}

// NewTimer creates a new Timer that will send the current time on its channel
// after at least duration d.
//
// A negative or zero duration fires the timer immediately.
func (s *Simulator) NewTimer(d time.Duration) *Timer {
	s.Lock()
	defer s.Unlock()
	return s.newTimerFunc(s.now.Add(d), nil)
}

//
func (s *Simulator) newTimerFunc(deadline time.Time, afterFunc func()) *Timer {
	c := make(chan time.Time, 1)
	runTask := func(t *task) *task {
		if afterFunc != nil {
			// NOTICE: AfterFunc 创建的 *Timer 不会发送 current time
			go afterFunc()
		} else {
			// Timer 的发送逻辑和 Tick 的不一样。
			// 必须发送到位
			c <- s.now
		}
		return nil
	}
	timer := &Timer{
		C:    c,
		task: newTask(deadline, runTask),
	}
	s.accept(timer.task)
	timer.Stop = func() bool {
		s.Lock()
		defer s.Unlock()
		isActive := !timer.task.hasStopped()
		s.heap.remove(timer.task)
		return isActive
	}
	timer.Reset = func(d time.Duration) bool {
		s.Lock()
		defer s.Unlock()
		isActive := !timer.hasStopped()
		s.heap.remove(timer.task)
		timer.deadline = s.now.Add(d)
		s.accept(timer.task)
		return isActive
	}
	return timer
}

func doNothing() {}
