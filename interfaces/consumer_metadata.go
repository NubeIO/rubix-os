package interfaces

type ConsumerMetadata struct {
	UUID    string            `json:"uuid"`
	Name    string            `json:"name"`
	Writers []*WriterMetadata `json:"writers"`
}
