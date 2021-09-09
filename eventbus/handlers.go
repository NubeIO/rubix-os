package eventbus

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/mustafaturan/bus/v3"
	"strings"
)

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