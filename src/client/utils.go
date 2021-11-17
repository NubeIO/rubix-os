package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/go-resty/resty/v2"
)

func (a *FlowClient) GetQuery(url string) (*[]byte, error) {
	resp, err := a.client.R().
		Get(url)
	e := checkError("GET", url, resp, err)
	if e != nil {
		return nil, *e
	}
	output := resp.Body()
	return &output, nil
}

func (a *FlowClient) PostQuery(url string, body interface{}) (*[]byte, error) {
	resp, err := a.client.R().
		SetBody(body).
		Post(url)
	e := checkError("POST", url, resp, err)
	if e != nil {
		return nil, *e
	}
	output := resp.Body()
	return &output, nil
}

func (a *FlowClient) PutQuery(url string, body interface{}) (*[]byte, error) {
	resp, err := a.client.R().
		SetBody(body).
		Patch(url)
	e := checkError("PUT", url, resp, err)
	if e != nil {
		return nil, *e
	}
	output := resp.Body()
	return &output, nil
}

func (a *FlowClient) PatchQuery(url string, body interface{}) (*[]byte, error) {
	resp, err := a.client.R().
		SetBody(body).
		Patch(url)
	e := checkError("PATCH", url, resp, err)
	if e != nil {
		return nil, *e
	}
	output := resp.Body()
	return &output, nil
}

func (a *FlowClient) DeleteQuery(url string) error {
	resp, err := a.client.R().
		Delete(url)
	e := checkError("DELETE", url, resp, err)
	if e != nil {
		return *e
	}
	return nil
}

func (a *FlowClient) GetQueryMarshal(url string, result interface{}) (interface{}, error) {
	resp, err := a.client.R().
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

func (a *FlowClient) PutQueryMarshal(url string, body interface{}, result interface{}) (interface{}, error) {
	resp, err := a.client.R().
		SetBody(body).
		SetResult(result).
		Put(url)
	e := checkError("PUT", url, resp, err)
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
