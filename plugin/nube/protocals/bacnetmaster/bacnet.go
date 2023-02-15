package main

import (
	"errors"
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/btypes/priority"
	"github.com/NubeDev/bacnet/btypes/segmentation"
	"github.com/NubeDev/bacnet/network"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/flow-framework/utils/priorityarray"
	"github.com/NubeIO/flow-framework/utils/writemode"
	address "github.com/NubeIO/lib-networking/ip"
	"github.com/NubeIO/lib-networking/networking"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func (inst *Instance) bacnetNetworkInit() {
	log.Infof("bacnet-master bacnetNetworkInit enable network plg-uuid: %s", inst.pluginUUID)
	networks, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{WithDevices: true})
	if err != nil {
		log.Errorln("bacnet-master bacnetNetworkInit err:", err.Error())
		return
	}
	for _, net := range networks {
		err := inst.makeBacnetStoreNetwork(net)
		if err != nil {
			log.Errorf("bacnet-master init network error: %s name: %s", err.Error(), net.Name)
			continue
		} else {
			log.Infof("bacnet-master init network: %s uuid: %s", net.Name, net.UUID)
		}
		for _, dev := range net.Devices {
			err := inst.bacnetStoreDevice(dev)
			if err != nil {
				log.Errorf("bacnet-master init device error: %s name: %s", err.Error(), dev.Name)
				continue
			} else {
				log.Infof("bacnet-master init device: %s uuid: %s", dev.Name, dev.UUID)
			}
		}
	}
}

func (inst *Instance) initBacStore() {
	if inst.BacStore == nil {
		log.Errorf("bacnet:master initBacStore bacnet store was empty")
		inst.BacStore = network.NewStore()
		inst.bacnetNetworkInit()
	} else {
		inst.bacnetNetworkInit()
	}
}

// makeBacnetStoreNetwork add or update an instance a bacnet network that is cached in bacnet lib
func (inst *Instance) makeBacnetStoreNetwork(net *model.Network) error {
	return inst.BacStore.NewNetwork(net.UUID, net.NetworkInterface, net.IP, integer.NonNil(net.Port), 0)
}

// getBacnetNetwork get an instance of a created bacnet network that is cached in bacnet lib
func (inst *Instance) getBacnetStoreNetwork(networkUUID string) (*network.Network, error) {
	net, err := inst.BacStore.GetNetwork(networkUUID)
	if err != nil {
		log.Errorf("bacnet-master: network getBacnetStoreNetwork() err: %s", err.Error())
	}
	if net == nil {
		log.Errorf("bacnet-master: network getBacnetStoreNetwork() network was empty")
	}
	return net, err
}

// closeBacnetNetwork delete the instance of a created bacnet network that is cached in bacnet lib
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
	if dev == nil {
		return errors.New("device can not be empty")
	}
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
		if err != nil {
			errMes := fmt.Sprintf("bacnet-master get network: %s name: %s", err.Error(), dev.Name)
			return errors.New(errMes)
		}
		err = inst.makeBacnetStoreNetwork(getNetwork)
		if err != nil {
			errMes := fmt.Sprintf("bacnet-master init device add/update bacnet error: %s name: %s", err.Error(), dev.Name)
			return errors.New(errMes)
		}
	}
	return inst.BacStore.UpdateDevice(dev.UUID, net, d)
}

func getNetworkIP(network string) (*networking.NetworkInterfaces, error) {
	net, err := networking.New().GetNetworkByIface(network)
	if err != nil {
		return nil, err
	}
	return &net, nil
}

