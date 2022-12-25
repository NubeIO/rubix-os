package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/interfaces/connection"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/nstring"
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
	_ = d.syncAfterCreateUpdateDevice(body.UUID, api.Args{})
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
	_ = d.syncAfterCreateUpdateDevice(body.UUID, api.Args{WithPoints: true})
	return deviceModel, nil
}

// UpdateDeviceErrors will only update the CommonFault properties of the device, all other properties won't be updated
// Does not update `LastOk`
func (d *GormDatabase) UpdateDeviceErrors(uuid string, body *model.Device) error {
	return d.DB.Model(&body).
		Where("uuid = ?", uuid).
		Select("InFault", "MessageLevel", "MessageCode", "Message", "LastFail", "InSync", "Connection").
		Updates(&body).
		Error
}

func (d *GormDatabase) DeleteDevice(uuid string) (bool, error) {
	var aType = api.ArgsType
	deviceModel, err := d.GetDevice(uuid, api.Args{WithPoints: true})
	if err != nil {
		return false, err
	}
	for _, point := range deviceModel.Points {
		_, _ = d.DeletePoint(point.UUID)
	}
	network, err := d.GetNetworkByDeviceUUID(deviceModel.UUID, api.Args{})
	if err != nil {
		return false, err
	}
	if boolean.IsTrue(network.AutoMappingEnable) {
		fn, _ := d.selectFlowNetwork(network.AutoMappingFlowNetworkName, network.AutoMappingFlowNetworkUUID)
		cli := client.NewFlowClientCliFromFN(fn)
		url := urls.SingularUrlByArg(urls.DeviceUrl, aType.AutoMappingUUID, deviceModel.UUID)
		_ = cli.DeleteQuery(url)
	}
	query := d.DB.Delete(&deviceModel)
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

func (d *GormDatabase) SyncDevicePoints(uuid string, args api.Args) ([]*interfaces.SyncModel, error) {
	device, _ := d.GetDevice(uuid, api.Args{WithPoints: true})
	var outputs []*interfaces.SyncModel
	channel := make(chan *interfaces.SyncModel)
	defer close(channel)
	// This is for syncing child descendants
	for _, point := range device.Points {
		go d.syncPoint(point, channel)
	}
	for range device.Points {
		outputs = append(outputs, <-channel)
	}
	return outputs, nil
}

func (d *GormDatabase) syncAfterCreateUpdateDevice(uuid string, args api.Args) error {
	device, err := d.GetDevice(uuid, args)
	if err != nil {
		return err
	}
	network, err := d.GetNetworkByDeviceUUID(device.UUID, api.Args{})
	if err != nil {
		return err
	}
	if boolean.IsTrue(network.AutoMappingEnable) {
		fn, err := d.selectFlowNetwork(network.AutoMappingFlowNetworkName, network.AutoMappingFlowNetworkUUID)
		if err != nil {
			return err
		}
		cli := client.NewFlowClientCliFromFN(fn)
		syncBody := model.SyncDevice{
			NetworkUUID:     network.UUID,
			NetworkName:     network.Name,
			DeviceUUID:      device.UUID,
			DeviceName:      device.Name,
			FlowNetworkUUID: fn.UUID,
			IsLocal:         boolean.IsFalse(fn.IsRemote) && boolean.IsFalse(fn.IsMasterSlave),
		}
		_, err = cli.SyncDevice(&syncBody)
		if err != nil {
			return err
		}
		if args.WithPoints {
			_, _ = d.SyncDevicePoints(device.UUID, args)
		}
	} else if device.AutoMappingUUID != "" {
		device.Connection = connection.Connected.String()
		device.Message = nstring.NotAvailable
		fnc, err := d.GetFlowNetworkClone(network.AutoMappingFlowNetworkUUID, api.Args{})
		if err != nil {
			device.Connection = connection.Broken.String()
			device.Message = "flow network clone not found"
		} else {
			cli := client.NewFlowClientCliFromFNC(fnc)
			_, err = cli.GetQueryMarshal(urls.SingularUrl(urls.DeviceUrl, device.AutoMappingUUID), model.Device{})
			if err != nil {
				device.Connection = connection.Broken.String()
				device.Message = err.Error()
			}
		}
		_ = d.UpdateDeviceErrors(device.UUID, device)
		if args.WithPoints {
			_, _ = d.SyncDevicePoints(device.UUID, args)
		}
	}
	return err
}

func (d *GormDatabase) syncPoint(point *model.Point, channel chan *interfaces.SyncModel) {
	err := d.UpdatePointAutoMapping(point)
	var output interfaces.SyncModel
	if err != nil {
		output = interfaces.SyncModel{UUID: point.UUID, IsError: true, Message: nstring.New(err.Error())}
	} else {
		output = interfaces.SyncModel{UUID: point.UUID, IsError: false}
	}
	channel <- &output
}
