// Package clock 可以模拟 time 和 context 标准库的部分行为。
//
// All methods are safe for concurrent use.
package clock

import (
	"context"
	"time"
)

// Clock 是对 time 和 context 标准库的部分 API 进行的封装
// 就是需要在时间轴上进行跳转那些部分。
type Clock interface {
	After(d time.Duration) <-chan time.Time
	AfterFunc(d time.Duration, f func()) *Timer
	NewTicker(d time.Duration) *Ticker
	NewTimer(d time.Duration) *Timer
	Now() time.Time
	Since(t time.Time) time.Duration
	Sleep(d time.Duration)
	Tick(d time.Duration) <-chan time.Time
	Until(t time.Time) time.Duration

	// ContextWithDeadline 与 context.ContextWithDeadline 具有相同的功能
	// 只是基于 Clock 的时间线
	ContextWithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc)
	// ContextWithTimeout 是 ContextWithDeadline(parent, Now(parent).Add(timeout)).
	ContextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
}
