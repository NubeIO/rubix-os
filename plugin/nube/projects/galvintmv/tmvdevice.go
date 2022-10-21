package main

type TMVDevices struct {
	Devices []TMVDevice `json:"tmv-devices`
}

type TMVDevice struct {
	DeviceLocation      string  `json:"Device Location"`
	DeviceName          string  `json:"Device Name"`
	TemperatureSetpoint float64 `json:"Temperature Setpoint"`
	TMVNumber           string  `json:"TMV No."`
	SolenoidRequired    string  `json:"Solenoid Required"`
	DeviceAddress       int     `json:"Device Address"`
	LoRaWANDeviceEUI    string  `json:"LoRaWAN Device EUI"`
}
