package decoder

import (
	"github.com/NubeIO/flow-framework/model"
	"math"
	"strconv"
)

type TMicroEdge struct {
	CommonValues
	Voltage float64 `json:"voltage"`
	Pulse   int     `json:"pulse"`
	AI1     float64 `json:"ai_1"`
	AI2     float64 `json:"ai_2"`
	AI3     float64 `json:"ai_3"`
}

func MicroEdge(data string, sensor TSensorType) TMicroEdge {
	d := Common(data, sensor)
	p := pulse(data)
	a1 := ai1(data)
	a2 := ai2(data)
	a3 := ai3(data)
	vol := voltage(data)
	v := TMicroEdge{
		CommonValues: d,
		Voltage:      vol,
		Pulse:        p,
		AI1:          a1,
		AI2:          a2,
		AI3:          a3,
	}
	return v
}

func pulse(data string) int {
	v, _ := strconv.ParseInt(data[8:16], 16, 0)
	return int(v)
}

func ai1(data string) float64 {
	v, _ := strconv.ParseInt(data[18:22], 16, 0)
	return float64(v)
}

func ai2(data string) float64 {
	v, _ := strconv.ParseInt(data[22:26], 16, 0)
	return float64(v)
}

func ai3(data string) float64 {
	v, _ := strconv.ParseInt(data[26:30], 16, 0)
	return float64(v)
}

func voltage(data string) float64 {
	v, _ := strconv.ParseInt(data[16:18], 16, 0)
	v_ := float64(v) / 50
	return v_
}

func MicroEdgePointType(sensorType string, value float64) float64 {
	switch sensorType {
	case model.IOType.RAW:
		return value
	case model.IOType.Digital:
		if value == 0 || value >= 1000 {
			return 0
		} else {
			return 1
		}
	case model.IOType.Thermistor10K:
		vlt := 3.34
		v := (value / 1024) * vlt
		R0 := 10000.0
		R := (R0 * v) / (vlt - v)
		t0 := 273.0 + 25.0
		b := 3850.0
		var ml float64
		ml = math.Log(R / R0)
		T := 1.0 / (1.0/t0 + (1.0/b)*ml)
		output := T - 273.15
		return output
	case model.IOType.VoltageDC:
		output := (value / 1024) * 10
		return output
	default:
		return value
	}
}
