package randutil

import (
	"math/rand"
	"time"
)

const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterIdxMask = 0x3F
)

func seed() {
	rand.Seed(time.Now().UnixNano())
}

// Bool ... 指定確率でbool値を生成する
func Bool(rate float32) bool {
	seed()
	return rand.Float32()*100 < rate
}

// Int ... 指定範囲の乱数を生成する
func Int(min int, max int) int {
	seed()
	return rand.Intn((max+1)-min) + min
}

// String ... nビットのランダムな文字列を生成する
func String(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i := 0; i < n; {
		idx := int(buf[i] & letterIdxMask)
		if idx < len(letters) {
			buf[i] = letters[idx]
			i++
		} else {
			if _, err := rand.Read(buf[i : i+1]); err != nil {
				return "", err
			}
		}
	}
	return string(buf), nil
}

// StringByChar ... nビットのランダムな文字列を生成する
func StringByChar(n int, cr string) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i := 0; i < n; {
		idx := int(buf[i] & letterIdxMask)
		if idx < len(cr) {
			buf[i] = cr[idx]
			i++
		} else {
			if _, err := rand.Read(buf[i : i+1]); err != nil {
				return "", err
			}
		}
	}
	return string(buf), nil
}
