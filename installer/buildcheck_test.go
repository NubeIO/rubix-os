package installer

import (
	"fmt"
	"testing"
)

func TestConfigEnv(t *testing.T) {
	pattern1 := "system-0.0.1-rc1-ea079f5b.amd64.zip"
	pattern2 := "system-0.0.1-ea079f5b.amd64.zip"
	pattern3 := "rubix-os-0.0.1-ea079f5b.amd64.zip"
	pattern4 := "rubix-os-1.0.1-rc1-ea079f5b.amd64.zip"

	installer := Installer{}

	buildDetails := installer.GetZipBuildDetails(pattern1)
	fmt.Println(buildDetails)

	buildDetails = installer.GetZipBuildDetails(pattern2)
	fmt.Println(buildDetails)

	buildDetails = installer.GetZipBuildDetails(pattern3)
	fmt.Println(buildDetails)

	buildDetails = installer.GetZipBuildDetails(pattern4)
	fmt.Println(buildDetails)
}
