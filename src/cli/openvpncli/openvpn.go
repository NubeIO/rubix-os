package openvpncli

import (
	"fmt"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nresty"
	log "github.com/sirupsen/logrus"
)

func (inst *OpenVPNClient) GetClients() (*map[string]interfaces.OpenVPNClient, error) {
	url := "/api/clients"
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&map[string]interfaces.OpenVPNClient{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*map[string]interfaces.OpenVPNClient), nil
}

func (inst *OpenVPNClient) GetOpenVPNConfig(name string) (*interfaces.OpenVPNConfig, error) {
	getURL := fmt.Sprintf("/api/openvpn/%s", name)
	resp, connectionErr, responseErr := nresty.FormatRestyV2Response(inst.Rest.R().
		SetResult(&interfaces.OpenVPNConfig{}).
		Get(getURL))
	if connectionErr != nil {
		return nil, connectionErr
	}
	if responseErr != nil {
		log.Info(fmt.Sprintf("OpenVPN is not found for %s, so generating for it", name))
		postURL := "/api/openvpn"
		_, err := nresty.FormatRestyResponse(inst.Rest.R().
			SetBody(interfaces.OpenVPNBody{Name: name}).
			SetResult(&interfaces.Message{}).
			Post(postURL))
		if err != nil {
			return nil, err
		}
		resp, err = nresty.FormatRestyResponse(inst.Rest.R().
			SetResult(&interfaces.OpenVPNConfig{}).
			Get(getURL))
		if err != nil {
			return nil, err
		}
		openVPNConfig := resp.Result().(*interfaces.OpenVPNConfig)
		return openVPNConfig, nil
	}
	openVPNConfig := resp.Result().(*interfaces.OpenVPNConfig)
	return openVPNConfig, nil
}
