package parameter_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/rabee-inc/go-pkg/parameter"
)

func newGetRequest(rawQuery string) *http.Request {
	return httptest.NewRequest(http.MethodGet, "/?"+rawQuery, nil)
}

func Test_Parameter_GetForm(t *testing.T) {
	type args struct {
		rawQuery string
		key      string
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
			name: "正常系: 値が取得できる",
			args: args{rawQuery: "name=taro", key: "name"},
			want: want{result: "taro"},
		},
		{
			name: "正常系: URLエスケープされた値",
			args: args{rawQuery: "name=%E5%A4%AA%E9%83%8E", key: "name"},
			want: want{result: "太郎"},
		},
		{
			name: "正常系: キーが存在しない場合は空文字",
			args: args{rawQuery: "name=taro", key: "age"},
			want: want{result: ""},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := parameter.GetForm(newGetRequest(tc.args.rawQuery), tc.args.key)

			if result != tc.want.result {
				t.Errorf("result = %v, want %v", result, tc.want.result)
			}
		})
	}
}

func Test_Parameter_GetFormByInt(t *testing.T) {
	type args struct {
		rawQuery string
		key      string
	}
	type want struct {
		result int
		err    bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 整数",
			args: args{rawQuery: "age=20", key: "age"},
			want: want{result: 20, err: false},
		},
		{
			name: "正常系: 負の整数",
			args: args{rawQuery: "age=-1", key: "age"},
			want: want{result: -1, err: false},
		},
		{
			name: "正常系: 未指定は0",
			args: args{rawQuery: "", key: "age"},
			want: want{result: 0, err: false},
		},
		{
			name: "異常系: 整数以外",
			args: args{rawQuery: "age=abc", key: "age"},
			want: want{result: 0, err: true},
		},
	}

	ctx := context.Background()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parameter.GetFormByInt(ctx, newGetRequest(tc.args.rawQuery), tc.args.key)

			if result != tc.want.result {
				t.Errorf("result = %v, want %v", result, tc.want.result)
			}
			if (err != nil) != tc.want.err {
				t.Errorf("err = %v, want %v", err, tc.want.err)
			}
		})
	}
}

func Test_Parameter_GetFormByBool(t *testing.T) {
	type args struct {
		rawQuery string
		key      string
	}
	type want struct {
		result bool
		err    bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: true",
			args: args{rawQuery: "flag=true", key: "flag"},
			want: want{result: true, err: false},
		},
		{
			name: "正常系: 1はtrue",
			args: args{rawQuery: "flag=1", key: "flag"},
			want: want{result: true, err: false},
		},
		{
			name: "正常系: false",
			args: args{rawQuery: "flag=false", key: "flag"},
			want: want{result: false, err: false},
		},
		{
			name: "正常系: 未指定はfalse",
			args: args{rawQuery: "", key: "flag"},
			want: want{result: false, err: false},
		},
		{
			name: "異常系: bool以外",
			args: args{rawQuery: "flag=abc", key: "flag"},
			want: want{result: false, err: true},
		},
	}

	ctx := context.Background()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parameter.GetFormByBool(ctx, newGetRequest(tc.args.rawQuery), tc.args.key)

			if result != tc.want.result {
				t.Errorf("result = %v, want %v", result, tc.want.result)
			}
			if (err != nil) != tc.want.err {
				t.Errorf("err = %v, want %v", err, tc.want.err)
			}
		})
	}
}

func Test_Parameter_GetFormBySlice(t *testing.T) {
	type args struct {
		rawQuery string
		key      string
	}
	type want struct {
		result []string
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 複数の値が取得できる",
			args: args{rawQuery: "ids%5B%5D=a&ids%5B%5D=b", key: "ids"},
			want: want{result: []string{"a", "b"}},
		},
		{
			name: "正常系: URLエスケープされた値",
			args: args{rawQuery: "ids%5B%5D=%E5%A4%AA%E9%83%8E", key: "ids"},
			want: want{result: []string{"太郎"}},
		},
		{
			name: "正常系: []なしのキーは対象外",
			args: args{rawQuery: "ids=a&ids=b", key: "ids"},
			want: want{result: []string{}},
		},
		{
			name: "正常系: 該当なしは空スライス",
			args: args{rawQuery: "names%5B%5D=a", key: "ids"},
			want: want{result: []string{}},
		},
	}

	ctx := context.Background()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := parameter.GetFormBySlice(ctx, newGetRequest(tc.args.rawQuery), tc.args.key)

			if !slices.Equal(result, tc.want.result) {
				t.Errorf("result = %v, want %v", result, tc.want.result)
			}
		})
	}
}

