package timefmt

import "time"

type TimeFormatter interface {
	// dayjs like format
	// REF: https://day.js.org/docs/en/display/format
	// ISO8601 format example: "YYYY-MM-DDTHH:mm:ss.SSSZ"
	// YYYY: 4 digit year (2024)
	// YY: 2 digit year (24)
	// MMMM: full month name (January)
	// MMM: short month name (Jan)
	// MM: 2 digit month (01-12)
	// M: 1 digit month (1-12)
	// DD: 2 digit day (01-31)
	// D: 1 digit day (1-31)
	// HH: 2 digit hour (00-23)
	// H: 1 digit hour (0-23)
	// hh: 2 digit hour (00-12)
	// h: 1 digit hour (0-12)
	// mm: 2 digit minute (00-59)
	// m: 1 digit minute (0-59)
	// ss: 2 digit second (00-59)
	// s: 1 digit second (0-59)
	// SSS: 3 digit millisecond (000-999)
	// A: upper case meridiem (AM/PM)
	// a: lower case meridiem (am/pm)
	// Z: time zone offset (+09:00/Z). if time zone is UTC, return "Z".
	// ZZ: time zone offset (+0900/Z). if time zone is UTC, return "Z".
	// dddd: full day of the week (Sunday)
	// ddd: short day of the week (Sun)
	// dd: min day of the week (Su)
	// d: 1 digit day of the week (0-6)
	Format(t time.Time, layout string) string

	// dddd で使用される曜日の full name を設定する。
	// len(names) != 7 の場合 panic。
	// ex) []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	// ex) []string{"日曜日", "月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日"}
	SetWeekdayFullNames(names []string)

	// dddd で使用される曜日の short name を取得する
	WeekdayFullNames() []string

	// ddd で使用される曜日の short name を設定する。
	// len(names) != 7 の場合 panic。
	// ex) []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	// ex) []string{"日", "月", "火", "水", "木", "金", "土"}
	SetWeekdayShortNames(names []string)

	// ddd で使用される曜日の short name を取得する。
	WeekdayShortNames() []string

	// dd で使用される曜日の min name を設定する。
	// len(names) != 7 の場合 panic。
	// ex) []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	// ex) []string{"日", "月", "火", "水", "木", "金", "土"}
	SetWeekdayMinNames(names []string)

	// dd で使用される曜日の min name を取得する。
	WeekdayMinNames() []string

	// MMMM で使用される month の name を設定する。
	// len(names) != 12 の場合 panic。
	// ex) []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	// ex) []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"}
	SetMonthNames(names []string)

	// MMMM で使用される month の name を取得する
	MonthNames() []string

	// MMM で使用される month の short name を設定する。
	// len(names) != 12 の場合 panic。
	// ex) []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	// ex) []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"}
	SetMonthShortNames(names []string)

	// MMM で使用される month の short name を取得する
	MonthShortNames() []string

	// AM/PM を判定する関数を設定する。
	// hours には 24 時間表記の時間が入る。(0〜23)
	SetMeridiem(f func(hours int) string)
}
