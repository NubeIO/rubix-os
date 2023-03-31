package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

func (d *GormDatabase) GetDevices(args api.Args) ([]*model.Device, error) {
	var devicesModel []*model.Device
	query := d.buildDeviceQuery(args)
	if err := query.Find(&devicesModel).Error; err != nil {
		return nil, err
	}
	marshallCacheDevices(devicesModel, args)
	return devicesModel, nil
}

func marshallCacheDevices(devices []*model.Device, args api.Args) {
	for _, device := range devices {
		marshallCachePoints(device.Points, args)
	}
}

func (d *GormDatabase) GetDevice(uuid string, args api.Args) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&deviceModel).Error; err != nil {
		return nil, err
	}
	marshallCachePoints(deviceModel.Points, args)
	return deviceModel, nil
}

func (d *GormDatabase) CreateDeviceTransaction(db *gorm.DB, body *model.Device) (*model.Device, error) {
	var net *model.Network
	query := db.Where("uuid = ? ", body.NetworkUUID).First(&net)
	if query.Error != nil {
		return nil, query.Error
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
	return d.CreateDeviceTransaction(d.DB, body)
}

func (d *GormDatabase) UpdateDeviceTransaction(db *gorm.DB, uuid string, body *model.Device) (*model.Device, error) {
	var deviceModel *model.Device
	query := db.Where("uuid = ?", uuid).First(&deviceModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&deviceModel, body.Tags); err != nil {
			return nil, err
		}
	}
	body.Name = strings.TrimSpace(body.Name)
	body.ThingClass = model.ThingClass.Device
	if err := db.Model(&deviceModel).Select("*").Updates(body).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
}

func (d *GormDatabase) UpdateDevice(uuid string, body *model.Device) (*model.Device, error) {
	return d.UpdateDeviceTransaction(d.DB, uuid, body)
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
	deviceModel, err := d.GetDevice(uuid, api.Args{WithPoints: true})
	if err != nil {
		return false, fmt.Errorf("failed to get device: %w", err)
	}

	if boolean.IsTrue(deviceModel.AutoMappingEnable) {
		var cli *client.FlowClient
		networkModel, _ := d.GetNetworkByDeviceUUID(deviceModel.UUID, api.Args{})

		fn, err := d.GetOneFlowNetworkByArgs(api.Args{Name: &networkModel.AutoMappingFlowNetworkName})
		if err != nil {
			log.Errorf("failed to find flow network with name %s", networkModel.AutoMappingFlowNetworkName)
		} else {
			cli = client.NewFlowClientCliFromFN(fn)
		}

		if cli != nil {
			stream, _ := d.GetOneStreamByArgs(api.Args{AutoMappingDeviceUUID: &deviceModel.UUID})
			if stream != nil { // todo: maybe replace it by d.DeleteStream(<uuid>)
				aType := api.ArgsType
				url := urls.SingularUrlByArg(urls.StreamCloneUrl, aType.SourceUUID, stream.UUID)
				_ = cli.DeleteQuery(url)
				d.DB.Delete(&stream)
			}
		}
	}

	if boolean.IsTrue(deviceModel.CreatedFromAutoMapping) {
		d.DB.
			Where("auto_mapping_device_uuid = ? AND created_from_auto_mapping IS TRUE", deviceModel.UUID).
			Delete(&model.StreamClone{})
	}

	query := d.DB.Delete(&deviceModel)
	go d.PublishPointsList("")
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeleteOneDeviceByArgs(args api.Args) (bool, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args)
	if err := query.First(&deviceModel).Error; err != nil {
		return false, err
	}
	query = query.Delete(&deviceModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) SyncDevicePoints(uuid string, network *model.Network, level interfaces.Level, args api.Args) error {
	if network == nil {
		network, _ = d.GetNetworkByDeviceUUID(uuid, api.Args{})
		if network == nil {
			return errors.New("network doesn't exist")
		}
		device, _ := d.GetDevice(uuid, args)
		network.Devices = append(network.Devices, device)
	}
	if network == nil {
		return errors.New("network doesn't exist")
	}
	return d.CreateNetworkAutoMappings(network, level)
}
