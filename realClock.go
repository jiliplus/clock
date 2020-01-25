package clock

import (
	"context"
	"time"
)

type realClock struct{}

// NewRealClock 返回标准库中真实时间的时钟。
// 并实现了 Clock 接口
func NewRealClock() Clock {
	return realClock{}
}

func (realClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (realClock) AfterFunc(d time.Duration, f func()) *Timer {
	t := time.AfterFunc(d, f)
	return &Timer{
		// 为了与 time.AfterFunc 保持一致性
		// AfterFunc 的 C 是 nil
		Stop:  t.Stop,
		Reset: t.Reset,
		timer: t,
	}
}

func (realClock) NewTicker(d time.Duration) *Ticker {
	t := time.NewTicker(d)
	return &Ticker{
		C:      t.C,
		Stop:   t.Stop,
		ticker: t,
	}
}

func (realClock) NewTimer(d time.Duration) *Timer {
	t := time.NewTimer(d)
	return &Timer{
		C:     t.C,
		Reset: t.Reset,
		Stop:  t.Stop,
		timer: t,
	}
}

func (realClock) Now() time.Time {
	return time.Now()
}

func (realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (realClock) Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

func (realClock) Until(t time.Time) time.Duration {
	return time.Until(t)
}

func (realClock) ContextWithDeadline(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, deadline)
}

func (realClock) ContextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}
