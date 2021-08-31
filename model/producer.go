package model

//WriterClone list of all the consumers
// a consumer
type WriterClone struct {
	CommonUUID
	ProducerUUID 		string  `json:"producer_uuid" gorm:"TYPE:string REFERENCES producers;not null;default:null"` // is the producer UUID
	WriteValue       	float64  `json:"write_value"` // for common use of points
	ConsumerUUID 		string 	`json:"consumer_uuid"`  // is the remote consumer UUID, ie: whatever is subscribing to this producer
	WriterUUID 			string 	`json:"writer_uuid"`  // is the remote consumer UUID, ie: whatever is subscribing to this producer
	ConsumerCOV 		float64 `json:"consumer_cov"`
	CommonCreated
}


//Producer a producer is a placeholder to register an object to enable consumers to
// A producer for example is a point, Something that makes data, and the subscriber would have a consumer to it, Like grafana reading and writing to it from edge to cloud or wires over rest(peer to peer)
type Producer struct {
	CommonProducer
	PresentValue 			float64  `json:"present_value"` //these fields are support as points is the most common use case for histories
	SLWriteUUID       		string  `json:"sl_write_uuid"`
	ProducerType 			string  `json:"producer_type"` //point, schedule, job, network
	EnableHistory 			bool 	`json:"enable_history"`
	ProducerApplication 	string 	`json:"producer_application"`
	StreamUUID     			string 	`json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	ProducerThingUUID 		string  `json:"producer_thing_uuid"` //a producer_uuid is the point uuid
	WriterClone				[]WriterClone `json:"writer_clones" gorm:"constraint:OnDelete:CASCADE;"`
	ProducerHistory			[]ProducerHistory `json:"producer_histories" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}

