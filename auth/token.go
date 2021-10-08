package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/config"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"path"
	"strings"
	"time"
)

var (
	secretKeyLocation  = "config/secret_key.txt"
	credentialLocation = "data/users.txt"
)

func GetToken(username string, password string, conf *config.Configuration) (*string, error) {
	if _, err := isValidCredential(username, password, conf); err != nil {
		log.Warn(fmt.Sprintf("Credential: %s", err))
		return nil, err
	}
	secretKey, err := GetSecretKey(path.Join(conf.Location.BiosDataDir, secretKeyLocation))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	c := jwt.MapClaims{
		"exp": time.Now().UTC().Add(time.Hour * time.Duration(12)).Unix(),
		"iat": time.Now().UTC().Unix(),
		"sub": "admin",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &signedToken, nil
}

func VerifyToken(tokenInput string, conf *config.Configuration) (bool, error) {
	secretKey, err := GetSecretKey(path.Join(conf.Location.BiosDataDir, secretKeyLocation))
	if err != nil {
		return false, err
	}
	token, err := jwt.Parse(tokenInput, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false, err
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}
	return false, nil
}

func isValidCredential(usernameInput string, password string, conf *config.Configuration) (bool, error) {
	username, hashedPassword, err := GetCredentials(path.Join(conf.Location.BiosDataDir, credentialLocation))
	if err != nil {
		return false, err
	}
	if usernameInput == *username {
		return false, errors.New("username not found")
	}
	if strings.Count(*hashedPassword, "$") < 2 {
		return false, errors.New("invalid hashed password")
	}
	hashedPasswordList := strings.Split(*hashedPassword, "$")
	if hashedPasswordList[2] == hashInternal(hashedPasswordList[1], password) {
		return true, nil
	}
	return false, errors.New("incorrect password")
}

func hashInternal(salt string, password string) string {
	mac := hmac.New(sha256.New, []byte(salt))
	mac.Write([]byte(password))
	expectedMAC := mac.Sum(nil)
	return hex.EncodeToString(expectedMAC)
}
