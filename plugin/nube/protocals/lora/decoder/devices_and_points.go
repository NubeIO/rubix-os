package decoder

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

type LoRaDeviceDescription struct {
	DeviceName      string
	Model           string
	SensorCode      string
	CheckLength     func(data string) bool
	Decode          func(data string, devDesc *LoRaDeviceDescription) (*CommonValues, interface{})
	GetPointsStruct func() interface{}
}

var NilLoRaDeviceDescription = LoRaDeviceDescription{
	DeviceName:      "",
	Model:           "",
	SensorCode:      "",
	CheckLength:     NilLoRaDeviceDescriptionCheckLength,
	Decode:          NilLoRaDeviceDescriptionDecode,
	GetPointsStruct: NilLoRaDeviceDescriptionGetPointsStruct,
}

func NilLoRaDeviceDescriptionCheckLength(data string) bool {
	return false
}

func NilLoRaDeviceDescriptionDecode(data string, devDesc *LoRaDeviceDescription) (*CommonValues, interface{}) {
	return &CommonValues{}, struct{}{}
}

func NilLoRaDeviceDescriptionGetPointsStruct() interface{} {
	return struct{}{}
}

var LoRaDeviceDescriptions = [...]LoRaDeviceDescription{
	{
		DeviceName:      "MicroEdge",
		Model:           "MicroEdge",
		CheckLength:     CheckPayloadLengthME,
		Decode:          DecodeME,
		GetPointsStruct: GetPointsStructME,
	},
	{
		DeviceName:      "Droplet",
		Model:           "THLM",
		CheckLength:     CheckPayloadLengthDroplet,
		Decode:          DecodeDropletTHLM,
		GetPointsStruct: GetPointsStructTHLM,
	},
	{
		DeviceName:      "Droplet",
		Model:           "TH",
		CheckLength:     CheckPayloadLengthDroplet,
		Decode:          DecodeDropletTH,
		GetPointsStruct: GetPointsStructTH,
	},
	{
		DeviceName:      "Droplet",
		Model:           "THL",
		CheckLength:     CheckPayloadLengthDroplet,
		Decode:          DecodeDropletTHL,
		GetPointsStruct: GetPointsStructTHL,
	},
	{
		DeviceName:      "Droplet",
		Model:           "THLM",
		CheckLength:     CheckPayloadLengthDroplet,
		Decode:          DecodeDropletTHLM,
		GetPointsStruct: GetPointsStructTHLM,
	},
	{
		DeviceName:      "ZipHydroTap",
		Model:           "ZipHydroTap",
		CheckLength:     CheckPayloadLengthZHT,
		Decode:          DecodeZHT,
		GetPointsStruct: GetPointsStructZHT,
	},
}

func GetDeviceDescription(device *model.Device) *LoRaDeviceDescription {
	for _, dev := range LoRaDeviceDescriptions {
		if device.Model == dev.Model {
			return &dev
		}
	}
	return &NilLoRaDeviceDescription
}

func GetDevicePointsStruct(device *model.Device) interface{} {
	return GetDeviceDescription(device).GetPointsStruct()
}
