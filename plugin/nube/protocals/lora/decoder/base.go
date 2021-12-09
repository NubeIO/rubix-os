package decoder

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

type CommonValues struct {
	Sensor string `json:"sensor"`
	Id     string `json:"id"`
	Rssi   int    `json:"rssi"`
}

func getDeviceDescriptionFromPayload(data string) *LoRaDeviceDescription {
	sensorCode := data[2:4]
	return GetLoRaDeviceDescription(sensorCode)
}

func CheckPayloadLength(data string) bool {
	log.Println(data)
	dl := len(data)
	if dl <= 4 {
		return false
	}
	if data == "!\r\n" || data == "!\n" {
		return false
	}
	dev := getDeviceDescriptionFromPayload(data)
	if dev == &NilLoRaDeviceDescription {
		return false
	}
	return dev.CheckLength(data)
}

func DecodePayload(data string) (*CommonValues, interface{}) {
	devDesc := getDeviceDescriptionFromPayload(data)
	if devDesc == &NilLoRaDeviceDescription {
		return &CommonValues{}, nil
	}

	cmn, payload := devDesc.Decode(data, devDesc)
	common(cmn, data, devDesc.Model)
	return cmn, payload
}

func common(payload *CommonValues, data string, sensor string) {
	payload.Sensor = sensor
	payload.Id = decodeID(data)
	payload.Rssi = rssi(data)
}

func decodeID(data string) string {
	return data[0:8]
}

func rssi(data string) int {
	dataLen := len(data)
	v, _ := strconv.ParseInt(data[dataLen-4:dataLen-2], 16, 0)
	v = v * -1
	return int(v)
}
