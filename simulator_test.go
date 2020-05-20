package clock

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Simulator_Add(t *testing.T) {
	Convey("新建模拟器 s", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		Convey("使用 Add 给 s 添加负时间", func() {
			expected := s.now
			d := -time.Second
			actual := s.Add(d)
			Convey("s 还是原来时间", func() {
				So(actual, ShouldEqual, expected)
			})
		})
		Convey("使用 Add 给 s 添加 0 时间", func() {
			expected := s.now
			actual := s.Add(0)
			Convey("不会改变 s 的时间", func() {
				So(actual, ShouldEqual, expected)
			})
		})
		Convey("使用 Add 给 s 添加正时间", func() {
			d := time.Second
			actual := s.Add(d)
			expected := now.Add(d)
			Convey("会改变 s 的时间", func() {
				So(actual, ShouldEqual, expected)
			})
		})
	})
}

func Test_Simulator_AddOrPanic(t *testing.T) {
	Convey("新建模拟器 s", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		Convey("使用 AddOrPanic 给 s 添加负时间会 panic", func() {
			So(func() {
				s.AddOrPanic(-time.Second)
			}, ShouldPanicWith, timeReversal)
		})
		Convey("使用 AddOrPanic 给 s 添加 0 时间", func() {
			actual := s.AddOrPanic(0)
			Convey("不会改变 s 的时间", func() {
				So(actual, ShouldEqual, now)
			})
		})
		Convey("使用 AddOrPanic 给 s 添加正时间", func() {
			d := time.Second
			actual := s.AddOrPanic(d)
			expected := now.Add(d)
			Convey("会改变 s 的时间", func() {
				So(actual, ShouldEqual, expected)
			})
		})
	})
}

func Test_Simulator_Set(t *testing.T) {
	Convey("新建模拟器 s", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		Convey("使用 Set 把 s 设置为过去的时间", func() {
			passedTime := s.now.Add(-time.Second)
			actual := s.Set(passedTime)
			Convey("s 还是原来时间", func() {
				So(s.now, ShouldEqual, now)
				So(actual, ShouldEqual, 0)
			})
		})
		Convey("使用 Set 把 s 设置为当前的时间", func() {
			actual := s.Set(s.now)
			Convey("不会改变 s 的时间", func() {
				So(s.now, ShouldEqual, now)
				So(actual, ShouldEqual, 0)
			})
		})
		Convey("使用 Set 把 s 设置为以后的时间", func() {
			d := time.Second
			expectTime := now.Add(d)
			actualDur := s.Set(expectTime)
			Convey("会改变 s 的时间", func() {
				So(actualDur, ShouldEqual, d)
				So(s.now, ShouldEqual, expectTime)
			})
		})
	})
}

func Test_Simulator_SetOrPanic(t *testing.T) {
	Convey("新建模拟器 s", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		Convey("使用 SetOrPanic 把 s 设置为过去的时间，会 panic", func() {
			passedTime := s.now.Add(-time.Second)
			So(func() {
				s.SetOrPanic(passedTime)
			}, ShouldPanicWith, timeReversal)
		})
		Convey("使用 SetOrPanic 把 s 设置为当前的时间", func() {
			actualDur := s.SetOrPanic(s.now)
			Convey("不会改变 s 的时间", func() {
				So(s.now, ShouldEqual, now)
				So(actualDur, ShouldEqual, 0)
			})
		})
		Convey("使用 SetOrPanic 把 s 设置为以后的时间", func() {
			d := time.Second
			expectTime := now.Add(d)
			actualDur := s.SetOrPanic(expectTime)
			Convey("会改变 s 的时间", func() {
				So(actualDur, ShouldEqual, d)
				So(s.now, ShouldEqual, expectTime)
			})
		})
	})
}

func Test_Simulator_Since(t *testing.T) {
	Convey("新建模拟器 s", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		Convey("Since 时间段的起点", func() {
			expectDur := time.Second
			startTime := s.now.Add(-expectDur)
			actualDur := s.Since(startTime)
			Convey("会得到正确的距离", func() {
				So(actualDur, ShouldEqual, expectDur)
			})
		})
	})
}

func Test_Simulator_Until(t *testing.T) {
	Convey("新建模拟器 s", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		Convey("Until 时间段的终点", func() {
			expectDur := time.Second
			startTime := s.now.Add(expectDur)
			actualDur := s.Until(startTime)
			Convey("会得到正确的距离", func() {
				So(actualDur, ShouldEqual, expectDur)
			})
		})
	})
}

