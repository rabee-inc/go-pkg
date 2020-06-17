package bytesutil

import "unsafe"

// ToStr ... バイト列を文字列に変換する
func ToStr(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
