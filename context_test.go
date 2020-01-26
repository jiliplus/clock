package clock

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_contextSim(t *testing.T) {
	Convey("利用 Simulator 新建 Context", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		parent, cancel := context.WithCancel(context.Background())
		deadline := now.Add(time.Second)
		cs := s.newContextSim(parent, deadline)
		Convey("子文已经具备了 deadline", func() {
			actual, ok := cs.Deadline()
			So(actual, ShouldEqual, deadline)
			So(ok, ShouldBeTrue)
		})
		Convey("新子文的还不会自己改变", func() {
			err := cs.Err()
			So(err, ShouldBeNil)
		})
		Convey("父文取消的话，子文也会结束", func() {
			cancel()
			<-cs.Done()
			So(s.Now().Before(deadline), ShouldBeTrue)
			err := cs.Err()
			So(err, ShouldEqual, context.Canceled)
		})
		Convey("子文到期的话，也会结束", func() {
			s.SetOrPanic(deadline)
			<-cs.Done()
			err := cs.Err()
			So(err.Error(), ShouldEqual, context.DeadlineExceeded.Error())
		})
	})
}

func Test_Set_and_Get(t *testing.T) {
	Convey("测试 Set 和 Get", t, func() {
		now := time.Now()
		s := NewSimulator(now)
		parent := context.Background()
		Convey("parent 中取出 realClock", func() {
			_, ok := Get(parent).(realClock)
			So(ok, ShouldBeTrue)
		})
		child := Set(parent, s)
		Convey("child 中可以取出 s", func() {
			actual, ok := Get(child).(*Simulator)
			So(ok, ShouldBeTrue)
			So(actual, ShouldEqual, s)
		})
	})
}

func Test_After(t *testing.T) {
	Convey("测试 After", t, func() {
		ctx := context.Background()
		Convey("返回值的类型应该符合预期", func() {
			actual := After(ctx, time.Second)
			expected := make(<-chan time.Time)
			So(actual, ShouldHaveSameTypeAs, expected)
		})
	})
}

func Test_AfterFunc(t *testing.T) {
	Convey("测试 AfterFunc", t, func() {
		ctx := context.Background()
		Convey("返回值的类型应该符合预期", func() {
			actual := AfterFunc(ctx, time.Second, func() {})
			expected := &Timer{}
			So(actual, ShouldHaveSameTypeAs, expected)
		})
	})
}

func Test_NewTicker(t *testing.T) {
	Convey("测试 NewTicker", t, func() {
		ctx := context.Background()
		Convey("返回值的类型应该符合预期", func() {
			actual := NewTicker(ctx, time.Second)
			expected := &Ticker{}
			So(actual, ShouldHaveSameTypeAs, expected)
		})
	})
}

func Test_NewTimer(t *testing.T) {
	Convey("测试 NewTimer", t, func() {
		ctx := context.Background()
		Convey("返回值的类型应该符合预期", func() {
			actual := NewTimer(ctx, time.Second)
			expected := &Timer{}
			So(actual, ShouldHaveSameTypeAs, expected)
		})
	})
}

func Test_Now(t *testing.T) {
	Convey("测试 Now", t, func() {
		ctx := context.Background()
		Convey("返回值的类型应该符合预期", func() {
			actual := Now(ctx)
			expected := time.Now()
			So(actual, ShouldHaveSameTypeAs, expected)
		})
	})
}

func Test_Since(t *testing.T) {
	Convey("测试 Since", t, func() {
		ctx := context.Background()
		Convey("返回值的类型应该符合预期", func() {
			actual := Since(ctx, time.Now())
			expected := time.Second
			So(actual, ShouldHaveSameTypeAs, expected)
		})
	})
}

func Test_Sleep(t *testing.T) {
	Convey("测试 Sleep", t, func() {
		ctx := context.Background()
		Convey("返回值的类型应该符合预期", func() {
			dur := time.Millisecond * 10
			delta := time.Millisecond
			start := time.Now()
			Sleep(ctx, dur)
			end := time.Now()
			So(end, ShouldHappenWithin, dur+delta, start)
		})
	})
}

func Test_Tick(t *testing.T) {
	Convey("测试 Tick", t, func() {
		ctx := context.Background()
		Convey("返回值的类型应该符合预期", func() {
			actual := Tick(ctx, time.Second)
			expected := make(<-chan time.Time)
			So(actual, ShouldHaveSameTypeAs, expected)
		})
	})
}

func Test_Until(t *testing.T) {
	Convey("测试 Until", t, func() {
		ctx := context.Background()
		Convey("返回值的类型应该符合预期", func() {
			actual := Until(ctx, time.Now())
			expected := time.Second
			So(actual, ShouldHaveSameTypeAs, expected)
		})
	})
}

func Test_ContextWithDeadline(t *testing.T) {
	Convey("测试 ContextWithDeadline", t, func() {
		ctx := context.Background()
		Convey("返回值的类型应该符合预期", func() {
			deadline := time.Now().Add(time.Second)
			ctx, cancel := ContextWithDeadline(ctx, deadline)
			eCtx, eCancel := context.WithDeadline(ctx, deadline)
			So(ctx, ShouldHaveSameTypeAs, eCtx)
			So(cancel, ShouldHaveSameTypeAs, eCancel)
		})
	})
}

func Test_ContextWithTimeout(t *testing.T) {
	Convey("测试 ContextWithTimeout", t, func() {
		ctx := context.Background()
		Convey("返回值的类型应该符合预期", func() {
			timeout := time.Second
			aCtx, aCancel := ContextWithTimeout(ctx, timeout)
			eCtx, eCancel := context.WithTimeout(ctx, timeout)
			So(aCtx, ShouldHaveSameTypeAs, eCtx)
			So(aCancel, ShouldHaveSameTypeAs, eCancel)
		})
	})
}
