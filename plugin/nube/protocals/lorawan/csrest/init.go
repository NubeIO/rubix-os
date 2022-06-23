package csrest

import (
	"fmt"

	"github.com/NubeIO/flow-framework/nresty"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type RestClient struct {
	client      *resty.Client
	ClientToken string
}

type Token struct {
	JWT string `json:"jwt"`
}

type user struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const csURLLogin = "/api/internal/login"

// CSLogin login to Chirpstack and get JWT token
func CSLogin(address string, port int, username string, password string) (*RestClient, error) {
	log.Infof("lorawan: Connecting to chirpstack at %s:%d", address, port)
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%d", address, port)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Content-Type", "application/json")
	var t Token
	resp, err := client.R().
		SetResult(&t).
		SetHeader("Content-Type", "application/json").
		SetBody(user{Email: username, Password: password}).
		Post(csURLLogin)
	err = checkResponse(resp, err)
	if err != nil {
		log.Warn("lorawan: Connection error: ", err)
		return nil, err
	}
	client.SetHeader("Grpc-Metadata-Authorization", t.JWT)
	return &RestClient{client: client, ClientToken: t.JWT}, nil
}
