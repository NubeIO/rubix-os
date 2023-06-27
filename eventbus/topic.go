package eventbus

// bus topics ...
const (
	All         = ".*"
	DevicesAll  = "device.*"
	PointsAll   = "point.*"
	JobsAll     = "job.*"
	NetworksAll = "network.*"
	PluginsAll  = "plugin.*"
	ScheduleAll = "schedule.*"
	NodesAll    = "node.*"
	MQTTAll     = "MQTT.*"

	MQTTCreated = "MQTT.created"
	MQTTUpdated = "MQTT.updated"

	JobTrigger = "job.trigger"
	JobCreated = "job.created"
	JobUpdated = "job.updated"
	JobDeleted = "job.deleted"

	SchTrigger = "sch.trigger"
	SchCreated = "sch.created"
	SchUpdated = "sch.updated"
	SchDeleted = "sch.deleted"

	PluginsCreated = "plugin.created"
	PluginsUpdated = "plugin.updated"
	PluginsDeleted = "plugin.deleted"

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
		NodesAll,
		MQTTAll,
		JobTrigger,
		JobCreated,
		JobUpdated,
		JobDeleted,
		ScheduleAll,
		SchTrigger,
		SchCreated,
		SchUpdated,
		SchDeleted,
		PluginsCreated,
		PluginsUpdated,
		PluginsDeleted,
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
		MQTTCreated,
		MQTTUpdated,
	}
}
