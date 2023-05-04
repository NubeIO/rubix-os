package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/NubeDev/bacnet/btypes/priority"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"io"
)

func (inst *Instance) commandPV(body *readBody) error {
	cli, _ := mqttclient.GetMQTT()
	err := cli.PublishNonBuffer(string(topicCommandRead), mqttclient.AtMostOnce, false, buildPayload(body))
	return err
}

func buildPayload(payload interface{}) interface{} {
	p, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(p)
}

func (inst *Instance) pvCallBack() {

}

func shortUUID(prefix ...string) string {
	u := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, u)
	if n != len(u) || err != nil {
		return "-error-uuid-"
	}
	uuid := fmt.Sprintf("%x%x", u[0:4], u[4:4])
	if len(prefix) > 0 {
		uuid = fmt.Sprintf("%s_%s", prefix[0], uuid)
	}
	return uuid
}

func (inst *Instance) doRead(point *model.Point, deviceUUID, networkUUID string) (currentBACServPriority *priority.Float32, highestPriorityValue *float64, readSuccess, writeSuccess bool, err error) {

	// fmt.Println(111111, deviceUUID, networkUUID)
	// inst.mqttGetPV("1")

	currentBACServPriority = &priority.Float32{
		P1:  nil,
		P2:  nil,
		P3:  nil,
		P4:  nil,
		P5:  nil,
		P6:  nil,
		P7:  nil,
		P8:  nil,
		P9:  nil,
		P10: nil,
		P11: nil,
		P12: nil,
		P13: nil,
		P14: nil,
		P15: nil,
		P16: nil,
	}

	return nil, nil, false, false, err

}
