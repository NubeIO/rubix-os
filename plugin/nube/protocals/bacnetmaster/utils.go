package main

import (
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/writemode"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"strconv"
	"strings"
)

func s2f(v string) (float64, error) {
	return strconv.ParseFloat(v, 64)
}

func s2i(v string) (int, error) {
	return strconv.Atoi(v)
}

func s2iNoErr(v string) int {
	i, _ := strconv.Atoi(v)
	return i
}

func decodeMac(mac string) (string, int) {
	parts := strings.Split(mac, ":")
	if len(parts) >= 2 {
		return parts[0], s2iNoErr(parts[1])
	} else {
		return "", 0
	}
}

func checkForBooleanType(ObjectType, DataType string) (isTypeBool bool) {
	isTypeBool = false
	if DataType == string(model.TypeDigital) {
		isTypeBool = true
	}
	switch ObjectType {
	case string(model.ObjTypeBinaryInput),
		string(model.ObjTypeBinaryOutput),
		string(model.ObjTypeBinaryValue),
		string(model.ObjTypeReadDiscreteInputs),
		string(model.ObjBinaryInput),
		string(model.ObjBinaryOutput),
		string(model.ObjBinaryValue):
		isTypeBool = true
	}
	return
}

func checkForOutputType(ObjectType string) (isOutput bool) {
	isOutput = false
	switch ObjectType {
	case string(model.ObjTypeAnalogOutput),
		string(model.ObjTypeAnalogValue),
		string(model.ObjTypeBinaryOutput),
		string(model.ObjTypeBinaryValue),
		string(model.ObjTypeEnumOutput),
		string(model.ObjTypeEnumValue),
		string(model.ObjAnalogOutput),
		string(model.ObjAnalogValue),
		string(model.ObjEnumOutput), // MSO
		string(model.ObjEnumValue),  // MSV
		string(model.ObjBinaryOutput),
		string(model.ObjBinaryValue):
		isOutput = true
	}
	return
}

func isWriteable(writeMode model.WriteMode, objectType string) bool {
	if isWriteableObjectType(objectType) && writemode.IsWriteable(writeMode) {
		return true
	} else {
		return false
	}
}

func isWriteableObjectType(objectType string) bool {
	switch objectType {
	case string(model.ObjTypeAnalogOutput),
		string(model.ObjTypeAnalogValue):
		return true
	case string(model.ObjTypeBinaryOutput),
		string(model.ObjTypeBinaryValue):
		return true
	case string(model.ObjTypeEnumOutput),
		string(model.ObjTypeEnumValue):
		return true
	case string(model.ObjAnalogOutput),
		string(model.ObjAnalogValue):
		return true
	case string(model.ObjBinaryOutput),
		string(model.ObjBinaryValue):
		return true
	case string(model.ObjEnumOutput),
		string(model.ObjEnumValue):
		return true
	}
	return false
}

type PriArray struct {
	P1  *float64 `json:"_1"`
	P2  *float64 `json:"_2"`
	P3  *float64 `json:"_3"`
	P4  *float64 `json:"_4"`
	P5  *float64 `json:"_5"`
	P6  *float64 `json:"_6"`
	P7  *float64 `json:"_7"`
	P8  *float64 `json:"_8"`
	P9  *float64 `json:"_9"`
	P10 *float64 `json:"_10"`
	P11 *float64 `json:"_11"`
	P12 *float64 `json:"_12"`
	P13 *float64 `json:"_13"`
	P14 *float64 `json:"_14"`
	P15 *float64 `json:"_15"`
	P16 *float64 `json:"_16"`
}

func set(part string) *float64 {
	if part == "Null" {
		return nil
	} else {
		f, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil
		}
		return float.New(f)
	}
}

func cleanArray(payload []string) *PriArray {
	if len(payload) != 16 {
		return nil
	}
	arr := &PriArray{
		P1:  set(payload[0]),
		P2:  set(payload[1]),
		P3:  set(payload[2]),
		P4:  set(payload[3]),
		P5:  set(payload[4]),
		P6:  set(payload[5]),
		P7:  set(payload[6]),
		P8:  set(payload[7]),
		P9:  set(payload[8]),
		P10: set(payload[9]),
		P11: set(payload[10]),
		P12: set(payload[11]),
		P13: set(payload[12]),
		P14: set(payload[13]),
		P15: set(payload[14]),
		P16: set(payload[15]),
	}
	return arr
}
