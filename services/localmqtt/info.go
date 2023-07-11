package localmqtt

import (
	"encoding/json"
	"github.com/NubeIO/rubix-os/rubixregistry"
)

func PublishInfo(deviceInfo *rubixregistry.DeviceInfo) {
	marshal, err := json.Marshal(deviceInfo)
	if err != nil {
		return
	}
	localMqtt.Client.Publish(DeviceInfoPublishTopic, localMqtt.QOS, localMqtt.Retain, string(marshal))
}
