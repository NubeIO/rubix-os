package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

/*
Mapping
lora to bacnet >> one to one
bacnet to edge >>  one to many
modbus to bacnet >> many to one
*/

const (
	edge28Bacnet = "edge28-to-bacnetserver"
	loRaToBACnet = "lora-to-bacnetserver"
)

type MappingNetwork struct {
	Name        string
	FromNetwork string
	ToNetwork   string
}

var MappingNetworks = struct {
	Edge28Bacnet MappingNetwork
	LoRaToBACnet MappingNetwork
}{
	Edge28Bacnet: MappingNetwork{Name: edge28Bacnet, FromNetwork: model.PluginsNames.Edge28.PluginName, ToNetwork: model.PluginsNames.BACnetServer.PluginName},
	LoRaToBACnet: MappingNetwork{Name: loRaToBACnet, FromNetwork: model.PluginsNames.LoRa.PluginName, ToNetwork: model.PluginsNames.BACnetServer.PluginName},
}

func selectMappingNetwork(selectedPlugin string) (pluginName string) {
	switch selectedPlugin {
	case MappingNetworks.Edge28Bacnet.Name:
		pluginName = MappingNetworks.Edge28Bacnet.ToNetwork
	case MappingNetworks.LoRaToBACnet.Name:
		pluginName = MappingNetworks.LoRaToBACnet.ToNetwork
	}
	return

}

func (d *GormDatabase) selectFlowNetwork(flowNetworkName, flowNetworkUUID string) (flowNetwork *model.FlowNetwork, err error) {
	if flowNetworkUUID != "" {
		flowNetwork, err = d.GetFlowNetwork(flowNetworkUUID, api.Args{})
		if err != nil || flowNetwork == nil {
			log.Errorln("mapping.db.selectFlowNetwork(): select by uuid missing flow network please add uuid:", flowNetworkUUID)
			return nil, errors.New(fmt.Sprintf("failed to find a flow-network with uuid: %s", flowNetworkUUID))
		}
	} else {
		name := "local"
		if flowNetworkName != "" {
			name = flowNetworkName
		}
		flowNetwork, err = d.GetOneFlowNetworkByArgs(api.Args{Name: nils.NewString(name)})
		if err != nil || flowNetwork == nil {
			log.Errorln("mapping.db.selectFlowNetwork(): select by name missing flow network please add name:", name)
			return nil, errors.New(fmt.Sprintf("failed to find a flow-network with name: %s", name))
		}
	}
	return
}

