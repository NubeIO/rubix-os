package client

import (
	"context"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/internaltoken"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/config"
	"github.com/NubeIO/rubix-os/nresty"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/go-resty/resty/v2"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	mutex       = &sync.RWMutex{}
	flowClients = map[string]*FlowClient{}
)

type FlowClient struct {
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

func GetFlowToken(ip string, port int, username string, password string) (*string, error) {
	flowClient := getNewFlowClientCli(ip, port)
	token, err := flowClient.Login(&model.LoginBody{Username: username, Password: password})
	if err != nil {
		return nil, err
	}
	return &token.AccessToken, nil
}

func getNewFlowClientCli(ip string, port int) *FlowClient {
	mutex.Lock()
	defer mutex.Unlock()
	url := fmt.Sprintf("%s://%s:%d", getSchema(port), ip, port)
	flowClient, found := flowClients[url]
	if !found {
		client := resty.New()
		client.SetDebug(false)
		client.SetBaseURL(url)
		client.SetError(&nresty.Error{})
		client.SetTransport(&transport)
		flowClient = &FlowClient{client: client}
		flowClients[url] = flowClient
	}
	return flowClient
}

func NewClient(ip string, port int, token string) *FlowClient {
	return newSessionWithToken(ip, port, token)
}

func NewLocalClient() *FlowClient {
	mutex.Lock()
	defer mutex.Unlock()
	port := config.Get().Server.Port
	url := fmt.Sprintf("%s://%s:%d", getSchema(port), "0.0.0.0", port)
	if flowClient, found := flowClients[url]; found {
		flowClient.client.SetHeader("Authorization", internaltoken.GetInternalToken(true))
		return flowClient
	}
	client := resty.New()
	client.SetDebug(false)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetTransport(&transport)
	client.SetHeader("Authorization", internaltoken.GetInternalToken(true))
	flowClient := &FlowClient{client: client}
	flowClients[url] = flowClient
	return flowClient
}

func NewFlowClientCliFromFN(fn *model.FlowNetwork) *FlowClient {
	if boolean.IsTrue(fn.IsMasterSlave) {
		return newSlaveToMasterCallSession()
	} else {
		if boolean.IsTrue(fn.IsRemote) {
			return newSessionWithToken(*fn.FlowIP, *fn.FlowPort, *fn.FlowToken)
		} else {
			return NewLocalClient()
		}
	}
}

func NewFlowClientCliFromFNC(fnc *model.FlowNetworkClone) *FlowClient {
	if boolean.IsTrue(fnc.IsMasterSlave) {
		return NewMasterToSlaveSession(fnc.GlobalUUID)
	} else {
		if boolean.IsTrue(fnc.IsRemote) {
			return newSessionWithToken(*fnc.FlowIP, *fnc.FlowPort, *fnc.FlowToken)
		} else {
			return NewLocalClient()
		}
	}
}

func newSessionWithToken(ip string, port int, token string) *FlowClient {
	mutex.Lock()
	defer mutex.Unlock()
	url := fmt.Sprintf("%s://%s:%d", getSchema(port), ip, port)
	token = fmt.Sprintf("External %s", token)
	if flowClient, found := flowClients[url]; found {
		flowClient.client.SetHeader("Authorization", token)
		return flowClient
	}
	client := resty.New()
	client.SetDebug(false)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Authorization", token)
	client.SetTransport(&transport)
	flowClient := &FlowClient{client: client}
	flowClients[url] = flowClient
	return flowClient
}

func NewMasterToSlaveSession(globalUUID string) *FlowClient {
	mutex.Lock()
	defer mutex.Unlock()
	conf := config.Get()
	url := fmt.Sprintf("http://%s:%d/slave/%s/ff", "0.0.0.0", conf.Server.RSPort, globalUUID)
	if flowClient, found := flowClients[url]; found {
		flowClient.client.SetHeader("Authorization", internaltoken.GetInternalToken(true))
		return flowClient
	}
	client := resty.New()
	client.SetDebug(false)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Authorization", internaltoken.GetInternalToken(true))
	client.SetTransport(&transport)
	flowClient := &FlowClient{client: client}
	flowClients[url] = flowClient
	return flowClient
}

func newSlaveToMasterCallSession() *FlowClient {
	mutex.Lock()
	defer mutex.Unlock()
	conf := config.Get()
	url := fmt.Sprintf("http://%s:%d/master/ff", "0.0.0.0", conf.Server.RSPort)
	if flowClient, found := flowClients[url]; found {
		flowClient.client.SetHeader("Authorization", internaltoken.GetInternalToken(true))
		return flowClient
	}
	client := resty.New()
	client.SetDebug(false)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Authorization", internaltoken.GetInternalToken(true))
	client.SetTransport(&transport)
	flowClient := &FlowClient{client: client}
	flowClients[url] = flowClient
	return flowClient
}

func getSchema(port int) string {
	if port == 443 {
		return "https"
	}
	return "http"
}
