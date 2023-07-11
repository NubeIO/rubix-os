package localmqtt

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	log "github.com/sirupsen/logrus"
	"strings"
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
	localMqtt.Client.Publish(PointPublishTopic, localMqtt.QOS, localMqtt.Retain, string(payload))
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
		topic = MakeTopic([]string{PointsPublishTopic})
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
	topic := MakeTopic([]string{PointValueTopic, CovTopic, AllTopic, network.PluginPath, network.UUID, networkName,
		device.UUID, deviceName, point.UUID, pointName})
	payload, err := json.Marshal(pointCovPayload)
	if err != nil {
		log.Error(err)
		return
	}
	localMqtt.Client.Publish(topic, localMqtt.QOS, localMqtt.Retain, string(payload))
}
