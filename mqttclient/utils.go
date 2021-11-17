package mqttclient

import (
	"github.com/NubeIO/flow-framework/utils"
	"strings"
)

func TopicParts(topic string) (clean, raw *utils.Array) {
	s := strings.SplitAfter(topic, "/")
	clean = utils.NewArray()
	raw = utils.NewArray()
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
