package timefmt

import (
	"reflect"
	"testing"
	"time"
)

var monthNamesJP = []string{"睦月", "如月", "弥生", "卯月", "皐月", "水無月", "文月", "葉月", "長月", "神無月", "霜月", "師走"}
var monthShortNamesJP = []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"}
var weekdayFullNamesJP = []string{"日曜日", "月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日"}
var weekdayShortNamesJP = []string{"日曜", "月曜", "火曜", "水曜", "木曜", "金曜", "土曜"}
var weekdayMinNamesJP = []string{"日", "月", "火", "水", "木", "金", "土"}
var meridiemFuncJP = func(hours int) string {
	if hours < 12 {
		return "午前"
	}
	return "午後"
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want TimeFormatter
	}{
		{
			name: "New",
			want: &timeFormatter{
				weekdayFullNames:  defaultWeekdayFullNames,
				weekdayShortNames: defaultWeekdayShortNames,
				weekdayMinNames:   defaultWeekdayMinNames,
				monthNames:        defaultMonthNames,
				monthShortNames:   defaultMonthShortNames,
				meridiemFunc:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeFormatter_Format(t *testing.T) {
	type fields struct {
		weekdayFullNames  []string
		weekdayShortNames []string
		weekdayMinNames   []string
		monthNames        []string
		monthShortNames   []string
		meridiemFunc      func(int) string
	}
	type args struct {
		tm     time.Time
		layout string
	}
	loc, _ := time.LoadLocation("Asia/Tokyo")
	tm := time.Date(2021, time.June, 1, 12, 34, 56, 789000000, loc)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "escape",
			fields: fields{},
			args: args{
				tm:     tm,
				layout: "YYYY/[YYYY]",
			},
			want: "2021/YYYY",
		},
		{
			name:   "ISO8601",
			fields: fields{},
			args: args{
				tm:     tm,
				layout: "YYYY-MM-DDTHH:mm:ss.SSSZ",
			},
			want: "2021-06-01T12:34:56.789+09:00",
		},
		{
			name:   "year",
			fields: fields{},
			args: args{
				tm:     tm,
				layout: "YYYYYY/YYYY/YY",
			},
			want: "202121/2021/21",
		},
		{
			name: "month",
			fields: fields{
				monthNames:      defaultMonthNames,
				monthShortNames: defaultMonthShortNames,
			},
			args: args{
				tm:     tm,
				layout: "MMMMMMM/MMMM/MMM/MM/M",
			},
			want: "JuneJun/June/Jun/06/6",
		},
		{
			name: "custom month",
			fields: fields{
				monthNames:      monthNamesJP,
				monthShortNames: monthShortNamesJP,
			},
			args: args{
				tm:     tm,
				layout: "MMMMMMM/MMMM/MMM/MM/M",
			},
			want: "水無月6月/水無月/6月/06/6",
		},
		{
			name:   "day",
			fields: fields{},
			args: args{
				tm:     tm,
				layout: "DDD/DD/D",
			},
			want: "011/01/1",
		},
		{
			name:   "12 hours",
			fields: fields{},
			args: args{
				tm:     tm,
				layout: "HHHhhh/HH/H/hh/h",
			},
			want: "12121212/12/12/12/12",
		},
		{
			name:   "0 hours",
			fields: fields{},
			args: args{
				tm:     time.Date(2021, time.June, 1, 0, 34, 56, 0, loc),
				layout: "HHHhhh/HH/H/hh/h",
			},
			want: "0001212/00/0/12/12",
		},
		{
			name:   "1 hours",
			fields: fields{},
			args: args{
				tm:     time.Date(2021, time.June, 1, 1, 34, 56, 0, loc),
				layout: "HHHhhh/HH/H/hh/h",
			},
			want: "011011/01/1/01/1",
		},
		{
			name:   "minutes",
			fields: fields{},
			args: args{
				tm:     time.Date(2021, time.June, 1, 1, 3, 56, 0, loc),
				layout: "mmm/mm/m",
			},
			want: "033/03/3",
		},
		{
			name:   "seconds",
			fields: fields{},
			args: args{
				tm:     time.Date(2021, time.June, 1, 1, 3, 4, 0, loc),
				layout: "sss/ss/s",
			},
			want: "044/04/4",
		},
		{
			name:   "milliseconds",
			fields: fields{},
			args: args{
				tm:     time.Date(2021, time.June, 1, 1, 3, 4, 789000000, loc),
				layout: "SSS",
			},
			want: "789",
		},
		{
			name:   "AM",
			fields: fields{},
			args: args{
				tm:     time.Date(2021, time.June, 1, 0, 3, 4, 789000000, loc),
				layout: "aA",
			},
			want: "amAM",
		},
		{
			name:   "PM",
			fields: fields{},
			args: args{
				tm:     time.Date(2021, time.June, 1, 12, 3, 4, 789000000, loc),
				layout: "aA",
			},
			want: "pmPM",
		},
		{
			name: "custom AM",
			fields: fields{
				meridiemFunc: meridiemFuncJP,
			},
			args: args{
				tm:     time.Date(2021, time.June, 1, 0, 3, 4, 789000000, loc),
				layout: "aA",
			},
			want: "午前午前",
		},
		{
			name: "custom PM",
			fields: fields{
				meridiemFunc: meridiemFuncJP,
			},
			args: args{
				tm:     time.Date(2021, time.June, 1, 12, 3, 4, 789000000, loc),
				layout: "aA",
			},
			want: "午後午後",
		},
		{
			name:   "timezone plus",
			fields: fields{},
			args: args{
				tm:     tm,
				layout: "ZZZ/Z/ZZ",
			},
			want: "+0900+09:00/+09:00/+0900",
		},
		{
			name: "weekdays",
			fields: fields{
				weekdayFullNames:  defaultWeekdayFullNames,
				weekdayShortNames: defaultWeekdayShortNames,
				weekdayMinNames:   defaultWeekdayMinNames,
			},
			args: args{
				tm:     tm,
				layout: "ddddddd/dddd/ddd/dd/d",
			},
			want: "TuesdayTue/Tuesday/Tue/Tu/2",
		},
		{
			name: "custom weekdays",
			fields: fields{
				weekdayFullNames:  weekdayFullNamesJP,
				weekdayShortNames: weekdayShortNamesJP,
				weekdayMinNames:   weekdayMinNamesJP,
			},
			args: args{
				tm:     tm,
				layout: "ddddddd/dddd/ddd/dd/d",
			},
			want: "火曜日火曜/火曜日/火曜/火/2",
		},
		{
			name:   "timezone minus",
			fields: fields{},
			args: args{
				tm:     time.Date(2021, time.June, 1, 12, 34, 56, 0, time.FixedZone("UTC", -9*60*60)),
				layout: "ZZZ/Z/ZZ",
			},
			want: "-0900-09:00/-09:00/-0900",
		},
		{
			name:   "timezone UTC",
			fields: fields{},
			args: args{
				tm:     time.Date(2021, time.June, 1, 12, 34, 56, 0, time.UTC),
				layout: "ZZZ/Z/ZZ",
			},
			want: "ZZ/Z/Z",
		},
		{
			name:   "timezone UTC",
			fields: fields{},
			args: args{
				tm:     time.Date(2021, time.June, 1, 12, 34, 56, 0, time.UTC),
				layout: "ZZZ/Z/ZZ",
			},
			want: "ZZ/Z/Z",
		},
		{
			name:   "unmatched",
			fields: fields{},
			args: args{
				tm:     tm,
				layout: "Y/YYY",
			},
			want: "Y/YYY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &timeFormatter{
				weekdayFullNames:  tt.fields.weekdayFullNames,
				weekdayShortNames: tt.fields.weekdayShortNames,
				weekdayMinNames:   tt.fields.weekdayMinNames,
				monthNames:        tt.fields.monthNames,
				monthShortNames:   tt.fields.monthShortNames,
				meridiemFunc:      tt.fields.meridiemFunc,
			}
			if got := tr.Format(tt.args.tm, tt.args.layout); got != tt.want {
				t.Errorf("timeFormatter.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeFormatter_MonthNames(t *testing.T) {
	type fields struct {
		monthNames []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "month names",
			fields: fields{
				monthNames: defaultMonthNames,
			},
			want: defaultMonthNames,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &timeFormatter{
				monthNames: tt.fields.monthNames,
			}
			if got := tr.MonthNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeFormatter.MonthNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeFormatter_SetMonthNames(t *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		isErr bool
	}{
		{
			name: "SetMonthNames:valid",
			args: args{
				names: monthNamesJP,
			},
			want: monthNamesJP,
		},
		{
			name: "SetMonthNames:invalid_length",
			args: args{
				names: []string{},
			},
			want:  nil,
			isErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.isErr {
						t.Errorf("timeFormatter.SetMonthNames() panic = %v, want %v", r, tt.isErr)
					}
				}
			}()
			tr := &timeFormatter{}
			tr.SetMonthNames(tt.args.names)
			if got := tr.monthNames; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeFormatter.monthNames = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeFormatter_MonthShortNames(t *testing.T) {
	type fields struct {
		monthShortNames []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "month short names",
			fields: fields{
				monthShortNames: defaultMonthShortNames,
			},
			want: defaultMonthShortNames,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &timeFormatter{
				monthShortNames: tt.fields.monthShortNames,
			}
			if got := tr.MonthShortNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeFormatter.MonthShortNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeFormatter_SetMonthShortNames(t *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		isErr bool
	}{
		{
			name: "SetMonthShortNames:valid",
			args: args{
				names: monthShortNamesJP,
			},
			want: monthShortNamesJP,
		},
		{
			name: "SetMonthShortNames:invalid_length",
			args: args{
				names: []string{},
			},
			want:  nil,
			isErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.isErr {
						t.Errorf("timeFormatter.SetMonthShortNames() panic = %v, want %v", r, tt.isErr)
					}
				}
			}()
			tr := &timeFormatter{}
			tr.SetMonthShortNames(tt.args.names)
			if got := tr.monthShortNames; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeFormatter.monthShortNames = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeFormatter_SetWeekdayFullNames(t *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		isErr bool
	}{
		{
			name: "SetWeekdayFullNames:valid",
			args: args{
				names: weekdayFullNamesJP,
			},
			want: weekdayFullNamesJP,
		},
		{
			name: "SetWeekdayFullNames:invalid_length",
			args: args{
				names: []string{},
			},
			want:  nil,
			isErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.isErr {
						t.Errorf("timeFormatter.SetWeekdayFullNames() panic = %v, want %v", r, tt.isErr)
					}
				}
			}()
			tr := &timeFormatter{}
			tr.SetWeekdayFullNames(tt.args.names)
			if got := tr.weekdayFullNames; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeFormatter.weekdayFullNames = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeFormatter_SetWeekdayMinNames(t *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		isErr bool
	}{
		{
			name: "SetWeekdayMinNames:valid",
			args: args{
				names: weekdayMinNamesJP,
			},
			want: weekdayMinNamesJP,
		},
		{
			name: "SetWeekdayMinNames:invalid_length",
			args: args{
				names: []string{},
			},
			want:  nil,
			isErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.isErr {
						t.Errorf("timeFormatter.SetWeekdayMinNames() panic = %v, want %v", r, tt.isErr)
					}
				}
			}()
			tr := &timeFormatter{}
			tr.SetWeekdayMinNames(tt.args.names)
			if got := tr.weekdayMinNames; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeFormatter.weekdayMinNames = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeFormatter_SetWeekdayShortNames(t *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		isErr bool
	}{
		{
			name: "SetWeekdayShortNames:valid",
			args: args{
				names: weekdayShortNamesJP,
			},
			want: weekdayShortNamesJP,
		},
		{
			name: "SetWeekdayShortNames:invalid_length",
			args: args{
				names: []string{},
			},
			want:  nil,
			isErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.isErr {
						t.Errorf("timeFormatter.SetWeekdayShortNames() panic = %v, want %v", r, tt.isErr)
					}
				}
			}()
			tr := &timeFormatter{}
			tr.SetWeekdayShortNames(tt.args.names)
			if got := tr.weekdayShortNames; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeFormatter.weekdayShortNames = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeFormatter_WeekdayFullNames(t *testing.T) {
	type fields struct {
		weekdayFullNames []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "weekday full names",
			fields: fields{
				weekdayFullNames: defaultWeekdayFullNames,
			},
			want: defaultWeekdayFullNames,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &timeFormatter{
				weekdayFullNames: tt.fields.weekdayFullNames,
			}
			if got := tr.WeekdayFullNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeFormatter.WeekdayFullNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeFormatter_WeekdayMinNames(t *testing.T) {
	type fields struct {
		weekdayMinNames []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "weekday min names",
			fields: fields{
				weekdayMinNames: defaultWeekdayMinNames,
			},
			want: defaultWeekdayMinNames,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &timeFormatter{
				weekdayMinNames: tt.fields.weekdayMinNames,
			}
			if got := tr.WeekdayMinNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeFormatter.WeekdayMinNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeFormatter_WeekdayShortNames(t *testing.T) {
	type fields struct {
		weekdayShortNames []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "weekday short names",
			fields: fields{
				weekdayShortNames: defaultWeekdayShortNames,
			},
			want: defaultWeekdayShortNames,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &timeFormatter{
				weekdayShortNames: tt.fields.weekdayShortNames,
			}
			if got := tr.WeekdayShortNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeFormatter.WeekdayShortNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeFormatter_SetMeridiemFunc(t *testing.T) {
	type args struct {
		f func(int) string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "set meridiem func",
			args: args{
				f: meridiemFuncJP,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &timeFormatter{}
			tr.SetMeridiem(tt.args.f)
			if got := tr.meridiemFunc; reflect.ValueOf(got).Pointer() != reflect.ValueOf(tt.args.f).Pointer() {
				t.Errorf("timeFormatter.meridiemFunc = %p, want %p", got, tt.args.f)
			}
		})
	}
}
