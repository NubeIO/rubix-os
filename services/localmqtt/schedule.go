package localmqtt

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

func PublishSchedule(schedule *model.Schedule) {
	if schedule == nil {
		return
	}
	payload, err := json.Marshal(schedule)
	if err != nil {
		log.Error(err)
		return
	}
	topic := fmt.Sprintf("rubix/platform/schedule/publish")
	pointMqtt.Client.Publish(topic, pointMqtt.QOS, pointMqtt.Retain, string(payload))
}

func PublishSchedules(schedules []*model.Schedule, topic string) {
	if topic == "" {
		topic = MakeTopic([]string{fetchSchedulesTopic})
	}
	payload, err := json.Marshal(schedules)
	if err != nil {
		log.Error(err)
		return
	}
	pointMqtt.Client.Publish(topic, pointMqtt.QOS, pointMqtt.Retain, string(payload))
}
