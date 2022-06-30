package main

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/btypes/segmentation"
	"github.com/NubeDev/bacnet/network"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) bacnetNetworkInit() {

	networks, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{WithDevices: true})
	if err != nil {
		return
	}
	for _, net := range networks {
		err := inst.bacnetNetwork(net)
		if err != nil {
			log.Errorln("bacnet-master init network error:", err)
			continue
		}
		for _, dev := range net.Devices {
			err := inst.bacnetDevice(dev)
			if err != nil {
				log.Errorln("bacnet-master init device error:", err)
				continue
			}
		}
	}

}

func (inst *Instance) initBacStore() {
	inst.BacStore = network.NewStore()
}

// bacnetNetwork add or update an instance a bacnet network that is cached in bacnet lib
func (inst *Instance) bacnetNetwork(net *model.Network) error {
	bacnetNet := &network.Network{
		Interface: net.NetworkInterface,
		Port:      integer.NonNil(net.Port),
		StoreID:   net.UUID,
	}
	return inst.BacStore.UpdateNetwork(net.UUID, bacnetNet)
}

// getBacnetNetwork get an instance of a created bacnet network that is cached in bacnet lib
func (inst *Instance) getBacnetNetwork(networkUUID string) (*network.Network, error) {
	return inst.BacStore.GetNetwork(networkUUID)
}

// closeBacnetNetwork delete the instance of a created bacnet network that is cached in bacnet lib
func (inst *Instance) closeBacnetNetwork(networkUUID string) (bool, error) {
	net, err := inst.BacStore.GetNetwork(networkUUID)
	if err != nil {
		return false, err
	}
	net.NetworkClose()
	return true, nil
}

// getBacnetDevice get an instance of a created bacnet device that is cached in bacnet lib
func (inst *Instance) getBacnetDevice(deviceUUID string) (*network.Device, error) {
	return inst.BacStore.GetDevice(deviceUUID)
}

// bacnetDevice add or update an instance of a created bacnet device that is cached in bacnet lib
func (inst *Instance) bacnetDevice(dev *model.Device) error {
	max := intToUnit32(integer.NonNil(dev.MaxADPU))
	seg := uint32(setSegmentation(dev.Segmentation))
	d := &network.Device{
		Ip:            dev.CommonIP.Host,
		Port:          dev.CommonIP.Port,
		DeviceID:      integer.NonNil(dev.DeviceObjectId),
		StoreID:       dev.UUID,
		NetworkNumber: integer.NonNil(dev.NetworkNumber),
		MacMSTP:       integer.NonNil(dev.DeviceMac),
		MaxApdu:       max,
		Segmentation:  seg,
	}

	net, _ := inst.getBacnetNetwork(dev.NetworkUUID)
	return inst.BacStore.UpdateDevice(dev.UUID, net, d)
}

