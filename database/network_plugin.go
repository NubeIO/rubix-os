package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

func (d *GormDatabase) CreateNetworkPlugin(body *model.Network) (network *model.Network, err error) {
	pluginName := body.PluginPath
	if pluginName == "system" {
		network, err = d.CreateNetwork(body)
		if err != nil {
			return nil, err
		}
		return
	}
	body.CommonFault.InFault = true
	body.CommonFault.MessageLevel = model.MessageLevel.NoneCritical
	body.CommonFault.MessageCode = model.CommonFaultCode.PluginNotEnabled
	body.CommonFault.Message = model.CommonFaultMessage.PluginNotEnabled
	body.CommonFault.LastFail = time.Now().UTC()
	body.CommonFault.LastOk = time.Now().UTC()
	// if plugin like bacnet then call the api direct on the plugin as the plugin knows best how to add a point to keep things in sync
	cli := client.NewLocalClient()
	network, err = cli.CreateNetworkPlugin(body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) UpdateNetworkPlugin(uuid string, body *model.Network) (network *model.Network, err error) {
	pluginName := body.PluginPath
	if pluginName == "system" {
		body, err = d.updateNetworkBody(err, body)
		if err != nil {
			return nil, err
		}
		network, err = d.UpdateNetwork(body.UUID, body)
		if err != nil {
			return nil, err
		}
		return
	}
	cli := client.NewLocalClient()
	body, err = d.updateNetworkBody(err, body)
	if err != nil {
		return nil, err
	}
	network, err = cli.UpdateNetworkPlugin(body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

// restrict to update critical fields of auto-mapped network
func (d *GormDatabase) updateNetworkBody(err error, body *model.Network) (*model.Network, error) {
	network, err := d.GetNetwork(body.UUID, api.Args{})
	if err != nil {
		return nil, err
	}
	if boolean.IsTrue(network.CreatedFromAutoMapping) {
		body.AutoMappingEnable = boolean.NewFalse()
		body.AutoMappingFlowNetworkName = network.AutoMappingFlowNetworkName
		body.CreatedFromAutoMapping = network.CreatedFromAutoMapping
		body.AutoMappingUUID = network.AutoMappingUUID
	}
	return body, nil
}

func (d *GormDatabase) DeleteNetworkPlugin(uuid string) (ok bool, err error) {
	network, err := d.GetNetwork(uuid, api.Args{})
	if err != nil {
		return false, err
	}
	pluginName := network.PluginPath
	if pluginName == "system" {
		ok, err = d.DeleteNetwork(uuid)
		if err != nil {
			return ok, err
		}
		return
	}
	cli := client.NewLocalClient()
	ok, err = cli.DeleteNetworkPlugin(network, pluginName)
	if err != nil {
		ok, err = d.DeleteNetwork(uuid)
		if err != nil {
			return ok, err
		}
	}
	return
}
