package model

//FlowNetwork flow network
type FlowNetwork struct {
	CommonFlowNetworkUUID
	CommonFlowNetworkName
	IsRemote 	bool
	FlowIP 		string
	FlowPort 	string
	FlowHTTPS 	bool
	FlowUsername string
	FlowPassword string
	MqttIP 		string
	MqttPort 	string
	MqttHTTPS 	bool
	MqttUsername string
	MqttPassword string

}

