package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

/*
WizardMasterSlavePointMapping
On producer (edge) side:
- Create network > device > point
- Create flow_network
- Create stream and assign that stream under that flow_network
- Create producer and link it with point

On consumer (normally server) side:
- Create network > device > point
- Get flow_network_clone by producer's global_uuid
- Get stream_clone by flow_network_clone.uuid
- Get producer from stream_clone.source_uuid
- Create consumer and link it with stream_clone & producer
- Create writer under the consumer
*/
func (d *GormDatabase) WizardMasterSlavePointMapping() (bool, error) {
	flow, err := d.wizardMasterSlavePointMappingOnProducerSide()
	if err != nil {
		return false, err
	}

	deviceInfo, err := deviceinfo.GetDeviceInfo()
	if err != nil {
		return false, fmt.Errorf("GetDeviceInfo: %s", deviceInfo)
	}

	cli := client.NewFlowClientCliFromFN(flow)
	_, err = cli.WizardMasterSlavePointMappingOnConsumerSideByProducerSide(deviceInfo.GlobalUUID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *GormDatabase) wizardMasterSlavePointMappingOnProducerSide() (*model.FlowNetwork, error) {
	var flowNetworkModel model.FlowNetwork
	var streamModel model.Stream
	var producerModel model.Producer

	pointModel := model.Point{}
	pointModel.Name = "ZATSP-PRO"
	point, err := d.WizardNewNetworkDevicePoint("system", nil, nil, &pointModel)
	if err != nil {
		return nil, fmt.Errorf("CreateNetworkDevicePoint: %s", err)
	}

	flowNetworkModel.IsMasterSlave = boolean.NewTrue()
	flowNetworkModel.IsMasterSlave = boolean.NewTrue()
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
	log.Info("Producer is created successfully: ", producer)

	return fn, nil
}

func (d *GormDatabase) WizardMasterSlavePointMappingOnConsumerSideByProducerSide(producerGlobalUUID string) (bool, error) {
	var consumerModel model.Consumer
	var writerModel model.Writer
	point, err := d.WizardNewNetworkDevicePoint("system", nil, nil, nil)
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

	cli := client.NewFlowClientCliFromFNC(fnc)
	url := urls.PluralUrlByArg(urls.ProducerUrl, "stream_uuid", streamClones[0].SourceUUID)
	rawProducers, err := cli.GetQueryMarshal(url, []model.Producer{})
	if err != nil {
		return false, err
	}
	producers := rawProducers.(*[]model.Producer)

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
	writerModel.WriterThingUUID = point.UUID
	writer, err := d.CreateWriter(&writerModel)
	if err != nil {
		return false, fmt.Errorf("writer creation failure: %s", err)
	}
	log.Info("Writer is created successfully: ", writer)

	return true, nil
}
