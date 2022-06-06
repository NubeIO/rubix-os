package interfaces

type StreamMetadata struct {
	UUID      string              `json:"uuid"`
	Name      string              `json:"name"`
	Producers []*ProducerMetadata `json:"producers"`
}

type StreamCloneMetadata struct {
	UUID      string              `json:"uuid"`
	Name      string              `json:"name"`
	Consumers []*ConsumerMetadata `json:"consumers"`
}
