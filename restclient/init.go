package restclient

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

// FlowClient is used to invoke Form3 Accounts API.
type FlowClient struct {
	client *resty.Client
	ClientToken string
}

// NewNoAuth returns a new instance
func NewNoAuth(address string, port string) *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%s", address, port)
	apiURL := url
	client.SetHostURL(apiURL)
	client.SetError(&Error{})
	client.SetHeader("Content-Type", "application/json")
	return &FlowClient{client: client}
}
