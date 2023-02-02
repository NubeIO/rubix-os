package main

import (
	"errors"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/btypes/priority"
	"github.com/NubeDev/bacnet/btypes/segmentation"
	"github.com/NubeDev/bacnet/network"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"reflect"
)

const (
	defaultPort = 47809
)

func (inst *Instance) bacnetNetworkInit() {
	networks, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{WithDevices: true})
	if err != nil {
		return
	}
	for _, net := range networks {
		err := inst.bacnetStoreNetwork(net)
		if err != nil {
			inst.bacnetErrorMsg("init network error:", err)
			continue
		}
		for _, dev := range net.Devices {
			err := inst.bacnetStoreDevice(dev)
			if err != nil {
				inst.bacnetErrorMsg("init device error:", err)
				continue
			}
		}
	}
}

func (inst *Instance) initBacStore() {
	if inst.BacStore == nil {
		inst.BacStore = network.NewStore()
		inst.bacnetNetworkInit()
	} else {
		inst.bacnetNetworkInit()
	}
}

// bacnetStoreNetwork add or update an instance a bacnet network that is cached in bacnet lib
func (inst *Instance) bacnetStoreNetwork(net *model.Network) error {
	bacnetNet := &network.Network{
		Interface: net.NetworkInterface,
		Port:      integer.NonNil(net.Port),
		StoreID:   net.UUID,
	}
	return inst.BacStore.UpdateNetwork(net.UUID, bacnetNet)
}

// getBacnetStoreNetwork get an instance of a created bacnet network that is cached in bacnet lib
func (inst *Instance) getBacnetStoreNetwork(networkUUID string) (*network.Network, error) {
	return inst.BacStore.GetNetwork(networkUUID)
}

// closeBacnetStoreNetwork delete the instance of a created bacnet network that is cached in bacnet lib
func (inst *Instance) closeBacnetStoreNetwork(networkUUID string) (bool, error) {
	net, err := inst.BacStore.GetNetwork(networkUUID)
	if err != nil {
		return false, err
	}
	net.NetworkClose(false)
	return true, nil
}

// getBacnetDevice get an instance of a created bacnet device that is cached in bacnet lib
func (inst *Instance) getBacnetStoreDevice(deviceUUID string) (*network.Device, error) {
	return inst.BacStore.GetDevice(deviceUUID)
}

// bacnetDevice add or update an instance of a created bacnet device that is cached in bacnet lib
func (inst *Instance) bacnetStoreDevice(dev *model.Device) error {
	max := intToUint32(integer.NonNil(dev.MaxADPU))
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

	net, _ := inst.getBacnetStoreNetwork(dev.NetworkUUID)
	if net == nil {
		getNetwork, err := inst.db.GetNetwork(dev.NetworkUUID, api.Args{})
		if getNetwork == nil {
			return errors.New("failed to find network to init bacnet network")
		}
		err = inst.bacnetStoreNetwork(getNetwork)
		if err != nil {
			return errors.New("network can not be empty")
		}

	}
	return inst.BacStore.UpdateDevice(dev.UUID, net, d)
}

