package utils

import (
	"fmt"
	"strings"
	"time"
)

// Seconds-based time units
const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

func computeTimeDiff(diff int64) (int64, string) {
	diffStr := ""
	switch {
	case diff <= 0:
		diff = 0
		diffStr = "now"
	case diff < 2:
		diff = 0
		diffStr = "1 second"
	case diff < 1*Minute:
		diffStr = fmt.Sprintf("%d seconds", diff)
		diff = 0

	case diff < 2*Minute:
		diff -= 1 * Minute
		diffStr = "1 minute"
	case diff < 1*Hour:
		diffStr = fmt.Sprintf("%d minutes", diff/Minute)
		diff -= diff / Minute * Minute

	case diff < 2*Hour:
		diff -= 1 * Hour
		diffStr = "1 hour"
	case diff < 1*Day:
		diffStr = fmt.Sprintf("%d hours", diff/Hour)
		diff -= diff / Hour * Hour

	case diff < 2*Day:
		diff -= 1 * Day
		diffStr = "1 day"
	case diff < 1*Week:
		diffStr = fmt.Sprintf("%d days", diff/Day)
		diff -= diff / Day * Day

	case diff < 2*Week:
		diff -= 1 * Week
		diffStr = "1 week"
	case diff < 1*Month:
		diffStr = fmt.Sprintf("%d weeks", diff/Week)
		diff -= diff / Week * Week

	case diff < 2*Month:
		diff -= 1 * Month
		diffStr = "1 month"
	case diff < 1*Year:
		diffStr = fmt.Sprintf("%d months", diff/Month)
		diff -= diff / Month * Month

	case diff < 2*Year:
		diff -= 1 * Year
		diffStr = "1 year"
	default:
		diffStr = fmt.Sprintf("%d years", diff/Year)
		diff = 0
	}
	return diff, diffStr
}

func TimeSincePro(then time.Time) string {
	now := time.Now()
	diff := now.Unix() - then.Unix()

	if then.After(now) {
		return "future"
	}

	var timeStr, diffStr string
	for {
		if diff == 0 {
			break
		}

		diff, diffStr = computeTimeDiff(diff)
		timeStr += ", " + diffStr
	}
	return strings.TrimPrefix(timeStr, ", ")
}


const TIMEFORMAT = "2006-01-02 15:04:05" // 固定时间，不可更改
const TIMEFORMAT1 = "2006-01-02T15:04:05.000Z"

const DATEORMAT = "2006-01-02" // 固定时间，不可更改

// 此时此刻 Time 类型 时间
func NowTime() time.Time {
	return time.Now()
}

// 此时此刻 秒级时间戳
func NowInt64() int64 {
	return TimeTime2Int64(NowTime())
}

// 此时此刻 YYYY-MM-DD hh:mm:ss
func NowFormat() string {
	return TimeTime2String(NowTime())
}

// 此时的日期
func NowDateFormat() string {
	return NowTime().Format(DATEORMAT)
}

// 此时此刻 周几
func NowWeek() int {
	return TimeTime2Week(NowTime())
}

// Time2Date 时间转日期
func Time2Date(t time.Time) string {
	return t.Format(DATEORMAT)
}

// 时间戳 转 time 类型
func TimeInt642Time(second int64) time.Time {
	return time.Unix(second, 0)
}

// 时间戳 转 string 类型
func TimeInt642String(second int64) string {
	return time.Unix(second, 0).Format(TIMEFORMAT)
}

// string 转 time 类型
func TimeString2Time(t string) (time.Time, error) {
	return time.ParseInLocation(TIMEFORMAT, t, time.Local)
}

// string 转 time 类型
func TimeString2Int64(t string) (int64, error) {
	timeType, err := time.ParseInLocation(TIMEFORMAT1, t, time.Local)
	if err != nil {
		return NowInt64(), err
	}
	return timeType.Unix(), nil
}

// time 类型 转 string
func TimeTime2String(t time.Time) string {
	return t.Format(TIMEFORMAT)
}

// time 类型 转 时间戳（秒级）
func TimeTime2Int64(t time.Time) int64 {
	return t.Unix()
}

// time 类型 转 周几
func TimeTime2Week(t time.Time) int {
	return int(t.Weekday())
}

