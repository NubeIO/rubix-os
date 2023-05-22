package interfaces

type Plugin struct {
	Name      string `json:"name"`
	Arch      string `json:"arch"`
	Version   string `json:"version,omitempty"`
	Extension string `json:"extension"`
}
