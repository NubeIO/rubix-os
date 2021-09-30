package model

type MqttConnection struct {
	CommonUUID
	Enabled                       bool   `json:"enabled,omitempty"`
	Master                        bool   `json:"master,omitempty"`
	Name                          string `json:"name,omitempty"`
	Host                          string `json:"host,omitempty"`
	Port                          int    `json:"port,omitempty"`
	Authentication                bool   `json:"authentication,omitempty"`
	Username                      string `json:"username,omitempty"`
	Password                      string `json:"password,omitempty"`
	Keepalive                     int    `json:"keepalive,omitempty"`
	Qos                           int    `json:"qos,omitempty"`
	Retain                        bool   `json:"retain,omitempty"`
	AttemptReconnectOnUnavailable bool   `json:"attempt_reconnect_on_unavailable,omitempty"`
	AttemptReconnectSecs          int    `json:"attempt_reconnect_secs,omitempty"`
	Timeout                       int    `json:"timeout,omitempty"`
	IntegrationUUID               string `json:"integration_uuid" gorm:"TYPE:varchar(255) REFERENCES integrations;null;default:null"`
}
