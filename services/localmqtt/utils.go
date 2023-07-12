package localmqtt

import (
	"strings"
)

func ifEmpty(in string) string {
	if in == "" {
		return "na"
	}
	return in
}

func MakeTopic(parts []string) string {
	// TODO: if localMqtt.GlobalBroadcast -> use location/group/host uuid and name
	return strings.Join(append(parts), Separator)
}

func PublishValue(topic, payload string) {
	localMqtt.Client.Publish(topic, localMqtt.QOS, localMqtt.Retain, payload)
}
