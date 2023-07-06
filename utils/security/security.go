package security

func EncryptDecrypt(plainText string) string {
	key := byte(64)
	cipherText := ""
	for i := 0; i < len(plainText); i++ {
		// XOR each character with the key
		cipherText += string(plainText[i] ^ key)
	}
	return cipherText
}
