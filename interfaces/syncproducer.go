package interfaces

type SyncProducer struct {
	ProducerUUID      string `json:"producer_uuid"`
	ProducerThingName string `json:"producer_thing_name"`
	ProducerThingUUID string `json:"producer_thing_uuid"`
}
