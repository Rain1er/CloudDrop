package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// Encrypt keeps the legacy request-side offset used by PHP payloads.
func Encrypt(code string, password string) string {
	return EncryptWithOffset([]byte(code), password, 1)
}

// Decrypt keeps the legacy response-side offset used by all current payloads.
func Decrypt(code string, password string) string {
	decrypted, _ := DecryptWithOffset(code, password, 5)
	return decrypted
}

func DeriveXORKey(password string) string {
	md5Hash := md5.Sum([]byte(password))
	return hex.EncodeToString(md5Hash[:])[:16]
}

func EncryptWithOffset(data []byte, password string, offset int) string {
	return EncryptWithKeyOffset(data, DeriveXORKey(password), offset)
}

func EncryptWithKeyOffset(data []byte, key string, offset int) string {
	out := make([]byte, len(data))
	copy(out, data)
	for i := range out {
		out[i] = out[i] ^ key[(i+offset)&15]
	}
	return base64.StdEncoding.EncodeToString(out)
}

func DecryptWithOffset(code string, password string, offset int) (string, error) {
	return DecryptWithKeyOffset(code, DeriveXORKey(password), offset)
}

func DecryptWithKeyOffset(code string, key string, offset int) (string, error) {
	out, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}
	for i := range out {
		out[i] = out[i] ^ key[(i+offset)&15]
	}
	return string(out), nil
}
