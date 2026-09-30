package validation_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/rabee-inc/go-pkg/errcode"
	"github.com/rabee-inc/go-pkg/validation"
)

func Test_Validation_IsZero(t *testing.T) {
	type args struct {
		val any
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
			name: "正常系: nil",
			args: args{val: nil},
			want: want{result: true},
		},
		{
			name: "正常系: int のゼロ値",
			args: args{val: 0},
			want: want{result: true},
		},
		{
			name: "正常系: int のゼロ値以外",
			args: args{val: 1},
			want: want{result: false},
		},
		{
			name: "正常系: int64 のゼロ値",
			args: args{val: int64(0)},
			want: want{result: true},
		},
		{
			name: "正常系: int64 のゼロ値以外",
			args: args{val: int64(-1)},
			want: want{result: false},
		},
		{
			name: "正常系: float64 のゼロ値",
			args: args{val: float64(0)},
			want: want{result: true},
		},
		{
			name: "正常系: float64 のゼロ値以外",
			args: args{val: 0.1},
			want: want{result: false},
		},
		{
			name: "正常系: string のゼロ値",
			args: args{val: ""},
			want: want{result: true},
		},
		{
			name: "正常系: string のゼロ値以外",
			args: args{val: "a"},
			want: want{result: false},
		},
		{
			name: "正常系: bool のゼロ値",
			args: args{val: false},
			want: want{result: true},
		},
		{
			name: "正常系: bool のゼロ値以外",
			args: args{val: true},
			want: want{result: false},
		},
		{
			name: "正常系: 対象外の型はゼロ値でもfalse",
			args: args{val: []string{}},
			want: want{result: false},
		},
		{
			name: "正常系: 型付きnilはtrueにならない",
			args: args{val: (*int)(nil)},
			want: want{result: false},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := validation.IsZero(tc.args.val)

			if result != tc.want.result {
				t.Errorf("result = %v, want %v", result, tc.want.result)
			}
		})
	}
}

func Test_Validation_MinimumString(t *testing.T) {
	type args struct {
		tag   string
		value string
	}
	type want struct {
		err bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 必要文字数を満たす",
			args: args{tag: "min_str=3", value: "abc"},
			want: want{err: false},
		},
		{
			name: "正常系: 必要文字数を超える",
			args: args{tag: "min_str=3", value: "abcd"},
			want: want{err: false},
		},
		{
			name: "正常系: ゼロ値は許容される",
			args: args{tag: "min_str=3", value: ""},
			want: want{err: false},
		},
		{
			name: "正常系: マルチバイトはルーン数で数える",
			args: args{tag: "min_str=3", value: "あいう"},
			want: want{err: false},
		},
		{
			name: "異常系: 必要文字数に満たない",
			args: args{tag: "min_str=3", value: "ab"},
			want: want{err: true},
		},
		{
			name: "異常系: マルチバイトで必要文字数に満たない",
			args: args{tag: "min_str=3", value: "あい"},
			want: want{err: true},
		},
		{
			name: "異常系: パラメータが数値でない",
			args: args{tag: "min_str=x", value: "abc"},
			want: want{err: true},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			vl := validator.New()
			if err := vl.RegisterValidation("min_str", validation.MinimumString); err != nil {
				t.Fatalf("RegisterValidation err = %v, want nil", err)
			}

			err := vl.Var(tc.args.value, tc.args.tag)

			if (err != nil) != tc.want.err {
				t.Errorf("err = %v, want %v", err, tc.want.err)
			}
		})
	}
}

type convertTarget struct {
	Name  string `validate:"required"`
	Email string `validate:"omitempty,email"`
	Age   int    `validate:"omitempty,gte=20"`
}

