package decoder

import (
	"strconv"
)

type TDropletTH struct {
	CommonValues
	Voltage     int     `json:"voltage"`
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
}

type TDropletTHL struct {
	TDropletTH
	Light int `json:"light"`
}

type TDropletTHLM struct {
	TDropletTHL
	Motion int `json:"motion"`
}

func DropletTH(data string, sensor TSensorType) TDropletTH {
	d := Common(data, sensor)
	temperature := dropletTemp(data)
	humidity := dropletHumidity(data)
	voltage := dropletVoltage(data)
	v := TDropletTH{
		CommonValues: d,
		Voltage:      voltage,
		Temperature:  temperature,
		Humidity:     humidity,
	}
	return v
}

func DropletTHL(data string, sensor TSensorType) TDropletTHL {
	d := DropletTH(data, sensor)
	light := dropletLight(data)
	v := TDropletTHL{
		TDropletTH: d,
		Light:      light,
	}
	return v
}

func DropletTHLM(data string, sensor TSensorType) TDropletTHLM {
	d := DropletTHL(data, sensor)
	motion := dropletMotion(data)
	v := TDropletTHLM{
		TDropletTHL: d,
		Motion:      motion,
	}
	return v
}

func dropletTemp(data string) float64 {
	v, _ := strconv.ParseInt(data[10:12]+data[8:10], 16, 0)
	v_ := float64(v) / 100
	return v_
}

func dropletHumidity(data string) int {
	v, _ := strconv.ParseInt(data[16:18], 16, 0)
	v_ := v & 127
	return int(v_)
}

func dropletVoltage(data string) int {
	v, _ := strconv.ParseInt(data[22:24], 16, 0)
	v_ := v / 50
	return int(v_)
}

func dropletLight(data string) int {
	v := data[20:22] + data[18:20]
	v_, _ := strconv.ParseInt(v, 16, 0)
	return int(v_)
}

func dropletMotion(data string) int {
	v_, _ := strconv.ParseInt(data[16:18], 16, 0)
	if v_ > 127 {
		return 1
	} else {
		return 0
	}
}
