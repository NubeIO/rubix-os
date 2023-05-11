package csrest

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/NubeIO/flow-framework/nresty"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type ChirpClient struct {
	client      *resty.Client
	ClientToken string
	csURL       *url.URL
	proxy       *httputil.ReverseProxy
	basePath    string
}

type CSApplications struct {
	Result []struct {
		ID string `json:"id"`
	} `json:"result"`
}

var csApplications CSApplications

const CsURLPrefix = "/cs"

// InitRest Set constant CS REST params
func InitRest(address string, port int, basePath string) ChirpClient {

	client := resty.New()
	chirpClient := ChirpClient{}
	chirpClient.csURL, _ = url.Parse(fmt.Sprintf("http://%s:%d/api", address, port))

	client.SetDebug(false)
	client.SetBaseURL(chirpClient.csURL.String())
	client.SetError(&nresty.Error{})
	client.SetHeader("Content-Type", "application/json")
	chirpClient.client = client
	chirpClient.basePath = basePath

	chirpClient.proxy = httputil.NewSingleHostReverseProxy(chirpClient.csURL)
	chirpClient.proxy.Director = chirpClient.proxyDirector

	return chirpClient
}

// SetToken set the REST auth token
func (inst *ChirpClient) SetToken(token string) {
	inst.ClientToken = token
	inst.client.SetHeader("Grpc-Metadata-Authorization", token)
}

// Connect test CS connection with API token
func (inst *ChirpClient) ConnectTest() error {
	log.Infof("lorawan: Connecting to chirpstack at %s", inst.client.BaseURL)
	csURLConnect := fmt.Sprintf("/applications?limit=%s", limit)
	resp, err := inst.client.R().
		SetResult(&csApplications).
		Get(csURLConnect)
	err = checkResponse(resp, err)
	if err != nil {
		log.Warn("lorawan: Connection error: ", err)
	}
	return err
}

func (inst *ChirpClient) proxyDirector(req *http.Request) {
	req.URL.Scheme = inst.csURL.Scheme
	req.URL.Host = inst.csURL.Host
	req.URL.Path = inst.csURL.Path + strings.TrimPrefix(req.URL.Path, inst.basePath+CsURLPrefix)
}
