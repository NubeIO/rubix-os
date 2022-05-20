package interfaces

type SyncModel struct {
	Id      string  `json:"id"`
	IsError bool    `json:"is_error"`
	Message *string `json:"message"`
}
