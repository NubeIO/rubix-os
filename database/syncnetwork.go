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
	d.mutex.Lock()
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
		network, err := d.CreateNetwork(networkModel)
		d.mutex.Unlock()
		return network, err
	}
	_, _ = d.CreateNetworkMetaTags(network.UUID, body.NetworkMetaTags)
	if network.Name != networkName || !reflect.DeepEqual(network.Tags, body.NetworkTags) {
		network.Name = networkName
		network.Tags = body.NetworkTags
		network, err := d.UpdateNetwork(network.UUID, network)
		d.mutex.Unlock()
		return network, err
	}
	d.mutex.Unlock()
	return network, nil
}
