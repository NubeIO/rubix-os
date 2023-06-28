package database

import (
	"encoding/json"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/constants"
	"github.com/NubeIO/rubix-os/module/common"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"strings"
	"time"
)

func (d *GormDatabase) CreatePointPlugin(body *model.Point) (point *model.Point, err error) {
	network, err := d.GetNetworkByPoint(body, api.Args{})
	if err != nil {
		return nil, err
	}
	pluginName := network.PluginPath
	if pluginName == "system" || pluginName == "module-core-system" {
		body.EnableWriteable = boolean.NewTrue()
		point, err = d.CreatePoint(body)
		if err != nil {
			return nil, err
		}
		point, err = d.UpdatePoint(point.UUID, point)
		return
	}

	body.CommonFault.MessageLevel = model.MessageLevel.NoneCritical
	body.CommonFault.MessageCode = model.CommonFaultCode.PluginNotEnabled
	body.CommonFault.Message = model.CommonFaultMessage.PluginNotEnabled
	body.CommonFault.LastFail = time.Now().UTC()
	body.CommonFault.LastOk = time.Now().UTC()
	body.CommonFault.InFault = true

	if strings.HasPrefix(pluginName, constants.ModulePrefix) {
		module := d.Modules[pluginName]
		if module == nil {
			return nil, moduleNotFoundError(pluginName)
		}
		bytes, _ := json.Marshal(body)
		bytes, err = module.Post(common.PointsURL, bytes)
		if err != nil {
			return nil, err
		}
		var pnt *model.Point
		_ = json.Unmarshal(bytes, &pnt)
		return pnt, nil
	}

	cli := client.NewLocalClient()
	point, err = cli.CreatePointPlugin(body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) UpdatePointPlugin(uuid string, body *model.Point) (point *model.Point, err error) {
	network, err := d.GetNetworkByPoint(body, api.Args{})
	if err != nil {
		return nil, err
	}
	pluginName := network.PluginPath
	if pluginName == "system" || pluginName == "module-core-system" {
		body.EnableWriteable = boolean.NewTrue()
		point, err = d.UpdatePoint(uuid, body)
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
		bytes, err = module.Patch(common.PointsURL, uuid, bytes)
		if err != nil {
			return nil, err
		}
		var pnt *model.Point
		_ = json.Unmarshal(bytes, &pnt)
		return pnt, nil
	}

	cli := client.NewLocalClient()
	point, err = cli.UpdatePointPlugin(body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) WritePointPlugin(uuid string, body *model.PointWriter) (point *model.Point, err error) {
	network, err := d.GetNetworkByPointUUID(uuid, api.Args{})
	if err != nil || network == nil {
		return nil, err
	}
	pluginName := network.PluginPath
	if pluginName == "system" || pluginName == "module-core-system" {
		point, _, _, _, err = d.PointWrite(uuid, body)
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
		bytes, err = module.Patch(common.PointsWriteURL, uuid, bytes)
		if err != nil {
			return nil, err
		}
		var pnt *model.Point
		_ = json.Unmarshal(bytes, &pnt)
		return pnt, nil
	}

	cli := client.NewLocalClient()
	point, err = cli.WritePointPlugin(uuid, body, pluginName)
	if err != nil {
		return nil, err
	}
	return
}

func (d *GormDatabase) DeletePointPlugin(uuid string) (ok bool, err error) {
	point, err := d.GetPoint(uuid, api.Args{})
	if err != nil {
		return ok, err
	}
	network, err := d.GetNetworkByPoint(point, api.Args{})
	if err != nil {
		return ok, err
	}
	pluginName := network.PluginPath
	if pluginName == "system" || pluginName == "module-core-system" {
		ok, err = d.DeletePoint(uuid)
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
		_, err = module.Delete(common.PointsURL, uuid)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	cli := client.NewLocalClient()
	ok, err = cli.DeletePointPlugin(point, pluginName)
	if err != nil {
		_, err := d.DeletePoint(uuid)
		if err != nil {
			return false, err
		}
		return true, err
	}
	return
}
