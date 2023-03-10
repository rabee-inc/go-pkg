package maputil_test

import (
	"fmt"
	"testing"

	"github.com/rabee-inc/go-pkg/maputil"
	"gopkg.in/go-playground/assert.v1"
)

func Test(t *testing.T) {

	// Values, Has, ValuesByKeys
	t.Run("Values, Has, ValuesByKeys", func(t *testing.T) {
		m1 := map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		}
		expected := []int{1, 2, 3}
		actual := maputil.ValuesByKeys(m1, []string{"one", "two", "three"})
		fmt.Println("actual:", actual, "expected:", expected)
		m2 := maputil.Map[string, int]{
			"one":   1,
			"two":   2,
			"three": 3,
		}
		assert.Equal(t, m2.Has("one"), true)
		assert.Equal(t, m2.Has("two"), true)
		assert.Equal(t, m2.Has("three"), true)
		assert.Equal(t, m2.Has("four"), false)
		m3 := maputil.NewMap(m1)
		m3["five"] = 5
		m3["four"] = 4
		ints := m3.ValuesByKeys([]string{"one", "two", "three"})
		assertSlice(t, ints, []int{1, 2, 3})
		fmt.Println("values: ", m3.Values())
	})

	// Keys
	t.Run("Keys", func(t *testing.T) {
		m := map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		}
		set := maputil.NewSet(maputil.Keys(m))
		assert.Equal(t, set.Has("one"), true)
		assert.Equal(t, set.Has("two"), true)
		assert.Equal(t, set.Has("three"), true)
		assert.Equal(t, set.Len(), 3)
		m2 := maputil.Map[string, int]{
			"one": 1,
		}
		m2["two"] = 2
		m2["three"] = 3
		set = maputil.NewSet(m2.Keys())
		assert.Equal(t, set.Has("one"), true)
		assert.Equal(t, set.Has("two"), true)
		assert.Equal(t, set.Has("three"), true)
		assert.Equal(t, set.Len(), 3)
	})

	// Delete
	t.Run("Delete", func(t *testing.T) {
		m := maputil.Map[string, int]{
			"one": 1,
		}
		m["two"] = 2
		m["three"] = 3
		m.Delete("two")
		assert.Equal(t, len(m), 2)
	})

	// Clear
	t.Run("Clear", func(t *testing.T) {
		m := map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		}
		maputil.Clear(m)
		assert.Equal(t, len(m), 0)

		m2 := maputil.Map[string, int]{
			"one": 1,
		}
		m2.Clear()
		assert.Equal(t, len(m2), 0)
	})

	// --- JSON ---
	t.Run("JSON", func(t *testing.T) {
		// json := map[string]any{}
		json := maputil.JSON{
			"one": 1,
			"two": "2",
		}
		// var m map[string]any = json
		json["three"] = maputil.Map[string, float64]{
			"three_one": 3.1,
		}
		assert.Equal(t, json.ToMap().Has("one"), true)
		assert.Equal(t, json["one"], 1)
		assert.Equal(t, json["two"], "2")
		assertMap(t, json["three"].(maputil.Map[string, float64]), map[string]float64{"three_one": 3.1})
	})

	// --- OrderedSet ---

	// Add, Delete, Keys
	t.Run("Add, Delete, Keys", func(t *testing.T) {
		s := maputil.NewOrderedSet([]int{0, 1, 2, 3})
		assertSlice(t, s.Keys(), []int{0, 1, 2, 3})

		s.Add(4)
		assertSlice(t, s.Keys(), []int{0, 1, 2, 3, 4})

		s.Delete(2)
		assertSlice(t, s.Keys(), []int{0, 1, 3, 4})

		s.Add(2)
		assertSlice(t, s.Keys(), []int{0, 1, 3, 4, 2})

		s.Add(1)
		assertSlice(t, s.Keys(), []int{0, 1, 3, 4, 2})
	})

}

func assertMap[T comparable, U comparable](t *testing.T, expected, actual map[T]U) {
	if !eqMap(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func assertSlice[T comparable](t *testing.T, expected, actual []T) {
	if !eqSlice(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func eqMap[T comparable, U comparable](a, b map[T]U) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if v != b[k] {
			return false
		}
	}
	return true
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
