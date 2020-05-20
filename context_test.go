package clock

import (
	"context"
	"testing"
	"time"

	. "github.com/golang/mock/gomock"
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
		//
		ctrl := NewController(t)
		defer ctrl.Finish()
		//
		timeChan := make(chan time.Time, 1)
		now := time.Now()
		go func() {
			timeChan <- now
			close(timeChan)
		}()
		//
		mockClock := NewMockClock(ctrl)
		mockClock.EXPECT().After(Any()).Return(timeChan)
		//
		ctx = Set(ctx, mockClock)
		//
		Convey("应该返回特定的时间", func() {
			actual, ok := <-After(ctx, time.Second)
			So(actual, ShouldEqual, now)
			So(ok, ShouldBeTrue)
		})
	})
}

func Test_AfterFunc(t *testing.T) {
	Convey("测试 After", t, func() {
		ctx := context.Background()
		//
		ctrl := NewController(t)
		defer ctrl.Finish()
		//
		timer := &Timer{}
		//
		mockClock := NewMockClock(ctrl)
		mockClock.EXPECT().AfterFunc(Any(), Any()).Return(timer)
		//
		ctx = Set(ctx, mockClock)
		//
		Convey("应该返回特定的时间", func() {
			actual := AfterFunc(ctx, time.Second, nil)
			So(actual, ShouldEqual, timer)
		})
	})
}

func Test_NewTicker(t *testing.T) {
	Convey("测试 NewTicker", t, func() {
		ctx := context.Background()
		//
		ctrl := NewController(t)
		defer ctrl.Finish()
		//
		ticker := &Ticker{}
		//
		mockClock := NewMockClock(ctrl)
		mockClock.EXPECT().NewTicker(Any()).Return(ticker)
		//
		ctx = Set(ctx, mockClock)
		//
		Convey("返回值的类型应该符合预期", func() {
			actual := NewTicker(ctx, time.Second)
			So(actual, ShouldEqual, ticker)
		})
	})
}

func Test_NewTimer(t *testing.T) {
	Convey("测试 NewTimer", t, func() {
		ctx := context.Background()
		//
		ctrl := NewController(t)
		defer ctrl.Finish()
		//
		timer := &Timer{}
		//
		mockClock := NewMockClock(ctrl)
		mockClock.EXPECT().NewTimer(Any()).Return(timer)
		//
		ctx = Set(ctx, mockClock)
		//
		Convey("返回值的类型应该符合预期", func() {
			actual := NewTimer(ctx, time.Second)
			So(actual, ShouldEqual, timer)
		})
	})
}

func Test_Now(t *testing.T) {
	Convey("测试 Now", t, func() {
		ctx := context.Background()
		//
		ctrl := NewController(t)
		defer ctrl.Finish()
		//
		now := time.Now()
		//
		mockClock := NewMockClock(ctrl)
		mockClock.EXPECT().Now().Return(now)
		//
		ctx = Set(ctx, mockClock)
		//
		Convey("返回值的类型应该符合预期", func() {
			actual := Now(ctx)
			So(actual, ShouldEqual, now)
		})
	})
}

func Test_Since(t *testing.T) {
	Convey("测试 Since", t, func() {
		ctx := context.Background()
		//
		ctrl := NewController(t)
		defer ctrl.Finish()
		//
		dur := time.Second
		//
		mockClock := NewMockClock(ctrl)
		mockClock.EXPECT().Since(Any()).Return(dur)
		//
		ctx = Set(ctx, mockClock)
		//
		Convey("返回值的类型应该符合预期", func() {
			actual := Since(ctx, time.Now())
			So(actual, ShouldEqual, dur)
		})
	})
}

func Test_Sleep(t *testing.T) {
	Convey("测试 Sleep", t, func() {
		ctx := context.Background()
		Convey("的确休眠了那么长时间", func() {
			dur := time.Millisecond * 100
			delta := dur / 10
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
		//
		ctrl := NewController(t)
		defer ctrl.Finish()
		//
		timeChan := make(<-chan time.Time)
		//
		mockClock := NewMockClock(ctrl)
		mockClock.EXPECT().Tick(Any()).Return(timeChan)
		//
		ctx = Set(ctx, mockClock)
		//
		Convey("返回值的类型应该符合预期", func() {
			actual := Tick(ctx, time.Second)
			So(actual, ShouldEqual, timeChan)
		})
	})
}

func Test_Until(t *testing.T) {
	Convey("测试 Until", t, func() {
		ctx := context.Background()
		//
		ctrl := NewController(t)
		defer ctrl.Finish()
		//
		dur := time.Second
		//
		mockClock := NewMockClock(ctrl)
		mockClock.EXPECT().Until(Any()).Return(dur)
		//
		ctx = Set(ctx, mockClock)
		//
		Convey("返回值的类型应该符合预期", func() {
			actual := Until(ctx, time.Now())
			So(actual, ShouldEqual, dur)
		})
	})
}

func Test_ContextWithDeadline(t *testing.T) {
	Convey("测试 ContextWithDeadline", t, func() {
		ctx := context.Background()
		//
		ctrl := NewController(t)
		defer ctrl.Finish()
		//
		var cancel context.CancelFunc
		//
		mockClock := NewMockClock(ctrl)
		mockClock.EXPECT().ContextWithDeadline(Any(), Any()).Return(context.TODO(), cancel)
		//
		ctx = Set(ctx, mockClock)
		//
		Convey("返回值的类型应该符合预期", func() {
			deadline := time.Now().Add(time.Second)
			actual, actualCancel := ContextWithDeadline(ctx, deadline)
			So(actual, ShouldEqual, context.TODO())
			So(actualCancel, ShouldEqual, cancel)
		})
	})
}

func Test_ContextWithTimeout(t *testing.T) {
	Convey("测试 ContextWithTimeout", t, func() {
		ctx := context.Background()
		//
		ctrl := NewController(t)
		defer ctrl.Finish()
		//
		var cancel context.CancelFunc
		//
		mockClock := NewMockClock(ctrl)
		mockClock.EXPECT().ContextWithTimeout(Any(), Any()).Return(context.TODO(), cancel)
		//
		ctx = Set(ctx, mockClock)
		//
		Convey("返回值的类型应该符合预期", func() {
			timeout := time.Second
			actual, actualCancel := ContextWithTimeout(ctx, timeout)
			So(actual, ShouldEqual, context.TODO())
			So(actualCancel, ShouldEqual, cancel)
		})
	})
}
