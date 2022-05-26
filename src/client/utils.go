package client

func (a *FlowClient) GetQuery(url string) (*[]byte, error) {
	resp, err := a.client.R().
		Get(url)
	e := checkError(resp, err)
	if e != nil {
		return nil, e
	}
	output := resp.Body()
	return &output, nil
}

func (a *FlowClient) PostQuery(url string, body interface{}) (*[]byte, error) {
	resp, err := a.client.R().
		SetBody(body).
		Post(url)
	e := checkError(resp, err)
	if e != nil {
		return nil, e
	}
	output := resp.Body()
	return &output, nil
}

func (a *FlowClient) PutQuery(url string, body interface{}) (*[]byte, error) {
	resp, err := a.client.R().
		SetBody(body).
		Patch(url)
	e := checkError(resp, err)
	if e != nil {
		return nil, e
	}
	output := resp.Body()
	return &output, nil
}

func (a *FlowClient) PatchQuery(url string, body interface{}) (*[]byte, error) {
	resp, err := a.client.R().
		SetBody(body).
		Patch(url)
	e := checkError(resp, err)
	if e != nil {
		return nil, e
	}
	output := resp.Body()
	return &output, nil
}

func (a *FlowClient) DeleteQuery(url string) error {
	resp, err := a.client.R().
		Delete(url)
	e := checkError(resp, err)
	if e != nil {
		return e
	}
	return nil
}

func (a *FlowClient) GetQueryMarshal(url string, result interface{}) (interface{}, error) {
	resp, err := a.client.R().
		SetResult(result).
		Get(url)
	e := checkError(resp, err)
	if e != nil {
		return nil, e
	}
	return resp.Result(), nil
}

func (a *FlowClient) PostQueryMarshal(url string, body interface{}, result interface{}) (interface{}, error) {
	resp, err := a.client.R().
		SetBody(body).
		SetResult(result).
		Post(url)
	e := checkError(resp, err)
	if e != nil {
		return nil, e
	}
	return resp.Result(), nil
}

func (a *FlowClient) PutQueryMarshal(url string, body interface{}, result interface{}) (interface{}, error) {
	resp, err := a.client.R().
		SetBody(body).
		SetResult(result).
		Put(url)
	e := checkError(resp, err)
	if e != nil {
		return nil, e
	}
	return resp.Result(), nil
}

func (a *FlowClient) PatchQueryMarshal(url string, body interface{}, result interface{}) (interface{}, error) {
	resp, err := a.client.R().
		SetBody(body).
		SetResult(result).
		Patch(url)
	e := checkError(resp, err)
	if e != nil {
		return nil, e
	}
	return resp.Result(), nil
}
