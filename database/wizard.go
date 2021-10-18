package database

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
)

func (d *GormDatabase) WizardRemotePointMapping() (bool, error) {
	var networkModel model.Network
	var deviceModel model.Device
	var pointModel model.Point
	var schModel model.Schedule

	var flowNetworkModel model.FlowNetwork
	var streamModel model.Stream
	var producerModel model.Producer
	var consumerModel model.Consumer
	var writerModel model.Writer

	p, err := d.GetPluginByPath("system")
	if err != nil {
		return false, errors.New("not valid plugin found")
	}
	networkModel.PluginConfId = p.UUID
	networkModel.TransportType = "ip"
	n, err := d.CreateNetwork(&networkModel)
	if err != nil {
		return false, fmt.Errorf("network creation failure: %s", err)
	}
	log.Info("Created a Network")

	deviceModel.NetworkUUID = n.UUID
	deviceModel.TransportType = "ip"
	dev, err := d.CreateDevice(&deviceModel)
	if err != nil {
		return false, fmt.Errorf("device creation failure: %s", err)
	}
	log.Info("Created a Device: ", dev)

	pointModel.DeviceUUID = dev.UUID
	pointModel.Name = "ZATSP-PRO"
	pointModel.IsProducer = utils.NewTrue()
	pointModel.ObjectType = "analogValue" // TODO: check
	pointProducer, err := d.CreatePoint(&pointModel, "")
	if err != nil {
		return false, fmt.Errorf("producer point creation failure: %s", err)
	}
	log.Info("Created a Point for Producer", pointProducer)

	pointModel.DeviceUUID = dev.UUID
	pointModel.Name = "ZATSP"
	pointModel.IsProducer = utils.NewFalse()
	pointModel.ObjectType = "analogValue" // TODO: check
	pointConsumer, err := d.CreatePoint(&pointModel, "")
	if err != nil {
		return false, fmt.Errorf("consumer point creation failure: %s", err)
	}
	log.Info("Created a Point for Consumer", pointConsumer)

	flowNetworkModel.FlowIP = utils.NewStringAddress("0.0.0.0")
	flowNetworkModel.Name = "network"
	fn, err := d.CreateFlowNetwork(&flowNetworkModel)
	if err != nil {
		return false, fmt.Errorf("FlowNetwork creation failure: %s", err)
	}
	log.Info("FlowNetwork is created successfully: ", fn)

	streamModel.FlowNetworks = []*model.FlowNetwork{fn}
	streamModel.Name = "stream"
	stream, err := d.CreateStream(&streamModel)
	if err != nil {
		return false, fmt.Errorf("stream creation failure: %s", err)
	}
	log.Info("Stream is created successfully: ", stream)

	producerModel.Name = "ZATSP"
	producerModel.StreamUUID = stream.UUID
	producerModel.ProducerThingUUID = pointProducer.UUID
	producerModel.ProducerThingClass = "point"
	producerModel.ProducerApplication = "mapping"
	producer, err := d.CreateProducer(&producerModel)
	if err != nil {
		return false, fmt.Errorf("producer creation failure: %s", err)
	}
	log.Info("Producer is created successfully: ", producer)

	streamUUID := stream.UUID
	streamClones, err := d.GetStreamClones(api.Args{SourceUUID: &streamUUID})
	if err != nil {
		return false, fmt.Errorf("StreamClone search failure: %s", err)
	}
	consumerModel.Name = "ZATSP_2"
	consumerModel.ProducerUUID = producer.UUID
	consumerModel.ConsumerApplication = "mapping"
	consumerModel.StreamCloneUUID = streamClones[0].UUID
	consumer, err := d.CreateConsumer(&consumerModel)
	if err != nil {
		return false, fmt.Errorf("consumer creation failure: %s", err)
	}
	log.Info("Consumer is created successfully: ", consumer)

	writerModel.ConsumerUUID = consumerModel.UUID
	writerModel.WriterThingClass = "point"
	writerModel.WriterThingType = "temp"
	writerModel.WriterThingUUID = pointConsumer.UUID
	writer, err := d.CreateWriter(&writerModel)
	if err != nil {
		return false, fmt.Errorf("writer creation failure: %s", err)
	}
	log.Info("Writer is created successfully: ", writer)

	sch, err := d.CreateSchedule(&schModel)
	if err != nil {
		return false, errors.New("CreateSchedule creation failure")
	}
	producerModel = model.Producer{}
	producerModel.Name = "SCH"
	producerModel.StreamUUID = stream.UUID
	producerModel.ProducerThingUUID = sch.UUID
	producerModel.ProducerThingClass = "schedule"
	producerModel.ProducerApplication = "schedule"
	producer2, err := d.CreateProducer(&producerModel)
	if err != nil {
		return false, errors.New("producer-sch creation failure")
	}
	log.Info("Producer-sch is created successfully: ", producer2)

	consumerModel.Name = "schedule"
	consumerModel.ProducerUUID = producer.UUID
	consumerModel.ConsumerApplication = "schedule"
	consumerModel.StreamCloneUUID = streamClones[0].UUID
	consumer2, err := d.CreateConsumer(&consumerModel)
	if err != nil {
		log.Error("CreateConsumer-sch creation failure: ", err)
		return false, errors.New("consumer creation failure")
	}
	log.Info("Consumer-sch is created successfully: ", consumer2)

	writerModel.ConsumerUUID = consumerModel.UUID
	writerModel.WriterThingClass = "schedule"
	writerModel.WriterThingType = "schedule"
	writerModel.WriterThingUUID = sch.UUID
	writer2, err := d.CreateWriter(&writerModel)
	if err != nil {
		log.Error("CreateWriter-sch creation failure: ", err)
		return false, errors.New("writer creation failure")
	}
	log.Info("Writer-sch is created successfully: ", writer2)

	return true, nil
}

