package stringutil_test

import (
	"math"
	"testing"

	"github.com/rabee-inc/go-pkg/stringutil"
)

func Test_Stringutil_ToBytes(t *testing.T) {
	type args struct {
		str string
	}
	type want struct {
		result string
		length int
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: ASCII",
			args: args{str: "hello"},
			want: want{result: "hello", length: 5},
		},
		{
			name: "正常系: マルチバイト",
			args: args{str: "こんにちは"},
			want: want{result: "こんにちは", length: 15},
		},
		{
			name: "正常系: 空文字",
			args: args{str: ""},
			want: want{result: "", length: 0},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := stringutil.ToBytes(tc.args.str)

			if string(result) != tc.want.result {
				t.Errorf("result = %v, want %v", string(result), tc.want.result)
			}
			if len(result) != tc.want.length {
				t.Errorf("len(result) = %v, want %v", len(result), tc.want.length)
			}
			// cap が len と一致していないと、書き込み時に元の文字列の外へはみ出す
			if cap(result) != tc.want.length {
				t.Errorf("cap(result) = %v, want %v", cap(result), tc.want.length)
			}
		})
	}
}

func Test_Stringutil_UniqueID(t *testing.T) {
	type want struct {
		length int
		count  int
	}
	type testCase struct {
		name string
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 20文字のユニークなIDが生成される",
			want: want{length: 20, count: 1000},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ids := map[string]struct{}{}
			for range tc.want.count {
				id := stringutil.UniqueID()
				if len(id) != tc.want.length {
					t.Fatalf("len(id) = %v, want %v", len(id), tc.want.length)
				}
				if _, ok := ids[id]; ok {
					t.Fatalf("id = %v, want unique", id)
				}
				ids[id] = struct{}{}
			}
			if len(ids) != tc.want.count {
				t.Errorf("len(ids) = %v, want %v", len(ids), tc.want.count)
			}
		})
	}
}

func Test_Stringutil_IsNumeric(t *testing.T) {
	type args struct {
		s string
	}
	type want struct {
		result bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 整数",
			args: args{s: "123"},
			want: want{result: true},
		},
		{
			name: "正常系: 負の小数",
			args: args{s: "-1.5"},
			want: want{result: true},
		},
		{
			name: "正常系: 指数表記",
			args: args{s: "1e3"},
			want: want{result: true},
		},
		{
			name: "異常系: 数字以外",
			args: args{s: "abc"},
			want: want{result: false},
		},
		{
			name: "異常系: 空文字",
			args: args{s: ""},
			want: want{result: false},
		},
		{
			name: "異常系: カンマ区切り",
			args: args{s: "1,000"},
			want: want{result: false},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := stringutil.IsNumeric(tc.args.s)

			if result != tc.want.result {
				t.Errorf("result = %v, want %v", result, tc.want.result)
			}
		})
	}
}

func Test_Stringutil_ToComma(t *testing.T) {
	type args struct {
		v int64
	}
	type want struct {
		result string
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 0",
			args: args{v: 0},
			want: want{result: "0"},
		},
		{
			name: "正常系: 3桁",
			args: args{v: 999},
			want: want{result: "999"},
		},
		{
			name: "正常系: 4桁",
			args: args{v: 1000},
			want: want{result: "1,000"},
		},
		{
			name: "正常系: 7桁",
			args: args{v: 1234567},
			want: want{result: "1,234,567"},
		},
		{
			name: "正常系: 負の数",
			args: args{v: -1234567},
			want: want{result: "-1,234,567"},
		},
		{
			name: "正常系: 0埋めが必要な桁",
			args: args{v: 1000001},
			want: want{result: "1,000,001"},
		},
		{
			name: "正常系: MaxInt64",
			args: args{v: math.MaxInt64},
			want: want{result: "9,223,372,036,854,775,807"},
		},
		{
			name: "正常系: MinInt64",
			args: args{v: math.MinInt64},
			want: want{result: "-9,223,372,036,854,775,808"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := stringutil.ToComma(tc.args.v)

			if result != tc.want.result {
				t.Errorf("result = %v, want %v", result, tc.want.result)
			}
		})
	}
}

func Test_Stringutil_ToCommaf(t *testing.T) {
	type args struct {
		v float64
	}
	type want struct {
		result string
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 0",
			args: args{v: 0},
			want: want{result: "0"},
		},
		{
			name: "正常系: 4桁",
			args: args{v: 1000},
			want: want{result: "1,000"},
		},
		{
			name: "正常系: 小数あり",
			args: args{v: 1234567.89},
			want: want{result: "1,234,567.89"},
		},
		{
			name: "正常系: 負の数",
			args: args{v: -1234.5},
			want: want{result: "-1,234.5"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := stringutil.ToCommaf(tc.args.v)

			if result != tc.want.result {
				t.Errorf("result = %v, want %v", result, tc.want.result)
			}
		})
	}
}

func Test_Stringutil_ReplaceLast(t *testing.T) {
	type args struct {
		s   string
		old string
		new string
	}
	type want struct {
		result string
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 最後の一致のみ置換される",
			args: args{s: "a/b/c", old: "/", new: "-"},
			want: want{result: "a/b-c"},
		},
		{
			name: "正常系: 一致が1つ",
			args: args{s: "hello world", old: "world", new: "go"},
			want: want{result: "hello go"},
		},
		{
			name: "正常系: 空文字に置換",
			args: args{s: "foo.bar.baz", old: ".baz", new: ""},
			want: want{result: "foo.bar"},
		},
		{
			name: "異常系: 一致しない場合は元の文字列を返す",
			args: args{s: "hello", old: "xyz", new: "-"},
			want: want{result: "hello"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := stringutil.ReplaceLast(tc.args.s, tc.args.old, tc.args.new)

			if result != tc.want.result {
				t.Errorf("result = %v, want %v", result, tc.want.result)
			}
		})
	}
}

func Test_Stringutil_TrimBOM(t *testing.T) {
	type args struct {
		str string
	}
	type want struct {
		result string
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 先頭のBOMが除去される",
			args: args{str: "\xef\xbb\xbfhello"},
			want: want{result: "hello"},
		},
		{
			name: "正常系: BOMが無い場合はそのまま",
			args: args{str: "hello"},
			want: want{result: "hello"},
		},
		{
			name: "正常系: 先頭以外のBOMは除去されない",
			args: args{str: "hello\xef\xbb\xbf"},
			want: want{result: "hello\xef\xbb\xbf"},
		},
		{
			name: "正常系: 空文字",
			args: args{str: ""},
			want: want{result: ""},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := stringutil.TrimBOM(tc.args.str)

			if result != tc.want.result {
				t.Errorf("result = %v, want %v", result, tc.want.result)
			}
		})
	}
}
