package rubixregistry

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

func (inst *RubixRegistry) CreateGlobalUUIDIfDoesNotExist() error {
	dirExist := DirExists(inst.Dir)
	if !dirExist {
		if err := os.MkdirAll(inst.Dir, os.FileMode(inst.FileMode)); err != nil {
			panic(err)
		}
	}
	fileExist := FileExists(inst.GlobalUUIDFile)
	if !fileExist {
		deprecatedDeviceInfo, err := inst.GetLegacyDeviceInfo()
		var globalUUID string
		if err != nil {
			globalUUID = ShortUUID("glb")
		} else {
			globalUUID = deprecatedDeviceInfo.GlobalUUID
		}
		err = os.WriteFile(inst.GlobalUUIDFile, []byte(globalUUID), os.FileMode(inst.FileMode))
		if err != nil {
			return err
		}
	}
	return nil
}

func DirExists(dirPath string) bool {
	f, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return f.IsDir()
}

func FileExists(filePath string) bool {
	f, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

func ShortUUID(prefix ...string) string {
	u := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, u)
	if n != len(u) || err != nil {
		return "-error-uuid-"
	}
	uuid := fmt.Sprintf("%x%x", u[0:4], u[4:6])
	if len(prefix) > 0 {
		uuid = fmt.Sprintf("%s_%s", prefix[0], uuid)
	}
	return uuid
}
