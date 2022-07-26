package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
)

// addNetwork add network
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	if body.NetworkInterface == "" {
		return nil, errors.New("network interface can not be empty try, eth0")
	}

	nets, err := inst.db.GetNetworksByPluginName(body.PluginPath, api.Args{})
	if err != nil {
		return nil, err
	}
	for _, net := range nets {
		if net != nil {
			errMsg := fmt.Sprintf("bacnet-server: only max one network is allowed with lora")
			log.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
	}
	body.Port = integer.New(defaultPort)
	network, err = inst.db.CreateNetwork(body, true)
	if err != nil {
		return nil, err
	}
	err = inst.bacnetNetwork(network)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("err:%s", err.Error()))
	}
	return body, nil
}

// addDevice add device
func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	network, err := inst.db.GetNetwork(body.NetworkUUID, api.Args{WithDevices: true})
	if err != nil {
		return nil, err
	}
	if len(network.Devices) >= 1 {
		errMsg := fmt.Sprintf("bacnet-server: only max one device is allowed")
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}

	body.NumberOfDevicesPermitted = integer.New(1)
	body.CommonIP.Host = inst.getIp(network.NetworkInterface)
	if integer.IsNil(body.DeviceObjectId) {
		body.DeviceObjectId = integer.New(2508)
	}
	body.Port = 47808
	device, err = inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}
	err = inst.bacnetDevice(device)
	if err != nil {
		return nil, errors.New("issue on add bacnet-device to store")
	}
	return device, nil
}

// addPoint add point
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body.ObjectType == "" {
		errMsg := fmt.Sprintf("bacnet-bserver: point object type can not be empty")
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	point, err = inst.db.CreatePoint(body, true, true)
	if err != nil {
		return nil, err
	}
	return point, nil
}

// updateNetwork update network
func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	network, err = inst.db.UpdateNetwork(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return network, nil
}

// updateDevice update device
func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	device, err = inst.db.UpdateDevice(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	err = inst.bacnetDevice(device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (inst *Instance) getNetworks() ([]*model.Network, error) {
	return inst.db.GetNetworks(api.Args{})
}

// deleteNetwork delete network
func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	ok, err = inst.closeBacnetNetwork(body.UUID)
	return ok, err
}

// writePoint update point. Called via API call.
func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	point, _, _, _, err = inst.db.PointWrite(pntUUID, body, false)
	return point, err
}

// deleteNetwork delete device
func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// deletePoint delete point
func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// pointWrite update point present value
func (inst *Instance) pointWrite(uuid string, value float64) error {
	priority := map[string]*float64{"_16": &value}
	pointWriter := model.PointWriter{Priority: &priority}
	_, _, _, _, err := inst.db.PointWrite(uuid, &pointWriter, true)
	if err != nil {
		log.Error("bacnet-server: pointWrite()", err)
	}
	return err
}

// pointUpdateSuccess update point present value
func (inst *Instance) pointUpdateSuccess(uuid string) error {
	var point model.Point
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	point.InSync = boolean.NewTrue()
	err := inst.db.UpdatePointErrors(uuid, &point)
	if err != nil {
		log.Error("bacnet-server: pointUpdateSuccess()", err)
	}
	return err
}

// pointUpdateErr update point present value
func (inst *Instance) pointUpdateErr(uuid string, err error) error {
	var point model.Point
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteError
	point.CommonFault.Message = fmt.Sprintf("error-time: %s msg:%s", utilstime.TimeStamp(), err.Error())
	point.CommonFault.LastFail = time.Now().UTC()
	point.InSync = boolean.NewFalse()
	err = inst.db.UpdatePointErrors(uuid, &point)
	if err != nil {
		log.Error("bacnet-server: pointUpdateErr()", err)
	}
	return err
}
