package client

import "github.com/NubeIO/flow-framework/nresty"

func (inst *FlowClient) GetQuery(url string) (*[]byte, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		Get(url))
	if err != nil {
		return nil, err
	}
	output := resp.Body()
	return &output, nil
}

func (inst *FlowClient) PostQuery(url string, body interface{}) (*[]byte, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		Post(url))
	if err != nil {
		return nil, err
	}
	output := resp.Body()
	return &output, nil
}

func (inst *FlowClient) PutQuery(url string, body interface{}) (*[]byte, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		Patch(url))
	if err != nil {
		return nil, err
	}
	output := resp.Body()
	return &output, nil
}

func (inst *FlowClient) PatchQuery(url string, body interface{}) (*[]byte, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		Patch(url))
	if err != nil {
		return nil, err
	}
	output := resp.Body()
	return &output, nil
}

func (inst *FlowClient) DeleteQuery(url string) error {
	_, err := nresty.FormatRestyResponse(inst.client.R().
		Delete(url))
	return err
}

func (inst *FlowClient) GetQueryMarshal(url string, result interface{}) (interface{}, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(result).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result(), nil
}

func (inst *FlowClient) PostQueryMarshal(url string, body interface{}, result interface{}) (interface{}, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		SetResult(result).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result(), nil
}

func (inst *FlowClient) PutQueryMarshal(url string, body interface{}, result interface{}) (interface{}, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		SetResult(result).
		Put(url))
	if err != nil {
		return nil, err
	}
	return resp.Result(), nil
}

func (inst *FlowClient) PatchQueryMarshal(url string, body interface{}, result interface{}) (interface{}, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		SetResult(result).
		Patch(url))
	if err != nil {
		return nil, err
	}
	return resp.Result(), nil
}
