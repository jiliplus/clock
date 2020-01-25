package clock

import (
	"fmt"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_realClock_x(t *testing.T) {
	Convey("新建一个 realClock", t, func() {
		c := NewRealClock()
		fmt.Println(c.Now())
	})
}

func Test_realClock_After(t *testing.T) {
	Convey("新建一个 realClock", t, func() {
		c := NewRealClock()
		fmt.Println(c.Now())
		expected := c.Now()
		patches := gomonkey.ApplyFunc(time.After, func(_ time.Duration) <-chan time.Time {
			res := make(chan time.Time)
			go func() {
				res <- expected
				close(res)
			}()
			return res
		})
		defer patches.Reset()
		actual := <-c.After(time.Second)
		Convey("realClock 应该调用了 time.After", func() {
			So(actual, ShouldEqual, expected)
		})
	})
}
