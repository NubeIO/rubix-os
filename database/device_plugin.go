package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/module/common"
	"github.com/NubeIO/rubix-os/src/client"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (d *GormDatabase) CreateDevicePlugin(body *model.Device) (device *model.Device, err error) {
	network, err := d.GetNetwork(body.NetworkUUID, api.Args{})
	if network == nil {
		errMsg := fmt.Sprintf("model.device failed to find a network with uuid:%s", body.NetworkUUID)
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	pluginName := network.PluginPath

	if pluginName == "system" {
		device, err = d.CreateDevice(body)
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

	if strings.HasPrefix(pluginName, "module") {
		module := d.Modules[pluginName]
		if module == nil {
			return nil, moduleNotFoundError(pluginName)
		}
		bytes, _ := json.Marshal(body)
		bytes, err = module.Post(common.DevicesURL, bytes)
		if err != nil {
			return nil, err
		}
		var dev *model.Device
		_ = json.Unmarshal(bytes, &dev)
		return dev, nil
	}

	cli := client.NewLocalClient()
	device, err = cli.CreateDevicePlugin(body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) UpdateDevicePlugin(uuid string, body *model.Device) (device *model.Device, err error) {
	network, err := d.GetNetwork(body.NetworkUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	pluginName := network.PluginPath
	if pluginName == "system" {
		device, err = d.UpdateDevice(uuid, body)
		if err != nil {
			return nil, err
		}
		return
	}

	if strings.HasPrefix(pluginName, "module") {
		module := d.Modules[pluginName]
		if module == nil {
			return nil, moduleNotFoundError(pluginName)
		}
		bytes, _ := json.Marshal(body)
		bytes, err = module.Patch(common.DevicesURL, uuid, bytes)
		if err != nil {
			return nil, err
		}
		var dev *model.Device
		_ = json.Unmarshal(bytes, &dev)
		return dev, nil
	}

	cli := client.NewLocalClient()
	device, err = cli.UpdateDevicePlugin(body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) DeleteDevicePlugin(uuid string) (ok bool, err error) {
	device, err := d.GetDevice(uuid, api.Args{})
	if err != nil {
		return false, err
	}
	getNetwork, err := d.GetNetwork(device.NetworkUUID, api.Args{})
	if err != nil {
		return false, err
	}
	pluginName := getNetwork.PluginPath
	if pluginName == "system" {
		ok, err = d.DeleteDevice(uuid)
		if err != nil {
			return ok, err
		}
		return
	}

	if strings.HasPrefix(pluginName, "module") {
		module := d.Modules[pluginName]
		if module == nil {
			return false, moduleNotFoundError(pluginName)
		}
		_, err = module.Delete(common.DevicesURL, uuid)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	cli := client.NewLocalClient()
	ok, err = cli.DeleteDevicePlugin(device, pluginName)
	if err != nil {
		_, err := d.DeleteDevice(uuid)
		if err != nil {
			return false, err
		}
		return true, err
	}
	return
}
