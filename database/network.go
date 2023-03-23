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

func (d *GormDatabase) UpdateNetworkConnectionErrorsByName(name string, network *model.Network) error {
	return d.DB.Model(&model.Network{}).
		Where("name = ?", name).
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
				networkName := getAutoMappedNetworkName(fn.Name, networkModel.Name)
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

func (d *GormDatabase) SyncNetworks(args api.Args) ([]*interfaces.AutoMappingNetworkError, error) {
	d.removeUnlinkedAutoMappedNetworks()
	d.removeUnlinkedAutoMappedDevices()
	d.removeUnlinkedAutoMappedPoints()
	d.removeUnlinkedAutoMappedStreams()
	networks, err := d.GetNetworks(args)
	if err != nil {
		return nil, err
	}
	var outputs []*interfaces.AutoMappingNetworkError
	for _, network := range networks {
		output, _ := d.SyncNetworkDevices(network.UUID, network, false, args) // we never get err on here
		outputs = append(outputs, output)
	}
	return outputs, nil
}

func (d *GormDatabase) SyncNetworkDevices(uuid string, network *model.Network, removeUnlinked bool, args api.Args) (*interfaces.AutoMappingNetworkError, error) {
	if removeUnlinked {
		d.removeUnlinkedAutoMappedDevices()
		d.removeUnlinkedAutoMappedPoints()
		d.removeUnlinkedAutoMappedStreams()
	}
	if network == nil {
		network, _ = d.GetNetwork(uuid, args)
	}
	if network == nil {
		return nil, errors.New("network doesn't exist")
	}
	return d.SyncDevicePoints(uuid, network, false, args)
}

func (d *GormDatabase) removeUnlinkedAutoMappedStreams() {
	streams, err := d.GetStreams(api.Args{})
	if err != nil {
		log.Errorf(err.Error())
	}
	for _, stream := range streams {
		if boolean.IsTrue(stream.CreatedFromAutoMapping) {
			parts := strings.Split(stream.Name, ":")
			if len(parts) == 3 {
				fn, _ := d.GetOneFlowNetworkByArgs(api.Args{Name: nstring.New(parts[0])})
				device, _ := d.GetDeviceByName(parts[1], parts[2], api.Args{})
				if device == nil || fn == nil {
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
			} else {
				_, err := d.DeleteStream(stream.UUID)
				if err != nil {
					log.Errorf(err.Error())
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
			networkName := getAutoMappedOriginalNetworkName(fnc.Name, network.Name)
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
			networkName := getAutoMappedOriginalNetworkName(fnc.Name, network.Name)
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

func (d *GormDatabase) removeUnlinkedAutoMappedPoints() {
	points, err := d.GetPoints(api.Args{})
	if err != nil {
		log.Errorf(err.Error())
	}
	for _, point := range points {
		if boolean.IsTrue(point.CreatedFromAutoMapping) {
			device, _ := d.GetDevice(point.DeviceUUID, api.Args{})
			network, _ := d.GetNetwork(device.NetworkUUID, api.Args{})
			fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{Name: nstring.New(network.AutoMappingFlowNetworkName)})
			if err != nil {
				_, err := d.DeletePoint(point.UUID, false)
				if err != nil {
					log.Errorf(err.Error())
				}
			}
			cli := client.NewFlowClientCliFromFNC(fnc)
			networkName := getAutoMappedOriginalNetworkName(fnc.Name, network.Name)
			remotePoint, connectionErr, _ := cli.GetPointByName(networkName, device.Name, point.Name)
			if connectionErr != nil {
				point.Connection = connection.Broken.String()
				point.ConnectionMessage = nstring.New(err.Error())
				_ = d.UpdatePointConnectionErrors(point.UUID, point)
			} else if remotePoint == nil || boolean.IsFalse(remotePoint.AutoMappingEnable) {
				_, err := d.DeletePoint(point.UUID, false)
				if err != nil {
					log.Errorf(err.Error())
				}
			}
		}
	}
	go d.PublishPointsList("")
}