// utc 时间 转 本地时间
func TimeUtc2Local(t time.Time) time.Time {
	return t.In(time.Local)
}

// 本地时间 转 utc 时间
func TimeLocal2Utc(t time.Time) time.Time {
	return t.UTC()
}

// int 转 time.Duration
func TimeInt2Duration(i int) time.Duration {
	return time.Duration(i)
}

// 每xx 的某个时间 开始时间 xx代表，分钟，秒钟，小时
func TimeBegin(t time.Time, duration time.Duration) time.Time {
	return t.Round(duration)
}

// 每xx 的某个时间 结束时间 xx代表，分钟，秒钟，小时
func TimeEnd(t time.Time, duration time.Duration) time.Time {
	return t.Truncate(duration)
}

// 每分钟开始时间 2006-01-02 15:04:05 --> 2006-01-02 15:04:00
func TimeMinuteBegin(t time.Time) time.Time {
	return TimeBegin(t, time.Minute)
}

// 每小时开始时间 2006-01-02 15:04:05 --> 2006-01-02 15:00:00
func TimeHourBegin(t time.Time) time.Time {
	return TimeBegin(t, time.Hour)
}

// 每天开始时间 2006-01-02 15:04:05 --> 2006-01-02 00:00:00
func TimeDayBegin(t time.Time) time.Time {
	currentYear, currentMonth, currentDay := t.Date()
	currentLocation := t.Location()
	return time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, currentLocation)
}

// 每月开始的时间 2006-01-02 15:04:05 --> 2006-01-01 00:00:00
func TimeMonthBegin(t time.Time) time.Time {
	currentYear, currentMonth, _ := t.Date()
	currentLocation := t.Location()
	return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
}

// TimeWeekBegin 每周1，开始时间
func TimeWeekBegin(t time.Time) time.Time {
	w := TimeTime2Week(t)
	if w == 1 {
		return TimeDayBegin(t)
	} else if w == 0 {
		w = 7
	}
	return t.AddDate(0, 0, (1 - w))
}

// 每分钟结束时间 2006-01-02 15:04:05 --> 2006-01-02 15:05:00
func TimeMinuteEnd(t time.Time) time.Time {
	return TimeEnd(t, time.Minute)
}

// 每小时结束时间 2006-01-02 15:04:05 --> 2006-01-02 16:00:00
func TimeHourEnd(t time.Time) time.Time {
	return TimeEnd(t, time.Hour)
}

// 每天结束时间 2006-01-02 15:04:05 --> 2006-01-03 00:00:00
func TimeDayEnd(t time.Time) time.Time {
	return TimeDayBegin(t).AddDate(0, 0, 1)
}

// 每月结束时间 2006-01-02 15:04:05 --> 2006-02-01 00:00:00
func TimeMonthEnd(t time.Time) time.Time {
	return TimeMonthBegin(t).AddDate(0, 1, 0)
}

// 几秒钟，分钟，小时，之前或之后
func TimeAdd(t time.Time, duration time.Duration, l int) time.Time {
	return t.Add(TimeInt2Duration(l) * duration)
}

// 几天，月，年，之前或之后
func TimeDateAdd(t time.Time, years int, months int, days int) time.Time {
	return t.AddDate(years, months, days)
}

// 几秒钟前后 l，正负值 l：几
func TimeSecondBeforeAfter(t time.Time, l int) time.Time {
	return TimeAdd(t, time.Second, l)
}

// 几分钟前后 l，正负值 l：几
func TimeMinuteBeforeAfter(t time.Time, l int) time.Time {
	return TimeAdd(t, time.Minute, l)
}

// 几小时前后 l，正负值 l：几
func TimeHourBeforeAfter(t time.Time, l int) time.Time {
	return TimeAdd(t, time.Hour, l)
}

// 几天前后 l，正负值 l：几
func TimeDayBeforeAfter(t time.Time, l int) time.Time {
	return TimeDateAdd(t, 0, 0, l)
}

// 几月前后 l，正负值 l：几
func TimeMonthBeforeAfter(t time.Time, l int) time.Time {
	return TimeDateAdd(t, 0, l, 0)
}

// 返回本周的周一的时间 2018-10-28
func WeekStarString() string {
	return Time2Date(TimeWeekBegin(NowTime()))
}
