package model

type CommonFlowNetwork struct {
	CommonUUID
	CommonSyncUUID
	CommonName
	CommonDescription
	GlobalUUID              string    `json:"global_uuid,omitempty"`
	ClientId                string    `json:"client_id,omitempty"`
	ClientName              string    `json:"client_name,omitempty"`
	SiteId                  string    `json:"site_id,omitempty"`
	SiteName                string    `json:"site_name,omitempty"`
	DeviceId                string    `json:"device_id,omitempty"`
	DeviceName              string    `json:"device_name,omitempty"`
	IsRemote                *bool     `json:"is_remote,omitempty"`
	FetchHistories          *bool     `json:"fetch_histories,omitempty"`
	FetchHistoriesFrequency int       `json:"fetch_hist_frequency,omitempty"`      //time example 15min
	DeleteHistoriesOnFetch  *bool     `json:"delete_histories_on_fetch,omitempty"` //drop the histories on the producer device on migration
	IsMQTT                  *bool     `json:"is_mqtt,omitempty"`
	FlowIP                  string    `json:"flow_ip,omitempty"`
	FlowPort                string    `json:"flow_port,omitempty"`
	FlowHTTPS               *bool     `json:"flow_https,omitempty"`
	FlowUsername            string    `json:"flow_username,omitempty"`
	FlowPassword            string    `json:"flow_password,omitempty"`
	FlowToken               string    `json:"flow_token,omitempty"`
	MqttIP                  string    `json:"mqtt_ip,omitempty"`
	MqttPort                string    `json:"mqtt_port,omitempty"`
	MqttHTTPS               *bool     `json:"mqtt_https,omitempty"`
	MqttUsername            string    `json:"mqtt_username,omitempty"`
	MqttPassword            string    `json:"mqtt_password,omitempty"`
	Streams                 []*Stream `json:"streams" gorm:"many2many:flow_networks_streams;constraint:OnDelete:CASCADE"`
	CommonCreated
}

type FlowNetwork struct {
	CommonFlowNetwork
	Streams []*Stream `json:"streams" gorm:"many2many:flow_networks_streams;constraint:OnDelete:CASCADE"`
}

type FlowNetworkClone struct {
	CommonFlowNetwork
	CommonSourceUUID
	StreamClones []*StreamClone `json:"stream_clones" gorm:"constraint:OnDelete:CASCADE;"`
}
