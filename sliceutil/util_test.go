package sliceutil_test

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/rabee-inc/go-pkg/sliceutil"
	"gopkg.in/go-playground/assert.v1"
)

func main() {

}

func Test(t *testing.T) {

	// ForEach
	t.Run("ForEach", func(t *testing.T) {
		// 関数
		expect := []int{1, 2, 3, 4}
		actual := []int{}
		s := sliceutil.ForEach(expect, func(v int) {
			actual = append(actual, v)
		})
		assertSlice(t, expect, actual)

		// メソッド
		actual = []int{}
		s.ForEach(func(v int) {
			actual = append(actual, v)
		})
		assertSlice(t, expect, actual)
	})

	// Filter
	t.Run("Filter", func(t *testing.T) {
		// 関数
		expect := []int{3, 4}
		actual := sliceutil.Filter([]int{1, 2, 3, 4}, func(v int) bool {
			return v > 2
		})
		assertSlice(t, expect, actual)

		// メソッド
		actual = actual.Filter(func(v int) bool {
			return v > 3
		})
		assertSlice(t, []int{4}, actual)
	})

	// FilterWithIndex
	t.Run("FilterWithIndex", func(t *testing.T) {
		// 関数
		expect := []int{1, 2}
		input := []int{1, 2, 3, 4}
		actual := sliceutil.FilterWithIndex(input, func(i int) bool {
			return input[i] < 3
		})
		assertSlice(t, expect, actual)

		// メソッド
		actual = actual.FilterWithIndex(func(i int) bool {
			return actual[i] < 2
		})
		assertSlice(t, []int{1}, actual)
	})

	// Map
	t.Run("Map", func(t *testing.T) {
		// 関数
		expect := []string{"1", "2", "3", "4"}
		actual := sliceutil.Map([]int{1, 2, 3, 4}, func(v int) string {
			return strconv.Itoa(v)
		})
		assertSlice(t, expect, actual)

		// メソッドチェーン
		actual = actual.Filter(func(v string) bool {
			return v != "2"
		})
		assertSlice(t, []string{"1", "3", "4"}, actual)
	})

	// MapWithIndex
	t.Run("MapWithIndex", func(t *testing.T) {
		// 関数
		expect := []string{"1", "2", "3", "4"}
		input := []int{1, 2, 3, 4}
		actual := sliceutil.MapWithIndex(input, func(i int) string {
			return strconv.Itoa(input[i])
		})
		assertSlice(t, expect, actual)

		// メソッドチェーン
		actual = actual.Filter(func(v string) bool {
			return v != "2"
		})
		assertSlice(t, []string{"1", "3", "4"}, actual)
	})

	// Reduce
	t.Run("Reduce", func(t *testing.T) {
		// 関数
		expect := 10
		actual := sliceutil.Reduce([]int{1, 2, 3, 4}, func(acc, v int) int {
			return acc + v
		})
		assert.Equal(t, expect, actual)
	})

	// Contains
	t.Run("Contains", func(t *testing.T) {
		// 関数
		expect := true
		actual := sliceutil.Contains([]int{1, 2, 3, 4}, 2)
		assert.Equal(t, expect, actual)

		expect = false
		actual = sliceutil.Contains([]int{1, 2, 3, 4}, 5)
		assert.Equal(t, expect, actual)
	})

	// ContainsFunc
	t.Run("ContainsFunc", func(t *testing.T) {
		// 関数
		expect := true
		actual := sliceutil.ContainsFunc([]int{1, 2, 3, 4}, func(v int) bool {
			return v == 2
		})
		assert.Equal(t, expect, actual)

		expect = false
		actual = sliceutil.ContainsFunc([]int{1, 2, 3, 4}, func(v int) bool {
			return v == 5
		})
		assert.Equal(t, expect, actual)

		actual = sliceutil.NewSlice([]int{1, 2, 3, 4}).ContainsFunc(func(v int) bool {
			return v == 2
		})
		assert.Equal(t, true, actual)

	})

	// Shuffle
	t.Run("Shuffle", func(t *testing.T) {
		// 関数
		input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		// 3000 回実行して、各数値のばらつきが誤差 30% 以内なら OK とする
		N := 3000

		count := make(map[int]map[int]int)
		// map 初期化
		for _, v := range input {
			count[v] = make(map[int]int)
			for i := range input {
				count[v][i] = 0
			}
		}

		for i := 0; i < N; i++ {
			// input をコピー
			ip := append([]int{}, input...)
			shuffled := sliceutil.Shuffle(ip)
			for i, v := range shuffled {
				count[v][i]++
			}
		}

		// 平均と誤差が 30% 以内かどうか
		for _, v := range count {
			average := float64(N) / float64(len(v))
			for _, v2 := range v {
				diff := math.Abs(float64(v2) - average)
				if diff > average*0.3 {
					t.Errorf("expected: %v, actual: %v", diff, average*0.3)
				}
			}
		}

		// メソッド
		count = make(map[int]map[int]int)

		// map 初期化
		for _, v := range input {
			count[v] = make(map[int]int)
			for i := range input {
				count[v][i] = 0
			}
		}

		for i := 0; i < N; i++ {
			// input をコピー
			ip := append([]int{}, input...)
			shuffled := sliceutil.NewSlice(ip).Shuffle()
			for i, v := range shuffled {
				count[v][i]++
			}
		}

		// 平均と誤差が 30% 以内かどうか
		for _, v := range count {
			average := float64(N) / float64(len(v))
			for _, v2 := range v {
				diff := math.Abs(float64(v2) - average)
				if diff > average*0.3 {
					t.Errorf("expected: %v, actual: %v", diff, average*0.3)
				}
			}
		}

	})

	// Sort
	t.Run("Sort", func(t *testing.T) {
		// 関数
		expect := []int{1, 2, 3, 4, 10}
		input := []int{10, 4, 3, 2, 1}
		actual := sliceutil.Sort(input, func(i, j int) bool {
			return input[i] < input[j]
		})
		assertSlice(t, expect, actual)

		// メソッド
		expect = []int{1, 2, 3, 4}
		input = []int{4, 3, 2, 1}
		actual = sliceutil.NewSlice(input).Sort(func(i, j int) bool {
			return input[i] < input[j]
		})
		assertSlice(t, expect, actual)
	})

	// Insert
	t.Run("Insert", func(t *testing.T) {
		// 関数
		expect := []int{1, 2, 3, 4, 5}
		input := []int{1, 2, 4, 5}
		actual := sliceutil.Insert(input, 2, 3)
		assertSlice(t, expect, actual)

		// メソッド
		expect = []int{1, 2, 3, 3, 4, 5}
		input = []int{1, 2, 4, 5}
		actual = sliceutil.NewSlice(input).Insert(2, 3, 3)
		assertSlice(t, expect, actual)
	})

	// Delete
	t.Run("Delete", func(t *testing.T) {
		// 関数
		expect := []int{1, 2, 4, 5}
		input := []int{1, 2, 3, 4, 5}
		actual := sliceutil.Delete(input, 2)
		assertSlice(t, expect, actual)

		// メソッド
		expect = []int{1, 2, 4, 5}
		input = []int{1, 2, 3, 4, 5}
		actual = sliceutil.NewSlice(input).Delete(2)
		assertSlice(t, expect, actual)
	})

	// First
	t.Run("First", func(t *testing.T) {
		// 関数
		expect := 1
		input := []int{1, 2, 3, 4, 5}
		actual := sliceutil.First(input)
		assert.Equal(t, expect, actual)
		assertSlice(t, []int{1, 2, 3, 4, 5}, input)

		// メソッド
		expect = 1
		input = []int{1, 2, 3, 4, 5}
		actual = sliceutil.NewSlice(input).First()
		assert.Equal(t, expect, actual)
		assertSlice(t, []int{1, 2, 3, 4, 5}, input)

	})

	// Last
	t.Run("Last", func(t *testing.T) {
		// 関数
		expect := 5
		input := []int{1, 2, 3, 4, 5}
		actual := sliceutil.Last(input)
		assert.Equal(t, expect, actual)
		assertSlice(t, []int{1, 2, 3, 4, 5}, input)

		// メソッド
		expect = 5
		input = []int{1, 2, 3, 4, 5}
		actual = sliceutil.NewSlice(input).Last()
		assert.Equal(t, expect, actual)
		assertSlice(t, []int{1, 2, 3, 4, 5}, input)
	})

	// Shift
	t.Run("Shift", func(t *testing.T) {
		// 関数
		expect := 1
		input := []int{1, 2, 3, 4, 5}
		actual, output := sliceutil.Shift(input)
		assert.Equal(t, expect, actual)
		fmt.Println("input: ", input, "actual: ", actual, "expect: ", expect)
		assertSlice(t, []int{2, 3, 4, 5}, output)

		// メソッド
		expect = 2
		actual, s := output.Shift()
		assert.Equal(t, expect, actual)
		assertSlice(t, []int{3, 4, 5}, s)
	})

	// Pop
	t.Run("Pop", func(t *testing.T) {
		// 関数
		expect := 5
		input := []int{1, 2, 3, 4, 5}
		actual, output := sliceutil.Pop(input)
		assert.Equal(t, expect, actual)
		assertSlice(t, []int{1, 2, 3, 4}, output)

		// メソッド
		expect = 4

		actual, s := output.Pop()
		assert.Equal(t, expect, actual)
		assertSlice(t, []int{1, 2, 3}, s)
	})

	// Chunk
	t.Run("Chunk", func(t *testing.T) {
		// 関数
		expect := [][]int{{1, 2}, {3, 4}, {5}}
		input := []int{1, 2, 3, 4, 5}
		actual := sliceutil.Chunk(input, 2)
		for i, s := range actual {
			assertSlice(t, expect[i], s)
		}

		// メソッド
		expect = [][]int{{1, 2}, {3, 4}, {5}}
		input = []int{1, 2, 3, 4, 5}
		actual = sliceutil.NewSlice(input).Chunk(2)
		for i, s := range actual {
			assertSlice(t, expect[i], s)
		}
	})

}

func eqSlice[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func assertSlice[T comparable](t *testing.T, expected, actual []T) {
	if !eqSlice(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}
