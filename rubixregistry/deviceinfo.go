package rubixregistry

import (
	"encoding/json"
	"os"
)

func (inst *RubixRegistry) GetDeviceInfo() (*DeviceInfo, error) {
	data, err := os.ReadFile(inst.GlobalUUIDFile)
	if err != nil {
		return nil, err
	}
	return &DeviceInfo{GlobalUUID: string(data)}, nil
}

func (inst *RubixRegistry) GetLegacyDeviceInfo() (*DeviceInfo, error) { // TODO: remove after migration done
	data, err := os.ReadFile(inst.LegacyDeviceInfoFile)
	if err != nil {
		return nil, err
	}
	deviceInfoDefault := DeviceInfoDefault{}
	err = json.Unmarshal(data, &deviceInfoDefault)
	if err != nil {
		return nil, err
	}
	return &deviceInfoDefault.DeviceInfoFirstRecord.DeviceInfo, nil
}
