package rubixregistry

type RubixRegistry struct {
	Dir                  string
	LegacyDeviceInfoFile string // TODO: remove after migration done
	GlobalUUIDFile       string
	FileMode             int
}

func New() *RubixRegistry {
	rubixRegistry := RubixRegistry{
		Dir:                  "/data/rubix-registry",
		LegacyDeviceInfoFile: "/data/rubix-registry/device_info.json", // TODO: remove after migration done
		GlobalUUIDFile:       "/data/rubix-registry/global_uuid.txt",
		FileMode:             0755,
	}
	return &rubixRegistry
}
