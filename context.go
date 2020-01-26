package clock

import (
	"context"
	"time"
)

// contextSim 实现了 context.Context 接口
type contextSim struct {
	// 所有与时间无关的方法，直接由此属性组合
	context.Context
	// 与时间相关的方法，通过以下属性改写
	deadline time.Time
	done     chan struct{}
	err      error
}

func (s *Simulator) newContextSim(parent context.Context, deadline time.Time) context.Context {
	ctx := &contextSim{
		Context:  parent,
		done:     make(chan struct{}),
		deadline: deadline,
	}
	t := s.newTimerFunc(deadline, nil)
	go func() {
		// 监控上下文的改变
		select {
		// 本上下文的时间到期了
		case <-t.C:
			ctx.err = context.DeadlineExceeded
		// 父上下文改变了
		case <-parent.Done():
			ctx.err = parent.Err()
			// 停止本上下文的时间监控
			defer t.Stop()
		}
		// 通知子上下文
		close(ctx.done)
	}()
	return ctx
}

func (ctx *contextSim) Deadline() (time.Time, bool) {
	return ctx.deadline, true
}

func (ctx *contextSim) Done() <-chan struct{} {
	return ctx.done
}

func (ctx *contextSim) Err() error {
	select {
	case <-ctx.done:
		return ctx.err
	default:
		return nil
	}
}

type clockKey struct{}

// Set 把 Clock 放入 ctx 中
// 这样的话，*Simulator 可以通过 Context 传递。
func Set(ctx context.Context, c Clock) context.Context {
	return context.WithValue(ctx, clockKey{}, c)
}

// Get 取出 ctx 中的 Clock
// 为了简化从 Context 中取出 Clock 后的相关操作，
// 对相关的操作进行了封装。
func Get(ctx context.Context) Clock {
	if c, ok := ctx.Value(clockKey{}).(Clock); ok {
		return c
	}
	return NewRealClock()
}

// After is a convenience wrapper for Get(ctx).After.
func After(ctx context.Context, d time.Duration) <-chan time.Time {
	return Get(ctx).After(d)
}

// AfterFunc is a convenience wrapper for Get(ctx).AfterFunc.
func AfterFunc(ctx context.Context, d time.Duration, f func()) *Timer {
	return Get(ctx).AfterFunc(d, f)
}

// NewTicker is a convenience wrapper for Get(ctx).NewTicker.
func NewTicker(ctx context.Context, d time.Duration) *Ticker {
	return Get(ctx).NewTicker(d)
}

// NewTimer is a convenience wrapper for Get(ctx).NewTimer.
func NewTimer(ctx context.Context, d time.Duration) *Timer {
	return Get(ctx).NewTimer(d)
}

// Now is a convenience wrapper for Get(ctx).Now.
func Now(ctx context.Context) time.Time {
	return Get(ctx).Now()
}

// Since is a convenience wrapper for Get(ctx).Since.
func Since(ctx context.Context, t time.Time) time.Duration {
	return Get(ctx).Since(t)
}

// Sleep is a convenience wrapper for Get(ctx).Sleep.
func Sleep(ctx context.Context, d time.Duration) {
	Get(ctx).Sleep(d)
}

// Tick is a convenience wrapper for Get(ctx).Tick.
func Tick(ctx context.Context, d time.Duration) <-chan time.Time {
	return Get(ctx).Tick(d)
}

// Until is a convenience wrapper for Get(ctx).Until.
func Until(ctx context.Context, t time.Time) time.Duration {
	return Get(ctx).Until(t)
}

// ContextWithDeadline is a convenience wrapper for Get(ctx).ContextWithDeadline.
func ContextWithDeadline(ctx context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return Get(ctx).ContextWithDeadline(ctx, d)
}

// ContextWithTimeout is a convenience wrapper for Get(ctx).ContextWithTimeout.
func ContextWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return Get(ctx).ContextWithTimeout(ctx, timeout)
}
