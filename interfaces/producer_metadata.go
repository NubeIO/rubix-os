package interfaces

type ProducerMetadata struct {
	UUID         string                 `json:"uuid"`
	Name         string                 `json:"name"`
	WriterClones []*WriterCloneMetadata `json:"writer_clones"`
}
