package client

import (
	"fmt"

)

// ClientGetRubixPlat an object
func (a *FlowClient) ClientGetRubixPlat() (*ResponseBody, error) {
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		Get("/api/wires/plat")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	fmt.Println(resp.String())
	return resp.Result().(*ResponseBody), nil
}

