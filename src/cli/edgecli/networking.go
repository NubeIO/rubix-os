package edgecli

import (
	"github.com/NubeIO/lib-networking/networking"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *Client) GetNetworking() ([]networking.NetworkInterfaces, error) {
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&[]networking.NetworkInterfaces{}).
		Get("api/networking"))
	if err != nil {
		return nil, err
	}
	data := resp.Result().(*[]networking.NetworkInterfaces)
	return *data, nil
}

func (inst *Client) GetInternetIP() (networking.Check, error) {
	checkResult := networking.Check{}
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&checkResult).
		Get("api/networking/internet"))
	if err != nil {
		return checkResult, err
	}
	data := resp.Result().(*networking.Check)
	return *data, nil
}
