package clock

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Simulator_NewTicker(t *testing.T) {
	Convey("对于 Simulator 来说", t, func() {
		s := NewSimulator(time.Now())
		Convey("创建非正时间的 Ticker 会 panic", func() {
			So(func() {
				s.NewTicker(-time.Second)
			}, ShouldPanicWith, "non-positive interval for NewTicker")
		})
		Convey("创建正时间的 Ticker 会返回 *Ticker 类型 ", func() {
			So(s.NewTicker(time.Second), ShouldHaveSameTypeAs, &Ticker2{})
		})
	})
}

func Test_Simulator_Tick(t *testing.T) {
	Convey("对于 Simulator 来说", t, func() {
		s := NewSimulator(time.Now())
		Convey("创建非正时间的 Tick 会返回 nil", func() {
			channel := s.Tick(-time.Second)
			So(channel, ShouldBeNil)
		})
		Convey("创建正时间的 Tick 会返回 <-chan time.Time 类型 ", func() {
			timeChan := make(<-chan time.Time)
			So(s.Tick(time.Second), ShouldHaveSameTypeAs, timeChan)
		})
	})
}

func Test_Simulator_newTicker(t *testing.T) {
	Convey("新建一个 *Ticker", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		interval := time.Second
		t := s.newTicker(interval)
		Convey("t.task 还是活的", func() {
			So(t.task.hasStopped(), ShouldBeFalse)
			Convey("立刻停止的话，t.task 会死掉", func() {
				t.Stop()
				So(t.task.hasStopped(), ShouldBeTrue)
			})
		})
		Convey("一个周期后，能接收到正确的当前时间", func() {
			go func() {
				s.Add(interval)
			}()
			expected := now.Add(interval)
			actual := <-t.C
			So(actual, ShouldEqual, expected)
		})
		Convey("3 个周期后接收一次，会是第一个周期的时间", func() {
			s.Add(3 * interval)
			expected := now.Add(interval)
			actual := <-t.C
			So(actual, ShouldEqual, expected)
			Convey("再过一个周期后，还是能接收到正确的当前时间", func() {
				go func() {
					s.Add(interval)
				}()
				expected := now.Add(4 * interval)
				actual := <-t.C
				So(actual, ShouldEqual, expected)
			})
		})
	})
}
