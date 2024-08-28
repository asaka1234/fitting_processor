package utils

import "time"

// 把一个小时 按照一段sliceLength分钟 进行切割
// 返回ctime所在段的起始时间
func GetSliceStartTime(ctime time.Time, sliceLength int) time.Time {
	sTime5 := (ctime.Minute() / sliceLength) * sliceLength
	minuteTime := time.Date(ctime.Year(), ctime.Month(), ctime.Day(), ctime.Hour(), sTime5, 0, 0, ctime.Location())
	return minuteTime
}
