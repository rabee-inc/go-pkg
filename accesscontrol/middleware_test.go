package accesscontrol_test

import (
	"testing"

	"github.com/rabee-inc/go-pkg/accesscontrol"
)

func Test_GetOriginValue(t *testing.T) {
	type args struct {
		origins       []string
		requestOrigin string
	}
	type want struct {
		gotOrigin string
	}
	type testCase struct {
		name string
		args args
		want want
	}

	// テストケース
	tcs := []testCase{
		{
			name: "未指定",
			args: args{
				origins:       []string{},
				requestOrigin: "https://example.com",
			},
			want: want{
				gotOrigin: "*",
			},
		},
		{
			name: "完全一致",
			args: args{
				origins: []string{
					"https://example.com",
				},
				requestOrigin: "https://example.com",
			},
			want: want{
				gotOrigin: "https://example.com",
			},
		},
		{
			name: "ワイルドカード",
			args: args{
				origins: []string{
					"https://aaa-*.example.com",
				},
				requestOrigin: "https://aaa-sub.example.com",
			},
			want: want{
				gotOrigin: "https://aaa-sub.example.com",
			},
		},
		{
			name: "不一致",
			args: args{
				origins: []string{
					"https://aaa-*.example.com",
				},
				requestOrigin: "https://bbb-sub.example.com",
			},
			want: want{
				gotOrigin: "",
			},
		},
	}

	// 実行
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mAccesscontrol := accesscontrol.NewMiddleware(tc.args.origins, []string{})
			gotOrigin := mAccesscontrol.GetOriginValue(tc.args.requestOrigin)
			if gotOrigin != tc.want.gotOrigin {
				t.Errorf("got=%v, want=%v", gotOrigin, tc.want.gotOrigin)
			}
		})
	}
}
