package interfaces

type NetworkTagForPostgresSync struct {
	NetworkUUID string `json:"network_uuid"`
	Tag         string `json:"tag"`
}
