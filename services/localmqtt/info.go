package localmqtt

import (
	"encoding/json"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	log "github.com/sirupsen/logrus"
)

const fetchDeviceInfo = "rubix/platform/info/publish"

func PublishInfo() {
	deviceInfo, err := deviceinfo.GetDeviceInfo()
	if err != nil {
		log.Error(err)
	}
	marshal, err := json.Marshal(deviceInfo)
	if err != nil {
		return
	}
	err = pointMqtt.Client.Publish(fetchDeviceInfo, pointMqtt.QOS, retainMessage, string(marshal))
	if err != nil {
		log.Error(err)
	}
}
