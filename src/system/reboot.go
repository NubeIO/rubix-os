package system

import (
	"github.com/NubeIO/flow-framework/src/system/command"
)

// RebootNow Reboot system
func RebootNow() (result string, err error) {
	return command.Run("sudo", "shutdown", "-r", "now")
}
