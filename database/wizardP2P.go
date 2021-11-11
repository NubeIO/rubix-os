package database

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/client"
	"github.com/NubeDev/flow-framework/urls"
	log "github.com/sirupsen/logrus"
)

/*
WizardP2PMapping
On producer (edge) side:
- Create network > device > point
- Create schedule
- Update local_storage_flow_network using local credential (we read this and send this credentials to it's paired device)
- Create flow_network using remote credential
- Create stream and assign that stream under that flow_network
- Create producers and link it with it's respective point & schedule

On consumer (normally server) side:
- Create network > device > point
- Create schedule
- Get flow_network_clone by producer's global_uuid
- Get stream_clone by flow_network_clone.uuid
- Get producers from stream_clone.source_uuid
- Create consumers and link it with stream_clone & producers
- Create writer under the consumer
*/
func (d *GormDatabase) WizardP2PMapping(body *model.P2PBody) (bool, error) {
	flow, err := d.wizardP2PPointMappingOnProducerSide(body)
	if err != nil {
		return false, err
	}

	deviceInfo, err := d.GetDeviceInfo()
	if err != nil {
		return false, fmt.Errorf("GetDeviceInfo: %s", deviceInfo)
	}

	cli := client.NewFlowClientCli(flow.FlowIP, flow.FlowPort, flow.FlowToken, flow.IsMasterSlave, flow.GlobalUUID, model.IsFNCreator(flow))
	url := fmt.Sprintf("/api/database/wizard/mapping/p2p/points/consumer/%s", deviceInfo.GlobalUUID)
	_, err = cli.PostQuery(url, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *GormDatabase) wizardP2PPointMappingOnProducerSide(body *model.P2PBody) (*model.FlowNetwork, error) {
	var flowNetworkModel model.FlowNetwork
	var streamModel model.Stream
	var producerModel model.Producer
	var pointModel model.Point
	var scheduleModel model.Schedule

	pointModel.Name = "ZATSP-PRO"
	scheduleModel.Name = "SCH-PRO"

	point, err := d.WizardNewNetworkDevicePoint("system", nil, nil, &pointModel)
	if err != nil {
		return nil, fmt.Errorf("CreateNetworkDevicePoint: %s", err)
	}

	schedule, err := d.CreateSchedule(&scheduleModel)
	if err != nil {
		return nil, errors.New("CreateSchedule creation failure")
	}

	lsfn := model.LocalStorageFlowNetwork{}
	if body.LocalFlowIP != nil {
		lsfn.FlowIP = *body.LocalFlowIP
	}
	if body.LocalFlowPort != nil {
		lsfn.FlowPort = *body.LocalFlowPort
	}
	if body.LocalFlowUsername != nil {
		lsfn.FlowUsername = *body.LocalFlowUsername
	}
	if body.LocalFlowPassword != nil {
		lsfn.FlowPassword = *body.LocalFlowPassword
	}
	_, err = d.UpdateLocalStorageFlowNetwork(&lsfn)
	if err != nil {
		return nil, err
	}

	flowNetworkModel.Name = "FlowNetwork"
	flowNetworkModel.FlowIP = body.RemoteFlowIP
	flowNetworkModel.FlowPort = body.RemoteFlowPort
	flowNetworkModel.FlowUsername = body.RemoteFlowUsername
	flowNetworkModel.FlowPassword = body.RemoteFlowPassword
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
	producerModel.ProducerThingUUID = point.UUID
	producerModel.ProducerThingClass = "point"
	producerModel.ProducerApplication = "mapping"
	producer, err := d.CreateProducer(&producerModel)
	if err != nil {
		return nil, fmt.Errorf("producer creation failure: %s", err)
	}
	log.Info("Producer point is created successfully: ", producer)

	producerModel = model.Producer{}
	producerModel.Name = "SCH"
	producerModel.StreamUUID = stream.UUID
	producerModel.ProducerThingUUID = schedule.UUID
	producerModel.ProducerThingClass = "schedule"
	producerModel.ProducerApplication = "schedule"
	producerSchedule, err := d.CreateProducer(&producerModel)
	if err != nil {
		return nil, errors.New("schedule producer creation failure")
	}
	log.Info("Producer schedule is created successfully: ", producerSchedule)

	return fn, nil
}

func (d *GormDatabase) WizardP2PMappingOnConsumerSideByProducerSide(producerGlobalUUID string) (bool, error) {
	var consumerModel model.Consumer
	var writerModel model.Writer
	var scheduleModel model.Schedule

	scheduleModel.Name = "SCH"

	point, err := d.WizardNewNetworkDevicePoint("system", nil, nil, nil)
	if err != nil {
		return false, fmt.Errorf("CreateNetworkDevicePoint: %s", err)
	}
	schedule, err := d.CreateSchedule(&scheduleModel)
	if err != nil {
		return false, errors.New("CreateSchedule creation failure")
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
	rawProducers, err := cli.GetQueryMarshal(urls.ProducerURLWithStream(streamClones[0].SourceUUID), []model.Producer{})
	if err != nil {
		return false, err
	}
	producers := rawProducers.(*[]model.Producer)

	consumerModel.Name = "ZATSP"
	consumerModel.ProducerUUID = (*producers)[0].UUID
	consumerModel.ConsumerApplication = "mapping"
	consumerModel.StreamCloneUUID = streamClones[0].UUID
	pointConsumer, err := d.CreateConsumer(&consumerModel)
	if err != nil {
		return false, fmt.Errorf("point consumer creation failure: %s", err)
	}
	log.Info("Point consumer is created successfully: ", pointConsumer)

	writerModel.ConsumerUUID = pointConsumer.UUID
	writerModel.WriterThingClass = "point"
	writerModel.WriterThingType = "temp"
	writerModel.WriterThingUUID = point.UUID
	writer, err := d.CreateWriter(&writerModel)
	if err != nil {
		return false, fmt.Errorf("writer creation failure: %s", err)
	}
	log.Info("Writer is created successfully: ", writer)

	consumerModel.Name = "schedule"
	consumerModel.ProducerUUID = (*producers)[1].UUID //todo check indexing
	consumerModel.ConsumerApplication = "schedule"
	consumerModel.StreamCloneUUID = streamClones[0].UUID
	scheduleConsumer, err := d.CreateConsumer(&consumerModel)
	if err != nil {
		return false, errors.New("schedule consumer creation failure")
	}
	log.Info("Schedule consumer is created successfully: ", scheduleConsumer)

	writerModel.ConsumerUUID = scheduleConsumer.UUID
	writerModel.WriterThingClass = "schedule"
	writerModel.WriterThingType = "schedule"
	writerModel.WriterThingUUID = schedule.UUID
	writerSchedule, err := d.CreateWriter(&writerModel)
	if err != nil {
		return false, errors.New("schedule writer creation failure")
	}
	log.Info("Schedule writer is created successfully: ", writerSchedule)
	return true, nil
}