func (d *GormDatabase) WizardRemoteSchedule() (bool, error) {
	var schModel model.Schedule

	var flowNetworkModel model.FlowNetwork
	var streamModel model.Stream
	var producerModel model.Producer
	var consumerModel model.Consumer
	var writerModel model.Writer

	_, err := d.GetPluginByPath("system")
	if err != nil {
		return false, errors.New("not valid plugin found")
	}

	schModel.Name = "sch"
	sch, err := d.CreateSchedule(&schModel)
	if err != nil {
		return false, errors.New("CreateSchedule creation failure")
	}

	flowNetworkModel.FlowIP = utils.NewStringAddress("0.0.0.0")
	flowNetworkModel.Name = "FlowNetwork1"
	fn, err := d.CreateFlowNetwork(&flowNetworkModel)
	if err != nil {
		log.Error("FlowNetwork creation failure: ", err)
		return false, errors.New("FlowNetwork creation failure")
	}
	log.Info("FlowNetwork is created successfully: ", fn)

	streamModel.FlowNetworks = []*model.FlowNetwork{fn}
	streamModel.Name = "Stream1"
	stream, err := d.CreateStream(&streamModel)
	if err != nil {
		log.Error("stream creation failure: ", err)
		return false, errors.New("stream creation failure")
	}
	log.Info("Stream is created successfully: ", stream)

	producerModel.Name = "Producer1"
	producerModel.StreamUUID = stream.UUID
	producerModel.ProducerThingUUID = sch.UUID
	producerModel.ProducerThingClass = "schedule"
	producerModel.ProducerApplication = "schedule"
	producer, err := d.CreateProducer(&producerModel)
	if err != nil {
		log.Error("producer creation failure: ", err)
		return false, errors.New("producer creation failure")
	}
	log.Info("Producer is created successfully: ", producer)

	streamUUID := stream.UUID
	streamClones, err := d.GetStreamClones(api.Args{SourceUUID: &streamUUID})
	if err != nil {
		log.Error("StreamClone creation failure: ", err)
		return false, errors.New("StreamClone search failure")
	}
	consumerModel.Name = "Consumer1"
	consumerModel.ProducerUUID = producer.UUID
	consumerModel.ConsumerApplication = "schedule"
	consumerModel.StreamCloneUUID = streamClones[0].UUID
	consumer, err := d.CreateConsumer(&consumerModel)
	if err != nil {
		log.Error("CreateConsumer creation failure: ", err)
		return false, errors.New("consumer creation failure")
	}
	log.Info("Consumer is created successfully: ", consumer)

	writerModel.ConsumerUUID = consumerModel.UUID
	writerModel.WriterThingClass = "schedule"
	writerModel.WriterThingType = "schedule"
	writerModel.WriterThingUUID = sch.UUID
	writer, err := d.CreateWriter(&writerModel)
	if err != nil {
		log.Error("CreateWriter creation failure: ", err)
		return false, errors.New("writer creation failure")
	}
	log.Info("Writer is created successfully: ", writer)

	return true, nil
}
