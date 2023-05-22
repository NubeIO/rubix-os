package interfaces

type DirExistence struct {
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
}
