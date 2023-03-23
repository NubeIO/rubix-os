package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/interfaces/connection"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"reflect"
)

func (d *GormDatabase) CreateNetworkAutoMappings(network *model.Network) *interfaces.AutoMappingNetworkError {
	if boolean.IsTrue(network.AutoMappingEnable) {
		fn, err := d.GetOneFlowNetworkByArgs(api.Args{Name: nstring.New(network.AutoMappingFlowNetworkName)})
		if err != nil {
			errMessage := fmt.Sprintf("failed to find flow network with name %s", network.AutoMappingFlowNetworkName)
			network.Connection = connection.Broken.String()
			network.ConnectionMessage = nstring.New(errMessage)
			_ = d.UpdateNetworkConnectionErrors(network.UUID, network)
			return &interfaces.AutoMappingNetworkError{Name: network.Name, Error: nstring.New(errMessage)}
		} else {
			network.Connection = connection.Connected.String()
			network.ConnectionMessage = nstring.New(nstring.NotAvailable)
			_ = d.UpdateNetworkConnectionErrors(network.UUID, network)
		}
		var amDevices []*interfaces.AutoMappingDevice
		for _, device := range network.Devices {
			if boolean.IsTrue(device.AutoMappingEnable) {
				var amPoints []*interfaces.AutoMappingPoint
				stream, err := d.createPointAutoMappingStream(fn, network.Name, device.Name)
				if err != nil {
					continue
				}
				for _, point := range device.Points {
					if boolean.IsTrue(point.AutoMappingEnable) {
						producer, err := d.createPointAutoMappingProducer(stream.UUID, point.UUID, point.Name)
						if err != nil {
							continue
						}
						amPoints = append(amPoints, &interfaces.AutoMappingPoint{
							Name:         point.Name,
							Tags:         point.Tags,
							MetaTags:     point.MetaTags,
							ProducerUUID: producer.UUID,
						})
					}
				}
				amDevices = append(amDevices, &interfaces.AutoMappingDevice{
					Name:       device.Name,
					Tags:       device.Tags,
					MetaTags:   device.MetaTags,
					Points:     amPoints,
					StreamUUID: stream.UUID,
				})
			}
		}
		cli := client.NewFlowClientCliFromFN(fn)
		deviceInfo, _ := deviceinfo.GetDeviceInfo()
		amNetworkError := cli.AddAutoMappings(&interfaces.AutoMappingNetwork{
			GlobalUUID:      deviceInfo.GlobalUUID,
			Name:            network.Name,
			Tags:            network.Tags,
			MetaTags:        network.MetaTags,
			Devices:         amDevices,
			FlowNetworkUUID: fn.UUID,
		})
		d.updateConnectionErrors(amNetworkError)
		return amNetworkError
	}
	return nil
}

func (d *GormDatabase) updateConnectionErrors(amNetworkError *interfaces.AutoMappingNetworkError) {
	networkModel := model.Network{}
	if amNetworkError.Error != nil {
		networkModel.Connection = connection.Broken.String()
		networkModel.ConnectionMessage = amNetworkError.Error
	} else {
		networkModel.Connection = connection.Connected.String()
		networkModel.ConnectionMessage = nstring.New(nstring.NotAvailable)
	}
	_ = d.UpdateNetworkConnectionErrorsByName(amNetworkError.Name, &networkModel)

	for _, amDeviceError := range amNetworkError.Devices {
		deviceModel := model.Device{}
		if amDeviceError.Error != nil {
			deviceModel.Connection = connection.Broken.String()
			deviceModel.ConnectionMessage = amDeviceError.Error
		} else {
			deviceModel.Connection = connection.Connected.String()
			deviceModel.ConnectionMessage = nstring.New(nstring.NotAvailable)
		}
		_ = d.UpdateDeviceConnectionErrorsByName(amDeviceError.Name, &deviceModel)

		for _, amPointError := range amDeviceError.Points {
			pointModel := model.Point{}
			if amPointError.Error != nil {
				pointModel.Connection = connection.Broken.String()
				pointModel.ConnectionMessage = amPointError.Error
			} else {
				pointModel.Connection = connection.Connected.String()
				pointModel.ConnectionMessage = nstring.New(nstring.NotAvailable)
			}
			_ = d.UpdatePointConnectionErrors(amPointError.Name, &pointModel)
		}
	}
}

