package model

type FlowNetwork struct {
	CommonUUID
	CommonName
	CommonDescription
	GlobalUUID              string    `json:"global_uuid"`
	ClientId                string    `json:"client_id"`
	ClientName              string    `json:"client_name"`
	SiteId                  string    `json:"site_id"`
	SiteName                string    `json:"site_name"`
	DeviceId                string    `json:"device_id"`
	DeviceName              string    `json:"device_name"`
	IsRemote                *bool     `json:"is_remote"`
	FetchHistories          *bool     `json:"fetch_histories"`
	FetchHistoriesFrequency int       `json:"fetch_hist_frequency"`      //time example 15min
	DeleteHistoriesOnFetch  *bool     `json:"delete_histories_on_fetch"` //drop the histories on the producer device on migration
	IsMQTT                  *bool     `json:"is_mqtt"`
	FlowIP                  string    `json:"flow_ip"`
	FlowPort                string    `json:"flow_port"`
	FlowHTTPS               *bool     `json:"flow_https"`
	FlowUsername            string    `json:"flow_username"`
	FlowPassword            string    `json:"flow_password"`
	FlowToken               string    `json:"flow_token"`
	MqttIP                  string    `json:"mqtt_ip"`
	MqttPort                string    `json:"mqtt_port"`
	MqttHTTPS               *bool     `json:"mqtt_https"`
	MqttUsername            string    `json:"mqtt_username"`
	MqttPassword            string    `json:"mqtt_password"`
	Streams                 []*Stream `json:"streams" gorm:"many2many:flow_networks_streams;constraint:OnDelete:CASCADE"`
	CommonCreated
}
