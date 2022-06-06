package interfaces

type StreamMetadata struct {
	UUID      string              `json:"uuid"`
	Name      string              `json:"name"`
	Producers []*ProducerMetadata `json:"producers"`
}
