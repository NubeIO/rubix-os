package eventbus

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/mustafaturan/bus/v3"
	"strings"
)

//IsThisPlugin check if this is the correct plugin
func IsThisPlugin(topic, pluginUUID string) (this bool) {
	s := strings.Split(topic, ".")
	if len(s) > 2 {
		if strings.Contains(s[2], "plg") {
			if s[2] == pluginUUID {
				return true
			}
		}
	}
	return
}

//GetTopicPart will split the topics
func GetTopicPart(topic string, index int, contains string) string {
	s := strings.Split(topic, ".")
	for i, e := range s {
		if i == index {
			if strings.Contains(e, contains) { // if topic has pnt (is uuid of point)
				return e
			}
		}
	}
	return ""
}

//IsNetwork check if the payload is of type device
func IsNetwork(topic string, payload bus.Event) (*model.Network, error) {
	if GetTopicPart(topic, 3, "net") != "" {
		p, _ := payload.Data.(*model.Network)
		return p, nil
	}
	return nil, nil
}

//IsDevice check if the payload is of type device
func IsDevice(topic string, payload bus.Event) (*model.Device, error) {
	if GetTopicPart(topic, 3, "dev") != "" {
		p, _ := payload.Data.(*model.Device)
		return p, nil
	}
	return nil, nil
}

//IsPoint check if the payload is of type point
func IsPoint(topic string, payload bus.Event) (*model.Point, error) {
	if GetTopicPart(topic, 3, "pnt") != "" {
		p, _ := payload.Data.(*model.Point)
		return p, nil
	}
	return nil, nil
}

//IsSchedule check if the payload is of type schedule
func IsSchedule(topic string, payload bus.Event) (*model.Schedule, error) {
	if GetTopicPart(topic, 3, "sch") != "" {
		p, _ := payload.Data.(*model.Schedule)
		return p, nil
	}
	return nil, nil
}

//IsJob check if the payload is of type job
func IsJob(topic string, payload bus.Event) (*model.Job, error) {
	if GetTopicPart(topic, 3, "job") != "" {
		p, _ := payload.Data.(*model.Job)
		return p, nil
	}
	return nil, nil
}

// DecodeBody  update it
func DecodeBody(thingType string, payload interface{}) (interface{}, error) {
	switch thingType {
	case model.ThingClass.Point:
		p := payload.(*model.Point)
		return p, nil
	}
	return nil, nil
}
