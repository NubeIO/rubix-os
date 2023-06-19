package rubixregistry

type DeviceInfo struct {
	GlobalUUID string `json:"global_uuid"`
}

type DeviceInfoFirstRecord struct { // TODO: remove after migration done
	DeviceInfo DeviceInfo `json:"1"`
}

type DeviceInfoDefault struct { // TODO: remove after migration done
	DeviceInfoFirstRecord DeviceInfoFirstRecord `json:"_default"`
}
