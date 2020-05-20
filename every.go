package clock

import (
	"time"
)

// nextDayTime 是为了 EveryDay 函数准备的
// 是为了每天在 hour:minute:second 提供 tick
// NOTICE: hour 是 24 小时制的小时
func nextDayTime(now time.Time, hour, minute, second int) time.Time {
	yyyy, mm, dd := now.Year(), now.Month(), now.Day()
	loc := time.Now().Location()
	next := time.Date(yyyy, mm, dd, hour, minute, second, 0, loc)
	if now.Before(next) {
		return next
	}
	return next.Add(24 * time.Hour)
}

// EveryDay returns a channel which
// output a time.Time by your setting
func (s *Simulator) EveryDay(hour, minute, second int) <-chan time.Time {
	s.Lock()
	defer s.Unlock()
	c := make(chan time.Time, 1)
	run := func(t *task) *task {
		c <- s.now // 必须要送达
		t.deadline = t.deadline.Add(24 * time.Hour)
		return t
	}
	task := newTask(nextDayTime(s.now, hour, minute, second), run)
	s.accept(task)
	return c
}
