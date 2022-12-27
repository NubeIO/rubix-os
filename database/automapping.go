package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/interfaces/connection"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) CreatePointAutoMapping(point *model.Point) error {
	return d.createUpdatePointAutoMapping(point)
}

func (d *GormDatabase) UpdatePointAutoMapping(point *model.Point) error {
	if err := d.createUpdatePointAutoMapping(point); err != nil {
		return err
	}
	if point.AutoMappingUUID != "" {
		point.Connection = connection.Connected.String()
		point.Message = nstring.NotAvailable
		writer, err := d.GetOneWriterByArgs(api.Args{WriterThingUUID: nils.NewString(point.UUID)})
		if err != nil {
			point.Connection = connection.Broken.String()
			point.Message = "writer not found"
		} else {
			consumer, _ := d.GetConsumer(writer.ConsumerUUID, api.Args{})
			streamClone, _ := d.GetStreamClone(consumer.StreamCloneUUID, api.Args{})
			flowNetworkClone, _ := d.GetFlowNetworkClone(streamClone.FlowNetworkCloneUUID, api.Args{})
			cli := client.NewFlowClientCliFromFNC(flowNetworkClone)
			_, err = cli.GetQueryMarshal(urls.SingularUrl(urls.PointUrl, point.AutoMappingUUID), model.Point{})
			if err != nil {
				point.Connection = connection.Broken.String()
				point.Message = err.Error()
			}
		}
		err = d.UpdatePointErrors(point.UUID, point)
	}
	return nil
}

func (d *GormDatabase) CreateAutoMapping(autoMapping *interfaces.AutoMapping) error {
	consumer, err := d.createPointAutoMappingConsumer(autoMapping.StreamUUID, autoMapping.ProducerUUID,
		autoMapping.PointName)
	if err != nil {
		return err
	}
	device, err := d.createPointAutoMappingDevice(autoMapping.NetworkUUID, autoMapping.NetworkName,
		autoMapping.DeviceUUID, autoMapping.DeviceName, autoMapping.FlowNetworkUUID, autoMapping.IsLocal)
	if err != nil {
		return err
	}
	point, err := d.createPointAutoMappingPoint(device.UUID, autoMapping.PointUUID, autoMapping.PointName)
	if err != nil {
		return err
	}
	_, err = d.createPointAutoMappingWriter(consumer.UUID, point.UUID, point.Name)
	if err != nil {
		return err
	}
	return nil
}

