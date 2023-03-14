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
	log "github.com/sirupsen/logrus"
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
		Select("InFault", "MessageLevel", "MessageCode", "Message", "LastFail", "InSync").
		Updates(&body).
		Error
}

func (d *GormDatabase) UpdateNetworkConnectionErrors(uuid string, network *model.Network) error {
	return d.DB.Model(&model.Network{}).
		Where("uuid = ?", uuid).
		Select("Connection", "ConnectionMessage").
		Updates(&network).
		Error
}

func (d *GormDatabase) DeleteNetwork(uuid string) (bool, error) {
	networkModel, err := d.GetNetwork(uuid, api.Args{WithDevices: true})
	if err != nil {
		return false, err
	}
	if boolean.IsTrue(networkModel.AutoMappingEnable) {
		fn, err := d.GetOneFlowNetworkByArgs(api.Args{Name: nstring.New(networkModel.AutoMappingFlowNetworkName)})
		if err != nil {
			log.Errorf("failed to find flow network with name %s", networkModel.AutoMappingFlowNetworkName)
			return false, err
		}
		cli := client.NewFlowClientCliFromFN(fn)
		var wg sync.WaitGroup
		for _, device := range networkModel.Devices {
			wg.Add(1)
			go func(device *model.Device) {
				defer wg.Done()
				networkName := getAutoMappedNetworkName(networkModel.Name, boolean.IsFalse(fn.IsRemote))
				url := urls.SingularUrl(urls.NetworkNameUrl, networkName)
				_ = cli.DeleteQuery(url)
				_, _ = d.DeleteDevice(device.UUID)
			}(device)
		}
		wg.Wait()
	}
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
	d.removeUnlinkedAutoMappedNetworks()
	d.removeUnlinkedAutoMappedDevices()
	d.removeUnlinkedAutoMappedStreams()
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
	syncModels, err := d.SyncNetworkDevices(networkUUID, false, args)
	output := interfaces.SyncModel{UUID: networkUUID, IsError: false}
	if err != nil {
		output = interfaces.SyncModel{UUID: networkUUID, IsError: true, Message: nstring.New(err.Error())}
	}
	for _, syncModel := range syncModels {
		if syncModel.IsError {
			output = interfaces.SyncModel{UUID: networkUUID, IsError: true, Message: syncModel.Message}
		}
	}
	networkModel := model.Network{}
	if output.IsError {
		networkModel.Connection = connection.Broken.String()
		networkModel.ConnectionMessage = output.Message
	} else {
		networkModel.Connection = connection.Connected.String()
		networkModel.ConnectionMessage = nstring.New(nstring.NotAvailable)
	}
	_ = d.UpdateNetworkConnectionErrors(output.UUID, &networkModel)
	channel <- &output
}

func (d *GormDatabase) SyncNetworkDevices(uuid string, removeUnlinked bool, args api.Args) ([]*interfaces.SyncModel, error) {
	if removeUnlinked {
		d.removeUnlinkedAutoMappedDevices()
		d.removeUnlinkedAutoMappedStreams()
	}
	network, _ := d.GetNetwork(uuid, args)
	var outputs []*interfaces.SyncModel
	if network == nil {
		return nil, errors.New("no network")
	}
	channel := make(chan *interfaces.SyncModel)
	defer close(channel)
	for _, device := range network.Devices {
		go d.syncDevice(network.AutoMappingFlowNetworkName, device, args, channel)
	}
	for range network.Devices {
		outputs = append(outputs, <-channel)
	}
	return outputs, nil
}