// getDev get an instance of a created bacnet device that is cached in bacnet lib
func (inst *Instance) doReadValue(pnt *model.Point, networkUUID, deviceUUID string) (float64, error) {
	object, _, isBool, _ := setObjectType(pnt.ObjectType)
	bp := &network.Point{
		ObjectID:         btypes.ObjectInstance(integer.NonNil(pnt.ObjectId)),
		ObjectType:       object,
		WriteValue:       nil,
		WriteNull:        false,
		WritePriority:    0,
		ReadPresentValue: false,
		ReadPriority:     false,
	}
	// get network
	net, err := inst.getBacnetNetwork(networkUUID)
	if err != nil {
		return 0, err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetDevice(deviceUUID)
	if err != nil {
		return 0, err
	}
	var outValue float64
	if isBool {
		readBool, err := dev.PointReadBool(bp)
		if err != nil {
			log.Errorln("bacnet-master-read-bool:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " error:", err)
			return 0, err
		}
		outValue = Unit32ToFloat64(readBool)

	} else {
		readFloat32, err := dev.PointReadFloat32(bp)
		if err != nil {
			log.Errorln("bacnet-master-read-analogue:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " error:", err)
			return 0, err
		}
		outValue = float32ToFloat64(readFloat32)
	}
	log.Infoln("bacnet-master-POINT-READ:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", outValue)
	return outValue, nil
}

func (inst *Instance) doWrite(pnt *model.Point, networkUUID, deviceUUID string) error {
	val := float.NonNil(pnt.WriteValue)
	object, isWrite, isBool, _ := setObjectType(pnt.ObjectType)
	writePriority := integer.NonNil(pnt.WritePriority)
	if writePriority == 0 {
		writePriority = 16
	}
	bp := &network.Point{
		ObjectID:         btypes.ObjectInstance(integer.NonNil(pnt.ObjectId)),
		ObjectType:       object,
		WriteNull:        false,
		WritePriority:    uint8(writePriority),
		ReadPresentValue: false,
		ReadPriority:     false,
	}
	net, err := inst.getBacnetNetwork(networkUUID)
	if err != nil {
		return err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetDevice(deviceUUID)
	if err != nil {
		return err
	}

	if isWrite {
		if isBool {
			err = dev.PointWriteBool(bp, float64ToUnit32(val))
			if err != nil {
				log.Errorln("bacnet-master-write-bool:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority, " error:", err)
				return err
			}
		} else {
			err = dev.PointWriteAnalogue(bp, float64ToFloat32(val))
			if err != nil {
				log.Errorln("bacnet-master-write-analogue:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority, " error:", err)
				return err
			}
		}
	}
	log.Infoln("bacnet-master-POINT-WRITE:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority)
	return nil

}

func setObjectType(object string) (obj btypes.ObjectType, isWritable, isBool bool, asString string) {
	object = strcase.ToSnake(object)
	switch object {
	case "analog_input":
		return btypes.AnalogInput, false, false, "analog_input"
	case "analog_output":
		return btypes.AnalogOutput, true, false, "analog_output"
	case "analog_value":
		return btypes.AnalogValue, true, false, "analog_value"
	case "binary_input":
		return btypes.BinaryInput, false, true, "binary_input"
	case "binary_output":
		return btypes.BinaryOutput, true, true, "binary_output"
	case "binary_value":
		return btypes.BinaryValue, true, true, "binary_value"
	case "multi_state_input":
		return btypes.MultiStateInput, false, false, "multi_state_input"
	case "multi_state_output":
		return btypes.MultiStateOutput, true, false, "multi_state_output"
	case "multi_state_value":
		return btypes.MultiStateValue, true, false, "multi_state_value"
	default:
		return btypes.AnalogInput, false, false, "analog_input"
	}
}

func convertObjectType(object btypes.ObjectType) string {
	switch object {
	case btypes.AnalogInput:
		return "analog_input"
	case btypes.AnalogOutput:
		return "analog_output"
	case btypes.AnalogValue:
		return "analog_value"
	case btypes.BinaryInput:
		return "binary_input"
	case btypes.BinaryOutput:
		return "binary_output"
	case btypes.BinaryValue:
		return "binary_value"
	case btypes.MultiStateInput:
		return "multi_state_input"
	case btypes.MultiStateOutput:
		return "multi_state_output"
	case btypes.MultiStateValue:
		return "multi_state_value"
	}
	return "analog_input"
}

type SegmentedType string

const (
	SegmentedBoth     SegmentedType = "segmentation_both"
	SegmentedTransmit SegmentedType = "segmentation_transmit"
	SegmentedReceive  SegmentedType = "segmentation_receive"
	NoSegmentation    SegmentedType = "no_segmentation"
)

func setSegmentation(SegmentedType string) (out segmentation.SegmentedType) {
	switch SegmentedType {
	case string(SegmentedBoth):
		return segmentation.SegmentedBoth
	case string(SegmentedTransmit):
		return segmentation.SegmentedTransmit
	case string(SegmentedReceive):
		return segmentation.SegmentedReceive
	case string(NoSegmentation):
		return segmentation.NoSegmentation
	default:
		return segmentation.NoSegmentation
	}
}

func convertSegmentation(segmentedType segmentation.SegmentedType) SegmentedType {
	switch segmentedType {
	case segmentation.SegmentedBoth:
		return SegmentedBoth
	case segmentation.SegmentedTransmit:
		return SegmentedTransmit
	case segmentation.SegmentedReceive:
		return SegmentedReceive
	case segmentation.NoSegmentation:
		return NoSegmentation
	default:
		return NoSegmentation
	}
}

func (inst *Instance) doWriteBool(networkUUID, deviceUUID string, pnt *network.Point, value uint32) error {
	net, err := inst.getBacnetNetwork(networkUUID)
	if err != nil {
		return err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetDevice(deviceUUID)
	if err != nil {
		return err
	}

	err = dev.PointWriteBool(pnt, value)
	if err != nil {
		return err
	}
	return nil

}

func (inst *Instance) whoIs(networkUUID string, opts *bacnet.WhoIsOpts, addDevices bool) (resp []btypes.Device, err error) {
	net, err := inst.getBacnetNetwork(networkUUID)
	if err != nil {
		return nil, err
	}
	go net.NetworkRun() //TODO: do we need to defer a NetworkClose()?
	defer net.NetworkClose()
	devices, err := net.Whois(opts)
	if err != nil {
		return nil, err
	}
	if addDevices {
		for _, device := range devices {
			newDevice := &model.Device{
				CommonName:     model.CommonName{Name: fmt.Sprintf("deviceId_%d_networkNum_%d", device.ID.Instance, device.NetworkNumber)},
				DeviceMac:      integer.New(device.MacMSTP),
				DeviceObjectId: integer.New(int(device.ID.Instance)),
				NetworkNumber:  integer.New(device.NetworkNumber),
				MaxADPU:        integer.New(int(device.MaxApdu)),
				Segmentation:   string(convertSegmentation(segmentation.SegmentedType(device.Segmentation))),
				NetworkUUID:    networkUUID,
			}
			addDevice, err := inst.addDevice(newDevice)
			if err != nil {
				log.Errorf("failed to add a new device from whois %d", device.ID.Instance)
				return nil, err
			}
			log.Infof("added new device from whois %s", addDevice.Name)
		}
	}

	return devices, err
}

func (inst *Instance) devicePoints(deviceUUID string, addPoints, writeablePoints bool) (resp []*network.PointDetails, err error) {
	getNetwork, err := inst.db.GetNetworkByDeviceUUID(deviceUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	net, err := inst.getBacnetNetwork(getNetwork.UUID)
	if err != nil {
		return nil, err
	}
	go net.NetworkRun() //TODO: do we need to defer a NetworkClose()?
	defer net.NetworkClose()
	dev, err := inst.getBacnetDevice(deviceUUID)
	if err != nil {
		return nil, err
	}

	resp, err = dev.GetDevicePoints(btypes.ObjectInstance(dev.DeviceID))
	if err != nil {
		return nil, err
	}
	if addPoints {
		for _, pnt := range resp {
			_, isWrite, _, objectType := setObjectType(convertObjectType(pnt.ObjectType))
			writeMode := model.ReadOnly
			if isWrite && writeablePoints {
				writeMode = model.WriteOnceThenRead
			}
			newPnt := &model.Point{
				CommonName: model.CommonName{Name: pnt.Name},
				DeviceUUID: deviceUUID,
				ObjectType: objectType,
				ObjectId:   integer.New(int(pnt.ObjectID)),
				WriteMode:  writeMode,
			}
			point, err := inst.addPoint(newPnt)
			if err != nil {
				log.Errorf("failed to add a new point from discover points %s", pnt.Name)
				return nil, err
			}
			log.Infof("added new point from discover points%s", point.Name)
		}
	}
	return
}

func intToUnit32(value int) uint32 {
	var y = uint32(value)
	return y
}

func Unit32ToFloat64(value uint32) float64 {
	var y = float64(value)
	return y
}
func float32ToFloat64(value float32) float64 {
	var y = float64(value)
	return y
}

func float64ToUnit32(value float64) uint32 {
	var y = uint32(value)
	return y
}

func float64ToFloat32(value float64) float32 {
	var y = float32(value)
	return y
}
