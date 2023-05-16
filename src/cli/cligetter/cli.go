package cligetter

import (
	"github.com/NubeIO/flow-framework/src/cli/edgebioscli"
	"github.com/NubeIO/flow-framework/src/cli/edgecli"
	"github.com/NubeIO/flow-framework/src/cli/openvpncli"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func GetEdgeClient(host *model.Host) *edgecli.Client {
	cli := edgecli.New(&edgecli.Client{
		Rest:          nil,
		Ip:            host.IP,
		Port:          host.Port,
		HTTPS:         host.HTTPS,
		ExternalToken: host.ExternalToken,
	})
	return cli
}

func GetEdgeClientFastTimeout(host *model.Host) *edgecli.Client {
	cli := edgecli.NewFastTimeout(&edgecli.Client{
		Rest:          nil,
		Ip:            host.IP,
		Port:          host.Port,
		HTTPS:         host.HTTPS,
		ExternalToken: host.ExternalToken,
	})
	return cli
}

func GetEdgeBiosClient(host *model.Host) *edgebioscli.BiosClient {
	cli := edgebioscli.New(&edgebioscli.BiosClient{
		Rest:          nil,
		Ip:            host.IP,
		Port:          host.BiosPort,
		HTTPS:         host.HTTPS,
		ExternalToken: host.ExternalToken,
	})
	return cli
}

func GetOpenVPNClient() (*openvpncli.OpenVPNClient, error) {
	return openvpncli.Get()
}
