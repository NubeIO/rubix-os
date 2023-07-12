package localmqtt

import (
	"encoding/json"
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
	localMqtt.Client.Publish(SchedulePublishTopic, localMqtt.QOS, localMqtt.Retain, string(payload))
}

func PublishSchedules(schedules []*model.Schedule, topic string) {
	if !localMqtt.PublishScheduleList {
		return
	}
	if topic == "" {
		topic = MakeTopic([]string{SchedulesPublishTopic})
	}
	payload, err := json.Marshal(schedules)
	if err != nil {
		log.Error(err)
		return
	}
	localMqtt.Client.Publish(topic, localMqtt.QOS, localMqtt.Retain, string(payload))
}
