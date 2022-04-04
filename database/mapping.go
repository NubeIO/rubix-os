package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

func (d *GormDatabase) CreatePointMapping(body *model.PointMapping) (*model.PointMapping, error) {

	device, err := d.GetDevice(body.Point.DeviceUUID, api.Args{})
	if err != nil {
		return nil, err
	}

	network, err := d.GetNetwork(device.NetworkUUID, api.Args{})
	if err != nil {
		return nil, err
	}

	flowNetwork, err := d.GetOneFlowNetworkByArgs(api.Args{Name: nils.NewString("local")})
	if err != nil {
		log.Errorln("mapping.db.CreatePointMapping(): missing flow network please add")
		return nil, errors.New("please add a flow-network named local")
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

	for _, plugin := range body.PluginsList { //j njjk,liunmjj code comment from my daughter lenny-j 3-apr-2022  //DONT TOUCH its her first one
		//make the new network
		network = &model.Network{}
		network, err = d.GetNetworkByPluginName(plugin, api.Args{})
		networkDidExist := false
		if network == nil {
			network = &model.Network{}
			network.Name = plugin
			network.PluginPath = plugin
			network, err = d.CreateNetwork(network, false)
			if err != nil {
				log.Errorln("mapping.db.CreatePointMapping(): failed to add network for plugin name:", plugin)
				return nil, errors.New("failed to add a new network for auto mapping")
			}
			networkDidExist = true
			log.Info("mapping.db.CreatePointMapping(): an new network was made with name:", network.Name, "for plugin:", network.PluginPath)
		} else {
			log.Info("mapping.db.CreatePointMapping(): an existing network with this name exists name:", network.Name, "for plugin:", network.PluginPath)
		}
		if networkDidExist {

		}
		deviceName := device.Name //
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
