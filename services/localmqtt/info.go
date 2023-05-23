package localmqtt

import (
	"encoding/json"
	"github.com/NubeIO/rubix-os/utils/deviceinfo"
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
	localMqtt.Client.Publish(fetchDeviceInfo, localMqtt.QOS, localMqtt.Retain, string(marshal))
}
