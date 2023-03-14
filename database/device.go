package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/interfaces/connection"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
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

func (d *GormDatabase) CreateDevice(body *model.Device) (*model.Device, error) {
	var net *model.Network
	query := d.DB.Where("uuid = ? ", body.NetworkUUID).First(&net)
	if query.Error != nil {
		return nil, query.Error
	}
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Device)
	body.Name = strings.TrimSpace(body.Name)
	body.ThingClass = model.ThingClass.Device
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, query.Error
}

func (d *GormDatabase) UpdateDevice(uuid string, body *model.Device) (*model.Device, error) {
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
	body.Name = strings.TrimSpace(body.Name)
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

func (d *GormDatabase) UpdateDeviceConnectionErrors(uuid string, device *model.Device) error {
	return d.DB.Model(&model.Device{}).
		Where("uuid = ?", uuid).
		Select("Connection", "ConnectionMessage").
		Updates(&device).
		Error
}

func (d *GormDatabase) DeleteDevice(uuid string) (bool, error) {
	deviceModel, err := d.GetDevice(uuid, api.Args{WithPoints: true})
	if err != nil {
		return false, err
	}
	var wg sync.WaitGroup
	for _, point := range deviceModel.Points {
		wg.Add(1)
		go func(point *model.Point) {
			defer wg.Done()
			_, _ = d.DeletePoint(point.UUID)
		}(point)
	}
	wg.Wait()

	if boolean.IsTrue(deviceModel.AutoMappingEnable) {
		networkModel, err := d.GetNetworkByDeviceUUID(deviceModel.UUID, api.Args{})
		if err != nil {
			return false, err
		}
		streamName := fmt.Sprintf("%s:%s", networkModel.Name, deviceModel.Name)
		stream, _ := d.GetStreamByArgs(api.Args{Name: nils.NewString(streamName)})
		if stream != nil {
			_, _ = d.DeleteStream(stream.UUID)
		}
		fn, err := d.GetOneFlowNetworkByArgs(api.Args{Name: nstring.New(networkModel.AutoMappingFlowNetworkName)})
		if err != nil {
			log.Errorf("failed to find flow network with name %s", networkModel.AutoMappingFlowNetworkName)
			return false, fmt.Errorf("failed to find flow network with name %s", networkModel.AutoMappingFlowNetworkName)
		}
		cli := client.NewFlowClientCliFromFN(fn)
		networkName := getAutoMappedNetworkName(networkModel.Name, boolean.IsFalse(fn.IsRemote))
		url := urls.SingularUrl(urls.DeviceNameUrl, fmt.Sprintf("%s/%s", networkName, deviceModel.Name))
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

func (d *GormDatabase) SyncDevicePoints(uuid string, removeUnlinked bool, args api.Args) ([]*interfaces.SyncModel, error) {
	if removeUnlinked {
		d.removeUnlinkedAutoMappedStreams()
	}
	device, _ := d.GetDevice(uuid, args)
	var outputs []*interfaces.SyncModel
	channel := make(chan *interfaces.SyncModel)
	defer close(channel)
	// This is for syncing child descendants
	for _, point := range device.Points {
		go d.syncPoint(device, point, channel)
	}
	for range device.Points {
		outputs = append(outputs, <-channel)
	}
	return outputs, nil
}

func (d *GormDatabase) syncPoint(device *model.Device, point *model.Point, channel chan *interfaces.SyncModel) {
	var output interfaces.SyncModel
	if boolean.IsTrue(point.CreatedFromAutoMapping) {
		point.Connection = connection.Connected.String()
		point.ConnectionMessage = nstring.New(nstring.NotAvailable)
		_, err := d.GetOneWriterByArgs(api.Args{WriterThingUUID: nils.NewString(point.UUID)})
		if err != nil {
			point.Connection = connection.Broken.String()
			point.ConnectionMessage = nstring.New("writer not found")
			point.Message = "writer not found"
		} else {
			network, _ := d.GetNetworkByDeviceUUID(device.UUID, api.Args{})
			fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{Name: nstring.New(network.AutoMappingFlowNetworkName)})
			if err != nil {
				point.Connection = connection.Broken.String()
				point.ConnectionMessage = nstring.New("flow network clone not found")
			} else {
				networkName := getAutoMappedOriginalNetworkName(network.Name, boolean.IsFalse(fnc.IsRemote))
				cli := client.NewFlowClientCliFromFNC(fnc)
				point_, connectionErr, _ := cli.GetPointByName(networkName, device.Name, point.Name)
				if connectionErr != nil {
					point.Connection = connection.Broken.String()
					point.ConnectionMessage = nstring.New(err.Error())
				} else {
					if point_ == nil || boolean.IsFalse(point_.AutoMappingEnable) {
						_, _ = d.DeletePoint(point.UUID)
						channel <- &output
						return
					}
				}
			}
		}
	}
	err := d.CreatePointAutoMapping(point)
	if err != nil {
		output = interfaces.SyncModel{UUID: point.UUID, IsError: true, Message: nstring.New(err.Error())}
	} else {
		output = interfaces.SyncModel{UUID: point.UUID, IsError: false}
	}
	pointModel := model.Point{}
	if output.IsError {
		pointModel.Connection = connection.Broken.String()
		pointModel.ConnectionMessage = output.Message
	} else {
		pointModel.Connection = connection.Connected.String()
		pointModel.ConnectionMessage = nstring.New(nstring.NotAvailable)
	}
	_ = d.UpdatePointConnectionErrors(output.UUID, &pointModel)
	channel <- &output
}
