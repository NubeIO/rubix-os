package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/client"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
)

func (d *GormDatabase) WizardRemotePointMapping(body *model.FlowNetworkCredential) (bool, error) {
	flow, err := d.wizardRemotePointMappingOnProducerSide(body)
	if err != nil {
		return false, err
	}

	deviceInfo, err := d.GetDeviceInfo()
	if err != nil {
		return false, fmt.Errorf("GetDeviceInfo: %s", deviceInfo)
	}
	cli := client.NewFlowClientCli(flow.FlowIP, flow.FlowPort, flow.FlowToken, flow.IsMasterSlave, flow.GlobalUUID, model.IsFNCreator(flow))
	_, err = cli.WizardRemotePointMappingOnConsumerSideByProducerSide(deviceInfo.GlobalUUID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *GormDatabase) wizardRemotePointMappingOnProducerSide(body *model.FlowNetworkCredential) (*model.FlowNetwork, error) {
	var flowNetworkModel model.FlowNetwork
	var streamModel model.Stream
	var producerModel model.Producer

	point := model.Point{}
	point.Name = "ZATSP-PRO"
	pointProducer, err := d.WizardNewNetworkDevicePoint("system", nil, nil, &point)
	if err != nil {
		return nil, fmt.Errorf("CreateNetworkDevicePoint: %s", err)
	}

	flowNetworkModel.Name = "FlowNetwork"
	flowNetworkModel.FlowIP = utils.NewStringAddress(body.FlowIP)
	flowNetworkModel.FlowPort = utils.NewInt(body.FlowPort)
	flowNetworkModel.FlowUsername = utils.NewStringAddress(body.FlowUsername)
	flowNetworkModel.FlowPassword = utils.NewStringAddress(body.FlowPassword)
	fn, err := d.CreateFlowNetwork(&flowNetworkModel)
	if err != nil {
		return nil, fmt.Errorf("FlowNetwork creation failure: %s", err)
	}
	log.Info("FlowNetwork is created successfully: ", fn)

	streamModel.FlowNetworks = []*model.FlowNetwork{fn}
	streamModel.Name = "stream"
	stream, err := d.CreateStream(&streamModel)
	if err != nil {
		return nil, fmt.Errorf("stream creation failure: %s", err)
	}
	log.Info("Stream is created successfully: ", stream)

	producerModel.Name = "ZATSP"
	producerModel.StreamUUID = stream.UUID
	producerModel.ProducerThingUUID = pointProducer.UUID
	producerModel.ProducerThingClass = "point"
	producerModel.ProducerApplication = "mapping"
	producer, err := d.CreateProducer(&producerModel)
	if err != nil {
		return nil, fmt.Errorf("producer creation failure: %s", err)
	}
	log.Info("Producer is created successfully: ", producer)

	return fn, nil
}

func (d *GormDatabase) WizardRemotePointMappingOnConsumerSideByProducerSide(producerGlobalUUID string) (bool, error) {
	var consumerModel model.Consumer
	var writerModel model.Writer
	pointConsumer, err := d.WizardNewNetworkDevicePoint("system", nil, nil, nil)
	if err != nil {
		return false, fmt.Errorf("CreateNetworkDevicePoint: %s", err)
	}
	fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{GlobalUUID: &producerGlobalUUID})
	if err != nil {
		return false, fmt.Errorf("FlowNetworkClone search failure: %s", err)
	}
	streamClones, err := d.GetStreamClones(api.Args{FlowNetworkCloneUUID: &fnc.UUID})
	if err != nil {
		return false, fmt.Errorf("StreamClone search failure: %s", err)
	}

	cli := client.NewFlowClientCli(fnc.FlowIP, fnc.FlowPort, fnc.FlowToken, fnc.IsMasterSlave, fnc.GlobalUUID, model.IsFNCreator(fnc))
	producers, err := cli.GetProducers(&streamClones[0].SourceUUID)
	if err != nil {
		return false, fmt.Errorf("producer search failure: %s", err)
	}

	consumerModel.Name = "ZATSP"
	consumerModel.ProducerUUID = (*producers)[0].UUID
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

	return true, nil
}
