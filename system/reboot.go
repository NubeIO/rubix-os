package system

import "github.com/NubeDev/flow-framework/system/command"

// RebootNow Reboot system
func RebootNow() (result string, err error) {
	return command.Run("sudo", "shutdown", "-r", "now")
}

// ShutdownNow Shutdown system
//func ShutdownNow() (result string, err error) {
//	return command.Run("sudo", "shutdown", "-h", "now")
//}
