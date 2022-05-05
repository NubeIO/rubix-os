package main

import (
	"encoding/json"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetmaster/bmmodel"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"time"
)

//bacnetUpdate listen on mqtt and then update the point in flow-framework
func (inst *Instance) bacnetUpdate(body mqtt.Message) {
	payload := new(bmmodel.MqttPayload)
	err := json.Unmarshal(body.Payload(), &payload)
	if err != nil {
		return
	}
	t, _ := mqttclient.TopicParts(body.Topic())
	const pointUUID = 14
	if t.Size() >= pointUUID {
		pUUID := t.Get(pointUUID)
		_pUUID := pUUID.(string)
		point, err := inst.db.GetOnePointByArgs(api.Args{AddressUUID: &_pUUID})
		if err != nil || point.UUID == "" {
			log.Error("bacnet-master-plugin: ERROR on get bacnetUpdate() failed to find point", err, _pUUID)
			return
		}
		priority := map[string]*float64{"_16": payload.Value}
		point.CommonFault.InFault = false
		point.CommonFault.MessageLevel = model.MessageLevel.Info
		point.CommonFault.MessageCode = model.CommonFaultCode.Ok
		point.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
		point.CommonFault.LastOk = time.Now().UTC()
		_, err = inst.db.UpdatePointValue(point.UUID, point, &priority, true)
		if err != nil {
			log.Error("BACNET UPDATE POINT issue on message from mqtt update point")
			return
		}
		return
	}
}

//writePoint update point. Called via API call.
func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	//TODO: check for PointWriteByName calls that might not flow through the plugin.
	if body == nil {
		return
	}
	point, err = inst.db.WritePoint(pntUUID, body, true)
	if err != nil || point == nil {
		return nil, err
	}
	return point, nil
}
