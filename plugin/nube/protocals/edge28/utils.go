package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/structs"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube/edge28"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube/thermistor"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"reflect"
	"strconv"
)

// GetGPIOValueForUOByType converts the point value to the correct edge28 UO GPIO value based on the IoType
func GetGPIOValueForUOByType(point *model.Point) (float64, error) {
	var err error
	var result float64
	if !structs.ExistsInStrut(UOTypes, point.IoType) {
		err = errors.New(fmt.Sprintf("skipping %v, IoType %v not recognized.", point.IoNumber, point.IoType))
		return 0, err
	}
	result = plugin.PointWrite(point)

	switch point.IoType {
	case UOTypes.DIGITAL:
		result, err = DigitalToGPIOValue(result, true)
	case UOTypes.PERCENT:
		result = edge28.PercentToGPIOValue(result)
	case UOTypes.VOLTSDC:
		result = edge28.VoltageToGPIOValue(result)
	default:
		err = errors.New("UO IoType is not a recognized type")
	}
	if err != nil {
		return 0, err
	} else {
		return result, nil
	}
}

// GetValueFromGPIOForUIByType converts the GPIO value to the scaled UI value based on the IoType
func GetValueFromGPIOForUIByType(point *model.Point, value float64) (float64, error) {
	var err error
	var result float64

	if !structs.ExistsInStrut(UITypes, point.IoType) {
		err = errors.New(fmt.Sprintf("skipping %v, IoType %v not recognized.", point.IoNumber, point.IoType))
		return 0, err
	}
	switch point.IoType {
	case UITypes.RAW:
		result = value
	case UITypes.DIGITAL:
		result = edge28.GPIOValueToDigital(value)
	case UITypes.PERCENT:
		result = edge28.GPIOValueToPercent(value)
	case UITypes.VOLTSDC:
		result = edge28.GPIOValueToVoltage(value)
	case UITypes.MILLIAMPS:
		result = edge28.ScaleGPIOValueTo420ma(value)
	case UITypes.RESISTANCE:
		result = edge28.ScaleGPIOValueToResistance(value)
	case UITypes.THERMISTOR10KT2:
		resistance := edge28.ScaleGPIOValueToResistance(value)
		result, err = thermistor.ResistanceToTemperature(resistance, thermistor.T210K)
	case UITypes.THERMISTOR10KT3:
		resistance := edge28.ScaleGPIOValueToResistance(value)
		result, err = thermistor.ResistanceToTemperature(resistance, thermistor.T310K)
	case UITypes.THERMISTOR20KT1:
		resistance := edge28.ScaleGPIOValueToResistance(value)
		result, err = thermistor.ResistanceToTemperature(resistance, thermistor.T120K)
	case UITypes.THERMISTORPT100:
		resistance := edge28.ScaleGPIOValueToResistance(value)
		result, err = thermistor.ResistanceToTemperature(resistance, thermistor.PT100)
	case UITypes.THERMISTORPT1000:
		resistance := edge28.ScaleGPIOValueToResistance(value)
		result, err = thermistor.ResistanceToTemperature(resistance, thermistor.PT1000)
	default:
		err = errors.New("UI IoType is not a recognized type")
		return 0, err
	}
	return result, nil
}

// DigitalToGPIOValue converts true/false values (all basic types allowed) to BBB GPIO 0/1 ON/OFF (FOR DOs/Relays) and to 100/0 (FOR UOs).  Note that the GPIO value for digital points is inverted.
func DigitalToGPIOValue(input interface{}, isUO bool) (float64, error) {
	var inputAsBool bool
	var err error = nil
	switch input.(type) {
	case string:
		inputAsBool, err = strconv.ParseBool(reflect.ValueOf(input).String())
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		inputAsBool = reflect.ValueOf(input).Int() != 0
	case float32, float64:
		inputAsBool = reflect.ValueOf(input).Float() != float64(0)
	case bool:
		inputAsBool = reflect.ValueOf(input).Bool()
	default:
		err = errors.New("edge28-polling: input is not a recognized type")
	}
	if err != nil {
		return 0, err
	} else if inputAsBool {
		if isUO {
			return 0, nil // 0 is the 12vdc/ON GPIO value for UOs
		} else {
			return 1, nil // 1 is the 12vdc/ON GPIO value for DOs/Relays
		}
	} else {
		if isUO {
			return 100, nil // 100 is the 0vdc/OF GPIO value for UOs
		} else {
			return 0, nil // 0 is the 0vdc/OFF GPIO value for DOs/Relays
		}
	}
}

