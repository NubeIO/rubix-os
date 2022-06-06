package interfaces

type FlowNetworkMetadata struct {
	UUID       string            `json:"uuid"`
	Name       string            `json:"name"`
	ClientName string            `json:"client_name"`
	SiteName   string            `json:"site_name"`
	DeviceName string            `json:"device_name"`
	Streams    []*StreamMetadata `json:"streams"`
}
