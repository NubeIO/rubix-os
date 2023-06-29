package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/utils/integer"
	log "github.com/sirupsen/logrus"
	"time"
)

// GetDeviceByPoint get a device by point object
func (d *GormDatabase) GetDeviceByPoint(point *model.Point) (*model.Device, error) {
	device, err := d.GetDevice(point.DeviceUUID, args.Args{})
	if err != nil {
		return nil, err
	}
	return device, nil
}

// GetDeviceByPointUUID get a device by its pointUUID
func (d *GormDatabase) GetDeviceByPointUUID(pntUUID string) (*model.Device, error) {
	point, err := d.GetPoint(pntUUID, args.Args{})
	if err != nil || point == nil {
		return nil, err
	}

	device, err := d.GetDevice(point.DeviceUUID, args.Args{})
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (d *GormDatabase) GetOneDeviceByArgs(args args.Args) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args)
	if err := query.First(&deviceModel).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
}

func (d *GormDatabase) GetDeviceByName(networkName string, deviceName string, args args.Args) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args)
	if err := query.Joins("JOIN networks ON devices.network_uuid = networks.uuid").
		Where("networks.name = ?", networkName).Where("devices.name = ?", deviceName).
		First(&deviceModel).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
}

func (d *GormDatabase) deviceNameExistsInNetwork(deviceName, networkUUID string) (device *model.Device, existing bool) {
	network, err := d.GetNetwork(networkUUID, args.Args{WithDevices: true})
	if err != nil {
		return nil, false
	}
	for _, dev := range network.Devices {
		if integer.NonNil(dev.NumberOfDevicesPermitted) == 1 { // some networks like the bacnet-server plugin only allow one network and device
			return dev, true
		}
		if dev.Name == deviceName {
			return dev, true
		}
	}
	return nil, false
}

// SetErrorsForAllPointsOnDevice sets the fault/error properties of all points for a specific device
// messageLevel = model.MessageLevel
// messageCode = model.CommonFaultCode
func (d *GormDatabase) SetErrorsForAllPointsOnDevice(deviceUUID string, message string, messageLevel string, messageCode string) error {
	device, err := d.GetDevice(deviceUUID, args.Args{WithPoints: true})
	if device != nil && err != nil {
		return err
	}
	for _, point := range device.Points {
		point.CommonFault.InFault = true
		point.CommonFault.MessageLevel = messageLevel
		point.CommonFault.MessageCode = messageCode
		point.CommonFault.Message = message
		point.CommonFault.LastFail = time.Now().UTC()
		err = d.UpdatePointErrors(point.UUID, point)
		if err != nil {
			log.Infof("setErrorsForAllPointsOnDevice() Error: %s\n", err.Error())
		}
	}
	return nil
}

// ClearErrorsForAllPointsOnDevice clears the fault/error properties of all points for a specific device
func (d *GormDatabase) ClearErrorsForAllPointsOnDevice(deviceUUID string) error {
	device, err := d.GetDevice(deviceUUID, args.Args{WithPoints: true})
	if device != nil && err != nil {
		return err
	}
	for _, point := range device.Points {
		point.CommonFault.InFault = false
		point.CommonFault.MessageLevel = model.MessageLevel.Normal
		point.CommonFault.MessageCode = model.CommonFaultCode.Ok
		point.CommonFault.Message = ""
		point.CommonFault.LastOk = time.Now().UTC()
		err = d.UpdatePointErrors(point.UUID, point)
		if err != nil {
			log.Infof("clearErrorsForAllPointsOnDevice() Error: %s\n", err.Error())
		}
	}
	return nil
}

func (d *GormDatabase) DeleteDeviceByName(networkName string, deviceName string, args args.Args) (bool, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args)
	if err := query.Joins("JOIN networks ON devices.network_uuid = networks.uuid").
		Where("networks.name = ?", networkName).Where("devices.name = ?", deviceName).
		First(&deviceModel).Error; err != nil {
		return false, err
	}
	query = query.Delete(&deviceModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) GetDevicesTagsForPostgresSync() ([]*interfaces.DeviceTagForPostgresSync, error) {
	var deviceTagsForPostgresModel []*interfaces.DeviceTagForPostgresSync
	query := d.DB.Table("devices_tags").
		Select("devices.source_uuid AS device_uuid, devices_tags.tag_tag AS tag").
		Joins("INNER JOIN devices ON devices.uuid = devices_tags.device_uuid").
		Where("IFNULL(devices.source_uuid,'') != ''").
		Scan(&deviceTagsForPostgresModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return deviceTagsForPostgresModel, nil
}