func (d *GormDatabase) createUpdatePointAutoMapping(point *model.Point) error {
	device, err := d.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil {
		return err
	}
	network, err := d.GetNetwork(device.NetworkUUID, api.Args{})
	if err != nil {
		return err
	}
	if boolean.IsTrue(network.AutoMappingEnable) {
		flowNetwork, err := d.selectFlowNetwork(network.AutoMappingFlowNetworkName, network.AutoMappingFlowNetworkUUID)
		if err != nil {
			return err
		}
		// edge
		stream, err := d.createPointAutoMappingStream(flowNetwork, network.UUID, network.Name)
		if err != nil {
			return err
		}
		producer, err := d.createPointAutoMappingProducer(stream.UUID, point.UUID, point.Name)
		if err != nil {
			return err
		}
		// cloud
		cli := client.NewFlowClientCliFromFN(flowNetwork)
		body := &interfaces.AutoMapping{
			FlowNetworkUUID: flowNetwork.UUID,
			StreamUUID:      stream.UUID,
			ProducerUUID:    producer.UUID,
			NetworkUUID:     network.UUID,
			NetworkName:     network.Name,
			DeviceUUID:      device.UUID,
			DeviceName:      device.Name,
			PointUUID:       point.UUID,
			PointName:       point.Name,
			IsLocal:         boolean.IsFalse(flowNetwork.IsRemote) && boolean.IsFalse(flowNetwork.IsMasterSlave),
		}
		_, err = cli.AddAutoMapping(body)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *GormDatabase) createPointAutoMappingStream(flowNetwork *model.FlowNetwork, networkUUID string,
	networkName string) (stream *model.Stream, err error) {
	stream, _ = d.GetStreamByArgs(api.Args{AutoMappingUUID: nils.NewString(networkUUID), WithFlowNetworks: true})
	if stream == nil {
		streamModel := &model.Stream{}
		streamModel.Enable = boolean.NewTrue()
		streamModel.FlowNetworks = []*model.FlowNetwork{flowNetwork}
		streamModel.Name = networkName
		streamModel.AutoMappingUUID = networkUUID
		return d.CreateStream(streamModel)
	}
	updateStream := stream.Name != networkName
	if !updateStream {
		for _, fn := range stream.FlowNetworks {
			updateStream = fn.UUID != flowNetwork.UUID
		}
	}
	if updateStream {
		stream.Name = networkName
		stream.FlowNetworks = []*model.FlowNetwork{flowNetwork}
		return d.UpdateStream(stream.UUID, stream)
	}
	return stream, nil
}

func (d *GormDatabase) createPointAutoMappingProducer(streamUUID string, pointUUID string, pointName string) (
	producer *model.Producer, err error) {
	producer, _ = d.GetOneProducerByArgs(api.Args{StreamUUID: nils.NewString(streamUUID), ProducerThingUUID: nils.NewString(pointUUID)})
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
	if producer.ProducerThingName != pointName {
		producer.Name = pointName
		producer.ProducerThingName = pointName
		return d.UpdateProducer(producer.UUID, producer)
	}
	return producer, nil
}

func (d *GormDatabase) createPointAutoMappingDevice(networkUUID string, networkName string, deviceUUID string,
	deviceName string, flowNetworkUUID string, isLocal bool) (*model.Device, error) {
	syncDevice := &model.SyncDevice{
		NetworkUUID:     networkUUID,
		NetworkName:     networkName,
		DeviceUUID:      deviceUUID,
		DeviceName:      deviceName,
		FlowNetworkUUID: flowNetworkUUID,
		IsLocal:         isLocal}
	return d.SyncDevice(syncDevice)
}

func (d *GormDatabase) createPointAutoMappingPoint(deviceUUID string, pointUUID string, pointName string) (
	point *model.Point, err error) {
	point, err = d.GetOnePointByArgs(api.Args{AutoMappingUUID: nils.NewString(pointUUID)})
	if point == nil {
		pointModel := &model.Point{}
		pointModel.Enable = boolean.NewTrue()
		pointModel.Name = pointName
		pointModel.DeviceUUID = deviceUUID
		pointModel.ThingClass = "point"
		pointModel.ThingType = ""
		pointModel.AutoMappingUUID = pointUUID
		return d.CreatePoint(pointModel, false)
	}
	if point.Name != pointName {
		point.Name = pointName
		return d.UpdatePoint(point.UUID, point, false, false)
	}
	return point, err
}

func (d *GormDatabase) createPointAutoMappingConsumer(streamUUID string, producerUUID string, pointName string) (
	consumer *model.Consumer, err error) {
	streamClone, err := d.GetStreamCloneByArg(api.Args{SourceUUID: nils.NewString(streamUUID)})
	if err != nil {
		return nil, err
	}
	consumer, _ = d.GetOneConsumerByArgs(api.Args{ProducerUUID: nils.NewString(producerUUID)})
	if consumer == nil {
		consumerModel := &model.Consumer{}
		consumerModel.StreamCloneUUID = streamClone.UUID
		consumerModel.Enable = boolean.NewTrue()
		consumerModel.Name = pointName
		consumerModel.ProducerUUID = producerUUID
		consumerModel.ProducerThingName = pointName
		consumerModel.ConsumerApplication = "mapping"
		return d.CreateConsumer(consumerModel)
	}
	if consumer.ProducerThingName != pointName {
		consumer.Name = pointName
		consumer.ProducerThingName = pointName
		return d.UpdateConsumer(consumer.UUID, consumer)
	}
	return consumer, nil
}

func (d *GormDatabase) createPointAutoMappingWriter(consumerUUID string, pointUUID string, pointName string) (
	writer *model.Writer, err error) {
	writer, err = d.GetOneWriterByArgs(api.Args{ConsumerUUID: nils.NewString(consumerUUID),
		WriterThingUUID: nils.NewString(pointUUID)})
	if err != nil {
		writerModel := &model.Writer{}
		writerModel.ConsumerUUID = consumerUUID
		writerModel.WriterThingClass = "point"
		writerModel.WriterThingUUID = pointUUID
		writerModel.WriterThingName = pointName
		return d.CreateWriter(writerModel)
	}
	if writer.WriterThingName != pointName {
		writer.WriterThingName = pointName
		return d.UpdateWriter(writer.UUID, writer)
	}
	return writer, nil
}
