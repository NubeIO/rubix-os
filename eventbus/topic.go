package eventbus

// bus topics ...
const (
	All = ".*"
	DevicesAll 		= "device.*"
	PointsAll 		= "point.*"
	JobsAll  		= "jobs.*"
	NetworksAll  	= "network.*"
	PluginsAll  	= "plugin.*"
	StreamsAll  	= "stream.*"

	PluginsCreated  = "plugin.created"
	PluginsUpdated  = "plugin.updated"

	StreamsCreated  = "stream.created"
	StreamsUpdated  = "stream.updated"

	NetworkCreated  = "network.created"
	NetworkUpdated  = "network.updated"

	PointUpdated  	= "point.updated"
	PointCreated  	= "point.created"
	PointCOV  		= "point.cov"

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
		PluginsCreated,
		PluginsUpdated,
		StreamsCreated,
		StreamsUpdated,
		NetworkCreated,
		NetworkUpdated,
		PointCreated,
		PointUpdated,
		PointCOV,
	}
}
