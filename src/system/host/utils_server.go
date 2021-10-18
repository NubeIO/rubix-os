package host

import (
	"fmt"
	"strings"
)

func getServerInfo(debug bool) string {
	// lsb_release -a
	// head -n 1 /etc/issue | sed 's/\\n//g' | sed 's/\\l//g'
	sh := "lsb_release -d|sed 's/'$(lsb_release -d|awk '{print $1}')'//g'"
	res, err := cmdRun(sh, debug)
	if err != nil {
		return ""
	} else {
		return strings.Trim(string(res), "\n\\n\\l\t\\t")
	}
}

const (
	BCM2708 = "BCM2708"
	BCM2709 = "BCM2709"
	BCM2711 = "BCM2711"
	BCM2835 = "BCM2835"
	BCM2836 = "BCM2836"
)

type IsNubeHardware struct {
	Type            string `json:"type"`
	IsNubeSupported bool   `json:"is_nube_supported"`
	Processor       string `json:"processor"`
}

func IsRaspberryPI() (model string, isPi bool, err error) {
	sh := "cat /proc/cpuinfo | grep \"Hard\""
	sys, err := cmdRun(sh, false)
	if err != nil {
		return "", false, err
	}
	out := string(sys)
	m := "unknown"
	if strings.Contains(out, BCM2708) {
		m = BCM2708
	} else if strings.Contains(out, BCM2709) {
		fmt.Println(BCM2709)
		m = BCM2709
	} else if strings.Contains(out, BCM2711) {
		m = BCM2711
	} else if strings.Contains(out, BCM2835) {
		m = BCM2835
	} else if strings.Contains(out, BCM2836) {
		m = BCM2836
	}
	return m, true, err
}
