package interfaces

type FileExistence struct {
	File   string `json:"file"`
	Exists bool   `json:"exists"`
}

type FileUploadResponse struct {
	Destination string `json:"destination"`
	File        string `json:"file"`
	Size        string `json:"size"`
	UploadTime  string `json:"upload_time"`
}

type WriteFile struct {
	FilePath     string      `json:"path"`
	Body         interface{} `json:"body"`
	BodyAsString string      `json:"body_as_string"`
}

type WriteFileData struct {
	Data string `json:"data"`
}

type WriteFormatFile struct {
	FilePath     string      `json:"path"`
	Body         interface{} `json:"body"`
	BodyAsString string      `json:"body_as_string"`
}
