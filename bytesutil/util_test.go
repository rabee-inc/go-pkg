package bytesutil_test

import (
	"testing"

	"github.com/rabee-inc/go-pkg/bytesutil"
)

func Test_Bytesutil_ToStr(t *testing.T) {
	type args struct {
		bytes []byte
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
			name: "正常系: ASCII",
			args: args{bytes: []byte("hello")},
			want: want{result: "hello"},
		},
		{
			name: "正常系: マルチバイト",
			args: args{bytes: []byte("こんにちは")},
			want: want{result: "こんにちは"},
		},
		{
			name: "正常系: 空のバイト列",
			args: args{bytes: []byte{}},
			want: want{result: ""},
		},
		{
			name: "正常系: nil",
			args: args{bytes: nil},
			want: want{result: ""},
		},
		{
			name: "正常系: NULL文字を含む",
			args: args{bytes: []byte{'a', 0x00, 'b'}},
			want: want{result: "a\x00b"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := bytesutil.ToStr(tc.args.bytes)

			if result != tc.want.result {
				t.Errorf("result = %v, want %v", result, tc.want.result)
			}
			if len(result) != len(tc.args.bytes) {
				t.Errorf("len(result) = %v, want %v", len(result), len(tc.args.bytes))
			}
		})
	}
}
