package timeutil

import (
	"time"
)

// Now ... 現在時刻をJSTのTimeで取得する
func Now() time.Time {
	return time.Now().In(ZoneJST())
}

// NowUnix ... 現在時刻をJSTのUnixtimestamp(ミリ秒)で取得する
func NowUnix() int64 {
	return time.Now().In(ZoneJST()).UnixNano() / int64(time.Millisecond)
}

// ByUnix ... Unixtimestamp(ミリ秒)からJSTのTimeを取得する
func ByUnix(u int64) time.Time {
	uNano := u * 1000 * 1000
	uSec := u / 1000
	return time.Unix(uSec, uNano-(uSec*1000*1000*1000)).In(ZoneJST())
}

// ToUnix ... TimeからUnixtimestamp(ミリ秒)に変換する
func ToUnix(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// ZoneJST ... 日本のタイムゾーンを取得する
func ZoneJST() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}

// SecondsToMiliseconds ... 指定秒数をミリ秒に変換する
func SecondsToMiliseconds(seconds int) int64 {
	return int64(seconds * 1000)
}

// MinutesToMiliseconds ... 指定分数をミリ秒に変換する
func MinutesToMiliseconds(minutes int) int64 {
	return int64(minutes * 60 * 1000)
}

// HoursToMiliseconds ... 指定時数をミリ秒に変換する
func HoursToMiliseconds(hours int) int64 {
	return int64(hours * 60 * 60 * 1000)
}

// DaysToMiliseconds ... 指定日数をミリ秒に変換する
func DaysToMiliseconds(days int) int64 {
	return int64(days * 24 * 60 * 60 * 1000)
}

// IsToday ... 本日か判定する
func IsToday(u int64) bool {
	t := Now()
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, ZoneJST())
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, ZoneJST())
	return ToUnix(startTime) <= u && u <= ToUnix(endTime)
}

// DayPeriod ... 日の範囲(0:00:00〜23:59:59)を取得する
func DayPeriod(at int64, diff int) (int64, int64) {
	t := ByUnix(at)
	t = t.AddDate(0, 0, diff)
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, ZoneJST())
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, ZoneJST())
	return ToUnix(startTime), ToUnix(endTime)
}

// MonthPeriod ... 月の範囲(1日0:00:00〜最終日23:59:59)を取得する
func MonthPeriod(at int64, diff int) (int64, int64) {
	t := ByUnix(at)
	t = t.AddDate(0, diff, 0)
	startTime := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, ZoneJST())
	endTime := time.Date(t.Year(), t.Month(), 1, 23, 59, 59, 0, ZoneJST()).AddDate(0, 1, -1)
	return ToUnix(startTime), ToUnix(endTime)
}

// YearPeriod ... 年の範囲(12月1日0:00:00〜12月31日23:59:59)を取得する
func YearPeriod(at int64, diff int) (int64, int64) {
	t := ByUnix(at)
	t = t.AddDate(diff, 0, 0)
	startTime := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, ZoneJST())
	endTime := time.Date(t.Year(), 1, 1, 23, 59, 59, 0, ZoneJST()).AddDate(1, 0, -1)
	return ToUnix(startTime), ToUnix(endTime)
}

// LastDayByMonth ... 月の最終日を取得する
func LastDayByMonth(at int64) int {
	t := ByUnix(at)
	t = time.Date(t.Year(), t.Month(), 1, 23, 59, 59, 0, ZoneJST()).AddDate(0, 1, -1)
	return t.Day()
}
