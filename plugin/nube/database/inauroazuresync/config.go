package main

import (
	"errors"
	"fmt"
	"time"
)

const descriptionIndexGatewayID = 0
const descriptionIndexKey = 1
const descriptionIndexHost = 2

const timestampRangeToCombine = 5 * time.Second

type Azure struct {
	HostName       string                    `json:"azure_host_name"`
	GatewayDetails map[string]GatewayDetails `json:"gateway_details"` // gateway details stored by host uuid
}

type GatewayDetails struct {
	AzureDeviceId  string `json:"azure_device_id"`
	AzureDeviceKey string `json:"azure_device_key"`
	RouterIMEI     string `json:"router_imei"`
	SIMICCID       string `json:"sim_iccid"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	NetworkType    string `json:"network_type"` // TODO: look into replacing with automatic update from host/edge network info
}

type Postgres struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`
	Port     int    `yaml:"port"`
	SslMode  string `yaml:"ssl_mode"`
}

type Job struct {
	SensorHistorySyncFrequency  string `yaml:"sensor_history_sync_frequency"`
	GatewayPayloadSyncFrequency string `yaml:"gateway_payload_sync_frequency"`
}

type Config struct {
	Azure    Azure  `yaml:"azure"`
	Job      Job    `yaml:"job"`
	LogLevel string `yaml:"log_level"`
}

func (inst *Instance) DefaultConfig() interface{} {
	azure := Azure{
		HostName: "",
		GatewayDetails: map[string]GatewayDetails{
			"host_xxxxxxxxxx": {
				AzureDeviceId:  "",
				AzureDeviceKey: "",
				SIMICCID:       "",
				Latitude:       "",
				Longitude:      "",
				NetworkType:    "",
			},
		},
	}
	job := Job{
		SensorHistorySyncFrequency:  "5m",
		GatewayPayloadSyncFrequency: "15m",
	}
	return &Config{
		Azure:    azure,
		Job:      job,
		LogLevel: "ERROR", // DEBUG or ERROR
	}
}

/*  EXAMPLE CONFIG YAML
azure:
  azure_host_name: "NubeTestHub1.azure-devices.net"
  gateway_details:
	host-uuid-1:
	  azure_device_id: "NubeTestDevice1"
	  azure_device_key: "DY/IrrTaQ+K/gU8S5v2B5HSaotirM0lMzaWbAqBJl3U="
	  sim_iccid: "123456789"
	  latitude: "123.45"
	  longitude: "67.89"
	  network_type: "Cellular"
	host-uuid-2:
	  azure_device_id: ""
	  azure_device_key: ""
	  sim_iccid: ""
	  latitude: ""
	  longitude: ""
	  network_type: ""
job:
  sensor_history_sync_frequency: "5m"
  gateway_payload_sync_frequency: "10m"
log_level: "DEBUG"
*/

func (inst *Instance) GetConfig() interface{} {
	return inst.config
}

func (inst *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	inst.config = newConfig
	return nil
}

func (inst *Instance) getAzureDetailsFromConfigByHostUUID(hostUUID string) (azureSAS, azureHost, azureDeviceId, azureDeviceKey string, err error) {
	if inst.config == nil {
		return azureSAS, azureHost, azureDeviceId, azureDeviceKey, errors.New("invalid inauroazuresync module configuration")
	}

	gatewayDetails, ok := inst.config.Azure.GatewayDetails[hostUUID]
	if !ok {
		return azureSAS, azureHost, azureDeviceId, azureDeviceKey, errors.New(fmt.Sprint("add gateway azure details in module config.  host: ", hostUUID))
	}

	azureHost = inst.config.Azure.HostName
	if azureHost == "" {
		return azureSAS, azureHost, azureDeviceId, azureDeviceKey, errors.New("set azure_host_name in module config")
	}

	azureDeviceId = gatewayDetails.AzureDeviceId
	if azureDeviceId == "" {
		return azureSAS, azureHost, azureDeviceId, azureDeviceKey, errors.New(fmt.Sprint("add gateway azure_device_id in module config.  host: ", hostUUID))
	}

	azureDeviceKey = gatewayDetails.AzureDeviceKey
	if azureDeviceKey == "" {
		return azureSAS, azureHost, azureDeviceId, azureDeviceKey, errors.New(fmt.Sprint("add gateway azure_device_key in module config.  host: ", hostUUID))
	}

	azureSAS = fmt.Sprintf("HostName=%s;DeviceId=%s;SharedAccessKey=%s", azureHost, azureDeviceId, azureDeviceKey)
	return azureSAS, azureHost, azureDeviceId, azureDeviceKey, nil
}

func (inst *Instance) getSAS(azureDeviceId, azureDeviceKey string) (sas string, err error) {
	if inst.config == nil {
		return "", errors.New("invalid inauroazuresync module configuration")
	}

	azureHost := inst.config.Azure.HostName
	if azureHost == "" {
		return "", errors.New("set azure_host_name in module config")
	}

	return fmt.Sprintf("HostName=%s;DeviceId=%s;SharedAccessKey=%s", azureHost, azureDeviceId, azureDeviceKey), nil
}

func (inst *Instance) getGatewayDetailsFromConfig() (gatewayDetails map[string]GatewayDetails, err error) {
	if inst.config == nil {
		return nil, errors.New("invalid inauroazuresync module configuration")
	}

	return inst.config.Azure.GatewayDetails, nil
}