func Test_Validation_ConvertErrorMessage(t *testing.T) {
	type args struct {
		target convertTarget
		prefix string
	}
	type want struct {
		lines []string
		code  int
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: バリデーションエラーが整形される",
			args: args{
				target: convertTarget{Name: "", Email: "invalid", Age: 10},
				prefix: "入力エラー",
			},
			want: want{
				lines: []string{"入力エラー", "required/Name/", "email/Email/", "gte/Age/20"},
				code:  http.StatusBadRequest,
			},
		},
		{
			name: "正常系: エラーが1件",
			args: args{
				target: convertTarget{Name: "", Email: "a@example.com", Age: 20},
				prefix: "入力エラー",
			},
			want: want{
				lines: []string{"入力エラー", "required/Name/"},
				code:  http.StatusBadRequest,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			srcErr := validator.New().Struct(tc.args.target)
			if srcErr == nil {
				t.Fatalf("srcErr = nil, want validation error")
			}

			dst := validation.ConvertErrorMessage(srcErr, tc.args.prefix, func(tag, field, value string) string {
				return tag + "/" + field + "/" + value
			})

			lines := strings.Split(dst.Error(), "\n")
			if len(lines) != len(tc.want.lines) {
				t.Fatalf("lines = %v, want %v", lines, tc.want.lines)
			}
			for i, line := range lines {
				if line != tc.want.lines[i] {
					t.Errorf("lines[%v] = %v, want %v", i, line, tc.want.lines[i])
				}
			}
			code, ok := errcode.Get(dst)
			if !ok {
				t.Errorf("errcode.Get ok = false, want true")
			}
			if code != tc.want.code {
				t.Errorf("code = %v, want %v", code, tc.want.code)
			}
		})
	}
}

func Test_Validation_ConvertErrorMessage_NotValidationError(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		sameErr bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "異常系: バリデーションエラー以外はそのまま返す",
			args: args{err: errors.New("other error")},
			want: want{sameErr: true},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			dst := validation.ConvertErrorMessage(tc.args.err, "prefix", func(tag, field, value string) string {
				return tag
			})

			if (dst == tc.args.err) != tc.want.sameErr {
				t.Errorf("dst = %v, want %v", dst, tc.args.err)
			}
		})
	}
}

func Test_Validation_ConvertErrorMessageByDefault(t *testing.T) {
	type args struct {
		target convertTarget
		useFn  bool
	}
	type want struct {
		prefix   string
		contains []string
		code     int
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: フィールド名変換関数ありでデフォルトメッセージが生成される",
			args: args{
				target: convertTarget{Name: "", Email: "invalid", Age: 10},
				useFn:  true,
			},
			want: want{
				prefix:   "以下の入力項目を確認してください",
				contains: []string{"[Name]", "[Email]", "[Age]"},
				code:     http.StatusBadRequest,
			},
		},
		{
			name: "正常系: フィールド名変換関数なし",
			args: args{
				target: convertTarget{Name: "", Email: "a@example.com", Age: 20},
				useFn:  false,
			},
			want: want{
				prefix:   "以下の入力項目を確認してください",
				contains: []string{"Name"},
				code:     http.StatusBadRequest,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			srcErr := validator.New().Struct(tc.args.target)
			if srcErr == nil {
				t.Fatalf("srcErr = nil, want validation error")
			}

			var fn func(field string) string
			if tc.args.useFn {
				fn = func(field string) string { return "[" + field + "]" }
			}
			dst := validation.ConvertErrorMessageByDefault(srcErr, fn)

			lines := strings.Split(dst.Error(), "\n")
			if lines[0] != tc.want.prefix {
				t.Errorf("lines[0] = %v, want %v", lines[0], tc.want.prefix)
			}
			for _, s := range tc.want.contains {
				if !strings.Contains(dst.Error(), s) {
					t.Errorf("dst = %v, want to contain %v", dst.Error(), s)
				}
			}
			code, ok := errcode.Get(dst)
			if !ok {
				t.Errorf("errcode.Get ok = false, want true")
			}
			if code != tc.want.code {
				t.Errorf("code = %v, want %v", code, tc.want.code)
			}
		})
	}
}