func Test_Simulator_Move(t *testing.T) {
	Convey("新建模拟器 s", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		Convey("让没有 task 的 s Move 一下", func() {
			expectTime, expectDur := s.now, time.Duration(0)
			actualTime, actualDur := s.Move()
			Convey("s 不会发生改变", func() {
				So(actualTime, ShouldEqual, expectTime)
				So(actualDur, ShouldEqual, expectDur)
			})
		})
		Convey("给 s 添加 task", func() {
			expectDur := time.Second
			expectTime := s.now.Add(expectDur)
			runTask := func(ts *task) *task { return nil }
			ts := newTask(expectTime, runTask)
			s.accept(ts)
			Convey("让 s Move 一下，会发生改变", func() {
				actualTime, actualDur := s.Move()
				So(actualTime, ShouldEqual, expectTime)
				So(actualDur, ShouldEqual, expectDur)
			})
		})
	})
}

func Test_Simulator_set_timerStyle(t *testing.T) {
	Convey("新建模拟器 s", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		num := 10
		actualOrder := make([]time.Time, 0, num)
		expectOrder := make([]time.Time, num)
		runTask := func(ts *task) *task {
			actualOrder = append(actualOrder, s.now)
			return nil
		}
		for i := num; i > 0; i-- {
			deadline := now.Add(time.Duration(i) * time.Second)
			ts := newTask(deadline, runTask)
			s.accept(ts)
			expectOrder[i-1] = deadline
		}
		Convey("s.heap 的长度应该等于 count", func() {
			So(len(*(s.heap)), ShouldEqual, num)
		})
		Convey("改变 s 的当前时间", func() {
			expectDur := time.Second * time.Duration(num)
			expectTime := now.Add(expectDur)
			s.Lock()
			actualTime, actualDur := s.set(expectTime)
			s.Unlock()
			Convey("s 被改变，并按照预定的顺序执行", func() {
				So(actualTime, ShouldEqual, expectTime)
				So(actualDur, ShouldEqual, expectDur)
				So(actualOrder, ShouldResemble, expectOrder)
			})
		})
	})
}

func Test_Simulator_set_tickerStyle(t *testing.T) {
	Convey("新建模拟器 s", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		num := 10
		actualOrder := make([]time.Time, 0, num)
		runTask := func(ts *task) *task {
			actualOrder = append(actualOrder, s.now)
			ts.deadline = ts.deadline.Add(time.Second)
			return ts
		}
		deadline := now.Add(time.Second)
		ts := newTask(deadline, runTask)
		s.accept(ts)
		Convey("s.heap 的长度应该等于 1", func() {
			So(len(*(s.heap)), ShouldEqual, 1)
		})
		expectOrder := make([]time.Time, num)
		for i := 0; i < num; i++ {
			deadline := now.Add(time.Duration(i+1) * time.Second)
			expectOrder[i] = deadline
		}
		Convey("改变 s 的当前时间", func() {
			expectDur := time.Second * time.Duration(num)
			expectTime := now.Add(expectDur)
			s.Lock()
			actualTime, actualDur := s.set(expectTime)
			s.Unlock()
			Convey("s 被改变，并按照预定的顺序执行", func() {
				So(actualTime, ShouldEqual, expectTime)
				So(actualDur, ShouldEqual, expectDur)
				So(actualOrder, ShouldResemble, expectOrder)
			})
		})
	})
}

func Test_Simulator_accept(t *testing.T) {
	Convey("新建模拟器 s", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		Convey("s.heap 的长度应该为 0", func() {
			So(len(*(s.heap)), ShouldEqual, 0)
		})
		Convey("往 s 中放入 nil task", func() {
			s.accept(nil)
			Convey("s.heap 的长度还是 0", func() {
				So(len(*(s.heap)), ShouldEqual, 0)
			})
		})
		isRunned := false
		ts := &task{}
		ts.runFunc = func(tk *task) *task {
			isRunned = true
			return nil
		}
		Convey("往 s 中放入过期的 task", func() {
			passedTime := now.Add(-1 * time.Minute)
			ts.deadline = passedTime
			s.accept(ts)
			Convey("s.heap 的长度是 1", func() {
				So(len(*(s.heap)), ShouldEqual, 1)
			})
			Convey("任务不会被执行", func() {
				So(isRunned, ShouldBeFalse)
			})
		})
		Convey("往 s 中放入未来的 task", func() {
			future := now.Add(1 * time.Minute)
			ts.deadline = future
			s.accept(ts)
			Convey("s.heap 的长度变成 1", func() {
				So(len(*(s.heap)), ShouldEqual, 1)
			})
			Convey("任务不会被执行", func() {
				So(isRunned, ShouldBeFalse)
			})
		})
	})

}

