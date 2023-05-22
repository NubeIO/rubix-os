package system

import (
	"errors"
	"github.com/NubeIO/lib-networking/scanner"
)

type Scanner struct {
	Count int      `json:"count"`
	Iface string   `json:"iface"`
	Ip    string   `json:"ip"`
	Ports []string `json:"ports"`
}

func (inst *System) RunScanner(body *Scanner) (*scanner.Hosts, error) {
	if body == nil {
		return nil, errors.New("scanner body can not be empty")
	}
	var count = body.Count
	var iface = body.Iface
	var ip = body.Ip
	var ports = body.Ports

	if count > 254 {
		count = 254
	}
	if count <= 0 {
		count = 254
	}
	if len(ports) == 0 {
		ports = []string{"22", "1414", "1883", "1615", "1616", "502", "1313", "1660", "1661", "1662"}
	}
	scan := scanner.New()
	address, err := scan.ResoleAddress(ip, count, iface)
	if err != nil {
		return nil, err
	}
	host := scan.IPScanner(address, ports, true)
	return host, nil
}
