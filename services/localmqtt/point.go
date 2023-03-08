package localmqtt

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
	separator           = "/"
	mqttTopic           = "rubix/points/value"
	fetchPointsTopic    = "rubix/platform/points/publish"
	fetchSchedulesTopic = "rubix/platform/schedules/publish"
	mqttTopicCov        = "cov"
	mqttTopicCovAll     = "all"
)

var pointMqtt *PointMqtt
var retainMessage bool
var globalBroadcast bool

func Init(ip string, conf *config.Configuration, onConnected interface{}) error {
	pm := new(PointMqtt)
	pm.QOS = mqttclient.QOS(conf.MQTT.QOS)
	c, err := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers:              []string{ip},
		Username:             conf.MQTT.Username,
		Password:             conf.MQTT.Password,
		SetKeepAlive:         conf.MQTT.SetKeepAlive,
		SetPingTimeout:       conf.MQTT.SetPingTimeout,
		ConnectRetry:         conf.MQTT.ConnectRetry,
		ConnectRetryInterval: conf.MQTT.ConnectRetryInterval,
		AutoReconnect:        conf.MQTT.AutoReconnect,
		MaxReconnectInterval: conf.MQTT.MaxReconnectInterval,
	}, onConnected)
	if err != nil {
		log.Info("MQTT connection error:", err)
		return err
	}
	pm.Client = c
	pointMqtt = pm
	err = pm.Client.Connect()
	if err != nil {
		return err
	}
	retainMessage = boolean.NonNil(conf.MQTT.Retain)
	globalBroadcast = boolean.NonNil(conf.MQTT.GlobalBroadcast)
	return nil
}

func GetPointMqtt() *PointMqtt {
	return pointMqtt
}

func PublishPoint(point *model.Point) {
	if point == nil {
		return
	}
	payload, err := json.Marshal(point)
	if err != nil {
		log.Error(err)
		return
	}
	topic := fmt.Sprintf("rubix/platform/point/publish")
	pointMqtt.Client.Publish(topic, pointMqtt.QOS, retainMessage, string(payload))
}

func PublishPointsList(networks []*model.Network, topic string) {
	var pointPayload []*PointListPayload
	for _, network := range networks {
		for _, device := range network.Devices {
			for _, point := range device.Points {
				pointPayload = append(pointPayload, &PointListPayload{UUID: point.UUID,
					Name: fmt.Sprintf("%s:%s:%s:%s", network.PluginPath, network.Name, device.Name, point.Name)})
			}
		}
	}
	if topic == "" {
		topic = MakeTopic([]string{fetchPointsTopic})
	}
	payload, err := json.Marshal(pointPayload)
	if err != nil {
		log.Error(err)
		return
	}
	pointMqtt.Client.Publish(topic, pointMqtt.QOS, retainMessage, string(payload))
}

func PublishPointCov(network *model.Network, device *model.Device, point *model.Point) {
	pointCovPayload := &PointCovPayload{
		Value:    point.PresentValue,
		ValueRaw: point.OriginalValue,
		Priority: point.CurrentPriority,
		Ts:       point.UpdatedAt.String(),
	}
	topic := MakeTopic([]string{mqttTopic, mqttTopicCov, mqttTopicCovAll, network.PluginPath, network.UUID,
		network.Name, device.UUID, device.Name, point.UUID, point.Name})
	payload, err := json.Marshal(pointCovPayload)
	if err != nil {
		log.Error(err)
		return
	}
	pointMqtt.Client.Publish(topic, pointMqtt.QOS, retainMessage, string(payload))
}

func ifEmpty(in string) string {
	if in == "" {
		return "na"
	}
	return in
}

func MakeTopic(parts []string) string {
	deviceInfo, _ := deviceinfo.GetDeviceInfo()
	clientId := deviceInfo.ClientId
	clientName := deviceInfo.ClientName
	siteId := deviceInfo.SiteId
	siteName := deviceInfo.SiteName
	deviceId := deviceInfo.DeviceId
	deviceName := deviceInfo.DeviceName
	prefixTopic := []string{ifEmpty(clientId), ifEmpty(clientName), ifEmpty(siteId), ifEmpty(siteName),
		ifEmpty(deviceId), ifEmpty(deviceName)}

	if globalBroadcast {
		return strings.Join(append(prefixTopic, parts...), separator)
	}
	return strings.Join(append(parts), separator)

}
