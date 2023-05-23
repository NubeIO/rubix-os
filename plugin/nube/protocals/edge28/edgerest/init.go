package edgerest

import (
	"fmt"
	"github.com/NubeIO/rubix-os/nresty"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// RestClient is used to invoke Form3 Accounts API.
type RestClient struct {
	client      *resty.Client
	ClientToken string
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

func (*RestClient) edge28ClientDebugMsg(args ...interface{}) {
	enable := false
	if enable {
		prefix := "Edge28 Client: "
		log.Info(prefix, args)
	}
}

func (*RestClient) edge28ClientErrorMsg(args ...interface{}) {
	prefix := "Edge28 Client: "
	log.Error(prefix, args)
}
