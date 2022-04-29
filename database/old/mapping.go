package database

//
//import (
//	"errors"
//	"fmt"
//	"github.com/NubeIO/flow-framework/api"
//	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
//	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
//	log "github.com/sirupsen/logrus"
//)
//
///*
//Mapping
//lora to bacnet >> one to one
//bacnet to edge >>  one to many
//modbus to bacnet >> many to one
//*/
//
//const (
//	edge28Bacnet  = "edge28-to-bacnetserver"
//	loRaToBACnet  = "lora-to-bacnetserver"
//	rubixIOBACnet = "rubix-io-to-bacnetserver"
//)
//
//type MappingNetwork struct {
//	Name        string
//	FromNetwork string
//	ToNetwork   string
//}
//
//var MappingNetworks = struct {
//	Edge28Bacnet  MappingNetwork
//	LoRaToBACnet  MappingNetwork
//	RubixIOBACnet MappingNetwork
//}{
//	Edge28Bacnet:  MappingNetwork{Name: edge28Bacnet, FromNetwork: model.PluginsNames.Edge28.PluginName, ToNetwork: model.PluginsNames.BACnetServer.PluginName},
//	LoRaToBACnet:  MappingNetwork{Name: loRaToBACnet, FromNetwork: model.PluginsNames.LoRa.PluginName, ToNetwork: model.PluginsNames.BACnetServer.PluginName},
//	RubixIOBACnet: MappingNetwork{Name: rubixIOBACnet, FromNetwork: model.PluginsNames.RubixIO.PluginName, ToNetwork: model.PluginsNames.BACnetServer.PluginName},
//}
//
//func selectMappingNetwork(selectedPlugin string) (pluginName string) {
//	switch selectedPlugin {
//	case MappingNetworks.Edge28Bacnet.Name:
//		pluginName = MappingNetworks.Edge28Bacnet.ToNetwork
//	case MappingNetworks.LoRaToBACnet.Name:
//		pluginName = MappingNetworks.LoRaToBACnet.ToNetwork
//	case MappingNetworks.RubixIOBACnet.Name:
//		pluginName = MappingNetworks.RubixIOBACnet.ToNetwork
//	}
//	return
//
//}
//
///*
//CreatePointMapping
//
//Producer:
//In our world a producer would be a lora temp sensor or a bacnet output
//
//Consumer:
//A consumer would be a bacnet point from the lora temp sensor or an edge-28 output as it needs to be written to from bacnet
//
//Example mapping for self mapping
//- user when making a new network selects "self-mapping"
//- when they add a new point, make a new producer as type point and select the newly made point
//- make a new consumer as type point and select the newly made point and newly made producer
//- add a writer to the producer
//*/
//func (d *GormDatabase) CreatePointMapping(body *model.PointMapping) (*model.PointMapping, error) {
//
//	network := &model.Network{}
//	device := &model.Device{}
//	point := &model.Point{}
//
//	log.Infoln("points.db.CreatePointMapping() try and make a new mapping pointMapping:", "AutoMappingFlowNetworkName:", body.AutoMappingFlowNetworkName, "AutoMappingFlowNetworkUUID:", body.AutoMappingFlowNetworkUUID)
//	device, err := d.GetDevice(body.Point.DeviceUUID, api.Args{})
//	if err != nil {
//		return nil, err
//	}
//	network, err = d.GetNetwork(device.NetworkUUID, api.Args{})
//	if err != nil {
//		return nil, err
//	}
//
//	flowNetwork, err := d.selectFlowNetwork(body.AutoMappingFlowNetworkName, body.AutoMappingFlowNetworkUUID)
//	if err != nil {
//		return nil, err
//	}
//
//	flowNetworkClone, err := d.GetOneFlowNetworkCloneByArgs(api.Args{SourceUUID: nils.NewString(flowNetwork.UUID)})
//	if err != nil || flowNetworkClone == nil {
//		log.Errorln("mapping.db.selectFlowNetwork(): missing flow network clone")
//		return nil, errors.New(fmt.Sprintf("missing flow network clone"))
//	}
//	//create a new stream or use existing
//	stream, err := d.createPointMappingStream(device.Name, network.Name, flowNetwork)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, plugin := range body.AutoMappingNetworksSelection { //j njjk,liunmjj code comment from my daughter lenny-j 3-apr-2022  //DONT DELETE its her first one :)
//		streamClone, err := d.GetStreamCloneByArg(api.Args{SourceUUID: nils.NewString(stream.UUID)})
//		if err != nil {
//			log.Errorln("mapping.db.CreatePointMapping(): failed to find stream clone with source uuid:", stream.UUID)
//			return nil, fmt.Errorf("stream-clone is missing please re-save the stream")
//		}
//		objectType := body.Point.ObjectType
//		isOutput := body.Point.IsOutput
//		if plugin != "self-mapping" {
//			//make the new network
//			network, err = d.createPointMappingNetwork(plugin)
//			if err != nil {
//				return nil, err
//			}
//			////Make a new device, check if device exits and if not make a new one
//			device, err = d.createPointMappingDevice(device.Name, network.UUID)
//			if err != nil {
//				return nil, err
//			}
//
//			point, err = d.createPointMappingPoint(objectType, body.Point.Name, device.UUID)
//			if err != nil {
//				return nil, err
//			}
//
//		} else {
//			point = body.Point
//		}
//
//		//example user wants to make edge-28 mapping to bacnet: if the source point is an output from the edge-28 then we need to make the bacnet point the producer so the edge-28 can be commanded over bacnet
//		if nils.BoolIsNil(isOutput) {
//			//example edge-28 UO to bacnet AO
//			//make the bacnet point the producer
//			producer, err := d.createPointMappingProducer(point.UUID, point.Name, device.Name, stream.UUID)
//			if err != nil {
//				return nil, err
//			}
//			//this would be the edge28 point
//			consumer, err := d.createPointMappingConsumer(body.Point.UUID, producer.Name, producer.UUID, streamClone.UUID)
//			if err != nil {
//				return nil, err
//			}
//			log.Info("Point mapping is done for TYPE-OUTPUT as the PRODUCER: ", consumer.Name)
//		} else { //if type is of input that make this the producer from the source point example: lora is the producer and bacnet is the consumer (as bacnet will read the value sent to it from lora)
//			//example edge-28 UI to bacnet AV
//			//make the edge-28 point the producer
//			producer, err := d.createPointMappingProducer(body.Point.UUID, body.Point.Name, device.Name, stream.UUID)
//			if err != nil {
//				return nil, err
//			}
//			//the consumer would be the bacnet-point
//			consumer, err := d.createPointMappingConsumer(point.UUID, producer.Name, producer.UUID, streamClone.UUID)
//			if err != nil {
//				return nil, err
//			}
//			log.Info("Point mapping is done for TYPE-INPUT as the PRODUCER: ", consumer.Name)
//		}
//
//	}
//	//make new point for the producer
//	return body, nil
//}
//
//func (d *GormDatabase) selectFlowNetwork(flowNetworkName, flowNetworkUUID string) (flowNetwork *model.FlowNetwork, err error) {
//	if flowNetworkUUID != "" {
//		flowNetwork, err = d.GetFlowNetwork(flowNetworkUUID, api.Args{})
//		if err != nil || flowNetwork == nil {
//			log.Errorln("mapping.db.selectFlowNetwork(): select by uuid missing flow network please add uuid:", flowNetworkUUID)
//			return nil, errors.New(fmt.Sprintf("failed to find a flow-network with uuid: %s", flowNetworkUUID))
//		}
//	} else {
//		name := "local"
//		if flowNetworkName != "" {
//			name = flowNetworkName
//		}
//		flowNetwork, err = d.GetOneFlowNetworkByArgs(api.Args{Name: nils.NewString(name)})
//		if err != nil || flowNetwork == nil {
//			log.Errorln("mapping.db.selectFlowNetwork(): select by name missing flow network please add name:", name)
//			return nil, errors.New(fmt.Sprintf("failed to find a flow-network with name: %s", name))
//		}
//	}
//	return
//}
//
//func (d *GormDatabase) createPointMappingStream(deviceName, networkName string, flowNetwork *model.FlowNetwork) (stream *model.Stream, err error) {
//	streamModel := &model.Stream{}
//	streamModel.FlowNetworks = []*model.FlowNetwork{flowNetwork}
//	streamModel.Name = fmt.Sprintf("%s_%s", networkName, deviceName)
//	log.Info("Try and make anew stream or select existing with name: ", streamModel.Name)
//	stream, err = d.GetStreamByArgs(api.Args{SourceUUID: nils.NewString(streamModel.Name)})
//	if stream != nil {
//		log.Warning("mapping.db.CreatePointMapping(): an existing stream with this name exists name:", stream.Name)
//		return stream, nil
//	}
//	stream, err = d.CreateStream(streamModel)
//	if stream == nil || err != nil {
//		log.Error("mapping.db.CreatePointMapping(): failed to make a new stream ", stream)
//		//return nil, nil
//	} else {
//		log.Info("mapping.db.CreatePointMapping(): stream  is created successfully: ", stream)
//		return stream, nil
//	}
//	return stream, nil
//}
//
//func (d *GormDatabase) createPointMappingStreamClone(deviceName, networkName string, flowNetwork *model.FlowNetwork) (stream *model.Stream, err error) {
//	streamModel := &model.Stream{}
//	streamModel.FlowNetworks = []*model.FlowNetwork{flowNetwork}
//	streamModel.Name = fmt.Sprintf("%s_%s", networkName, deviceName)
//	log.Info("Try and make anew stream or select existing with name: ", streamModel.Name)
//	stream, err = d.GetStreamByArgs(api.Args{SourceUUID: nils.NewString(streamModel.Name)})
//	if stream != nil {
//		log.Warning("mapping.db.CreatePointMapping(): an existing stream with this name exists name:", stream.Name)
//		return stream, nil
//	}
//	stream, err = d.CreateStream(streamModel)
//	if stream == nil || err != nil {
//		log.Error("mapping.db.CreatePointMapping(): failed to make a new stream ", stream)
//		//return nil, nil
//	} else {
//		log.Info("mapping.db.CreatePointMapping(): stream  is created successfully: ", stream)
//		return stream, nil
//	}
//	return stream, nil
//}
//
//func (d *GormDatabase) createPointMappingPoint(pointObjectType, pointName, deviceUUID string) (point *model.Point, err error) {
//	point = &model.Point{}
//	//make pnt, first check
//	point.Name = pointName
//	point.ObjectType = pointObjectType
//	point.DeviceUUID = deviceUUID
//	point, err = d.CreatePointPlugin(point)
//	if err != nil {
//		log.Errorln("mapping.db.CreatePointMapping(): failed to add point for point name:", pointName)
//		return nil, errors.New("failed to add a new point for auto mapping")
//	}
//	return point, err
//}
//
//func (d *GormDatabase) createPointMappingNetwork(plugin string) (network *model.Network, err error) {
//	plugin = selectMappingNetwork(plugin)
//	if plugin == "" {
//		log.Errorln("mapping.db.CreatePointMapping(): no valid plugin selection, try lora-to-bacnetserver")
//		return nil, errors.New("no valid mapping selecting between 2x networks")
//	}
//	network = &model.Network{}
//	network, err = d.GetNetworkByPluginName(plugin, api.Args{})
//	if network == nil {
//		network = &model.Network{}
//		network.Name = plugin
//		network.PluginPath = plugin
//		network, err = d.CreateNetwork(network, false)
//		if err != nil {
//			log.Errorln("mapping.db.CreatePointMapping(): failed to add network for plugin name:", plugin)
//			return nil, errors.New("failed to add a new network for auto mapping")
//		}
//		log.Info("mapping.db.CreatePointMapping(): an new network was made with name:", network.Name, "for plugin:", network.PluginPath)
//		return network, err
//	} else {
//		log.Info("mapping.db.CreatePointMapping(): an existing network with this name exists name:", network.Name, "for plugin:", network.PluginPath)
//		return network, err
//	}
//}
//
//func (d *GormDatabase) createPointMappingDevice(deviceName, networkUUID string) (device *model.Device, err error) {
//	//Make a new device, check if device exits and if not make a new one
//	device, existing := d.deviceNameExistsInNetwork(deviceName, networkUUID)
//	if !existing {
//		newDevice := &model.Device{}
//		newDevice.Name = deviceName
//		newDevice.NetworkUUID = networkUUID
//		device, err = d.CreateDevice(newDevice)
//		if err != nil {
//			log.Errorln("mapping.db.CreatePointMapping(): failed to add new device:", newDevice.Name)
//			return nil, err
//		}
//		log.Info("mapping.db.CreatePointMapping(): an new device was added name::", deviceName, "for network_uuid:", networkUUID)
//		return device, nil
//	} else {
//		log.Info("mapping.db.CreatePointMapping(): an existing device with this name exists name:", deviceName, "for network_uuid:", networkUUID)
//		return device, nil
//	}
//}
//
//func (d *GormDatabase) createPointMappingProducer(pointUUID, pointName, deviceName, streamUUID string) (producer *model.Producer, err error) {
//
//	//make a producer
//	producer = &model.Producer{}
//	producer.Name = fmt.Sprintf("%s_%s", deviceName, pointName)
//	producer.StreamUUID = streamUUID
//	producer.ProducerThingUUID = pointUUID
//	producer.ProducerThingClass = "point"
//	producer.ProducerApplication = "mapping"
//	producer, err = d.CreateProducer(producer)
//	if err != nil {
//		return nil, fmt.Errorf("createPointMappingProducer() producer creation failure: %s", err)
//	}
//	log.Info("Producer point is created successfully: ", producer)
//	return
//
//}
//
//func (d *GormDatabase) createPointMappingConsumer(pointUUID, producerName, producerUUID, streamCloneUUID string) (consumer *model.Consumer, err error) {
//
//	consumer = &model.Consumer{}
//	consumer.Name = producerName
//	consumer.ProducerUUID = producerUUID
//	consumer.ConsumerApplication = "mapping"
//	consumer.StreamCloneUUID = streamCloneUUID
//	consumer, err = d.CreateConsumer(consumer)
//	if err != nil {
//		return nil, fmt.Errorf("createPointMappingConsumer() point consumer creation failure: %s", err)
//	}
//	log.Info("Point consumer is created successfully: ", consumer.Name)
//
//	writer := &model.Writer{}
//	writer.ConsumerUUID = consumer.UUID
//	writer.WriterThingClass = "point"
//	writer.WriterThingType = "temp"
//	writer.WriterThingUUID = pointUUID
//	writer, err = d.CreateWriter(writer)
//	if err != nil {
//		return nil, fmt.Errorf("writer creation failure: %s", err)
//	}
//	log.Info("Writer is created successfully: ", writer.WriterThingName)
//	return
//
//}
