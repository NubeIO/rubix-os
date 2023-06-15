package database

import (
	"encoding/json"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/constants"
	"github.com/NubeIO/rubix-os/module/common"
	"github.com/NubeIO/rubix-os/src/client"
	"strings"
	"time"
)

func (d *GormDatabase) CreateNetworkPlugin(body *model.Network) (network *model.Network, err error) {
	pluginName := body.PluginPath
	if pluginName == "system" || pluginName == "module-core-system" {
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

	if strings.HasPrefix(pluginName, constants.ModulePrefix) {
		module := d.Modules[pluginName]
		if module == nil {
			return nil, moduleNotFoundError(pluginName)
		}
		bytes, _ := json.Marshal(body)
		bytes, err = module.Post(common.NetworksURL, bytes)
		if err != nil {
			return nil, err
		}
		var net *model.Network
		_ = json.Unmarshal(bytes, &net)
		return net, nil
	}

	cli := client.NewLocalClient()
	network, err = cli.CreateNetworkPlugin(body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) UpdateNetworkPlugin(uuid string, body *model.Network) (network *model.Network, err error) {
	pluginName := body.PluginPath
	if pluginName == "system" || pluginName == "module-core-system" {
		if err != nil {
			return nil, err
		}
		network, err = d.UpdateNetwork(body.UUID, body)
		if err != nil {
			return nil, err
		}
		return
	}

	if strings.HasPrefix(pluginName, constants.ModulePrefix) {
		module := d.Modules[pluginName]
		if module == nil {
			return nil, moduleNotFoundError(pluginName)
		}
		bytes, _ := json.Marshal(body)
		bytes, err = module.Patch(common.NetworksURL, uuid, bytes)
		if err != nil {
			return nil, err
		}
		var net *model.Network
		_ = json.Unmarshal(bytes, &net)
		return net, nil
	}

	cli := client.NewLocalClient()
	if err != nil {
		return nil, err
	}
	network, err = cli.UpdateNetworkPlugin(body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) DeleteNetworkPlugin(uuid string) (ok bool, err error) {
	network, err := d.GetNetwork(uuid, api.Args{})
	if err != nil {
		return false, err
	}
	pluginName := network.PluginPath
	if pluginName == "system" || pluginName == "module-core-system" {
		ok, err = d.DeleteNetwork(uuid)
		if err != nil {
			return ok, err
		}
		return
	}

	if strings.HasPrefix(pluginName, constants.ModulePrefix) {
		module := d.Modules[pluginName]
		if module == nil {
			return false, moduleNotFoundError(pluginName)
		}
		_, err = module.Delete(common.NetworksURL, uuid)
		if err != nil {
			return false, err
		}
		return true, nil
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
