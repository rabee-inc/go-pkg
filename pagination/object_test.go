package pagination_test

import (
	"testing"

	"github.com/rabee-inc/go-pkg/pagination"
)

func Test_Pagination_New(t *testing.T) {
	type args struct {
		currentPage int
		perPage     int
	}
	type want struct {
		currentPage int
		perPage     int
		isFirstPage bool
		isLastPage  bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 指定した値が設定される",
			args: args{currentPage: 2, perPage: 20},
			want: want{currentPage: 2, perPage: 20, isFirstPage: false, isLastPage: false},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := pagination.New(tc.args.currentPage, tc.args.perPage)

			if result.CurrentPage != tc.want.currentPage {
				t.Errorf("CurrentPage = %v, want %v", result.CurrentPage, tc.want.currentPage)
			}
			if result.PerPage != tc.want.perPage {
				t.Errorf("PerPage = %v, want %v", result.PerPage, tc.want.perPage)
			}
			if result.IsFirstPage != tc.want.isFirstPage {
				t.Errorf("IsFirstPage = %v, want %v", result.IsFirstPage, tc.want.isFirstPage)
			}
			if result.IsLastPage != tc.want.isLastPage {
				t.Errorf("IsLastPage = %v, want %v", result.IsLastPage, tc.want.isLastPage)
			}
		})
	}
}

func Test_Pagination_Object_Set(t *testing.T) {
	type args struct {
		currentPage int
		perPage     int
		totalCount  int
	}
	type want struct {
		totalCount  int
		offsetCount int
		totalPage   int
		nextPage    int
		prevPage    int
		isFirstPage bool
		isLastPage  bool
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 最初のページ",
			args: args{currentPage: 1, perPage: 10, totalCount: 100},
			want: want{
				totalCount:  100,
				offsetCount: 0,
				totalPage:   10,
				nextPage:    2,
				prevPage:    0,
				isFirstPage: true,
				isLastPage:  false,
			},
		},
		{
			name: "正常系: 中間のページ",
			args: args{currentPage: 5, perPage: 10, totalCount: 100},
			want: want{
				totalCount:  100,
				offsetCount: 40,
				totalPage:   10,
				nextPage:    6,
				prevPage:    4,
				isFirstPage: false,
				isLastPage:  false,
			},
		},
		{
			name: "正常系: 最後のページ",
			args: args{currentPage: 10, perPage: 10, totalCount: 100},
			want: want{
				totalCount:  100,
				offsetCount: 90,
				totalPage:   10,
				nextPage:    0,
				prevPage:    9,
				isFirstPage: false,
				isLastPage:  true,
			},
		},
		{
			name: "正常系: 端数があるので総ページ数は切り上げ",
			args: args{currentPage: 1, perPage: 10, totalCount: 101},
			want: want{
				totalCount:  101,
				offsetCount: 0,
				totalPage:   11,
				nextPage:    2,
				prevPage:    0,
				isFirstPage: true,
				isLastPage:  false,
			},
		},
		{
			name: "正常系: 1ページに収まる場合は最初かつ最後のページ",
			args: args{currentPage: 1, perPage: 10, totalCount: 5},
			want: want{
				totalCount:  5,
				offsetCount: 0,
				totalPage:   1,
				nextPage:    0,
				prevPage:    0,
				isFirstPage: true,
				isLastPage:  true,
			},
		},
		{
			name: "正常系: 総件数が0の場合は何も設定されない",
			args: args{currentPage: 1, perPage: 10, totalCount: 0},
			want: want{
				totalCount:  0,
				offsetCount: 0,
				totalPage:   0,
				nextPage:    0,
				prevPage:    0,
				isFirstPage: false,
				isLastPage:  false,
			},
		},
		{
			name: "正常系: 総ページ数を超えるページ指定は最後のページ扱い",
			args: args{currentPage: 20, perPage: 10, totalCount: 100},
			want: want{
				totalCount:  100,
				offsetCount: 190,
				totalPage:   10,
				nextPage:    0,
				prevPage:    19,
				isFirstPage: false,
				isLastPage:  true,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := pagination.New(tc.args.currentPage, tc.args.perPage)
			result.Set(tc.args.totalCount)

			if result.TotalCount != tc.want.totalCount {
				t.Errorf("TotalCount = %v, want %v", result.TotalCount, tc.want.totalCount)
			}
			if result.OffsetCount != tc.want.offsetCount {
				t.Errorf("OffsetCount = %v, want %v", result.OffsetCount, tc.want.offsetCount)
			}
			if result.TotalPage != tc.want.totalPage {
				t.Errorf("TotalPage = %v, want %v", result.TotalPage, tc.want.totalPage)
			}
			if result.NextPage != tc.want.nextPage {
				t.Errorf("NextPage = %v, want %v", result.NextPage, tc.want.nextPage)
			}
			if result.PrevPage != tc.want.prevPage {
				t.Errorf("PrevPage = %v, want %v", result.PrevPage, tc.want.prevPage)
			}
			if result.IsFirstPage != tc.want.isFirstPage {
				t.Errorf("IsFirstPage = %v, want %v", result.IsFirstPage, tc.want.isFirstPage)
			}
			if result.IsLastPage != tc.want.isLastPage {
				t.Errorf("IsLastPage = %v, want %v", result.IsLastPage, tc.want.isLastPage)
			}
		})
	}
}
