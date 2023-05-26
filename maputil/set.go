package maputil

type empty struct{}
type Set[T comparable] map[T]empty

// NewSet ... slice から Set を作成する
func NewSet[T comparable](s []T) Set[T] {
	set := Set[T]{}
	for _, v := range s {
		set[v] = empty{}
	}
	return set
}

// Len ... 要素数を返す
func (s Set[T]) Len() int {
	return len(s)
}

// Has ... key が存在するかどうか
func (s Set[T]) Has(key T) bool {
	return Has(s, key)
}

// Keys ... キーのslice を返す
func (s Set[T]) Keys() []T {
	return Keys(s)
}

// Add ... キーを追加する
func (s Set[T]) Add(key T) {
	s[key] = empty{}
}

// Delete ... キーを削除する
func (s Set[T]) Delete(key T) {
	delete(s, key)
}

// Clear ... Set を空にする
func (s Set[T]) Clear() {
	Clear(s)
}
