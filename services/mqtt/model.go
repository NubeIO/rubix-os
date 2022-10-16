package mqtt

import "github.com/NubeIO/flow-framework/mqttclient"

type PointMqtt struct {
	client *mqttclient.Client
	QOS    mqttclient.QOS
}

type PointListPayload struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type PointCovPayload struct {
	Value    *float64 `json:"value"`
	ValueRaw *float64 `json:"value_raw"`
	Ts       string   `json:"ts"`
	Priority *float64 `json:"priority"`
}