func (d *GormDatabase) CreatePointMapping(body *model.PointMapping) (*model.PointMapping, error) {
	log.Infoln("points.db.CreatePointMapping() try and make a new mapping pointMapping:", "AutoMappingFlowNetworkName:", body.AutoMappingFlowNetworkName, "AutoMappingFlowNetworkUUID:", body.AutoMappingFlowNetworkUUID)
	device, err := d.GetDevice(body.Point.DeviceUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	network, err := d.GetNetwork(device.NetworkUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	flowNetwork, err := d.selectFlowNetwork(body.AutoMappingFlowNetworkName, body.AutoMappingFlowNetworkUUID)
	if err != nil {
		return nil, err
	}
	streamModel := &model.Stream{}
	streamModel.FlowNetworks = []*model.FlowNetwork{flowNetwork}
	streamModel.Name = fmt.Sprintf("%s_%s", network.Name, device.Name)
	log.Info("Try and make stream: ", streamModel.Name)
	stream, _ := d.CreateStream(streamModel)
	if stream != nil {
		log.Warning("mapping.db.CreatePointMapping(): an existing stream with this name exists name:", stream.Name)
	} else {
		log.Info("Stream is created successfully: ", stream)
	}
	for _, plugin := range body.AutoMappingNetworksSelection { //j njjk,liunmjj code comment from my daughter lenny-j 3-apr-2022  //DONT DELETE its her first one :)
		//make the new network
		plugin = selectMappingNetwork(plugin)
		network = &model.Network{}
		network, err = d.GetNetworkByPluginName(plugin, api.Args{})
		if network == nil {
			network = &model.Network{}
			network.Name = plugin
			network.PluginPath = plugin
			network, err = d.CreateNetwork(network, false)
			if err != nil {
				log.Errorln("mapping.db.CreatePointMapping(): failed to add network for plugin name:", plugin)
				return nil, errors.New("failed to add a new network for auto mapping")
			}
			log.Info("mapping.db.CreatePointMapping(): an new network was made with name:", network.Name, "for plugin:", network.PluginPath)
		} else {
			log.Info("mapping.db.CreatePointMapping(): an existing network with this name exists name:", network.Name, "for plugin:", network.PluginPath)
		}
		deviceName := device.Name
		//Make a new device, check if device exits and if not make a new one
		device, existing := d.deviceNameExistsInNetwork(device.Name, network.UUID)
		deviceUUID := ""
		if !existing {
			newDevice := &model.Device{}
			newDevice.Name = deviceName
			newDevice.NetworkUUID = network.UUID
			device, err = d.CreateDevice(newDevice)
			if err != nil {
				log.Errorln("mapping.db.CreatePointMapping(): failed to add new device:", newDevice.Name)
			}
			deviceUUID = device.UUID
			log.Info("mapping.db.CreatePointMapping(): an new device was added name::", deviceName, "for network_uuid:", network.UUID)

		} else {
			deviceUUID = device.UUID
			log.Info("mapping.db.CreatePointMapping(): an existing device with this name exists name:", deviceName, "for network_uuid:", network.UUID)
		}
		//make pnt, first check
		point := &model.Point{}
		point.Name = body.Point.Name
		point.DeviceUUID = deviceUUID
		point, err = d.CreatePointPlugin(point)
		if err != nil {
			log.Errorln("mapping.db.CreatePointMapping(): failed to add point for point name:", body.Point.Name)
			return nil, errors.New("failed to add a new point for auto mapping")
		}

		//make a producer
		producer := &model.Producer{}
		producer.Name = fmt.Sprintf("%s_%s", device.Name, body.Point.Name)
		producer.StreamUUID = stream.UUID
		producer.ProducerThingUUID = body.Point.UUID
		producer.ProducerThingClass = "point"
		producer.ProducerApplication = "mapping"
		producer, err := d.CreateProducer(producer)
		if err != nil {
			return nil, fmt.Errorf("producer creation failure: %s", err)
		}
		log.Info("Producer point is created successfully: ", producer)

		streamClone, err := d.GetStreamCloneByArg(api.Args{SourceUUID: nils.NewString(stream.UUID)})
		if err != nil {
			log.Errorln("mapping.db.CreatePointMapping(): failed to find stream clone with source uuid:", stream.UUID)
			return nil, fmt.Errorf("failed to get stream-clone: %s", err)
		}

		consumer := &model.Consumer{}
		consumer.Name = producer.Name
		consumer.ProducerUUID = producer.UUID
		consumer.ConsumerApplication = "mapping"
		consumer.StreamCloneUUID = streamClone.UUID
		consumer, err = d.CreateConsumer(consumer)
		if err != nil {
			return nil, fmt.Errorf("point consumer creation failure: %s", err)
		}
		log.Info("Point consumer is created successfully: ", consumer.Name)

		writer := &model.Writer{}
		writer.ConsumerUUID = consumer.UUID
		writer.WriterThingClass = "point"
		writer.WriterThingType = "temp"
		writer.WriterThingUUID = point.UUID
		writer, err = d.CreateWriter(writer)
		if err != nil {
			return nil, fmt.Errorf("writer creation failure: %s", err)
		}
		log.Info("Writer is created successfully: ", writer.WriterThingName)

	}
	//make new point for the producer
	return body, nil
}
