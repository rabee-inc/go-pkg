package stringutil

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"unsafe"

	"github.com/rs/xid"
)

const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterIdxMask = 0x3F
)

// ToMD5 ... 文字列のハッシュ(MD5)を取得する
func ToMD5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ToSHA256 ... 文字列のハッシュ(SHA256)を取得する
func ToSHA256(str string) string {
	c := sha256.Sum256([]byte(str))
	return hex.EncodeToString(c[:])
}

// ToBytes ... 文字列をバイト列に変換する
func ToBytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}

// Rand ... nビットのランダムな文字列を生成する。
func Rand(n int) (string, error) {
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

// UniqueID ... ユニークでソータブルなIDを作成する
func UniqueID() string {
	guid := xid.New()
	return guid.String()
}

// IsNumeric ... 数字か確認する
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
