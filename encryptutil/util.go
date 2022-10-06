package encryptutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/rabee-inc/go-pkg/stringutil"
)

// ToSHA256 ... 文字列のハッシュ(SHA256)を取得する
func ToSHA256(str string) string {
	c := sha256.Sum256([]byte(str))
	return hex.EncodeToString(c[:])
}

func ToHmacSHA256(str string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(str))
	return hex.EncodeToString(mac.Sum(nil))
}

// Encrypt ... 暗号化
func Encrypt(plainText string, key string) (string, error) {
	if plainText == "" {
		return "", nil
	}
	block, err := aes.NewCipher(stringutil.ToBytes(key))
	if err != nil {
		return "", err
	}
	paddedPlaintext := padByPkcs7([]byte(plainText))
	cipherText := make([]byte, aes.BlockSize+len(paddedPlaintext))
	iv := cipherText[:aes.BlockSize]
	_, err = rand.Read(iv)
	if err != nil {
		return "", err
	}
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(cipherText[aes.BlockSize:], paddedPlaintext)
	cipherTextBase64 := base64.StdEncoding.EncodeToString(cipherText)
	return cipherTextBase64, nil
}

// Decrypt ... 復号化
func Decrypt(encText string, key string) (string, error) {
	if encText == "" {
		return "", nil
	}
	block, err := aes.NewCipher(stringutil.ToBytes(key))
	if err != nil {
		return "", err
	}
	cipherText, err := base64.StdEncoding.DecodeString(string(encText))
	if err != nil {
		return "", err
	}
	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipher text must be longer than blocksize")
	}
	if len(cipherText)%aes.BlockSize != 0 {
		return "", errors.New("cipher text must be multiple of blocksize(256bit)")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	plainText := make([]byte, len(cipherText))
	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(plainText, cipherText)
	return string(unPadByPkcs7(plainText)), nil
}

func padByPkcs7(data []byte) []byte {
	padSize := aes.BlockSize
	if len(data)%aes.BlockSize != 0 {
		padSize = aes.BlockSize - (len(data))%aes.BlockSize
	}
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(data, pad...)
}

func unPadByPkcs7(data []byte) []byte {
	padSize := int(data[len(data)-1])
	return data[:len(data)-padSize]
}
