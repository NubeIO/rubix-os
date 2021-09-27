package main

import (
	"fmt"
	unit "github.com/NubeDev/flow-framework/src/units"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"net"
)

func GetInternalIP(nic string) (port string, found bool) {
	itf, _ := net.InterfaceByName(nic)
	item, _ := itf.Addrs()
	var ip net.IP
	for _, addr := range item {
		switch v := addr.(type) {
		case *net.IPNet:
			if !v.IP.IsLoopback() {
				if v.IP.To4() != nil {
					ip = v.IP
				}
			}
		}
	}
	if ip != nil {
		return ip.String(), true
	} else {
		return "na", false
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Info(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	_, res := unit.Process(5000, "meter", "kilometer")
	fmt.Println(res.String())
	fmt.Println(res.AsFloat())

	fmt.Println(utils.RandInt(1, 11))
	fmt.Println(utils.RandInt(1, 11))
	fmt.Println(utils.RandFloat(1, 1011))
	fmt.Println(utils.RandFloat(1, 1011))
	fmt.Println(GetOutboundIP())
	fmt.Println(GetInternalIP("ee"))
}
