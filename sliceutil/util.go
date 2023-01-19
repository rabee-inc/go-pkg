package sliceutil

import "github.com/rabee-inc/go-pkg/randutil"

// 配列をシャッフルする
func Shuffle[T any](srcs []T) []T {
	n := len(srcs)
	for i := n - 1; i >= 0; i-- {
		j := randutil.Int(0, i)
		srcs[i], srcs[j] = srcs[j], srcs[i]
	}
	return srcs
}

// 配列の任意の場所に挿入する
func Insert[T any](srcs []T, v T, i int) []T {
	return append(srcs[:i], append([]T{v}, srcs[i:]...)...)
}

// 配列の任意の値を削除する
func Delete[T any](srcs []T, i int) []T {
	return append(srcs[:i], srcs[i+1:]...)
}

// 配列の先頭を切り取る
func Shift[T any](srcs []T) (T, []T) {
	return srcs[0], srcs[1:]
}

// 配列の後尾を切り取る
func Back[T any](srcs []T) (T, []T) {
	return srcs[len(srcs)-1], srcs[:len(srcs)-1]
}

// 配列の重複を排除する
func Uniq[T comparable](srcs []T) []T {
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

// 配列の値の存在確認
func Contains[T comparable](srcs []T, e T) bool {
	for _, v := range srcs {
		if e == v {
			return true
		}
	}
	return false
}

// 配列の分割
func Chunk[T any](srcs []T, size int) [][]T {
	var chunks [][]T
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

// 配列で指定したbaseの中でtargetに含まれない値の配列を取得
func Excludes[T comparable](base []T, target []T) []T {
	dsts := []T{}
	for _, b := range base {
		if !Contains(target, b) {
			dsts = append(dsts, b)
		}
	}
	return dsts
}
