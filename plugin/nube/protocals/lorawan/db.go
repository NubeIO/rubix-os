package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csmodel"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) csConvertDevice(dev *model.Device, csDev *csmodel.Device) {
	dev.NetworkUUID = inst.networkUUID
	dev.CommonName.Name = csDev.Name
	dev.CommonDescription.Description = csDev.Description
	dev.CommonDevice.AddressUUID = &csDev.DevEUI
}

func (inst *Instance) getNetwork() (network *model.Network, err error) {
	net, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{})
	if len(net) == 0 {
		return nil, err
	}
	return net[0], err
}

// addNetwork add network
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	if len(nets) > 0 {
		errMsg := "lorawan: only max one network is allowed with lora"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	body, err = inst.db.CreateNetwork(body, true)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// addPoint add point
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body.ObjectType == "" {
		errMsg := "lorawan: point object type can not be empty"
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
// func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
//     network, err = inst.db.UpdateNetwork(body.UUID, body, true)
//     if err != nil {
//         return nil, err
//     }
//     return network, nil
// }

// updateDevice update device
// func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
//     device, err = inst.db.UpdateDevice(body.UUID, body, true)
//     if err != nil {
//         return nil, err
//     }

//     return device, nil
// }

// updatePoint update point
// func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
//     point, err = inst.db.UpdatePoint(body.UUID, body, true)
//     if err != nil {
//         return nil, err
//     }
//     return point, nil
// }

// writePoint update point. Called via API call.
func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	// TODO: check for PointWriteByName calls that might not flow through the plugin.
	if body == nil {
		return
	}
	point, err = inst.db.WritePoint(pntUUID, body, true)
	if err != nil || point == nil {
		return nil, err
	}
	return point, nil
}

// pointUpdate update point present value
func (inst *Instance) pointUpdate(uuid string) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	point.InSync = boolean.NewTrue()
	_, err := inst.db.UpdatePoint(uuid, &point, true)
	if err != nil {
		log.Error("lorawan: UpdatePoint()", err)
		return nil, err
	}
	return nil, nil
}

// pointUpdate update point present value
func (inst *Instance) pointUpdateValue(uuid string, value float64) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	priority := map[string]*float64{"_16": &value}
	point.InSync = boolean.NewTrue()
	_, err := inst.db.UpdatePointValue(uuid, &point, &priority, true)
	if err != nil {
		log.Error("lorawan: pointUpdateValue()", err)
		return nil, err
	}
	return nil, nil
}

// pointUpdate update point present value
func (inst *Instance) pointUpdateErr(uuid string, err error) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteError
	point.CommonFault.Message = fmt.Sprintf("error-time: %s msg:%s", utilstime.TimeStamp(), err.Error())
	point.CommonFault.LastFail = time.Now().UTC()
	point.InSync = boolean.NewFalse()
	_, err = inst.db.UpdatePoint(uuid, &point, true)
	if err != nil {
		log.Error("lorawan: pointUpdateErr()", err)
		return nil, err
	}
	return nil, nil
}
