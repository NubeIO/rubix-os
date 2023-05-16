package installer

import (
	"fmt"
	"strings"
)

type BuildDetails struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Arch    string `json:"arch,omitempty"`
	ZipName string `json:"zip_name,omitempty"`
}

func (inst *Installer) GetZipBuildDetails(zipName string) *BuildDetails {
	parts := strings.Split(zipName, "-")
	if len(parts) > 2 {
		version := parts[len(parts)-2]
		if !strings.Contains(version, "v") {
			version = fmt.Sprintf("v%s", version)
		}
		archContent := parts[len(parts)-1]
		archParts := strings.Split(archContent, ".")
		arch := ""
		if len(parts) > 1 {
			arch = archParts[1]
		}
		nameParts := parts[:len(parts)-2]
		name := strings.Join(nameParts, "-")
		return &BuildDetails{
			Name:    name,
			Version: version,
			Arch:    arch,
			ZipName: zipName,
		}
	}
	return &BuildDetails{
		ZipName: zipName,
	}
}
