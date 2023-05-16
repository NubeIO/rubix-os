package interfaces

import (
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/lib-systemctl-go/systemd"
)

type AppResponse struct {
	Name              string                     `json:"name"`
	Version           string                     `json:"version,omitempty"`
	State             *systemctl.SystemState     `json:"state,omitempty"`
	Error             string                     `json:"error,omitempty"`
	UninstallResponse *systemd.UninstallResponse `json:"remove_response"`
}

type Apps struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Path    string `json:"path,omitempty"`
}

type AppsStatus struct {
	Name        string                 `json:"name,omitempty"`
	Version     string                 `json:"version,omitempty"`
	ServiceName string                 `json:"service_name,omitempty"`
	State       *systemctl.SystemState `json:"state,omitempty"`
}
