package systemctl

import (
	"fmt"
	"github.com/NubeIO/rubix-os/installer"
	"github.com/NubeIO/rubix-os/src/cli/constants"
	"testing"
)

func TestStore_generateServiceFile(t *testing.T) {
	tmpDir, absoluteServiceFileName, err := GenerateServiceFile(&ServiceFile{
		Name:                        constants.RubixOs,
		Version:                     "v0.0.1",
		ExecStart:                   "app -p 1660 -g /data/rubix-os -d data --prod",
		AttachWorkingDirOnExecStart: true,
	}, installer.New(&installer.Installer{}))
	fmt.Println(tmpDir, absoluteServiceFileName, err)
	if err != nil {
		return
	}
}
