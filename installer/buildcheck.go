package installer

import (
	"fmt"
	"regexp"
)

type BuildDetails struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Arch    string `json:"arch,omitempty"`
	ZipName string `json:"zip_name,omitempty"`
}

func (inst *Installer) GetZipBuildDetails(zipName string) *BuildDetails {
	infoRegex := regexp.MustCompile(`^([\w-]+)-([\d.]+(?:-[\w\d]+)?)-(\w+)\.(\w+)\.zip$`)
	match := infoRegex.FindStringSubmatch(zipName)
	if match != nil {
		name := match[1]
		version := match[2]
		arch := match[4]
		return &BuildDetails{
			Name:    name,
			Version: fmt.Sprintf("v%s", version),
			Arch:    arch,
			ZipName: zipName,
		}
	}
	return &BuildDetails{
		ZipName: zipName,
	}
}
