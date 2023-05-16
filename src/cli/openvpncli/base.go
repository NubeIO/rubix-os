package openvpncli

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

var openVPNClient = &OpenVPNClient{}

type OpenVPNClient struct {
	Rest *resty.Client
}

func Get() (*OpenVPNClient, error) {
	if os.Getenv("OPENVPN_ENABLED") != "true" {
		return nil, errors.New("OpenVPN is not enabled")
	}
	if openVPNClient.Rest == nil {
		rest := resty.New()
		rest.SetBaseURL(getBaseUrl())
		openVPNClient.Rest = rest
	}
	return openVPNClient, nil
}

func getBaseUrl() string {
	openvpnHost := os.Getenv("OPENVPN_HOST")
	openvpnPort := os.Getenv("OPENVPN_PORT")
	return fmt.Sprintf("http://%s:%s", openvpnHost, openvpnPort)
}
