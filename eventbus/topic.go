package eventbus

// bus topics ...
const (
	All = ".*"
	JobsAll  		= "jobs.*"
	NetworksAll  	= "network.*"
	NetworkCreated  = "network.created"
	NetworkEnabled  = "network.enabled"
	DevicesAll 		= "device.*"

	PointsAll 		= "point.*"
	PointEnabled  	= "point.enabled"
	PointCreated  	= "point.created"
	PointUpdated  	= "point.updated"

)


// BusTopics return all bus topics
func BusTopics() []string {

	return []string{
		All,
		JobsAll,
		NetworksAll,
		NetworkCreated,
		DevicesAll,
		PointsAll,
		PointCreated,
		PointUpdated,
		PointEnabled,
	}
}
