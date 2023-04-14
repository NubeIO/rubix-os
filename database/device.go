package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm"
	"strings"
)

func (d *GormDatabase) GetDevices(args api.Args) ([]*model.Device, error) {
	var devicesModel []*model.Device
	query := d.buildDeviceQuery(args)
	if err := query.Find(&devicesModel).Error; err != nil {
		return nil, err
	}
	return devicesModel, nil
}

func (d *GormDatabase) GetDevice(uuid string, args api.Args) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&deviceModel).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
}

func (d *GormDatabase) CreateDeviceTransaction(db *gorm.DB, body *model.Device, checkAm bool) (*model.Device, error) {
	var network *model.Network
	query := db.Where("uuid = ? ", body.NetworkUUID).First(&network)
	if query.Error != nil {
		return nil, fmt.Errorf("no such parent network with uuid %s", body.NetworkUUID)
	}
	if boolean.IsTrue(network.CreatedFromAutoMapping) && checkAm {
		return nil, errors.New("can't create a device for the auto-mapped network")
	}
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Device)
	body.Name = strings.TrimSpace(body.Name)
	body.ThingClass = model.ThingClass.Device
	if err := db.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, query.Error
}

func (d *GormDatabase) CreateDevice(body *model.Device) (*model.Device, error) {
	return d.CreateDeviceTransaction(d.DB, body, true)
}

func (d *GormDatabase) UpdateDeviceTransaction(db *gorm.DB, uuid string, body *model.Device, checkAm bool) (*model.Device, error) {
	var deviceModel *model.Device
	query := db.Where("uuid = ?", uuid).First(&deviceModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if boolean.IsTrue(deviceModel.CreatedFromAutoMapping) && checkAm {
		return nil, errors.New("can't update auto-mapped device")
	}
	if err := updateTagsTransaction(db, &deviceModel, body.Tags); err != nil {
		return nil, err
	}
	body.Name = strings.TrimSpace(body.Name)
	body.ThingClass = model.ThingClass.Device
	if err := db.Model(&deviceModel).Select("*").Updates(body).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
}

func (d *GormDatabase) UpdateDevice(uuid string, body *model.Device) (*model.Device, error) {
	return d.UpdateDeviceTransaction(d.DB, uuid, body, true)
}

// UpdateDeviceErrors will only update the CommonFault properties of the device, all other properties won't be updated
// Does not update `LastOk`
func (d *GormDatabase) UpdateDeviceErrors(uuid string, body *model.Device) error {
	return d.DB.Model(&body).
		Where("uuid = ?", uuid).
		Select("InFault", "MessageLevel", "MessageCode", "Message", "LastFail", "InSync").
		Updates(&body).
		Error
}

func UpdateDeviceConnectionErrorsTransaction(db *gorm.DB, uuid string, device *model.Device) error {
	return db.Model(&model.Device{}).
		Where("uuid = ?", uuid).
		Select("Connection", "ConnectionMessage").
		Updates(&device).
		Error
}

func (d *GormDatabase) UpdateDeviceConnectionErrors(uuid string, device *model.Device) error {
	return UpdateDeviceConnectionErrorsTransaction(d.DB, uuid, device)
}

func (d *GormDatabase) UpdateDeviceConnectionErrorsByName(name string, device *model.Device) error {
	return d.DB.Model(&model.Device{}).
		Where("name = ?", name).
		Select("Connection", "ConnectionMessage").
		Updates(&device).
		Error
}

func (d *GormDatabase) DeleteDevice(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.Device{})
	go d.PublishPointsList("")
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeleteOneDeviceByArgs(args api.Args) (bool, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args).Delete(&deviceModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) SyncDevicePoints(uuid string) error {
	device, err := d.GetDevice(uuid, api.Args{WithPoints: true, WithPriority: true, WithTags: true, WithMetaTags: true})
	if err != nil {
		return err
	}
	devices := make([]*model.Device, 0)
	devices = append(devices, device)

	network, _ := d.GetNetwork(device.NetworkUUID, api.Args{WithTags: true, WithMetaTags: true})
	network.Devices = devices // doing this for just to sync one device points
	networks := make([]*model.Network, 0)
	networks = append(networks, network)
	return d.CreateNetworksAutoMappings(network.AutoMappingFlowNetworkName, networks, interfaces.Point)
}
