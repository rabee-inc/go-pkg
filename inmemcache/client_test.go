package inmemcache_test

import (
	"errors"
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

func Test_GetMultiOrSet(t *testing.T) {
	type Item struct {
		Value string
	}

	t.Run("キャッシュミスしたキーのみfnで取得し、両方の結果を返す", func(t *testing.T) {
		cache := inmemcache.NewClient[*Item]()

		// 初回: 全てキャッシュミスなので全キーがfnに渡される
		var calledKeys []string
		got, err := cache.GetMultiOrSet([]string{"a", "b"}, func(keys []string) (map[string]*Item, int, error) {
			calledKeys = keys
			res := map[string]*Item{}
			for _, k := range keys {
				res[k] = &Item{Value: "v_" + k}
			}
			return res, 10, nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(calledKeys) != 2 {
			t.Errorf("calledKeys want 2 got %d (%v)", len(calledKeys), calledKeys)
		}
		if len(got) != 2 || got["a"].Value != "v_a" || got["b"].Value != "v_b" {
			t.Errorf("got unexpected: %+v", got)
		}

		// 2回目: "a" はキャッシュ済み、"c" のみミスなので "c" だけfnに渡される
		calledKeys = nil
		got, err = cache.GetMultiOrSet([]string{"a", "c"}, func(keys []string) (map[string]*Item, int, error) {
			calledKeys = keys
			res := map[string]*Item{}
			for _, k := range keys {
				res[k] = &Item{Value: "new_" + k}
			}
			return res, 10, nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(calledKeys) != 1 || calledKeys[0] != "c" {
			t.Errorf("calledKeys want [c] got %v", calledKeys)
		}
		// "a" はキャッシュから取得されるので値は変わらない
		if got["a"].Value != "v_a" {
			t.Errorf(`got["a"] want v_a got %s`, got["a"].Value)
		}
		if got["c"].Value != "new_c" {
			t.Errorf(`got["c"] want new_c got %s`, got["c"].Value)
		}
	})

	t.Run("全てキャッシュヒットの場合fnは呼ばれない", func(t *testing.T) {
		cache := inmemcache.NewClient[*Item]()

		_, err := cache.GetMultiOrSet([]string{"a", "b"}, func(keys []string) (map[string]*Item, int, error) {
			res := map[string]*Item{}
			for _, k := range keys {
				res[k] = &Item{Value: "v_" + k}
			}
			return res, 10, nil
		})
		if err != nil {
			t.Fatal(err)
		}

		called := false
		got, err := cache.GetMultiOrSet([]string{"a", "b"}, func(keys []string) (map[string]*Item, int, error) {
			called = true
			return map[string]*Item{}, 10, nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if called {
			t.Error("fn should not be called when all keys hit cache")
		}
		if got["a"].Value != "v_a" || got["b"].Value != "v_b" {
			t.Errorf("got unexpected: %+v", got)
		}
	})

	t.Run("有効期限切れ後は再度fnで取得する", func(t *testing.T) {
		cache := inmemcache.NewClient[*Item]()

		_, err := cache.GetMultiOrSet([]string{"a"}, func(keys []string) (map[string]*Item, int, error) {
			return map[string]*Item{"a": {Value: "before"}}, 1, nil
		})
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(2 * time.Second)

		called := false
		got, err := cache.GetMultiOrSet([]string{"a"}, func(keys []string) (map[string]*Item, int, error) {
			called = true
			return map[string]*Item{"a": {Value: "after"}}, 10, nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Error("fn should be called after cache expired")
		}
		if got["a"].Value != "after" {
			t.Errorf(`got["a"] want after got %s`, got["a"].Value)
		}
	})

	t.Run("fnがエラーを返す場合はキャッシュ済みの値とエラーを返す", func(t *testing.T) {
		cache := inmemcache.NewClient[*Item]()

		_, err := cache.GetMultiOrSet([]string{"a"}, func(keys []string) (map[string]*Item, int, error) {
			return map[string]*Item{"a": {Value: "v_a"}}, 10, nil
		})
		if err != nil {
			t.Fatal(err)
		}

		wantErr := errors.New("fetch error")
		got, err := cache.GetMultiOrSet([]string{"a", "b"}, func(keys []string) (map[string]*Item, int, error) {
			return nil, 0, wantErr
		})
		if !errors.Is(err, wantErr) {
			t.Errorf("err want %v got %v", wantErr, err)
		}
		// キャッシュヒットした "a" は返る
		if got["a"] == nil || got["a"].Value != "v_a" {
			t.Errorf("got unexpected: %+v", got)
		}
	})
}
