package mqttclient

import (
	"github.com/NubeDev/flow-framework/utils"
	"strings"
)

func TopicParts(topic string) *utils.Array {
	s := strings.SplitAfter(topic, "/")
	arr := utils.NewArray()
	for _, e := range s {
		res := strings.ReplaceAll(e, "/", "")
		if res != "" {
			arr.Add(res)
		}
	}
	return arr
}
