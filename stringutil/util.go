package stringutil

import (
	"bytes"
	"math"
	"strconv"
	"strings"
	"unsafe"

	"github.com/rs/xid"
)

// ToBytes ... 文字列をバイト列に変換する
func ToBytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
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

// ToComma ... 数字を金額表記にする
func ToComma(v int64) string {
	sign := ""
	if v == math.MinInt64 {
		return "-9,223,372,036,854,775,808"
	}
	if v < 0 {
		sign = "-"
		v = 0 - v
	}
	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1
	for v > 999 {
		parts[j] = strconv.FormatInt(v%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
}

// ToCommaf ... 数字を金額表記にする
func ToCommaf(v float64) string {
	buf := &bytes.Buffer{}
	if v < 0 {
		buf.Write([]byte{'-'})
		v = 0 - v
	}
	comma := []byte{','}
	parts := strings.Split(strconv.FormatFloat(v, 'f', -1, 64), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)
	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}
