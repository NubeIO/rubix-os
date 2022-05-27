package client

import "fmt"

// GetToken gets the token
func (a *FlowClient) GetToken(name string, password string) (*Token, error) {
	// Validate the name
	if name == "" {
		return nil, fmt.Errorf("provide a name in the body")
	}
	resp, err := FormatRestyResponse(a.client.R().
		SetResult(&Token{}).
		SetBody(map[string]string{"name": name}).
		SetBasicAuth(name, password).
		Post("/client"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*Token), nil
}
