package interfaces

type WriterCloneMetadata struct {
	UUID             string `json:"uuid"`
	WriterThingClass string `json:"writer_thing_class"`
	WriterThingUUID  string `json:"writer_thing_uuid"`
	WriterThingName  string `json:"writer_thing_name"`
}

type WriterMetadata struct {
	UUID             string `json:"uuid"`
	WriterThingClass string `json:"writer_thing_class"`
	WriterThingUUID  string `json:"writer_thing_uuid"`
	WriterThingName  string `json:"writer_thing_name"`
}
