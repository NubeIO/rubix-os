package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/plugin/compat"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

func (d *GormDatabase) GetNetworksTransaction(db *gorm.DB, args api.Args) ([]*model.Network, error) {
	var networksModel []*model.Network
	query := buildNetworkQueryTransaction(db, args)
	if err := query.Find(&networksModel).Error; err != nil {
		return nil, err
	}
	return networksModel, nil
}

func (d *GormDatabase) GetNetworks(args api.Args) ([]*model.Network, error) {
	return d.GetNetworksTransaction(d.DB, args)
}

func (d *GormDatabase) GetNetwork(uuid string, args api.Args) (*model.Network, error) {
	var networkModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&networkModel).Error; err != nil {
		return nil, err
	}
	return networkModel, nil
}

func (d *GormDatabase) GetOneNetworkByArgsTransaction(db *gorm.DB, args api.Args) (*model.Network, error) {
	var networkModel *model.Network
	query := buildNetworkQueryTransaction(db, args)
	if err := query.First(&networkModel).Error; err != nil {
		return nil, err
	}
	return networkModel, nil
}

func (d *GormDatabase) GetOneNetworkByArgs(args api.Args) (*model.Network, error) {
	return d.GetOneNetworkByArgsTransaction(d.DB, args)
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
			plugin, err := d.GetPluginByPathTransaction(db, body.PluginPath)
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

func (d *GormDatabase) UpdateNetworkTransaction(db *gorm.DB, uuid string, body *model.Network, checkAm bool) (*model.Network, error) {
	var networkModel *model.Network
	query := db.Where("uuid = ?", uuid).First(&networkModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if boolean.IsTrue(networkModel.CreatedFromAutoMapping) && checkAm {
		return nil, errors.New("can't update auto-mapped network")
	}
	if err := updateTagsTransaction(db, &networkModel, body.Tags); err != nil {
		return nil, err
	}
	body.Name = strings.TrimSpace(body.Name)
	query = db.Model(&networkModel).Select("*").Updates(&body)
	if query.Error != nil {
		return nil, query.Error
	}
	return networkModel, nil
}

func (d *GormDatabase) UpdateNetwork(uuid string, body *model.Network) (*model.Network, error) {
	return d.UpdateNetworkTransaction(d.DB, uuid, body, true)
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
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.Network{})
	go d.PublishPointsList("")
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeleteOneNetworkByArgs(args api.Args) (bool, error) {
	var networkModel *model.Network
	query := d.buildNetworkQuery(args).Delete(&networkModel)
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

func (d *GormDatabase) SyncNetworks() error {
	networks, err := d.GetNetworks(api.Args{WithDevices: true, WithPoints: true, WithPriority: true, WithTags: true, WithMetaTags: true})
	var firstErr error
	if err != nil {
		return err
	}
	uniqueAutoMappingFlowNetworkNames := GetUniqueAutoMappingFlowNetworkNames(networks)
	for _, fnName := range uniqueAutoMappingFlowNetworkNames {
		err = d.CreateNetworksAutoMappings(fnName, networks, interfaces.Network)
		if err != nil {
			log.Error("Auto mapping error: ", err)
		}
	}
	return firstErr
}

func (d *GormDatabase) SyncNetworkDevices(uuid string) error {
	return nil
	network, err := d.GetNetwork(uuid, api.Args{WithDevices: true, WithPoints: true, WithPriority: true, WithTags: true, WithMetaTags: true})
	if err != nil {
		return err
	}
	networks := make([]*model.Network, 0)
	networks = append(networks, network)
	return d.CreateNetworksAutoMappings(network.AutoMappingFlowNetworkName, networks, interfaces.Device)
}

func GetUniqueAutoMappingFlowNetworkNames(networks []*model.Network) []string {
	uniqueAutoMappingFlowNetworkNamesMap := make(map[string]struct{})
	var uniqueAutoMappingFlowNetworkNames []string

	for _, network := range networks {
		if _, ok := uniqueAutoMappingFlowNetworkNamesMap[network.AutoMappingFlowNetworkName]; !ok {
			uniqueAutoMappingFlowNetworkNamesMap[network.AutoMappingFlowNetworkName] = struct{}{}
			uniqueAutoMappingFlowNetworkNames = append(uniqueAutoMappingFlowNetworkNames, network.AutoMappingFlowNetworkName)
		}
	}

	return uniqueAutoMappingFlowNetworkNames
}
