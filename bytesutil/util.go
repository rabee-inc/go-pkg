package bytesutil

import "unsafe"

// バイト列を文字列に変換する
// NOTE: 変換後に元のバイト列を書き換えてはいけない(メモリを共有している)
func ToStr(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}
