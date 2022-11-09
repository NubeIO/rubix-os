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

var csApplications CSApplications

// InitRest Set constant CS REST params
func InitRest(address string, port int, token string) RestClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%d/api", address, port)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Grpc-Metadata-Authorization", token)
	return RestClient{client: client, ClientToken: token}
}

// NewNoAuth returns a new instance
func NewNoAuth(address string, port int) *RestClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%d", address, port)
	apiURL := url
	client.SetBaseURL(apiURL)
	client.SetError(&nresty.Error{})
	client.SetHeader("Content-Type", "application/json")
	return &RestClient{client: client}
}

// Connect test CS connection with API token
func (a *RestClient) Connect() error {
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
