// Package decoder contains all device specific payload decoding and checks
package decoder

import (
	"strconv"
)

type CommonValues struct {
	Sensor string  `json:"sensor"`
	ID     string  `json:"id"`
	Rssi   int     `json:"rssi"`
	Snr    float32 `json:"snr"`
}

func DecodePayload(data string, devDesc *LoRaDeviceDescription) (*CommonValues, interface{}) {
	if !devDesc.CheckLength(data) {
		return nil, nil
	}
	cmn, payload := devDesc.Decode(data, devDesc)
	if cmn == nil {
		return nil, nil
	}
	decodeCommonValues(cmn, data, devDesc.Model)
	return cmn, payload
}

func ValidPayload(data string) bool {
	return !(len(data) <= 8)
}

func DecodeAddress(data string) string {
	return data[:8]
}

func decodeCommonValues(payload *CommonValues, data string, sensor string) {
	payload.Sensor = sensor
	payload.ID = DecodeAddress(data)
	payload.Rssi = decodeRSSI(data)
	payload.Snr = decodeSNR(data)
}

func decodeRSSI(data string) int {
	dataLen := len(data)
	v, _ := strconv.ParseInt(data[dataLen-4:dataLen-2], 16, 0)
	v = v * -1
	return int(v)
}

func decodeSNR(data string) float32 {
	dataLen := len(data)
	v, _ := strconv.ParseInt(data[dataLen-2:], 16, 0)
	var f float32
	if v > 127 {
		f = float32(v - 256)
	} else {
		f = float32(v) / 4.
	}
	return f
}