func (d *GormDatabase) syncDevice(flowNetworkName string, device *model.Device, args api.Args, channel chan *interfaces.SyncModel) {
	output := interfaces.SyncModel{UUID: device.UUID, IsError: false}
	if boolean.IsTrue(device.CreatedFromAutoMapping) {
		device.Connection = connection.Connected.String()
		device.ConnectionMessage = nstring.New(nstring.NotAvailable)
		fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{Name: nstring.New(flowNetworkName)})
		if err != nil {
			device.Connection = connection.Broken.String()
			device.ConnectionMessage = nstring.New("flow network clone not found")
		} else {
			network, _ := d.GetNetworkByDeviceUUID(device.UUID, api.Args{})
			cli := client.NewFlowClientCliFromFNC(fnc)
			rawDevice, err := cli.GetQueryMarshal(urls.SingularUrl(urls.DeviceNameUrl, fmt.Sprintf("%s/%s",
				strings.Replace(network.Name, "mapping_", "", -1), device.Name)), model.Device{})
			if err != nil {
				device.Connection = connection.Broken.String()
				device.ConnectionMessage = nstring.New(err.Error())
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
	syncModels, err := d.SyncDevicePoints(device.UUID, false, args)
	if err != nil {
		output = interfaces.SyncModel{UUID: device.UUID, IsError: true, Message: nstring.New(err.Error())}
	}
	for _, syncModel := range syncModels {
		if syncModel.IsError {
			output = interfaces.SyncModel{UUID: device.UUID, IsError: true, Message: syncModel.Message}
		}
	}
	deviceModel := model.Device{}
	if output.IsError {
		deviceModel.Connection = connection.Broken.String()
		deviceModel.ConnectionMessage = output.Message
	} else {
		deviceModel.Connection = connection.Connected.String()
		deviceModel.ConnectionMessage = nstring.New(nstring.NotAvailable)
	}
	_ = d.UpdateDeviceConnectionErrors(output.UUID, &deviceModel)
	channel <- &output
}

func (d *GormDatabase) removeUnlinkedAutoMappedStreams() {
	streams, err := d.GetStreams(api.Args{})
	if err != nil {
		log.Errorf(err.Error())
	}
	for _, stream := range streams {
		if boolean.IsTrue(stream.CreatedFromAutoMapping) {
			parts := strings.Split(stream.Name, ":")
			if len(parts) == 2 {
				device, _ := d.GetDeviceByName(parts[0], parts[1], api.Args{})
				if device == nil {
					_, err := d.DeleteStream(stream.UUID)
					if err != nil {
						log.Errorf(err.Error())
					}
				} else {
					network, _ := d.GetNetwork(device.NetworkUUID, api.Args{})
					if boolean.IsFalse(network.AutoMappingEnable) || boolean.IsFalse(device.AutoMappingEnable) {
						_, err := d.DeleteStream(stream.UUID)
						if err != nil {
							log.Errorf(err.Error())
						}
					}
				}
			}
		}
	}
}

func (d *GormDatabase) removeUnlinkedAutoMappedNetworks() {
	networks, err := d.GetNetworks(api.Args{})
	if err != nil {
		log.Errorf(err.Error())
	}
	for _, network := range networks {
		if boolean.IsTrue(network.CreatedFromAutoMapping) {
			fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{Name: nstring.New(network.AutoMappingFlowNetworkName)})
			if err != nil {
				_, err := d.DeleteNetwork(network.UUID)
				if err != nil {
					log.Errorf(err.Error())
				}
			}
			cli := client.NewFlowClientCliFromFNC(fnc)
			networkName := getAutoMappedOriginalNetworkName(network.Name, boolean.IsFalse(fnc.IsRemote))
			remoteNetwork, connectionErr, _ := cli.GetNetworkByName(networkName)
			if connectionErr != nil {
				network.Connection = connection.Broken.String()
				network.ConnectionMessage = nstring.New(err.Error())
				_ = d.UpdateNetworkConnectionErrors(network.UUID, network)
			} else if remoteNetwork == nil ||
				boolean.IsFalse(remoteNetwork.AutoMappingEnable) ||
				remoteNetwork.AutoMappingFlowNetworkName != network.AutoMappingFlowNetworkName {
				_, err := d.DeleteNetwork(network.UUID)
				if err != nil {
					log.Errorf(err.Error())
				}
			}
		}
	}
}

func (d *GormDatabase) removeUnlinkedAutoMappedDevices() {
	devices, err := d.GetDevices(api.Args{})
	if err != nil {
		log.Errorf(err.Error())
	}
	for _, device := range devices {
		if boolean.IsTrue(device.CreatedFromAutoMapping) {
			network, err := d.GetNetwork(device.NetworkUUID, api.Args{})
			fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{Name: nstring.New(network.AutoMappingFlowNetworkName)})
			if err != nil {
				_, err := d.DeleteDevice(device.UUID)
				if err != nil {
					log.Errorf(err.Error())
				}
			}
			cli := client.NewFlowClientCliFromFNC(fnc)
			networkName := getAutoMappedOriginalNetworkName(network.Name, boolean.IsFalse(fnc.IsRemote))
			remoteDevice, connectionErr, _ := cli.GetDeviceByName(networkName, device.Name)
			if connectionErr != nil {
				device.Connection = connection.Broken.String()
				device.ConnectionMessage = nstring.New(err.Error())
				_ = d.UpdateDeviceConnectionErrors(device.UUID, device)
			} else if remoteDevice == nil || boolean.IsFalse(remoteDevice.AutoMappingEnable) {
				_, err := d.DeleteDevice(device.UUID)
				if err != nil {
					log.Errorf(err.Error())
				}
			}
		}
	}
}
