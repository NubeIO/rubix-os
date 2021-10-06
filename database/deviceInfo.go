package database

import (
	"encoding/json"
	"github.com/NubeDev/flow-framework/model"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func (d *GormDatabase) GetDeviceInfo() (*model.DeviceInfo, error) {
	file, err := os.Open("/data/rubix-registry/device_info.json")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	deviceInfo, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var rawDeviceInfo map[string]map[string]interface{}
	if err = json.Unmarshal(deviceInfo, &rawDeviceInfo); err != nil {
		return nil, err
	}
	extractedRawDeviceInfo := rawDeviceInfo["_default"]["1"]
	marshalledExtractedRawDeviceInfo, err := json.Marshal(extractedRawDeviceInfo)
	if err != nil {
		return nil, err
	}
	deviceInfoModel := model.DeviceInfo{}
	if err = json.Unmarshal(marshalledExtractedRawDeviceInfo, &deviceInfoModel); err != nil {
		return nil, err
	}
	return &deviceInfoModel, nil
}
