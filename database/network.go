package database

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/compat"
	"github.com/NubeDev/flow-framework/utils"
	"time"
)

const ip = "IpConnection"
const serial = "SerialConnection"

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

// GetNetworkByPlugin returns the network for the given id or nil.
func (d *GormDatabase) GetNetworkByPlugin(pluginUUID string, args api.Args) (*model.Network, error) {
	var networkModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("plugin_conf_id = ? ", pluginUUID).First(&networkModel).Error; err != nil {
		return nil, err
	}
	return networkModel, nil
}

// GetNetworksByPlugin returns the network for the given id or nil.
func (d *GormDatabase) GetNetworksByPlugin(pluginUUID string, args api.Args) ([]*model.Network, error) {
	var networksModel []*model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("plugin_conf_id = ? ", pluginUUID).Find(&networksModel).Error; err != nil {
		return nil, err
	}
	return networksModel, nil
}

// GetNetworksByName returns the network for the given id or nil.
func (d *GormDatabase) GetNetworksByName(name string, args api.Args) ([]*model.Network, error) {
	var networksModel []*model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Find(&networksModel).Where("name = ? ", name).Error; err != nil {
		return nil, err
	}
	return networksModel, nil
}

// CreateNetwork creates a device.
func (d *GormDatabase) CreateNetwork(body *model.Network) (*model.Network, error) {
	body.UUID = utils.MakeTopicUUID(model.ThingClass.Network)
	body.Name = nameIsNil(body.Name)
	body.ThingClass = model.ThingClass.Network
	body.CommonEnable.Enable = utils.NewTrue()
	body.CommonFault.InFault = true
	body.CommonFault.MessageLevel = model.MessageLevel.NoneCritical
	body.CommonFault.MessageCode = model.CommonFaultCode.PluginNotEnabled
	body.CommonFault.Message = model.CommonFaultMessage.PluginNotEnabled
	body.CommonFault.LastFail = time.Now().UTC()
	body.CommonFault.LastOk = time.Now().UTC()
	t := body.TransportType
	s := model.TransType.Serial
	host := model.TransType.IP
	transport, err := checkTransport(body.TransportType) //set to ip by default
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
		switch t {
		case s:
			if body.SerialConnection == nil {
				body.SerialConnection = &model.SerialConnection{}
			}
			body.SerialConnection.UUID = utils.MakeTopicUUID(model.TransType.Serial)
		case host:
			if body.IpConnection == nil {
				body.IpConnection = &model.IpConnection{}
			}
			body.IpConnection.UUID = utils.MakeTopicUUID(model.TransType.IP)
		}
	} else {
		return nil, errors.New("provide a plugin name ie: system, lora, modbus, lorawan, bacnet")
	}

	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateNetwork(uuid string, body *model.Network) (*model.Network, error) {
	var networkModel *model.Network
	query := d.DB.Where("uuid = ?", uuid).Preload(serial).Preload(ip).Find(&networkModel)
	if query.Error != nil {
		return nil, query.Error
	}
	switch networkModel.TransportType {
	case model.TransType.Serial:
		d.DB.Model(&networkModel.SerialConnection).Updates(body.SerialConnection)
	case model.TransType.IP:
		d.DB.Model(&networkModel.IpConnection).Updates(body.IpConnection)
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&networkModel, body.Tags); err != nil {
			return nil, err
		}
	}
	query = d.DB.Model(&networkModel).Updates(&body)
	if query.Error != nil {
		return nil, query.Error
	}
	t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsUpdated, networkModel.PluginConfId, networkModel.UUID)
	d.Bus.RegisterTopic(t)
	err := d.Bus.Emit(eventbus.CTX(), t, networkModel)
	if err != nil {
		return nil, errors.New("error on network eventbus")
	}
	return networkModel, nil

}

func (d *GormDatabase) DeleteNetwork(uuid string) (bool, error) {
	var networkModel *model.Network
	query := d.DB.Where("uuid = ? ", uuid).Delete(&networkModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (d *GormDatabase) DropNetworks() (bool, error) {
	var networkModel *model.Network
	query := d.DB.Where("1 = 1").Delete(&networkModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (d *GormDatabase) getPluginConf(body *model.Network) compat.Info {
	var pluginConf *model.PluginConf
	query := d.DB.Where("uuid = ?", body.PluginConfId).Find(&pluginConf)
	if query.Error != nil {
		return compat.Info{}
	}
	info := d.PluginManager.PluginInfo(pluginConf.ModulePath)
	return info
}
