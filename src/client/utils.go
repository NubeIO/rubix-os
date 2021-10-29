package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/go-resty/resty/v2"
)

func (a *FlowClient) GetQuery(url string, params map[string]string) (*string, error) {
	resp, err := a.client.R().
		SetPathParams(params).
		Get(url)
	e := checkError("GET", url, resp, err)
	if e != nil {
		return nil, *e
	}
	return utils.NewStringAddress(resp.String()), nil
}

func (a *FlowClient) PostQuery(url string, body interface{}) (*string, error) {
	resp, err := a.client.R().
		SetBody(body).
		Post(url)
	e := checkError("POST", url, resp, err)
	if e != nil {
		return nil, *e
	}
	return utils.NewStringAddress(resp.String()), nil
}

func (a *FlowClient) PatchQuery(url string, body interface{}) (*string, error) {
	resp, err := a.client.R().
		SetBody(body).
		Patch(url)
	e := checkError("PATCH", url, resp, err)
	if e != nil {
		return nil, *e
	}
	return utils.NewStringAddress(resp.String()), nil
}

func (a *FlowClient) DeleteQuery(url string) error {
	resp, err := a.client.R().
		Patch(url)
	e := checkError("DELETE", url, resp, err)
	if e != nil {
		return *e
	}
	return nil
}

func (a *FlowClient) GetQueryMarshal(url string, params map[string]string, result interface{}) (interface{}, error) {
	resp, err := a.client.R().
		SetQueryParams(params).
		SetResult(result).
		Get(url)
	e := checkError("GET", url, resp, err)
	if e != nil {
		return nil, *e
	}
	return resp.Result(), nil
}

func (a *FlowClient) PostQueryMarshal(url string, body interface{}, result interface{}) (interface{}, error) {
	resp, err := a.client.R().
		SetBody(body).
		SetResult(result).
		Post(url)
	e := checkError("POST", url, resp, err)
	if e != nil {
		return nil, *e
	}
	return resp.Result(), nil
}

func (a *FlowClient) PatchQueryMarshal(url string, body interface{}, result interface{}) (interface{}, error) {
	resp, err := a.client.R().
		SetBody(body).
		SetResult(result).
		Patch(url)
	e := checkError("PATCH", url, resp, err)
	if e != nil {
		return nil, *e
	}
	return resp.Result(), nil
}

func checkError(method string, url string, resp *resty.Response, err error) *error {
	if err != nil {
		if resp == nil || resp.String() == "" {
			return utils.NewError(fmt.Errorf("%s %s: %s", method, url, err))
		} else {
			return utils.NewError(fmt.Errorf("%s %s: %s", method, url, resp))
		}
	}
	if resp.IsError() {
		return utils.NewError(fmt.Errorf("%s %s: %s", method, url, resp))
	}
	return nil
}
