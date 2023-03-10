package maputil

type empty struct{}
type set[T comparable] map[T]empty
type orderedSet[T comparable] map[T]int
type Map[T comparable, U any] map[T]U
type JSON map[string]any

// --- Set 関数 ---

// NewSet ... slice から Set を作成する
func NewSet[T comparable](s []T) set[T] {
	set := set[T]{}
	for _, v := range s {
		set[v] = empty{}
	}
	return set
}

// --- OrderedSet 関数 ---

// NewOrderedSet ... slice から orderedSet を作成する。orderedSet は Set と違い Keys() の戻り値が挿入した順になります。そのため、uniqなsliceのような扱い方もできます。
func NewOrderedSet[T comparable](s []T) orderedSet[T] {
	set := orderedSet[T]{}
	for _, v := range s {
		set.Add(v)
	}
	return set
}

// --- Map 関数 ---

// NewMap ... Map を作成する
func NewMap[T comparable, U any](m map[T]U) Map[T, U] {
	return Map[T, U](m)
}

// Has ... key が存在するかどうか
func Has[T comparable, U any](m map[T]U, key T) bool {
	_, ok := m[key]
	return ok
}

// Keys ... キーのslice を返す
func Keys[T comparable, U any](m map[T]U) []T {
	keys := make([]T, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// Values ... 値のslice を返す
func Values[T comparable, U any](m map[T]U) []U {
	values := make([]U, len(m))
	i := 0
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

// ValuesByKeys ... keys slice の順番で value のslice を返す(map に存在するもののみ)
func ValuesByKeys[T comparable, U any](m map[T]U, keys []T) []U {
	values := []U{}
	for _, k := range keys {
		if v, ok := m[k]; ok {
			values = append(values, v)
		}
	}
	return values
}

// Clear ... map を空にする
func Clear[T comparable, U any](m map[T]U) {
	for k := range m {
		delete(m, k)
	}
}

// --- Set メソッド ---

// Len ... 要素数を返す
func (s set[T]) Len() int {
	return len(s)
}

// Has ... key が存在するかどうか
func (s set[T]) Has(key T) bool {
	return Has(s, key)
}

// Keys ... キーのslice を返す
func (s set[T]) Keys() []T {
	return Keys(s)
}

// Add ... キーを追加する
func (s set[T]) Add(key T) {
	s[key] = empty{}
}

// Delete ... キーを削除する
func (s set[T]) Delete(key T) {
	delete(s, key)
}

// Clear ... Set を空にする
func (s set[T]) Clear() {
	Clear(s)
}

// --- OrderedSet メソッド ---

// Len ... 要素数を返す
func (s orderedSet[T]) Len() int {
	return len(s)
}

// Has ... key が存在するかどうか
func (s orderedSet[T]) Has(key T) bool {
	return Has(s, key)
}

// Keys ... キーのslice を返す
func (s orderedSet[T]) Keys() []T {
	keys := make([]T, len(s))
	for k, v := range s {
		keys[v] = k
	}
	return keys
}

// Add ... キーを追加する
func (s orderedSet[T]) Add(key T) {
	if !s.Has(key) {
		s[key] = len(s)
	}
}

// Delete ... キーを削除する
func (s orderedSet[T]) Delete(key T) {
	if i, ok := s[key]; ok {
		for k, v := range s {
			if v > i {
				s[k] = v - 1
			}
		}
		delete(s, key)
	}
}

// Clear ... Set を空にする
func (s orderedSet[T]) Clear() {
	Clear(s)
}

// --- Map メソッド ---

// Len ... 要素数を返す
func (m Map[T, U]) Len() int {
	return len(m)
}

// Has ... key が存在するかどうか
func (m Map[T, U]) Has(key T) bool {
	return Has(m, key)
}

// Keys ... キーのslice を返す
func (m Map[T, U]) Keys() []T {
	return Keys(m)
}

// Values ... 値のslice を返す
func (m Map[T, U]) Values() []U {
	return Values(m)
}

// ValuesByKeys ... keys slice の順番で value のslice を返す(map に存在するもののみ)
func (m Map[T, U]) ValuesByKeys(keys []T) []U {
	return ValuesByKeys(m, keys)
}

// Delete ... キーを削除する
func (m Map[T, U]) Delete(key T) {
	delete(m, key)
}

// Clear ... Map を空にする
func (m Map[T, U]) Clear() {
	Clear(m)
}

// --- JSON メソッド ---

// ToMap ... JSON 型を Map 型に変換する
func (j JSON) ToMap() Map[string, any] {
	return Map[string, any](j)
}
