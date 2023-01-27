package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"reflect"
)

func (d *GormDatabase) SyncNetwork(body *interfaces.SyncNetwork) (*model.Network, error) {
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
		networkModel.MetaTags = body.NetworkMetaTags
		return d.CreateNetwork(networkModel, false)
	}
	_, _ = d.CreateNetworkMetaTags(network.UUID, body.NetworkMetaTags)
	if network.Name != networkName || !reflect.DeepEqual(network.Tags, body.NetworkTags) {
		network.Name = networkName
		network.Tags = body.NetworkTags
		return d.UpdateNetwork(network.UUID, network, false)
	}
	return network, nil
}
