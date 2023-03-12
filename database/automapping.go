package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"reflect"
)

func (d *GormDatabase) CreatePointAutoMapping(point *model.Point) error {
	if boolean.IsTrue(point.AutoMappingEnable) {
		device, err := d.GetDevice(point.DeviceUUID, api.Args{WithTags: true, WithMetaTags: true})
		if err != nil {
			return err
		}
		if boolean.IsTrue(device.AutoMappingEnable) {
			err := d.createUpdatePointAutoMapping(device, point)
			if err != nil {
				log.Errorln("points.db.CreatePointAutoMapping() failed to make auto mapping", err)
				return err
			}
			log.Println("points.db.CreatePointAutoMapping() added point new mapping")
			return nil
		}
	}
	return nil
}

func (d *GormDatabase) CreateAutoMapping(autoMapping *interfaces.AutoMapping) error {
	networkName := autoMapping.NetworkName
	if autoMapping.IsLocal {
		networkName = generateLocalNetworkName(networkName)
	}
	network, err := d.GetNetworkByName(networkName, api.Args{})
	if network != nil && network.GlobalUUID != autoMapping.NetworkGlobalUUID {
		return fmt.Errorf("network.name %s already exists", network.Name)
	}
	consumer, err := d.createPointAutoMappingConsumer(autoMapping.StreamUUID, autoMapping.ProducerUUID,
		autoMapping.PointName)
	if err != nil {
		return err
	}
	network, device, err := d.createPointAutoMappingDevice(autoMapping.NetworkGlobalUUID, autoMapping.NetworkName,
		autoMapping.NetworkTags, autoMapping.NetworkMetaTags, autoMapping.DeviceName, autoMapping.DeviceTags,
		autoMapping.DeviceMetaTags, autoMapping.FlowNetworkUUID, autoMapping.IsLocal)
	if err != nil {
		return err
	}
	point, err := d.createPointAutoMappingPoint(network.Name, device.UUID, device.Name, autoMapping.PointName,
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
	fn, err := d.GetOneFlowNetworkByArgs(api.Args{Name: nstring.New(device.AutoMappingFlowNetworkName)})
	if err != nil {
		log.Errorf("failed to find flow network with name %s", device.AutoMappingFlowNetworkName)
		return fmt.Errorf("failed to find flow network with name %s", device.AutoMappingFlowNetworkName)
	}
	network, err := d.GetNetwork(device.NetworkUUID, api.Args{WithTags: true, WithMetaTags: true})
	if err != nil {
		return err
	}
	// edge
	stream, err := d.createPointAutoMappingStream(fn, network.Name, device.Name)
	if err != nil {
		return err
	}
	producer, err := d.createPointAutoMappingProducer(stream.UUID, point.UUID, point.Name)
	if err != nil {
		return err
	}
	// cloud
	deviceInfo, _ := deviceinfo.GetDeviceInfo()
	cli := client.NewFlowClientCliFromFN(fn)
	body := &interfaces.AutoMapping{
		FlowNetworkUUID:   fn.UUID,
		StreamUUID:        stream.UUID,
		ProducerUUID:      producer.UUID,
		NetworkGlobalUUID: deviceInfo.GlobalUUID,
		NetworkName:       network.Name,
		NetworkTags:       network.Tags,
		NetworkMetaTags:   network.MetaTags,
		DeviceName:        device.Name,
		DeviceTags:        device.Tags,
		DeviceMetaTags:    device.MetaTags,
		PointName:         point.Name,
		PointTags:         point.Tags,
		PointMetaTags:     point.MetaTags,
		IsLocal:           boolean.IsFalse(fn.IsRemote) && boolean.IsFalse(fn.IsMasterSlave),
	}
	_, err = cli.AddAutoMapping(body)
	if err != nil {
		return err
	}
	return nil
}

func (d *GormDatabase) createPointAutoMappingStream(flowNetwork *model.FlowNetwork, networkName string,
	deviceName string) (stream *model.Stream, err error) {
	name := fmt.Sprintf("%s:%s", networkName, deviceName)
	stream, _ = d.GetStreamByArgs(api.Args{Name: nils.NewString(name), WithFlowNetworks: true})
	if stream == nil {
		streamModel := &model.Stream{}
		streamModel.Enable = boolean.NewTrue()
		streamModel.FlowNetworks = []*model.FlowNetwork{flowNetwork}
		streamModel.Name = name
		return d.CreateStream(streamModel)
	}
	stream.Name = name
	stream.FlowNetworks = []*model.FlowNetwork{flowNetwork}
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

func (d *GormDatabase) createPointAutoMappingDevice(networkGlobalUUID string, networkName string, networkTags []*model.Tag,
	networkMetaTags []*model.NetworkMetaTag, deviceName string, deviceTags []*model.Tag,
	deviceMetaTags []*model.DeviceMetaTag, flowNetworkUUID string, isLocal bool) (*model.Network, *model.Device, error) {
	syncDevice := &interfaces.SyncDevice{
		NetworkGlobalUUID: networkGlobalUUID,
		NetworkName:       networkName,
		NetworkTags:       networkTags,
		NetworkMetaTags:   networkMetaTags,
		DeviceName:        deviceName,
		DeviceTags:        deviceTags,
		DeviceMetaTags:    deviceMetaTags,
		FlowNetworkUUID:   flowNetworkUUID,
		IsLocal:           isLocal}
	return d.SyncDevice(syncDevice)
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
		return d.CreatePoint(pointModel)
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
		WriterThingName: nils.NewString(pointName)})
	if err != nil {
		writerModel := &model.Writer{}
		writerModel.ConsumerUUID = consumerUUID
		writerModel.WriterThingClass = "point"
		writerModel.WriterThingUUID = pointUUID
		writerModel.WriterThingName = pointName
		return d.CreateWriter(writerModel)
	}
	writer.WriterThingName = pointName
	writer.WriterThingUUID = pointUUID
	return d.UpdateWriter(writer.UUID, writer) // note: to create writer clone in case of it does not exist
}
