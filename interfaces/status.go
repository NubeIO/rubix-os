package interfaces

type SystemCtlProperty struct {
	Property string `json:"property"`
}

type SystemCtlStatus struct {
	Status string `json:"status"`
}

type SystemCtlState struct {
	State bool `json:"state"`
}

type SystemCtlStateStatus struct {
	State  bool   `json:"state"`
	Status string `json:"status"`
}

type CreateStatus string

const (
	CreateNotAvailable CreateStatus = "N/A"
	Creating           CreateStatus = "Creating"
	Created            CreateStatus = "Created"
	CreateFailed       CreateStatus = "Failed"
)

type RestoreStatus string

const (
	RestoreNotAvailable RestoreStatus = "N/A"
	Restoring           RestoreStatus = "Restoring"
	Restored            RestoreStatus = "Restored"
	RestoreFailed       RestoreStatus = "Failed"
)

type SnapshotStatus struct {
	CreateStatus  CreateStatus  `json:"create_status"`
	RestoreStatus RestoreStatus `json:"restore_status"`
}
