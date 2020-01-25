package clock

import (
	"fmt"
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
		second := time.Second
		var mutex sync.Mutex
		timer := c.AfterFunc(second, func() {
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
				}, ShouldNotPanic)
			})
		})
		Convey("立刻重置为两倍原来的时候的话", func() {
			isActive := timer.Reset(second * 2)
			Convey("timer 仍然是活的", func() {
				So(isActive, ShouldBeTrue)
			})
			Convey("func 未被执行，所以，可以上锁", func() {
				So(func() {
					mutex.Lock()
				}, ShouldNotPanic)
			})
			Convey("等待两倍的时间后", func() {
				time.Sleep(second * 2)
				Convey("func 已执行，此时可以解锁", func() {
					So(func() {
						mutex.Unlock()
					}, ShouldNotPanic)
				})
			})
		})
		SkipConvey("等待一段时间后", func() {
			time.Sleep(second)
			Convey("func 已执行，此时可以解锁", func() {
				So(func() {
					mutex.Unlock()
				}, ShouldNotPanic)
			})
		})
	})
}

func Test_realClock_Now(t *testing.T) {
	Convey("新建一个 realClock", t, func() {
		c := NewRealClock()
		Convey("realClock.Now() 和 time.Now() 应该返回差不多的时间", func() {
			delta := time.Microsecond
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
		fmt.Println(c.Now())
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

func Test_realClock_x(t *testing.T) {
	Convey("新建一个 realClock", t, func() {
		c := NewRealClock()
		fmt.Println(c.Now())
	})
}
