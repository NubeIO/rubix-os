package interfaces

type SyncModel struct {
	UUID    string  `json:"uuid"`
	IsError bool    `json:"is_error"`
	Message *string `json:"message"`
}
