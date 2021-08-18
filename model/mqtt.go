package model


//TODO at binod to sort out
type MasterMqttConnection struct {
	Enabled                       bool   `json:"enabled"`
	Master                        bool   `json:"master"`
	Name                          string `json:"name"`
	Host                          string `json:"host"`
	Port                          int    `json:"port"`
	Authentication                bool   `json:"authentication"`
	Username                      string `json:"username"`
	Password                      string `json:"password"`
	Keepalive                     int    `json:"keepalive"`
	Qos                           int    `json:"qos"`
	Retain                        bool   `json:"retain"`
	AttemptReconnectOnUnavailable bool   `json:"attempt_reconnect_on_unavailable"`
	AttemptReconnectSecs          int    `json:"attempt_reconnect_secs"`
	Timeout                       int    `json:"timeout"`

}


type MqttConnection struct {
	Enabled                       bool   `json:"enabled"`
	Master                        bool   `json:"master"`
	Name                          string `json:"name"`
	Host                          string `json:"host"`
	Port                          int    `json:"port"`
	Authentication                bool   `json:"authentication"`
	Username                      string `json:"username"`
	Password                      string `json:"password"`
	Keepalive                     int    `json:"keepalive"`
	Qos                           int    `json:"qos"`
	Retain                        bool   `json:"retain"`
	AttemptReconnectOnUnavailable bool   `json:"attempt_reconnect_on_unavailable"`
	AttemptReconnectSecs          int    `json:"attempt_reconnect_secs"`
	Timeout                       int    `json:"timeout"`

}