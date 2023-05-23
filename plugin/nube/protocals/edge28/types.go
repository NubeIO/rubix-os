package main

import (
	"github.com/NubeIO/rubix-os/utils/boolean"
)

var UOTypes = struct {
	// RAW  string
	DIGITAL string
	// PERCENT string
	VOLTSDC string
	// MILLIAMPS  string
}{
	// RAW:  "RAW",
	DIGITAL: "digital",
	// PERCENT: "PERCENT",
	VOLTSDC: "0-10VDC",
	// MILLIAMPS:  "4-20mA",
}

var UITypes = struct {
	RAW             string
	DIGITAL         string
	PERCENT         string
	VOLTSDC         string
	MILLIAMPS       string
	RESISTANCE      string
	THERMISTOR10KT2 string
	// THERMISTOR10KT3  string
	// THERMISTOR20KT1  string
	// THERMISTORPT100  string
	// THERMISTORPT1000 string
}{
	RAW:       "raw",
	DIGITAL:   "digital",
	PERCENT:   "PERCENT",
	VOLTSDC:   "voltage_dc",
	MILLIAMPS: "current",
	// RESISTANCE:       "RESISTANCE",
	THERMISTOR10KT2: "thermistor_10k_type_2",
	// THERMISTOR10KT3:  "THERMISTOR_10K_TYPE3",
	// THERMISTOR20KT1:  "THERMISTOR_20K_TYPE1",
	// THERMISTORPT100:  "THERMISTOR_PT100",
	// THERMISTORPT1000: "THERMISTOR_PT1000",
}

var Rls = []string{"R1", "R2"}
var DOs = []string{"DO1", "DO2", "DO3", "DO4", "DO5"}
var UOs = []string{"UO1", "UO2", "UO3", "UO4", "UO5", "UO6", "UO7"}
var UIs = []string{"UI1", "UI2", "UI3", "UI4", "UI5", "UI6", "UI7"}
var DIs = []string{"DI1", "DI2", "DI3", "DI4", "DI5", "DI6", "DI7"}

type Point struct {
	IoNumber   string // R1
	ObjectType string // binary_output
	IsOutput   *bool
	IsTypeBool *bool
}

