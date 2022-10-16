package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	separator       = "/"
	mqttTopic       = "rubix/points/value"
	mqttTopicCov    = "cov"
	mqttTopicCovAll = "all"
)

var pointMqtt *PointMqtt

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
					Name: fmt.Sprintf("%s:%s:%s", network.Name, device.Name, point.Name)})
			}
		}
	}
	topic := makeTopic([]string{mqttTopic, "points"})
	payload, err := json.Marshal(pointPayload)
	if err != nil {
		log.Error(err)
		return
	}
	err = pointMqtt.client.Publish(topic, pointMqtt.QOS, false, string(payload))
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
	err = pointMqtt.client.Publish(topic, pointMqtt.QOS, false, string(payload))
	if err != nil {
		log.Error(err)
	}
}

func makeTopic(parts []string) string {
	deviceInfo, _ := deviceinfo.GetDeviceInfo()
	prefixTopic := []string{deviceInfo.ClientId, deviceInfo.ClientName, deviceInfo.SiteId, deviceInfo.SiteName,
		deviceInfo.DeviceId, deviceInfo.DeviceName}
	return strings.Join(append(prefixTopic, parts...), separator)
}