func Test_Simulator_setNowTo(t *testing.T) {
	Convey("假设存在模拟器 s", t, func() {
		now := time.Now()
		s := &Simulator{
			now: now,
		}
		Convey("如果想把 s.now 设置成过去的时间点", func() {
			last := s.now
			s.setNowTo(s.now.Add(-time.Second))
			Convey("s.now 还是会等于原来的值", func() {
				So(s.now, ShouldEqual, last)
			})
		})
		Convey("如果想把 s.now 设置成以后的时间点", func() {
			last := s.now
			now := s.now.Add(time.Second)
			s.setNowTo(now)
			Convey("s.now 会被设置成新值", func() {
				So(last.Before(s.now), ShouldBeTrue)
				So(s.now, ShouldEqual, now)
			})
		})
	})
}

func Test_Simulator_ContextWithDeadline(t *testing.T) {
	// 把 cancel 赋值到 _ 是因为
	// 它来至于 contest.WithCancel，我 100% 信任标准库
	Convey("对于 Simulator 来说", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		ctx := context.Background()
		deadline := now.Add(time.Second * 1)
		Convey("放入 ctx", func() {
			child, _ := s.ContextWithDeadline(ctx, deadline)
			Convey("child 的 deadline 符合预期", func() {
				actual, ok := child.Deadline()
				So(actual, ShouldEqual, deadline)
				So(ok, ShouldBeTrue)
			})
		})
	})
}

func Test_Simulator_ContextWithTimeout(t *testing.T) {
	// 把 cancel 赋值到 _ 是因为
	// 它来至于 contest.WithCancel，我 100% 信任标准库
	Convey("对于 Simulator 来说", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		ctx := context.Background()
		timeout := time.Second
		deadline := now.Add(timeout)
		Convey("放入 ctx", func() {
			child, _ := s.ContextWithTimeout(ctx, timeout)
			Convey("child 的 deadline 符合预期", func() {
				actual, ok := child.Deadline()
				So(actual, ShouldEqual, deadline)
				So(ok, ShouldBeTrue)
			})
		})
	})
}

func Test_Simulator_contextWithDeadline(t *testing.T) {
	// 把 cancel 赋值到 _ 是因为
	// 它来至于 contest.WithCancel，我 100% 信任标准库
	Convey("对于 Simulator 来说", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		ctxWithoutDeadline := context.Background()
		oneSecondLater := now.Add(time.Second * 1)
		twoSecondLater := now.Add(time.Second * 2)
		threeSecondLater := now.Add(time.Second * 3)
		Convey("如果放入 ctxWithoutDeadline", func() {
			child, _ := s.contextWithDeadline(ctxWithoutDeadline, oneSecondLater)
			Convey("child 应该是 *contextSim 类型", func() {
				_, ok := child.(*contextSim)
				So(ok, ShouldBeTrue)
			})
			Convey("child 的 deadline 符合预期", func() {
				actual, ok := child.Deadline()
				So(actual, ShouldEqual, oneSecondLater)
				So(ok, ShouldBeTrue)
			})
		})
		Convey("生成 ctxDeadTwoSecondLater", func() {
			ctxDeadTwoSecondLater, cancel := context.WithDeadline(context.Background(), twoSecondLater)
			defer cancel()
			Convey("如果 child 的 deadline 更早", func() {
				child, _ := s.contextWithDeadline(ctxDeadTwoSecondLater, oneSecondLater)
				Convey("child 应该是 *contextSim 类型", func() {
					_, ok := child.(*contextSim)
					So(ok, ShouldBeTrue)
				})
				Convey("child 的 deadline 会在 1 秒后", func() {
					actual, ok := child.Deadline()
					So(actual, ShouldEqual, oneSecondLater)
					So(ok, ShouldBeTrue)
				})
			})
			Convey("如果 child 的 deadline 更晚", func() {
				child, _ := s.contextWithDeadline(ctxDeadTwoSecondLater, threeSecondLater)
				Convey("child 不应该是 *contextSim 类型", func() {
					_, ok := child.(*contextSim)
					So(ok, ShouldBeFalse)
				})
				Convey("child 的 deadline 会在 2 秒后", func() {
					actual, ok := child.Deadline()
					So(actual, ShouldEqual, twoSecondLater)
					So(ok, ShouldBeTrue)
				})
			})
		})
	})
}
