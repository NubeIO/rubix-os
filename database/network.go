package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/plugin/compat"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"sync"
)

func (d *GormDatabase) GetNetworks(args api.Args) ([]*model.Network, error) {
	var networksModel []*model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Find(&networksModel).Error; err != nil {
		return nil, err
	}
	return networksModel, nil
}

func (d *GormDatabase) GetNetwork(uuid string, args api.Args) (*model.Network, error) {
	var networkModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&networkModel).Error; err != nil {
		return nil, err
	}
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

func (d *GormDatabase) CreateNetwork(body *model.Network, fromPlugin bool) (*model.Network, error) {
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Network)
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
	if err = d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateNetwork(uuid string, body *model.Network, fromPlugin bool) (*model.Network, error) {
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
	query = d.DB.Model(&networkModel).Select("*").Updates(&body)
	if query.Error != nil {
		return nil, query.Error
	}
	_ = d.syncAfterUpdateNetwork(body.UUID, api.Args{WithTags: true, WithDevices: true,
		WithPoints: true})
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
	var aType = api.ArgsType
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
				fn, err := d.selectFlowNetwork(device.AutoMappingFlowNetworkName, device.AutoMappingFlowNetworkUUID)
				if err != nil {
					return
				}
				cli := client.NewFlowClientCliFromFN(fn)
				url := urls.SingularUrlByArg(urls.NetworkUrl, aType.AutoMappingUUID, networkModel.UUID)
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

func (d *GormDatabase) syncAfterUpdateNetwork(uuid string, args api.Args) error {
	network, err := d.GetNetwork(uuid, args)
	if err != nil {
		return err
	}
	if args.WithDevices {
		_, err = d.SyncNetworkDevices(network.UUID, args)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *GormDatabase) SyncNetworks(args api.Args) ([]*interfaces.SyncModel, error) {
	networks, _ := d.GetNetworks(args)
	var outputs []*interfaces.SyncModel
	channel := make(chan *interfaces.SyncModel)
	defer close(channel)
	for _, network := range networks {
		go d.syncNetwork(network, channel)
	}
	for range networks {
		outputs = append(outputs, <-channel)
	}
	return outputs, nil
}

func (d *GormDatabase) syncNetwork(network *model.Network, channel chan *interfaces.SyncModel) {
	_, err := d.UpdateNetwork(network.UUID, network, false)
	var output interfaces.SyncModel
	if err != nil {
		output = interfaces.SyncModel{UUID: network.UUID, IsError: true, Message: nstring.New(err.Error())}
	} else {
		output = interfaces.SyncModel{UUID: network.UUID, IsError: false}
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
	err := d.syncAfterCreateUpdateDevice(device.UUID, args)
	var output interfaces.SyncModel
	if err != nil {
		output = interfaces.SyncModel{UUID: device.UUID, IsError: true, Message: nstring.New(err.Error())}
	} else {
		output = interfaces.SyncModel{UUID: device.UUID, IsError: false}
	}
	channel <- &output
}
