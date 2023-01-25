package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"reflect"
)

func (d *GormDatabase) SyncNetwork(body *model.SyncNetwork) (*model.Network, error) {
	network, _ := d.GetNetworkByAutoMappingUUID(body.NetworkUUID, api.Args{WithTags: true})
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
		networkModel.Tags = body.NetworkTags
		return d.CreateNetwork(networkModel, false)
	}
	if network.Name != networkName || !reflect.DeepEqual(network.Tags, body.NetworkTags) {
		network.Name = networkName
		network.Tags = body.NetworkTags
		return d.UpdateNetwork(network.UUID, network, false)
	}
	return network, nil
}
