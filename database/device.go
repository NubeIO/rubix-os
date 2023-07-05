package database

import (
	"fmt"

	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetDevices(args argspkg.Args) ([]*model.Device, error) {
	var devicesModel []*model.Device
	query := d.buildDeviceQuery(args)
	if err := query.Find(&devicesModel).Error; err != nil {
		return nil, err
	}
	return devicesModel, nil
}

func (d *GormDatabase) GetDevice(uuid string, args argspkg.Args) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&deviceModel).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
}

func (d *GormDatabase) CreateDevice(body *model.Device) (*model.Device, error) {
	name, err := validateName(body.Name)
	if err != nil {
		return nil, err
	}
	var network *model.Network
	query := d.DB.Where("uuid = ? ", body.NetworkUUID).First(&network)
	if query.Error != nil {
		return nil, fmt.Errorf("no such parent network with uuid %s", body.NetworkUUID)
	}
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Device)
	body.Name = name
	body.ThingClass = model.ThingClass.Device
	if body.HistoryEnable == nil {
		body.HistoryEnable = boolean.NewTrue()
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, query.Error
}

func (d *GormDatabase) UpdateDevice(uuid string, body *model.Device) (*model.Device, error) {
	name, err := validateName(body.Name)
	if err != nil {
		return nil, err
	}
	var deviceModel *model.Device
	query := d.DB.Where("uuid = ?", uuid).First(&deviceModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if err := d.updateTags(&deviceModel, body.Tags); err != nil {
		return nil, err
	}
	body.Name = name
	body.ThingClass = model.ThingClass.Device
	if err := d.DB.Model(&deviceModel).Select("*").Updates(body).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
}

// UpdateDeviceErrors will only update the CommonFault properties of the device, all other properties won't be updated
// Does not update `LastOk`
func (d *GormDatabase) UpdateDeviceErrors(uuid string, body *model.Device) error {
	if body.InFault {
		return d.DB.Model(&body).
			Where("uuid = ?", uuid).
			Select("InFault", "MessageLevel", "MessageCode", "Message", "LastFail", "InSync").
			Updates(&body).
			Error
	} else {
		return d.DB.Model(&body).
			Where("uuid = ?", uuid).
			Select("InFault", "MessageLevel", "MessageCode", "Message", "LastOk", "InSync").
			Updates(&body).
			Error
	}
}

func (d *GormDatabase) DeleteDevice(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.Device{})
	go d.PublishPointsList("")
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeleteOneDeviceByArgs(args argspkg.Args) (bool, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args).Delete(&deviceModel)
	return d.deleteResponseBuilder(query)
}
