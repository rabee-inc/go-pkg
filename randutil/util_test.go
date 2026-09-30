package randutil_test

import (
	"strings"
	"testing"

	"github.com/rabee-inc/go-pkg/randutil"
)

const defaultLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func Test_Randutil_Bool(t *testing.T) {
	type args struct {
		rate  float32
		count int
	}
	type want struct {
		allTrue  bool
		allFalse bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 確率100%なら常にtrue",
			args: args{rate: 100, count: 1000},
			want: want{allTrue: true, allFalse: false},
		},
		{
			name: "正常系: 確率0%なら常にfalse",
			args: args{rate: 0, count: 1000},
			want: want{allTrue: false, allFalse: true},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			trueCount := 0
			for range tc.args.count {
				if randutil.Bool(tc.args.rate) {
					trueCount++
				}
			}

			if tc.want.allTrue && trueCount != tc.args.count {
				t.Errorf("trueCount = %v, want %v", trueCount, tc.args.count)
			}
			if tc.want.allFalse && trueCount != 0 {
				t.Errorf("trueCount = %v, want 0", trueCount)
			}
		})
	}
}

func Test_Randutil_Int(t *testing.T) {
	type args struct {
		min   int
		max   int
		count int
	}
	type testCase struct {
		name string
		args args
	}

	tcs := []testCase{
		{
			name: "正常系: 範囲内の値が返る",
			args: args{min: 1, max: 10, count: 1000},
		},
		{
			name: "正常系: minとmaxが同じ",
			args: args{min: 5, max: 5, count: 100},
		},
		{
			name: "正常系: 負の範囲",
			args: args{min: -10, max: -1, count: 1000},
		},
		{
			name: "正常系: 0を跨ぐ範囲",
			args: args{min: -5, max: 5, count: 1000},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			for range tc.args.count {
				result := randutil.Int(tc.args.min, tc.args.max)

				if result < tc.args.min || tc.args.max < result {
					t.Fatalf("result = %v, want between %v and %v", result, tc.args.min, tc.args.max)
				}
			}
		})
	}
}

func Test_Randutil_Int_Boundary(t *testing.T) {
	type args struct {
		min   int
		max   int
		count int
	}
	type want struct {
		hitMin bool
		hitMax bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: min と max の両端が返りうる",
			args: args{min: 0, max: 1, count: 1000},
			want: want{hitMin: true, hitMax: true},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			hitMin, hitMax := false, false
			for range tc.args.count {
				switch randutil.Int(tc.args.min, tc.args.max) {
				case tc.args.min:
					hitMin = true
				case tc.args.max:
					hitMax = true
				}
			}

			if hitMin != tc.want.hitMin {
				t.Errorf("hitMin = %v, want %v", hitMin, tc.want.hitMin)
			}
			if hitMax != tc.want.hitMax {
				t.Errorf("hitMax = %v, want %v", hitMax, tc.want.hitMax)
			}
		})
	}
}

func Test_Randutil_String(t *testing.T) {
	type args struct {
		n int
	}
	type want struct {
		length int
		err    bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 指定した長さの文字列が返る",
			args: args{n: 32},
			want: want{length: 32, err: false},
		},
		{
			name: "正常系: 1文字",
			args: args{n: 1},
			want: want{length: 1, err: false},
		},
		{
			name: "正常系: 0文字",
			args: args{n: 0},
			want: want{length: 0, err: false},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := randutil.String(tc.args.n)

			if len(result) != tc.want.length {
				t.Errorf("len(result) = %v, want %v", len(result), tc.want.length)
			}
			if (err != nil) != tc.want.err {
				t.Errorf("err = %v, want %v", err, tc.want.err)
			}
			for _, r := range result {
				if !strings.ContainsRune(defaultLetters, r) {
					t.Errorf("result = %v, want only %v", result, defaultLetters)
					break
				}
			}
		})
	}
}

func Test_Randutil_String_Unique(t *testing.T) {
	type args struct {
		n     int
		count int
	}
	type testCase struct {
		name string
		args args
	}

	tcs := []testCase{
		{
			name: "正常系: 呼び出しごとに異なる文字列が返る",
			args: args{n: 32, count: 1000},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			results := map[string]struct{}{}
			for range tc.args.count {
				result, err := randutil.String(tc.args.n)
				if err != nil {
					t.Fatalf("err = %v, want nil", err)
				}
				if _, ok := results[result]; ok {
					t.Fatalf("result = %v, want unique", result)
				}
				results[result] = struct{}{}
			}

			if len(results) != tc.args.count {
				t.Errorf("len(results) = %v, want %v", len(results), tc.args.count)
			}
		})
	}
}

func Test_Randutil_StringByChar(t *testing.T) {
	type args struct {
		n  int
		cr string
	}
	type want struct {
		length int
		err    bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 指定した文字集合のみが使われる",
			args: args{n: 64, cr: "0123456789"},
			want: want{length: 64, err: false},
		},
		{
			name: "正常系: 1種類の文字集合",
			args: args{n: 16, cr: "a"},
			want: want{length: 16, err: false},
		},
		{
			name: "正常系: 64文字ちょうどの文字集合(マスク上限)",
			args: args{n: 64, cr: defaultLetters + "@#"},
			want: want{length: 64, err: false},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := randutil.StringByChar(tc.args.n, tc.args.cr)

			if len(result) != tc.want.length {
				t.Errorf("len(result) = %v, want %v", len(result), tc.want.length)
			}
			if (err != nil) != tc.want.err {
				t.Errorf("err = %v, want %v", err, tc.want.err)
			}
			for _, r := range result {
				if !strings.ContainsRune(tc.args.cr, r) {
					t.Errorf("result = %v, want only %v", result, tc.args.cr)
					break
				}
			}
		})
	}
}
