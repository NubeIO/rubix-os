package edgecli

import (
	"fmt"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *Client) CreateDir(path string) (*interfaces.Message, error) {
	url := fmt.Sprintf("/api/dirs/create?path=%s", path)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&interfaces.Message{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*interfaces.Message), nil
}

func (inst *Client) DirExists(path string) (*interfaces.DirExistence, error) {
	url := fmt.Sprintf("/api/dirs/exists?path=%s", path)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&interfaces.DirExistence{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*interfaces.DirExistence), nil
}
