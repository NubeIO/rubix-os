package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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

func (d *GormDatabase) CreateDevice(body *model.Device) (*model.Device, error) {
	var net *model.Network
	query := d.DB.Where("uuid = ? ", body.NetworkUUID).First(&net)
	if query.Error != nil {
		return nil, query.Error
	}
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Device)
	body.ThingClass = model.ThingClass.Device
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, query.Error
}

func (d *GormDatabase) UpdateDevice(uuid string, body *model.Device, fromPlugin bool) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.DB.Where("uuid = ?", uuid).First(&deviceModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&deviceModel, body.Tags); err != nil {
			return nil, err
		}
	}
	body.ThingClass = model.ThingClass.Device
	if err := d.DB.Model(&deviceModel).Select("*").Updates(body).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
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

func (d *GormDatabase) DeleteDevice(uuid string) (bool, error) {
	var deviceModel *model.Device
	query := d.DB.Where("uuid = ? ", uuid).Delete(&deviceModel)
	return d.deleteResponseBuilder(query)
}
