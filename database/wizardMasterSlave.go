package database

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/client"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
)

/*
WizardRemotePointMapping
On producer (edge) side:
- Create network, device, point
- Create flow_network
- Create stream and assign that stream under that flow_network
- Create producer and link it with point

On consumer (normally server) side:
- Create network, device, point
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

	deviceInfo, err := d.GetDeviceInfo()
	if err != nil {
		return false, fmt.Errorf("GetDeviceInfo: %s", deviceInfo)
	}

	cli := client.NewFlowClientCli(flow.FlowIP, flow.FlowPort, flow.FlowToken, flow.IsMasterSlave, flow.GlobalUUID, model.IsFNCreator(flow))
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

	point := model.Point{}
	point.Name = "ZATSP-PRO"
	pointProducer, err := d.WizardNewNetworkDevicePoint("system", nil, nil, &point)
	if err != nil {
		return nil, fmt.Errorf("CreateNetworkDevicePoint: %s", err)
	}

	flowNetworkModel.IsMasterSlave = utils.NewTrue()
	flowNetworkModel.IsMasterSlave = utils.NewTrue()
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

	return &flowNetworkModel, nil
}

func (d *GormDatabase) WizardMasterSlavePointMappingOnConsumerSideByProducerSide(producerGlobalUUID string) (bool, error) {
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

func (d *GormDatabase) WizardNewNetworkDevicePoint(plugin string, network *model.Network, device *model.Device, point *model.Point) (*model.Point, error) {
	if network == nil {
		network = &model.Network{
			TransportType: "ip",
		}
	}
	if device == nil {
		device = &model.Device{}
		device.TransportType = "ip"
	}
	if point == nil {
		point = &model.Point{
			IsProducer: utils.NewTrue(),
			ObjectType: "analogValue",
		}
		point.Name = "ZATSP"
	}
	if point.IsProducer != nil {
		point.IsProducer = utils.NewTrue()
	}
	if point.ObjectType == "" {
		point.ObjectType = "analogValue"
	}
	if point.Name == "" {
		point.Name = "ZATSP"
	}

	p, err := d.GetPluginByPath(plugin)
	if err != nil {
		return nil, errors.New("not valid plugin found")
	}

	network.PluginConfId = p.UUID
	n, err := d.CreateNetwork(network)
	if err != nil {
		return nil, fmt.Errorf("network creation failure: %s", err)
	}
	log.Info("Created a Network")

	device.NetworkUUID = n.UUID
	dev, err := d.CreateDevice(device)
	if err != nil {
		return nil, fmt.Errorf("device creation failure: %s", err)
	}
	log.Info("Created a Device: ", dev)

	point.DeviceUUID = dev.UUID
	_, err = d.CreatePoint(point, "")
	if err != nil {
		return nil, fmt.Errorf("consumer point creation failure: %s", err)
	}
	log.Info("Created a Point for Consumer", point)
	return point, nil
}
