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
		getPnt, err := inst.db.GetOnePointByArgs(api.Args{AddressUUID: &_pUUID})
		if err != nil || getPnt.UUID == "" {
			log.Error("bacnet-master-plugin: ERROR on get bacnetUpdate() failed to find point", err, _pUUID)
			return
		}
		var pri model.Priority
		pri.P16 = payload.Value
		getPnt.Priority = &pri
		getPnt.CommonFault.InFault = false
		getPnt.CommonFault.MessageLevel = model.MessageLevel.Info
		getPnt.CommonFault.MessageCode = model.CommonFaultCode.Ok
		getPnt.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
		getPnt.CommonFault.LastOk = time.Now().UTC()
		_, err = inst.db.UpdatePointValue(getPnt.UUID, getPnt, true)
		if err != nil {
			log.Error("BACNET UPDATE POINT issue on message from mqtt update point")
			return
		}
		return
	}
}

//writePoint update point. Called via API call.
func (inst *Instance) writePoint(pntUUID string, body *model.Point) (point *model.Point, err error) {
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
