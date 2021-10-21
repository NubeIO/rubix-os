package model

type FlowNetworkCredential struct {
	FlowIP       string `json:"flow_ip"`
	FlowPort     int    `json:"flow_port"`
	FlowUsername string `json:"flow_username"`
	FlowPassword string `json:"flow_password"`
}
