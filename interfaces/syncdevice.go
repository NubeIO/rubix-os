package interfaces

type SyncDevice struct {
	NetworkUUID     string
	NetworkName     string
	DeviceUUID      string
	DeviceName      string
	FlowNetworkUUID string
	IsLocal         bool
}
