package localmqtt

import (
	"encoding/json"
	"github.com/NubeIO/rubix-registry-go/rubixregistry"
)

const fetchDeviceInfo = "rubix/platform/info/publish"

func PublishInfo(deviceInfo *rubixregistry.DeviceInfo) {
	marshal, err := json.Marshal(deviceInfo)
	if err != nil {
		return
	}
	localMqtt.Client.Publish(fetchDeviceInfo, localMqtt.QOS, localMqtt.Retain, string(marshal))
}
