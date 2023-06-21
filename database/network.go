package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/plugin/compat"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	name, err := validateName(body.Name)
	if err != nil {
		return nil, err
	}
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Network)
	body.Name = name
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
	if body.HistoryEnable == nil {
		body.HistoryEnable = boolean.NewFalse()
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
	name, err := validateName(body.Name)
	if err != nil {
		return nil, err
	}
	var networkModel *model.Network
	query := db.Where("uuid = ?", uuid).First(&networkModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if err := updateTagsTransaction(db, &networkModel, body.Tags); err != nil {
		return nil, err
	}
	body.Name = name
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

func (d *GormDatabase) DeleteNetworkClonesByHostUUIDTransaction(db *gorm.DB, hostUUID string) (bool, error) {
	query := db.Where("host_uuid = ?", hostUUID).Where("is_clone IS TRUE").Delete(&model.Network{})
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

func (d *GormDatabase) CreateBulkNetworksTransaction(db *gorm.DB, networks []*model.Network) (bool, error) {
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(networks, 1000).Error; err != nil {
		log.Error("Issue on creating bulk networks")
		return false, err
	}
	return true, nil
}
