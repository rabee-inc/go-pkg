package sliceutil

import "github.com/rabee-inc/go-pkg/randutil"

// Shuffle ... string配列をシャッフルする
func Shuffle(arr []string) []string {
	n := len(arr)
	for i := n - 1; i >= 0; i-- {
		j := randutil.Int(0, i)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// IntShuffle ... int配列をシャッフルする
func IntShuffle(arr []int) []int {
	n := len(arr)
	for i := n - 1; i >= 0; i-- {
		j := randutil.Int(0, i)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// Int64Shuffle ... int64配列をシャッフルする
func Int64Shuffle(arr []int64) []int64 {
	n := len(arr)
	for i := n - 1; i >= 0; i-- {
		j := randutil.Int(0, i)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// Insert ... string配列の任意の場所に挿入する
func Insert(arr []string, v string, i int) []string {
	return append(arr[:i], append([]string{v}, arr[i:]...)...)
}

// IntInsert ... int配列の任意の場所に挿入する
func IntInsert(arr []int, v int, i int) []int {
	return append(arr[:i], append([]int{v}, arr[i:]...)...)
}

// Int64Insert ... int64配列の任意の場所に挿入する
func Int64Insert(arr []int64, v int64, i int) []int64 {
	return append(arr[:i], append([]int64{v}, arr[i:]...)...)
}

// Delete ... string配列の任意の値を削除する
func Delete(arr []string, i int) []string {
	return append(arr[:i], arr[i+1:]...)
}

// IntDelete ... int配列の任意の値を削除する
func IntDelete(arr []int, i int) []int {
	return append(arr[:i], arr[i+1:]...)
}

// Int64Delete ... int64配列の任意の値を削除する
func Int64Delete(arr []int64, i int) []int64 {
	return append(arr[:i], arr[i+1:]...)
}

// Shift ... string配列の先頭を切り取る
func Shift(arr []string) (string, []string) {
	return arr[0], arr[1:]
}

// IntShift ... int配列の先頭を切り取る
func IntShift(arr []int) (int, []int) {
	return arr[0], arr[1:]
}

// Int64Shift ... int64配列の先頭を切り取る
func Int64Shift(arr []int64) (int64, []int64) {
	return arr[0], arr[1:]
}

// Back ... string配列の後尾を切り取る
func Back(arr []string) (string, []string) {
	return arr[len(arr)-1], arr[:len(arr)-1]
}

// IntBack ... int配列の後尾を切り取る
func IntBack(arr []int) (int, []int) {
	return arr[len(arr)-1], arr[:len(arr)-1]
}

// Int64Back ... int64配列の後尾を切り取る
func Int64Back(arr []int64) (int64, []int64) {
	return arr[len(arr)-1], arr[:len(arr)-1]
}

// Uniq ... string配列の重複を排除する
func Uniq(arr []string) []string {
	m := make(map[string]bool)
	uniq := []string{}
	for _, v := range arr {
		if !m[v] {
			m[v] = true
			uniq = append(uniq, v)
		}
	}
	return uniq
}

// IntUniq ... int配列の重複を排除する
func IntUniq(arr []int) []int {
	m := make(map[int]bool)
	uniq := []int{}
	for _, v := range arr {
		if !m[v] {
			m[v] = true
			uniq = append(uniq, v)
		}
	}
	return uniq
}

// Int64Uniq ... int64配列の重複を排除する
func Int64Uniq(arr []int64) []int64 {
	m := make(map[int64]bool)
	uniq := []int64{}
	for _, v := range arr {
		if !m[v] {
			m[v] = true
			uniq = append(uniq, v)
		}
	}
	return uniq
}

// Contains ... string配列の値の存在確認
func Contains(arr []string, e string) bool {
	for _, v := range arr {
		if e == v {
			return true
		}
	}
	return false
}

// IntContains ... int配列の値の存在確認
func IntContains(arr []int, e int) bool {
	for _, v := range arr {
		if e == v {
			return true
		}
	}
	return false
}

// Int64Contains ... int64配列の値の存在確認
func Int64Contains(arr []int64, e int64) bool {
	for _, v := range arr {
		if e == v {
			return true
		}
	}
	return false
}

// Chunk ... string配列の分割
func Chunk(arr []string, size int) [][]string {
	var chunks [][]string
	arrSize := len(arr)
	for i := 0; i < arrSize; i += size {
		end := i + size
		if arrSize < end {
			end = arrSize
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}

// IntChunk ... int配列の分割
func IntChunk(arr []int, size int) [][]int {
	var chunks [][]int
	arrSize := len(arr)
	for i := 0; i < arrSize; i += size {
		end := i + size
		if arrSize < end {
			end = arrSize
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}

// Int64Chunk ... int64配列の分割
func Int64Chunk(arr []int64, size int) [][]int64 {
	var chunks [][]int64
	arrSize := len(arr)
	for i := 0; i < arrSize; i += size {
		end := i + size
		if arrSize < end {
			end = arrSize
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}

// Mix ... string配列を交互に結合する(空文字は無視)
func Mix(arrList [][]string) []string {
	maxLen := 0
	for _, arr := range arrList {
		len := len(arr)
		if maxLen < len {
			maxLen = len
		}
	}
	dst := []string{}
	for i := 0; i < maxLen; i++ {
		for _, arr := range arrList {
			if len(arr) < i+1 || arr[i] == "" {
				continue
			}
			dst = append(dst, arr[i])
		}
	}
	return dst
}

// IntMix ... int配列を交互に結合する
func IntMix(arrList [][]int) []int {
	maxLen := 0
	for _, arr := range arrList {
		len := len(arr)
		if maxLen < len {
			maxLen = len
		}
	}
	dst := []int{}
	for i := 0; i < maxLen; i++ {
		for _, arr := range arrList {
			if len(arr) < i+1 {
				continue
			}
			dst = append(dst, arr[i])
		}
	}
	return dst
}

// Int64Mix ... int64配列を交互に結合する
func Int64Mix(arrList [][]int64) []int64 {
	maxLen := 0
	for _, arr := range arrList {
		len := len(arr)
		if maxLen < len {
			maxLen = len
		}
	}
	dst := []int64{}
	for i := 0; i < maxLen; i++ {
		for _, arr := range arrList {
			if len(arr) < i+1 {
				continue
			}
			dst = append(dst, arr[i])
		}
	}
	return dst
}

// Excludes ... string配列で指定したbaseの中でtargetに含まれない値の配列を取得
func Excludes(base []string, target []string) []string {
	dsts := []string{}
	for _, b := range base {
		if !Contains(target, b) {
			dsts = append(dsts, b)
		}
	}
	return dsts
}

// IntExcludes ... int配列で指定したbaseの中でtargetに含まれない値の配列を取得
func IntExcludes(base []int, target []int) []int {
	dsts := []int{}
	for _, b := range base {
		if !IntContains(target, b) {
			dsts = append(dsts, b)
		}
	}
	return dsts
}

// Int64Excludes ... int64配列で指定したbaseの中でtargetに含まれない値の配列を取得
func Int64Excludes(base []int64, target []int64) []int64 {
	dsts := []int64{}
	for _, b := range base {
		if !Int64Contains(target, b) {
			dsts = append(dsts, b)
		}
	}
	return dsts
}