var PointsList = struct {
	R1  Point `json:"R1"`
	R2  Point `json:"R2"`
	DO1 Point `json:"DO1"`
	DO2 Point `json:"DO2"`
	DO3 Point `json:"DO3"`
	DO4 Point `json:"DO4"`
	DO5 Point `json:"DO5"`
	UO1 Point `json:"UO1"`
	UO2 Point `json:"UO2"`
	UO3 Point `json:"UO3"`
	UO4 Point `json:"UO4"`
	UO5 Point `json:"UO5"`
	UO6 Point `json:"UO6"`
	UO7 Point `json:"UO7"`
	UI1 Point `json:"UI1"`
	UI2 Point `json:"UI2"`
	UI3 Point `json:"UI3"`
	UI4 Point `json:"UI4"`
	UI5 Point `json:"UI5"`
	UI6 Point `json:"UI6"`
	UI7 Point `json:"UI7"`
	DI1 Point `json:"DI1"`
	DI2 Point `json:"DI2"`
	DI3 Point `json:"DI3"`
	DI4 Point `json:"DI4"`
	DI5 Point `json:"DI5"`
	DI6 Point `json:"DI6"`
	DI7 Point `json:"DI7"`
}{
	R1:  Point{IoNumber: "R1", ObjectType: "binary_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewTrue()},
	R2:  Point{IoNumber: "R2", ObjectType: "binary_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewTrue()},
	DO1: Point{IoNumber: "DO1", ObjectType: "binary_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewTrue()},
	DO2: Point{IoNumber: "DO2", ObjectType: "binary_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewTrue()},
	DO3: Point{IoNumber: "DO3", ObjectType: "binary_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewTrue()},
	DO4: Point{IoNumber: "DO4", ObjectType: "binary_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewTrue()},
	DO5: Point{IoNumber: "DO5", ObjectType: "binary_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewTrue()},
	UO1: Point{IoNumber: "UO1", ObjectType: "analog_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewFalse()},
	UO2: Point{IoNumber: "UO2", ObjectType: "analog_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewFalse()},
	UO3: Point{IoNumber: "UO3", ObjectType: "analog_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewFalse()},
	UO4: Point{IoNumber: "UO4", ObjectType: "analog_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewFalse()},
	UO5: Point{IoNumber: "UO5", ObjectType: "analog_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewFalse()},
	UO6: Point{IoNumber: "UO6", ObjectType: "analog_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewFalse()},
	UO7: Point{IoNumber: "UO7", ObjectType: "analog_output", IsOutput: boolean.NewTrue(), IsTypeBool: boolean.NewFalse()},
	UI1: Point{IoNumber: "UI1", ObjectType: "analog_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewFalse()},
	UI2: Point{IoNumber: "UI2", ObjectType: "analog_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewFalse()},
	UI3: Point{IoNumber: "UI3", ObjectType: "analog_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewFalse()},
	UI4: Point{IoNumber: "UI4", ObjectType: "analog_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewFalse()},
	UI5: Point{IoNumber: "UI5", ObjectType: "analog_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewFalse()},
	UI6: Point{IoNumber: "UI6", ObjectType: "analog_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewFalse()},
	UI7: Point{IoNumber: "UI7", ObjectType: "analog_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewFalse()},
	DI1: Point{IoNumber: "DI1", ObjectType: "binary_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewTrue()},
	DI2: Point{IoNumber: "DI2", ObjectType: "binary_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewTrue()},
	DI3: Point{IoNumber: "DI3", ObjectType: "binary_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewTrue()},
	DI4: Point{IoNumber: "DI4", ObjectType: "binary_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewTrue()},
	DI5: Point{IoNumber: "DI5", ObjectType: "binary_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewTrue()},
	DI6: Point{IoNumber: "DI6", ObjectType: "binary_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewTrue()},
	DI7: Point{IoNumber: "DI7", ObjectType: "binary_input", IsOutput: boolean.NewFalse(), IsTypeBool: boolean.NewTrue()},
}

var pointList = struct {
	R1  string `json:"R1"`
	R2  string `json:"R2"`
	DO1 string `json:"DO1"`
	DO2 string `json:"DO2"`
	DO3 string `json:"DO3"`
	DO4 string `json:"DO4"`
	DO5 string `json:"DO5"`
	UO1 string `json:"UO1"`
	UO2 string `json:"UO2"`
	UO3 string `json:"UO3"`
	UO4 string `json:"UO4"`
	UO5 string `json:"UO5"`
	UO6 string `json:"UO6"`
	UO7 string `json:"UO7"`
	UI1 string `json:"UI1"`
	UI2 string `json:"UI2"`
	UI3 string `json:"UI3"`
	UI4 string `json:"UI4"`
	UI5 string `json:"UI5"`
	UI6 string `json:"UI6"`
	UI7 string `json:"UI7"`
	DI1 string `json:"DI1"`
	DI2 string `json:"DI2"`
	DI3 string `json:"DI3"`
	DI4 string `json:"DI4"`
	DI5 string `json:"DI5"`
	DI6 string `json:"DI6"`
	DI7 string `json:"DI7"`
}{
	R1:  "R1",
	R2:  "R2",
	DO1: "DO1",
	DO2: "DO2",
	DO3: "DO3",
	DO4: "DO4",
	DO5: "DO5",
	UO1: "UO1",
	UO2: "UO2",
	UO3: "UO3",
	UO4: "UO4",
	UO5: "UO5",
	UO6: "UO6",
	UO7: "UO7",
	UI1: "UI1",
	UI2: "UI2",
	UI3: "UI3",
	UI4: "UI4",
	UI5: "UI5",
	UI6: "UI6",
	UI7: "UI7",
	DI1: "DI1",
	DI2: "DI2",
	DI3: "DI3",
	DI4: "DI4",
	DI5: "DI5",
	DI6: "DI6",
	DI7: "DI7",
}

func pointsAll() []string {
	out := append(Rls, DOs...)
	out = append(out, UOs...)
	out = append(out, DIs...)
	out = append(out, UIs...)
	return out
}
