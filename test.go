package main

import (
	"github.com/NubeDev/flow-framework/src/system/networking"
)

func main() {
	//fmt.Println(networking.GetValidNetInterfacesForWeb())
	networking.CheckInternetStatus()

	//exec.Command("bash", "-c", `if ping -I eth0 -c 2 google.com; then echo OK; else echo DEAD ;fi`)
	//
}
