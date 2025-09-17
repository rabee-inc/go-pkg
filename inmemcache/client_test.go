package inmemcache_test

import (
	"testing"
	"time"

	"github.com/rabee-inc/go-pkg/inmemcache"
)

func Test_GetOrSet(t *testing.T) {
	type args struct {
		key          string
		expireSecond int
		waitSecond   time.Duration
	}
	type want struct {
		gotValues []string
	}
	type testCase struct {
		name string
		args args
		want want
	}

	// 準備
	type Item struct {
		Value string
	}
	beforeValue := "before_value"
	afterValue := "after_value"

	// テストケース
	tcs := []testCase{
		{
			name: "全てキャッシュから読み込む",
			args: args{
				key:          "key",
				expireSecond: 10,
				waitSecond:   1,
			},
			want: want{
				gotValues: []string{beforeValue, beforeValue, beforeValue, beforeValue},
			},
		},
		{
			name: "キャッシュ有効期限切れ後、新たな値を読み込んでキャッシュする",
			args: args{
				key:          "key",
				expireSecond: 1,
				waitSecond:   2,
			},
			want: want{
				gotValues: []string{beforeValue, beforeValue, afterValue, afterValue},
			},
		},
	}

	// 実行
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			beforeValueFunc := func() (*Item, int, error) {
				return &Item{Value: beforeValue}, tc.args.expireSecond, nil
			}
			afterValueFunc := func() (*Item, int, error) {
				return &Item{Value: afterValue}, tc.args.expireSecond, nil
			}

			cache := inmemcache.NewClient[*Item]()

			// 初回読み込み(オリジナルを取得)
			item, err := cache.GetOrSet(tc.args.key, beforeValueFunc)
			if err != nil {
				t.Fatal(err)
			}
			if item == nil {
				t.Fatal("item is nil")
			}
			if item.Value != tc.want.gotValues[0] {
				t.Errorf("item.Value 0 want %s got %s", tc.want.gotValues[0], item.Value)
			}

			// 2回目読み込み(after valueを設定するが、キャッシュを取得するので before valueが取得される)
			item, err = cache.GetOrSet(tc.args.key, afterValueFunc)
			if err != nil {
				t.Fatal(err)
			}
			if item == nil {
				t.Fatal("item is nil")
			}
			if item.Value != tc.want.gotValues[1] {
				t.Errorf("item.Value 1 want %s got %s", tc.want.gotValues[1], item.Value)
			}

			// 待機
			time.Sleep(tc.args.waitSecond * time.Second)

			// 3回目読み込み(test caseによる)
			item, err = cache.GetOrSet(tc.args.key, afterValueFunc)
			if err != nil {
				t.Fatal(err)
			}
			if item == nil {
				t.Fatal("item is nil")
			}
			if item.Value != tc.want.gotValues[2] {
				t.Errorf("item.Value 2 want %s got %s", tc.want.gotValues[2], item.Value)
			}

			// 4回目読み込み(新しい値になっている)
			item, err = cache.GetOrSet(tc.args.key, afterValueFunc)
			if err != nil {
				t.Fatal(err)
			}
			if item == nil {
				t.Fatal("item is nil")
			}
			if item.Value != tc.want.gotValues[3] {
				t.Errorf("item.Value 3 want %s got %s", tc.want.gotValues[3], item.Value)
			}
		})
	}
}
