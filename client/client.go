package client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

const defaultBaseURL = "http://localhost:1660"

// FlowClient is used to invoke Form3 Accounts API.
type FlowClient struct {
	client *resty.Client
	ClientToken string
}

// NewSession returns a new instance of FlowClient.
func NewSession(name string, password string, address string, port string) *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%s", address, port)
	apiURL := url
	client.SetHostURL(apiURL)
	client.SetError(&Error{})
	client.SetHeader("Content-Type", "application/json")
	//set token in header
	var t Token
	getToken, err := client.R().
		SetResult(&t).
		SetHeader("Content-Type", "application/json").
		SetBody(`{"name":"admin"}`).
		SetBasicAuth(name, password).
		Post("/client")
	if err != nil {
		log.Println("getToken err:", err, getToken.Status())
	}
	client.SetHeader("Authorization", t.Token)
	return &FlowClient{client: client, ClientToken: t.Token}
}


// NewSessionWithToken returns a new instance of FlowClient.
func NewSessionWithToken(token string, address string, port string) *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%s", address, port)
	apiURL := url
	client.SetHostURL(apiURL)
	client.SetError(&Error{})
	client.SetHeader("Content-Type", "application/json")
	//set token in header
	client.SetHeader("Authorization", token)
	return &FlowClient{client: client}
}