func (inst *Instance) doReadPriority(pnt *model.Point, networkUUID, deviceUUID string) (pri *priority.Float32, err error) {
	object, _, _ := setObjectType(pnt.ObjectType)
	bp := &network.Point{
		ObjectID:   btypes.ObjectInstance(integer.NonNil(pnt.ObjectId)),
		ObjectType: object,
	}
	// get network
	net, err := inst.getBacnetStoreNetwork(networkUUID)
	if err != nil {
		return nil, err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetStoreDevice(deviceUUID)
	if err != nil {
		return nil, err
	}
	return dev.PointReadPriority(bp)
}

// getDev get an instance of a created bacnet device that is cached in bacnet lib
func (inst *Instance) doReadValue(pnt *model.Point, networkUUID, deviceUUID string) (float64, error) {
	object, _, isBool := setObjectType(pnt.ObjectType)
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
	net, err := inst.getBacnetStoreNetwork(networkUUID)
	if err != nil {
		return 0, err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetStoreDevice(deviceUUID)
	if err != nil {
		return 0, err
	}
	var outValue float64
	if isBool {
		readBool, err := dev.PointReadBool(bp)
		if err != nil {
			inst.bacnetErrorMsg(" read-bool:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " error:", err)
			return 0, err
		}
		outValue = Unit32ToFloat64(readBool)

	} else {
		readFloat32, err := dev.PointReadFloat32(bp)
		if err != nil {
			inst.bacnetErrorMsg(" read-analogue:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " error:", err)
			return 0, err
		}
		outValue = float32ToFloat64(readFloat32)
	}
	inst.bacnetDebugMsg(" POINT-READ:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", outValue)
	return outValue, nil
}

func (inst *Instance) doWrite(pnt *model.Point, networkUUID, deviceUUID string, writeValue float64, priority uint8) error {
	object, isWrite, isBool := setObjectType(pnt.ObjectType)
	if !isWrite || priority <= 0 || priority > 16 {
		priority = 16
	}
	bp := &network.Point{
		ObjectID:         btypes.ObjectInstance(integer.NonNil(pnt.ObjectId)),
		ObjectType:       object,
		WriteNull:        false,
		WritePriority:    priority,
		ReadPresentValue: false,
		ReadPriority:     false,
	}
	net, err := inst.getBacnetStoreNetwork(networkUUID)
	if err != nil {
		return err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetStoreDevice(deviceUUID)
	if err != nil {
		return err
	}
	if isBool {
		err = dev.PointWriteBool(bp, float64ToUint32(writeValue))
		if err != nil {
			inst.bacnetErrorMsg("write-bool:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", writeValue, " writePriority", priority, " error:", err)
			return err
		}
	} else {
		err = dev.PointWriteAnalogue(bp, float64ToFloat32(writeValue))
		if err != nil {
			inst.bacnetErrorMsg("write-analog:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", writeValue, " writePriority", priority, " error:", err)
			return err
		}
	}
	inst.bacnetDebugMsg("POINT-WRITE:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", writeValue, " writePriority", priority)
	return nil
}

func (inst *Instance) doRelease(pnt *model.Point, networkUUID, deviceUUID string, priority uint8) error {
	object, _, _ := setObjectType(pnt.ObjectType)
	if priority <= 0 || priority > 16 {
		return errors.New("invalid priority to doRelease()")
	}
	bp := &network.Point{
		ObjectID:         btypes.ObjectInstance(integer.NonNil(pnt.ObjectId)),
		ObjectType:       object,
		WriteNull:        false,
		WritePriority:    priority,
		ReadPresentValue: false,
		ReadPriority:     false,
	}
	net, err := inst.getBacnetStoreNetwork(networkUUID)
	if err != nil {
		return err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetStoreDevice(deviceUUID)
	if err != nil {
		return err
	}
	err = dev.PointReleasePriority(bp, priority)
	if err != nil {
		inst.bacnetErrorMsg("release-priority:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " priority:", priority, " error:", err)
		return err
	}
	inst.bacnetDebugMsg("POINT-RELEASE:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " priority:", priority)
	return nil
}

func (inst *Instance) pingDevice(network *model.Network, device *model.Device) error {
	if inst.BacStore != nil {
		netUUID := network.UUID
		devUUID := device.UUID
		storeNetwork, err := inst.getBacnetStoreNetwork(netUUID)
		if err != nil {
			return err
		}
		go storeNetwork.NetworkRun()
		storeDevice, err := inst.getBacnetStoreDevice(devUUID)
		if err != nil {
			return err
		}
		obj := integer.NonNil(device.DeviceObjectId)
		deviceName, err := storeDevice.ReadDeviceName(btypes.ObjectInstance(obj))
		if err != nil {
			return err
		}
		log.Infof("bacnet-server ping:%s", deviceName)
	} else {
		log.Infof("bacnet-server ping:%s device-id:%d", device.Name, integer.NonNil(device.DeviceObjectId))
	}
	return nil
}

func (inst *Instance) massUpdateServer(network *model.Network, device *model.Device) error {
	if inst.BacStore != nil {
		netUUID := network.UUID
		devUUID := device.UUID
		for _, point := range device.Points {
			err := inst.writeBacnetPointName(point, point.Name, netUUID, devUUID)
			if err != nil {
			} else {
				log.Infof("bacnet-server mass update names:%s", point.Name)
			}
		}
	} else {
		log.Infof("bacnet-server ping:%s device-id:%d", device.Name, integer.NonNil(device.DeviceObjectId))
	}
	return nil
}

/*
writeBacnetPointName
examples from the lib
go run main.go read --interface=wlp0s20f3 --address=192.168.15.191 --device=2508 --objectID=1 --objectType=1 --property=77
go run main.go write --interface=wlp0s20f3 --address=192.168.15.191 --device=2508 --objectID=1 --objectType=1 --property=77 --priority=16 --value=testing
*/
func (inst *Instance) writeBacnetPointName(pnt *model.Point, name, networkUUID, deviceUUID string) error {
	object, _, _ := setObjectType(pnt.ObjectType)
	bp := &network.Point{
		ObjectID:   btypes.ObjectInstance(integer.NonNil(pnt.ObjectId)),
		ObjectType: object,
	}
	net, err := inst.getBacnetStoreNetwork(networkUUID)
	if err != nil {
		return err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetStoreDevice(deviceUUID)
	if err != nil {
		return err
	}
	err = dev.WritePointName(bp, name)
	if err != nil {
		inst.bacnetErrorMsg("bacnet-server-write-name:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " name:", name, " error:", err)
		return err
	}
	inst.bacnetDebugMsg("bacnet-server-write-name:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " name:", name)
	return nil
}

func setObjectType(object string) (obj btypes.ObjectType, isWritable, isBool bool) {
	switch object {
	case "analog_input":
		return btypes.AnalogInput, false, false
	case "analog_output":
		return btypes.AnalogOutput, true, false
	case "analog_value":
		return btypes.AnalogValue, true, false
	case "binary_input":
		return btypes.BinaryInput, false, true
	case "binary_output":
		return btypes.BinaryOutput, true, true
	case "binary_value":
		return btypes.BinaryValue, true, true
	case "multi_state_input":
		return btypes.MultiStateInput, false, false
	case "multi_state_output":
		return btypes.MultiStateOutput, true, false
	case "multi_state_value":
		return btypes.MultiStateValue, true, false
	default:
		return btypes.AnalogInput, false, false
	}
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

func (inst *Instance) doWriteBool(networkUUID, deviceUUID string, pnt *network.Point, value uint32) error {
	net, err := inst.getBacnetStoreNetwork(networkUUID)
	if err != nil {
		return err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetStoreDevice(deviceUUID)
	if err != nil {
		return err
	}

	err = dev.PointWriteBool(pnt, value)
	if err != nil {
		return err
	}
	return nil

}

func (inst *Instance) whoIs(networkUUID string, opts *bacnet.WhoIsOpts) (resp []btypes.Device, err error) {
	net, err := inst.getBacnetStoreNetwork(networkUUID)
	if err != nil {
		return nil, err
	}
	go net.NetworkRun()
	devices, err := net.Whois(opts)
	if err != nil {
		return nil, err
	}
	return devices, err
}

func (inst *Instance) devicePoints(deviceUUID string) (resp []*network.PointDetails, err error) {
	getNetwork, err := inst.db.GetNetworkByDeviceUUID(deviceUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	net, err := inst.getBacnetStoreNetwork(getNetwork.UUID)
	if err != nil {
		return nil, err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetStoreDevice(deviceUUID)
	if err != nil {
		return nil, err
	}

	resp, err = dev.GetDevicePoints(btypes.ObjectInstance(dev.DeviceID))
	if err != nil {
		return nil, err
	}
	return

}

func ConvertPriorityToMap(priority priority.Float32) map[string]*float32 {
	priorityMap := map[string]*float32{}
	priorityValue := reflect.ValueOf(priority)
	typeOfPriority := priorityValue.Type()
	for i := 0; i < priorityValue.NumField(); i++ {
		if priorityValue.Field(i).Type().Kind().String() == "ptr" {
			key := typeOfPriority.Field(i).Tag.Get("json")
			val := priorityValue.Field(i).Interface().(*float32)
			priorityMap[key] = val
		}
	}
	return priorityMap
}

func FFPointAndBACnetServerPointAreEqual(FFPointValPntr *float64, BACServPointValPntr *float32) bool {
	BACServValIsNil := float.IsNil32(BACServPointValPntr)
	BACServVal64 := float64(float.NonNil32(BACServPointValPntr))
	FFPointValIsNil := float.IsNil(FFPointValPntr)
	FFPointVal64 := float.NonNil(FFPointValPntr)

	if BACServValIsNil && FFPointValIsNil {
		return true
	}
	if (BACServValIsNil && !FFPointValIsNil) || (!BACServValIsNil && FFPointValIsNil) {
		return false
	}
	if BACServVal64 == FFPointVal64 {
		return true
	}
	return false
}

func intToUint32(value int) uint32 {
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

func float64ToUint32(value float64) uint32 {
	var y = uint32(value)
	return y
}

func float64ToFloat32(value float64) float32 {
	var y = float32(value)
	return y
}
