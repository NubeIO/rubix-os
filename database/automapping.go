package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/interfaces/connection"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"reflect"
)

func (d *GormDatabase) CreateNetworkAutoMappings(network *model.Network) error {
	if boolean.IsTrue(network.AutoMappingEnable) {
		fn, err := d.GetOneFlowNetworkByArgs(api.Args{Name: nstring.New(network.AutoMappingFlowNetworkName)})
		if err != nil {
			errMessage := fmt.Sprintf("failed to find flow network with name %s", network.AutoMappingFlowNetworkName)
			network.Connection = connection.Broken.String()
			network.ConnectionMessage = nstring.New(errMessage)
			_ = d.UpdateNetworkConnectionErrors(network.UUID, network)
			return errors.New(errMessage)
		} else {
			network.Connection = connection.Connected.String()
			network.ConnectionMessage = nstring.New(nstring.NotAvailable)
			_ = d.UpdateNetworkConnectionErrors(network.UUID, network)
		}
		var amDevices []*interfaces.AutoMappingDevice
		deviceUUIDToStreamUUIDMap, err := d.createPointAutoMappingStreams(fn, network.Name, network.Devices)
		if err != nil {
			log.Error(err)
			return err
		}
		for _, device := range network.Devices {
			if boolean.IsTrue(device.AutoMappingEnable) {
				streamUUID, _ := deviceUUIDToStreamUUIDMap[device.UUID]
				pointUUIDToProducerUUIDMap, err := d.createPointsAutoMappingProducers(streamUUID, device.Points)
				if err != nil {
					log.Error(err)
					return err
				}

				var amPoints []*interfaces.AutoMappingPoint
				for _, point := range device.Points {
					amPoints = append(amPoints, &interfaces.AutoMappingPoint{
						UUID:         point.UUID,
						Name:         point.Name,
						Tags:         point.Tags,
						MetaTags:     point.MetaTags,
						ProducerUUID: pointUUIDToProducerUUIDMap[point.UUID],
					})
				}
				amDevices = append(amDevices, &interfaces.AutoMappingDevice{
					UUID:       device.UUID,
					Name:       device.Name,
					Tags:       device.Tags,
					MetaTags:   device.MetaTags,
					Points:     amPoints,
					StreamUUID: streamUUID,
				})
			}
		}
		cli := client.NewFlowClientCliFromFN(fn)
		deviceInfo, _ := deviceinfo.GetDeviceInfo()
		amNetwork := &interfaces.AutoMappingNetwork{
			GlobalUUID:      deviceInfo.GlobalUUID,
			UUID:            network.UUID,
			Name:            network.Name,
			Tags:            network.Tags,
			MetaTags:        network.MetaTags,
			Devices:         amDevices,
			FlowNetworkUUID: fn.UUID,
		}
		amNetworkError := cli.AddAutoMappings(amNetwork)
		d.updateConnectionErrors(amNetworkError)
		return nil
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
	_ = d.UpdateNetworkConnectionErrors(amNetworkError.UUID, &networkModel)

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
	amNetworkError := &interfaces.AutoMappingNetworkError{UUID: amNetwork.UUID, Error: nil}
	fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{SourceUUID: nstring.New(amNetwork.FlowNetworkUUID)})
	if err != nil {
		amNetworkError.Error = nstring.New(err.Error())
		return amNetworkError
	}
	network, err := d.GetNetworkByName(getAutoMappedNetworkName(fnc.Name, amNetwork.Name), api.Args{})

	if network != nil {
		if network.GlobalUUID != amNetwork.GlobalUUID {
			amNetworkError.Error = nstring.New(fmt.Sprintf("network.name %s already exists in fnc side with different global_uuid", network.Name))
			return amNetworkError
		} else {
			if network.GlobalUUID != amNetwork.GlobalUUID {

			}
			if boolean.IsTrue(network.CreatedFromAutoMapping) {

			}
		}
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
func (d *GormDatabase) createPointAutoMappingStreams(flowNetwork *model.FlowNetwork, networkName string, devices []*model.Device) (map[string]string, error) {
	deviceUUIDToStreamUUIDMap := map[string]string{}
	tx := d.DB.Begin()
	for _, device := range devices {
		if boolean.IsTrue(device.AutoMappingEnable) {
			streamName := getAutoMappedStreamName(flowNetwork.Name, networkName, device.Name)
			stream, _ := d.GetStreamByArgs(api.Args{Name: nils.NewString(streamName)})
			if stream != nil {
				if boolean.IsFalse(stream.CreatedFromAutoMapping) {
					tx.Commit()
					errMessage := fmt.Sprintf("manually created stream_name %s already exists", streamName)
					device.Connection = connection.Broken.String()
					device.ConnectionMessage = nstring.New(errMessage)
					_ = d.UpdateDeviceConnectionErrorsTransaction(tx, device.UUID, device)
					return nil, errors.New(fmt.Sprintf("manually created stream_name %s already exists", streamName))
				} else {
					if stream.AutoMappingUUID != nil && *stream.AutoMappingUUID != device.UUID {
						log.Warnf("mismatch between stream auto mapping uuid %s and device uuid %s", *stream.AutoMappingUUID, device.UUID)
						if err := tx.Delete(&stream).Error; err != nil {
							tx.Rollback()
							device.Connection = connection.Broken.String()
							device.ConnectionMessage = nstring.New(err.Error())
							_ = d.UpdateDeviceConnectionErrorsTransaction(tx, device.UUID, device)
						}
					}
				}
			}
			stream, _ = d.GetStreamByArgs(api.Args{AutoMappingUUID: nstring.New(device.UUID), WithFlowNetworks: true})
			if stream == nil {
				streamModel := model.Stream{}
				streamModel.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Stream)
				d.setStreamModel(flowNetwork, device, streamName, &streamModel)
				if err := tx.Create(&streamModel).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
				deviceUUIDToStreamUUIDMap[device.UUID] = streamModel.UUID
			} else {
				device.Connection = connection.Connected.String()
				device.ConnectionMessage = nstring.New(nstring.NotAvailable)
				_ = d.UpdateDeviceConnectionErrorsTransaction(tx, device.UUID, device)
				if err := tx.Model(&stream).Association("FlowNetworks").Replace([]*model.FlowNetwork{flowNetwork}); err != nil {
					tx.Rollback()
					return nil, err
				}
				d.setStreamModel(flowNetwork, device, streamName, stream)
				stream.Name = getAutoMapperName(streamName)
				stream.CreatedFromAutoMapping = boolean.NewTrue()
				stream.AutoMappingUUID = nstring.New(device.UUID)
				if err := tx.Model(&stream).Updates(&stream).Error; err != nil {
					tx.Rollback()
					device.Connection = connection.Broken.String()
					device.ConnectionMessage = nstring.New(err.Error())
					_ = d.UpdateDeviceConnectionErrorsTransaction(tx, device.UUID, device)
					return nil, err
				}
				deviceUUIDToStreamUUIDMap[device.UUID] = stream.UUID
			}
		}
	}
	for _, device := range devices {
		if boolean.IsTrue(device.AutoMappingEnable) {
			streamName := getAutoMappedStreamName(flowNetwork.Name, networkName, device.Name)
			tempStreamName := getAutoMapperName(streamName)
			q := tx.Model(&model.Stream{}).Where("name = ?", tempStreamName).Update("name", streamName)
			if q.Error != nil {
				tx.Rollback()
				return nil, q.Error
			}
		}
	}
	tx.Commit()
	return deviceUUIDToStreamUUIDMap, nil
}

func (d *GormDatabase) setStreamModel(flowNetwork *model.FlowNetwork, device *model.Device, streamName string, streamModel *model.Stream) {
	streamModel.FlowNetworks = []*model.FlowNetwork{flowNetwork}
	streamModel.Name = getAutoMapperName(streamName)
	streamModel.CreatedFromAutoMapping = boolean.NewTrue()
	streamModel.AutoMappingUUID = nstring.New(device.UUID)
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

func (d *GormDatabase) createPointsAutoMappingProducers(streamUUID string, points []*model.Point) (map[string]string, error) {
	pointUUIDToProducerUUIDMap := map[string]string{}
	tx := d.DB.Begin()
	for _, point := range points {
		if boolean.IsTrue(point.AutoMappingEnable) {
			producer, _ := d.GetOneProducerByArgs(api.Args{StreamUUID: nils.NewString(streamUUID), ProducerThingUUID: nils.NewString(point.UUID)})
			if producer == nil {
				producerModel := model.Producer{}
				producerModel.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Producer)
				d.setProducerModel(streamUUID, point, &producerModel)
				if err := tx.Create(&producerModel).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
				pointUUIDToProducerUUIDMap[point.UUID] = producerModel.UUID
			} else {
				d.setProducerModel(streamUUID, point, producer)
				if err := tx.Save(producer).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
				pointUUIDToProducerUUIDMap[point.UUID] = producer.UUID
			}
		}
	}
	for _, point := range points {
		if boolean.IsTrue(point.AutoMappingEnable) {
			q := tx.Model(&model.Producer{}).Where("producer_thing_uuid = ?", point.UUID).Update("name", point.Name)
			if q.Error != nil {
				tx.Rollback()
				return nil, q.Error
			}
		}
	}
	tx.Commit()
	return pointUUIDToProducerUUIDMap, nil
}

func (d *GormDatabase) setProducerModel(streamUUID string, point *model.Point, producerModel *model.Producer) {
	producerModel.Enable = boolean.NewTrue()
	producerModel.Name = getAutoMapperName(point.Name)
	producerModel.StreamUUID = streamUUID
	producerModel.ProducerThingUUID = point.UUID
	producerModel.ProducerThingName = point.Name
	producerModel.ProducerThingClass = "point"
	producerModel.ProducerApplication = "mapping"
	producerModel.EnableHistory = point.HistoryEnable
	producerModel.HistoryType = point.HistoryType
	producerModel.HistoryInterval = point.HistoryInterval
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
