package model

//FlowNetwork flow network
type FlowNetwork struct {
	CommonUUID
	CommonName
	CommonDescription
	GlobalFlowID            string    `json:"global_flow_id"`
	GlobalRemoteFlowID      string    `json:"global_remote_flow_id"` //if is a remote
	RemoteFlowUUID          string    `json:"remote_flow_uuid"`      //if is a remote
	IsRemote                bool      `json:"is_remote"`
	FetchHistories          bool      `json:"fetch_histories"`
	FetchHistoriesFrequency int       `json:"fetch_hist_frequency"`      //time example 15min
	DeleteHistoriesOnFetch  bool      `json:"delete_histories_on_fetch"` //drop the histories on the producer device on migration
	IsMQTT                  bool      `json:"is_mqtt"`
	FlowIP                  string    `json:"flow_ip"`
	FlowPort                string    `json:"flow_port"`
	FlowHTTPS               bool      `json:"flow_https"`
	FlowUsername            string    `json:"flow_username"`
	FlowPassword            string    `json:"flow_password"`
	FlowToken               string    `json:"flow_token"`
	MqttIP                  string    `json:"mqtt_ip"`
	MqttPort                string    `json:"mqtt_port"`
	MqttHTTPS               bool      `json:"mqtt_https"`
	MqttUsername            string    `json:"mqtt_username"`
	MqttPassword            string    `json:"mqtt_password"`
	Streams                 []*Stream `gorm:"many2many:streams_flow_networks;"`
	CommonCreated
}