func (d *GormDatabase) CreateAutoMapping(amNetwork *interfaces.AutoMappingNetwork) *interfaces.AutoMappingNetworkError {
	amNetworkError := &interfaces.AutoMappingNetworkError{Name: amNetwork.Name, Error: nil}
	fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{SourceUUID: nstring.New(amNetwork.FlowNetworkUUID)})
	if err != nil {
		amNetworkError.Error = nstring.New(err.Error())
		return amNetworkError
	}
	network, err := d.GetNetworkByName(getAutoMappedNetworkName(fnc.Name, amNetwork.Name), api.Args{})
	if network != nil && network.GlobalUUID != amNetwork.GlobalUUID {
		amNetworkError.Error = nstring.New(fmt.Sprintf("network.name %s already exists", network.Name))
		return amNetworkError
	}
	for _, amDevice := range amNetwork.Devices {
		amDeviceError := &interfaces.AutoMappingDeviceError{Name: amDevice.Name}
		network, err := d.createPointAutoMappingNetwork(amNetwork)
		if err != nil {
			amDeviceError.Error = nstring.New(err.Error())
			amNetworkError.Devices = append(amNetworkError.Devices, amDeviceError)
			continue
		}
		device, err := d.createPointAutoMappingDevice(network, amDevice)
		if err != nil {
			amDeviceError.Error = nstring.New(err.Error())
			amNetworkError.Devices = append(amNetworkError.Devices, amDeviceError)
			continue
		}
		amDeviceError.Consumers = d.createPointAutoMappingConsumers(amDevice)
		amDeviceError.Points = d.createPointAutoMappingPoints(network.Name, device.UUID, device.Name, amDevice)
		amDeviceError.Writers = d.createPointAutoMappingWriters(network.Name, device.Name, amDevice)
		amNetworkError.Devices = append(amNetworkError.Devices, amDeviceError)
	}
	return amNetworkError
}

func (d *GormDatabase) createPointAutoMappingStream(flowNetwork *model.FlowNetwork, networkName string,
	deviceName string) (stream *model.Stream, err error) {
	streamName := getAutoMappedStreamName(flowNetwork.Name, networkName, deviceName)
	stream, _ = d.GetStreamByArgs(api.Args{Name: nils.NewString(streamName), WithFlowNetworks: true})
	if stream == nil {
		streamModel := &model.Stream{}
		streamModel.Enable = boolean.NewTrue()
		streamModel.FlowNetworks = []*model.FlowNetwork{flowNetwork}
		streamModel.Name = streamName
		streamModel.CreatedFromAutoMapping = boolean.NewTrue()
		return d.CreateStream(streamModel)
	}
	stream.Name = streamName
	stream.FlowNetworks = []*model.FlowNetwork{flowNetwork}
	stream.CreatedFromAutoMapping = boolean.NewTrue()
	return d.UpdateStream(stream.UUID, stream) // note: to create stream clone in case of it does not exist
}

func (d *GormDatabase) createPointAutoMappingProducer(streamUUID string, pointUUID string, pointName string) (
	producer *model.Producer, err error) {
	producer, _ = d.GetOneProducerByArgs(api.Args{StreamUUID: nils.NewString(streamUUID), Name: nils.NewString(pointName)})
	if producer == nil {
		producerModel := &model.Producer{}
		producerModel.Enable = boolean.NewTrue()
		producerModel.Name = pointName
		producerModel.StreamUUID = streamUUID
		producerModel.ProducerThingUUID = pointUUID
		producerModel.ProducerThingName = pointName
		producerModel.ProducerThingClass = "point"
		producerModel.ProducerApplication = "mapping"
		producerModel.EnableHistory = boolean.NewFalse()
		producerModel.HistoryType = model.HistoryTypeInterval
		producerModel.HistoryInterval = integer.New(15)
		producerModel.EnableHistory = boolean.NewFalse()
		return d.CreateProducer(producerModel)
	}
	if producer.Name != pointName || producer.ProducerThingUUID != pointUUID {
		producer.Name = pointName
		producer.ProducerThingUUID = pointUUID
		return d.UpdateProducer(producer.UUID, producer)
	}
	return producer, nil
}

