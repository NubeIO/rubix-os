package main

var Rls = []string{"R2", "R2"}
var DOs = []string{"DO1", "DO2", "DO3", "DO4", "DO5"}
var UOs = []string{"UO1", "UO2", "UO3", "UO4", "UO5", "UO6", "UO7"}
var UIs = []string{"UI1", "UI2", "UI3", "UI4", "UI5", "UI6", "UI7"}
var DIs = []string{"DI1", "DI2", "DI3", "DI4", "DI5", "DI6", "DI7"}

var pointList = struct {
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
