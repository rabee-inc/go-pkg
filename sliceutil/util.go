package sliceutil

import (
	"sort"

	"github.com/rabee-inc/go-pkg/randutil"
)

type Slice[T any] []T

// NewSlice ... constructor
func NewSlice[T any](srcs []T) Slice[T] {
	return srcs
}

// ForEach ... ループ
func ForEach[T any](srcs []T, fn func(src T)) Slice[T] {
	for _, src := range srcs {
		fn(src)
	}
	return srcs
}

// Filter ... 条件に合う要素のみを抽出する
func Filter[T any](srcs []T, fn func(src T) bool) Slice[T] {
	dsts := []T{}
	for _, src := range srcs {
		if fn(src) {
			dsts = append(dsts, src)
		}
	}
	return dsts
}

// FilterWithIndex ... 条件に合う要素のみを抽出する
func FilterWithIndex[T any](srcs []T, fn func(index int) bool) Slice[T] {
	dsts := []T{}
	for i, src := range srcs {
		if fn(i) {
			dsts = append(dsts, src)
		}
	}
	return dsts
}

// Map ... 配列の要素を変換する
func Map[T, E any](srcs []T, fn func(src T) E) Slice[E] {
	dsts := []E{}
	for _, src := range srcs {
		dsts = append(dsts, fn(src))
	}
	return dsts
}

// MapWithIndex ... 配列の要素を変換する
func MapWithIndex[T, E any](srcs []T, fn func(index int) E) Slice[E] {
	dsts := []E{}
	for i := range srcs {
		dsts = append(dsts, fn(i))
	}
	return dsts
}

// Reduce ... reduce
func Reduce[T, E any](srcs []T, fn func(dst E, src T) E) E {
	var dst E
	for _, src := range srcs {
		dst = fn(dst, src)
	}
	return dst
}

// Contains ... 配列に要素が含まれているか
func Contains[T comparable](srcs []T, e T) bool {
	for _, v := range srcs {
		if e == v {
			return true
		}
	}
	return false
}

// ContainsFunc ... 配列に要素が含まれているか
func ContainsFunc[T any](srcs []T, fn func(src T) bool) bool {
	for _, src := range srcs {
		if fn(src) {
			return true
		}
	}
	return false
}

// Some ... ContainsFunc のエイリアス
func Some[T any](srcs []T, fn func(src T) bool) bool {
	return ContainsFunc(srcs, fn)
}

// Shuffle ... 配列をシャッフルする
func Shuffle[T any](srcs []T) Slice[T] {
	n := len(srcs)
	for i := n - 1; i >= 0; i-- {
		j := randutil.Int(0, i)
		srcs[i], srcs[j] = srcs[j], srcs[i]
	}
	return srcs
}

// Sort ... ソート
func Sort[T any](srcs []T, fn func(i, j int) bool) Slice[T] {
	sort.SliceStable(srcs, fn)
	return srcs
}

// Insert ... 配列の任意の場所に挿入する
func Insert[T any](srcs []T, i int, v ...T) Slice[T] {
	return append(srcs[:i], append(v, srcs[i:]...)...)
}

// Delete ... 配列の任意の値をindex指定で削除する
func Delete[T any](srcs []T, i int) Slice[T] {
	return append(srcs[:i], srcs[i+1:]...)
}

// First ... 配列の先頭の要素を取得
func First[T any](srcs []T) T {
	return srcs[0]
}

// Last ... 配列の最後の要素を取得
func Last[T any](srcs []T) T {
	return srcs[len(srcs)-1]
}

// Shift ... 配列の先頭を切り取る (破壊)
func Shift[T any](srcs []T) (T, Slice[T]) {
	return srcs[0], srcs[1:]
}

// Pop ... 配列の後尾を切り取る (破壊)
func Pop[T any](srcs []T) (T, Slice[T]) {
	return srcs[len(srcs)-1], srcs[:len(srcs)-1]
}

// Uniq ... 配列の重複を排除する
func Uniq[T comparable](srcs []T) Slice[T] {
	dsts := make([]T, 0, len(srcs))
	m := make(map[T]bool)
	for _, src := range srcs {
		if _, ok := m[src]; !ok {
			m[src] = true
			dsts = append(dsts, src)
		}
	}
	return dsts
}

// Chunk ... 配列の分割
func Chunk[T any](srcs []T, size int) Slice[Slice[T]] {
	var chunks Slice[Slice[T]]
	srcsSize := len(srcs)
	for i := 0; i < srcsSize; i += size {
		end := i + size
		if srcsSize < end {
			end = srcsSize
		}
		chunks = append(chunks, srcs[i:end])
	}
	return chunks
}

// Excludes ... 配列で指定したbaseの中でtargetに含まれない値の配列を取得
func Excludes[T comparable](base []T, target []T) Slice[T] {
	dsts := []T{}
	for _, b := range base {
		if !Contains(target, b) {
			dsts = append(dsts, b)
		}
	}
	return dsts
}

// --- メソッド ---

// ForEach ... ループ
func (s Slice[T]) ForEach(fn func(src T)) Slice[T] {
	return ForEach(s, fn)
}

// Len ... 要素数を取得する
func (s Slice[T]) Len() int {
	return len(s)
}

// Filter ... 条件に合う要素のみを抽出する
func (s Slice[T]) Filter(fn func(src T) bool) Slice[T] {
	return Filter(s, fn)
}

// FilterWithIndex ... 条件に合う要素のみを抽出する
func (s Slice[T]) FilterWithIndex(fn func(index int) bool) Slice[T] {
	return FilterWithIndex(s, fn)
}

// Shuffle ... 配列をシャッフルする
func (s Slice[T]) Shuffle() Slice[T] {
	return Shuffle(s)
}

// Sort ... ソート
func (s Slice[T]) Sort(fn func(i, j int) bool) Slice[T] {
	return Sort(s, fn)
}

// ContainsFunc ... 配列に要素が含まれているか
func (s Slice[T]) ContainsFunc(fn func(src T) bool) bool {
	return ContainsFunc(s, fn)
}

// Some ... ContainsFunc のエイリアス
func (s Slice[T]) Some(fn func(src T) bool) bool {
	return s.ContainsFunc(fn)
}

// Insert ... 配列の任意の場所に挿入する
func (s Slice[T]) Insert(i int, v ...T) Slice[T] {
	return Insert(s, i, v...)
}

// Delete ... 配列の任意の値をindex指定で削除する
func (s Slice[T]) Delete(i int) Slice[T] {
	return Delete(s, i)
}

// First ... 配列の先頭の要素を取得
func (s Slice[T]) First() T {
	return First(s)
}

// Last ... 配列の最後の要素を取得
func (s Slice[T]) Last() T {
	return Last(s)
}

// Shift ... 配列の先頭を切り取る
func (s Slice[T]) Shift() (T, Slice[T]) {
	return Shift(s)
}

// Pop ... 配列の後尾を切り取る
func (s Slice[T]) Pop() (T, Slice[T]) {
	return Pop(s)
}
