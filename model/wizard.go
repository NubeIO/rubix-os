package model

type P2PBody struct {
	LocalFlowIP        *string `json:"local_flow_ip"`
	LocalFlowPort      *int    `json:"local_flow_port"`
	LocalFlowUsername  *string `json:"local_flow_username"`
	LocalFlowPassword  *string `json:"local_flow_password"`
	RemoteFlowIP       *string `json:"remote_flow_ip"`
	RemoteFlowPort     *int    `json:"remote_flow_port"`
	RemoteFlowUsername *string `json:"remote_flow_username"`
	RemoteFlowPassword *string `json:"remote_flow_password"`
}
