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

	//pass in networks
	//for each

	//select the local flow-network

	//make a new stream with the network_name_device_name
	//add each point for the device under this stream

	//example map from modbus to system and bacnet
	//for _, plugin := range body.PluginsToMap {
	//	network, err := d.GetNetworksByPluginName(plugin, api.Args{})
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//}

	//body.Point
	//network, err := d.GetNetworkByPointUUID(body.Point, api.Args{WithDevices: true})
	//if err != nil {
	//	return nil, err
	//}

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
	//d.GetStreams()

	stream, _ := d.CreateStream(streamModel)
	if stream != nil {
		log.Info("mapping.db.CreatePointMapping(): an existing stream with this name exists name:", stream.Name)
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
		}
		if networkDidExist {

		}
		//Make a new device, check if device exits and if not make a new one
		device, existing := d.deviceNameExistsInNetwork(device.Name, network.UUID)
		deviceUUID := device.UUID
		if !existing {
			newDevice := &model.Device{}
			newDevice.Name = device.Name
			newDevice.NetworkUUID = network.UUID
			device, err = d.CreateDevice(newDevice)
			if err != nil {
				log.Errorln("mapping.db.CreatePointMapping(): failed to add new device:", newDevice.Name)
			}
			deviceUUID = device.UUID
		}
		//make pnt, first check
		point := &model.Point{}
		point.Name = body.Point.Name
		point.DeviceUUID = deviceUUID
		fmt.Println(99999, point.Name, deviceUUID)
		point, err = d.CreatePoint(point, false)
		if err != nil {
			log.Errorln("mapping.db.CreatePointMapping(): failed to add point for point name:", body.Point.Name)
			return nil, errors.New("failed to add a new point for auto mapping")
		}

		//fmt.Println(existing, "existing point", pointName)
		//if !existing {
		//	point := &model.Point{}
		//	point.Name = body.Point.Name
		//	point.DeviceUUID = device.UUID
		//	point, err = d.CreatePoint(point, false)
		//} else {
		//	log.Errorln("mapping.db.CreatePointMapping(): failed to create a new point as an existing point with the same name exists", pointName)
		//	return nil, errors.New("failed to create a new point as a point with same name exists")
		//}
		//
		//fmt.Println(network)

		fmt.Println(333)
		fmt.Println(network, err)
		fmt.Println(444)

	}

	//make new point for the producer

	return body, nil
}
