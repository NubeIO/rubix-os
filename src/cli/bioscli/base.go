package bioscli

import (
	"context"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/internaltoken"
	"github.com/NubeIO/rubix-os/nresty"
	"github.com/go-resty/resty/v2"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	mutex       = &sync.RWMutex{}
	biosClients = map[string]*BiosClient{}
)

type BiosClient struct {
	client *resty.Client
}

// The dialTimeout normally catches: when the server is unreachable and returns i/o timeout within 2 seconds.
// Otherwise, the i/o timeout takes 1.3 minutes on default; which is a very long time for waiting.
// It uses the DialTimeout function of the net package which connects to a server address on a named network before
// a specified timeout.
func dialTimeout(_ context.Context, network, addr string) (net.Conn, error) {
	timeout := 2 * time.Second
	return net.DialTimeout(network, addr, timeout)
}

var transport = http.Transport{
	DialContext: dialTimeout,
}

func NewLocalBiosClient() *BiosClient {
	mutex.Lock()
	defer mutex.Unlock()
	url := "http://localhost:1659"
	if flowClient, found := biosClients[url]; found {
		flowClient.client.SetHeader("Authorization", internaltoken.GetInternalToken(true))
		return flowClient
	}
	client := resty.New()
	client.SetDebug(false)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetTransport(&transport)
	client.SetHeader("Authorization", internaltoken.GetInternalToken(true))
	flowClient := &BiosClient{client: client}
	biosClients[url] = flowClient
	return flowClient
}
