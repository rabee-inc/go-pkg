package timeutil

import (
	"time"
)

// 現在時刻をJSTのTimeで取得する
func Now() time.Time {
	return time.Now().In(ZoneJST())
}

// 現在時刻をJSTのUnixTime(ミリ秒)で取得する
func NowUnix() int64 {
	return time.Now().In(ZoneJST()).UnixNano() / int64(time.Millisecond)
}

// UnixTime(ミリ秒)からJSTのTimeを取得する
func ByUnix(u int64) time.Time {
	uNano := u * 1000 * 1000
	uSec := u / 1000
	return time.Unix(uSec, uNano-(uSec*1000*1000*1000)).In(ZoneJST())
}

// TimeからUnixTime(ミリ秒)に変換する
func ToUnix(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// 日本のタイムゾーンを取得する
func ZoneJST() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}

// 指定秒数をミリ秒に変換する
func SecondsToMilliseconds(seconds int) int64 {
	return int64(seconds * 1000)
}

// 指定分数をミリ秒に変換する
func MinutesToMilliseconds(minutes int) int64 {
	return int64(minutes * 60 * 1000)
}

// 指定時数をミリ秒に変換する
func HoursToMilliseconds(hours int) int64 {
	return int64(hours * 60 * 60 * 1000)
}

// 指定日数をミリ秒に変換する
func DaysToMilliseconds(days int) int64 {
	return int64(days * 24 * 60 * 60 * 1000)
}

// 本日か判定する
func IsToday(u int64) bool {
	t := Now()
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, ZoneJST())
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, ZoneJST())
	return ToUnix(startTime) <= u && u <= ToUnix(endTime)
}

// 日の範囲(0:00:00〜23:59:59)を取得する
func DayPeriod(at int64, diff int) (int64, int64) {
	t := ByUnix(at)
	t = t.AddDate(0, 0, diff)
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, ZoneJST())
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, ZoneJST())
	return ToUnix(startTime), ToUnix(endTime)
}

// 月の範囲(1日0:00:00〜最終日23:59:59)を取得する
func MonthPeriod(at int64, diff int) (int64, int64) {
	t := ByUnix(at)
	t = t.AddDate(0, diff, 0)
	startTime := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, ZoneJST())
	endTime := time.Date(t.Year(), t.Month(), 1, 23, 59, 59, 0, ZoneJST()).AddDate(0, 1, -1)
	return ToUnix(startTime), ToUnix(endTime)
}

// 年の範囲(12月1日0:00:00〜12月31日23:59:59)を取得する
func YearPeriod(at int64, diff int) (int64, int64) {
	t := ByUnix(at)
	t = t.AddDate(diff, 0, 0)
	startTime := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, ZoneJST())
	endTime := time.Date(t.Year(), 1, 1, 23, 59, 59, 0, ZoneJST()).AddDate(1, 0, -1)
	return ToUnix(startTime), ToUnix(endTime)
}

// 月の最終日を取得する
func LastDayByMonth(at int64) int {
	t := ByUnix(at)
	t = time.Date(t.Year(), t.Month(), 1, 23, 59, 59, 0, ZoneJST()).AddDate(0, 1, -1)
	return t.Day()
}

// 曜日(日本語)を取得する
func GetWeekJP(t time.Time) string {
	var week string
	switch t.Weekday() {
	case time.Sunday:
		week = "日"
	case time.Monday:
		week = "月"
	case time.Tuesday:
		week = "火"
	case time.Wednesday:
		week = "水"
	case time.Thursday:
		week = "木"
	case time.Friday:
		week = "金"
	case time.Saturday:
		week = "土"
	}
	return week
}
