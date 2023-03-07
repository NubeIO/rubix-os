package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/interfaces/connection"
	"github.com/NubeIO/flow-framework/plugin/compat"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"strings"
	"sync"
)

func marshallCacheNetworks(networks []*model.Network, args api.Args) {
	for _, network := range networks {
		for _, device := range network.Devices {
			marshallCachePoints(device.Points, args)
		}
	}
}

func (d *GormDatabase) GetNetworks(args api.Args) ([]*model.Network, error) {
	var networksModel []*model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Find(&networksModel).Error; err != nil {
		return nil, err
	}
	marshallCacheNetworks(networksModel, args)
	return networksModel, nil
}

func (d *GormDatabase) GetNetwork(uuid string, args api.Args) (*model.Network, error) {
	var networkModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&networkModel).Error; err != nil {
		return nil, err
	}
	marshallCacheDevices(networkModel.Devices, args)
	return networkModel, nil
}

// GetNetworkByField returns the network for the given field ie name or nil.
func (d *GormDatabase) GetNetworkByField(field string, value string, withDevices bool) (*model.Network, error) {
	var networkModel *model.Network
	f := fmt.Sprintf("%s = ? ", field)
	if withDevices { // drop child to reduce json size
		query := d.DB.Where(f, value).Preload("Devices").First(&networkModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return networkModel, nil
	} else {
		query := d.DB.Where(f, value).First(&networkModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return networkModel, nil
	}
}

func (d *GormDatabase) CreateNetwork(body *model.Network) (*model.Network, error) {
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Network)
	body.Name = strings.TrimSpace(body.Name)
	body.ThingClass = model.ThingClass.Network
	transport, err := checkTransport(body.TransportType)
	if err != nil {
		return nil, err
	}
	body.TransportType = transport
	if body.PluginPath != "" || body.PluginConfId != "" {
		if body.PluginConfId == "" {
			plugin, err := d.GetPluginByPath(body.PluginPath)
			if err != nil {
				return nil, errors.New("failed to find a valid plugin")
			}
			if plugin.UUID == "" && body.PluginConfId != "" {
				return nil, errors.New("failed to find a valid plugin uuid")
			}
			body.PluginConfId = plugin.UUID
		}
	} else {
		return nil, errors.New("provide a plugin name ie: system, lora, modbus, lorawan, bacnet")
	}
	if body.GlobalUUID == "" {
		deviceInfo, err := deviceinfo.GetDeviceInfo()
		if err != nil {
			return nil, err
		}
		body.GlobalUUID = deviceInfo.GlobalUUID
	}
	if err = d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateNetwork(uuid string, body *model.Network) (*model.Network, error) {
	var networkModel *model.Network
	query := d.DB.Where("uuid = ?", uuid).First(&networkModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&networkModel, body.Tags); err != nil {
			return nil, err
		}
	}
	body.Name = strings.TrimSpace(body.Name)
	query = d.DB.Model(&networkModel).Select("*").Updates(&body)
	if query.Error != nil {
		return nil, query.Error
	}
	return networkModel, nil
}

// UpdateNetworkErrors will only update the CommonFault properties of the network, all other properties won't be updated.
// Does not update `LastOk`.
func (d *GormDatabase) UpdateNetworkErrors(uuid string, body *model.Network) error {
	return d.DB.Model(&body).
		Where("uuid = ?", uuid).
		Select("InFault", "MessageLevel", "MessageCode", "Message", "LastFail", "InSync", "Connection").
		Updates(&body).
		Error
}

func (d *GormDatabase) DeleteNetwork(uuid string) (ok bool, err error) {
	networkModel, err := d.GetNetwork(uuid, api.Args{WithDevices: true})
	if err != nil {
		return false, err
	}
	var wg sync.WaitGroup
	for _, device := range networkModel.Devices {
		wg.Add(1)
		device := device
		go func() {
			defer wg.Done()
			if boolean.IsTrue(device.AutoMappingEnable) {
				fn, err := d.selectFlowNetwork(device.AutoMappingFlowNetworkName, "")
				if err != nil {
					return
				}
				cli := client.NewFlowClientCliFromFN(fn)
				networkName := networkModel.Name
				if boolean.IsFalse(fn.IsRemote) && boolean.IsFalse(fn.IsMasterSlave) {
					networkName = generateLocalNetworkName(networkName)
				}
				url := urls.SingularUrl(urls.NetworkNameUrl, networkName)
				_ = cli.DeleteQuery(url)
			}
			_, _ = d.DeleteDevice(device.UUID)
		}()
	}
	wg.Wait()
	query := d.DB.Delete(&networkModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeleteOneNetworkByArgs(args api.Args) (bool, error) {
	var networkModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.First(&networkModel).Error; err != nil {
		return false, err
	}
	query = query.Delete(&networkModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) getPluginConf(body *model.Network) compat.Info {
	var pluginConf *model.PluginConf
	query := d.DB.Where("uuid = ?", body.PluginConfId).First(&pluginConf)
	if query.Error != nil {
		return compat.Info{}
	}
	info := d.PluginManager.PluginInfo(pluginConf.ModulePath)
	return info
}

func (d *GormDatabase) SyncNetworks(args api.Args) ([]*interfaces.SyncModel, error) {
	networks, _ := d.GetNetworks(args)
	var outputs []*interfaces.SyncModel
	channel := make(chan *interfaces.SyncModel)
	defer close(channel)
	for _, network := range networks {
		go d.syncNetwork(network.UUID, args, channel)
	}
	for range networks {
		outputs = append(outputs, <-channel)
	}
	return outputs, nil
}

func (d *GormDatabase) syncNetwork(networkUUID string, args api.Args, channel chan *interfaces.SyncModel) {
	// This is for syncing child descendants
	syncModels, err := d.SyncNetworkDevices(networkUUID, args)
	output := interfaces.SyncModel{UUID: networkUUID, IsError: false}
	if err != nil {
		output = interfaces.SyncModel{UUID: networkUUID, IsError: true, Message: nstring.New(err.Error())}
	}
	for _, syncModel := range syncModels {
		if syncModel.IsError {
			output = interfaces.SyncModel{UUID: networkUUID, IsError: true, Message: syncModel.Message}
		}
	}
	channel <- &output
}

func (d *GormDatabase) SyncNetworkDevices(uuid string, args api.Args) ([]*interfaces.SyncModel, error) {
	network, _ := d.GetNetwork(uuid, args)
	var outputs []*interfaces.SyncModel
	if network == nil {
		return nil, errors.New("no network")
	}
	channel := make(chan *interfaces.SyncModel)
	defer close(channel)
	for _, device := range network.Devices {
		go d.syncDevice(device, args, channel)
	}
	for range network.Devices {
		outputs = append(outputs, <-channel)
	}
	return outputs, nil
}

func (d *GormDatabase) syncDevice(device *model.Device, args api.Args, channel chan *interfaces.SyncModel) {
	output := interfaces.SyncModel{UUID: device.UUID, IsError: false}
	if boolean.IsTrue(device.CreatedFromAutoMapping) {
		device.Connection = connection.Connected.String()
		device.Message = nstring.NotAvailable
		fnc, err := d.selectFlowNetworkClone(device.AutoMappingFlowNetworkName, "")
		if err != nil {
			device.Connection = connection.Broken.String()
			device.Message = "flow network clone not found"
		} else {
			network, _ := d.GetNetworkByDeviceUUID(device.UUID, api.Args{})
			cli := client.NewFlowClientCliFromFNC(fnc)
			rawDevice, err := cli.GetQueryMarshal(urls.SingularUrl(urls.DeviceNameUrl, fmt.Sprintf("%s/%s",
				strings.Replace(network.Name, "mapping_", "", -1), device.Name)), model.Device{})
			if err != nil {
				device.Connection = connection.Broken.String()
				device.Message = err.Error()
			} else {
				if boolean.IsFalse(rawDevice.(*model.Device).AutoMappingEnable) {
					_, _ = d.DeleteDevice(device.UUID)
					channel <- &output
					return
				}
			}
		}
		_ = d.UpdateDeviceErrors(device.UUID, device)
	}
	// This is for syncing child descendants
	syncModels, err := d.SyncDevicePoints(device.UUID, args)
	if err != nil {
		output = interfaces.SyncModel{UUID: device.UUID, IsError: true, Message: nstring.New(err.Error())}
	}
	for _, syncModel := range syncModels {
		if syncModel.IsError {
			output = interfaces.SyncModel{UUID: device.UUID, IsError: true, Message: syncModel.Message}
		}
	}
	channel <- &output
}