func Test_Parameter_GetFormByIntSlice(t *testing.T) {
	type args struct {
		rawQuery string
		key      string
	}
	type want struct {
		result []int
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 複数の整数が取得できる",
			args: args{rawQuery: "ids%5B%5D=1&ids%5B%5D=2", key: "ids"},
			want: want{result: []int{1, 2}},
		},
		{
			name: "正常系: 整数に変換できない値はスキップされる",
			args: args{rawQuery: "ids%5B%5D=1&ids%5B%5D=abc&ids%5B%5D=3", key: "ids"},
			want: want{result: []int{1, 3}},
		},
		{
			name: "正常系: 該当なしは空スライス",
			args: args{rawQuery: "", key: "ids"},
			want: want{result: []int{}},
		},
	}

	ctx := context.Background()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := parameter.GetFormByIntSlice(ctx, newGetRequest(tc.args.rawQuery), tc.args.key)

			if !slices.Equal(result, tc.want.result) {
				t.Errorf("result = %v, want %v", result, tc.want.result)
			}
		})
	}
}

type formParam struct {
	Name    string   `form:"name"`
	Age     int      `form:"age"`
	Rate    float64  `form:"rate"`
	Enabled bool     `form:"enabled"`
	Tags    []string `form:"tags"`
	IDs     []int    `form:"ids"`
	Ignored string   // form タグなし
}

func Test_Parameter_GetForms(t *testing.T) {
	type args struct {
		rawQuery string
	}
	type want struct {
		param formParam
		err   bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 全ての型が設定される",
			args: args{rawQuery: "name=taro&age=20&rate=1.5&enabled=true&tags%5B%5D=a&tags%5B%5D=b&ids%5B%5D=1&ids%5B%5D=2"},
			want: want{
				param: formParam{
					Name:    "taro",
					Age:     20,
					Rate:    1.5,
					Enabled: true,
					Tags:    []string{"a", "b"},
					IDs:     []int{1, 2},
				},
				err: false,
			},
		},
		{
			name: "正常系: 未指定の項目はゼロ値",
			args: args{rawQuery: "name=taro"},
			want: want{
				param: formParam{Name: "taro"},
				err:   false,
			},
		},
		{
			name: "異常系: 整数に変換できない",
			args: args{rawQuery: "age=abc"},
			want: want{param: formParam{}, err: true},
		},
	}

	ctx := context.Background()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var param formParam
			err := parameter.GetForms(ctx, newGetRequest(tc.args.rawQuery), &param)

			if (err != nil) != tc.want.err {
				t.Errorf("err = %v, want %v", err, tc.want.err)
			}
			if err != nil {
				return
			}
			if param.Name != tc.want.param.Name {
				t.Errorf("Name = %v, want %v", param.Name, tc.want.param.Name)
			}
			if param.Age != tc.want.param.Age {
				t.Errorf("Age = %v, want %v", param.Age, tc.want.param.Age)
			}
			if param.Rate != tc.want.param.Rate {
				t.Errorf("Rate = %v, want %v", param.Rate, tc.want.param.Rate)
			}
			if param.Enabled != tc.want.param.Enabled {
				t.Errorf("Enabled = %v, want %v", param.Enabled, tc.want.param.Enabled)
			}
			if !slices.Equal(param.Tags, tc.want.param.Tags) {
				t.Errorf("Tags = %v, want %v", param.Tags, tc.want.param.Tags)
			}
			if !slices.Equal(param.IDs, tc.want.param.IDs) {
				t.Errorf("IDs = %v, want %v", param.IDs, tc.want.param.IDs)
			}
			if param.Ignored != tc.want.param.Ignored {
				t.Errorf("Ignored = %v, want %v", param.Ignored, tc.want.param.Ignored)
			}
		})
	}
}

func Test_Parameter_GetForms_NotPointer(t *testing.T) {
	type args struct {
		dst any
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
			name: "異常系: ポインタ以外はエラー",
			args: args{dst: formParam{}},
			want: want{err: true},
		},
	}

	ctx := context.Background()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := parameter.GetForms(ctx, newGetRequest(""), tc.args.dst)

			if (err != nil) != tc.want.err {
				t.Errorf("err = %v, want %v", err, tc.want.err)
			}
		})
	}
}

func Test_Parameter_GetJSON(t *testing.T) {
	type jsonParam struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type args struct {
		body string
	}
	type want struct {
		name string
		age  int
		err  bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: JSONがデコードされる",
			args: args{body: `{"name":"taro","age":20}`},
			want: want{name: "taro", age: 20, err: false},
		},
		{
			name: "正常系: 未指定の項目はゼロ値",
			args: args{body: `{"name":"taro"}`},
			want: want{name: "taro", age: 0, err: false},
		},
		{
			name: "異常系: 不正なJSON",
			args: args{body: `{invalid}`},
			want: want{name: "", age: 0, err: true},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.args.body))

			var param jsonParam
			err := parameter.GetJSON(r, &param)

			if (err != nil) != tc.want.err {
				t.Errorf("err = %v, want %v", err, tc.want.err)
			}
			if param.Name != tc.want.name {
				t.Errorf("Name = %v, want %v", param.Name, tc.want.name)
			}
			if param.Age != tc.want.age {
				t.Errorf("Age = %v, want %v", param.Age, tc.want.age)
			}
		})
	}
}
