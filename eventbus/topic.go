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

	PluginsCreated = "plugin.created"
	PluginsUpdated = "plugin.updated"

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
		PluginsCreated,
		PluginsUpdated,
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
	}
}
