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

type CSApplications struct {
	Result []struct {
		ID string `json:"id"`
	} `json:"result"`
}

type CSCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CSLoginToken struct {
	Token string `json:"jwt"`
}

var csApplications CSApplications

// InitRest Set constant CS REST params
func InitRest(address string, port int) RestClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%d/api", address, port)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Content-Type", "application/json")
	return RestClient{client: client}
}

// SetToken set the REST auth token
func (a *RestClient) SetToken(token string) {
	a.ClientToken = token
	a.client.SetHeader("Grpc-Metadata-Authorization", token)
}

// Login login to CS with username and password to get token if not provided in config
func (a *RestClient) Login(user string, pass string) error {
	token := CSLoginToken{}
	csURLConnect := "/internal/login"
	resp, err := a.client.R().
		SetBody(CSCredentials{
			Email:    user,
			Password: pass,
		}).
		SetResult(&token).
		Post(csURLConnect)
	err = checkResponse(resp, err)
	if err != nil {
		log.Warn("lorawan: Login error: ", err)
	} else {
		a.SetToken(token.Token)
	}
	return err
}

// Connect test CS connection with API token
func (a *RestClient) ConnectTest() error {
	log.Infof("lorawan: Connecting to chirpstack at %s", a.client.BaseURL)
	csURLConnect := fmt.Sprintf("/applications?limit=%s", limit)
	resp, err := a.client.R().
		SetResult(&csApplications).
		Get(csURLConnect)
	err = checkResponse(resp, err)
	if err != nil {
		log.Warn("lorawan: Connection error: ", err)
	}
	return err
}
