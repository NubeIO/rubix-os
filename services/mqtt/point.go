package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	separator = "/"
	mqttTopic = "rubix/points/value"
	// +/+/+/+/+/+/rubix/points/value/points
	// +/+/+/+/+/+/rubix/points/value/cov/all/system/+/net/+/dev/+/pnt
	mqttTopicCov    = "cov"
	mqttTopicCovAll = "all"
)

var pointMqtt *PointMqtt
var retainMessage bool

func Init(ip string, conf *config.Configuration) error {
	pm := new(PointMqtt)
	pm.QOS = mqttclient.QOS(conf.MQTT.QOS)
	c, err := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers:        []string{ip},
		Username:       conf.MQTT.Username,
		Password:       conf.MQTT.Password,
		SetKeepAlive:   conf.MQTT.SetKeepAlive,
		SetPingTimeout: conf.MQTT.SetPingTimeout,
		AutoReconnect:  conf.MQTT.AutoReconnect,
	})
	if err != nil {
		log.Println("MQTT connection error:", err)
		return err
	}
	pm.client = c
	pointMqtt = pm
	err = pm.client.Connect()
	if err != nil {
		return err
	}
	retainMessage = boolean.NonNil(conf.MQTT.Retain)
	return nil
}

func GetPointMqtt() *PointMqtt {
	return pointMqtt
}

func PublishPointsList(networks []*model.Network) {
	var pointPayload []*PointListPayload
	for _, network := range networks {
		for _, device := range network.Devices {
			for _, point := range device.Points {
				pointPayload = append(pointPayload, &PointListPayload{UUID: point.UUID,
					Name: fmt.Sprintf("%s:%s:%s:%s", network.PluginPath, network.Name, device.Name, point.Name)})
			}
		}
	}
	topic := makeTopic([]string{mqttTopic, "points"})
	payload, err := json.Marshal(pointPayload)
	if err != nil {
		log.Error(err)
		return
	}
	err = pointMqtt.client.Publish(topic, pointMqtt.QOS, retainMessage, string(payload))
	if err != nil {
		log.Error(err)
	}
}

func PublishPointCov(network *model.Network, device *model.Device, point *model.Point, priority *float64) {
	pointCovPayload := &PointCovPayload{
		Value:    point.PresentValue,
		ValueRaw: point.OriginalValue,
		Priority: priority,
		Ts:       point.UpdatedAt.String(),
	}
	topic := makeTopic([]string{mqttTopic, mqttTopicCov, mqttTopicCovAll, network.PluginPath, network.UUID,
		network.Name, device.UUID, device.Name, point.UUID, point.Name})
	payload, err := json.Marshal(pointCovPayload)
	if err != nil {
		log.Error(err)
		return
	}
	err = pointMqtt.client.Publish(topic, pointMqtt.QOS, retainMessage, string(payload))
	if err != nil {
		log.Error(err)
	}
}

func ifEmpty(in string) string {
	if in == "" {
		return "na"
	}
	return in
}

func makeTopic(parts []string) string {
	deviceInfo, _ := deviceinfo.GetDeviceInfo()
	clientId := deviceInfo.ClientId
	clientName := deviceInfo.ClientName
	siteId := deviceInfo.SiteId
	siteName := deviceInfo.SiteName
	deviceId := deviceInfo.DeviceId
	deviceName := deviceInfo.DeviceName
	prefixTopic := []string{ifEmpty(clientId), ifEmpty(clientName), ifEmpty(siteId), ifEmpty(siteName),
		ifEmpty(deviceId), ifEmpty(deviceName)}
	return strings.Join(append(prefixTopic, parts...), separator)
}