func (d *GormDatabase) createPointAutoMappingNetwork(amNetwork *interfaces.AutoMappingNetwork) (
	*model.Network, error) {
	syncNetwork := &interfaces.SyncNetwork{
		NetworkGlobalUUID: amNetwork.GlobalUUID,
		NetworkName:       amNetwork.Name,
		NetworkTags:       amNetwork.Tags,
		NetworkMetaTags:   amNetwork.MetaTags,
		FlowNetworkUUID:   amNetwork.FlowNetworkUUID,
	}
	return d.SyncNetwork(syncNetwork)
}

func (d *GormDatabase) createPointAutoMappingDevice(network *model.Network, amDevice *interfaces.AutoMappingDevice) (
	*model.Device, error) {
	syncDevice := &interfaces.SyncDevice{
		DeviceName:     amDevice.Name,
		DeviceTags:     amDevice.Tags,
		DeviceMetaTags: amDevice.MetaTags,
	}
	return d.SyncDevice(syncDevice, network)
}

func (d *GormDatabase) createPointAutoMappingConsumers(amDevice *interfaces.AutoMappingDevice) []*interfaces.AutoMappingConsumerError {
	var apConsumerErrors []*interfaces.AutoMappingConsumerError
	channel := make(chan *interfaces.AutoMappingConsumerError)
	defer close(channel)
	for _, amPoint := range amDevice.Points {
		go d.createPointAutoMappingConsumer(amDevice.StreamUUID, amPoint.ProducerUUID, amPoint.Name, channel)
	}
	for range amDevice.Points {
		apConsumerErrors = append(apConsumerErrors, <-channel)
	}
	return apConsumerErrors
}

func (d *GormDatabase) createPointAutoMappingConsumer(streamUUID string, producerUUID string, pointName string,
	channel chan *interfaces.AutoMappingConsumerError) {
	var amConsumerError interfaces.AutoMappingConsumerError
	streamClone, err := d.GetStreamCloneByArg(api.Args{SourceUUID: nils.NewString(streamUUID)})
	if err != nil {
		amConsumerError.Error = nstring.New(err.Error())
	} else {
		consumer, _ := d.GetOneConsumerByArgs(api.Args{ProducerUUID: nils.NewString(producerUUID)})
		if consumer == nil {
			consumerModel := &model.Consumer{}
			consumerModel.StreamCloneUUID = streamClone.UUID
			consumerModel.Enable = boolean.NewTrue()
			consumerModel.Name = pointName
			consumerModel.ProducerUUID = producerUUID
			consumerModel.ProducerThingName = pointName
			consumerModel.ConsumerApplication = "mapping"
			_, err := d.CreateConsumer(consumerModel)
			if err != nil {
				amConsumerError.Error = nstring.New(err.Error())
			}
		} else if consumer.Name != pointName {
			consumer.Name = pointName
			_, err := d.UpdateConsumer(consumer.UUID, consumer)
			if err != nil {
				amConsumerError.Error = nstring.New(err.Error())
			}
		}
	}
	channel <- &amConsumerError
}

func (d *GormDatabase) createPointAutoMappingPoints(networkName string, deviceUUID string, deviceName string,
	amDevice *interfaces.AutoMappingDevice) []*interfaces.AutoMappingPointError {
	var amPointErrors []*interfaces.AutoMappingPointError
	for _, amPoint := range amDevice.Points {
		apPointError := &interfaces.AutoMappingPointError{Name: amPoint.Name}
		_, err := d.createPointAutoMappingPoint(networkName, deviceUUID, deviceName, amPoint.Name,
			amPoint.Tags, amPoint.MetaTags)
		if err != nil {
			apPointError.Error = nstring.New(err.Error())
			amPointErrors = append(amPointErrors, apPointError)
			continue
		}
		amPointErrors = append(amPointErrors, apPointError)
	}
	d.PublishPointsList("")
	return amPointErrors
}

