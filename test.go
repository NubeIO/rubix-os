package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/src/utilstime"
)

func main() {

	tz, err := utilstime.GetHardwareTZ()
	if err != nil {
		return
	}
	fmt.Println(tz, err)
}
