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
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
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

func (d *GormDatabase) GetOneNetworkByArgs(args api.Args) (*model.Network, error) {
	var networkModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.First(&networkModel).Error; err != nil {
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

func (d *GormDatabase) CreateNetworkTransaction(db *gorm.DB, body *model.Network) (*model.Network, error) {
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
	if err = db.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) CreateNetwork(body *model.Network) (*model.Network, error) {
	return d.CreateNetworkTransaction(d.DB, body)
}

func (d *GormDatabase) UpdateNetworkTransaction(db *gorm.DB, uuid string, body *model.Network) (*model.Network, error) {
	var networkModel *model.Network
	query := db.Where("uuid = ?", uuid).First(&networkModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&networkModel, body.Tags); err != nil {
			return nil, err
		}
	}
	body.Name = strings.TrimSpace(body.Name)
	query = db.Model(&networkModel).Select("*").Updates(&body)
	if query.Error != nil {
		return nil, query.Error
	}
	return networkModel, nil
}

func (d *GormDatabase) UpdateNetwork(uuid string, body *model.Network) (*model.Network, error) {
	return d.UpdateNetworkTransaction(d.DB, uuid, body)
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

func UpdateNetworkConnectionErrorsTransaction(db *gorm.DB, uuid string, network *model.Network) error {
	return db.Model(&model.Network{}).
		Where("uuid = ?", uuid).
		Select("Connection", "ConnectionMessage").
		Updates(&network).
		Error
}

func (d *GormDatabase) UpdateNetworkConnectionErrors(uuid string, network *model.Network) error {
	return UpdateNetworkConnectionErrorsTransaction(d.DB, uuid, network)
}

func (d *GormDatabase) DeleteNetwork(uuid string) (bool, error) {
	networkModel, err := d.GetNetwork(uuid, api.Args{WithDevices: true})
	if err != nil {
		return false, fmt.Errorf("failed to get network: %w", err)
	}

	if boolean.IsTrue(networkModel.AutoMappingEnable) {
		var cli *client.FlowClient

		fn, err := d.GetOneFlowNetworkByArgs(api.Args{Name: nstring.New(networkModel.AutoMappingFlowNetworkName)})
		if err != nil {
			log.Errorf("failed to find flow network with name %s", networkModel.AutoMappingFlowNetworkName)
		} else {
			cli = client.NewFlowClientCliFromFN(fn)
		}

		if cli != nil {
			aType := api.ArgsType
			url := urls.SingularUrlByArg(urls.NetworksUrl, aType.AutoMappingUUID, networkModel.UUID)
			_ = cli.DeleteQuery(url)

			streams, _ := d.GetStreamByArgs(api.Args{AutoMappingNetworkUUID: &networkModel.UUID})
			if streams != nil {
				for _, stream := range streams { // todo: create bulk stream delete API
					url := urls.SingularUrlByArg(urls.StreamCloneUrl, aType.SourceUUID, stream.UUID)
					_ = cli.DeleteQuery(url)
				}
				d.DB.Delete(&streams)
			}
		}
	}

	if boolean.IsTrue(networkModel.CreatedFromAutoMapping) {
		d.DB.
			Where("auto_mapping_network_uuid = ? AND created_from_auto_mapping IS TRUE", networkModel.UUID).
			Delete(&model.StreamClone{})
	}

	query := d.DB.Delete(&networkModel)
	go d.PublishPointsList("")
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

func (d *GormDatabase) SyncNetworks(level interfaces.Level, args api.Args) error {
	networks, err := d.GetNetworks(args)
	if err != nil {
		return err
	}
	for _, network := range networks {
		err = d.SyncNetworkDevices(network.UUID, network, level, args)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *GormDatabase) SyncNetworkDevices(uuid string, network *model.Network, level interfaces.Level, args api.Args) error {
	if network == nil {
		network, _ = d.GetNetwork(uuid, args)
	}
	if network == nil {
		return errors.New("network doesn't exist")
	}
	return d.SyncDevicePoints(uuid, network, level, args)
}
