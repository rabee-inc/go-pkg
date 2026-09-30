package randutil

import (
	crand "crypto/rand"
	"math/rand/v2"
)

const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterIdxMask = 0x3F
)

// 指定確率でbool値を生成する
func Bool(rate float32) bool {
	return rand.Float32()*100 < rate
}

// 指定範囲の乱数を生成する
func Int(min int, max int) int {
	return rand.IntN((max+1)-min) + min
}

// nビットのランダムな文字列を生成する
func String(n int) (string, error) {
	return StringByChar(n, letters)
}

// nビットのランダムな文字列を生成する
func StringByChar(n int, cr string) (string, error) {
	buf := make([]byte, n)
	if _, err := crand.Read(buf); err != nil {
		return "", err
	}
	for i := 0; i < n; {
		idx := int(buf[i] & letterIdxMask)
		if idx < len(cr) {
			buf[i] = cr[idx]
			i++
		} else {
			if _, err := crand.Read(buf[i : i+1]); err != nil {
				return "", err
			}
		}
	}
	return string(buf), nil
}
