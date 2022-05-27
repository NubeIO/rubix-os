package client

func (a *FlowClient) GetQuery(url string) (*[]byte, error) {
	resp, err := FormatRestyResponse(a.client.R().
		Get(url))
	if err != nil {
		return nil, err
	}
	output := resp.Body()
	return &output, nil
}

func (a *FlowClient) PostQuery(url string, body interface{}) (*[]byte, error) {
	resp, err := FormatRestyResponse(a.client.R().
		SetBody(body).
		Post(url))
	if err != nil {
		return nil, err
	}
	output := resp.Body()
	return &output, nil
}

func (a *FlowClient) PutQuery(url string, body interface{}) (*[]byte, error) {
	resp, err := FormatRestyResponse(a.client.R().
		SetBody(body).
		Patch(url))
	if err != nil {
		return nil, err
	}
	output := resp.Body()
	return &output, nil
}

func (a *FlowClient) PatchQuery(url string, body interface{}) (*[]byte, error) {
	resp, err := FormatRestyResponse(a.client.R().
		SetBody(body).
		Patch(url))
	if err != nil {
		return nil, err
	}
	output := resp.Body()
	return &output, nil
}

func (a *FlowClient) DeleteQuery(url string) error {
	_, err := FormatRestyResponse(a.client.R().
		Delete(url))
	return err
}

func (a *FlowClient) GetQueryMarshal(url string, result interface{}) (interface{}, error) {
	resp, err := FormatRestyResponse(a.client.R().
		SetResult(result).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result(), nil
}

func (a *FlowClient) PostQueryMarshal(url string, body interface{}, result interface{}) (interface{}, error) {
	resp, err := FormatRestyResponse(a.client.R().
		SetBody(body).
		SetResult(result).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result(), nil
}

func (a *FlowClient) PutQueryMarshal(url string, body interface{}, result interface{}) (interface{}, error) {
	resp, err := FormatRestyResponse(a.client.R().
		SetBody(body).
		SetResult(result).
		Put(url))
	if err != nil {
		return nil, err
	}
	return resp.Result(), nil
}

func (a *FlowClient) PatchQueryMarshal(url string, body interface{}, result interface{}) (interface{}, error) {
	resp, err := FormatRestyResponse(a.client.R().
		SetBody(body).
		SetResult(result).
		Patch(url))
	if err != nil {
		return nil, err
	}
	return resp.Result(), nil
}
