package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) SyncNetwork(body *model.SyncNetwork) (*model.Network, error) {
	network, _ := d.GetNetworkByAutoMappingUUID(body.NetworkUUID, api.Args{})
	networkName := body.NetworkName
	if body.IsLocal {
		networkName = fmt.Sprintf("mapping_%s", networkName)
	}
	if network == nil {
		networkModel := &model.Network{}
		networkModel.Name = networkName
		networkModel.AutoMappingUUID = body.NetworkUUID
		networkModel.Enable = boolean.NewTrue()
		networkModel.PluginPath = "system"
		return d.CreateNetwork(networkModel, false)
	}
	if network.Name != networkName {
		network.Name = networkName
		return d.UpdateNetwork(network.UUID, network, false)
	}
	return network, nil
}
