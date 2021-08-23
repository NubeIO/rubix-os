package client


import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
)

const defaultBaseURL = "http://localhost:1660"

// FlowClient is used to invoke Form3 Accounts API.
type FlowClient struct {
	client *resty.Client
}

// NewFlowRestClient returns a new instance of FlowClient.
func NewFlowRestClient(name string, password string) *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	// Try getting Accounts API base URL from env var
	apiURL := os.Getenv("API_ADDR")
	if apiURL == "" {
		apiURL = defaultBaseURL
	}
	client.SetHostURL(apiURL)
	// Setting global error struct that maps to Form3's error response
	client.SetError(&Error{})
	client.SetHeader("Content-Type", "application/json")

	//set token in header
	var t Token
	getToken, err := client.R().
		SetResult(&t).
		SetHeader("Content-Type", "application/json").
		SetBody(`{"name":"admin"}`).
		SetBasicAuth("admin", "admin").
		Post("/client")
	if err != nil {
		log.Println("getToken err:", err, getToken.Status())
	}
	fmt.Println("token:", t.Token)
	client.SetHeader("Authorization", t.Token)
	return &FlowClient{client: client}
}


