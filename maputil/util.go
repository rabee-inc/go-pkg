package maputil

type Map[T comparable, U any] map[T]U
type JSON map[string]any

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
