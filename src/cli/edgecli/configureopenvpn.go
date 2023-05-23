package edgecli

import (
	"fmt"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nresty"
	log "github.com/sirupsen/logrus"
	"path"
)

func (inst *Client) ConfigureOpenVPN(openvpnConfig *interfaces.OpenVPNConfig) (*interfaces.Message, error) {
	const configDir = "/etc/openvpn"
	const configName = "client.conf"
	const serviceName = "openvpn@client.service"

	url := fmt.Sprintf("/api/dirs/create?path=%s", configDir)
	_, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&interfaces.Message{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	log.Info("Created OpenVPN config directory")

	url = fmt.Sprintf("/api/files/write?file=%s", path.Join(configDir, configName))
	_, err = nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&interfaces.Message{}).
		SetBody(openvpnConfig).
		Put(url))
	if err != nil {
		return nil, err
	}
	log.Info("Configured OpenVPN config")

	url = fmt.Sprintf("/api/systemctl/enable?unit=%s", serviceName)
	if _, err = nresty.FormatRestyResponse(inst.Rest.R().Post(url)); err != nil {
		log.Error(err)
	}
	log.Info("Enabled OpenVPN")

	url = fmt.Sprintf("/api/systemctl/restart?unit=%s", serviceName)
	if _, err = nresty.FormatRestyResponse(inst.Rest.R().Post(url)); err != nil {
		log.Error(err)
	}
	log.Info("Restarted OpenVPN")
	return &interfaces.Message{Message: "Configured OpenVPN config"}, nil
}
