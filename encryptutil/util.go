package encryptutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"strconv"
	"strings"

	"github.com/rabee-inc/go-pkg/bytesutil"
	"github.com/rabee-inc/go-pkg/stringutil"
)

const (
	sep = "::"
)

var (
	iv = []byte{34, 35, 35, 57, 68, 4, 35, 36, 7, 8, 35, 23, 35, 86, 35, 23}
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

// ToNumHash ... 文字列を数字ハッシュに変換する
func ToNumHash(str string) (uint32, error) {
	b := stringutil.ToBytes(str)
	h := fnv.New32()
	_, err := h.Write(b)
	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}

// EncryptByExpiredAt ... 有効期限付き暗号化
func EncryptByExpiredAt(text string, key string, expiredAt int64) (string, error) {
	exText := fmt.Sprintf("%d%s%s", expiredAt, sep, text)
	encText, err := Encrypt(exText, key)
	if err != nil {
		return "", err
	}
	return encText, nil
}

// DecryptByExpiredAt ... 有効期限付き復号化
func DecryptByExpiredAt(encText string, key string, now int64) (bool, string, error) {
	decText, err := Decrypt(encText, key)
	if err != nil {
		return false, "", err
	}
	index := strings.Index(decText, sep)
	expiredAt, err := strconv.ParseInt(decText[:index], 10, 64)
	if err != nil {
		return false, "", err
	}
	if expiredAt < now {
		return true, "", nil
	}
	text := decText[index+len(sep):]
	return false, text, nil
}

// Encrypt ... 暗号化
func Encrypt(text string, key string) (string, error) {
	bText := stringutil.ToBytes(text)
	bKey := stringutil.ToBytes(key)
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}
	b := base64.StdEncoding.EncodeToString(bText)
	cip := make([]byte, aes.BlockSize+len(b))
	iv := cip[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cip[aes.BlockSize:], []byte(b))
	return base64.StdEncoding.EncodeToString(cip), nil
}

// Decrypt ... 復号化
func Decrypt(encText string, key string) (string, error) {
	bKey := stringutil.ToBytes(key)
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}
	text, err := base64.StdEncoding.DecodeString(encText)
	if err != nil {
		return "", err
	}
	if len(text) < aes.BlockSize {
		return "", errors.New("too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(bytesutil.ToStr(text))
	if err != nil {
		return "", err
	}
	return bytesutil.ToStr(data), nil
}
