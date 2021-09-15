package pkgmodel

// MqttPayload payload from the bacnet server
type MqttPayload struct {
	Value    float64
	Priority int
}

type BacnetPoint struct {
	ObjectType           string  `json:"object_type"`
	ObjectName           string  `json:"object_name"`
	Address              int     `json:"address"`
	RelinquishDefault    float64 `json:"relinquish_default"`
	EventState           string  `json:"event_state"`
	Units                string  `json:"units"`
	Description          string  `json:"description"`
	Enable               bool    `json:"enable"`
	Fault                bool    `json:"fault"`
	DataRound            float64 `json:"data_round"`
	DataOffset           float64 `json:"data_offset"`
	UseNextAvailableAddr *bool   `json:"use_next_available_address"`
	COV                  float32 `json:"cov"`
	Priority             `json:"priority_array_write"`
}

type Priority struct {
	P1        float64 `json:"_1"` //would be better if we stored the TS and where it was written from, for example from a Remote Producer
	//P2        null.Float `json:"_2"`
	//P3        null.Float `json:"_3"`
	//P4        null.Float `json:"_4"`
	//P5        null.Float `json:"_5"`
	//P6        null.Float `json:"_6"`
	//P7        null.Float `json:"_7"`
	//P8        null.Float `json:"_8"`
	//P9        null.Float `json:"_9"`
	//P10       null.Float `json:"_10"`
	//P11       null.Float `json:"_11"`
	//P12       null.Float `json:"_12"`
	//P13       null.Float `json:"_13"`
	//P14       null.Float `json:"_14"`
	//P15       null.Float `json:"_15"`
	//P16       null.Float `json:"_16"` //removed and added to the point to save one DB write

}
