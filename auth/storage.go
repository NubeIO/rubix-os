package auth

import (
	"errors"
	"io/ioutil"
	"strings"
)

func GetSecretKey(filename string) ([]byte, error) {
	secretKey, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return secretKey, nil
}

func GetCredentials(filename string) (*string, *string, error) {
	credential, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	parts := strings.Split(string(credential), ":")
	if len(parts) != 2 {
		return nil, nil, errors.New("check users file")
	}
	username := parts[0]
	hashedPassword := parts[1]
	return &username, &hashedPassword, nil
}
