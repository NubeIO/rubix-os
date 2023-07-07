package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigEnv(t *testing.T) {
	key := "AAAA8N82UWY:APA9ddEQZxR3NKKBwWUaFTktMjpzFqD63oaOhFyQsynKK90M0SdPwmxCPwUwqvKmjPxVavULOZWL_qjGagD7LZaC-sHZkDbYRLL37UFC8MQagX2oRmyN3TIys9dQtRW8O2SMNbmpTfnf"
	encrypted := Encrypt(key)
	decrypted := Decrypt(encrypted)
	assert.Equal(t, key, decrypted)
}
