package clock

import "time"

// nextTime 是为了 EveryDay 函数准备的
// 是为了每天在 hour:minute:second 提供 tick
// NOTICE: hour 是 24 小时制的小时
func nextTime(now time.Time, hour, minute, second int) time.Time {
	yyyy, mm, dd := now.Year(), now.Month(), now.Day()
	loc := time.Now().Location()
	next := time.Date(yyyy, mm, dd, hour, minute, second, 0, loc)
	if now.Before(next) {
		return next
	}
	return next.Add(24 * time.Hour)
}
