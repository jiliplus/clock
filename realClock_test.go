package clock

import (
	"context"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_realClock_After(t *testing.T) {
	Convey("新建一个 realClock", t, func() {
		c := NewRealClock()
		Convey("after 应该在规定的时间内，触发", func() {
			now := c.Now()
			dur := time.Millisecond * 500
			delta := time.Millisecond * 5
			after := <-c.After(dur)
			So(after, ShouldHappenWithin, dur+delta, now)
		})
	})
}

func Test_realClock_AfterFunc(t *testing.T) {
	Convey("利用 realClock.AfterFunc 生成 *Timer", t, func() {
		c := NewRealClock()
		secs := time.Second
		var mutex sync.Mutex
		timer := c.AfterFunc(secs, func() {
			mutex.Lock()
		})
		Convey("立刻停止 timer 的话", func() {
			isActive := timer.Stop()
			Convey("timer 仍然是活的", func() {
				So(isActive, ShouldBeTrue)
			})
			Convey("func 未被执行，所以，可以上锁", func() {
				So(func() {
					mutex.Lock()
					defer mutex.Unlock()
				}, ShouldNotPanic)
			})
		})
		Convey("立刻重置为两倍原来的时候的话", func() {
			isActive := timer.Reset(secs * 2)
			Convey("timer 仍然是活的", func() {
				So(isActive, ShouldBeTrue)
			})
			Convey("func 未被执行，所以，可以上锁", func() {
				So(func() {
					mutex.Lock()
					defer mutex.Unlock()
				}, ShouldNotPanic)
			})
			Convey("超过 timer 的时限后", func() {
				// * 110 / 100 是为了确保 timer 触发
				time.Sleep(secs * 2 * 110 / 100)
				Convey("func 已执行，此时可以解锁", func() {
					So(func() {
						mutex.Unlock()
					}, ShouldNotPanic)
				})
			})
		})
		Convey("超过 timer 的时限后", func() {
			// * 110 / 100 是为了确保 timer 触发
			time.Sleep(secs * 110 / 100)
			Convey("func 已执行，此时可以解锁", func() {
				So(func() {
					mutex.Unlock()
				}, ShouldNotPanic)
			})
		})
	})
}

func Test_realClock_NewTicker(t *testing.T) {
	Convey("利用 realClock 创建 ticker", t, func() {
		c := NewRealClock()
		t := c.NewTicker(time.Second)
		Convey("应该是对 time.Ticker 的封装", func() {
			So(t.Stop, ShouldEqual, t.ticker.Stop)
		})
	})
}

func Test_realClock_NewTimer(t *testing.T) {
	Convey("利用 realClock 创建 timer", t, func() {
		c := NewRealClock()
		t := c.NewTimer(time.Second)
		Convey("应该是对 time.Timer 的封装", func() {
			So(t.C, ShouldEqual, t.timer.C)
			So(t.Reset, ShouldEqual, t.timer.Reset)
			So(t.Stop, ShouldEqual, t.timer.Stop)
		})
	})
}

func Test_realClock_Now(t *testing.T) {
	Convey("新建一个 realClock", t, func() {
		c := NewRealClock()
		Convey("realClock.Now() 和 time.Now() 应该返回差不多的时间", func() {
			delta := 10 * time.Microsecond
			So(c.Now(), ShouldHappenWithin, delta, time.Now())
		})
	})
}

func Test_realClock_Since(t *testing.T) {
	Convey("新建一个 realClock", t, func() {
		c := NewRealClock()
		expected := time.Second
		oneSecondBefore := time.Now().Add(-expected)
		actual := c.Since(oneSecondBefore)
		Convey("realClock.Since 应该返回 1 秒钟", func() {
			So(actual, ShouldAlmostEqual, expected, time.Millisecond)
		})
	})
}

func Test_realClock_Sleep(t *testing.T) {
	Convey("新建一个 realClock", t, func() {
		c := NewRealClock()
		dur := time.Millisecond * 10
		delta := time.Millisecond
		start := time.Now()
		c.Sleep(dur)
		end := time.Now()
		Convey("起止时间，应该差不多就相差了一个 dur", func() {
			So(end, ShouldHappenWithin, dur+delta, start)
		})
	})
}

func Test_realClock_Tick(t *testing.T) {
	Convey("新建一个 realClock", t, func() {
		c := NewRealClock()
		dur := time.Second
		delta := dur / 10
		start := time.Now()
		end := <-c.Tick(dur)
		Convey("起止时间，应该差不多就相差了一个 dur", func() {
			So(end, ShouldHappenWithin, dur+delta, start)
		})
	})
}

func Test_realClock_Until(t *testing.T) {
	Convey("新建一个 realClock", t, func() {
		c := NewRealClock()
		expected := time.Second
		oneSecondAfter := time.Now().Add(expected)
		actual := c.Until(oneSecondAfter)
		Convey("realClock.Since 应该返回 1 秒钟", func() {
			So(actual, ShouldAlmostEqual, expected, time.Millisecond)
		})
	})
}

func Test_realClock_ContextWithDeadline(t *testing.T) {
	Convey("利用 realClock 创建 context", t, func() {
		c := NewRealClock()
		dur := time.Millisecond * 200
		delta := time.Millisecond
		deadline := time.Now().Add(dur)
		ctx, _ := c.ContextWithDeadline(context.Background(), deadline)
		<-ctx.Done()
		endTime := time.Now()
		Convey("应该大约在预定的时间内完成", func() {
			So(endTime, ShouldHappenWithin, delta, deadline)
		})
	})
}

func Test_realClock_ContextWithTimeout(t *testing.T) {
	Convey("利用 realClock 创建 context", t, func() {
		c := NewRealClock()
		timeout := time.Millisecond * 200
		delta := time.Millisecond
		deadline := time.Now().Add(timeout)
		ctx, _ := c.ContextWithTimeout(context.Background(), timeout)
		<-ctx.Done()
		endTime := time.Now()
		Convey("应该大约在预定的时间内完成", func() {
			So(endTime, ShouldHappenWithin, delta, deadline)
		})
	})
}