func (d *GormDatabase) createPointAutoMappingPoint(networkName, deviceUUID string, deviceName string, pointName string,
	pointTags []*model.Tag, pointMetaTags []*model.PointMetaTag) (
	point *model.Point, err error) {
	point, err = d.GetPointByName(networkName, deviceName, pointName, api.Args{})
	if point == nil {
		pointModel := &model.Point{}
		pointModel.Enable = boolean.NewTrue()
		pointModel.Name = pointName
		pointModel.DeviceUUID = deviceUUID
		pointModel.ThingClass = "point"
		pointModel.ThingType = ""
		pointModel.Tags = pointTags
		pointModel.MetaTags = pointMetaTags
		pointModel.CreatedFromAutoMapping = boolean.NewTrue()
		pointModel.EnableWriteable = boolean.NewTrue()
		return d.CreatePoint(pointModel, false)
	}
	_, _ = d.CreatePointMetaTags(point.UUID, pointMetaTags)
	if !reflect.DeepEqual(point.Tags, pointTags) {
		point.Tags = pointTags
		if err := d.updateTags(&point, pointTags); err != nil {
			return nil, err
		}
	}
	return point, err
}

func (d *GormDatabase) createPointAutoMappingWriters(networkName string, deviceName string,
	amDevice *interfaces.AutoMappingDevice) []*interfaces.AutoMappingWriterError {
	var amWriterErrors []*interfaces.AutoMappingWriterError
	channel := make(chan *interfaces.AutoMappingWriterError)
	defer close(channel)
	for _, amPoint := range amDevice.Points {
		amWriterError := &interfaces.AutoMappingWriterError{Name: amPoint.Name}
		consumer, err := d.GetOneConsumerByArgs(api.Args{ProducerUUID: nstring.New(amPoint.ProducerUUID)})
		if err != nil {
			amWriterError.Error = nstring.New(err.Error())
			amWriterErrors = append(amWriterErrors, amWriterError)
			continue
		}
		point, err := d.GetPointByName(networkName, deviceName, amPoint.Name, api.Args{})
		if err != nil {
			amWriterError.Error = nstring.New(err.Error())
			amWriterErrors = append(amWriterErrors, amWriterError)
			continue
		}
		go d.createPointAutoMappingWriter(consumer.UUID, point.UUID, point.Name, channel)
	}
	for range amDevice.Points {
		amWriterErrors = append(amWriterErrors, <-channel)
	}
	return amWriterErrors
}

func (d *GormDatabase) createPointAutoMappingWriter(consumerUUID string, pointUUID string, pointName string,
	channel chan *interfaces.AutoMappingWriterError) {
	var amWriterError interfaces.AutoMappingWriterError
	writer, err := d.GetOneWriterByArgs(api.Args{ConsumerUUID: nils.NewString(consumerUUID),
		WriterThingName: nils.NewString(pointName)})
	if err != nil {
		writerModel := &model.Writer{}
		writerModel.ConsumerUUID = consumerUUID
		writerModel.WriterThingClass = "point"
		writerModel.WriterThingUUID = pointUUID
		writerModel.WriterThingName = pointName
		_, err := d.CreateWriter(writerModel)
		if err != nil {
			amWriterError.Error = nstring.New(err.Error())
		}
	} else {
		writer.WriterThingName = pointName
		writer.WriterThingUUID = pointUUID
		_, err := d.UpdateWriter(writer.UUID, writer) // note: to create writer clone in case of it does not exist
		if err != nil {
			amWriterError.Error = nstring.New(err.Error())
		}
	}
	channel <- &amWriterError
}
