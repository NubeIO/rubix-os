package main

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func checkForBooleanType(ObjectType, DataType string) (isTypeBool bool) {
	isTypeBool = false
	if DataType == string(model.TypeDigital) {
		isTypeBool = true
	}
	switch ObjectType {
	case string(model.ObjTypeBinaryInput), string(model.ObjTypeBinaryOutput), string(model.ObjTypeBinaryValue), string(model.ObjTypeReadDiscreteInputs), string(model.ObjBinaryInput), string(model.ObjBinaryOutput), string(model.ObjBinaryValue):
		isTypeBool = true
	}
	return
}

func checkForOutputType(ObjectType string) (isOutput bool) {
	isOutput = false
	switch ObjectType {
	case string(model.ObjTypeAnalogOutput), string(model.ObjTypeAnalogValue), string(model.ObjTypeBinaryOutput), string(model.ObjTypeBinaryValue), string(model.ObjAnalogOutput), string(model.ObjAnalogValue), string(model.ObjBinaryOutput), string(model.ObjBinaryValue):
		isOutput = true
	}
	return
}

func isWriteable(writeMode model.WriteMode, objectType string) bool {
	if isWriteableObjectType(objectType) {
		return true
	} else {
		return false
	}
}

func isWriteableObjectType(objectType string) bool {
	switch objectType {
	case string(model.ObjTypeAnalogOutput), string(model.ObjTypeAnalogValue):
		return true
	case string(model.ObjTypeBinaryOutput), string(model.ObjTypeBinaryValue):
		return true
	case string(model.ObjAnalogOutput), string(model.ObjAnalogValue):
		return true
	case string(model.ObjBinaryOutput), string(model.ObjBinaryValue):
		return true
	}
	return false
}
