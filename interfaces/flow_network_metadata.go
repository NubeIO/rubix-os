package interfaces

type FlowNetworkMetadata struct {
	UUID       string            `json:"uuid"`
	Name       string            `json:"name"`
	ClientName string            `json:"client_name"`
	SiteName   string            `json:"site_name"`
	DeviceName string            `json:"device_name"`
	Streams    []*StreamMetadata `json:"streams"`
}

type FlowNetworkCloneMetadata struct {
	UUID         string                 `json:"uuid"`
	Name         string                 `json:"name"`
	ClientName   string                 `json:"client_name"`
	SiteName     string                 `json:"site_name"`
	DeviceName   string                 `json:"device_name"`
	StreamClones []*StreamCloneMetadata `json:"stream_clones"`
}
