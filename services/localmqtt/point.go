package localmqtt

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
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
	localMqtt.Client.Publish(topic, localMqtt.QOS, localMqtt.Retain, string(payload))
}

func PublishPointsList(publishPointList []*interfaces.PublishPointList, topic string) {
	if !localMqtt.PublishPointList {
		return
	}
	var pointPayload []*PointListPayload
	for _, publishPoint := range publishPointList {
		pointPayload = append(pointPayload, &PointListPayload{UUID: publishPoint.PointUUID,
			Name: fmt.Sprintf("%s:%s:%s:%s", publishPoint.PluginPath, publishPoint.NetworkName,
				publishPoint.DeviceName, publishPoint.PointName)})
	}
	if topic == "" {
		topic = MakeTopic([]string{fetchPointsTopic})
	}
	payload, err := json.Marshal(pointPayload)
	if err != nil {
		log.Error(err)
		return
	}
	localMqtt.Client.Publish(topic, localMqtt.QOS, localMqtt.Retain, string(payload))
}

func PublishPointCov(network *model.Network, device *model.Device, point *model.Point) {
	if !localMqtt.PublishPointCOV {
		return
	}
	pointCovPayload := &PointCovPayload{
		Value:    point.PresentValue,
		ValueRaw: point.OriginalValue,
		Priority: point.CurrentPriority,
		Ts:       point.UpdatedAt.String(),
	}
	networkName := strings.Trim(strings.Trim(network.Name, " "), "\t")
	deviceName := strings.Trim(strings.Trim(device.Name, " "), "\t")
	pointName := strings.Trim(strings.Trim(point.Name, " "), "\t")
	topic := MakeTopic([]string{mqttTopic, mqttTopicCov, mqttTopicCovAll, network.PluginPath, network.UUID, networkName,
		device.UUID, deviceName, point.UUID, pointName})
	payload, err := json.Marshal(pointCovPayload)
	if err != nil {
		log.Error(err)
		return
	}
	localMqtt.Client.Publish(topic, localMqtt.QOS, localMqtt.Retain, string(payload))
}

func ifEmpty(in string) string {
	if in == "" {
		return "na"
	}
	return in
}

func MakeTopic(parts []string) string {
	// TODO: if localMqtt.GlobalBroadcast -> use location/group/host uuid and name
	return strings.Join(append(parts), separator)
}
