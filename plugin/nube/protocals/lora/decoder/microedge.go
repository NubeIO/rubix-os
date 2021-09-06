package decoder

import (
	"strconv"
)

type TMicroEdge struct {
	CommonValues
	Voltage int     `json:"voltage"`
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

func voltage(data string) int {
	v, _ := strconv.ParseInt(data[16:18], 16, 0)
	v_ := v / 50
	return int(v_)
}
