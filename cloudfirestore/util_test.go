// nolint
package cloudfirestore_test

import (
	"reflect"
	"testing"

	"github.com/rabee-inc/go-pkg/cloudfirestore"
)

func Test_SetEmptyByMap(t *testing.T) {
	type source struct {
		Str string
		Map map[string]*struct{}
	}

	type args struct {
		source *source
	}
	type want struct {
		isNil bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	testCases := []testCase{
		{
			name: "補完される",
			args: args{
				source: &source{
					Str: "test",
					Map: nil,
				},
			},
			want: want{
				isNil: false,
			},
		},
		{
			name: "変化なし",
			args: args{
				source: &source{
					Str: "test",
					Map: map[string]*struct{}{},
				},
			},
			want: want{
				isNil: false,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cloudfirestore.SetEmptyByMap(testCase.args.source)
			if testCase.args.source.Map == nil {
				if !testCase.want.isNil {
					t.Errorf("testCase.args.source.Map is nil")
				}
			} else {
				if testCase.want.isNil {
					t.Errorf("testCase.args.source.Map is not nil")
				}
			}
		})
	}
}

func Test_SetEmptyByMaps(t *testing.T) {
	type source struct {
		Str string
		Map map[string]*struct{}
	}

	type args struct {
		sources []*source
	}
	type want struct {
		isNil bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	testCases := []testCase{
		{
			name: "補完される",
			args: args{
				sources: []*source{
					{
						Str: "test",
						Map: nil,
					},
				},
			},
			want: want{
				isNil: false,
			},
		},
		{
			name: "変化なし",
			args: args{
				sources: []*source{
					{
						Str: "test",
						Map: map[string]*struct{}{},
					},
				},
			},
			want: want{
				isNil: false,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dsts := []*source{}
			rv := reflect.Indirect(reflect.ValueOf(dsts))
			rrt := rv.Type().Elem().Elem()
			for _, source := range testCase.args.sources {
				rrv := reflect.ValueOf(source)
				cloudfirestore.SetEmptyByMaps(rrv, rrt)
				if source.Map == nil {
					if !testCase.want.isNil {
						t.Errorf("source.Map is nil")
					}
				} else {
					if testCase.want.isNil {
						t.Errorf("source.Map is not nil")
					}
				}
			}
		})
	}
}
