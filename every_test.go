package clock

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_nextDayTime(t *testing.T) {
	Convey("测试 nextDayTime 函数", t, func() {
		yyyy, mm, dd := 2020, time.Month(5), 18
		h, m, s, ns := 21, 40, 34, 0
		Convey("设置时间晚于 now 的时间段", func() {
			loc := time.Now().Location()
			now := time.Date(yyyy, mm, dd, h, m, s, ns, loc)
			So(m, ShouldBeLessThan, 60)
			Convey("则下一个时间点在当天", func() {
				actual := nextDayTime(now, h, m, s)
				m++
				expected := now.Add(time.Minute)
				So(actual, ShouldEqual, expected)
			})
		})
		Convey("设置时间早于 now 的时间段", func() {
			loc := time.Now().Location()
			now := time.Date(yyyy, mm, dd, h, m, s, ns, loc)
			So(m, ShouldBeGreaterThanOrEqualTo, 0)
			Convey("则下一个时间点在下一天", func() {
				actual := nextDayTime(now, h, m, s)
				m--
				expected := now.Add(-time.Minute).Add(24 * time.Hour)
				So(actual, ShouldEqual, expected)
			})
		})
		Convey("设置时间等于 now 的时间段", func() {
			loc := time.Now().Location()
			ns = 0
			now := time.Date(yyyy, mm, dd, h, m, s, ns, loc)
			Convey("则下一个时间点在下一天", func() {
				actual := nextDayTime(now, h, m, s)
				expected := now.Add(24 * time.Hour)
				So(actual, ShouldEqual, expected)
			})
		})
	})
}