func selectObjectType(IoNumber string) (objectType string, isOutput bool) {
	isOutput = false
	switch IoNumber {
	case PointsList.R1.IoNumber, PointsList.R2.IoNumber:
		objectType = PointsList.R1.ObjectType
		isOutput = true
	case PointsList.UO1.IoNumber, PointsList.UO2.IoNumber, PointsList.UO3.IoNumber, PointsList.UO4.IoNumber, PointsList.UO5.IoNumber, PointsList.UO6.IoNumber, PointsList.UO7.IoNumber:
		objectType = PointsList.UO1.ObjectType
		isOutput = true
	case PointsList.DO1.IoNumber, PointsList.DO2.IoNumber, PointsList.DO3.IoNumber, PointsList.DO4.IoNumber, PointsList.DO5.IoNumber:
		objectType = PointsList.DO1.ObjectType
		isOutput = true
	case PointsList.UI1.IoNumber, PointsList.UI2.IoNumber, PointsList.UI3.IoNumber, PointsList.UI4.IoNumber, PointsList.UI5.IoNumber, PointsList.UI6.IoNumber, PointsList.UI7.IoNumber:
		objectType = PointsList.UI1.ObjectType
	case PointsList.DI1.IoNumber, PointsList.DI2.IoNumber, PointsList.DI3.IoNumber, PointsList.DI4.IoNumber, PointsList.DI5.IoNumber, PointsList.DI6.IoNumber, PointsList.DI7.IoNumber:
		objectType = PointsList.DI1.ObjectType
	}
	return
}

func checkForBooleanType(ioType string) (isTypeBool bool) {
	isTypeBool = false
	switch ioType {
	case UOTypes.DIGITAL, UITypes.DIGITAL:
		isTypeBool = true
	}
	return
}

func limitValueByEdge28Type(ioType string, inputVal *float64) (outputVal *float64) {
	if inputVal == nil {
		return nil
	}
	inputValFloat := *inputVal
	switch ioType {
	case UOTypes.DIGITAL, UITypes.DIGITAL:
		if inputValFloat <= 0 {
			outputVal = float.New(0)
		} else {
			outputVal = float.New(1)
		}
	case UOTypes.VOLTSDC, UITypes.VOLTSDC:
		if inputValFloat <= 0 {
			outputVal = float.New(0)
		} else if inputValFloat >= 10 {
			outputVal = float.New(10)
		} else {
			outputVal = inputVal
		}
	case UOTypes.PERCENT, UITypes.PERCENT:
		if inputValFloat <= 0 {
			outputVal = float.New(0)
		} else if inputValFloat >= 100 {
			outputVal = float.New(100)
		} else {
			outputVal = inputVal
		}
	default:
		outputVal = inputVal
	}
	return outputVal
}

func limitPriorityArrayByEdge28Type(ioType string, priority *model.PointWriter) *map[string]*float64 {
	priorityMap := map[string]*float64{}

	for key, val := range *priority.Priority {

		var outputVal *float64
		if val == nil {
			outputVal = nil
			continue
		}
		fmt.Println(`limitPriorityArrayByEdge28Type(): `, key, *val)
		switch ioType {
		case UOTypes.DIGITAL, UITypes.DIGITAL:
			if *val <= 0 {
				outputVal = float.New(0)
			} else {
				outputVal = float.New(1)
			}

		case UOTypes.VOLTSDC, UITypes.VOLTSDC:
			if *val <= 0 {
				outputVal = float.New(0)
			} else if *val >= 10 {
				outputVal = float.New(10)
			} else {
				outputVal = float.New(*val)
			}

		case UOTypes.PERCENT, UITypes.PERCENT:
			if *val <= 0 {
				outputVal = float.New(0)
			} else if *val >= 100 {
				outputVal = float.New(100)
			} else {
				outputVal = float.New(*val)
			}

		default:
			outputVal = float.New(*val)
		}
		fmt.Println(`limitPriorityArrayByEdge28Type(): `, key, *outputVal)
		priorityMap[key] = outputVal
	}
	return &priorityMap
}
