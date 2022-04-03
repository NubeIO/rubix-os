package bmmodel

// MqttPayload payload from the bacnet server
type MqttPayload struct {
	Value    *float64
	Priority int
}
