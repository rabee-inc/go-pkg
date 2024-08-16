package timefmt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type timeFormatter struct {
	weekdayFullNames  []string
	weekdayShortNames []string
	weekdayMinNames   []string
	monthNames        []string
	monthShortNames   []string
	meridiemFunc      func(int) string
}

var defaultWeekdayFullNames = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var defaultWeekdayShortNames = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
var defaultWeekdayMinNames = []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
var defaultMonthNames = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
var defaultMonthShortNames = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
var defaultMeridiemFunc = func(hours int) string {
	if hours < 12 {
		return "am"
	}
	return "pm"
}

// REF: https://github.com/iamkun/dayjs/blob/dev/src/constant.js#L30C30-L30C112
var reFormat = regexp.MustCompile(`\[([^\]]+)]|Y{1,4}|M{1,4}|D{1,2}|d{1,4}|H{1,2}|h{1,2}|a|A|m{1,2}|s{1,2}|Z{1,2}|SSS`)

var weekLen = 7
var monthLen = 12

// New returns a new instance of TimeFormatter.
func New() TimeFormatter {
	return &timeFormatter{
		weekdayFullNames:  defaultWeekdayFullNames,
		weekdayShortNames: defaultWeekdayShortNames,
		weekdayMinNames:   defaultWeekdayMinNames,
		monthNames:        defaultMonthNames,
		monthShortNames:   defaultMonthShortNames,
		meridiemFunc:      nil,
	}
}

func (t *timeFormatter) Format(tm time.Time, layout string) string {
	return reFormat.ReplaceAllStringFunc(layout, func(s string) string {
		// [] で囲まれている場合は、フォーマットせずに[]の中身を返す
		if s[0] == '[' {
			return s[1 : len(s)-1]
		}
		return t.format(tm, s)
	})
}

func (t *timeFormatter) MonthNames() []string {
	return t.monthNames
}

func (t *timeFormatter) MonthShortNames() []string {
	return t.monthShortNames
}

func (t *timeFormatter) SetMonthNames(names []string) {
	mustLen(names, monthLen)
	t.monthNames = names
}

func (t *timeFormatter) SetMonthShortNames(names []string) {
	mustLen(names, monthLen)
	t.monthShortNames = names
}

func (t *timeFormatter) SetWeekdayFullNames(names []string) {
	mustLen(names, weekLen)
	t.weekdayFullNames = names
}

func (t *timeFormatter) SetWeekdayMinNames(names []string) {
	mustLen(names, weekLen)
	t.weekdayMinNames = names
}

func (t *timeFormatter) SetWeekdayShortNames(names []string) {
	mustLen(names, weekLen)
	t.weekdayShortNames = names
}

func (t *timeFormatter) WeekdayFullNames() []string {
	return t.weekdayFullNames
}

func (t *timeFormatter) WeekdayMinNames() []string {
	return t.weekdayMinNames
}

func (t *timeFormatter) WeekdayShortNames() []string {
	return t.weekdayShortNames
}

func (t *timeFormatter) SetMeridiem(f func(hours int) string) {
	t.meridiemFunc = f
}

// mustLen ... names の長さが l でない場合 panic にする
func mustLen(names []string, l int) {
	if len(names) != l {
		panic(fmt.Sprintf("len error: names length must be %d but %d", l, len(names)))
	}
}

// chunk に応じたフォーマットを行う
func (t *timeFormatter) format(tm time.Time, chunk string) string {
	switch chunk {
	// year
	case "YYYY":
		return strconv.Itoa(tm.Year())
	case "YY":
		return fmt.Sprintf("%02d", tm.Year()%100)

	// month
	case "M":
		return strconv.Itoa(int(tm.Month()))
	case "MM":
		return fmt.Sprintf("%02d", int(tm.Month()))
	case "MMM":
		return t.monthShortNames[int(tm.Month())-1]
	case "MMMM":
		return t.monthNames[int(tm.Month())-1]

	// day
	case "D":
		return strconv.Itoa(tm.Day())
	case "DD":
		return fmt.Sprintf("%02d", tm.Day())

	// hours
	case "H":
		return strconv.Itoa(tm.Hour())
	case "HH":
		return fmt.Sprintf("%02d", tm.Hour())
	case "h", "hh":
		h := tm.Hour() % 12
		if h == 0 {
			h = 12
		}
		if chunk == "h" {
			return strconv.Itoa(h)
		}
		return fmt.Sprintf("%02d", h)

	// minutes
	case "m":
		return strconv.Itoa(tm.Minute())
	case "mm":
		return fmt.Sprintf("%02d", tm.Minute())

	// seconds
	case "s":
		return strconv.Itoa(tm.Second())
	case "ss":
		return fmt.Sprintf("%02d", tm.Second())

	// milliseconds
	case "SSS":
		return fmt.Sprintf("%03d", tm.Nanosecond()/1e6)

	// AM/PM
	case "a", "A":
		meridiemFunc := t.meridiemFunc
		if meridiemFunc == nil {
			meridiemFunc = defaultMeridiemFunc
		}
		meridiem := meridiemFunc(tm.Hour())
		if chunk == "a" {
			return strings.ToLower(meridiem)
		}
		return strings.ToUpper(meridiem)

	// weekday
	case "d":
		return strconv.Itoa(int(tm.Weekday()))
	case "dd":
		return t.weekdayMinNames[int(tm.Weekday())]
	case "ddd":
		return t.weekdayShortNames[int(tm.Weekday())]
	case "dddd":
		return t.weekdayFullNames[int(tm.Weekday())]

	// timezone
	case "Z", "ZZ":
		_, offset := tm.Zone()
		if offset == 0 {
			return "Z"
		}
		sign := "+"
		if offset < 0 {
			sign = "-"
			offset = -offset
		}
		hour := offset / 3600
		min := (offset % 3600) / 60
		if chunk == "Z" {
			return fmt.Sprintf("%s%02d:%02d", sign, hour, min)
		}
		return fmt.Sprintf("%s%02d%02d", sign, hour, min)

	// match しない場合はそのまま返す
	default:
		return chunk
	}
}