func (inst *Instance) buildBacnetForAction(networkUUID, action string) (*network.Network, error) {
	// get network
	net, err := inst.getBacnetStoreNetwork(networkUUID)
	if err != nil {
		return nil, err
	}
	if net.Ip == "" { // TODO make this update the actual bacnet inst
		ip, err := getNetworkIP(net.Interface)
		if err != nil {
		} else {
			net.Ip = ip.IP
		}
	}
	err = address.New().IsIPAddrErr(net.Ip)
	if err != nil {
		log.Errorf("bacnet-master-%s: ip-address is invalid ip:%s err:%s", net.Ip, err.Error(), action)
		return nil, err
	}
	return net, err

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
	net, err := inst.buildBacnetForAction(networkUUID, "read")
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
			log.Errorln("bacnet-master-read-bool:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " error:", err)
			return 0, err
		}
		outValue = Unit32ToFloat64(readBool)

	} else if pnt.ObjectType == "multi_state_input" || pnt.ObjectType == "multi_state_output" || pnt.ObjectType == "multi_state_value" {
		readFloat32, err := dev.PointReadMultiState(bp)
		if err != nil {
			log.Errorln("bacnet-master-read-multi-state:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " error:", err)
			return 0, err
		}
		outValue = Unit32ToFloat64(readFloat32)
	} else {
		readFloat32, err := dev.PointReadFloat32(bp)
		if err != nil {
			log.Errorln("bacnet-master-read-analogue:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " error:", err)
			return 0, err
		}
		outValue = float32ToFloat64(readFloat32)
	}
	inst.bacnetDebugMsg("bacnet-master-POINT-READ:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", outValue)
	return outValue, nil
}

func (inst *Instance) doReadPriority(pnt *model.Point, networkUUID, deviceUUID string) (pri *priority.Float32, err error) {
	object, _, _, _ := setObjectType(pnt.ObjectType)
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

func (inst *Instance) doRelease(pnt *model.Point, networkUUID, deviceUUID string, priority uint8) error {
	object, _, _, _ := setObjectType(pnt.ObjectType)
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

func (inst *Instance) doWrite(pnt *model.Point, networkUUID, deviceUUID string) error {
	if pnt.WriteValue == nil {
		return errors.New("bacnet-write: point has no WriteValue")
	}
	val := float.NonNil(pnt.WriteValue)
	object, isWrite, isBool, _ := setObjectType(pnt.ObjectType)
	writePriority := integer.NonNil(pnt.WritePriority)
	if writePriority <= 0 || writePriority < 16 {
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
	// get network
	net, err := inst.buildBacnetForAction(networkUUID, "read")
	if err != nil {
		return err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetStoreDevice(deviceUUID)
	if err != nil {
		return errors.New("bacnet-write: error getting BACnet device details")
	}

	if isWrite {
		if isBool {
			err = dev.PointWriteBool(bp, float64ToUint32(val))
			if err != nil {
				inst.bacnetErrorMsg("bacnet-master-write-bool:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority, " error:", err)
				return err
			}
		} else if pnt.ObjectType == "multi_state_output" || pnt.ObjectType == "multi_state_value" {
			err = dev.PointWriteMultiState(bp, float64ToUint32(val))
			if err != nil {
				inst.bacnetErrorMsg("bacnet-master-write-multi-state:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority, " error:", err)
				return err
			}
		} else {
			err = dev.PointWriteAnalogue(bp, float64ToFloat32(val))
			if err != nil {
				inst.bacnetErrorMsg("bacnet-master-write-analogue:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority, " error:", err)
				return err
			}
		}
	}
	log.Infoln("bacnet-master-POINT-WRITE:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority)
	return nil
}

func (inst *Instance) doWriteAtPriority(pnt *model.Point, networkUUID, deviceUUID string, val *float64, writePriority int) error {
	if writePriority <= 0 || writePriority > 16 {
		writePriority = 16
	}
	object, isWrite, isBool, _ := setObjectType(pnt.ObjectType)
	bp := &network.Point{
		ObjectID:         btypes.ObjectInstance(integer.NonNil(pnt.ObjectId)),
		ObjectType:       object,
		WriteNull:        false,
		WritePriority:    uint8(writePriority),
		ReadPresentValue: false,
		ReadPriority:     false,
	}
	// get network
	net, err := inst.buildBacnetForAction(networkUUID, "read")
	if err != nil {
		return err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetStoreDevice(deviceUUID)
	if err != nil {
		return errors.New("bacnet-write: error getting BACnet device details")
	}

	if isWrite {
		if val == nil {
			err = dev.PointReleasePriority(bp, uint8(writePriority))
			if err != nil {
				inst.bacnetErrorMsg("bacnet-master-write-release-priority:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority, " error:", err)
				return err
			}
			log.Infoln("bacnet-master-point-release-priority:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority)
		} else {
			if isBool {
				err = dev.PointWriteBool(bp, float64ToUint32(*val))
				if err != nil {
					inst.bacnetErrorMsg("bacnet-master-write-bool:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", *val, " writePriority", writePriority, " error:", err)
					return err
				}
			} else if pnt.ObjectType == "multi_state_output" || pnt.ObjectType == "multi_state_value" {
				err = dev.PointWriteMultiState(bp, float64ToUint32(*val))
				if err != nil {
					inst.bacnetErrorMsg("bacnet-master-write-multi-state:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", *val, " writePriority", writePriority, " error:", err)
					return err
				}
			} else {
				err = dev.PointWriteAnalogue(bp, float64ToFloat32(*val))
				if err != nil {
					inst.bacnetErrorMsg("bacnet-master-write-analogue:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", *val, " writePriority", writePriority, " error:", err)
					return err
				}
			}
			log.Infoln("bacnet-master-point-write:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", *val, " writePriority", writePriority)
		}
	}
	log.Infoln("bacnet-master-POINT-WRITE:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority)
	return nil
}

func (inst *Instance) doReadAllThenWriteDiff7141516(pnt *model.Point, networkUUID, deviceUUID string) (currentBACServPriority *priority.Float32, highestPriorityValue *float64, readSuccess, writeSuccess bool, err error) {
	object, isWrite, isBool, _ := setObjectType(pnt.ObjectType)

	// Get Priority Array of BACnet Point
	currentBACServPriority, err = inst.doReadPriority(pnt, networkUUID, deviceUUID)
	if err != nil {
		inst.bacnetDebugMsg("BACnetMasterPolling(): doReadAllThenWriteDiff141516() doReadPriority() error:", err)
		if strings.Contains(err.Error(), "receive timed out") {
			return nil, nil, false, false, nil
		}
		inst.bacnetErrorMsg("BACnetMasterPolling(): failed to read BACnet priority array")
		// If the priority array read fails, then write only in16 (if it's a writable point), then just read the present value.
		if pnt.Priority != nil && writemode.IsWriteable(pnt.WriteMode) && isWrite {
			err = inst.doWriteAtPriority(pnt, networkUUID, deviceUUID, pnt.Priority.P16, 16)
			if err != nil {
				inst.bacnetErrorMsg("BACnetMasterPolling(): doReadAllThenWriteDiff141516() priority 16 write error:", err)
			} else {
				inst.bacnetDebugMsg("BACnetMasterPolling(): successfully wrote to priority 16")
				writeSuccess = true
			}
		}
		readValue, err := inst.doReadValue(pnt, networkUUID, deviceUUID)
		if err != nil {
			inst.bacnetErrorMsg("BACnetMasterPolling(): doReadAllThenWriteDiff141516() doReadValue() error:", err)
		} else {
			inst.bacnetDebugMsg("Successfully read bacnet point present value")
			readSuccess = true
		}
		highestPriorityValue = float.New(readValue)
		return nil, highestPriorityValue, readSuccess, writeSuccess, nil
	}
	readSuccess = true

	// get the highest priority value from the read priority array
	currentValueF32 := currentBACServPriority.HighestFloat32()
	if currentValueF32 == nil {
		highestPriorityValue = nil
	} else {
		highestPriorityValue = float.New(float64(*currentValueF32))
	}

	if writemode.IsWriteable(pnt.WriteMode) && isWrite {
		net, err := inst.buildBacnetForAction(networkUUID, "write")
		if err != nil {
			return nil, nil, readSuccess, writeSuccess, err
		}
		go net.NetworkRun()
		dev, err := inst.getBacnetStoreDevice(deviceUUID)
		if err != nil {
			inst.bacnetErrorMsg(fmt.Sprintf("get bacnet device err: %s", err.Error()))
			return nil, nil, readSuccess, writeSuccess, errors.New("bacnet-write: error getting BACnet device details")
		}
		priorityArray := pnt.Priority
		if priorityArray == nil {
			return nil, nil, readSuccess, writeSuccess, errors.New("bacnet-write: point has no PriorityArray")
		}
		priorityMap := priorityarray.ConvertToMap(*pnt.Priority)
		var writePriority int
		var bp *network.Point

		for key, val := range priorityMap {
			writePriority, _ = strconv.Atoi(regexp.MustCompile("[0-9]+").FindString(key))
			if writePriority != 7 && writePriority != 14 && writePriority != 15 && writePriority != 16 { // We only write to priorities 7, 14, 15, and 16
				continue
			}

			var currentBACnetServPriorityVal *float64
			var pointerToBACServPriority **float32
			switch writePriority {
			case 7:
				pointerToBACServPriority = &currentBACServPriority.P7
				if currentBACServPriority.P7 == nil {
					currentBACnetServPriorityVal = nil
				} else {
					currentBACnetServPriorityVal = float.New(float64(*currentBACServPriority.P7))
				}
			case 14:
				pointerToBACServPriority = &currentBACServPriority.P14
				if currentBACServPriority.P14 == nil {
					currentBACnetServPriorityVal = nil
				} else {
					currentBACnetServPriorityVal = float.New(float64(*currentBACServPriority.P14))
				}
			case 15:
				pointerToBACServPriority = &currentBACServPriority.P15
				if currentBACServPriority.P15 == nil {
					currentBACnetServPriorityVal = nil
				} else {
					currentBACnetServPriorityVal = float.New(float64(*currentBACServPriority.P15))
				}
			case 16:
				pointerToBACServPriority = &currentBACServPriority.P16
				if currentBACServPriority.P16 == nil {
					currentBACnetServPriorityVal = nil
				} else {
					currentBACnetServPriorityVal = float.New(float64(*currentBACServPriority.P16))
				}

			}
			if val == nil && currentBACnetServPriorityVal == nil || (val != nil && currentBACnetServPriorityVal != nil && *val == *currentBACnetServPriorityVal) {
				writeSuccess = true // write not required
				continue
			}

			bp = &network.Point{
				ObjectID:         btypes.ObjectInstance(integer.NonNil(pnt.ObjectId)),
				ObjectType:       object,
				WriteNull:        false,
				WritePriority:    uint8(writePriority),
				ReadPresentValue: false,
				ReadPriority:     false,
			}

			if val == nil {
				err = dev.PointReleasePriority(bp, uint8(writePriority))
				if err != nil {
					inst.bacnetErrorMsg("bacnet-master-write-release-priority:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority, " error:", err)
					continue
				}
				*pointerToBACServPriority = nil
				inst.bacnetDebugMsg("bacnet-master-point-release-priority:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority)
			} else {
				if isBool {
					err = dev.PointWriteBool(bp, float64ToUint32(*val))
					if err != nil {
						inst.bacnetErrorMsg("bacnet-master-write-bool:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", *val, " writePriority", writePriority, " error:", err)
						continue
					}
					writeSuccess = true
				} else if pnt.ObjectType == "multi_state_output" || pnt.ObjectType == "multi_state_value" {
					err = dev.PointWriteMultiState(bp, float64ToUint32(*val))
					if err != nil {
						inst.bacnetErrorMsg("bacnet-master-write-multi-state:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", *val, " writePriority", writePriority, " error:", err)
						continue
					}
					writeSuccess = true
				} else {
					err = dev.PointWriteAnalogue(bp, float64ToFloat32(*val))
					if err != nil {
						inst.bacnetErrorMsg("bacnet-master-write-analogue:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", *val, " writePriority", writePriority, " error:", err)
						continue
					}
					writeSuccess = true
				}
				inst.bacnetDebugMsg("bacnet-master-point-write:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", *val, " writePriority", writePriority)
				newF32 := float32(*val)
				*pointerToBACServPriority = &newF32
			}
		}
	}
	return currentBACServPriority, highestPriorityValue, readSuccess, writeSuccess, nil
}

func (inst *Instance) doWriteAllValues(pnt *model.Point, networkUUID, deviceUUID string) (highestPriorityValue *float64, err error) {
	object, isWrite, isBool, _ := setObjectType(pnt.ObjectType)
	if !isWrite {
		return nil, errors.New("bacnet-write: point is not a writeable type (AI/BI)")
	}

	// get network // TODO: Maybe this has to be run for each write?
	net, err := inst.buildBacnetForAction(networkUUID, "write")
	if err != nil {
		return nil, err
	}
	go net.NetworkRun()
	dev, err := inst.getBacnetStoreDevice(deviceUUID)
	if err != nil {
		log.Error(fmt.Sprintf("get bacnet device err: %s", err.Error()))
		return nil, errors.New("bacnet-write: error getting BACnet device details")
	}
	priorityArray := pnt.Priority
	if priorityArray == nil {
		return nil, errors.New("bacnet-write: point has no PriorityArray")
	}
	priorityMap := priorityarray.ConvertToMap(*pnt.Priority)
	var writePriority int
	var bp *network.Point
	highestPriorityWithValue := 16

	for key, val := range priorityMap {
		writePriority, _ = strconv.Atoi(regexp.MustCompile("[0-9]+").FindString(key))
		if writePriority <= 0 || writePriority > 16 {
			writePriority = 16
		}
		bp = &network.Point{
			ObjectID:         btypes.ObjectInstance(integer.NonNil(pnt.ObjectId)),
			ObjectType:       object,
			WriteNull:        false,
			WritePriority:    uint8(writePriority),
			ReadPresentValue: false,
			ReadPriority:     false,
		}

		if val == nil {
			err = dev.PointReleasePriority(bp, uint8(writePriority))
			if err != nil {
				inst.bacnetErrorMsg("bacnet-master-write-release-priority:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority, " error:", err)
				continue
			}
			log.Infoln("bacnet-master-point-release-priority:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", val, " writePriority", writePriority)
		} else {
			if isBool {
				err = dev.PointWriteBool(bp, float64ToUint32(*val))
				if err != nil {
					inst.bacnetErrorMsg("bacnet-master-write-bool:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", *val, " writePriority", writePriority, " error:", err)
					continue
				}
			} else if pnt.ObjectType == "multi_state_output" || pnt.ObjectType == "multi_state_value" {
				err = dev.PointWriteMultiState(bp, float64ToUint32(*val))
				if err != nil {
					inst.bacnetErrorMsg("bacnet-master-write-multi-state:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", *val, " writePriority", writePriority, " error:", err)
					continue
				}
			} else {
				err = dev.PointWriteAnalogue(bp, float64ToFloat32(*val))
				if err != nil {
					inst.bacnetErrorMsg("bacnet-master-write-analogue:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", *val, " writePriority", writePriority, " error:", err)
					continue
				}
			}
			log.Infoln("bacnet-master-point-write:", "type:", pnt.ObjectType, "id", integer.NonNil(pnt.ObjectId), " value:", *val, " writePriority", writePriority)
			if writePriority < highestPriorityWithValue {
				highestPriorityWithValue = writePriority
				highestPriorityValue = val
			}
		}
	}
	return highestPriorityValue, nil
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
		return btypes.AnalogValue, true, false, "analog_value"
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

func (inst *Instance) whoIs(networkUUID string, opts *bacnet.WhoIsOpts, addDevices bool) (resp []*model.Device, err error) {
	// get network
	net, err := inst.buildBacnetForAction(networkUUID, "read")
	if err != nil {
		return nil, err
	}
	devices, err := net.Whois(opts)
	if err != nil {
		return nil, err
	}
	var devicesList []*model.Device

	for _, device := range devices {
		newDevice := &model.Device{
			CommonUUID: model.CommonUUID{
				UUID: uuid.SmallUUID(),
			},
			CommonEnable: model.CommonEnable{
				Enable: boolean.NewTrue(),
			},
			Name: fmt.Sprintf("deviceId_%d_networkNum_%d", device.DeviceID, device.NetworkNumber),
			CommonDevice: model.CommonDevice{
				CommonIP: model.CommonIP{
					Host: device.Ip,
					Port: device.Port,
				},
				Manufacture: strconv.Itoa(int(device.Vendor)),
			},

			DeviceMac:      integer.New(device.MacMSTP),
			DeviceObjectId: integer.New(device.DeviceID),
			NetworkNumber:  integer.New(device.NetworkNumber),
			MaxADPU:        integer.New(int(device.MaxApdu)),
			Segmentation:   string(convertSegmentation(segmentation.SegmentedType(device.Segmentation))),
			NetworkUUID:    networkUUID,
		}
		if addDevices {
			addDevice, err := inst.addDevice(newDevice)
			if err != nil {
				log.Errorf("failed to add a new device from whois %d", device.ID.Instance)
				// return nil, err
			}
			log.Infof("added new device from whois %s", addDevice.Name)
		}
		devicesList = append(devicesList, newDevice)
	}
	return devicesList, nil
}

func (inst *Instance) devicePoints(deviceUUID string, addPoints, writeablePoints bool) (resp []*model.Point, err error) {
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

	bacnetPoints, err := dev.GetDevicePoints(btypes.ObjectInstance(dev.DeviceID))
	if err != nil {
		return nil, err
	}
	var pointsList []*model.Point
	for _, pnt := range bacnetPoints {
		_, isWrite, _, objectType := setObjectType(convertObjectType(pnt.ObjectType))
		writeMode := model.ReadOnly
		if isWrite && writeablePoints {
			writeMode = model.WriteOnceThenRead
		}
		newPnt := &model.Point{
			CommonUUID: model.CommonUUID{
				UUID: uuid.SmallUUID(),
			},
			CommonEnable: model.CommonEnable{
				Enable: boolean.NewTrue(),
			},
			ScaleEnable: boolean.NewFalse(),
			Name:        pnt.Name,
			DeviceUUID:  deviceUUID,
			ObjectType:  objectType,
			ObjectId:    integer.New(int(pnt.ObjectID)),
			WriteMode:   writeMode,
		}
		if addPoints {
			point, err := inst.addPoint(newPnt)
			if err != nil {
				log.Errorf("failed to add a new point from discover points %s", point.Name)
				// return nil, err
			} else {
				log.Infof("added new point from discover points%s", point.Name)
			}
		}
		pointsList = append(pointsList, newPnt)
	}
	return pointsList, nil
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

func ConvertPriorityToMap(priority priority.Float32) map[string]*float64 {
	priorityMap := map[string]*float64{}
	priorityValue := reflect.ValueOf(priority)
	typeOfPriority := priorityValue.Type()
	for i := 0; i < priorityValue.NumField(); i++ {
		if priorityValue.Field(i).Type().Kind().String() == "ptr" {
			key := typeOfPriority.Field(i).Tag.Get("json")
			val := priorityValue.Field(i).Interface().(*float32)
			var val64 *float64
			if val == nil {
				val64 = nil
			} else {
				val64 = float.New(float64(*val))
			}
			priorityMap[key] = val64
		}
	}
	return priorityMap
}
