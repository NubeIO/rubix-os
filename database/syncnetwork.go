package database

import (
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"reflect"
)

func (d *GormDatabase) SyncNetwork(body *interfaces.SyncNetwork) (*model.Network, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	body.NetworkName = getAutoMappedNetworkName(body.NetworkName, body.IsLocal)
	network, _ := d.GetNetworkByName(body.NetworkName, api.Args{WithTags: true})
	fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{SourceUUID: nils.NewString(body.FlowNetworkUUID)})
	if err != nil {
		return nil, err
	}
	if network == nil {
		networkModel := &model.Network{}
		networkModel.Name = body.NetworkName
		networkModel.Enable = boolean.NewTrue()
		networkModel.PluginPath = "system"
		networkModel.GlobalUUID = body.NetworkGlobalUUID
		networkModel.AutoMappingFlowNetworkName = fnc.Name
		networkModel.CreatedFromAutoMapping = boolean.NewTrue()
		networkModel.Tags = body.NetworkTags
		networkModel.MetaTags = body.NetworkMetaTags
		network, err = d.CreateNetwork(networkModel)
		return network, err
	}
	if network.GlobalUUID != body.NetworkGlobalUUID {
		return nil, errors.New("network.name already exists")
	}
	_, _ = d.CreateNetworkMetaTags(network.UUID, body.NetworkMetaTags)
	if network.Name != body.NetworkName || !reflect.DeepEqual(network.Tags, body.NetworkTags) {
		network.Name = body.NetworkName
		network.Tags = body.NetworkTags
		network, err := d.UpdateNetwork(network.UUID, network)
		return network, err
	}
	return network, nil
}
