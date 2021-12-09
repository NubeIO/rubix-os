package decoder

// TODO: load from JSON file or something foolish

type LoRaDeviceDescription struct {
	DeviceName  string
	Model       string
	SensorCode  string
	CheckLength func(data string) bool
	Decode      func(data string, devDesc *LoRaDeviceDescription) (*CommonValues, interface{})
}

var NilLoRaDeviceDescription LoRaDeviceDescription = LoRaDeviceDescription{
	DeviceName: "",
	Model:      "",
	SensorCode: "",
}

var LoRaDeviceDescriptions = [...]LoRaDeviceDescription{
	{
		DeviceName:  "MicroEdge",
		Model:       "MicroEdge",
		SensorCode:  "AA",
		CheckLength: CheckPayloadLengthME,
		Decode:      DecodeME,
	},
	{
		DeviceName:  "Droplet",
		Model:       "THLM",
		SensorCode:  "AB",
		CheckLength: CheckPayloadLengthDroplet,
		Decode:      DecodeDropletTHLM,
	},
	{
		DeviceName:  "Droplet",
		Model:       "TH",
		SensorCode:  "B0",
		CheckLength: CheckPayloadLengthDroplet,
		Decode:      DecodeDropletTH,
	},
	{
		DeviceName:  "Droplet",
		Model:       "THL",
		SensorCode:  "B1",
		CheckLength: CheckPayloadLengthDroplet,
		Decode:      DecodeDropletTHL,
	},
	{
		DeviceName:  "Droplet",
		Model:       "THLM",
		SensorCode:  "B2",
		CheckLength: CheckPayloadLengthDroplet,
		Decode:      DecodeDropletTHLM,
	},
	{
		DeviceName:  "ZipHydroTap",
		Model:       "ZipHydroTap",
		SensorCode:  "D1",
		CheckLength: CheckPayloadLengthZHT,
		Decode:      DecodeZHT,
	},
}

func GetLoRaDeviceDescriptionFromID(devID string) *LoRaDeviceDescription {
	return GetLoRaDeviceDescription(devID[2:4])
}

func GetLoRaDeviceDescription(sensorCode string) *LoRaDeviceDescription {
	for _, dev := range LoRaDeviceDescriptions {
		if sensorCode == dev.SensorCode {
			return &dev
		}
	}
	return &NilLoRaDeviceDescription
}
