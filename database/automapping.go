package database

import (
	"fmt"
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
	log "github.com/sirupsen/logrus"
	"reflect"
)

func (d *GormDatabase) CreatePointAutoMapping(point *model.Point) {
	device, err := d.GetDevice(point.DeviceUUID, api.Args{WithTags: true, WithMetaTags: true})
	if err != nil {
		return
	}
	if boolean.IsTrue(device.AutoMappingEnable) {
		err = d.createUpdatePointAutoMapping(device, point)
		if err != nil {
			log.Errorln("points.db.CreatePointAutoMapping() failed to make auto mapping")
		} else {
			log.Println("points.db.CreatePointAutoMapping() added point new mapping")
		}
	}
}

func (d *GormDatabase) UpdatePointAutoMapping(point *model.Point) error {
	device, err := d.GetDevice(point.DeviceUUID, api.Args{WithTags: true, WithMetaTags: true})
	if err != nil {
		return err
	}
	if boolean.IsTrue(device.AutoMappingEnable) {
		err := d.createUpdatePointAutoMapping(device, point)
		if err != nil {
			log.Errorln("points.db.UpdatePointAutoMapping() failed to make auto mapping")
			return err
		} else {
			log.Println("points.db.UpdatePointAutoMapping() added point new mapping")
		}
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
		autoMapping.NetworkTags, autoMapping.NetworkMetaTags, autoMapping.DeviceUUID, autoMapping.DeviceName,
		autoMapping.DeviceTags, autoMapping.DeviceMetaTags, autoMapping.FlowNetworkUUID, autoMapping.IsLocal)
	if err != nil {
		return err
	}
	point, err := d.createPointAutoMappingPoint(device.UUID, autoMapping.PointUUID, autoMapping.PointName,
		autoMapping.PointTags, autoMapping.PointMetaTags)
	if err != nil {
		return err
	}
	_, err = d.createPointAutoMappingWriter(consumer.UUID, point.UUID, point.Name)
	if err != nil {
		return err
	}
	return nil
}

func (d *GormDatabase) createUpdatePointAutoMapping(device *model.Device, point *model.Point) error {
	flowNetwork, err := d.selectFlowNetwork(device.AutoMappingFlowNetworkName, device.AutoMappingFlowNetworkUUID)
	if err != nil {
		return err
	}
	network, err := d.GetNetwork(device.NetworkUUID, api.Args{WithTags: true, WithMetaTags: true})
	if err != nil {
		return err
	}
	// edge
	stream, err := d.createPointAutoMappingStream(flowNetwork, network.UUID, network.Name, device.UUID, device.Name)
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
		NetworkTags:     network.Tags,
		NetworkMetaTags: network.MetaTags,
		DeviceUUID:      device.UUID,
		DeviceName:      device.Name,
		DeviceTags:      device.Tags,
		DeviceMetaTags:  device.MetaTags,
		PointUUID:       point.UUID,
		PointName:       point.Name,
		PointTags:       point.Tags,
		PointMetaTags:   point.MetaTags,
		IsLocal:         boolean.IsFalse(flowNetwork.IsRemote) && boolean.IsFalse(flowNetwork.IsMasterSlave),
	}
	_, err = cli.AddAutoMapping(body)
	if err != nil {
		return err
	}
	return nil
}

func (d *GormDatabase) createPointAutoMappingStream(flowNetwork *model.FlowNetwork, networkUUID string,
	networkName string, deviceUUID string, deviceName string) (stream *model.Stream, err error) {
	name := fmt.Sprintf("%s:%s", networkName, deviceName)
	autoMappingUUID := fmt.Sprintf("%s:%s", networkUUID, deviceUUID)
	stream, _ = d.GetStreamByArgs(api.Args{AutoMappingUUID: nils.NewString(autoMappingUUID), WithFlowNetworks: true})
	if stream == nil {
		streamModel := &model.Stream{}
		streamModel.Enable = boolean.NewTrue()
		streamModel.FlowNetworks = []*model.FlowNetwork{flowNetwork}
		streamModel.Name = name
		streamModel.AutoMappingUUID = autoMappingUUID
		return d.CreateStream(streamModel)
	}
	stream.Name = name
	stream.FlowNetworks = []*model.FlowNetwork{flowNetwork}
	return d.UpdateStream(stream.UUID, stream) // note: to create stream clone in case of it does not exist
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
	if producer.Name != pointName {
		producer.Name = pointName
		return d.UpdateProducer(producer.UUID, producer)
	}
	return producer, nil
}

func (d *GormDatabase) createPointAutoMappingDevice(networkUUID string, networkName string, networkTags []*model.Tag,
	networkMetaTags []*model.NetworkMetaTag, deviceUUID string, deviceName string, deviceTags []*model.Tag,
	deviceMetaTags []*model.DeviceMetaTag, flowNetworkUUID string, isLocal bool) (
	*model.Device, error) {
	syncDevice := &interfaces.SyncDevice{
		NetworkUUID:     networkUUID,
		NetworkName:     networkName,
		NetworkTags:     networkTags,
		NetworkMetaTags: networkMetaTags,
		DeviceUUID:      deviceUUID,
		DeviceName:      deviceName,
		DeviceTags:      deviceTags,
		DeviceMetaTags:  deviceMetaTags,
		FlowNetworkUUID: flowNetworkUUID,
		IsLocal:         isLocal}
	return d.SyncDevice(syncDevice)
}

func (d *GormDatabase) createPointAutoMappingPoint(deviceUUID string, pointUUID string, pointName string,
	pointTags []*model.Tag, pointMetaTags []*model.PointMetaTag) (
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
		pointModel.Tags = pointTags
		pointModel.MetaTags = pointMetaTags
		return d.CreatePoint(pointModel, false)
	}
	_, _ = d.CreatePointMetaTags(point.UUID, pointMetaTags)
	if point.Name != pointName || !reflect.DeepEqual(point.Tags, pointTags) {
		point.Name = pointName
		point.Tags = pointTags
		return d.UpdatePoint(point.UUID, point, false)
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
	if consumer.Name != pointName {
		consumer.Name = pointName
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
	writer.WriterThingName = pointName
	return d.UpdateWriter(writer.UUID, writer) // note: to create writer clone in case of it does not exist
}
