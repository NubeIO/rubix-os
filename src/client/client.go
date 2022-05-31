package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/auth"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/go-resty/resty/v2"
)

var flowClients = map[string]*FlowClient{}

type FlowClient struct {
	client *resty.Client
}

func GetFlowToken(ip string, port int, username string, password string) (*string, error) {
	url := fmt.Sprintf("%s://%s:%d", getSchema(port), ip, port)
	flowClient, found := flowClients[url]
	if !found {
		client := resty.New()
		client.SetDebug(false)
		client.SetBaseURL(url)
		client.SetError(&nresty.Error{})
		flowClient = &FlowClient{client: client}
	}
	token, err := flowClient.Login(&model.LoginBody{Username: username, Password: password})
	if err != nil {
		return nil, err
	}
	return &token.AccessToken, nil
}

func NewLocalClient() *FlowClient {
	port := config.Get().Server.Port
	url := fmt.Sprintf("%s://%s:%d", getSchema(port), "0.0.0.0", port)
	if flowClient, found := flowClients[url]; found {
		return flowClient
	}
	client := resty.New()
	client.SetDebug(false)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
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
		return newMasterToSlaveSession(fnc.GlobalUUID)
	} else {
		if boolean.IsTrue(fnc.IsRemote) {
			return newSessionWithToken(*fnc.FlowIP, *fnc.FlowPort, *fnc.FlowToken)
		} else {
			return NewLocalClient()
		}
	}
}

func newSessionWithToken(ip string, port int, token string) *FlowClient {
	url := fmt.Sprintf("%s://%s:%d/ff", getSchema(port), ip, port)
	if flowClient, found := flowClients[url]; found {
		return flowClient
	}
	client := resty.New()
	client.SetDebug(false)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Authorization", token)
	flowClient := &FlowClient{client: client}
	flowClients[url] = flowClient
	return flowClient
}

func newMasterToSlaveSession(globalUUID string) *FlowClient {
	conf := config.Get()
	url := fmt.Sprintf("http://%s:%d/slave/%s/ff", "0.0.0.0", conf.Server.RSPort, globalUUID)
	if flowClient, found := flowClients[url]; found {
		return flowClient
	}
	client := resty.New()
	client.SetDebug(false)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Authorization", auth.GetRubixServiceInternalToken())
	flowClient := &FlowClient{client: client}
	flowClients[url] = flowClient
	return flowClient
}

func newSlaveToMasterCallSession() *FlowClient {
	conf := config.Get()
	url := fmt.Sprintf("http://%s:%d/master/ff", "0.0.0.0", conf.Server.RSPort)
	if flowClient, found := flowClients[url]; found {
		return flowClient
	}
	client := resty.New()
	client.SetDebug(false)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Authorization", auth.GetRubixServiceInternalToken())
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
