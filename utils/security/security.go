package security

import (
	"encoding/base64"
)

func Encrypt(plainText string) string {
	key := byte(64)
	result := make([]byte, len(plainText))
	for i := 0; i < len(plainText); i++ {
		result[i] = plainText[i] ^ key
	}
	return base64.StdEncoding.EncodeToString(result)
}

func Decrypt(plainText string) string {
	decodeString, _ := base64.StdEncoding.DecodeString(plainText)
	key := byte(64)
	result := make([]byte, len(decodeString))
	for i := 0; i < len(decodeString); i++ {
		result[i] = decodeString[i] ^ key
	}
	return string(result)
}
