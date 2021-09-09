package decoder

import (
	log "github.com/sirupsen/logrus"
	"strconv"
)

type TSensorType string

const (
	ME   TSensorType = "ME"
	TH   TSensorType = "TH"
	THL  TSensorType = "THL"
	THLM TSensorType = "THLM"
	ZHT  TSensorType = "ZHT"
)

type TSensorCode string

const (
	MicroAA       TSensorCode = "AA"
	DropletAB     TSensorCode = "AB"
	DropletB0     TSensorCode = "B0"
	DropletB1     TSensorCode = "B1"
	DropletB2     TSensorCode = "B2"
	ZipHydrotapD1 TSensorCode = "D1"
)

func GetSensorType(data string) TSensorType {
	sensor := data[2:4]
	switch sensor {
	case string(MicroAA):
		return ME
	case string(DropletB0):
		return TH
	case string(DropletB1):
		return THL
	case string(DropletB2):
		return THLM
	case string(ZipHydrotapD1):
		return ZHT
	default:
		return "None"
	}
}

func CheckSensorCode(data string) TSensorCode {
	sensor := data[2:4]

	switch sensor {
	case string(MicroAA):
		return MicroAA
	case string(DropletB2):
		return DropletB2
	case string(ZipHydrotapD1):
		return ZipHydrotapD1
	default:
		return "None"
	}
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
	switch id := CheckSensorCode(data); id {
	case MicroAA:
	case DropletB2:
		return dl == 36 || dl == 32 || dl == 44
	case ZipHydrotapD1:
		return ZHtCheckPayloadLength(data)
	default:
		return false
	}
	return false
}

func DecodePayload(data string) (CommonValues, interface{}) {
	s := GetSensorType(data)
	var payload interface{} = nil
	var common *CommonValues = nil
	switch s {
	case ME:
		payloadFull := MicroEdge(data, ME)
		common = &payloadFull.CommonValues
		payload = payloadFull
	case TH:
		payloadFull := DropletTH(data, TH)
		common = &payloadFull.CommonValues
		payload = payloadFull
	case THL:
		payloadFull := DropletTHL(data, THL)
		common = &payloadFull.CommonValues
		payload = payloadFull
	case THLM:
		payloadFull := DropletTHLM(data, THLM)
		common = &payloadFull.CommonValues
		payload = payloadFull
	case ZHT:
		base, payloadFull := ZipHydrotap(data, ZHT)
		common = &base.CommonValues
		payload = payloadFull
	default:
		log.Printf("ERROR! No decoder for sensor type: %s", s)
		return CommonValues{}, nil
	}
	return *common, payload
}

type CommonValues struct {
	Sensor string `json:"sensor"`
	Id     string `json:"id"`
	Rssi   int    `json:"rssi"`
}

func Common(data string, sensor TSensorType) CommonValues {
	_id := decodeID(data)
	_rssi := rssi(data)
	_v := CommonValues{
		Sensor: string(sensor),
		Id:     _id,
		Rssi:   _rssi,
	}
	return _v
}

func DataLength(data string) int {
	return len(data)
}

func decodeID(data string) string {
	return data[0:8]
}

func rssi(data string) int {
	_len := DataLength(data)
	v, _ := strconv.ParseInt(data[_len-4:_len-2], 16, 0)
	_v := v * -1
	return int(_v)
}
