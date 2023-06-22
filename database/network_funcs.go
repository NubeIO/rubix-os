package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/interfaces"
	log "github.com/sirupsen/logrus"
	"time"
)

// GetNetworkByPluginName returns the network for the given id or nil.
func (d *GormDatabase) GetNetworkByPluginName(name string, args api.Args) (*model.Network, error) {
	var networkModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("plugin_path = ? ", name).First(&networkModel).Error; err != nil {
		return nil, err
	}
	return networkModel, nil
}

// GetNetworksByPluginName returns the network for the given id or nil.
func (d *GormDatabase) GetNetworksByPluginName(name string, args api.Args) ([]*model.Network, error) {
	var networksModel []*model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("plugin_path = ? ", name).Find(&networksModel).Error; err != nil {
		return nil, err
	}
	return networksModel, nil
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

// GetNetworkByName returns the network for the given id or nil.
func (d *GormDatabase) GetNetworkByName(name string, args api.Args) (*model.Network, error) {
	var networksModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("name = ? ", name).First(&networksModel).Error; err != nil {
		return nil, err
	}
	return networksModel, nil
}

// GetNetworkByPoint returns a network by passing in the pointUUID.
func (d *GormDatabase) GetNetworkByPoint(point *model.Point, args api.Args) (network *model.Network, err error) {
	device, err := d.GetDeviceByPoint(point)
	if err != nil {
		return nil, err
	}
	network, err = d.GetNetwork(device.NetworkUUID, args)
	if err != nil {
		return nil, err
	}
	return
}

// GetNetworkByPointUUID returns a network by passing in the pointUUID.
func (d *GormDatabase) GetNetworkByPointUUID(pntUUID string, args api.Args) (network *model.Network, err error) {
	device, err := d.GetDeviceByPointUUID(pntUUID)
	if err != nil {
		return nil, err
	}
	network, err = d.GetNetwork(device.NetworkUUID, args)
	if err != nil {
		return nil, err
	}
	return
}

// GetNetworkByDeviceUUID returns a network by passing in the device UUID.
func (d *GormDatabase) GetNetworkByDeviceUUID(devUUID string, args api.Args) (network *model.Network, err error) {
	device, err := d.GetDevice(devUUID, args)
	if err != nil {
		return nil, err
	}
	return d.GetNetwork(device.NetworkUUID, args)
}

// SetErrorsForAllDevicesOnNetwork sets the fault/error properties of all devices for a specific network. Optional to set the points of each device also.
// messageLevel = model.MessageLevel
// messageCode = model.CommonFaultCode
func (d *GormDatabase) SetErrorsForAllDevicesOnNetwork(networkUUID string, message string, messageLevel string, messageCode string, doPoints bool) error {
	network, err := d.GetNetwork(networkUUID, api.Args{WithDevices: true, WithPoints: doPoints})
	if err != nil {
		return err
	}
	for _, device := range network.Devices {
		device.CommonFault.InFault = true
		device.CommonFault.MessageLevel = messageLevel
		device.CommonFault.MessageCode = messageCode
		device.CommonFault.Message = message
		device.CommonFault.LastFail = time.Now().UTC()
		err = d.UpdateDeviceErrors(device.UUID, device)
		if err != nil {
			log.Infof("setErrorsForAllDevicesOnNetwork() Error: %s\n", err.Error())
		}
		if doPoints {
			err = d.SetErrorsForAllPointsOnDevice(device.UUID, message, messageLevel, messageCode)
		}
	}
	return nil
}

// ClearErrorsForAllDevicesOnNetwork clears the fault/error properties of all devices for a specific network. Optional to clear the points of each device also.
func (d *GormDatabase) ClearErrorsForAllDevicesOnNetwork(networkUUID string, doPoints bool) error {
	network, err := d.GetNetwork(networkUUID, api.Args{WithDevices: true, WithPoints: doPoints})
	if network != nil && err != nil {
		return err
	}
	for _, device := range network.Devices {
		device.CommonFault.InFault = false
		device.CommonFault.MessageLevel = model.MessageLevel.Normal
		device.CommonFault.MessageCode = model.CommonFaultCode.Ok
		device.CommonFault.Message = ""
		device.CommonFault.LastOk = time.Now().UTC()
		err = d.UpdateDeviceErrors(device.UUID, device)
		if err != nil {
			log.Infof("clearErrorsForAllDevicesOnNetwork() Error: %s\n", err.Error())
		}
		if doPoints {
			err = d.ClearErrorsForAllPointsOnDevice(device.UUID)
		}
	}
	return nil
}

func (d *GormDatabase) DeleteNetworkByName(name string, args api.Args) (bool, error) {
	var networksModel *model.Network
	query := d.buildNetworkQuery(args)
	if err := query.Where("name = ? ", name).First(&networksModel).Error; err != nil {
		return false, err
	}
	query = query.Delete(networksModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) GetPublishPointList() ([]*interfaces.PublishPointList, error) {
	var publishPointListModel []*interfaces.PublishPointList
	query := d.DB.Table("networks").
		Select("networks.plugin_path, networks.name AS network_name, devices.name AS device_name, " +
			"points.uuid AS point_uuid, points.name AS point_name").
		Joins("INNER JOIN devices ON devices.network_uuid = networks.uuid").
		Joins("INNER JOIN points ON points.device_uuid = devices.uuid").
		Scan(&publishPointListModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return publishPointListModel, nil
}

func (d *GormDatabase) GetNetworksTagsForPostgresSync() ([]*interfaces.NetworkTagForPostgresSync, error) {
	var networkTagsForPostgresModel []*interfaces.NetworkTagForPostgresSync
	query := d.DB.Table("networks_tags").
		Select("networks.source_uuid AS network_uuid, networks_tags.tag_tag AS tag").
		Joins("INNER JOIN networks ON networks.uuid = networks_tags.network_uuid").
		Where("IFNULL(networks.source_uuid,'') != ''").
		Scan(&networkTagsForPostgresModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return networkTagsForPostgresModel, nil
}
