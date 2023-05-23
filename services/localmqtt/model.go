package localmqtt

import "github.com/NubeIO/rubix-os/mqttclient"

type LocalMqtt struct {
	Client                *mqttclient.Client
	QOS                   mqttclient.QOS
	Retain                bool
	GlobalBroadcast       bool
	PublishPointCOV       bool
	PublishPointList      bool
	PointWriteListener    bool
	PublishScheduleCOV    bool
	PublishScheduleList   bool
	ScheduleWriteListener bool
}

type PointListPayload struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type ScheduleList struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type PointCovPayload struct {
	Value    *float64 `json:"value"`
	ValueRaw *float64 `json:"value_raw"`
	Ts       string   `json:"ts"`
	Priority *int     `json:"priority"`
}
