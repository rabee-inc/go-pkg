package maputil

type OrderedSet[T comparable] map[T]int

// NewOrderedSet ... slice から orderedSet を作成する。orderedSet は Set と違い Keys() の戻り値が挿入した順になります。そのため、uniqなsliceのような扱い方もできます。
func NewOrderedSet[T comparable](s []T) OrderedSet[T] {
	set := OrderedSet[T]{}
	for _, v := range s {
		set.Add(v)
	}
	return set
}

// Len ... 要素数を返す
func (s OrderedSet[T]) Len() int {
	return len(s)
}

// Has ... key が存在するかどうか
func (s OrderedSet[T]) Has(key T) bool {
	return Has(s, key)
}

// Keys ... キーのslice を返す
func (s OrderedSet[T]) Keys() []T {
	keys := make([]T, len(s))
	for k, v := range s {
		keys[v] = k
	}
	return keys
}

// Add ... キーを追加する
func (s OrderedSet[T]) Add(key T) {
	if !s.Has(key) {
		s[key] = len(s)
	}
}

// Delete ... キーを削除する
func (s OrderedSet[T]) Delete(key T) {
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
func (s OrderedSet[T]) Clear() {
	Clear(s)
}
