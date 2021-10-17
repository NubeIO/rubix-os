package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/src/system/command"
)

func main() {
	//fmt.Println(networking.GetValidNetInterfacesForWeb())
	//fmt.Println(system.ProgramUptime())
	a, err := command.Run("if ping -I eth0 -c 2 google.com ; then echo OK ; else echo DEAD ; fi")
	fmt.Println(a, err)
	//
}
