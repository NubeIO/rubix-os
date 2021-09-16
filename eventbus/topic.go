package eventbus

// bus topics ...
const (
	All         = ".*"
	DevicesAll  = "device.*"
	PointsAll   = "point.*"
	JobsAll     = "jobs.*"
	NetworksAll = "network.*"
	PluginsAll  = "plugin.*"
	StreamsAll  = "stream.*"
	NodesAll    = "node.*"
	ProducerAll = "producer.*"
	MQTTAll     = "MQTT.*"

	MQTTCreated = "MQTT.created"
	MQTTUpdated = "MQTT.updated"

	PluginsCreated = "plugin.created"
	PluginsUpdated = "plugin.updated"
	PluginsDeleted = "plugin.deleted"

	ProducerCreated = "producer.created"
	ProducerUpdated = "producer.updated"
	ProducerEvent   = "producer.event"

	StreamsCreated = "stream.created"
	StreamsUpdated = "stream.updated"

	Network        = "network"
	NetworkCreated = "network.created"
	NetworkUpdated = "network.updated"

	NetDevUpdated = "network.dev.updated"

	PointUpdated = "point.updated"
	PointCreated = "point.created"
	PointCOV     = "point.cov"

	NodeUpdated  = "node.updated"
	NodeCreated  = "node.created"
	NodeEvent    = "node.event"
	NodeEventIn  = "node.event.in"
	NodeEventOut = "node.event.out"
)

// BusTopics return all bus topics
func BusTopics() []string {
	return []string{
		All,
		JobsAll,
		NetworksAll,
		DevicesAll,
		PointsAll,
		PluginsAll,
		StreamsAll,
		NodesAll,
		MQTTAll,
		PluginsCreated,
		PluginsUpdated,
		PluginsDeleted,
		StreamsCreated,
		StreamsUpdated,
		Network,
		NetworkCreated,
		NetworkUpdated,
		NetDevUpdated,
		PointCreated,
		PointUpdated,
		PointCOV,
		NodeUpdated,
		NodeCreated,
		NodeEvent,
		NodeEventIn,
		NodeEventOut,
		ProducerAll,
		ProducerCreated,
		ProducerUpdated,
		ProducerEvent,
		MQTTCreated,
		MQTTUpdated,
	}
}
