package interfaces

type Action string

const (
	Enable  Action = "enable"
	Disable Action = "disable"
	Start   Action = "start"
	Stop    Action = "stop"
	Restart Action = "restart"
)
