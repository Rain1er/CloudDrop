package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

// 加密code
func Encrypt(code string, password string) (EnCode string) {
	// MD5 hash the password and take the first 16 characters
	md5Hash := md5.Sum([]byte(password))
	key := hex.EncodeToString(md5Hash[:])[:16]
	encryptedCode := make([]byte, len(code))
	copy(encryptedCode, code)

	for i := range encryptedCode {
		encryptedCode[i] = encryptedCode[i] ^ key[(i+1)&15] // Key offset by one digit
	}
	// Base64 encode the encrypted code
	encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedCode)
	return encryptedBase64
}

// 解密
func Decrypt(code string, password string) (DeCode string) {
	md5Hash := md5.Sum([]byte(password))
	key := hex.EncodeToString(md5Hash[:])[:16]

	decryptedCode := make([]byte, len(code))
	copy(decryptedCode, code)

	// base64 decode and xor
	decryptedCode, _ = base64.StdEncoding.DecodeString(string(decryptedCode))
	for i := range len(decryptedCode) {
		decryptedCode[i] = decryptedCode[i] ^ key[(i+1)&15] // Key offset by one digit
	}
	return string(decryptedCode)
}
