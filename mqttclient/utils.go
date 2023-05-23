package mqttclient

import (
	"github.com/NubeIO/rubix-os/utils/array"
	"strings"
)

func TopicParts(topic string) (clean, raw *array.Array) {
	s := strings.SplitAfter(topic, "/")
	clean = array.NewArray()
	raw = array.NewArray()
	for _, t := range s {
		if t == "/" {
			clean.Add("EMPTY-TOPIC-SPACE")
		}
		res := strings.ReplaceAll(t, "/", "")
		if res != "" {
			clean.Add(res)
		}
		raw.Add(t)
	}
	return clean, raw
}
