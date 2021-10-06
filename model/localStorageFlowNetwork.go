package model

type LocalStorageFlowNetwork struct {
	FlowIP       string `json:"flow_ip,omitempty" gorm:"type:varchar(255);primaryKey"`
	FlowPort     int    `json:"flow_port,omitempty"`
	FlowHTTPS    *bool  `json:"flow_https,omitempty"`
	FlowUsername string `json:"flow_username,omitempty"`
	FlowPassword string `json:"flow_password,omitempty"`
	FlowToken    string `json:"flow_token,omitempty"`
}
